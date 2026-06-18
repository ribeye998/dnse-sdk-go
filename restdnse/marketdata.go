package restdnse

import (
	"context"
	"fmt"
	"net/http"
)

// SecurityDefinition chứa thông tin biên độ và đặc tính của mã chứng khoán
type SecurityDefinition struct {
	Symbol     string  `json:"symbol"`
	Market     string  `json:"market"` // HOSE, HNX, UPCOM, DERIVATIVE
	FloorPrice float64 `json:"floorPrice"`
	CeilPrice  float64 `json:"ceilPrice"`
	RefPrice   float64 `json:"refPrice"`
}

// OHLC đại diện cho một thanh nến lịch sử dữ liệu đồ thị
type OHLC struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// GetSecurityDefinition tương đương marketdata-api/get_security_definition.py
func (c *Client) GetSecurityDefinition(ctx context.Context, symbol string) (*SecurityDefinition, error) {
	path := fmt.Sprintf("/v1/market/securities/%s", symbol)
	var res SecurityDefinition
	if err := c.sendRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetOHLC tương đương marketdata-api/get_ohlc.py
func (c *Client) GetOHLC(ctx context.Context, symbol, resolution string, from, to int64) ([]OHLC, error) {
	path := fmt.Sprintf("/v1/market/ohlc?symbol=%s&resolution=%s&from=%d&to=%d", symbol, resolution, from, to)
	var res []OHLC
	if err := c.sendRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}