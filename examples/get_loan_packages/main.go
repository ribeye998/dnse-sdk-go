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

	result, err := client.GetLoanPackages(context.Background(), accountNo, dnse.MarketStock, "VIC")
	if err != nil {
		log.Fatalf("GetLoanPackages: %v", err)
	}

	fmt.Println(string(result))
}
