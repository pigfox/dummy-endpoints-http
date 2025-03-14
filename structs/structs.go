package structs

import "math/rand"

const TokensPerServerMax = 4
const TokensPerServerMin = 1 // Simulating returned number of tokens by DEX
const PriceDifferencePct = 5
const RequestTimeOut = 5000   // Timeout in milliseconds
const ResponseDelayMax = 5000 //
const MinPort = 10001
const MaxPort = 10003
const RequestSleepInterval = 2

// var FailedPorts = []int{10002, 10003, 10006, 10011, 11012} // simulate failed ports
var FailedPorts = []int{}

type Swap struct {
}

type Ports struct {
	Min    int   `json:"min"`
	Max    int   `json:"max"`
	Failed []int `json:"excluded"`
}

type Response struct {
	Dex    string `json:"dex"`
	Tokens []Token
}

type Token struct {
	Timestamp string  `json:"timestamp"`
	Price     float64 `json:"price"`
	Supply    int     `json:"supply"`
	Address   string  `json:"address"`
}

func GetPorts() Ports {
	return Ports{
		Min:    MinPort,
		Max:    MaxPort,
		Failed: FailedPorts, // []int{10002, 10003, 10010},
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

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
