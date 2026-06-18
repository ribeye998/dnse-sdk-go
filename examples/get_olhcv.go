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

	// Cấu hình tham số truy vấn nến lịch sử giống hệt Python SDK
	queryParams := map[string]string{
		"symbol":     "HPG",
		"resolution": "1",          // Độ rộng nến: 1 phút
		"from":       "1735689600", // Thời gian Unix Timestamp bắt đầu
		"to":         "1735776000", // Thời gian Unix Timestamp kết thúc
	}

	var response interface{}
	// Hàm GetOHLC tự động truyền thêm trường "type": "STOCK" vào query parameter
	err := client.GetOHLC(context.Background(), "STOCK", queryParams, &response)
	if err != nil {
		log.Fatalf("Lỗi tải lịch sử đồ thị nến OHLC: %v", err)
	}

	fmt.Println("--- Dữ liệu đồ thị nến lịch sử (OHLC) ---")
	fmt.Printf("%+v\n", response)
}
