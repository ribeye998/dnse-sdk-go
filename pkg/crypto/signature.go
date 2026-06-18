package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"net/url"
	"strings"
	"time"
)

// FormatDateHeader trả về chuỗi thời gian định dạng GMT/UTC tuân thủ chuẩn RFC1123,
// khớp chính xác với định dạng `%a, %d %b %Y %H:%M:%S +0000` của Python SDK.
func FormatDateHeader(t time.Time) string {
	return t.UTC().Format("Mon, 02 Jan 2006 15:04:05 +0000")
}

// BuildSignature thực hiện tính toán chữ ký HMAC theo đặc tả HTTP Signature của DNSE OpenAPI.
//
// Tham số:
//   - secret: Khóa bí mật (apiSecret) dùng để ký dữ liệu.
//   - method: Phương thức HTTP (GET, POST, PUT, DELETE,...) hoặc action đối với WebSocket.
//   - path: URI Path của endpoint (ví dụ: "/accounts") hoặc URI đối với WebSocket stream.
//   - dateValue: Chuỗi thời gian đã được định dạng qua hàm FormatDateHeader.
//   - algorithm: Thuật toán băm mã hóa hỗ trợ (mặc định: "hmac-sha256", hỗ trợ "hmac-sha384", "hmac-sha512").
//   - nonce: Chuỗi ngẫu nhiên (UUID hex) để chống lặp lại gói tin. Nếu trống "", tham số nonce sẽ không được đưa vào chuỗi ký.
//
// Trả về:
//   - headers: Danh sách các trường tiêu đề tham gia vào chuỗi ký (Ví dụ: "(request-target) date").
//   - signature: Chuỗi chữ ký đã được mã hóa Base64 và mã hóa URL (URL Escaped).
func BuildSignature(secret, method, path, dateValue, algorithm, nonce string) (string, string) {
	headers := "(request-target) date"
	
	// Tạo chuỗi ký cơ sở (signature string) theo chuẩn của DNSE
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("(request-target): %s %s\n", strings.ToLower(method), path))
	sb.WriteString(fmt.Sprintf("date: %s", dateValue))
	
	if nonce != "" {
		sb.WriteString(fmt.Sprintf("\nnonce: %s", nonce))
	}

	// Xác định hàm băm mã hóa dựa trên thuật toán được chỉ định
	var h func() hash.Hash
	switch strings.ToLower(algorithm) {
	case "hmac-sha384":
		h = sha512.New384
	case "hmac-sha512":
		h = sha512.New
	case "hmac-sha256":
		h = sha256.New
	default:
		h = sha256.New // Mặc định sử dụng sha256 giống SDK hệ thống
	}

	// Thực hiện ký HMAC-SHA