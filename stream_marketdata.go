package dnse

import (
	"fmt"
	"log"
)

// Board IDs used in WebSocket channel names.
const (
	BoardG1          = "G1"  // Round-lot, primary HOSE/HNX board (lots of 10)
	BoardG3          = "G3"  // PLO put-through after close (lots of 10)
	BoardG4          = "G4"  // Odd-lot 1–99 shares (absolute shares, no ×10 scaling)
	BoardT1          = "T1"  // Put-through round-lot 09:00–14:45
	BoardT3          = "T3"  // Put-through round-lot 14:45–15:00
	BoardT4          = "T4"  // Put-through odd-lot 09:00–14:45 (absolute shares)
	BoardT6          = "T6"  // Put-through odd-lot 14:45–15:00 (absolute shares)
	BoardAll         = "AL"  // Wildcard: subscribe to all boards at once
)

// OHLC resolution constants for stream and REST.
const (
	Resolution1m  = "1"
	Resolution5m  = "5"
	Resolution15m = "15"
	Resolution1h  = "60"
	Resolution1d  = "D"
)

// chanName builds a channel name in the DNSE format "{type}.{qualifier}.{encoding}".
func chanName(chanType, qualifier, encoding string) string {
	return chanType + "." + qualifier + "." + encoding
}

// ChanTicks returns the trade tick channel name for the given board and encoding.
// Example: ChanTicks(BoardG1, "json") → "tick.G1.json"
func ChanTicks(board, encoding string) string { return chanName("tick", board, encoding) }

// ChanTicksExtra returns the detailed trade channel (buy/sell vol aggregation).
func ChanTicksExtra(board, encoding string) string { return chanName("tick_extra", board, encoding) }

// ChanTopPrice returns the bid/ask order-book (top-of-book) channel name.
// Note: wire key for ask side is "offer", not "ask".
func ChanTopPrice(board, encoding string) string { return chanName("top_price", board, encoding) }

// ChanExpectedPrice returns the expected/indicative price channel name.
func ChanExpectedPrice(board, encoding string) string {
	return chanName("expected_price", board, encoding)
}

// ChanSecurityDefinition returns the security definition stream channel name.
func ChanSecurityDefinition(board, encoding string) string {
	return chanName("security_definition", board, encoding)
}

// ChanForeign returns the foreign investor data channel name.
func ChanForeign(board, encoding string) string { return chanName("foreign", board, encoding) }

// ChanMarketIndex returns the market index channel name.
// index examples: "VN30", "HNX", "EST-VN30".
// IMPORTANT: when subscribing to this channel the symbols list MUST be empty [].
// Passing index names inside symbols silences the feed.
func ChanMarketIndex(index, encoding string) string { return chanName("market_index", index, encoding) }

// ChanOHLC returns the real-time (open) OHLC candlestick channel name.
// resolution examples: "1" (1-minute), "5", "15", "60", "D".
func ChanOHLC(resolution, encoding string) string { return chanName("ohlc", resolution, encoding) }

// ChanOHLCClosed returns the completed (closed) OHLC candlestick channel name.
func ChanOHLCClosed(resolution, encoding string) string {
	return chanName("ohlc_closed", resolution, encoding)
}

// ChanOrder returns the private order-event channel for a market type.
// marketType: "STOCK" or "DERIVATIVE".
func ChanOrder(marketType, encoding string) string { return chanName("order", marketType, encoding) }

// Subscribe sends a single subscribe request covering multiple channels at once.
// channels maps channel name → symbols list. Pass an empty slice for channels
// that do not filter by symbol (market index, private order/position/account channels).
func (s *StreamClient) Subscribe(channels map[string][]string) error {
	list := make([]map[string]interface{}, 0, len(channels))
	for name, symbols := range channels {
		if symbols == nil {
			symbols = []string{}
		}
		list = append(list, map[string]interface{}{
			"name":    name,
			"symbols": symbols,
		})
	}
	return s.writeJSON(map[string]interface{}{
		"action":   "subscribe",
		"channels": list,
	})
}

// SubscribeMarketData subscribes to a single named channel for the given symbols.
func (s *StreamClient) SubscribeMarketData(channelName string, symbols []string) error {
	return s.Subscribe(map[string][]string{channelName: symbols})
}

// SubscribeMarketIndex subscribes to one or more market index channels.
// The index name goes into the channel name itself; symbols must stay empty.
// indices examples: []string{"VN30", "HNX", "EST-VN30"}
func (s *StreamClient) SubscribeMarketIndex(indices []string, encoding string) error {
	channels := make(map[string][]string, len(indices))
	for _, idx := range indices {
		channels[ChanMarketIndex(idx, encoding)] = []string{} // symbols MUST be empty
	}
	return s.Subscribe(channels)
}

// StartMarketData connects and subscribes to the most common market data channels
// for the given symbols on board G1 (primary HOSE/HNX round-lot board) using JSON encoding.
//
//   - includeTicks: trade ticks (ChanTicks)
//   - includeOrderBook: bid/ask depth (ChanTopPrice)
//   - includeOHLC: real-time 1-minute candlesticks (ChanOHLC)
func (s *StreamClient) StartMarketData(symbols []string, includeTicks, includeOrderBook, includeOHLC bool) error {
	if err := s.Connect(); err != nil {
		return fmt.Errorf("dnse: stream connect: %w", err)
	}
	channels := map[string][]string{}
	if includeTicks {
		channels[ChanTicks(BoardG1, "json")] = symbols
	}
	if includeOrderBook {
		channels[ChanTopPrice(BoardG1, "json")] = symbols
	}
	if includeOHLC {
		channels[ChanOHLC(Resolution1m, "json")] = symbols
	}
	if len(channels) > 0 {
		if err := s.Subscribe(channels); err != nil {
			return fmt.Errorf("dnse: subscribe: %w", err)
		}
	}
	log.Printf("[dnse] market data stream active: %v", symbols)
	return nil
}
