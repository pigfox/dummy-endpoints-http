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
	beginPort := structs.GetPorts().Min
	endPort := structs.GetPorts().Max

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allResponses []structs.Response // Slice to hold structs.Response objects
	for {
		beginTime := time.Now()
		// Iterate through each port and send requests concurrently
		for port := beginPort; port <= endPort; port++ {
			wg.Add(1)

			go func(port int) {
				defer wg.Done()

				url := fmt.Sprintf("http://localhost:%d", port)
				response, err := requester.MakeWG(url) // Assuming MakeWG returns a *structs.Response
				if err != nil {
					log.Printf("Error for port %d: %v", port, err)
					return
				}

				// Lock before modifying the shared slice
				mu.Lock()
				allResponses = append(allResponses, *response) // Append the dereferenced response
				mu.Unlock()
			}(port)
		}

		// Wait for all Go routines to complete
		wg.Wait()

		// Group responses by address for price comparison
		groupedByAddress := make(map[string][]structs.Response)

		for _, response := range allResponses {
			for _, resp := range response.Tokens { // Access nested responses
				groupedByAddress[resp.Address] = append(groupedByAddress[resp.Address], response)
			}
		}

		// Compare prices by address and check for threshold
		priceDifferenceThreshold := structs.PriceDifferencePct / 100.0

		for address, responses := range groupedByAddress {
			if len(responses) < 2 {
				continue // Not enough responses to compare
			}

			// Compare prices across all DEXes for this address
			for i := 0; i < len(responses); i++ {
				for j := i + 1; j < len(responses); j++ {
					price1 := float64(responses[i].Tokens[0].Price)
					price2 := float64(responses[j].Tokens[0].Price)
					diffPct := math.Abs(price2-price1) / price1

					if diffPct > priceDifferenceThreshold && responses[i].Dex != responses[j].Dex {
						if responses[i].Tokens[0].Address == responses[j].Tokens[0].Address {
							// Log the price difference with DEX names
							var fromDex, toDex string
							if responses[i].Tokens[0].Price < responses[j].Tokens[0].Price {
								fromDex = responses[i].Dex
								toDex = responses[j].Dex
							} else {
								fromDex = responses[j].Dex
								toDex = responses[i].Dex
							}

							swap(address, fromDex, toDex, price1, price2, diffPct)
						}
					}
				}
			}
		}

		fmt.Println("Total time taken: ", time.Since(beginTime))
		fmt.Println("Total number of responses: ", len(allResponses))
		fmt.Println("Total number of ports (servers): ", endPort-beginPort+1)
		fmt.Println("----- Sleeping for ", structs.RequestSleepInterval, " seconds -----")
		time.Sleep(structs.RequestSleepInterval * time.Second)
	}
}

func swap(address, fromDex, toDex string, price1, price2, diffPct float64) {
	fmt.Printf("Price difference found for Address: %s\n", address)
	fmt.Printf("Lowest Price: %s ---> Highest Price: %s\n", fromDex, toDex)
	fmt.Printf("Price1: %.2f, Price2: %.2f, Difference: %.2f%%\n",
		price1, price2, diffPct*100)
}
