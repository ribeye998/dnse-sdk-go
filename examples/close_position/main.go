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

	// In production, configure trading token first:
	// client.SetTradingToken("...")

	positionID := "2150974"

	err = client.ClosePosition(context.Background(), accountID, positionID, dnse.MarketDerivative, nil)
	if err != nil {
		log.Fatalf("ClosePosition: %v", err)
	}

	fmt.Println("Position close request sent successfully")
}
