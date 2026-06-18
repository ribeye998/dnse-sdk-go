package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"trading-client-go/restdnse"
)

func main() {
	apiKey := os.Getenv("DNSE_API_KEY")
	apiSecret := os.Getenv("DNSE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		apiKey = "replace-with-api-key"
		apiSecret = "replace-with-api-secret"
	}

	// Khởi tạo client tới môi trường OpenAPI Production/Staging của DNSE
	client := restdnse.NewClient("https://openapi.dnse.com.vn", apiKey, apiSecret)

	var response interface{}
	// Endpoint /accounts tương đương client.get_accounts() trong Python SDK
	err := client.GetAccounts(context.Background(), &response)
	if err != nil {
		log.Fatalf("Lỗi gọi API GetAccounts: %v", err)
	}

	fmt.Printf("Trạng thái danh sách tài khoản thành công: %+v\n", response)
}
