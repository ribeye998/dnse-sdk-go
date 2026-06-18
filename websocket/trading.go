package websocket

import "encoding/json"

// TradingStream bọc StreamClient để phục vụ riêng cho các luồng dữ liệu tài khoản cá nhân bảo mật cao
type TradingStream struct {
	client *StreamClient
}

func NewTradingStream(client *StreamClient) *TradingStream {
	return &TradingStream{client: client}
}

// OrderUpdate đại diện cho gói tin thay đổi trạng thái của một lệnh (order.py)
type OrderUpdate struct {
	AccountNo      string  `json:"accountNo"`
	OrderID        string  `json:"orderId"`
	OrderSeq       int64   `json:"orderSeq"`
	Symbol         string  `json:"symbol"`
	Side           string  `json:"side"`
	Status         string  `json:"status"` // PENDING, FILLED, CANCELED, REJECTED
	Price          float64 `json:"price"`
	Quantity       int64   `json:"quantity"`
	FilledQuantity int64   `json:"filledQuantity"`
	RejectReason   string  `json:"rejectReason,omitempty"`
}

// PositionUpdate đại diện thông tin thay đổi vị thế tài sản/danh mục (position.py)
type PositionUpdate struct {
	AccountNo    string  `json:"accountNo"`
	Symbol       string  `json:"symbol"`
	AvailQty     int64   `json:"availQty"`
	AvgBuyPrice  float64 `json:"avgBuyPrice"`
	CurrentPrice float64 `json:"currentPrice"`
}

// SubscribeOrderUpdates đăng ký nhận luồng thay đổi trạng thái lệnh của tiểu khoản giao dịch
func (ts *TradingStream) SubscribeOrderUpdates(accountNo string) error {
	return ts.client.Subscribe("order:" + accountNo)
}

// SubscribePositionUpdates đăng ký nhận luồng biến động danh mục tài sản của tiểu khoản
func (ts *TradingStream) SubscribePositionUpdates(accountNo string) error {
	return ts.client.Subscribe("position:" + accountNo)
}

// ParseOrderUpdate giải mã gói tin thô thành cấu trúc dữ liệu OrderUpdate cụ thể
func (ts *TradingStream) ParseOrderUpdate(msg WSMessage) (*OrderUpdate, error) {
	var data OrderUpdate
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// ParsePositionUpdate giải mã gói tin thô thành cấu trúc dữ liệu PositionUpdate cụ thể
func (ts *TradingStream) ParsePositionUpdate(msg WSMessage) (*PositionUpdate, error) {
	var data PositionUpdate
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}