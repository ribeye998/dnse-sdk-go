package restdnse

import (
	"context"
	"fmt"
	"net/http"
)

// TradingTokenResponse trả về mã token giao dịch ngắn hạn sau khi xác thực
type TradingTokenResponse struct {
	TradingToken string `json:"tradingToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

// OrderParams chứa các tham số đầu vào để đẩy một lệnh mới lên sở
type OrderParams struct {
	AccountNo    string  `json:"accountNo"`
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`      // BUY, SELL
	OrderType    string  `json:"orderType"` // LO, MP, ATO, ATC
	Price        float64 `json:"price"`
	Quantity     int64   `json:"quantity"`
	TradingToken string  `json:"tradingToken"`
}

// OrderResponse chứa mã định danh lệnh trả về từ sàn
type OrderResponse struct {
	OrderID  string `json:"orderId"`
	OrderSeq int64  `json:"orderSeq"`
	Status   string `json:"status"`
}

// AccountBalance chứa thông tin số dư tài sản thực tế và sức mua
type AccountBalance struct {
	AccountNo       string  `json:"accountNo"`
	CashAvailable   float64 `json:"cashAvailable"`
	PurchasingPower float64 `json:"purchasingPower"`
}

// CreateTradingToken sinh token đặt lệnh dựa trên số tài khoản (create_trading_token.py)
func (c *Client) CreateTradingToken(ctx context.Context, accountNo string) (*TradingTokenResponse, error) {
	path := "/v1/trading/tokens"
	body := map[string]string{"accountNo": accountNo}
	var res TradingTokenResponse
	if err := c.sendRequest(ctx, http.MethodPost, path, body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// PostOrder gửi một yêu cầu đặt lệnh mới (post_order.py)
func (c *Client) PostOrder(ctx context.Context, params OrderParams) (*OrderResponse, error) {
	path := "/v1/trading/orders"
	var res OrderResponse
	if err := c.sendRequest(ctx, http.MethodPost, path, params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetBalances truy vấn số dư và sức mua của tài khoản mục tiêu (get_balances.py)
func (c *Client) GetBalances(ctx context.Context, accountNo string) (*AccountBalance, error) {
	path := fmt.Sprintf("/v1/trading/balances?accountNo=%s", accountNo)
	var res AccountBalance
	if err := c.sendRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}