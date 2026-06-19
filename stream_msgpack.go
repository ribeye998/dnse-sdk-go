package dnse

import (
	"github.com/vmihailenco/msgpack/v5"
)

// decodeMsgpack decodes one MsgPack binary WebSocket frame into a slice of messages.
// A frame may contain a single map object or a batch array (fixarray/array16/array32)
// where each element is a separate market data message.
func decodeMsgpack(data []byte) ([]map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}
	b := data[0]
	// Batch array: fixarray (0x90–0x9f), array16 (0xdc), array32 (0xdd).
	if (b >= 0x90 && b <= 0x9f) || b == 0xdc || b == 0xdd {
		var arr []map[string]interface{}
		if err := msgpack.Unmarshal(data, &arr); err != nil {
			return nil, err
		}
		return arr, nil
	}
	var msg map[string]interface{}
	if err := msgpack.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return []map[string]interface{}{msg}, nil
}
