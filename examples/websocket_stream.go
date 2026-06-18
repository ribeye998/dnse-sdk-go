package main

import (
	"context"
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
		wsURL = "wss://ws-openapi.dnse.com.vn" // Cổng wss streaming chính thức
	}

	// Hàm callback xử lý tin nhắn bất đồng bộ nhận về từ hệ thống
	onMessage := func(channel string, data []byte) {
		fmt.Printf("[WS RECEIVED] Kênh: %s | Dữ liệu thô: %s\n", channel, string(data))
	}

	wsClient := websocket.NewWSClient(wsURL, apiKey, apiSecret, onMessage)

	log.Println("Đang mở kết nối WebSocket tới DNSE...")
	err := wsClient.Connect(context.Background())
	if err != nil {
		log.Fatalf("Kết nối hoặc Xác thực (Auth) thất bại: %v", err)
	}
	log.Println("Kết nối và Đăng nhập thiết lập luồng thành công!")

	// Đăng ký nhận dữ liệu bảng giá (quote) giống hệt cấu trúc Python SDK
	err = wsClient.SubscribeMarketData("quote", []string{"HPG", "FPT", "SSI"})
	if err != nil {
		log.Printf("Lỗi đăng ký nhận dữ liệu bảng giá: %v", err)
	}

	// Đóng luồng chính sau 30 giây chạy thử nghiệm
	time.Sleep(30 * time.Second)
	wsClient.Close()
	log.Println("Đã ngắt kết nối WebSocket an toàn.")
}
