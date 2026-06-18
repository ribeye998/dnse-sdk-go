package websocket

import "fmt"

func (c *DNSEStreamClient) SubscribeMarketData(channelName string, symbols []string) error {
	payload := map[string]interface{}{
		"action": "subscribe",
		"channels": []map[string]interface{}{
			{
				"name":    channelName,
				"symbols": symbols,
			},
		},
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("connection is closed")
	}
	return c.conn.WriteJSON(payload)
}
