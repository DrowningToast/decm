package event

import (
	"testing"

	"apps/backend/core-api/internal/usecase/cyptoutils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

// TestVerifyHostSignature_ImportCertificateReceivers verifies signatures created
// during the import certificate receivers flow, which uses HashMessage() (no Ethereum prefix)
// This matches the actual signing process in import_certificate_receivers.go
func TestVerifyHostSignature_ImportCertificateReceivers(t *testing.T) {
	// Data from user's verification request
	signMessage := `{"eventContractAddress":"0x2425f1A838e3ee00d0F3AffA2C4D2Ecb87001ad7","receivers":["0x148d532c97fb3f21940c9f6923ab7b6a7df0489091da9fcfd4925fe05bdc49af"]}`
	hostSignature := "0x16cf0698ec4f7a48d5e29509aa6ccff7f414193b435d63a37c7bdd90196df2072388eb2038c71a6c4ccb02c5153d53b9a4b199c5636ddf5fc6969b6bfc45c03e1c"
	expectedHostAddress := "0x7836f1b8B0FDf5Fb86A7617eF167EbeC23aa4e8E"

	t.Run("Verify signature from import certificate receivers flow", func(t *testing.T) {
		// Step 1: Hash the message using HashMessage() (no Ethereum prefix)
		// This matches the signing process in import_certificate_receivers.go line 286
		messageHash := cyptoutils.HashMessage(signMessage)
		messageHashCommon := common.BytesToHash(messageHash)

		// Step 2: Decode the signature from hex
		signatureBytes := hexutil.MustDecode(hostSignature)

		// Step 3: Verify using VerifySignatureByDigest (matches the pattern used elsewhere)
		expectedAddress := common.HexToAddress(expectedHostAddress)
		isValid, err := cyptoutils.VerifySignatureByDigest(expectedAddress, messageHashCommon, signatureBytes)
		if err != nil {
			t.Logf("❌ Failed to verify signature: %v", err)
			t.FailNow()
		}

		t.Logf("📝 Sign Message: %s", signMessage)
		t.Logf("🔐 Signature: %s", hostSignature)
		t.Logf("📦 Message Hash: 0x%x", messageHash)
		t.Logf("🎯 Expected Host: %s", expectedHostAddress)
		t.Logf("")

		if isValid {
			t.Logf("✅ SUCCESS: Signature is VALID! Host created this signature.")
			t.Logf("   The signature matches the expected wallet address.")
		} else {
			t.Logf("❌ FAILURE: Signature is INVALID!")
			t.Logf("   Signature does not match the expected wallet address.")
		}

		// Assert signature is valid
		assert.True(t, isValid, "Signature verification failed - signature does not match the expected wallet address")
	})
}
