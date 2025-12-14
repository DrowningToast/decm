package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseToken(t *testing.T) {
	accessToken := "test-access-token-12345"

	token, err := ParseToken(accessToken)
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, accessToken, token.AccessToken)
}

func TestParseToken_EmptyToken(t *testing.T) {
	token, err := ParseToken("")
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.Empty(t, token.AccessToken)
}







