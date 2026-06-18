package websocket

import "encoding/json"

// MarketDataStream bọc StreamClient để quản lý chuyên biệt các luồng dữ liệu thị trường công khai
type MarketDataStream struct {
	client *StreamClient
}

func NewMarketDataStream(client *StreamClient) *MarketDataStream {
	return &MarketDataStream{client: client}
}

// QuoteData đại diện gói tin cập nhật bước giá mua/bán thời gian thực (quote.py)
type QuoteData struct {
	Symbol        string    `json:"symbol"`
	Sequence      int64     `json:"sequence"`
	MatchPrice    float64   `json:"matchPrice"`
	MatchQuantity int64     `json:"matchQuantity"`
	TotalVolume   int64     `json:"totalVolume"`
	BidPrices     []float64 `json:"bidPrices"`
	AskPrices     []float64 `json:"askPrices"`
}

// TradeData đại diện thông tin một giao dịch khớp lệnh trực tiếp vừa phát sinh trên sàn (trade.py)
type TradeData struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Quantity  int64   `json:"quantity"`
	Side      string  `json:"side"` // BUY hoặc SELL
	MatchedAt int64   `json:"matchedAt"`
}

// SubscribeQuotes đăng ký luồng cập nhật bước giá cho danh sách mã cổ phiếu
func (md *MarketDataStream) SubscribeQuotes(symbols []string) error {
	for _, sym := range symbols {
		if err := md.client.Subscribe("quote:" + sym); err != nil {
			return err
		}
	}
	return nil
}

// SubscribeTrades đăng ký luồng cập nhật khớp lệnh tức thời cho danh sách mã cổ phiếu
func (md *MarketDataStream) SubscribeTrades(symbols []string) error {
	for _, sym := range symbols {
		if err := md.client.Subscribe("trade:" + sym); err != nil {
			return err
		}
	}
	return nil
}

// ParseQuote bóc tách gói tin WS thô thành cấu trúc dữ liệu QuoteData
func (md *MarketDataStream) ParseQuote(msg WSMessage) (*QuoteData, error) {
	var data QuoteData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// ParseTrade bóc tách gói tin WS thô thành cấu trúc dữ liệu TradeData
func (md *MarketDataStream) ParseTrade(msg WSMessage) (*TradeData, error) {
	var data TradeData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}