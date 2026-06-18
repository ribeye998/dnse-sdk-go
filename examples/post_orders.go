package main

import (
	"context"
	"fmt"
	"log"

	"trading-client-go/restdnse"
)

func main() {
	client := restdnse.NewClient("https://openapi.dnse.com.vn", "your-api-key", "your-api-secret")

	// Payload cấu hình khớp lệnh mua chứng khoán cơ sở (STOCK) giống ví dụ python
	payload := map[string]interface{}{
		"accountNo":     "0001000115",
		"symbol":        "HPG",
		"side":          "NB", // Mua
		"orderType":     "LO", // Lệnh Giới hạn
		"price":         25950,
		"quantity":      100,
		"loanPackageId": 2396,
	}

	tradingToken := "replace-with-actual-trading-token"

	// Gọi hàm đặt lệnh
	res, err := client.PostOrder(context.Background(), "STOCK", payload, tradingToken)
	if err != nil {
		log.Fatalf("Đặt lệnh thất bại: %v", err)
	}

	fmt.Printf("Kết quả phản hồi đặt lệnh từ sàn DNSE: %+v\n", res)
}
