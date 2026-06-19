// Package dnse is a Go SDK for the DNSE OpenAPI — Vietnam's stock and
// derivative trading platform. It provides a REST client (Client) for
// account management, trading, and market data, and a WebSocket streaming
// client (StreamClient) for real-time market and trading events.
package dnse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ribeye998/dnse-sdk-go/internal/signing"
)

const (
	DefaultBaseURL    = "https://openapi.dnse.com.vn"
	defaultAPIVersion = "2026-05-07"
	defaultTimeout    = 15 * time.Second
)

// Client is the DNSE REST API client. Use NewClient to create one.
type Client struct {
	baseURL      string
	apiKey       string
	apiSecret    string
	apiVersion   string
	tradingToken string
	hc           *http.Client
}

// Option configures a Client at construction time.
type Option func(*Client)

// WithTimeout sets the HTTP client timeout (default 15 s).
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.hc.Timeout = d }
}

// WithHTTPClient replaces the underlying http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.hc = hc }
}

// WithAPIVersion overrides the API version header (default "2026-05-07").
func WithAPIVersion(v string) Option {
	return func(c *Client) { c.apiVersion = v }
}

// NewClient creates a new DNSE REST client.
func NewClient(baseURL, apiKey, apiSecret string, opts ...Option) *Client {
	c := &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		apiVersion: defaultAPIVersion,
		hc:         &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// SetTradingToken stores the trading token used for order placement and
// other authenticated trading operations.
func (c *Client) SetTradingToken(token string) {
	c.tradingToken = token
}

// TradingToken returns the currently configured trading token.
func (c *Client) TradingToken() string {
	return c.tradingToken
}

func (c *Client) sendRequest(ctx context.Context, method, path string, query url.Values, body interface{}, dst interface{}) error {
	reqURL := c.baseURL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("dnse: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("dnse: build request: %w", err)
	}

	dateVal, sigHeader := signing.BuildRESTHeaders(c.apiKey, c.apiSecret, method, path, time.Now())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Date", dateVal)
	req.Header.Set("X-Signature", sigHeader)
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("version", c.apiVersion)
	if c.tradingToken != "" {
		req.Header.Set("Trading-Token", c.tradingToken)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("dnse: http: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("dnse: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: string(respBytes)}
	}

	if dst != nil && len(respBytes) > 0 {
		if err := json.Unmarshal(respBytes, dst); err != nil {
			return fmt.Errorf("dnse: unmarshal response: %w", err)
		}
	}
	return nil
}
