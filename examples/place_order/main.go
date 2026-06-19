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

	accountNo := os.Getenv("DNSE_ACCOUNT_ID")
	if accountNo == "" {
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

	token, err := client.CreateTradingToken(ctx, "smart_otp", pin)
	if err != nil {
		log.Fatalf("CreateTradingToken: %v", err)
	}
	fmt.Printf("Trading token: %s\n", token)

	// Replace Symbol, Price, Quantity, and LoanPackageID with your own values.
	// For SELL orders, LoanPackageID must match the package where your position is held —
	// using the wrong package causes a TRADE_QUANTITY_NOT_ENOUGH error.
	req := dnse.OrderRequest{
		AccountNo:     accountNo,
		Symbol:        "VIC",   // e.g. "VIC" for stock, "VN30F2507" for derivative
		Side:          dnse.SideBuy,
		OrderType:     "LO",
		Price:         45000,
		Quantity:      100,
		LoanPackageID: 0, // set to your loan package ID for margin orders; 0 for cash
	}

	result, err := client.PlaceOrder(ctx, dnse.MarketStock, "NORMAL", req)
	if err != nil {
		log.Fatalf("PlaceOrder: %v", err)
	}

	fmt.Printf("Order placed: %s\n", string(result))
}
