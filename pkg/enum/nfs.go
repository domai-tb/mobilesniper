package enum

import (
	"fmt"
	"net/http"

	"github.com/awareseven/mobilesniper/pkg/models"
)

func TestOfNetworkFunction(ip string, port int, openapi *models.OpenAPI, nfrChan chan<- models.NetworkFunctionResult) (*models.NetworkFunctionResult, error) {
	var reachableCount int
	totalEndpoints := len(openapi.Paths)
	client := &http.Client{}

	for path := range openapi.Paths {
		url := fmt.Sprintf("http://%s:%d%s", ip, port, path)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			reachableCount++
		}
	}

	accuracy := (float64(reachableCount) / float64(totalEndpoints)) * 100
	apiName := "Unknown API"

	result := models.NetworkFunctionResult{
		IP:              ip,
		Port:            port,
		NetworkFunction: apiName,
		Accuracy:        accuracy,
	}
	nfrChan <- result

	return &result, nil
}
