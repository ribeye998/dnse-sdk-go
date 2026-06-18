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

	accountNo := os.Getenv("DNSE_ACCOUNT_ID")
	marketType := "STOCK"
	symbol := "TCB"

	var response interface{}
	err := client.GetLoanPackages(context.Background(), accountNo, marketType, symbol, &response)
	if err != nil {
		log.Fatalf("Lỗi tính toán Sức mua (PPSE): %v", err)
	}

	fmt.Println("--- Thông tin Sức mua / Sức bán (PPSE) ---")
	fmt.Printf("%+v\n", response)
}
