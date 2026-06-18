package websocket

import "fmt"

func (c *DNSEStreamClient) SubscribeTrading(channels []string) error {
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
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("connection is closed")
	}
	return c.conn.WriteJSON(payload)
}
