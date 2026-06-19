package signing

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

// FormatHTTPDate formats a time value to the HTTP-date format required by DNSE OpenAPI.
func FormatHTTPDate(t time.Time) string {
	return t.UTC().Format("Mon, 02 Jan 2006 15:04:05 +0000")
}

func generateNonce() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// BuildRESTHeaders computes the Date and X-Signature headers for a REST request.
func BuildRESTHeaders(apiKey, apiSecret, method, path string, now time.Time) (dateValue, signatureHeader string) {
	dateValue = FormatHTTPDate(now)
	nonce := generateNonce()
	method = strings.ToLower(method)

	signingString := fmt.Sprintf("(request-target): %s %s\ndate: %s\nnonce: %s", method, path, dateValue, nonce)

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(signingString))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	signatureHeader = fmt.Sprintf(
		`Signature keyId="%s",algorithm="hmac-sha256",headers="(request-target) date",signature="%s",nonce="%s"`,
		apiKey, url.QueryEscape(signature), nonce,
	)
	return
}

// ComputeWSSignature computes the HMAC-SHA256 hex signature for WebSocket authentication.
// Format: HMAC-SHA256(apiSecret, "apiKey:timestamp:nonce")
func ComputeWSSignature(apiKey, apiSecret string, timestamp int64, nonce string) string {
	msg := fmt.Sprintf("%s:%d:%s", apiKey, timestamp, nonce)
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}

// BuildWSAuthMessage builds the application-level auth payload for WebSocket streams.
func BuildWSAuthMessage(apiKey, apiSecret, tradingToken, accountNo string) map[string]interface{} {
	timestamp := time.Now().Unix()
	nonce := fmt.Sprintf("%d", time.Now().UnixNano()/1000)
	sig := ComputeWSSignature(apiKey, apiSecret, timestamp, nonce)

	msg := map[string]interface{}{
		"action":    "auth",
		"api_key":   apiKey,
		"signature": sig,
		"timestamp": timestamp,
		"nonce":     nonce,
	}
	if tradingToken != "" {
		msg["trading_token"] = tradingToken
	}
	if accountNo != "" {
		msg["account_no"] = accountNo
	}
	return msg
}
