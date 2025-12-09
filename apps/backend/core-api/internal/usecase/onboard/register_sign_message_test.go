package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnboardUsecase_GetRegisterSignMessageDetails(t *testing.T) {
	tests := []struct {
		name            string
		expectedMessage string
	}{
		{
			name:            "should_return_consistent_register_sign_message",
			expectedMessage: "Please sign this message to prove your ownership of the wallet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			uc := NewOnboardUsecase(nil, nil, nil, nil)

			// Act
			message := uc.GetRegisterSignMessage()

			// Assert
			assert.Equal(t, tt.expectedMessage, message)
			assert.NotEmpty(t, message)
		})
	}
}

func TestOnboardUsecase_GetRegisterSignMessage_Consistency(t *testing.T) {
	t.Run("should_return_same_message_on_multiple_calls", func(t *testing.T) {
		// Arrange
		uc := NewOnboardUsecase(nil, nil, nil, nil)

		// Act
		message1 := uc.GetRegisterSignMessage()
		message2 := uc.GetRegisterSignMessage()
		message3 := uc.GetRegisterSignMessage()

		// Assert
		assert.Equal(t, message1, message2)
		assert.Equal(t, message2, message3)
		assert.NotEmpty(t, message1)
	})
}

func TestOnboardUsecase_Initialization(t *testing.T) {
	t.Run("should_initialize_with_default_register_sign_message", func(t *testing.T) {
		// Act
		uc := NewOnboardUsecase(nil, nil, nil, nil)

		// Assert
		assert.NotNil(t, uc)
		assert.NotEmpty(t, uc.registerSignMessage)
		assert.Equal(t, "Please sign this message to prove your ownership of the wallet", uc.registerSignMessage)
	})

	t.Run("should_initialize_with_nil_dependencies", func(t *testing.T) {
		// Act
		uc := NewOnboardUsecase(nil, nil, nil, nil)

		// Assert
		assert.NotNil(t, uc)
		// Dependencies can be nil for certain operations
		assert.NotEmpty(t, uc.GetRegisterSignMessage())
	})
}







