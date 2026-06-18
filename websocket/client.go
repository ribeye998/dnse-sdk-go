package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	// Sửa lại đúng module name trong file go.mod của bạn (dnse-sdk-go hoặc trading-client-go)
	"dnse-sdk-go/pkg/crypto"
)

// DNSEStreamClient đại diện cho hệ thống kết nối WebSocket streaming
type DNSEStreamClient struct {
	BaseURL           string
	ApiKey            string
	ApiSecret         string
	TradingToken      string
	AccountNo         string
	Encoding          string // "json" hoặc "msgpack"
	HeartbeatInterval time.Duration

	conn          *websocket.Conn
	mu            sync.Mutex
	writeMu       sync.Mutex
	isClosing     bool
	lastHeartbeat time.Time

	// Các hàm callback xử lý dữ liệu truyền ra ngoài
	OnQuote           func(symbol string, data map[string]interface{})
	OnTick            func(symbol string, data map[string]interface{})
	OnForeignInvestor func(symbol string, data map[string]interface{})
	OnOrderUpdate     func(data map[string]interface{})
	OnOrderEvent      func(data map[string]interface{})
	OnAccountUpdate   func(data map[string]interface{})
	OnPositionUpdate  func(data map[string]interface{})
}

// NewDNSEStreamClient khởi tạo thực thể Stream Client
func NewDNSEStreamClient(baseURL, apiKey, apiSecret string) *DNSEStreamClient {
	return &DNSEStreamClient{
		BaseURL:           baseURL,
		ApiKey:            apiKey,
		ApiSecret:         apiSecret,
		Encoding:          "json",
		HeartbeatInterval: 10 * time.Second,
	}
}

// Connect thực hiện handshake bảo mật và quản lý vòng đời luồng mạng
// Connect thực hiện handshake bảo mật và quản lý vòng đời luồng mạng
func (c *DNSEStreamClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	// 1. Định hình đường dẫn URL và toàn bộ Query Params TRƯỚC khi ký
	urlStr := c.BaseURL
	if !strings.Contains(urlStr, "/v1/stream") {
		urlStr = strings.TrimSuffix(urlStr, "/") + "/v1/stream"
	}

	var queryParams []string
	if c.Encoding != "" {
		queryParams = append(queryParams, "encoding="+c.Encoding)
	}
	if c.AccountNo != "" {
		queryParams = append(queryParams, "accountNo="+c.AccountNo)
	}

	// Khởi tạo path mục tiêu tham gia vào chuỗi ký chữ ký
	pathTarget := "/v1/stream"
	if len(queryParams) > 0 {
		queryString := strings.Join(queryParams, "&")
		urlStr += "?" + queryString
		// QUAN TRỌNG: Đính kèm luôn cả query string vào path ký nếu Gateway yêu cầu toàn vẹn URI
		pathTarget += "?" + queryString
	}

	// 2. Tạo chữ ký bảo mật dựa trên Path hoàn chỉnh
	// Sử dụng hàm BuildSignatureHeader chuẩn từ file hmac.go của bạn
	dateVal, sigHeader := crypto.BuildSignatureHeader(c.ApiKey, c.ApiSecret, "GET", pathTarget, time.Now())

	// 3. Khởi tạo Header bằng http.Header chuẩn (tự động xử lý canonical format)
	httpHeaders := make(http.Header)
	httpHeaders.Set("Date", dateVal)
	httpHeaders.Set("Accept", "application/json")

	// Thiết lập các trường định danh viết thường theo tiêu chuẩn OpenAPI Gateway v2
	httpHeaders.Set("x-api-key", c.ApiKey)
	httpHeaders.Set("x-signature", sigHeader)

	if c.TradingToken != "" {
		httpHeaders.Set("trading-token", c.TradingToken)
	}

	dialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   1024 * 1024,
		WriteBufferSize:  1024 * 1024,
	}

	// 4. Thực thi kết nối mạng (Dial) thiết lập Socket mở
	conn, _, err := dialer.Dial(urlStr, httpHeaders)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %w", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(c.HeartbeatInterval * 2))

	// 5. Sinh gói tin đăng nhập Auth Message cấp độ Application từ file hmac.go của bạn
	authMsg := crypto.CreateWSAuthMessage(c.ApiKey, c.ApiSecret, c.TradingToken, c.AccountNo)

	c.writeMu.Lock()
	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = conn.WriteJSON(authMsg)
	_ = conn.SetWriteDeadline(time.Time{})
	c.writeMu.Unlock()

	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to send json auth message: %w", err)
	}

	c.conn = conn
	c.isClosing = false

	go c.readLoop(c.conn)
	go c.heartbeatLoop(conn)

	return nil
}
func (c *DNSEStreamClient) dispatch(data []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	channel, _ := msg["channel"].(string)
	symbol, _ := msg["symbol"].(string)

	switch channel {
	case "quote":
		if c.OnQuote != nil {
			c.OnQuote(symbol, msg)
		}
	case "trade":
		if c.OnTick != nil {
			c.OnTick(symbol, msg)
		}
	case "foreign":
		if c.OnForeignInvestor != nil {
			c.OnForeignInvestor(symbol, msg)
		}
	case "order":
		if c.OnOrderUpdate != nil {
			c.OnOrderUpdate(msg)
		}
	case "order_event":
		if c.OnOrderEvent != nil {
			c.OnOrderEvent(msg)
		}
	case "account":
		if c.OnAccountUpdate != nil {
			c.OnAccountUpdate(msg)
		}
	case "position":
		if c.OnPositionUpdate != nil {
			c.OnPositionUpdate(msg)
		}
	}
}

func (c *DNSEStreamClient) readLoop(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		_ = conn.SetReadDeadline(time.Now().Add(c.HeartbeatInterval * 2))
		c.dispatch(message)
	}
}

func (c *DNSEStreamClient) heartbeatLoop(conn *websocket.Conn) {
	ticker := time.NewTicker(c.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.writeMu.Lock()
			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			err := conn.WriteMessage(websocket.PingMessage, []byte{})
			_ = conn.SetWriteDeadline(time.Time{})
			c.writeMu.Unlock()
			if err != nil {
				return
			}
		}
	}
}

func (c *DNSEStreamClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isClosing = true
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
