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

	symbol := "HPG"
	boardID := "" // Để trống tương đương giá trị None của Python SDK để hệ thống tự nhận diện bảng giá

	var response interface{}
	err := client.GetSecurityDefinition(context.Background(), symbol, boardID, &response)
	if err != nil {
		log.Fatalf("Lỗi lấy thông số kỹ thuật của mã: %v", err)
	}

	fmt.Printf("--- Biên độ trần/sàn & Thông tin mã %s ---\n", symbol)
	fmt.Printf("%+v\n", response)
}
