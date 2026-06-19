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

	stream.OnForeign(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[foreign] %s: %v\n", symbol, data)
	})

	if err := stream.Connect(); err != nil {
		log.Fatalf("Connect: %v", err)
	}

	channels := map[string][]string{
		dnse.ChanForeign(dnse.BoardG1, "json"): {"VIC", "VHM"},
	}
	if err := stream.Subscribe(channels); err != nil {
		log.Fatalf("Subscribe: %v", err)
	}

	fmt.Println("Streaming foreign investor flow — press Ctrl+C to stop")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stream.Close()
}
