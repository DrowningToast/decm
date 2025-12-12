package cyptoutils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"apps/backend/common/encryptutils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAddressFromPrivateKey(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	address, err := GetAddressFromPrivateKey(privateKey)
	require.NoError(t, err)
	assert.NotNil(t, address)

	// Verify the address matches what crypto.PubkeyToAddress would produce
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	expectedAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	assert.Equal(t, expectedAddress, *address)
}

func TestGetPublicKeyFromPrivateKey(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey, err := GetPublicKeyFromPrivateKey(privateKey)
	require.NoError(t, err)
	assert.NotNil(t, publicKey)
	assert.Equal(t, privateKey.Public(), publicKey)
}

func TestRecoverPublicKeyFromSignature(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	messageHash := common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234")
	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	require.NoError(t, err)

	recoveredPublicKey, err := RecoverPublicKeyFromSignature(messageHash, signature)
	require.NoError(t, err)
	assert.NotNil(t, recoveredPublicKey)

	// Verify the recovered public key matches the original
	originalPublicKey := privateKey.Public().(*ecdsa.PublicKey)
	assert.Equal(t, originalPublicKey.X, recoveredPublicKey.X)
	assert.Equal(t, originalPublicKey.Y, recoveredPublicKey.Y)
}

func TestRecoverPublicKeyFromSignature_EthereumFormat(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	messageHash := common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234")
	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	require.NoError(t, err)

	// Convert to Ethereum format (recovery ID 27/28)
	ethereumSig := make([]byte, len(signature))
	copy(ethereumSig, signature)
	ethereumSig[64] += 27

	recoveredPublicKey, err := RecoverPublicKeyFromSignature(messageHash, ethereumSig)
	require.NoError(t, err)
	assert.NotNil(t, recoveredPublicKey)
}

func TestPublicKeyToAddress(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := PublicKeyToAddress(publicKey)

	expectedAddress := crypto.PubkeyToAddress(*publicKey)
	assert.Equal(t, expectedAddress, address)
}

func TestParsePrivateKey(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	parsed, err := ParsePrivateKey(privateKeyHex)
	require.NoError(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, privateKey.D, parsed.D)
}

func TestParsePrivateKey_InvalidHex(t *testing.T) {
	_, err := ParsePrivateKey("invalid-hex")
	assert.Error(t, err)
}

func TestParsePrivateKey_InvalidLength(t *testing.T) {
	_, err := ParsePrivateKey("1234") // Too short
	assert.Error(t, err)
}

func TestParseAddress(t *testing.T) {
	validAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"

	address, err := ParseAddress(validAddress)
	require.NoError(t, err)
	assert.NotNil(t, address)
	assert.Equal(t, common.HexToAddress(validAddress), *address)
}

func TestParseAddress_Invalid(t *testing.T) {
	tests := []string{
		"invalid-address",
		"0x",
		"0x123",
		"not-an-address",
		"",
	}

	for _, addr := range tests {
		t.Run(addr, func(t *testing.T) {
			_, err := ParseAddress(addr)
			assert.Error(t, err)
		})
	}
}

func TestDecryptPrivateKey(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	password := "test-password-123"
	encrypted, err := encryptutils.EncryptAESGCM(privateKeyHex, password)
	require.NoError(t, err)

	decryptedKey, address, err := DecryptPrivateKey(encrypted, password)
	require.NoError(t, err)
	assert.NotNil(t, decryptedKey)
	assert.NotNil(t, address)

	// Verify the decrypted key matches the original
	assert.Equal(t, privateKey.D, decryptedKey.D)

	// Verify the address matches
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	assert.Equal(t, expectedAddress, *address)
}

func TestDecryptPrivateKey_InvalidPassword(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	password := "test-password-123"
	wrongPassword := "wrong-password"
	encrypted, err := encryptutils.EncryptAESGCM(privateKeyHex, password)
	require.NoError(t, err)

	_, _, err = DecryptPrivateKey(encrypted, wrongPassword)
	assert.Error(t, err)
}

func TestDecryptPrivateKey_InvalidCiphertext(t *testing.T) {
	_, _, err := DecryptPrivateKey("invalid-ciphertext", "password")
	assert.Error(t, err)
}

func TestDecodeRevertReason_ErrorString(t *testing.T) {
	// Error(string) selector: 0x08c379a0
	errorSelector := []byte{0x08, 0xc3, 0x79, 0xa0}
	errorMessage := "Custom error message"

	// Build error data: selector + offset (32 bytes) + length (32 bytes) + message
	errData := make([]byte, 4+32+32+len(errorMessage))
	copy(errData[0:4], errorSelector)
	// Offset to data: 0x20 (32 bytes)
	errData[35] = 0x20
	// Length of string
	lengthBytes := make([]byte, 32)
	lengthBytes[31] = byte(len(errorMessage))
	copy(errData[36:68], lengthBytes)
	// String data
	copy(errData[68:], []byte(errorMessage))

	reason := DecodeRevertReason(errData)
	assert.Equal(t, errorMessage, reason)
}

func TestDecodeRevertReason_KnownCustomError(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		expected string
	}{
		{
			name:     "Event__InvalidEventName",
			selector: "8e4a23d6",
			expected: "Event__InvalidEventName()",
		},
		{
			name:     "Event__CannotReduceSeatsCount",
			selector: "4af0bf99",
			expected: "Event__CannotReduceSeatsCount()",
		},
		{
			name:     "Event__SeatsCountReached",
			selector: "85ab58b1",
			expected: "Event__SeatsCountReached()",
		},
		{
			name:     "Event__ParticipantIsNotJoined",
			selector: "d773224f",
			expected: "Event__ParticipantIsNotJoined()",
		},
		{
			name:     "Event__ParticipantIsAlreadyJoined",
			selector: "de3f615d",
			expected: "Event__ParticipantIsAlreadyJoined()",
		},
		{
			name:     "Event__AddressCannotBeZero",
			selector: "d92e233d",
			expected: "Event__AddressCannotBeZero()",
		},
		{
			name:     "Event__InvalidSignature",
			selector: "8baa579f",
			expected: "Event__InvalidSignature()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selectorBytes, err := hex.DecodeString(tt.selector)
			require.NoError(t, err)

			errData := make([]byte, 4)
			copy(errData, selectorBytes)

			reason := DecodeRevertReason(errData)
			assert.Equal(t, tt.expected, reason)
		})
	}
}

func TestDecodeRevertReason_UnknownError(t *testing.T) {
	unknownSelector := []byte{0x12, 0x34, 0x56, 0x78}
	reason := DecodeRevertReason(unknownSelector)
	assert.Contains(t, reason, "unknown error")
	assert.Contains(t, reason, "0x12345678")
}

func TestDecodeRevertReason_TooShort(t *testing.T) {
	shortData := []byte{0x12, 0x34}
	reason := DecodeRevertReason(shortData)
	assert.Equal(t, "insufficient data", reason)
}

func TestDecodeRevertReason_InvalidErrorDataLength(t *testing.T) {
	// Error(string) selector with insufficient data
	errorSelector := []byte{0x08, 0xc3, 0x79, 0xa0}
	errData := make([]byte, 50) // Less than 68 bytes required
	copy(errData[0:4], errorSelector)

	reason := DecodeRevertReason(errData)
	assert.Equal(t, "invalid error data length", reason)
}

func TestDecodeRevertReason_TruncatedErrorData(t *testing.T) {
	// Error(string) selector with truncated message
	errorSelector := []byte{0x08, 0xc3, 0x79, 0xa0}
	errData := make([]byte, 100) // Less than required for full message
	copy(errData[0:4], errorSelector)
	errData[35] = 0x20 // Offset
	lengthBytes := make([]byte, 32)
	lengthBytes[31] = 100 // Length longer than available data
	copy(errData[36:68], lengthBytes)

	reason := DecodeRevertReason(errData)
	assert.Equal(t, "truncated error data", reason)
}






