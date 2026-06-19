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

	pin := os.Getenv("DNSE_PIN")
	if pin == "" {
		fmt.Print("Enter PIN: ")
		_, err := fmt.Scanln(&pin)
		if err != nil {
			log.Fatalf("read PIN: %v", err)
		}
	}

	ctx := context.Background()

	token, err := client.CreateTradingToken(ctx, "smart_otp", string(pin))
	if err != nil {
		log.Fatalf("CreateTradingToken: %v", err)
	}
	fmt.Printf("Trading token: %s\n", token)

	// Replace with a valid active order ID
	orderID := "15742"

	req := dnse.AmendOrderRequest{
		Price:    33100,
		Quantity: 100,
	}

	result, err := client.AmendOrder(ctx, accountID, orderID, dnse.MarketStock, "NORMAL", req)
	if err != nil {
		log.Fatalf("AmendOrder (Replace): %v", err)
	}

	fmt.Printf("Order amended successfully: %s\n", string(result))
}
