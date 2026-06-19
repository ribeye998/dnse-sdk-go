package dnse

import (
	"fmt"
	"log"
)

// Private channel names (not parameterised by board/encoding).
const (
	// ChanOrders is the private order-update channel ("orders").
	ChanOrders = "orders"
	// ChanPositions is the private position-update channel ("positions").
	ChanPositions = "positions"
	// ChanAccount is the private account balance update channel ("account").
	ChanAccount = "account"
)

// StartTradingData connects and subscribes to all private trading channels for
// the given market type (order events, orders, positions, account updates).
// Requires SetTradingToken and SetAccountNo to be called before Connect.
func (s *StreamClient) StartTradingData(marketType MarketType) error {
	if err := s.Connect(); err != nil {
		return fmt.Errorf("dnse: stream connect: %w", err)
	}
	channels := map[string][]string{
		ChanOrder(string(marketType), "json"): {},
		ChanOrders:    {},
		ChanPositions: {},
		ChanAccount:   {},
	}
	if err := s.Subscribe(channels); err != nil {
		return fmt.Errorf("dnse: subscribe trading channels: %w", err)
	}
	log.Printf("[dnse] trading data stream active (%s)", marketType)
	return nil
}
