package main

import (
	"context"
	"fmt"
	"log"

	dnse "github.com/ribeye998/dnse-sdk-go"
	"github.com/ribeye998/dnse-sdk-go/config"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client := dnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)

	// Replace with a valid position ID (from get_positions)
	positionID := "2150974"

	pos, err := client.GetPositionByID(context.Background(), positionID, dnse.MarketStock)
	if err != nil {
		log.Fatalf("GetPositionByID: %v", err)
	}

	fmt.Printf("Position Detail:\n")
	fmt.Printf("- ID:           %d\n", pos.ID)
	fmt.Printf("- Symbol:       %s\n", pos.Symbol)
	fmt.Printf("- Status:       %s\n", pos.Status)
	fmt.Printf("- Cost Price:   %f\n", pos.CostPrice)
	fmt.Printf("- Market Price: %f\n", pos.MarketPrice)
	fmt.Printf("- Open Qty:     %d\n", pos.OpenQuantity)
}
