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

	orderID := "57527"

	executions, err := client.GetExecutions(context.Background(), accountID, orderID, dnse.MarketDerivative, "NORMAL")
	if err != nil {
		log.Fatalf("GetExecutions: %v", err)
	}

	fmt.Println(string(executions))
}
