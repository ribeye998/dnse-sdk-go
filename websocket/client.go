package websocket

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	conn       *websocket.Conn
	baseURL    string
	apiKey     string
	apiSecret  string
	mu         sync.Mutex
	stopChan   chan struct{}
	onMessage  func(msgType string, data []byte)
}

func NewWSClient(baseURL, apiKey, apiSecret string, onMessage func(msgType string, data []byte)) *WSClient {
	return &WSClient{
		baseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		stopChan:  make(chan struct{}),
		onMessage: onMessage,
	}
}

func generateNonce() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// Connect thực hiện quay số WebSocket và tự động kích hoạt cơ chế Xác thực (Auth)
func (ws *WSClient) Connect(ctx context.Context) error {
	u, err := url.Parse(ws.baseURL)
	if err != nil {
		return err
	}
	u.Path = "/v1/stream"
	u.RawQuery = "encoding=json"

	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	conn, _, err := dialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %w", err)
	}
	ws.conn = conn

	// Kích hoạt thủ tục đăng nhập (Auth) giống hệt auth.py bên Python SDK
	if err := ws.authenticate(); err != nil {
		ws.conn.Close()
		return fmt.Errorf("websocket authentication failed: %w", err)
	}

	go ws.readLoop()
	return nil
}

// authenticate thiết lập cấu trúc mã ký dựa trên chuỗi rỗng api_key:timestamp:nonce
func (ws *WSClient) authenticate() error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	nonce := generateNonce()
	
	// Công thức ký WebSocket riêng của DNSE: hmac_sha256(secret, "apiKey:timestamp:nonce")
	rawSig := fmt.Sprintf("%s:%s:%s", ws.apiKey, timestamp, nonce)
	mac := hmac.New(sha256.New, []byte(ws.apiSecret))
	mac.Write([]byte(rawSig))
	signature := hex.EncodeToString(mac.Sum(nil))

	authPayload := map[string]interface{}{
		"action": "auth",
		"params": map[string]string{
			"keyId":     ws.apiKey,
			"timestamp": timestamp,
			"nonce":     nonce,
			"signature": signature,
		},
	}

	return ws.send(authPayload)
}

func (ws *WSClient) send(v interface{}) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.conn == nil {
		return fmt.Errorf("connection is closed")
	}
	return ws.conn.WriteJSON(v)
}

func (ws *WSClient) readLoop() {
	defer func() {
		ws.conn.Close()
	}()

	for {
		select {
		case <-ws.stopChan:
			return
		default:
			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				log.Printf("Websocket read error: %v", err)
				return
			}

			// Phân tích sơ bộ bọc gói tin dựa theo cấu trúc "channel" hoặc "event"
			var generic map[string]interface{}
			if err := json.Unmarshal(message, &generic); err == nil {
				if channel, ok := generic["channel"].(string); ok {
					ws.onMessage(channel, message)
				} else {
					ws.onMessage("system", message)
				}
			}
		}
	}
}

func (ws *WSClient) Close() {
	close(ws.stopChan)
	ws.mu.