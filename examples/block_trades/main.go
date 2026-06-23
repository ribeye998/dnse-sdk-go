package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	dnse "github.com/ribeye998/dnse-sdk-go"
	"github.com/ribeye998/dnse-sdk-go/config"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client := dnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)

	// 1. Fetch historical block/put-through trades using limit counts
	fmt.Println("--- Fetching historical block (put-through) trades ---")
	now := time.Now().Unix()
	twoHoursAgo := time.Now().Add(-2 * time.Hour).Unix()

	params := dnse.TradeHistoryParams{
		BoardID: dnse.BoardT1, // Put-through board (T1)
		From:    twoHoursAgo,
		To:      now,
		Limit:   5, // Limit count of entries returned
	}

	result, err := client.GetTrades(context.Background(), "HDB", params)
	if err != nil {
		log.Printf("GetTrades for block board T1 failed: %v", err)
	} else {
		// Print the raw JSON response
		fmt.Printf("Raw historical trades:\n%s\n\n", string(result))
	}

	// 2. Subscribe to live block (put-through) trades using WebSockets
	fmt.Println("--- Starting WebSocket to stream block (put-through) trades ---")
	stream := dnse.NewStreamClient(cfg.WSURL, cfg.APIKey, cfg.APISecret)

	// Handle trade tick events
	stream.OnTick(func(symbol string, data map[string]interface{}) {
		board, _ := data["boardId"].(string)

		// Parse the quantity
		var rawQty int64
		if q, ok := data["matchQtty"].(float64); ok {
			rawQty = int64(q)
		}

		// 3. Applying the scale quantity multiplier based on the board type
		scaledQty := dnse.ScaleQuantity(rawQty, board)

		fmt.Printf("[TRADE TICK] Symbol: %s | Board: %s | Match Price: %v | Raw Qty: %d | Scaled Shares: %d\n",
			symbol, board, data["matchPrice"], rawQty, scaledQty)
	})

	if err := stream.Connect(); err != nil {
		log.Fatalf("Connect: %v", err)
	}
	defer stream.Close()

	// Subscribe to put-through (block trade) boards for VIC and VHM
	channels := map[string][]string{
		dnse.ChanTicks(dnse.BoardT1, "json"): {"VIC", "VHM"}, // Round-lot Put-through
		dnse.ChanTicks(dnse.BoardT4, "json"): {"VIC", "VHM"}, // Odd-lot Put-through
	}

	if err := stream.Subscribe(channels); err != nil {
		log.Fatalf("Subscribe: %v", err)
	}

	fmt.Println("Streaming block trades. Press Ctrl+C to stop...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down.")
}
