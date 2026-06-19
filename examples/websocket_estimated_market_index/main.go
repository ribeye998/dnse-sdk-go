package main

import (
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

	stream := dnse.NewStreamClient(cfg.WSURL, cfg.APIKey, cfg.APISecret)

	stream.OnEstimatedMarketIndex(func(data map[string]interface{}) {
		fmt.Printf("[estimated_market_index] %v\n", data)
	})

	if err := stream.Connect(); err != nil {
		log.Fatalf("Connect: %v", err)
	}

	// Subscribe to estimated index feed (e.g. "VN30")
	if err := stream.SubscribeEstimatedMarketIndex([]string{"VN30"}, "json"); err != nil {
		log.Fatalf("SubscribeEstimatedMarketIndex: %v", err)
	}

	fmt.Println("Streaming estimated market index data — press Ctrl+C to stop")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stream.Close()
	fmt.Println("Disconnected.")
}
