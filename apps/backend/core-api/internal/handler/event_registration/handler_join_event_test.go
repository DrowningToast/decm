package event_registration

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDecodeSignatureHex verifies the 0x-prefix stripping logic used in JoinEvent.
// The handler does: sigHex := strings.TrimPrefix(*req.Signature, "0x"); hex.DecodeString(sigHex)
// so both prefixed and unprefixed hex strings must decode successfully.
func TestDecodeSignatureHex(t *testing.T) {
	decodeSignature := func(sig string) ([]byte, error) {
		return hex.DecodeString(strings.TrimPrefix(sig, "0x"))
	}

	t.Run("decodes 0x-prefixed hex signature", func(t *testing.T) {
		sig := "0xaabbccdd1234567890abcdef"
		b, err := decodeSignature(sig)
		require.NoError(t, err)
		assert.Equal(t, []byte{0xaa, 0xbb, 0xcc, 0xdd, 0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}, b)
	})

	t.Run("decodes unprefixed hex signature", func(t *testing.T) {
		sig := "aabbccdd1234567890abcdef"
		b, err := decodeSignature(sig)
		require.NoError(t, err)
		assert.Equal(t, []byte{0xaa, 0xbb, 0xcc, 0xdd, 0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}, b)
	})

	t.Run("both formats produce identical bytes", func(t *testing.T) {
		withPrefix, err := decodeSignature("0xdeadbeef")
		require.NoError(t, err)
		withoutPrefix, err := decodeSignature("deadbeef")
		require.NoError(t, err)
		assert.Equal(t, withPrefix, withoutPrefix)
	})

	t.Run("returns error for non-hex characters", func(t *testing.T) {
		_, err := decodeSignature("0xZZZZ")
		require.Error(t, err)
	})

	t.Run("returns error for bare invalid hex without prefix", func(t *testing.T) {
		_, err := decodeSignature("ZZZZ")
		require.Error(t, err)
	})
}
