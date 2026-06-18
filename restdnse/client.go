package restdnse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	// Gọi chính xác gói crypto trong dự án của bạn (sửa trading-client-go hoặc dnse-sdk-go tùy thuộc go.mod của bạn)
	"dnse-sdk-go/pkg/crypto"
)

// Client đại diện cho bộ cấu hình REST API kết nối đến DNSE OpenAPI
type Client struct {
	baseURL      string
	apiKey       string
	apiSecret    string
	apiVersion   string
	tradingToken string // Bổ sung trường này để lưu token sau khi lấy thành công
	hc           *http.Client
}

// NewClient khởi tạo đối tượng Client mới phục vụ các cuộc gọi REST API
func NewClient(baseURL, apiKey, apiSecret string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		apiVersion: "2026-05-07",
		hc: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SetTradingToken cho phép gán token giao dịch thủ công hoặc từ hàm gọi API
func (c *Client) SetTradingToken(token string) {
	c.tradingToken = token
}

// sendRequest là hàm lõi xử lý việc ký Signature, thêm Header và thực thi HTTP Request
func (c *Client) sendRequest(ctx context.Context, method, path string, queryParams map[string]string, body interface{}, resTarget interface{}) error {
	// 1. Xây dựng đường dẫn URL đầy đủ kèm Query Parameters (nếu có)
	reqURL := fmt.Sprintf("%s%s", c.baseURL, path)
	if len(queryParams) > 0 {
		var q []string
		for k, v := range queryParams {
			q = append(q, fmt.Sprintf("%s=%s", k, v))
		}
		reqURL = fmt.Sprintf("%s?%s", reqURL, strings.Join(q, "&"))
	}

	// 2. Chuyển đổi dữ liệu Body sang dạng mảng Bytes JSON (đối với POST/PUT)
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("lỗi cấu trúc dữ liệu body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	// 3. Khởi tạo đối tượng HTTP Request kèm ngữ cảnh Context
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("lỗi khởi tạo request: %w", err)
	}

	// 4. Gọi hàm từ file hmac.go đúng của bạn để tự động tính toán Chữ ký tích hợp
	dateValue, sigHeaderValue := crypto.BuildSignatureHeader(c.apiKey, c.apiSecret, method, path, time.Now())

	// 5. Nạp cấu hình các tiêu đề chuẩn spec hệ thống DNSE OpenAPI yêu cầu
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Date", dateValue)
	req.Header.Set("X-Signature", sigHeaderValue)
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("version", c.apiVersion)

	if c.tradingToken != "" {
		req.Header.Set("Trading-Token", c.tradingToken)
	}
	// 6. Thực thi gửi gói tin lên máy chủ DNSE
	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("lỗi kết nối mạng: %w", err)
	}
	defer resp.Body.Close()

	// 7. Đọc toàn bộ nội dung phản hồi (Response Body) từ API gửi về
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("lỗi đọc dữ liệu trả về từ api: %w", err)
	}

	// 8. Nếu Status Code khác 2xx, bóc tách thông tin lỗi trả về trực tiếp cho người dùng
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("api trả về mã lỗi status %d: %s", resp.StatusCode, string(respBytes))
	}

	// 9. Giải mã (Unmarshal) chuỗi JSON trả về sang biến cấu trúc đích (resTarget) được truyền vào
	if resTarget != nil && len(respBytes) > 0 {
		if err := json.Unmarshal(respBytes, resTarget); err != nil {
			return fmt.Errorf("lỗi giải mã cấu trúc json kết quả: %w", err)
		}
	}

	return nil
}

type TradingTokenResponse struct {
	TradingToken string `json:"tradingToken"`
}

// CreateTradingToken gửi OTP/PIN lên hệ thống DNSE để sinh Token xác thực giao dịch
// - otpType: Thường là "PIN" hoặc "OTP"
// - passcode: Mã PIN (6 số) hoặc mã OTP từ SMS/SmartOTP của bạn
func (c *Client) CreateTradingToken(ctx context.Context, otpType, passcode string) (string, error) {
	path := "/registration/trading-token"

	payload := map[string]string{
		"otpType":  otpType,
		"passcode": passcode,
	}

	var result TradingTokenResponse

	// Gọi hàm sendRequest lõi (tự động xử lý ký hmac-sha256 và nhận JSON)
	err := c.sendRequest(ctx, "POST", path, nil, payload, &result)
	if err != nil {
		return "", fmt.Errorf("lấy trading token thất bại: %w", err)
	}

	// Lưu trực tiếp vào trạng thái nội tại của Client để các lệnh REST gọi sau tự động dùng
	c.SetTradingToken(result.TradingToken)

	return result.TradingToken, nil
}

// DNSEOrderRequest cấu trúc payload khớp 100% với Adapter của bạn
type DNSEOrderRequest struct {
	AccountNo     string `json:"accountNo"`
	LoanPackageID int64  `json:"loanPackageId,omitempty"` // Tự động thêm omitempty nếu phái sinh không cần
	OrderType     string `json:"orderType"`
	Price         int64  `json:"price"` // Đã được scale thành int64 (ví dụ: 33100)
	Quantity      int64  `json:"quantity"`
	Side          string `json:"side"`
	Symbol        string `json:"symbol"`
	Market        string `json:"market"`
}

// PlaceOrder nhận thêm tham số marketType để động hóa endpoint
func (c *Client) PlaceOrder(ctx context.Context, marketType string, payload DNSEOrderRequest) (string, error) {
	// Endpoint chuẩn cho cổ phiếu cơ sở cơ bản
	path := "/accounts/orders"

	// Nếu sau này bạn chơi phái sinh (DERIVATIVES), có thể map path khác tại đây
	if marketType == "DERIVATIVES" {
		path = "/derivatives/orders" // Ví dụ theo spec DNSE nếu có
	}

	var result struct {
		OrderID string `json:"orderId"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	err := c.sendRequest(ctx, "POST", path, nil, payload, &result)
	if err != nil {
		return "", err
	}

	// Trả về OrderID kiểu chuỗi chuẩn khớp với mong đợi của Adapter: result.OrderID
	return result.OrderID, nil
}
