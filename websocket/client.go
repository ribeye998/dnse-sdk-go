package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSMessage định nghĩa cấu trúc khung gói tin cơ bản từ WebSocket của DNSE
type WSMessage struct {
	Topic string          `json:"topic"`
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// StreamClient quản lý kết nối TCP WebSocket hạ tầng
type StreamClient struct {
	wsURL     string
	apiKey    string
	conn      *websocket.Conn
	mu        sync.Mutex
	isClosed  bool
	outChan   chan WSMessage
	errChan   chan error
}

// NewStreamClient khởi tạo một bộ kết nối WebSocket với kích thước bộ đệm tùy chọn
func NewStreamClient(wsURL, apiKey string, bufferSize int) *StreamClient {
	return &StreamClient{
		wsURL:   wsURL,
		apiKey:  apiKey,
		outChan: make(chan WSMessage, bufferSize),
		errChan: make(chan error, 64),
	}
}

// Connect thiết lập kết nối Handshake và kích hoạt luồng đọc ngầm
func (sc *StreamClient) Connect(ctx context.Context) error {
	u, err := url.Parse(sc.wsURL)
	if err != nil {
		return fmt.Errorf("invalid websocket url: %w", err)
	}

	q := u.Query()
	q.Set("apiKey", sc.apiKey)
	u.RawQuery = q.Encode()

	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %w", err)
	}

	sc.conn = conn
	sc.isClosed = false

	go sc.readLoop()

	return nil
}

// Subscribe gửi yêu cầu đăng ký theo dõi một Topic nghiệp vụ cụ thể
func (sc *StreamClient) Subscribe(topic string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.conn == nil || sc.isClosed {
		return fmt.Errorf("websocket connection is inactive")
	}

	req := map[string]string{
		"action": "subscribe",
		"topic":  topic,
	}
	return sc.conn.WriteJSON(req)
}

func (sc *StreamClient) readLoop() {
	defer sc.Close()

	for {
		_, message, err := sc.conn.ReadMessage()
		if err != nil {
			if !sc.isClosed {
				sc.errChan <- fmt.Errorf("websocket read connection broken: %w", err)
			}
			return
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue // Bỏ qua nếu gói tin thô lỗi định dạng hệ thống
		}

		// Đẩy dữ liệu vào channel theo mô hình Non-blocking để bảo vệ hiệu năng
		select {
		case sc.outChan <- wsMsg:
		default:
			// Bộ đệm đầy, bỏ qua tin cũ để tránh nghẽn toàn luồng hệ thống
		}
	}
}

// Messages xuất kênh dữ liệu ra ngoài cho các module tầng trên tiêu thụ
func (sc *StreamClient) Messages() <-chan WSMessage { return sc.outChan }
func (sc *StreamClient) Errors() <-chan error       { return sc.errChan }

// Close đóng kết nối an toàn bảo vệ tài nguyên mạng
func (sc *StreamClient) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.isClosed {
		return
	}
	sc.isClosed = true
	if sc.conn != nil {
		sc.conn.Close()
	}
}