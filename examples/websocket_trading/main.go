package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	passcode := os.Getenv("DNSE_PASSCODE")
	if passcode == "" {
		fmt.Print("Enter Smart OTP/Passcode: ")
		_, err := fmt.Scanln(&passcode)
		if err != nil {
			log.Fatalf("read passcode: %v", err)
		}
	}

	// 1. Get Trading Token first via REST Client
	restClient := dnse.NewClient(cfg.BaseURL, cfg.APIKey, cfg.APISecret)
	token, err := restClient.CreateTradingToken(context.Background(), "smart_otp", passcode)
	if err != nil {
		log.Fatalf("CreateTradingToken: %v", err)
	}
	fmt.Printf("Obtained trading token: %s\n", token)

	// 2. Configure WebSocket Stream Client.
	// Optionally pass dnse.WithMsgPack() for binary MsgPack encoding (more compact).
	stream := dnse.NewStreamClient(cfg.WSURL, cfg.APIKey, cfg.APISecret)
	stream.SetAccountNo(accountNo)
	stream.SetTradingToken(token)

	// Register callbacks for private trading events
	stream.OnOrderUpdate(func(data map[string]interface{}) {
		fmt.Printf("[Private Order] Update received: %v\n", data)
	})

	stream.OnPositionUpdate(func(data map[string]interface{}) {
		fmt.Printf("[Private Position] Update received: %v\n", data)
	})

	stream.OnAccountUpdate(func(data map[string]interface{}) {
		fmt.Printf("[Private Account] Update received: %v\n", data)
	})

	// 3. Connect and subscribe to private trading channels
	if err := stream.StartTradingData(dnse.MarketStock); err != nil {
		log.Fatalf("StartTradingData: %v", err)
	}

	fmt.Println("Streaming private trading events — press Ctrl+C to stop")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stream.Close()
	fmt.Println("Disconnected.")
}
