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

	params := dnse.OrderHistoryParams{
		From:     "2026-06-01",
		To:       "2026-06-19",
		PageSize: 10,
	}

	history, err := client.GetOrderHistory(context.Background(), accountID, dnse.MarketStock, params)
	if err != nil {
		log.Fatalf("GetOrderHistory: %v", err)
	}

	fmt.Println(string(history))
}
