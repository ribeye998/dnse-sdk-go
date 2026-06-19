# DNSE SDK for Go (dnse-sdk-go)

An open-source Go SDK for connecting and interacting with the DNSE Securities OpenAPI, supporting both REST APIs and WebSocket Streams.

## Key Features

- **REST Client (`dnse.Client`)**: 
  - Authenticate and request/renew a `Trading Token` using Smart OTP or PIN.
  - Query trading sub-accounts list and detailed asset balances.
  - Retrieve security definitions (Ceiling/Floor/Reference prices) and margin loan packages.
  - Calculate buying and selling power (`PPSE`) before placing an order.
  - Place, amend, or cancel trading orders with automatic cryptographic signing (`HMAC-SHA256`).
- **WebSocket Client (`dnse.StreamClient`)**:
  - Automatically handles connection lifecycle, WebSocket authentication, and heartbeats.
  - Subscribe to public real-time market data streams (Quotes/Depth, Trade ticks, Market Indices, OHLC Candlesticks).
  - Subscribe to private account streams (Order status changes, Account updates, Position updates).
- **Security**: Built-in automatic calculation of HMAC-SHA256 signatures required by the DNSE OpenAPI v2 gateway.

## Installation

Go 1.21.0 or higher is required.

Clone or import the module:
```bash
go get github.com/ribeye998/dnse-sdk-go
```

Set up your environment variables by copying the example file:
```bash
cp .env_example .env
```

Fill in your DNSE credentials in the `.env` file:
```env
DNSE_BASE_URL=https://openapi.dnse.com.vn
DNSE_WS_URL=wss://ws-openapi.dnse.com.vn
DNSE_API_KEY=your_api_key
DNSE_API_SECRET=your_api_secret
DNSE_ACCOUNT_ID=your_account_id
```

## Quick Start Examples

### 1. REST Client (Create Trading Token and Place Order)

```go
package main

import (
	"context"
	"fmt"
	"log"

	dnse "github.com/ribeye998/dnse-sdk-go"
)

func main() {
	client := dnse.NewClient("https://openapi.dnse.com.vn", "YOUR_API_KEY", "YOUR_API_SECRET")
	ctx := context.Background()

	// Exchange Smart OTP for a Trading Token
	token, err := client.CreateTradingToken(ctx, "smart_otp", "123456")
	if err != nil {
		log.Fatalf("Failed to obtain trading token: %v", err)
	}
	fmt.Printf("Trading Token: %s\n", token)

	// Place a Stock Order
	orderReq := dnse.OrderRequest{
		AccountNo:     "YOUR_ACCOUNT_ID",
		LoanPackageID: 1769,
		Symbol:        "TCB",
		Side:          dnse.SideBuy, // BUY
		OrderType:     "LO",
		Price:         33000,
		Quantity:      100,
	}

	result, err := client.PlaceOrder(ctx, dnse.MarketStock, "NORMAL", orderReq)
	if err != nil {
		log.Fatalf("PlaceOrder failed: %v", err)
	}
	fmt.Printf("Order Placed: %s\n", string(result))
}
```

### 2. WebSocket Client (Stream Live Quotes & Ticks)

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	dnse "github.com/ribeye998/dnse-sdk-go"
)

func main() {
	stream := dnse.NewStreamClient("wss://ws-openapi.dnse.com.vn", "YOUR_API_KEY", "YOUR_API_SECRET")

	// Register callbacks for quotes and ticks
	stream.OnQuote(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[Quote] Symbol: %s | Data: %v\n", symbol, data)
	})
	stream.OnTick(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[Tick] Symbol: %s | Data: %v\n", symbol, data)
	})

	// Start streaming HPG and FPT quotes and trade ticks on G1 board (JSON encoding)
	err := stream.StartMarketData([]string{"HPG", "FPT"}, true, true, false)
	if err != nil {
		log.Fatalf("Stream start error: %v", err)
	}

	fmt.Println("Streaming live data — press Ctrl+C to stop")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stream.Close()
	fmt.Println("Disconnected.")
}
```

---

## Example Scenarios (`examples/`)

Refer to the `examples/` directory for detailed runnable Go scripts matching the DNSE OpenAPI:

- **Trading APIs**:
  - `get_accounts`: Query the list of trading sub-accounts.
  - `get_balances`: Query asset balances and purchasing power.
  - `get_loan_packages`: Query available margin loan packages.
  - `get_ppse`: Calculate purchasing and selling power.
  - `get_orders`: Query the intraday order book.
  - `get_order_detail`: Query details of a specific order.
  - `get_order_history`: Query historical orders.
  - `get_corporate_action_history`: Query corporate action history.
  - `get_execution_detail`: Query execution/match reports (Derivative only).
  - `get_positions`: Query open positions list.
  - `get_positions_by_id`: Query specific position details.
  - `close_position`: Close a derivative position.
  - `create_trading_token`: Authenticate and request a trading token.
  - `place_order`: Place a new order (equivalent to `post_order`).
  - `cancel_order`: Cancel an unmatched active order.
  - `replace_order`: Amend price/quantity of an unmatched order.

- **Market Data APIs**:
  - `get_secdef`: Retrieve security metadata (Floor/Ceiling/Reference prices).
  - `get_instruments`: Search and list tradeable securities.
  - `get_trades`: Query historical trade logs.
  - `get_latest_trade`: Query the latest trade of a symbol.
  - `get_ohlc`: Query historical technical charts/candlesticks.
  - `get_close_price`: Query the latest close price of a symbol.
  - `get_working_dates`: Retrieve the market working calendar.

- **WebSocket Channels**:
  - `websocket_sec_def`: Subscribe to security definition updates.
  - `websocket_quote`: Subscribe to top-of-book best bid/ask quotes.
  - `websocket_trade`: Subscribe to match trade ticks.
  - `websocket_trade_extra`: Subscribe to detailed ticks with buy/sell volumes.
  - `websocket_ohlc`: Subscribe to real-time open OHLC candlesticks.
  - `websocket_ohlc_closed`: Subscribe to completed closed candlesticks.
  - `websocket_expected_price`: Subscribe to expected prices during auction sessions.
  - `websocket_foreign_investor`: Subscribe to foreign investor flows.
  - `websocket_market_index`: Subscribe to indices (e.g., `VN30`).
