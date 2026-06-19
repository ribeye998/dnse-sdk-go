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

	stream.OnQuote(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[quote] %s: %v\n", symbol, data)
	})
	stream.OnTick(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[tick]  %s: %v\n", symbol, data)
	})

	symbols := []string{"VIC", "VHM", "GAS"}
	if err := stream.StartMarketData(symbols, true, true, false); err != nil {
		log.Fatalf("StartMarketData: %v", err)
	}
	fmt.Printf("Streaming %v — press Ctrl+C to stop\n", symbols)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stream.Close()
	fmt.Println("Disconnected.")
}
