package main

import (
	"context"
	"log"
	
	"dnse-sdk-go/config"
	"dnse-sdk-go/restdnse"
	"dnse-sdk-go/websocket"
)

func main() {
	// 1. Đọc và kiểm tra cấu hình
	cfg := config.NewConfigFromEnv()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Cấu hình không hợp lệ: %v", err)
	}

	// 2. Khởi tạo REST Client
	restClient := restdnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)

	// 3. Khởi tạo WebSocket Stream Client
	wsClient := websocket.NewStreamClient(cfg.WSURL, cfg.APIKey, 1000)
	
	// Tiếp tục thực hiện logic giao dịch...
}