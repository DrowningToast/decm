package main

import (
	"context"
	"fmt"
	"log"

	"apps/backend/core-api/config"
	"apps/backend/core-api/internal/usecase/cyptoutils"
)

func main() {
	cfg := config.LoadConfig()

	// 2. Connect to Ethereum Client
	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		log.Fatalf("Failed to connect to blockchain: %v", err)
	}
	defer client.Close()

	fmt.Printf("Connected to Network: %s (Chain ID: %d)\n", cfg.Blockchain.Network, cfg.Blockchain.ChainID)
	fmt.Printf("RPC URL: %s\n", cfg.Blockchain.RPCURL)
	fmt.Println("-------------------------------------------")

	// 3. Fetch Gas Price Info
	ctx := context.Background()
	gasInfo, err := cyptoutils.GetCurrentGasPrice(ctx, client)
	if err != nil {
		log.Fatalf("Failed to fetch gas price: %v", err)
	}

	// 4. Print Results
	fmt.Printf("Base Fee:           %.4f Gwei\n", gasInfo.BaseFeeGwei)
	fmt.Printf("Priority Fee (Tip): %.4f Gwei\n", gasInfo.MaxPriorityFeePerGasGwei)
	fmt.Printf("Current Total:      %.4f Gwei\n", gasInfo.MaxFeePerGasGwei)
	fmt.Println("-------------------------------------------")

	maxAllowed := cfg.Blockchain.MaxGasPriceGwei
	softCap := cfg.Blockchain.SoftCapGasPriceGwei

	fmt.Printf("Soft Cap Fee:       %.4f Gwei\n", softCap)
	fmt.Printf("Hard Cap Fee:       %.4f Gwei\n", maxAllowed)

	if gasInfo.MaxFeePerGasGwei > maxAllowed {
		fmt.Printf("\n🔴 WARNING: Current gas price (%.2f) EXCEEDS your HARD cap (%.2f)!\n", gasInfo.MaxFeePerGasGwei, maxAllowed)
		fmt.Println("   Transactions sent now will likely stay pending or be rejected by the system.")
	} else if gasInfo.MaxFeePerGasGwei > softCap {
		fmt.Printf("\n🟡 WARNING: Current gas price (%.2f) exceeds your SOFT cap (%.2f).\n", gasInfo.MaxFeePerGasGwei, softCap)
		fmt.Println("   Transactions will still go through, but costs are higher than preferred.")
	} else {
		fmt.Printf("\n🟢 Status: Gas prices are within the allowed range.\n")
		margin := ((maxAllowed - gasInfo.MaxFeePerGasGwei) / maxAllowed) * 100
		fmt.Printf("   You have a %.1f%% safety margin remaining (Hard Cap).\n", margin)
	}
	fmt.Println("")
}
