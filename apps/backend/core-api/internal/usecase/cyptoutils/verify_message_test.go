package cyptoutils

import (
	"bytes"
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestVerifyMessage_WithGivenInputs verifies a raw Keccak256(message) signature
// using a fixed private key, message, and expected signature (as in the example).
func TestVerifyMessage_WithGivenInputs(t *testing.T) {
	// 1) Inputs (declare variables as requested)
	privateKeyHex := "106ca3abf57b49d551bdb637051724e3e3e5c58f3ae7526d52e54e2845cbefca"

	signMessage := `{"eventContractAddress":"0x3a5d67932d4670ff4E6B9C2B039C20F57a04f665","receivers":["0xe6916d285536eb37bbcedb957204d09285ace05ffa5212949e504451e19fd061","0x2dca29106272e7d038cef7d74f8038621c5c31252ed07418a277c05eaf28d8c1","0xee09f7c2890ddc914e48259524c1b0ba53128e10b8eb228c8f21005803b2ddaf"]}`

	// Expected values from reference example
	expectedHashHex := "0x20093951289c4f1bc64aea6cbd970137314de09fbd5badb987e9c826676c3917"

	// expectedSignatureHex :=
	expectedSignatureHex := "0xacc6c8ffd5d45ebc3ea8cd49b57e5d988aa74d23c2e9c41b09dafdd50d45ec2c555267b63126e300720c9c894ffbfe66319330ac2f736a0bedbda984df8bd72001"

	t.Logf("Step 1 - Inputs\n  PrivateKeyHex: %s\n  Message: %q", privateKeyHex, signMessage)

	// 2) Load private key and derive public key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		t.Fatalf("HexToECDSA error: %v", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatalf("failed to cast public key to *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	t.Logf("Step 2 - Public Key\n  PublicKey (uncompressed hex): %s\n  BytesLen: %d", hexutil.Encode(publicKeyBytes), len(publicKeyBytes))

	// 3) Hash the message with Keccak256 (no Ethereum prefix in this test)
	hash := crypto.Keccak256Hash([]byte(signMessage))
	t.Logf("Step 3 - Hash\n  Keccak256(message) hex: %s", hexutil.Encode(hash[:]))
	if hexutil.Encode(hash[:]) != expectedHashHex {
		t.Fatalf("unexpected hash: got %s want %s", hexutil.Encode(hash[:]), expectedHashHex)
	}

	// 4) Sign the hash with the private key
	signature, err := crypto.Sign(hash.Bytes(), privateKey)

	if err != nil {
		t.Fatalf("crypto.Sign error: %v", err)
	}
	gotSignatureHex := hexutil.Encode(signature)
	t.Logf("Step 4 - Signature\n  Signature (R||S||V) hex: %s\n  BytesLen: %d", gotSignatureHex, len(signature))
	if gotSignatureHex != expectedSignatureHex {
		// Some environments may differ; fail hard to ensure reproducibility here
		t.Fatalf("unexpected signature: got %s want %s", gotSignatureHex, expectedSignatureHex)
	}

	// 5) Recover public key and compare with original
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		t.Fatalf("Ecrecover error: %v", err)
	}
	t.Logf("Step 5 - Ecrecover\n  RecoveredPubKey (hex): %s\n  BytesLen: %d", hexutil.Encode(sigPublicKey), len(sigPublicKey))
	if !bytes.Equal(sigPublicKey, publicKeyBytes) {
		t.Fatalf("Ecrecover public key mismatch")
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		t.Fatalf("SigToPub error: %v", err)
	}
	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	t.Logf("Step 5.1 - SigToPub\n  RecoveredPubKey (hex): %s\n  BytesLen: %d", hexutil.Encode(sigPublicKeyBytes), len(sigPublicKeyBytes))
	if !bytes.Equal(sigPublicKeyBytes, publicKeyBytes) {
		t.Fatalf("SigToPub public key mismatch")
	}

	// 6) Direct verification (VerifySignature expects R||S without recovery byte)
	signatureNoRecoverID := signature[:len(signature)-1]
	t.Logf("Step 6 - Verification\n  Signature without recovery (hex): %s\n  BytesLen: %d", hexutil.Encode(signatureNoRecoverID), len(signatureNoRecoverID))
	okVerify := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	t.Logf("  VerifySignature result: %v", okVerify)
	if !okVerify {
		t.Fatalf("VerifySignature returned false, want true")
	}
}
