package main

import (
	"fmt"
	"time"
)

// APIError represents an error related to the API request
type APIError struct {
	Message string
}

// Error returns the error message for APIError
func (e *APIError) Error() string {
	return e.Message
}

// CurrencyError represents an error related to currency conversion
type CurrencyError struct {
	Message string
}

// Error returns the error message for CurrencyError
func (e *CurrencyError) Error() string {
	return e.Message
}

// retryRequest retries a function with exponential backoff
func retryRequest(attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		time.Sleep(sleep)
		sleep *= 2
	}

	return &APIError{Message: "Max retries reached"}
}

// userFriendlyError provides user-friendly error messages
func userFriendlyError(err error) string {
	switch e := err.(type) {
	case *APIError:
		return "There was an issue with the currency conversion service. Please try again later."
	case *CurrencyError:
		return fmt.Sprintf("Currency error: %s", e.Message)
	default:
		return "An unexpected error occurred. Please try again."
	}
}

// localCache stores exchange rates locally
var localCache = map[string]float64{
	"USD": 1.0,
	"EUR": 0.85,
	"GBP": 0.75,
}

// handleAPIUnavailability handles API unavailability by using local cache
func handleAPIUnavailability() map[string]float64 {
	return localCache
}
