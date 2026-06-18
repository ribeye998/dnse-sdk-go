package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"trading-client-go/websocket"
)

func main() {
	// Khởi tạo callback xử lý sự kiện bất đồng bộ nhận từ sàn DNSE
	onMessage := func(channel string, data []byte) {
		fmt.Printf("[WS RECEIVED] Kênh: %s | Nội dung: %s\n", channel, string(data))
	}

	// Cổng WebSocket DNSE (WSS)
	wsClient := websocket.NewWSClient("wss://openapi.dnse.com.vn", "your-api-key", "your-api-secret", onMessage)

	log.Println("Đang kết nối tới hệ thống DNSE Stream...")
	err := wsClient.Connect(context.Background())
	if err != nil {
		log.Fatalf("Kết nối WebSocket thất bại: %v", err)
	}
	log.Println("Kết nối và xác thực (Auth) thành công!")

	// Đăng ký nhận giá trực tuyến (Quote) của mã HPG và FPT giống như Python SDK
	err = wsClient.SubscribeMarketData("quote", []string{"HPG", "FPT"})
	if err != nil {
		log.Printf("Lỗi đăng ký nhận dữ liệu bảng giá: %v", err)
	}

	// Treo luồng chính trong 1 phút để hứng dữ liệu trả về trước khi đóng
	time.Sleep(1 * time.Minute)
	wsClient.Close()
	log.Println("Đã đóng kết nối an toàn.")
}
