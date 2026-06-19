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

	orders, err := client.GetOrders(context.Background(), accountID, dnse.MarketStock, "NORMAL")
	if err != nil {
		log.Fatalf("GetOrders: %v", err)
	}

	fmt.Printf("Found %d order(s) today:\n", len(orders))
	for _, order := range orders {
		fmt.Printf("- ID=%d Side=%s Symbol=%s Price=%f Qty=%d Status=%s\n",
			order.ID, order.Side, order.Symbol, order.Price, order.Quantity, order.OrderStatus)
	}
}
