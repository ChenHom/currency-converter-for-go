package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFetchExchangeRates(t *testing.T) {
	// Create a test server to mock the API response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rates":{"USD":1.0,"EUR":0.85},"base":"USD","date":"2021-01-01"}`))
	}))
	defer ts.Close()

	// Set the environment variables for the test
	os.Setenv("API_KEY", "test_api_key")
	os.Setenv("BASE_URL", ts.URL)

	// Fetch the exchange rates
	exchangeRate, err := fetchExchangeRates(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the exchange rates
	if exchangeRate.Rates["USD"] != 1.0 {
		t.Errorf("Expected USD rate to be 1.0, got %v", exchangeRate.Rates["USD"])
	}
	if exchangeRate.Rates["EUR"] != 0.85 {
		t.Errorf("Expected EUR rate to be 0.85, got %v", exchangeRate.Rates["EUR"])
	}
}

func TestConvertCurrency(t *testing.T) {
	rates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
	}

	amount := 100.0
	from := "USD"
	to := "EUR"
	expected := 85.0

	convertedAmount, err := convertCurrency(amount, from, to, rates)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if convertedAmount != expected {
		t.Errorf("Expected %v, got %v", expected, convertedAmount)
	}
}

func TestIntegration(t *testing.T) {
	// Create a test server to mock the API response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rates":{"USD":1.0,"EUR":0.85},"base":"USD","date":"2021-01-01"}`))
	}))
	defer ts.Close()

	// Set the environment variables for the test
	os.Setenv("API_KEY", "test_api_key")
	os.Setenv("BASE_URL", ts.URL)

	// Fetch the exchange rates
	exchangeRate, err := fetchExchangeRates(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Convert currency
	amount := 100.0
	from := "USD"
	to := "EUR"
	expected := 85.0

	convertedAmount, err := convertCurrency(amount, from, to, exchangeRate.Rates)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if convertedAmount != expected {
		t.Errorf("Expected %v, got %v", expected, convertedAmount)
	}
}
