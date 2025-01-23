package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"currency-converter/currency"

	"github.com/joho/godotenv"
)

func main() {
	go currency.ResetRateLimit()

	// 載入 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatal("載入 .env 文件時出錯")
	}

	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("BASE_URL")
	if apiKey == "" || baseURL == "" {
		log.Fatal("必須設置 API_KEY 和 BASE_URL")
	}

	apiURL := baseURL

	// 從命令行獲取輸入的金額和幣種
	amount := flag.Float64("amount", 0, "要轉換的金額")
	from := flag.String("from", "", "要轉換的貨幣")
	to := flag.String("to", "USD", "目標貨幣")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "使用方法 :\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -amount float\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        要轉換的金額 (必填)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -from string\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        要轉換的貨幣 (必填)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -to string\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        目標貨幣 (預設：USD)\n")
	}
	flag.Parse()

	if *amount == 0 || *from == "" || *to == "" {
		flag.Usage()
		log.Println("所有參數 -amount, -from 都是必填的。")
		return
	}

	exchangeRate, fetchErr := currency.GetExchangeRate(apiURL)
	if fetchErr != nil {
		log.Fatalf("獲取匯率時出錯: %v", fetchErr)
	}

	convertedAmount, convErr := currency.ConvertCurrency(*amount, *from, *to, exchangeRate.ConversionRates)
	if convErr != nil {
		log.Fatalf("轉換貨幣時出錯: %v", convErr)
	}

	fmt.Printf("%.2f %s 等於 %.2f %s\n", *amount, *from, convertedAmount, *to)
}
