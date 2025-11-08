package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnboardUsecase_GetRegisterSignMessage(t *testing.T) {
	t.Run("should return consistent register sign message", func(t *testing.T) {
		// Arrange
		uc := NewOnboardUsecase(nil, nil, nil, nil)

		// Act
		message := uc.GetRegisterSignMessage()

		// Assert
		require.NotEmpty(t, message)
		assert.Equal(t, "Please sign this message to prove your ownership of the wallet", message)
	})

	t.Run("should return same message on multiple calls", func(t *testing.T) {
		// Arrange
		uc := NewOnboardUsecase(nil, nil, nil, nil)

		// Act
		message1 := uc.GetRegisterSignMessage()
		message2 := uc.GetRegisterSignMessage()

		// Assert
		assert.Equal(t, message1, message2)
	})
}
