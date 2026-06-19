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

	accountNo := os.Getenv("DNSE_ACCOUNT_NO")
	if accountNo == "" {
		log.Fatal("set DNSE_ACCOUNT_NO in your .env file")
	}

	client := dnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)

	bal, err := client.GetBalances(context.Background(), accountNo)
	if err != nil {
		log.Fatalf("GetBalances: %v", err)
	}

	fmt.Printf("Account:          %s\n", bal.AccountNo)
	fmt.Printf("Cash available:   %.2f\n", bal.CashAvailable)
	fmt.Printf("Purchasing power: %.2f\n", bal.PurchasingPower)
}
