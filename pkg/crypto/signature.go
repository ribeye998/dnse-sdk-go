package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FormatHTTPDate định dạng thời gian khớp chính xác yêu cầu của DNSE OpenAPI
func FormatHTTPDate(t time.Time) string {
	return t.UTC().Format("Mon, 02 Jan 2006 15:04:05 +0000")
}

// GenerateNonce sinh chuỗi ngẫu nhiên uuid hex sạch
func GenerateNonce() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// BuildSignatureHeader tính toán chuỗi header X-Signature cho các yêu cầu REST API
func BuildSignatureHeader(apiKey, apiSecret, method, path string, now time.Time) (dateValue, signatureHeader string) {
	dateValue = FormatHTTPDate(now)
	nonce := GenerateNonce()

	method = strings.ToLower(method)
	signingString := fmt.Sprintf("(request-target): %s %s\ndate: %s\nnonce: %s", method, path, dateValue, nonce)

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(signingString))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	escapedSignature := url.QueryEscape(signature)

	headersList := "(request-target) date"
	signatureHeader = fmt.Sprintf(`Signature keyId="%s",algorithm="hmac-sha256",headers="%s",signature="%s",nonce="%s"`,
		apiKey, headersList, escapedSignature, nonce)

	return dateValue, signatureHeader
}

// ComputeWSSignature tính toán chữ ký hex HMAC-SHA256 phục vụ WebSocket Auth.
// Logic: HMAC-SHA256(api_secret, "api_key:timestamp:nonce").hexdigest()
func ComputeWSSignature(apiKey, apiSecret string, timestamp int64, nonce string) string {
	message := fmt.Sprintf("%s:%d:%s", apiKey, timestamp, nonce)
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateWSAuthMessage đóng gói cấu trúc Object JSON phục vụ đăng nhập luồng WebSocket
func CreateWSAuthMessage(apiKey, apiSecret, tradingToken, accountNo string) map[string]interface{} {
	timestamp := time.Now().Unix()
	// Sử dụng microsecond timestamp làm chuỗi nonce, khớp chính xác Python reference
	nonce := fmt.Sprintf("%d", time.Now().UnixNano()/1000)
	signature := ComputeWSSignature(apiKey, apiSecret, timestamp, nonce)

	authMsg := map[string]interface{}{
		"action":    "auth",
		"api_key":   apiKey,
		"signature": signature,
		"timestamp": timestamp,
		"nonce":     nonce,
	}

	// Nếu cấu hình có chứa Token giao dịch hoặc Số tài khoản, đính kèm vào gói tin
	if tradingToken != "" {
		authMsg["trading_token"] = tradingToken
	}
	if accountNo != "" {
		authMsg["account_no"] = accountNo
	}

	return authMsg
}
