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
	// Nạp biến môi trường từ file .env
	_ = godotenv.Load()

	apiKey := os.Getenv("DNSE_API_KEY")
	apiSecret := os.Getenv("DNSE_API_SECRET")
	baseURL := os.Getenv("DNSE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openapi.dnse.com.vn"
	}

	if apiKey == "" || apiSecret == "" {
		log.Fatalf("Lỗi: Vui lòng cấu hình DNSE_API_KEY và DNSE_API_SECRET trong file .env")
	}

	client := restdnse.NewClient(baseURL, apiKey, apiSecret)

	var response interface{}
	err := client.GetAccounts(context.Background(), &response)
	if err != nil {
		log.Fatalf("Gọi API GetAccounts thất bại: %v", err)
	}

	fmt.Println("--- DANH SÁCH TIỂU KHOẢN DNSE ---")
	fmt.Printf("%+v\n", response)
}
