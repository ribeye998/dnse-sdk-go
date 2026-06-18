package restdnse

import (
	"context"
	"fmt"
	"net/http"
)

// GetSecurityDefinition lấy thông tin trần/sàn/tham chiếu của mã (client.get_security_definition)
func (c *Client) GetSecurityDefinition(ctx context.Context, symbol, boardID string, resTarget interface{}) error {
	path := fmt.Sprintf("/price/%s/secdef", symbol)
	var query map[string]string
	if boardID != "" {
		query = map[string]string{"boardId": boardID}
	}
	return c.sendRequest(ctx, http.MethodGet, path, query, nil, resTarget)
}

// GetOHLC tải dữ liệu nến lịch sử đồ thị kỹ thuật (client.get_ohlc)
func (c *Client) GetOHLC(ctx context.Context, barType string, queryParams map[string]string, resTarget interface{}) error {
	path := "/price/ohlc"
	// Đóng gói barType (STOCK/DERIVATIVE) vào query giống Python SDK
	if queryParams == nil {
		queryParams = make(map[string]string)
	}
	queryParams["type"] = barType
	return c.sendRequest(ctx, http.MethodGet, path, queryParams, nil, resTarget)
}
