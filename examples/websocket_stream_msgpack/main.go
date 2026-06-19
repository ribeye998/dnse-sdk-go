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

	stream := dnse.NewStreamClient(cfg.WSURL, cfg.APIKey, cfg.APISecret, dnse.WithMsgPack())
	fmt.Printf("Encoding: %s\n", stream.Encoding())

	stream.OnQuote(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[quote] %s | bid=%v offer=%v\n", symbol, data["bid"], data["offer"])
	})
	stream.OnTick(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[tick]  %s | price=%v qty=%v\n", symbol, data["matchPrice"], data["matchQtty"])
	})

	if err := stream.StartMarketData([]string{"VIC", "VHM"}, true, true, false); err != nil {
		log.Fatalf("StartMarketData: %v", err)
	}
	fmt.Println("Streaming VIC, VHM over MsgPack — press Ctrl+C to stop")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stream.Close()
	fmt.Println("Disconnected.")
}
