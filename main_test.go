package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"currency-converter/currency"
)

func TestFetchExchangeRates(t *testing.T) {
	// 建立一個測試伺服器來模擬 API 回應
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"conversion_rates":{"USD":1.0,"EUR":0.85},"base_code":"USD"}`))
	}))
	defer ts.Close()

	// 設置測試的環境變數
	os.Setenv("API_KEY", "test_api_key")
	os.Setenv("BASE_URL", ts.URL)

	// 獲取匯率
	exchangeRate, err := currency.GetExchangeRate(ts.URL)
	if err != nil {
		t.Fatalf("預期沒有錯誤，得到 %v", err)
	}

	// 檢查匯率
	if exchangeRate.ConversionRates["USD"] != 1.0 {
		t.Errorf("預期 USD 匯率為 1.0，得到 %v", exchangeRate.ConversionRates["USD"])
	}
	if exchangeRate.ConversionRates["EUR"] != 0.85 {
		t.Errorf("預期 EUR 匯率為 0.85，得到 %v", exchangeRate.ConversionRates["EUR"])
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

	convertedAmount, err := currency.ConvertCurrency(amount, from, to, rates)
	if err != nil {
		t.Fatalf("預期沒有錯誤，得到 %v", err)
	}

	if convertedAmount != expected {
		t.Errorf("預期 %v，得到 %v", expected, convertedAmount)
	}
}

func TestIntegration(t *testing.T) {
	// 建立一個測試伺服器來模擬 API 回應
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"conversion_rates":{"USD":1.0,"EUR":0.85},"base_code":"USD"}`))
	}))
	defer ts.Close()

	// 設置測試的環境變數
	os.Setenv("API_KEY", "test_api_key")
	os.Setenv("BASE_URL", ts.URL)

	// 獲取匯率
	exchangeRate, err := currency.GetExchangeRate(ts.URL)
	if err != nil {
		t.Fatalf("預期沒有錯誤，得到 %v", err)
	}

	// 轉換貨幣
	amount := 100.0
	from := "USD"
	to := "EUR"
	expected := 85.0

	convertedAmount, err := currency.ConvertCurrency(amount, from, to, exchangeRate.ConversionRates)
	if err != nil {
		t.Fatalf("預期沒有錯誤，得到 %v", err)
	}

	if convertedAmount != expected {
		t.Errorf("預期 %v，得到 %v", expected, convertedAmount)
	}
}
