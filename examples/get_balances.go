package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dnse-sdk-go/restdnse"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("DNSE_API_KEY")
	apiSecret := os.Getenv("DNSE_API_SECRET")
	baseURL := os.Getenv("DNSE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openapi.dnse.com.vn"
	}

	client := restdnse.NewClient(baseURL, apiKey, apiSecret)

	// Tài khoản mục tiêu cần kiểm tra số dư
	accountNo := os.Getenv("DNSE_ACCOUNT_ID")

	var response interface{}
	err := client.GetBalances(context.Background(), accountNo, &response)
	if err != nil {
		log.Fatalf("Lỗi truy vấn số dư tài khoản: %v", err)
	}

	fmt.Printf("--- Số dư tài khoản %s ---\n", accountNo)
	fmt.Printf("%+v\n", response)
}
