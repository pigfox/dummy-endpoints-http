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
		w.Header().Set("Content-Type", "application/json")
		var out []structs.Response
		for i := 0; i < structs.ResponseRowsPerServer; i++ {
			res := structs.Response{
				Message:   fmt.Sprintf("This is port: %d", port),
				TimeStamp: time.Now().Format(time.RFC3339),
				Price:     structs.RandomInt(1, 100),
				Address:   "0x" + fmt.Sprintf("%d", i),
			}
			out = append(out, res)
		}

		if err := json.NewEncoder(w).Encode(out); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
