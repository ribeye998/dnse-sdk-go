package websocket

// SubscribeTrading đăng ký theo dõi trạng thái lệnh và vị thế của các tài khoản mục tiêu
func (ws *WSClient) SubscribeTrading(channels []string) error {
	// channels có thể truyền vào: ["order", "position"]
	var channelList []map[string]interface{}
	for _, ch := range channels {
		channelList = append(channelList, map[string]interface{}{
			"name": ch,
		})
	}

	payload := map[string]interface{}{
		"action":   "subscribe",
		"channels": channelList,
	}
	return ws.send(payload)
}
