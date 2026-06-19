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

	stream.OnMarketIndex(func(data map[string]interface{}) {
		fmt.Printf("[market_index] %v\n", data)
	})

	if err := stream.Connect(); err != nil {
		log.Fatalf("Connect: %v", err)
	}

	// Index names example: "VN30", "HNX", "EST-VN30"
	// Note: Subscribing to market_index channels requires specifying the index name 
	// in the channel name prefix, and keeping the symbols payload list empty [].
	if err := stream.SubscribeMarketIndex([]string{"VN30"}, "json"); err != nil {
		log.Fatalf("SubscribeMarketIndex: %v", err)
	}

	fmt.Println("Streaming market index data — press Ctrl+C to stop")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stream.Close()
}
