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
				response, err := requester.Make(url) // Assuming Make returns a *structs.Response
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

		// Wait for all goroutines to complete
		wg.Wait()

		// Map to group tokens by address
		groupedByAddress := make(map[string]map[string]float64) // map[address]map[DEX]Price

		for _, response := range allResponses {
			for _, token := range response.Tokens {
				if _, exists := groupedByAddress[token.Address]; !exists {
					groupedByAddress[token.Address] = make(map[string]float64)
				}
				groupedByAddress[token.Address][response.Dex] = token.Price
			}
		}

		// Compare prices by address
		priceDifferenceThreshold := structs.PriceDifferencePct / 100.0

		for address, dexPrices := range groupedByAddress {
			dexList := []string{}
			prices := []float64{}

			// Convert map to slice for comparison
			for dex, price := range dexPrices {
				dexList = append(dexList, dex)
				prices = append(prices, price)
			}

			if len(prices) < 2 {
				continue // Not enough DEXes to compare
			}

			// Compare prices across all DEXes for this token
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

func swap(address, fromDex, toDex string, price1, price2, diffPct float64) {
	fmt.Printf("Price difference found for Address: %s\n", address)
	fmt.Printf("Lowest Price: %s ---> Highest Price: %s\n", fromDex, toDex)
	fmt.Printf("Price1: %.2f, Price2: %.2f, Difference: %.2f%%\n",
		price1, price2, diffPct*100)
}
