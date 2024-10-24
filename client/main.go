package main

import (
	"dummy-endpoints-http/requester"
	"dummy-endpoints-http/structs"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

func main() {
	beginTime := time.Now()
	beginPort := structs.GetPorts().Min
	endPort := structs.GetPorts().Max

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allResponses []structs.Response // Slice to hold structs.Response objects

	// Iterate through each port and send requests concurrently
	for port := beginPort; port <= endPort; port++ {
		wg.Add(1)

		go func(port int) {
			defer wg.Done()

			url := fmt.Sprintf("http://localhost:%d", port)
			response, err := requester.MakeWG(url) // Create an empty structs.Response object and unmarshal
			if err != nil {
				log.Printf("Error for port %d: %v", port, err)
				return
			}

			// Lock before modifying the shared slice
			mu.Lock()
			allResponses = append(allResponses, *response)
			mu.Unlock()
		}(port)
	}

	// Wait for all Go routines to complete
	wg.Wait()

	// Group responses by address for price comparison
	groupedByAddress := make(map[string][]structs.ResponseRow)

	for _, response := range allResponses {
		for _, resp := range response.Responses { // Access nested responses
			groupedByAddress[resp.Address] = append(groupedByAddress[resp.Address], resp)
		}
	}

	// Compare prices by address and check for threshold
	priceDifferenceThreshold := structs.PriceDifferencePct / 100.0

	for address, responseRows := range groupedByAddress {
		if len(responseRows) < 2 {
			continue // Not enough responses
		}

		// Sort response rows by price within the group
		sort.Slice(responseRows, func(i, j int) bool {
			return responseRows[i].Price < responseRows[j].Price
		})

		// Compare prices across all DEXes for this address
		for i := 0; i < len(responseRows); i++ {
			for j := i + 1; j < len(responseRows); j++ {
				price1 := float64(responseRows[i].Price)
				price2 := float64(responseRows[j].Price)
				diffPct := math.Abs(price2-price1) / price1

				if diffPct > priceDifferenceThreshold {
					// Log the price difference with DEX names
					fmt.Printf("Price difference found for Address: %s\n", address)
					var fromAddress, toAddress string
					var lowPrice, highPrice int
					if responseRows[i].Price < responseRows[j].Price {
						fromAddress = responseRows[i].Address
						toAddress = responseRows[j].Address
						lowPrice = responseRows[i].Price
						highPrice = responseRows[j].Price
					} else {
						fromAddress = responseRows[j].Address
						toAddress = responseRows[i].Address
						lowPrice = responseRows[j].Price
						highPrice = responseRows[i].Price
					}
					//Need to print the address of the DEXes, from Response.Dex <---------
					fmt.Printf("Lowest Price: %s (Address: %s) ---> Highest Price: %s (Address: %s)\n", lowPrice, fromAddress, highPrice, toAddress)
					fmt.Printf("Price1: %.2f, Price2: %.2f, Difference: %.2f%%\n", price1, price2, diffPct*100)
				}
			}
		}
	}

	fmt.Println("Total time taken: ", time.Since(beginTime))
	fmt.Println("Total number of responses: ", len(allResponses))
	fmt.Println("Total number of ports(servers): ", endPort-beginPort+1)
}
