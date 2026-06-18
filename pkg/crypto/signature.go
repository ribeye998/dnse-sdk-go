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

func FormatDateHeader(t time.Time) string {
	return t.UTC().Format("Mon, 02 Jan 2006 15:04:05 +0000")
}

func BuildSignature(secret, method, path, dateValue, algorithm, nonce string) (string, string) {
	headers := "(request-target) date"
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("(request-target): %s %s\n", strings.ToLower(method), path))
	sb.WriteString(fmt.Sprintf("date: %s", dateValue))

	if nonce != "" {
		sb.WriteString(fmt.Sprintf("\nnonce: %s", nonce))
	}

	var h func() hash.Hash
	switch strings.ToLower(algorithm) {
	case "hmac-sha384":
		h = sha512.New384
	case "hmac-sha512":
		h = sha512.New
	case "hmac-sha256":
		h = sha256.New
	default:
		h = sha256.New
	}

	mac := hmac.New(h, []byte(secret))
	mac.Write([]byte(sb.String()))
	encodedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	escapedSignature := url.QueryEscape(encodedSignature)

	return headers, escapedSignature
}
