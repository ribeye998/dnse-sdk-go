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

	result, err := client.GetPPSE(context.Background(), accountNo, dnse.MarketStock, "VIC", 45000, 0)
	if err != nil {
		log.Fatalf("GetPPSE: %v", err)
	}

	fmt.Println(string(result))
}
