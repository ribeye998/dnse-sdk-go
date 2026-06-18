package restdnse

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	// Import package crypto dùng chung
	"trading-client-go/pkg/crypto"
)

type Client struct {
	hc         *http.Client
	baseURL    string
	apiKey     string
	apiSecret  string
	apiVersion string
}

// NewClient khởi tạo REST Client kết nối với DNSE OpenAPI
func NewClient(baseURL, apiKey, apiSecret string) *Client {
	return &Client{
		hc: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		apiVersion: "2026-05-07", // Đồng bộ phiên bản với hệ thống Python SDK
	}
}

// generateNonce sinh chuỗi ngẫu nhiên hex (tương đương uuid4().hex trong Python)
func generateNonce() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// sendRequest đóng gói toàn bộ logic tạo HTTP Signature, gán Header và giải mã Response
func (c *Client) sendRequest(ctx context.Context, method, path string, query map[string]string, reqBody interface{}, resTarget interface{}) error {
	// 1. Xây dựng URL kèm tham số Query Params
	reqURL := fmt.Sprintf("%s%s", c.baseURL, path)
	if len(query) > 0 {
		var qStr []string
		for k, v := range query {
			if v != "" {
				qStr = append(qStr, fmt.Sprintf("%s=%s", k, v))
			}
		}
		if len(qStr) > 0 {
			reqURL += "?" + strings.Join(qStr, "&")
		}
	}

	// 2. Mã hóa dữ liệu Body nếu có
	var serializedBody []byte
	if reqBody != nil {
		var err error
		serializedBody, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bytes.NewBuffer(serializedBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Tính toán HTTP Signature từ package crypto chung
	dateValue := crypto.FormatDateHeader(time.Now())
	nonce := generateNonce()
	headersList, signature := crypto.BuildSignature(c.apiSecret, method, path, dateValue, "hmac-sha256", nonce)

	// Đóng gói theo tiêu chuẩn Authorization Signature Header
	sigHeaderValue := fmt.Sprintf(`Signature keyId="%s",algorithm="hmac-sha256",headers="%s",signature="%s",nonce="%s"`,
		c.apiKey, headersList, signature, nonce)

	// 4. Thiết lập các tiêu đề Header bắt buộc
	req.Header.Set("Date", dateValue)
	req.Header.Set("X-Signature", sigHeaderValue)
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("version", c.apiVersion)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 5. Thực thi Request
	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("api returned error status: %d", resp.StatusCode)
	}

	// 6. Decode kết quả trả về
	if resTarget != nil {
		return json.NewDecoder(resp.Body).Decode(resTarget)
	}
	return nil
}
