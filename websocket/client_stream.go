package websocket

import (
	"fmt"
	"log"
)

// StartMarketData thực thi kết nối và đăng ký nhận dữ liệu luồng thị trường public
func (c *DNSEStreamClient) StartMarketData(symbols []string, includeQuotes, includeTicks, includeIndices bool) error {
	// Hàm Connect() không nhận tham số theo cấu trúc thiết kế hmac.go mới của bạn
	if err := c.Connect(); err != nil {
		return fmt.Errorf("không thể thiết lập kết nối streaming: %w", err)
	}

	if includeQuotes {
		_ = c.SubscribeMarketData("quote", symbols)
	}
	if includeTicks {
		_ = c.SubscribeMarketData("trade", symbols)
	}
	if includeIndices {
		_ = c.SubscribeMarketData("market_index", []string{"VN30", "HNX"})
	}

	log.Printf("[DNSE] Kích hoạt thành công luồng Market Data cho các mã: %v\n", symbols)
	return nil
}

// StartTradingData thực thi kết nối và đăng ký nhận luồng quản lý tài sản private
func (c *DNSEStreamClient) StartTradingData() error {
	if err := c.Connect(); err != nil {
		return fmt.Errorf("không thể thiết lập kết nối private streaming: %w", err)
	}

	privateChannels := []string{"order", "position", "account"}
	if err := c.SubscribeTrading(privateChannels); err != nil {
		return fmt.Errorf("lỗi đăng ký luồng private trading: %w", err)
	}

	log.Println("[DNSE] Kích hoạt thành công luồng dữ liệu Trading Data Stream.")
	return nil
}
