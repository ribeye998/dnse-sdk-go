package main

import (
	"context"
	"fmt"
	"log"
	"os"

	dnse "github.com/ribeye998/dnse-sdk-go"
	"github.com/ribeye998/dnse-sdk-go/config"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	accountID := os.Getenv("DNSE_ACCOUNT_ID")
	if accountID == "" {
		log.Fatal("set DNSE_ACCOUNT_ID in your .env file")
	}

	client := dnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)

	positions, err := client.GetPositions(context.Background(), accountID, dnse.MarketStock)
	if err != nil {
		log.Fatalf("GetPositions: %v", err)
	}

	fmt.Printf("Found %d position(s):\n", len(positions))
	for _, pos := range positions {
		fmt.Printf("- %s: CostPrice=%f, MarketPrice=%f, Qty=%d\n",
			pos.Symbol, pos.CostPrice, pos.MarketPrice, pos.OpenQuantity)
	}
}
