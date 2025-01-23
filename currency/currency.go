package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type ExchangeRate struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
	BaseCode        string             `json:"base_code"`
}

var (
	cache          = make(map[string]float64)
	cacheMutex     sync.Mutex
	cacheTimestamp time.Time
	rateLimit      = 10
	rateLimitMutex sync.Mutex
	requestCount   = 0
)

// fetchExchangeRates 從第三方 API 獲取匯率，使用 HTTP 和 JSON
func fetchExchangeRates(apiURL string) (*ExchangeRate, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("獲取匯率失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("意外的狀態碼: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取回應內文失敗: %v", err)
	}

	var exchangeRate ExchangeRate
	if err := json.Unmarshal(body, &exchangeRate); err != nil {
		return nil, fmt.Errorf("解析 JSON 失敗: %v", err)
	}

	return &exchangeRate, nil
}

// ConvertCurrency 使用獲取的匯率將金額從一種貨幣轉換為另一種貨幣
func ConvertCurrency(amount float64, from, to string, rates map[string]float64) (float64, error) {
	fromRate, ok := rates[from]
	if !ok {
		return 0, fmt.Errorf("不支援的貨幣: %s", from)
	}

	toRate, ok := rates[to]
	if !ok {
		return 0, fmt.Errorf("不支援的貨幣: %s", to)
	}

	return amount * (toRate / fromRate), nil
}

// getCachedRates 返回快取的匯率，如果它們仍然有效
func getCachedRates() map[string]float64 {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(cacheTimestamp) < 15*time.Minute {
		return cache
	}

	return nil
}

// updateCache 使用最新的匯率更新本地快取
func updateCache(rates map[string]float64) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cache = rates
	cacheTimestamp = time.Now()
}

// rateLimitExceeded 檢查 API 請求的速率限制是否已超過
func rateLimitExceeded() bool {
	rateLimitMutex.Lock()
	defer rateLimitMutex.Unlock()

	if requestCount >= rateLimit {
		return true
	}

	requestCount++
	return false
}

// ResetRateLimit 每分鐘重置速率限制計數器
func ResetRateLimit() {
	for {
		time.Sleep(1 * time.Minute)
		rateLimitMutex.Lock()
		requestCount = 0
		rateLimitMutex.Unlock()
	}
}

// GetExchangeRate 處理獲取和快取匯率的邏輯
func GetExchangeRate(apiURL string) (*ExchangeRate, error) {
	var exchangeRate *ExchangeRate
	var fetchErr error

	for i := 0; i < 3; i++ {
		if rateLimitExceeded() {
			log.Println("速率限制已超過，使用快取的匯率")
			exchangeRate = &ExchangeRate{ConversionRates: getCachedRates()}
			break
		}

		exchangeRate, fetchErr = fetchExchangeRates(apiURL)
		if fetchErr == nil {
			updateCache(exchangeRate.ConversionRates)
			break
		}

		time.Sleep(time.Duration(i*i) * time.Second)
	}

	return exchangeRate, fetchErr
}
