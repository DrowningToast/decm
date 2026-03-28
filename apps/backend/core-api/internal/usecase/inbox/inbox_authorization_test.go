package inbox

import (
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Tests for isAuthorizedToReadMessage - Case-Insensitive Email Comparison
// =============================================================================

func TestInboxUsecase_isAuthorizedToReadMessage_CaseInsensitiveEmail(t *testing.T) {
	t.Run("should authorize when receiver email matches with different case", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		messageEmail := "User@Example.Com"
		userEmail := "user@example.com"

		message := &entity.InboxMessage{
			ReceiverEmail: &messageEmail,
		}

		user := auth.JwtClaims{
			Email: &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.True(t, result, "Should authorize with case-insensitive email match")
	})

	t.Run("should authorize when emails match in uppercase", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		messageEmail := "TEST@EXAMPLE.COM"
		userEmail := "test@example.com"

		message := &entity.InboxMessage{
			ReceiverEmail: &messageEmail,
		}

		user := auth.JwtClaims{
			Email: &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.True(t, result)
	})

	t.Run("should authorize when emails match in lowercase", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		messageEmail := "admin@company.com"
		userEmail := "ADMIN@COMPANY.COM"

		message := &entity.InboxMessage{
			ReceiverEmail: &messageEmail,
		}

		user := auth.JwtClaims{
			Email: &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.True(t, result)
	})

	t.Run("should authorize when emails match in mixed case", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		messageEmail := "JoHn.DoE@ExAmPlE.CoM"
		userEmail := "john.doe@example.com"

		message := &entity.InboxMessage{
			ReceiverEmail: &messageEmail,
		}

		user := auth.JwtClaims{
			Email: &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.True(t, result)
	})

	t.Run("should not authorize when emails are completely different", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		messageEmail := "alice@example.com"
		userEmail := "bob@example.com"

		message := &entity.InboxMessage{
			ReceiverEmail: &messageEmail,
		}

		user := auth.JwtClaims{
			Email: &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.False(t, result, "Should not authorize with different emails")
	})

	t.Run("should authorize by credential ID regardless of email case", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		credentialID := uuid.New()
		messageEmail := "WRONG@EMAIL.COM"
		userEmail := "correct@email.com"

		message := &entity.InboxMessage{
			ReceiverCredentialId: &credentialID,
			ReceiverEmail:        &messageEmail,
		}

		user := auth.JwtClaims{
			UserId: credentialID,
			Email:  &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.True(t, result, "Should authorize by credential ID even if email doesn't match")
	})

	t.Run("should authorize by wallet address regardless of email case", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		walletAddress := "0x1234567890abcdef"
		messageEmail := "WRONG@EMAIL.COM"
		userEmail := "correct@email.com"

		message := &entity.InboxMessage{
			ReceiverWalletAddress: &walletAddress,
			ReceiverEmail:         &messageEmail,
		}

		user := auth.JwtClaims{
			WalletAddress: walletAddress,
			Email:         &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.True(t, result, "Should authorize by wallet address even if email doesn't match")
	})

	t.Run("should authorize when wallet address stored as lowercase but user has EIP-55 checksum", func(t *testing.T) {
		// Simulates: host imports participant with lowercase address, user registered with EIP-55 checksum
		uc := &InboxUsecase{}
		storedAddress := "0xf466e7ce6b06f9b3071557a790bd45f051c1c60a" // lowercase (as stored)
		userAddress := "0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A"   // EIP-55 checksum (from JWT)

		message := &entity.InboxMessage{
			ReceiverWalletAddress: &storedAddress,
		}

		user := auth.JwtClaims{
			WalletAddress: userAddress,
		}

		result := uc.isAuthorizedToReadMessage(message, user)

		assert.True(t, result, "Should authorize with case-insensitive wallet address match")
	})

	t.Run("should authorize when wallet address stored as EIP-55 but user has lowercase", func(t *testing.T) {
		uc := &InboxUsecase{}
		storedAddress := "0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A" // EIP-55 (as stored)
		userAddress := "0xf466e7ce6b06f9b3071557a790bd45f051c1c60a"   // lowercase (from JWT)

		message := &entity.InboxMessage{
			ReceiverWalletAddress: &storedAddress,
		}

		user := auth.JwtClaims{
			WalletAddress: userAddress,
		}

		result := uc.isAuthorizedToReadMessage(message, user)

		assert.True(t, result, "Should authorize with case-insensitive wallet address match")
	})

	t.Run("should not authorize when wallet addresses are completely different", func(t *testing.T) {
		uc := &InboxUsecase{}
		storedAddress := "0xf466e7ce6b06f9b3071557a790bd45f051c1c60a"
		userAddress := "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"

		message := &entity.InboxMessage{
			ReceiverWalletAddress: &storedAddress,
		}

		user := auth.JwtClaims{
			WalletAddress: userAddress,
		}

		result := uc.isAuthorizedToReadMessage(message, user)

		assert.False(t, result, "Should not authorize with different wallet addresses")
	})

	t.Run("should not authorize when all fields are nil or don't match", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		someEmail := "test@example.com"

		message := &entity.InboxMessage{
			ReceiverCredentialId:  nil,
			ReceiverEmail:         nil,
			ReceiverWalletAddress: nil,
		}

		user := auth.JwtClaims{
			UserId: uuid.New(),
			Email:  &someEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.False(t, result, "Should not authorize when no fields match")
	})

	t.Run("should handle nil user email gracefully", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		messageEmail := "test@example.com"

		message := &entity.InboxMessage{
			ReceiverEmail: &messageEmail,
		}

		user := auth.JwtClaims{
			Email: nil, // User doesn't have email
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.False(t, result, "Should not authorize when user email is nil")
	})

	t.Run("should handle nil message email gracefully", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		userEmail := "test@example.com"

		message := &entity.InboxMessage{
			ReceiverEmail: nil, // Message doesn't have receiver email
		}

		user := auth.JwtClaims{
			Email: &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.False(t, result, "Should not authorize when message email is nil")
	})

	t.Run("should handle email with special characters case-insensitively", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		messageEmail := "User.Name+Tag@Example-Domain.Com"
		userEmail := "user.name+tag@example-domain.com"

		message := &entity.InboxMessage{
			ReceiverEmail: &messageEmail,
		}

		user := auth.JwtClaims{
			Email: &userEmail,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.True(t, result, "Should handle special characters case-insensitively")
	})

	t.Run("should authorize with exact case match", func(t *testing.T) {
		// Arrange
		uc := &InboxUsecase{}
		email := "user@example.com"

		message := &entity.InboxMessage{
			ReceiverEmail: &email,
		}

		user := auth.JwtClaims{
			Email: &email,
		}

		// Act
		result := uc.isAuthorizedToReadMessage(message, user)

		// Assert
		assert.True(t, result, "Should authorize with exact case match")
	})
}
