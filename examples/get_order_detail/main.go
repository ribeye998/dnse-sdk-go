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

	// Replace with a valid order ID from your account or query get_orders first
	orderID := "15742"

	order, err := client.GetOrderDetail(context.Background(), accountID, orderID, dnse.MarketStock, "NORMAL")
	if err != nil {
		log.Fatalf("GetOrderDetail: %v", err)
	}

	fmt.Printf("Order Details:\n")
	fmt.Printf("- ID:           %d\n", order.ID)
	fmt.Printf("- Side:         %s\n", order.Side)
	fmt.Printf("- Symbol:       %s\n", order.Symbol)
	fmt.Printf("- Price:        %f\n", order.Price)
	fmt.Printf("- Quantity:     %d\n", order.Quantity)
	fmt.Printf("- Status:       %s\n", order.OrderStatus)
}
