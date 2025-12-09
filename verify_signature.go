package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Input data
	signatureHex := "0x16cf0698ec4f7a48d5e29509aa6ccff7f414193b435d63a37c7bdd90196df2072388eb2038c71a6c4ccb02c5153d53b9a4b199c5636ddf5fc6969b6bfc45c03e1c"
	walletAddress := "0x7836f1b8B0FDf5Fb86A7617eF167EbeC23aa4e8E"
	signMessage := `{"eventContractAddress":"0x2425f1A838e3ee00d0F3AffA2C4D2Ecb87001ad7","receivers":["0x148d532c97fb3f21940c9f6923ab7b6a7df0489091da9fcfd4925fe05bdc49af"]}`

	fmt.Println("🔍 Verifying Signature")
	fmt.Println("====================")
	fmt.Printf("Wallet Address: %s\n", walletAddress)
	fmt.Printf("Sign Message: %s\n", signMessage)
	fmt.Printf("Signature: %s\n", signatureHex)
	fmt.Println()

	// Step 1: Hash the message (same as HashMessage in the codebase - Keccak256 without Ethereum prefix)
	messageHash := crypto.Keccak256([]byte(signMessage))
	fmt.Printf("Message Hash: 0x%x\n", messageHash)
	fmt.Println()

	// Step 2: Decode the signature
	sigBytes := hexutil.MustDecode(signatureHex)
	if len(sigBytes) != 65 {
		log.Fatalf("❌ Invalid signature length: expected 65 bytes, got %d", len(sigBytes))
	}

	fmt.Printf("Signature bytes length: %d\n", len(sigBytes))
	fmt.Printf("Recovery ID (v): %d\n", sigBytes[64])

	// Step 3: Adjust recovery ID if needed (from Ethereum format 27/28 to go-ethereum format 0/1)
	sigBytesAdjusted := make([]byte, len(sigBytes))
	copy(sigBytesAdjusted, sigBytes)
	if sigBytesAdjusted[64] == 27 || sigBytesAdjusted[64] == 28 {
		sigBytesAdjusted[64] -= 27
		fmt.Printf("Adjusted recovery ID: %d\n", sigBytesAdjusted[64])
	}

	// Step 4: Recover the public key from signature
	publicKey, err := crypto.SigToPub(messageHash, sigBytesAdjusted)
	if err != nil {
		log.Fatalf("❌ Failed to recover public key: %v", err)
	}

	// Step 5: Get address from public key
	recoveredAddress := crypto.PubkeyToAddress(*publicKey)
	fmt.Printf("Recovered Address: %s\n", recoveredAddress.Hex())
	fmt.Println()

	// Step 6: Compare addresses
	expectedAddress := common.HexToAddress(walletAddress)
	addressMatch := recoveredAddress.Hex() == expectedAddress.Hex()

	fmt.Println("====================")
	if addressMatch {
		fmt.Println("✅ SIGNATURE VERIFICATION: SUCCESS")
		fmt.Printf("   Recovered address matches expected wallet address\n")
	} else {
		fmt.Println("❌ SIGNATURE VERIFICATION: FAILED")
		fmt.Printf("   Expected: %s\n", expectedAddress.Hex())
		fmt.Printf("   Recovered: %s\n", recoveredAddress.Hex())
	}
}
