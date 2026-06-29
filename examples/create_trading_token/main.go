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

	// otpType can be "smart_otp" or "email_otp"
	otpType := "smart_otp"
	passcode := "123456"

	token, err := client.CreateTradingToken(context.Background(), otpType, passcode)
	if err != nil {
		log.Fatalf("CreateTradingToken: %v", err)
	}

	fmt.Printf("Trading token created: %s\n", token)
}
