package restdnse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client đại diện cho HTTP REST Client giao tiếp với DNSE OpenAPI
type Client struct {
	hc        *http.Client
	baseURL   string
	apiKey    string
	apiSecret string
}

// NewClient khởi tạo một Client mới với cấu hình timeout chuẩn
func NewClient(baseURL, apiKey, apiSecret string) *Client {
	return &Client{
		hc:        &http.Client{Timeout: 10 * time.Second},
		baseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// sendRequest là hàm internal helper thực hiện đóng gói header, ký signature và decode dữ liệu
func (c *Client) sendRequest(ctx context.Context, method, path string, reqBody interface{}, resTarget interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	
	var serializedBody []byte
	if reqBody != nil {
		var err error
		serializedBody, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(serializedBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Tạo thông tin bảo mật và cấu hình Header theo đặc tả DNSE
	timestamp := GetCurrentTimestampMilli()
	signature := ComputeSignature(c.apiSecret, timestamp, serializedBody)
	SetupCommonHeaders(req, c.apiKey, timestamp, signature)

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("api returned error status: %d", resp.StatusCode)
	}

	if resTarget != nil {
		if err := json.NewDecoder(resp.Body).Decode(resTarget); err != nil {
			return fmt.Errorf("failed to decode response payload: %w", err)
		}
	}

	return nil
}