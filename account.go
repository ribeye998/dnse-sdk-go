package dnse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetAccounts returns the list of trading sub-accounts as raw JSON.
func (c *Client) GetAccounts(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, "/accounts", nil, nil, &result)
}

// GetBalances returns asset balances and purchasing power for an account.
func (c *Client) GetBalances(ctx context.Context, accountNo string) (*Balance, error) {
	var result Balance
	path := fmt.Sprintf("/accounts/%s/balances", accountNo)
	return &result, c.sendRequest(ctx, http.MethodGet, path, nil, nil, &result)
}

// GetLoanPackages returns available margin loan packages for an account.
// symbol is optional — pass an empty string to omit it.
func (c *Client) GetLoanPackages(ctx context.Context, accountNo string, marketType MarketType, symbol string) (json.RawMessage, error) {
	path := fmt.Sprintf("/accounts/%s/loan-packages", accountNo)
	q := url.Values{"marketType": {string(marketType)}}
	if symbol != "" {
		q.Set("symbol", symbol)
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetPPSE estimates purchasing/selling power for a prospective order.
func (c *Client) GetPPSE(ctx context.Context, accountNo string, marketType MarketType, symbol string, price float64, loanPackageID int64) (json.RawMessage, error) {
	path := fmt.Sprintf("/accounts/%s/ppse", accountNo)
	q := url.Values{
		"marketType":    {string(marketType)},
		"symbol":        {symbol},
		"price":         {strconv.FormatFloat(price, 'f', -1, 64)},
		"loanPackageId": {strconv.FormatInt(loanPackageID, 10)},
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetPositions returns open portfolio positions for an account.
func (c *Client) GetPositions(ctx context.Context, accountNo string, marketType MarketType) ([]Position, error) {
	path := fmt.Sprintf("/accounts/%s/positions", accountNo)
	q := url.Values{"marketType": {string(marketType)}}
	var wrapper struct {
		Positions []Position `json:"positions"`
	}
	err := c.sendRequest(ctx, http.MethodGet, path, q, nil, &wrapper)
	return wrapper.Positions, err
}

// GetPositionByID returns the details of a specific position.
func (c *Client) GetPositionByID(ctx context.Context, positionID string, marketType MarketType) (*Position, error) {
	path := fmt.Sprintf("/accounts/positions/%s", positionID)
	q := url.Values{"marketType": {string(marketType)}}
	var result Position
	return &result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetOrders returns the intraday order book for an account.
// orderCategory is optional (e.g. "NORMAL"); pass empty string to omit.
func (c *Client) GetOrders(ctx context.Context, accountNo string, marketType MarketType, orderCategory string) ([]Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders", accountNo)
	q := url.Values{"marketType": {string(marketType)}}
	if orderCategory != "" {
		q.Set("orderCategory", orderCategory)
	}
	var wrapper struct {
		Orders []Order `json:"orders"`
	}
	err := c.sendRequest(ctx, http.MethodGet, path, q, nil, &wrapper)
	return wrapper.Orders, err
}

// GetOrderDetail returns the details of a single intraday order.
func (c *Client) GetOrderDetail(ctx context.Context, accountNo, orderID string, marketType MarketType, orderCategory string) (*Order, error) {
	path := fmt.Sprintf("/accounts/%s/orders/%s", accountNo, orderID)
	q := url.Values{"marketType": {string(marketType)}}
	if orderCategory != "" {
		q.Set("orderCategory", orderCategory)
	}
	var result Order
	return &result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetOrderHistory returns the historical order book for an account.
func (c *Client) GetOrderHistory(ctx context.Context, accountNo string, marketType MarketType, p OrderHistoryParams) (json.RawMessage, error) {
	path := fmt.Sprintf("/accounts/%s/orders/history", accountNo)
	q := url.Values{"marketType": {string(marketType)}}
	if p.From != "" {
		q.Set("from", p.From)
	}
	if p.To != "" {
		q.Set("to", p.To)
	}
	if p.PageSize > 0 {
		q.Set("pageSize", strconv.Itoa(p.PageSize))
	}
	if p.PageIndex > 0 {
		q.Set("pageIndex", strconv.Itoa(p.PageIndex))
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetExecutions returns the fill details for a specific order.
func (c *Client) GetExecutions(ctx context.Context, accountNo, orderID string, marketType MarketType, orderCategory string) (json.RawMessage, error) {
	path := fmt.Sprintf("/accounts/%s/executions/%s", accountNo, orderID)
	q := url.Values{"marketType": {string(marketType)}}
	if orderCategory != "" {
		q.Set("orderCategory", orderCategory)
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetCorporateActionHistory returns the corporate action history for an account.
func (c *Client) GetCorporateActionHistory(ctx context.Context, accountNo string, p CorporateActionHistoryParams) (json.RawMessage, error) {
	path := fmt.Sprintf("/accounts/%s/corporate-action-history", accountNo)
	q := url.Values{}
	if p.Symbol != "" {
		q.Set("symbol", p.Symbol)
	}
	if p.CAType != "" {
		q.Set("caType", p.CAType)
	}
	if p.CAStatus != "" {
		q.Set("caStatus", p.CAStatus)
	}
	if p.PageIndex > 0 {
		q.Set("pageIndex", strconv.Itoa(p.PageIndex))
	}
	if p.PageSize > 0 {
		q.Set("pageSize", strconv.Itoa(p.PageSize))
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetBrokerAccounts returns the accounts managed by the authenticated broker.
func (c *Client) GetBrokerAccounts(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, "/brokers/accounts/care-by", nil, nil, &result)
}
