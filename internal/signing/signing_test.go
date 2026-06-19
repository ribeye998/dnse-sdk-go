package signing_test

import (
	"strings"
	"testing"
	"time"

	"github.com/ribeye998/dnse-sdk-go/internal/signing"
)

func TestFormatHTTPDate(t *testing.T) {
	fixed := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	got := signing.FormatHTTPDate(fixed)
	want := "Mon, 15 Jan 2024 10:30:00 +0000"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildRESTHeaders(t *testing.T) {
	fixed := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	date, sig := signing.BuildRESTHeaders("test-key", "test-secret", "GET", "/accounts", fixed)

	if date == "" {
		t.Error("date must not be empty")
	}
	if !strings.HasPrefix(sig, "Signature ") {
		t.Errorf("unexpected signature header format: %q", sig)
	}
	for _, part := range []string{`keyId="test-key"`, `algorithm="hmac-sha256"`} {
		if !strings.Contains(sig, part) {
			t.Errorf("signature header missing %q: %s", part, sig)
		}
	}
}

func TestComputeWSSignature_Deterministic(t *testing.T) {
	sig1 := signing.ComputeWSSignature("key", "secret", 1700000000, "nonce123")
	sig2 := signing.ComputeWSSignature("key", "secret", 1700000000, "nonce123")
	if sig1 != sig2 {
		t.Error("signature must be deterministic for same inputs")
	}
	if sig1 == "" {
		t.Error("signature must not be empty")
	}
}

func TestBuildWSAuthMessage(t *testing.T) {
	msg := signing.BuildWSAuthMessage("mykey", "mysecret", "tok", "123456")
	if msg["action"] != "auth" {
		t.Errorf("expected action=auth, got %v", msg["action"])
	}
	if msg["trading_token"] != "tok" {
		t.Errorf("expected trading_token=tok, got %v", msg["trading_token"])
	}
	if msg["account_no"] != "123456" {
		t.Errorf("expected account_no=123456, got %v", msg["account_no"])
	}
}

func TestBuildWSAuthMessage_OmitsEmptyOptionals(t *testing.T) {
	msg := signing.BuildWSAuthMessage("k", "s", "", "")
	if _, ok := msg["trading_token"]; ok {
		t.Error("trading_token must be omitted when empty")
	}
	if _, ok := msg["account_no"]; ok {
		t.Error("account_no must be omitted when empty")
	}
}
