package main

import (
	"dummy-endpoints-http/requester"
	"dummy-endpoints-http/structs"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

func main() {
	// Initialize credential manager
	cm := NewCredentialManager()

	// Get port range
	beginPort := structs.GetPorts().Min
	endPort := structs.GetPorts().Max

	// Add credentials for all ports
	for port := beginPort; port <= endPort; port++ {
		url := fmt.Sprintf("http://localhost:%d", port)
		cm.AddCredential(APICredential{
			EndpointURL: url,
			Name:        fmt.Sprintf("Port%d", port),
		})
	}

	for {
		beginTime := time.Now()

		// Channel to collect responses
		results := make(chan APIResponse, endPort-beginPort+1)
		var wg sync.WaitGroup
		var allResponses []structs.Response

		// Launch concurrent requests
		for port := beginPort; port <= endPort; port++ {
			wg.Add(1)
			cred, _ := cm.GetCredential(fmt.Sprintf("Port%d", port))

			go func(c APICredential) {
				defer wg.Done()

				response, err := requester.Make(c.EndpointURL)
				if err != nil {
					results <- APIResponse{
						EndpointName: c.Name,
						Error:        err,
					}
					return
				}

				results <- APIResponse{
					EndpointName: c.Name,
					Data:         *response,
				}
			}(cred)
		}

		// Close results channel when all goroutines are done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Process results
		for result := range results {
			if result.Error != nil {
				log.Printf("Error for %s: %v", result.EndpointName, result.Error)
				continue
			}

			if resp, ok := result.Data.(structs.Response); ok {
				allResponses = append(allResponses, resp)
			}
		}

		// Existing price comparison logic
		groupedByAddress := make(map[string]map[string]float64)
		for _, response := range allResponses {
			for _, token := range response.Tokens {
				if _, exists := groupedByAddress[token.Address]; !exists {
					groupedByAddress[token.Address] = make(map[string]float64)
				}
				groupedByAddress[token.Address][response.Dex] = token.Price
			}
		}

		priceDifferenceThreshold := structs.PriceDifferencePct / 100.0
		for address, dexPrices := range groupedByAddress {
			dexList := []string{}
			prices := []float64{}

			for dex, price := range dexPrices {
				dexList = append(dexList, dex)
				prices = append(prices, price)
			}

			if len(prices) < 2 {
				continue
			}

			for i := 0; i < len(prices); i++ {
				for j := i + 1; j < len(prices); j++ {
					price1, price2 := prices[i], prices[j]
					diffPct := math.Abs(price2-price1) / price1

					if diffPct > priceDifferenceThreshold {
						var fromDex, toDex string
						if price1 < price2 {
							fromDex = dexList[i]
							toDex = dexList[j]
						} else {
							fromDex = dexList[j]
							toDex = dexList[i]
						}
						swap(address, fromDex, toDex, price1, price2, diffPct)
					}
				}
			}
		}

		fmt.Println("Total time taken: ", time.Since(beginTime), "@", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("Total number of responses: ", len(allResponses))
		fmt.Println("Total number of ports (servers): ", endPort-beginPort+1)
		fmt.Println("----- Sleeping for ", structs.RequestSleepInterval, " seconds -----")
		time.Sleep(structs.RequestSleepInterval * time.Second)
	}
}

// Add the supporting types from the previous answer
type APICredential struct {
	EndpointURL string
	APIKey      string
	APIToken    string
	Name        string
}

type CredentialManager struct {
	credentials map[string]APICredential
	mu          sync.RWMutex
}

type APIResponse struct {
	EndpointName string
	Data         interface{}
	Error        error
}

func NewCredentialManager() *CredentialManager {
	return &CredentialManager{
		credentials: make(map[string]APICredential),
	}
}

func (cm *CredentialManager) AddCredential(cred APICredential) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.credentials[cred.Name] = cred
}

func (cm *CredentialManager) GetCredential(name string) (APICredential, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cred, exists := cm.credentials[name]
	return cred, exists
}

func swap(address, fromDex, toDex string, price1, price2, diffPct float64) {
	fmt.Printf("Price difference found for Address: %s\n", address)
	fmt.Printf("Lowest Price: %s ---> Highest Price: %s\n", fromDex, toDex)
	fmt.Printf("Price1: %.2f, Price2: %.2f, Difference: %.2f%%\n",
		price1, price2, diffPct*100)
}
