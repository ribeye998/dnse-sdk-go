package dnse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	gorilla "github.com/gorilla/websocket"
	"github.com/ribeye998/dnse-sdk-go/internal/signing"
)

const DefaultWSURL = "wss://ws-openapi.dnse.com.vn"

// StreamClient manages a persistent WebSocket connection to the DNSE stream server.
// Register callbacks with OnQuote, OnTick, etc. before calling Connect or Start*.
type StreamClient struct {
	baseURL           string
	apiKey            string
	apiSecret         string
	heartbeatInterval time.Duration

	mu           sync.Mutex
	tradingToken string
	accountNo    string
	encoding     string

	conn      *gorilla.Conn
	writeMu   sync.Mutex
	isClosing bool

	// Callbacks keyed on the message type code sent by the server
	// (e.g. "q", "t", "te", "b", "mi", "f", "o", "p", "a", "sd", "e").
	onQuote              func(symbol string, data map[string]interface{}) // "q"  — bid/ask depth
	onTick               func(symbol string, data map[string]interface{}) // "t"  — trade tick
	onTickExtra          func(symbol string, data map[string]interface{}) // "te" — trade extra
	onOHLC               func(symbol string, data map[string]interface{}) // "b"  — OHLC bar
	onMarketIndex        func(data map[string]interface{})                // "mi" / "idx"
	onForeign            func(symbol string, data map[string]interface{}) // "f"  — foreign investor
	onExpectedPrice      func(symbol string, data map[string]interface{}) // "e" / "ep"
	onSecurityDefinition func(symbol string, data map[string]interface{}) // "sd"
	onOrderUpdate        func(data map[string]interface{})                // "o"  — order event
	onPositionUpdate     func(data map[string]interface{})                // "p"  — position update
	onAccountUpdate      func(data map[string]interface{})                // "a"  — account balance
}

// NewStreamClient creates a WebSocket streaming client.
func NewStreamClient(baseURL, apiKey, apiSecret string) *StreamClient {
	return &StreamClient{
		baseURL:           strings.TrimSuffix(baseURL, "/"),
		apiKey:            apiKey,
		apiSecret:         apiSecret,
		encoding:          "json",
		heartbeatInterval: 10 * time.Second,
	}
}

// SetTradingToken sets the trading token for private channel subscriptions.
// Must be called before Connect.
func (s *StreamClient) SetTradingToken(token string) {
	s.mu.Lock()
	s.tradingToken = token
	s.mu.Unlock()
}

// SetAccountNo sets the account number used for private channel filtering.
// Must be called before Connect.
func (s *StreamClient) SetAccountNo(accountNo string) {
	s.mu.Lock()
	s.accountNo = accountNo
	s.mu.Unlock()
}

// OnQuote registers a handler for bid/ask depth messages (type code "q").
// Note: the wire key for the ask side is "offer", not "ask".
func (s *StreamClient) OnQuote(fn func(symbol string, data map[string]interface{})) {
	s.mu.Lock(); s.onQuote = fn; s.mu.Unlock()
}

// OnTick registers a handler for trade tick messages (type code "t").
func (s *StreamClient) OnTick(fn func(symbol string, data map[string]interface{})) {
	s.mu.Lock(); s.onTick = fn; s.mu.Unlock()
}

// OnTickExtra registers a handler for detailed trade messages with buy/sell volume
// aggregation (type code "te").
func (s *StreamClient) OnTickExtra(fn func(symbol string, data map[string]interface{})) {
	s.mu.Lock(); s.onTickExtra = fn; s.mu.Unlock()
}

// OnOHLC registers a handler for OHLC bar messages (type code "b"), covering
// both real-time (ChanOHLC) and closed (ChanOHLCClosed) candlestick channels.
func (s *StreamClient) OnOHLC(fn func(symbol string, data map[string]interface{})) {
	s.mu.Lock(); s.onOHLC = fn; s.mu.Unlock()
}

// OnMarketIndex registers a handler for market index update messages (type code "mi"/"idx").
func (s *StreamClient) OnMarketIndex(fn func(data map[string]interface{})) {
	s.mu.Lock(); s.onMarketIndex = fn; s.mu.Unlock()
}

// OnForeign registers a handler for foreign investor flow messages (type code "f").
func (s *StreamClient) OnForeign(fn func(symbol string, data map[string]interface{})) {
	s.mu.Lock(); s.onForeign = fn; s.mu.Unlock()
}

// OnExpectedPrice registers a handler for expected/indicative price messages (type code "e"/"ep").
func (s *StreamClient) OnExpectedPrice(fn func(symbol string, data map[string]interface{})) {
	s.mu.Lock(); s.onExpectedPrice = fn; s.mu.Unlock()
}

// OnSecurityDefinition registers a handler for security definition messages (type code "sd").
func (s *StreamClient) OnSecurityDefinition(fn func(symbol string, data map[string]interface{})) {
	s.mu.Lock(); s.onSecurityDefinition = fn; s.mu.Unlock()
}

// OnOrderUpdate registers a handler for order event messages (type code "o").
func (s *StreamClient) OnOrderUpdate(fn func(data map[string]interface{})) {
	s.mu.Lock(); s.onOrderUpdate = fn; s.mu.Unlock()
}

// OnPositionUpdate registers a handler for position update messages (type code "p").
func (s *StreamClient) OnPositionUpdate(fn func(data map[string]interface{})) {
	s.mu.Lock(); s.onPositionUpdate = fn; s.mu.Unlock()
}

// OnAccountUpdate registers a handler for account balance update messages (type code "a").
func (s *StreamClient) OnAccountUpdate(fn func(data map[string]interface{})) {
	s.mu.Lock(); s.onAccountUpdate = fn; s.mu.Unlock()
}

// Connect establishes the WebSocket connection and starts the read and heartbeat loops.
// Calling Connect on an already-connected client is a no-op.
func (s *StreamClient) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		return nil
	}

	// Build the target URL and path for signing.
	urlStr := s.baseURL
	if !strings.Contains(urlStr, "/v1/stream") {
		urlStr = strings.TrimSuffix(urlStr, "/") + "/v1/stream"
	}
	queryParts := []string{"encoding=" + s.encoding}
	if s.accountNo != "" {
		queryParts = append(queryParts, "accountNo="+s.accountNo)
	}
	qs := strings.Join(queryParts, "&")
	urlStr += "?" + qs
	pathTarget := "/v1/stream?" + qs

	dateVal, sigHeader := signing.BuildRESTHeaders(s.apiKey, s.apiSecret, "GET", pathTarget, time.Now())

	headers := make(http.Header)
	headers.Set("Date", dateVal)
	headers.Set("Accept", "application/json")
	headers.Set("x-api-key", s.apiKey)
	headers.Set("x-signature", sigHeader)
	if s.tradingToken != "" {
		headers.Set("trading-token", s.tradingToken)
	}

	dialer := &gorilla.Dialer{
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   1024 * 1024,
		WriteBufferSize:  1024 * 1024,
	}
	conn, _, err := dialer.Dial(urlStr, headers)
	if err != nil {
		return fmt.Errorf("dnse: stream dial: %w", err)
	}

	// Reset read deadline on every pong so a live-but-quiet connection stays open.
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(s.heartbeatInterval * 2))
	})
	_ = conn.SetReadDeadline(time.Now().Add(s.heartbeatInterval * 2))

	authMsg := signing.BuildWSAuthMessage(s.apiKey, s.apiSecret, s.tradingToken, s.accountNo)
	s.writeMu.Lock()
	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = conn.WriteJSON(authMsg)
	_ = conn.SetWriteDeadline(time.Time{})
	s.writeMu.Unlock()
	if err != nil {
		conn.Close()
		return fmt.Errorf("dnse: stream auth: %w", err)
	}

	s.conn = conn
	s.isClosing = false

	go s.readLoop(conn)
	go s.heartbeatLoop(conn)

	return nil
}

// Close gracefully shuts down the WebSocket connection.
func (s *StreamClient) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isClosing = true
	if s.conn != nil {
		_ = s.conn.WriteMessage(
			gorilla.CloseMessage,
			gorilla.FormatCloseMessage(gorilla.CloseNormalClosure, ""),
		)
		s.conn.Close()
		s.conn = nil
	}
}

func (s *StreamClient) writeJSON(v interface{}) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	if s.conn == nil {
		return fmt.Errorf("dnse: stream not connected")
	}
	_ = s.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err := s.conn.WriteJSON(v)
	_ = s.conn.SetWriteDeadline(time.Time{})
	return err
}

// dispatch routes an incoming JSON message to the registered callback.
// Routing is done on the "id" (message type code) field, not the channel name.
// Known type codes: "q" (quote/depth), "t" (tick), "te" (tick_extra),
// "b" (ohlc bar), "mi"/"idx" (market index), "f" (foreign), "e"/"ep" (expected price),
// "sd" (security definition), "o" (order), "p" (position), "a" (account).
func (s *StreamClient) dispatch(data []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	msgType, _ := msg["id"].(string)
	symbol, _ := msg["symbol"].(string)

	s.mu.Lock()
	onQuote := s.onQuote
	onTick := s.onTick
	onTickExtra := s.onTickExtra
	onOHLC := s.onOHLC
	onMarketIndex := s.onMarketIndex
	onForeign := s.onForeign
	onExpectedPrice := s.onExpectedPrice
	onSecDef := s.onSecurityDefinition
	onOrderUpdate := s.onOrderUpdate
	onPositionUpdate := s.onPositionUpdate
	onAccountUpdate := s.onAccountUpdate
	s.mu.Unlock()

	switch msgType {
	case "q":
		if onQuote != nil {
			onQuote(symbol, msg)
		}
	case "t":
		if onTick != nil {
			onTick(symbol, msg)
		}
	case "te":
		if onTickExtra != nil {
			onTickExtra(symbol, msg)
		}
	case "b":
		if onOHLC != nil {
			onOHLC(symbol, msg)
		}
	case "mi", "idx":
		if onMarketIndex != nil {
			onMarketIndex(msg)
		}
	case "f":
		if onForeign != nil {
			onForeign(symbol, msg)
		}
	case "e", "ep":
		if onExpectedPrice != nil {
			onExpectedPrice(symbol, msg)
		}
	case "sd":
		if onSecDef != nil {
			onSecDef(symbol, msg)
		}
	case "o":
		if onOrderUpdate != nil {
			onOrderUpdate(msg)
		}
	case "p":
		if onPositionUpdate != nil {
			onPositionUpdate(msg)
		}
	case "a":
		if onAccountUpdate != nil {
			onAccountUpdate(msg)
		}
	}
}

func (s *StreamClient) readLoop(conn *gorilla.Conn) {
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		s.dispatch(msg)
	}
}

func (s *StreamClient) heartbeatLoop(conn *gorilla.Conn) {
	ticker := time.NewTicker(s.heartbeatInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.writeMu.Lock()
		_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		err := conn.WriteMessage(gorilla.PingMessage, nil)
		_ = conn.SetWriteDeadline(time.Time{})
		s.writeMu.Unlock()
		if err != nil {
			return
		}
	}
}
