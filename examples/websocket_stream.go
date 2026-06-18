package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"dnse-sdk-go/websocket"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("DNSE_API_KEY")
	apiSecret := os.Getenv("DNSE_API_SECRET")
	wsURL := os.Getenv("DNSE_WS_URL")
	if wsURL == "" {
		wsURL = "wss://ws-openapi.dnse.com.vn"
	}

	// 1. Khởi tạo Stream Client cấp cao
	streamClient := websocket.NewDNSEStreamClient(wsURL, apiKey, apiSecret)

	// 2. Định nghĩa các hàm callback hứng dữ liệu sự kiện đã định nghĩa trong client_stream.go
	// Hàm callback xử lý dữ liệu sự kiện chuẩn map[string]interface{}
	streamClient.OnQuote = func(symbol string, data map[string]interface{}) {
		fmt.Printf("[SỰ KIỆN QUOTE] Mã: %s | Dữ liệu: %v\n", symbol, data)
	}

	streamClient.OnOrderUpdate = func(data map[string]interface{}) {
		fmt.Printf("[SỰ KIỆN LỆNH] Trạng thái lệnh thay đổi: %v\n", data)
	}

	// 3. Kích hoạt nhận dữ liệu thị trường (Market Data) cho các mã mục tiêu
	symbols := []string{"HPG", "FPT", "SSI"}
	log.Println("Đang chạy StartMarketData...")
	err := streamClient.StartMarketData(symbols, true, true, false)
	if err != nil {
		log.Fatalf("Khởi chạy Market Data lỗi: %v", err)
	}

	// 4. Kích hoạt nhận dữ liệu tài khoản đặt lệnh (Trading Data)
	log.Println("Đang chạy StartTradingData...")
	err = streamClient.StartTradingData()
	if err != nil {
		log.Printf("Khởi chạy Trading Data lỗi (Bỏ qua nếu chỉ dùng Key Market): %v", err)
	}

	// Duy trì ứng dụng hứng luồng sự kiện trong 45 giây trước khi giải phóng tài nguyên
	time.Sleep(45 * time.Second)
	streamClient.Close()
	log.Println("Hệ thống đóng kết nối stream an toàn.")
}
