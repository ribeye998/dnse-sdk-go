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

	to := time.Now().Unix()
	from := time.Now().AddDate(0, 0, -30).Unix()

	result, err := client.GetOHLC(context.Background(), "STOCK", "VIC", "D", from, to)
	if err != nil {
		log.Fatalf("GetOHLC: %v", err)
	}

	fmt.Println(string(result))
}
