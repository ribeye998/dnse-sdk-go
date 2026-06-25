package main

import (
	"context"
	"fmt"
	"log"
	"time"

	dnse "github.com/ribeye998/dnse-sdk-go"
	"github.com/ribeye998/dnse-sdk-go/config"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client := dnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)

	now := time.Now().Unix()
	fiveMinsAgo := time.Now().Add(-5 * time.Minute).Unix()

	params := dnse.TradeHistoryParams{
		BoardID: "G1",
		From:    fiveMinsAgo,
		To:      now,
		Limit:   10,
	}

	result, err := client.GetForeignTrading(context.Background(), "VIC", params)
	if err != nil {
		log.Fatalf("GetForeignTrading: %v", err)
	}

	fmt.Println(string(result))
}
