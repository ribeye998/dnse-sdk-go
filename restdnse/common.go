package restdnse

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// ComputeSignature tính toán chữ ký HMAC-SHA256 dựa trên timestamp và body thô
// Công thức DNSE: HMAC_SHA256(apiSecret, timestamp + body)
func ComputeSignature(apiSecret string, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(timestamp))
	if len(body) > 0 {
		mac.Write(body)
	}
	return hex.EncodeToString(mac.Sum(nil))
}

// SetupCommonHeaders điền các tham số định danh bắt buộc vào HTTP Header
func SetupCommonHeaders(req *http.Request, apiKey string, timestamp string, signature string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", signature)
}

// GetCurrentTimestampMilli trả về chuỗi thời gian Unix Epoch tính bằng mili-giây
func GetCurrentTimestampMilli() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}