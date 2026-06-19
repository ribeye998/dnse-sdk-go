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

	result, err := client.GetSecurityDefinition(context.Background(), "VIC", "")
	if err != nil {
		log.Fatalf("GetSecurityDefinition: %v", err)
	}

	fmt.Println(string(result))
}
