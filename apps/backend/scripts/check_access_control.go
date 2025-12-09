package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"apps/backend/contracts/eventaccessmanager"
)

func main() {
	rpcURL := os.Getenv("BLOCKCHAIN_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth-sepolia.g.alchemy.com/v2/p-1v26f8urtKJAgsL3ry5"
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to blockchain:", err)
	}

	eventAccessManagerAddr := common.HexToAddress("0x0EADEa85ee88907c46A0EF6AF2f0C530D55BEBB6")
	recoveredSigner := common.HexToAddress("0x9F6f8ef7c3CD068e96C6643754dbF57A49ee13aB")
	systemTransactor := common.HexToAddress("0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A")

	fmt.Println("🔍 Checking Access Control...")
	fmt.Printf("EventAccessManager: %s\n", eventAccessManagerAddr.Hex())
	fmt.Printf("Recovered Signer: %s\n", recoveredSigner.Hex())
	fmt.Printf("System Transactor: %s\n\n", systemTransactor.Hex())

	accessManager, err := eventaccessmanager.NewEventAccessManager(eventAccessManagerAddr, client)
	if err != nil {
		log.Fatal("Failed to instantiate EventAccessManager:", err)
	}

	// Check 1: Get DECM_ACCESS_MANAGER address
	fmt.Println("━━━ CHECK 1: DECM_ACCESS_MANAGER Address ━━━")
	decmAccessManagerAddr, err := accessManager.DECMACCESSMANAGER(nil)
	if err != nil {
		fmt.Printf("❌ Error getting DECM_ACCESS_MANAGER: %v\n", err)
	} else {
		fmt.Printf("DECM_ACCESS_MANAGER: %s\n", decmAccessManagerAddr.Hex())
		if decmAccessManagerAddr == (common.Address{}) {
			fmt.Println("⚠️  WARNING: DECM_ACCESS_MANAGER is zero address!")
		}
	}

	// Check 2: Is recovered signer a host?
	fmt.Println("\n━━━ CHECK 2: Is Recovered Signer a Host? ━━━")
	isHost, err := accessManager.CheckIsHost(nil, recoveredSigner)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		if isHost {
			fmt.Println("✅ YES - Recovered signer IS a host")
		} else {
			fmt.Println("❌ NO - Recovered signer is NOT a host")
		}
	}

	// Check 3: Is recovered signer host or admin?
	fmt.Println("\n━━━ CHECK 3: Is Recovered Signer Host/Admin? ━━━")
	isHostOrAdmin, err := accessManager.CheckIsHostOrAdmin(nil, recoveredSigner)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		if isHostOrAdmin {
			fmt.Println("✅ YES - Recovered signer IS host or admin")
		} else {
			fmt.Println("❌ NO - Recovered signer is NOT host or admin")
		}
	}

	// Check 4: Is system transactor allowed?
	fmt.Println("\n━━━ CHECK 4: Is System Transactor Allowed? ━━━")
	isAllowed, err := accessManager.CheckIsAllowedMsgSender(nil, systemTransactor)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		if isAllowed {
			fmt.Println("✅ YES - System transactor IS allowed")
		} else {
			fmt.Println("❌ NO - System transactor is NOT allowed")
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("VERDICT:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if isHostOrAdmin || isAllowed {
		fmt.Println("✅ Access control is satisfied")
		if isHostOrAdmin {
			fmt.Println("   - Recovered signer is host/admin")
		}
		if isAllowed {
			fmt.Println("   - System transactor is allowed")
		}
	} else {
		fmt.Println("❌ Access control is NOT satisfied")
		fmt.Println("   - Recovered signer is NOT host/admin")
		fmt.Println("   - System transactor is NOT allowed")
		fmt.Println("\n💡 FIX NEEDED:")
		fmt.Printf("   1. Add %s as host in EventAccessManager\n", recoveredSigner.Hex())
		fmt.Printf("   2. OR add %s to allowed senders in DecmAccessManager\n", systemTransactor.Hex())
	}
}
