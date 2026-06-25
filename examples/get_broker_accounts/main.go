package main

import (
	"context"
	"fmt"
	"log"

	dnse "github.com/ribeye998/dnse-sdk-go"
	"github.com/ribeye998/dnse-sdk-go/config"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client := dnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)

	ctx := context.Background()

	// Query broker-managed accounts (typically used for broker/admin integration)
	result, err := client.GetBrokerAccounts(ctx)
	if err != nil {
		log.Fatalf("GetBrokerAccounts: %v", err)
	}

	fmt.Printf("Broker managed accounts:\n%s\n", string(result))
}
