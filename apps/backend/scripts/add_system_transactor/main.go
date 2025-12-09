package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"apps/backend/contracts/decm"
)

func main() {
	// Configuration from .env
	rpcURL := os.Getenv("BLOCKCHAIN_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth-sepolia.g.alchemy.com/v2/p-1v26f8urtKJAgsL3ry5"
	}

	privateKeyHex := os.Getenv("BLOCKCHAIN_PRIVATE_KEY")
	if privateKeyHex == "" {
		privateKeyHex = "ba03f9ec5d706071ce41ceb8edf9c860d9b77458cbec70fd67163635675c820c"
	}

	decmAccessManagerAddr := common.HexToAddress("0x123453863fC927f7E9d76F0ba06EE7425db5C47B")
	systemTransactor := common.HexToAddress("0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A")
	chainID := int64(11155111) // Sepolia

	fmt.Println("🔧 Adding System Transactor to Allowed Senders...")
	fmt.Printf("DecmAccessManager: %s\n", decmAccessManagerAddr.Hex())
	fmt.Printf("System Transactor: %s\n", systemTransactor.Hex())
	fmt.Printf("Chain ID: %d\n\n", chainID)

	// Connect to blockchain
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to blockchain:", err)
	}

	// Load private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("Failed to load private key:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Signing from: %s\n", fromAddress.Hex())

	// Get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Failed to get gas price:", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	if err != nil {
		log.Fatal("Failed to create transactor:", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	// Instantiate contract
	decmAccessManager, err := decm.NewDecmAccessManager(decmAccessManagerAddr, client)
	if err != nil {
		log.Fatal("Failed to instantiate DecmAccessManager:", err)
	}

	// Check if already allowed
	fmt.Println("\n━━━ Checking Current Status ━━━")
	isAllowed, err := decmAccessManager.CheckIsAllowedMsgSender(nil, systemTransactor)
	if err != nil {
		log.Fatal("Failed to check if allowed:", err)
	}

	if isAllowed {
		fmt.Println("✅ System transactor is ALREADY in allowed senders list")
		fmt.Println("No action needed!")
		return
	}

	fmt.Println("❌ System transactor is NOT in allowed senders list")

	// Check if we have admin role
	isAdmin, err := decmAccessManager.CheckIsAdmin(nil, fromAddress)
	if err != nil {
		log.Fatal("Failed to check admin status:", err)
	}

	if !isAdmin {
		fmt.Printf("\n❌ ERROR: Address %s does NOT have admin role\n", fromAddress.Hex())
		fmt.Println("Cannot add allowed sender without admin role!")
		fmt.Println("\n💡 You need to:")
		fmt.Println("   1. Find the admin private key (whoever deployed DecmAccessManager)")
		fmt.Println("   2. Set BLOCKCHAIN_PRIVATE_KEY to that admin's private key")
		fmt.Println("   3. Run this script again")
		os.Exit(1)
	}

	fmt.Printf("✅ Address %s HAS admin role\n", fromAddress.Hex())

	// Add system transactor to allowed senders
	fmt.Println("\n━━━ Adding System Transactor to Allowed Senders ━━━")
	tx, err := decmAccessManager.AddAllowedMsgSender(auth, systemTransactor)
	if err != nil {
		log.Fatal("Failed to add allowed sender:", err)
	}

	fmt.Printf("✅ Transaction sent: %s\n", tx.Hash().Hex())
	fmt.Println("⏳ Waiting for confirmation...")

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal("Failed to wait for transaction:", err)
	}

	if receipt.Status == 1 {
		fmt.Println("✅ Transaction confirmed!")
		fmt.Printf("   Block: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("   Gas used: %d\n", receipt.GasUsed)
	} else {
		fmt.Println("❌ Transaction failed!")
		os.Exit(1)
	}

	// Verify
	fmt.Println("\n━━━ Verifying ━━━")
	isAllowed, err = decmAccessManager.CheckIsAllowedMsgSender(nil, systemTransactor)
	if err != nil {
		log.Fatal("Failed to verify:", err)
	}

	if isAllowed {
		fmt.Println("✅ SUCCESS! System transactor is now in allowed senders list")
		fmt.Println("\n🎉 Certificate claiming should now work!")
	} else {
		fmt.Println("❌ FAILED! System transactor is still not allowed")
	}
}

