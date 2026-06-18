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

	accountNo := "0001000115"
	marketType := "STOCK"
	symbol := "HPG"
	var price float64 = 26450      // Mức giá dự kiến đặt mua
	var loanPackageID int64 = 2396 // Gói margin dự kiến chọn

	var response interface{}
	err := client.GetPpse(context.Background(), accountNo, marketType, symbol, price, loanPackageID, &response)
	if err != nil {
		log.Fatalf("Lỗi tính toán Sức mua (PPSE): %v", err)
	}

	fmt.Println("--- Thông tin Sức mua / Sức bán (PPSE) ---")
	fmt.Printf("%+v\n", response)
}
