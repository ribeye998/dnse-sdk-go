package dnse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CreateTradingToken requests a trading token via smart_otp or email_otp authentication.
// otpType is the authentication type ("smart_otp" or "email_otp"); passcode is the value.
// The token is stored automatically on the client for subsequent calls.
func (c *Client) CreateTradingToken(ctx context.Context, otpType, passcode string) (string, error) {
	payload := map[string]string{"otpType": otpType, "passcode": passcode}
	var result struct {
		TradingToken string `json:"tradingToken"`
	}
	if err := c.sendRequest(ctx, http.MethodPost, "/registration/trading-token", nil, payload, &result); err != nil {
		return "", err
	}
	c.SetTradingToken(result.TradingToken)
	return result.TradingToken, nil
}

// SendEmailOTP requests a one-time password sent to the registered email address.
func (c *Client) SendEmailOTP(ctx context.Context, email string) error {
	payload := map[string]string{"email": email, "otpType": "email_otp"}
	return c.sendRequest(ctx, http.MethodPost, "/registration/send-email-otp", nil, payload, nil)
}

// PlaceOrder submits a new order. orderCategory defaults to "NORMAL" when empty.
func (c *Client) PlaceOrder(ctx context.Context, marketType MarketType, orderCategory string, req OrderRequest) (json.RawMessage, error) {
	if orderCategory == "" {
		orderCategory = "NORMAL"
	}
	q := url.Values{
		"marketType":    {string(marketType)},
		"orderCategory": {orderCategory},
	}
	var result json.RawMessage
	if err := c.sendRequest(ctx, http.MethodPost, "/accounts/orders", q, req, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AmendOrder modifies the price and/or quantity of an existing unmatched order.
func (c *Client) AmendOrder(ctx context.Context, accountNo, orderID string, marketType MarketType, orderCategory string, req AmendOrderRequest) (json.RawMessage, error) {
	path := fmt.Sprintf("/accounts/%s/orders/%s", accountNo, orderID)
	q := url.Values{"marketType": {string(marketType)}}
	if orderCategory != "" {
		q.Set("orderCategory", orderCategory)
	}
	var result json.RawMessage
	if err := c.sendRequest(ctx, http.MethodPut, path, q, req, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CancelOrder cancels an unmatched order.
func (c *Client) CancelOrder(ctx context.Context, accountNo, orderID string, marketType MarketType, orderCategory string) error {
	path := fmt.Sprintf("/accounts/%s/orders/%s", accountNo, orderID)
	q := url.Values{"marketType": {string(marketType)}}
	if orderCategory != "" {
		q.Set("orderCategory", orderCategory)
	}
	return c.sendRequest(ctx, http.MethodDelete, path, q, nil, nil)
}

// ClosePosition closes an open derivative position.
// payload is the optional request body as defined by the API; pass nil if none.
func (c *Client) ClosePosition(ctx context.Context, accountNo, positionID string, marketType MarketType, payload interface{}) error {
	path := fmt.Sprintf("/accounts/%s/positions/%s/close", accountNo, positionID)
	q := url.Values{"marketType": {string(marketType)}}
	return c.sendRequest(ctx, http.MethodPost, path, q, payload, nil)
}
