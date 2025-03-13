package main

import (
	"dummy-endpoints-http/structs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// handler function that returns the current port
func portHandler(port int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simulate random delay in response
		delay := structs.RandomInt(0, structs.ResponseDelayMax)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		var row []structs.Token
		responseRows := structs.RandomInt(structs.TokensPerServerMin, structs.TokensPerServerMax)

		for i := 0; i < responseRows; i++ {
			res := structs.Token{
				Timestamp: time.Now().Format(time.RFC3339),
				Price:     structs.RandomFloat(1, 100),
				Supply:    structs.RandomInt(1000, 100000000),
				Address:   "0x" + fmt.Sprintf("%d", i),
			}
			row = append(row, res)
		}

		response := structs.Response{
			Dex:    fmt.Sprintf("DEX %d", port),
			Tokens: row,
		}

		// Encode the response before setting headers
		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		// Set headers and write response
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseData)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	beginPort := structs.GetPorts().Min
	endPort := structs.GetPorts().Max
	fmt.Println("Total number of ports(servers): ", endPort-beginPort+1)

	if beginPort > endPort {
		log.Fatalf("Begin port should be less than or equal to end port")
	}

	for port := beginPort; port <= endPort; port++ {
		if !structs.Contains(structs.GetPorts().Failed, port) {
			go func(p int) {
				// Create a new mux for each server
				mux := http.NewServeMux()
				mux.HandleFunc("/", portHandler(p))

				addr := fmt.Sprintf(":%d", p)
				log.Printf("Starting server on port %d", p)

				// Start server on the specified port
				if err := http.ListenAndServe(addr, mux); err != nil {
					log.Fatalf("Failed to start server on port %d: %v", p, err)
				}
			}(port)
		}
	}

	// Block main goroutine so servers can continue running
	select {}
}
