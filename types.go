package dnse

// MarketType identifies the market segment.
type MarketType string

const (
	MarketStock      MarketType = "STOCK"
	MarketDerivative MarketType = "DERIVATIVE"
)

// OrderSide is the direction of an order.
type OrderSide string

const (
	SideBuy  OrderSide = "BUY"
	SideSell OrderSide = "SELL"
	SideNB   OrderSide = "NB"
	SideNS   OrderSide = "NS"
)

// Balance holds asset and purchasing power details for an account.
type Balance struct {
	AccountNo       string  `json:"accountNo"`
	CashAvailable   float64 `json:"cashAvailable"`
	PurchasingPower float64 `json:"purchasingPower"`
}

// Order represents a trading order as defined by the DNSE OpenAPI spec.
type Order struct {
	ID               string     `json:"id"`
	Side             OrderSide  `json:"side"`
	AccountNo        string     `json:"accountNo"`
	Symbol           string     `json:"symbol"`
	Price            float64    `json:"price"`
	PriceSecure      float64    `json:"priceSecure"`
	AveragePrice     float64    `json:"averagePrice"`
	Quantity         int64      `json:"quantity"`
	FillQuantity     int64      `json:"fillQuantity"`
	CanceledQuantity int64      `json:"canceledQuantity"`
	LeaveQuantity    int64      `json:"leaveQuantity"`
	OrderType        string     `json:"orderType"`
	OrderStatus      string     `json:"orderStatus"`
	LoanPackageID    int64      `json:"loanPackageId"`
	MarketType       MarketType `json:"marketType"`
	TransDate        string     `json:"transDate"`
	CreatedDate      string     `json:"createdDate"`
	ModifiedDate     string     `json:"modifiedDate"`
}

// Position represents an open portfolio position.
type Position struct {
	ID                 int64      `json:"id"`
	AccountNo          string     `json:"accountNo"`
	Symbol             string     `json:"symbol"`
	Status             string     `json:"status"`
	LoanPackageID      int64      `json:"loanPackageId"`
	Side               OrderSide  `json:"side"`
	AccumulateQuantity int64      `json:"accumulateQuantity"`
	TradeQuantity      int64      `json:"tradeQuantity"`
	ClosedQuantity     int64      `json:"closedQuantity"`
	CostPrice          float64    `json:"costPrice"`
	MarketPrice        float64    `json:"marketPrice"`
	BreakEvenPrice     float64    `json:"breakEvenPrice"`
	OpenQuantity       int64      `json:"openQuantity"`
	OverNightQuantity  int64      `json:"overNightQuantity"`
	AverageClosePrice  float64    `json:"averageClosePrice"`
	MarketType         MarketType `json:"marketType"`
}

// OrderRequest is the payload for placing a new order.
type OrderRequest struct {
	AccountNo     string    `json:"accountNo"`
	Symbol        string    `json:"symbol"`
	Side          OrderSide `json:"side"`
	OrderType     string    `json:"orderType"`
	Price         float64   `json:"price"`
	Quantity      int64     `json:"quantity"`
	LoanPackageID int64     `json:"loanPackageId,omitempty"`
	Market        string    `json:"market,omitempty"`
}

// AmendOrderRequest is the payload for modifying an existing order.
type AmendOrderRequest struct {
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}

// OrderHistoryParams holds optional filters for GetOrderHistory.
type OrderHistoryParams struct {
	From      string
	To        string
	PageSize  int
	PageIndex int
}

// CorporateActionHistoryParams holds optional filters for GetCorporateActionHistory.
type CorporateActionHistoryParams struct {
	Symbol    string
	CAType    string
	CAStatus  string
	PageIndex int
	PageSize  int
}

// TradeHistoryParams holds optional filters for GetTrades and GetQuotes.
type TradeHistoryParams struct {
	BoardID       string
	From          int64
	To            int64
	Limit         int
	Order         string
	NextPageToken string
}

// ScaleQuantity converts a raw wire quantity value to actual shares, applying
// the board-specific multiplier.
//
// Boards G1, G3, T1, T3 report quantity in lots-of-10; all others (G4, T4, T6)
// report absolute share counts. Using a flat ×10 for odd-lot boards inflates
// volume 10× (confirmed bug 2026-06-18).
func ScaleQuantity(qty int64, board string) int64 {
	switch board {
	case BoardG4, BoardT4, BoardT6:
		return qty // odd-lot boards: absolute share count
	default:
		return qty * 10 // round-lot boards: lots of 10
	}
}

// InstrumentParams holds optional filters for GetInstruments.
type InstrumentParams struct {
	Symbol          string
	MarketID        string
	SecurityGroupID string
	IndexName       string
	Limit           int
	Page            int
}
