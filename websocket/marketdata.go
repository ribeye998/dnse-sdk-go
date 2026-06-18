package websocket

// SubscribeMarketData đăng ký nhận dữ liệu luồng trực tuyến (Ví dụ: quote, trade, ohlc)
func (ws *WSClient) SubscribeMarketData(channelName string, symbols []string) error {
	payload := map[string]interface{}{
		"action": "subscribe",
		"channels": []map[string]interface{}{
			{
				"name":    channelName, // e.g., "quote", "trade", "ohlc"
				"symbols": symbols,     // e.g., ["HPG", "FPT"]
			},
		},
	}
	return ws.send(payload)
}

// UnsubscribeMarketData hủy nhận dữ liệu dòng
func (ws *WSClient) UnsubscribeMarketData(channelName string, symbols []string) error {
	payload := map[string]interface{}{
		"action": "unsubscribe",
		"channels": []map[string]interface{}{
			{
				"name":    channelName,
				"symbols": symbols,
			},
		},
	}
	return ws.send(payload)
}
