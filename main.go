package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type ExchangeRate struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

var (
	cache          = make(map[string]float64)
	cacheMutex     sync.Mutex
	cacheTimestamp time.Time
	rateLimit      = 10
	rateLimitMutex sync.Mutex
	requestCount   = 0
)

// fetchExchangeRates fetches exchange rates from a third-party API using HTTP and JSON
func fetchExchangeRates(apiURL string) (*ExchangeRate, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var exchangeRate ExchangeRate
	if err := json.Unmarshal(body, &exchangeRate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &exchangeRate, nil
}

// convertCurrency converts an amount from one currency to another using the fetched exchange rates
func convertCurrency(amount float64, from, to string, rates map[string]float64) (float64, error) {
	fromRate, ok := rates[from]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", from)
	}

	toRate, ok := rates[to]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", to)
	}

	return amount * (toRate / fromRate), nil
}

// getCachedRates returns the cached exchange rates if they are still valid
func getCachedRates() map[string]float64 {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(cacheTimestamp) < 15*time.Minute {
		return cache
	}

	return nil
}

// updateCache updates the local cache with the latest exchange rates
func updateCache(rates map[string]float64) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cache = rates
	cacheTimestamp = time.Now()
}

// rateLimitExceeded checks if the rate limit for API requests has been exceeded
func rateLimitExceeded() bool {
	rateLimitMutex.Lock()
	defer rateLimitMutex.Unlock()

	if requestCount >= rateLimit {
		return true
	}

	requestCount++
	return false
}

// resetRateLimit resets the rate limit counter every minute
func resetRateLimit() {
	for {
		time.Sleep(1 * time.Minute)
		rateLimitMutex.Lock()
		requestCount = 0
		rateLimitMutex.Unlock()
	}
}

func main() {
	go resetRateLimit()

	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("BASE_URL")
	if apiKey == "" || baseURL == "" {
		log.Fatal("API_KEY and BASE_URL must be set")
	}

	apiURL := fmt.Sprintf("%s/latest?access_key=%s", baseURL, apiKey)

	var exchangeRate *ExchangeRate
	var err error

	for i := 0; i < 3; i++ {
		if rateLimitExceeded() {
			log.Println("Rate limit exceeded, using cached rates")
			exchangeRate = &ExchangeRate{Rates: getCachedRates()}
			break
		}

		exchangeRate, err = fetchExchangeRates(apiURL)
		if err == nil {
			updateCache(exchangeRate.Rates)
			break
		}

		time.Sleep(time.Duration(i*i) * time.Second)
	}

	if err != nil {
		log.Fatalf("Error fetching exchange rates: %v", err)
	}

	amount := 100.0
	from := "USD"
	to := "EUR"
	convertedAmount, err := convertCurrency(amount, from, to, exchangeRate.Rates)
	if err != nil {
		log.Fatalf("Error converting currency: %v", err)
	}

	fmt.Printf("%.2f %s is %.2f %s\n", amount, from, convertedAmount, to)
}
