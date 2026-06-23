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

	// Register market index callback
	stream.OnMarketIndex(func(data map[string]interface{}) {
		indexName, _ := data["indexName"].(string)
		value, _ := data["valueIndexes"].(float64)

		fmt.Printf("\n=== INDEX UPDATE: %s (Value: %.2f) ===\n", indexName, value)

		// 1. Extract Market Breadth (Counts)
		upCount := data["fluctuationUpIssueCount"]
		downCount := data["fluctuationDownIssueCount"]
		steadyCount := data["fluctuationSteadinessIssueCount"]
		ceilingCount := data["fluctuationUpperLimitIssueCount"]
		floorCount := data["fluctuationLowerLimitIssueCount"]

		fmt.Printf("Breadth counts: Up: %v | Down: %v | Unchanged: %v | Ceiling: %v | Floor: %v\n",
			upCount, downCount, steadyCount, ceilingCount, floorCount)

		// 2. Extract Volume Breadth
		upVol := data["fluctuationUpIssueVolume"]
		downVol := data["fluctuationDownIssueVolume"]
		steadyVol := data["fluctuationSteadinessIssueVolume"]

		fmt.Printf("Volume Breadth: Up Vol: %v | Down Vol: %v | Unchanged Vol: %v\n",
			upVol, downVol, steadyVol)

		// 3. Extract Continuous Matching vs Block Trades (Put-through)
		contVal := data["contauctAccTrdVal"]
		contVol := data["contauctAccTrdVol"]
		blkVal := data["blkTrdAccTrdVal"]
		blkVol := data["blkTrdAccTrdVol"]

		fmt.Printf("Continuous: Vol: %v, Val: %v\n", contVol, contVal)
		fmt.Printf("Block/Put-through: Vol: %v, Val: %v\n", blkVol, blkVal)
	})

	if err := stream.Connect(); err != nil {
		log.Fatalf("Connect: %v", err)
	}
	defer stream.Close()

	// Subscribe to the VN30 index stream (json encoding)
	if err := stream.SubscribeMarketIndex([]string{"VN30"}, "json"); err != nil {
		log.Fatalf("SubscribeMarketIndex: %v", err)
	}

	fmt.Println("Streaming market breadth index data — press Ctrl+C to stop")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
