package enum

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/awareseven/mobilesniper/pkg/models"
)

func DiscoverNetworkFunctions(target models.Target, openapiPath string, nfrChan chan<- models.NetworkFunctionResult, wg *sync.WaitGroup, maxConcurrency int, verbose bool) error {

	semaphore := make(chan struct{}, maxConcurrency)

	if verbose {
		log.Printf("Enumerating NF definitions on %s:%d", target.IP, target.Port)
	}

	err := filepath.Walk(openapiPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		wg.Add(1)
		semaphore <- struct{}{} // add to channel

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }() // remove from channel

			openapi, err := models.ValidateOpenAPIFile(path)
			if err != nil {
				if verbose {
					log.Printf("Invalid OpenAPI file: %s - %v\n", path, err)
				}
				return // continue processing other files
			}

			var reachableCount int
			var totalExpectedResponses int
			client := &http.Client{}

			// Iterate over all paths defined in the OpenAPI specification
			for path, methods := range openapi.Paths {
				// Iterate over all HTTP methods (GET, POST, etc.) for each path
				for method, operationInterface := range methods {
					operation, ok := operationInterface.(map[string]interface{})
					if !ok {
						continue // Skip if operation cannot be asserted to the expected type
					}

					// Construct the request URL for the given path and method
					url := fmt.Sprintf("http://%s:%d%s", target.IP, target.Port, path)
					req, err := http.NewRequest(strings.ToUpper(method), url, nil)
					if err != nil {
						log.Printf("Error creating request for %s %s: %v", method, path, err)
						continue // Skip to the next method if there's an error creating the request
					}

					// Execute the request
					resp, err := client.Do(req)
					if err != nil {
						log.Printf("Error executing request for %s %s: %   v", method, path, err)
						continue // Skip to the next method if there's an error executing the request
					}
					defer resp.Body.Close()

					// Check if the returned status code is one of the expected codes
					_, ok = operation["responses"].(map[string]interface{})[fmt.Sprintf("%d", resp.StatusCode)]

					if ok {
						// Increment the reachable count if the response matches one of the expected status codes
						reachableCount++

						if verbose {
							log.Printf("Matched expected status code %d for %s %s", resp.StatusCode, method, path)
						}
					}

					// Increment the total expected responses count
					totalExpectedResponses++
				}
			}

			// Calculate accuracy based on the number of expected responses
			accuracy := (float64(reachableCount) / float64(totalExpectedResponses)) * 100

			apiName := openapi.Info.Title

			result := models.NetworkFunctionResult{
				Target:          target,
				NetworkFunction: apiName,
				Accuracy:        accuracy,
			}

			if verbose {
				log.Printf("%f%% NF %s on %s:%d",
					result.Accuracy, result.NetworkFunction,
					result.Target.IP, result.Target.Port,
				)
			}

			nfrChan <- result
		}()

		return nil // Return nil to continue processing other files
	})

	return err
}
