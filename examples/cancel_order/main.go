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

	passcode := os.Getenv("DNSE_PASSCODE")
	if passcode == "" {
		fmt.Print("Enter Smart OTP/Passcode: ")
		_, err := fmt.Scanln(&passcode)
		if err != nil {
			log.Fatalf("read passcode: %v", err)
		}
	}

	ctx := context.Background()

	token, err := client.CreateTradingToken(ctx, "smart_otp", passcode)
	if err != nil {
		log.Fatalf("CreateTradingToken: %v", err)
	}
	fmt.Printf("Trading token: %s\n", token)

	orderID := "123472"

	err = client.CancelOrder(ctx, accountID, orderID, dnse.MarketStock, "NORMAL")
	if err != nil {
		log.Fatalf("CancelOrder: %v", err)
	}

	fmt.Printf("Cancel order request sent for order ID: %s\n", orderID)
}
