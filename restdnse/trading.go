package restdnse

import (
	"context"
	"fmt"
	"net/http"
)

// AccountBalance chứa thông tin tài sản thực tế và sức mua
type AccountBalance struct {
	AccountNo       string  `json:"accountNo"`
	CashAvailable   float64 `json:"cashAvailable"`
	PurchasingPower float64 `json:"purchasingPower"`
}

// GetAccounts truy vấn danh sách tiểu khoản giao dịch (client.get_accounts)
func (c *Client) GetAccounts(ctx context.Context, resTarget interface{}) error {
	return c.sendRequest(ctx, http.MethodGet, "/accounts", nil, nil, resTarget)
}

// GetBalances truy vấn số dư chi tiết của một tiểu khoản (client.get_balances)
func (c *Client) GetBalances(ctx context.Context, accountNo string, resTarget interface{}) error {
	path := fmt.Sprintf("/accounts/%s/balances", accountNo)
	return c.sendRequest(ctx, http.MethodGet, path, nil, nil, resTarget)
}

// GetLoanPackages lấy danh sách mã gói vay ký quỹ khả dụng cho mã chứng khoán (client.get_loan_packages)
func (c *Client) GetLoanPackages(ctx context.Context, accountNo, marketType, symbol string, resTarget interface{}) error {
	path := fmt.Sprintf("/accounts/%s/loan-packages", accountNo)
	query := map[string]string{"marketType": marketType}
	if symbol != "" {
		query["symbol"] = symbol
	}
	return c.sendRequest(ctx, http.MethodGet, path, query, nil, resTarget)
}

// GetPpse tính toán sức mua và sức bán trước khi đặt lệnh (client.get_ppse)
func (c *Client) GetPpse(ctx context.Context, accountNo, marketType, symbol string, price float64, loanPackageID int64, resTarget interface{}) error {
	path := fmt.Sprintf("/accounts/%s/ppse", accountNo)
	query := map[string]string{
		"marketType":    marketType,
		"symbol":        symbol,
		"price":         fmt.Sprintf("%.2f", price),
		"loanPackageId": fmt.Sprintf("%d", loanPackageID),
	}
	return c.sendRequest(ctx, http.MethodGet, path, query, nil, resTarget)
}

// PostOrder thực hiện đẩy lệnh mới lên sàn giao dịch (client.post_order)
func (c *Client) PostOrder(ctx context.Context, marketType string, payload map[string]interface{}, tradingToken string, resTarget interface{}) error {
	// Để truyền trading-token độc lập qua header, ta gọi trực tiếp sendRequest
	// nhưng payload của DNSE yêu cầu token đính kèm trong logic nghiệp vụ hoặc custom request.
	// Thiết lập query và gửi body dữ liệu
	query := map[string]string{"marketType": marketType}
	return c.sendRequest(ctx, http.MethodPost, "/accounts/orders", query, payload, resTarget)
}
