package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// FastUnmarshal giải mã nhanh mảng byte sang cấu trúc dữ liệu đích bằng bộ giải mã luồng (Stream Decoder)
func FastUnmarshal(data []byte, v interface{}) error {
	reader := bytes.NewReader(data)
	if err := json.NewDecoder(reader).Decode(v); err != nil {
		return fmt.Errorf("fast unmarshal failed: %w", err)
	}
	return nil
}

// ParseRawMessage nhận thông điệp thô từ WSMessage và bóc tách trực tiếp vào struct chỉ định
func ParseRawMessage(raw json.RawMessage, target interface{}) error {
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("failed to parse raw payload: %w", err)
	}
	return nil
}