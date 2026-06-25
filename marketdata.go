package dnse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetSecurityDefinition returns ceiling/floor/reference prices for a symbol.
// boardID is optional; pass an empty string to omit.
func (c *Client) GetSecurityDefinition(ctx context.Context, symbol, boardID string) (json.RawMessage, error) {
	path := fmt.Sprintf("/price/%s/secdef", symbol)
	q := url.Values{}
	if boardID != "" {
		q.Set("boardId", boardID)
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetOHLC returns OHLC candlestick data.
// barType is "STOCK" or "DERIVATIVE"; resolution examples: "1", "5", "D".
// from and to are Unix timestamps.
func (c *Client) GetOHLC(ctx context.Context, barType, symbol, resolution string, from, to int64) (json.RawMessage, error) {
	q := url.Values{
		"type":       {barType},
		"symbol":     {symbol},
		"resolution": {resolution},
		"from":       {strconv.FormatInt(from, 10)},
		"to":         {strconv.FormatInt(to, 10)},
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, "/price/ohlc", q, nil, &result)
}

// GetTrades returns historical trade ticks for a symbol.
// TradeHistoryParams.From and .To must be epoch-seconds (not milliseconds) —
// passing milliseconds returns an empty result without an error.
// BoardID must be specified (e.g. "G1" for round-lot, "G4" for odd-lot); the
// endpoint does not mix boards in a single call.
// The server returns at most ~1000 trades per call; use 5-minute windows on busy symbols.
func (c *Client) GetTrades(ctx context.Context, symbol string, p TradeHistoryParams) (json.RawMessage, error) {
	path := fmt.Sprintf("/price/%s/trades", symbol)
	q := buildTradeQuery(p)
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetLatestTrade returns the most recent trade tick for a symbol.
func (c *Client) GetLatestTrade(ctx context.Context, symbol, boardID string) (json.RawMessage, error) {
	path := fmt.Sprintf("/price/%s/trades/latest", symbol)
	q := url.Values{}
	if boardID != "" {
		q.Set("boardId", boardID)
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetQuotes returns historical bid/ask quote book entries for a symbol.
func (c *Client) GetQuotes(ctx context.Context, symbol string, p TradeHistoryParams) (json.RawMessage, error) {
	path := fmt.Sprintf("/price/%s/quotes", symbol)
	q := buildTradeQuery(p)
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetLatestQuote returns the latest bid/ask quote for a symbol.
func (c *Client) GetLatestQuote(ctx context.Context, symbol, boardID string) (json.RawMessage, error) {
	path := fmt.Sprintf("/price/%s/quotes/latest", symbol)
	q := url.Values{}
	if boardID != "" {
		q.Set("boardId", boardID)
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetClosePrice returns the most recent official close price for a symbol.
func (c *Client) GetClosePrice(ctx context.Context, symbol, boardID string) (json.RawMessage, error) {
	path := fmt.Sprintf("/price/%s/close", symbol)
	q := url.Values{}
	if boardID != "" {
		q.Set("boardId", boardID)
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

// GetInstruments returns available securities and their metadata.
func (c *Client) GetInstruments(ctx context.Context, p InstrumentParams) (json.RawMessage, error) {
	q := url.Values{}
	if p.Symbol != "" {
		q.Set("symbol", p.Symbol)
	}
	if p.MarketID != "" {
		q.Set("marketId", p.MarketID)
	}
	if p.SecurityGroupID != "" {
		q.Set("securityGroupId", p.SecurityGroupID)
	}
	if p.IndexName != "" {
		q.Set("indexName", p.IndexName)
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, "/instruments", q, nil, &result)
}

// GetMarketWorkingDates returns the list of market trading calendar dates.
func (c *Client) GetMarketWorkingDates(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, "/market/working-dates", nil, nil, &result)
}

func buildTradeQuery(p TradeHistoryParams) url.Values {
	q := url.Values{}
	if p.BoardID != "" {
		q.Set("boardId", p.BoardID)
	}
	if p.From > 0 {
		q.Set("from", strconv.FormatInt(p.From, 10))
	}
	if p.To > 0 {
		q.Set("to", strconv.FormatInt(p.To, 10))
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		q.Set("order", p.Order)
	}
	if p.NextPageToken != "" {
		q.Set("nextPageToken", p.NextPageToken)
	}
	return q
}

// GetForeignTrading returns historical foreign trading data for a symbol.
func (c *Client) GetForeignTrading(ctx context.Context, symbol string, p TradeHistoryParams) (json.RawMessage, error) {
	path := fmt.Sprintf("/price/%s/foreign-trading", symbol)
	q := buildTradeQuery(p)
	var result json.RawMessage
	return result, c.sendRequest(ctx, http.MethodGet, path, q, nil, &result)
}

