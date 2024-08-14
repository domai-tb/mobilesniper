package enum

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/awareseven/mobilesniper/pkg/models"
)

func DiscoverNetworkFunctions(target models.Target, openapiPath string, nfrChan chan<- models.NetworkFunctionResult, wg *sync.WaitGroup, maxConcurrency int) error {

	semaphore := make(chan struct{}, maxConcurrency)

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
				// log.Printf("Invalid OpenAPI file: %s - %v\n", path, err)
				return // continue processing other files
			}

			var reachableCount int
			totalEndpoints := len(openapi.Paths)
			client := &http.Client{}

			for path := range openapi.Paths {
				url := fmt.Sprintf("http://%s:%d%s", target.IP, target.Port, path)
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					continue // skip request on error
				}

				resp, err := client.Do(req)
				if err == nil && resp.StatusCode == http.StatusOK {
					reachableCount++
				}
			}

			accuracy := (float64(reachableCount) / float64(totalEndpoints)) * 100
			apiName := openapi.Info.Title

			result := models.NetworkFunctionResult{
				Target:          target,
				NetworkFunction: apiName,
				Accuracy:        accuracy,
			}

			nfrChan <- result
		}()

		return err
	})

	return err
}
