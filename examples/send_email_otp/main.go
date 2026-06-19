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

	email := "user@example.com"
	err = client.SendEmailOTP(context.Background(), email)
	if err != nil {
		log.Fatalf("SendEmailOTP: %v", err)
	}

	fmt.Printf("Email OTP sent successfully to %s\n", email)
}
