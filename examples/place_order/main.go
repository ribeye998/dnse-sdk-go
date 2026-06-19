package main

import (
	"context"
	"fmt"
	"log"
	"os"

	dnse "github.com/ribeye998/dnse-sdk-go"
	"github.com/ribeye998/dnse-sdk-go/config"
	"golang.org/x/term"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	accountNo := os.Getenv("DNSE_ACCOUNT_NO")
	if accountNo == "" {
		log.Fatal("set DNSE_ACCOUNT_NO in your .env file")
	}

	client := dnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)

	fmt.Print("Enter PIN: ")
	pin, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		log.Fatalf("read PIN: %v", err)
	}

	ctx := context.Background()

	token, err := client.CreateTradingToken(ctx, "PIN", string(pin))
	if err != nil {
		log.Fatalf("CreateTradingToken: %v", err)
	}
	fmt.Printf("Trading token: %s\n", token)

	req := dnse.OrderRequest{
		AccountNo: accountNo,
		Symbol:    "VIC",
		Side:      dnse.SideBuy,
		OrderType: "LO",
		Price:     45000,
		Quantity:  100,
	}

	result, err := client.PlaceOrder(ctx, dnse.MarketStock, "NORMAL", req)
	if err != nil {
		log.Fatalf("PlaceOrder: %v", err)
	}

	fmt.Printf("Order placed: %s\n", string(result))
}
