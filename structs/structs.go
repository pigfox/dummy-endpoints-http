package structs

import "math/rand"

const ResponseRowsPerServer = 2000 // Simulating returned number of tokens by DEX
const PriceDifferencePct = 5
const RequestTimeOut = 5000 // Timeout in milliseconds
const MinPort = 10001
const MaxPort = 10003 //12000

type Ports struct {
	Min    int   `json:"min"`
	Max    int   `json:"max"`
	Failed []int `json:"excluded"`
}

type Response struct {
	Message   string `json:"message"`
	TimeStamp string `json:"timestamp"`
	Price     int    `json:"random_int"`
	Address   string `json:"address"`
}

func GetPorts() Ports {
	return Ports{
		Min:    MinPort,
		Max:    MaxPort,
		Failed: []int{}, //Simulating failed ports 10002, 10003, 10010
	}
}

func Contains(arr []int, num int) bool {
	for _, value := range arr {
		if value == num {
			return true
		}
	}
	return false
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
