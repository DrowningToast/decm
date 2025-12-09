package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"apps/backend/contracts/certificate"
	"apps/backend/contracts/accessmanager"
)

func main() {
	rpcURL := os.Getenv("BLOCKCHAIN_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth-sepolia.g.alchemy.com/v2/p-1v26f8urtKJAgsL3ry5"
	}

	certificateAddr := common.HexToAddress("0x3B328a7049a374a5D728d575383d38486f1c727c")
	eventAccessManagerAddr := common.HexToAddress("0x0EADEa85ee88907c46A0EF6AF2f0C530D55BEBB6")
	expectedDecmAddr := common.HexToAddress("0x123453863fC927f7E9d76F0ba06EE7425db5C47B")

	fmt.Println("🔍 Verifying Deployed Contracts...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to blockchain:", err)
	}

	// ============================================
	// CHECK 1: EventCertificate Contract
	// ============================================
	fmt.Println("━━━ CHECK 1: EventCertificate Contract ━━━")
	fmt.Printf("Address: %s\n\n", certificateAddr.Hex())

	// Check if contract exists
	code, err := client.CodeAt(context.Background(), certificateAddr, nil)
	if err != nil {
		log.Fatal("Failed to get certificate contract code:", err)
	}
	if len(code) == 0 {
		fmt.Println("❌ CRITICAL: No code at this address!")
		fmt.Println("   Contract was not deployed or address is wrong")
		os.Exit(1)
	}
	fmt.Printf("✅ Contract exists (code length: %d bytes)\n", len(code))

	// Instantiate
	certificateContract, err := certificate.NewEventCertificate(certificateAddr, client)
	if err != nil {
		log.Fatal("Failed to instantiate certificate contract:", err)
	}

	// Get EVENT_ACCESS_MANAGER from contract
	fmt.Println("\n📍 Reading EVENT_ACCESS_MANAGER address...")
	actualEventAccessManager, err := certificateContract.EVENTACCESSMANAGER(nil)
	if err != nil {
		fmt.Printf("❌ ERROR reading EVENT_ACCESS_MANAGER: %v\n", err)
		fmt.Println("   This suggests contract ABI mismatch or wrong contract deployed")
	} else {
		fmt.Printf("   Actual: %s\n", actualEventAccessManager.Hex())
		fmt.Printf("   Expected: %s\n", eventAccessManagerAddr.Hex())
		if actualEventAccessManager.Hex() == eventAccessManagerAddr.Hex() {
			fmt.Println("   ✅ MATCH")
		} else {
			fmt.Println("   ❌ MISMATCH!")
		}
	}

	// ============================================
	// CHECK 2: EventAccessManager Contract
	// ============================================
	fmt.Println("\n━━━ CHECK 2: EventAccessManager Contract ━━━")
	fmt.Printf("Address: %s\n\n", eventAccessManagerAddr.Hex())

	// Check if contract exists
	code, err = client.CodeAt(context.Background(), eventAccessManagerAddr, nil)
	if err != nil {
		log.Fatal("Failed to get access manager contract code:", err)
	}
	if len(code) == 0 {
		fmt.Println("❌ CRITICAL: No code at this address!")
		fmt.Println("   Contract was not deployed or address is wrong")
		os.Exit(1)
	}
	fmt.Printf("✅ Contract exists (code length: %d bytes)\n", len(code))

	// Instantiate
	accessManager, err := accessmanager.NewEventAccessManager(eventAccessManagerAddr, client)
	if err != nil {
		log.Fatal("Failed to instantiate access manager:", err)
	}

	// Get DECM_ACCESS_MANAGER from contract
	fmt.Println("\n📍 Reading DECM_ACCESS_MANAGER address...")
	actualDecmAddr, err := accessManager.DECMACCESSMANAGER(nil)
	if err != nil {
		fmt.Printf("❌ ERROR reading DECM_ACCESS_MANAGER: %v\n", err)
		fmt.Println("   This suggests contract ABI mismatch or wrong contract deployed")
	} else {
		fmt.Printf("   Actual: %s\n", actualDecmAddr.Hex())
		fmt.Printf("   Expected (from .env): %s\n", expectedDecmAddr.Hex())
		if actualDecmAddr.Hex() == expectedDecmAddr.Hex() {
			fmt.Println("   ✅ MATCH")
		} else {
			fmt.Println("   ❌ MISMATCH!")
		}

		// Check if it's zero address
		if actualDecmAddr == (common.Address{}) {
			fmt.Println("   ❌ CRITICAL: DECM_ACCESS_MANAGER is ZERO ADDRESS!")
			fmt.Println("   This will cause all access control checks to fail")
		}
	}

	// ============================================
	// CHECK 3: Test Access Control Functions
	// ============================================
	fmt.Println("\n━━━ CHECK 3: Access Control Functions ━━━")

	recoveredSigner := common.HexToAddress("0x9F6f8ef7c3CD068e96C6643754dbF57A49ee13aB")
	systemTransactor := common.HexToAddress("0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A")

	fmt.Printf("\nTesting with:\n")
	fmt.Printf("  Recovered Signer: %s\n", recoveredSigner.Hex())
	fmt.Printf("  System Transactor: %s\n\n", systemTransactor.Hex())

	// Test checkIsHost
	fmt.Println("Testing: checkIsHost(recoveredSigner)...")
	isHost, err := accessManager.CheckIsHost(nil, recoveredSigner)
	if err != nil {
		fmt.Printf("  ❌ ERROR: %v\n", err)
	} else {
		fmt.Printf("  Result: %v\n", isHost)
	}

	// Test checkIsHostOrAdmin
	fmt.Println("\nTesting: checkIsHostOrAdmin(recoveredSigner)...")
	isHostOrAdmin, err := accessManager.CheckIsHostOrAdmin(nil, recoveredSigner)
	if err != nil {
		fmt.Printf("  ❌ ERROR: %v\n", err)
		fmt.Println("  This is the same error your transaction gets!")
	} else {
		fmt.Printf("  Result: %v\n", isHostOrAdmin)
		if !isHostOrAdmin {
			fmt.Println("  ⚠️  Recovered signer is NOT host or admin")
		}
	}

	// Test checkIsAllowedMsgSender
	fmt.Println("\nTesting: checkIsAllowedMsgSender(systemTransactor)...")
	isAllowed, err := accessManager.CheckIsAllowedMsgSender(nil, systemTransactor)
	if err != nil {
		fmt.Printf("  ❌ ERROR: %v\n", err)
		fmt.Println("  This is the same error your transaction gets!")
	} else {
		fmt.Printf("  Result: %v\n", isAllowed)
		if !isAllowed {
			fmt.Println("  ⚠️  System transactor is NOT in allowed senders")
		}
	}

	// ============================================
	// CHECK 4: Verify DECM_ACCESS_MANAGER Exists
	// ============================================
	if err == nil && actualDecmAddr != (common.Address{}) {
		fmt.Println("\n━━━ CHECK 4: DECM_ACCESS_MANAGER Contract ━━━")
		fmt.Printf("Address: %s\n\n", actualDecmAddr.Hex())

		code, err = client.CodeAt(context.Background(), actualDecmAddr, nil)
		if err != nil {
			fmt.Printf("❌ ERROR checking DECM contract: %v\n", err)
		} else if len(code) == 0 {
			fmt.Println("❌ CRITICAL: DECM_ACCESS_MANAGER has NO CODE!")
			fmt.Println("   The address in EventAccessManager points to nothing")
			fmt.Println("   This is why all access checks are failing")
		} else {
			fmt.Printf("✅ DECM contract exists (code length: %d bytes)\n", len(code))
		}
	}

	// ============================================
	// SUMMARY
	// ============================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("SUMMARY:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err != nil {
		fmt.Println("\n❌ PROBLEM IDENTIFIED:")
		fmt.Println("   Access control function calls are REVERTING")
		fmt.Println("\n🔍 Possible causes:")
		fmt.Println("   1. DECM_ACCESS_MANAGER is zero address or invalid")
		fmt.Println("   2. DECM_ACCESS_MANAGER contract is not deployed")
		fmt.Println("   3. Contract ABI mismatch (wrong version deployed)")
		fmt.Println("\n💡 Next steps:")
		fmt.Println("   - Check DECM_ACCESS_MANAGER address in EventAccessManager")
		fmt.Println("   - Verify DECM_ACCESS_MANAGER contract is deployed")
		fmt.Println("   - May need to redeploy EventAccessManager with correct DECM address")
	} else {
		fmt.Println("\n✅ Contracts are correctly configured")
		fmt.Println("\n📋 Access control status:")
		if isHostOrAdmin {
			fmt.Println("   ✅ Recovered signer IS host/admin (should work)")
		} else {
			fmt.Println("   ❌ Recovered signer is NOT host/admin")
		}
		if isAllowed {
			fmt.Println("   ✅ System transactor IS allowed (should work)")
		} else {
			fmt.Println("   ❌ System transactor is NOT allowed")
		}

		if !isHostOrAdmin && !isAllowed {
			fmt.Println("\n🚨 BOTH checks fail - transaction will revert!")
			fmt.Println("\n💡 Fix: Run the add_system_transactor_to_allowed_senders script")
		} else {
			fmt.Println("\n✅ At least one check passes - transaction should work!")
		}
	}
}

