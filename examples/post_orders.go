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

	client := restdnse.NewClient("https://openapi.dnse.com.vn", os.Getenv("DNSE_API_KEY"), os.Getenv("DNSE_API_SECRET"))

	// Định dạng map payload đặt lệnh chuẩn khớp với ví dụ python trading-api/post_order.py
	payload := map[string]interface{}{
		"accountNo":     "0001000115",
		"symbol":        "HPG",
		"side":          "NB", // NB nghĩa là Mua (New Buy)
		"orderType":     "LO",
		"price":         25950,
		"quantity":      100,
		"loanPackageId": 2396,
	}

	tradingToken := "replace-with-actual-trading-token"

	var response interface{}
	err := client.PostOrder(context.Background(), "STOCK", payload, tradingToken, &response)
	if err != nil {
		log.Fatalf("Đặt lệnh thất bại: %v", err)
	}

	fmt.Printf("Kết quả phản hồi đặt lệnh từ sàn DNSE:\n%v\n", response)
}
