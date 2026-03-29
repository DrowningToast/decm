package event

import (
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCertificateImportSignMessage(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()
	currentUser := &auth.JwtClaims{UserId: userID}

	t.Run("fails when receiver has neither email nor wallet address", func(t *testing.T) {
		uc := &EventUsecase{
			cfg: createMockConfig(),
		}

		_, err := uc.GetCertificateImportSignMessage(ctx, eventID, []ImportCertificateReceiversRequest{
			{FirstName: strPtr("John")}, // Missing both email and wallet
		}, currentUser)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "must have exactly one of email or wallet_address")
	})

	t.Run("fails when receiver has both email and wallet address", func(t *testing.T) {
		uc := &EventUsecase{
			cfg: createMockConfig(),
		}

		_, err := uc.GetCertificateImportSignMessage(ctx, eventID, []ImportCertificateReceiversRequest{
			{Email: strPtr("a@b.com"), WalletAddress: strPtr("0x123")}, // Has both
		}, currentUser)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "must have exactly one of email or wallet_address")
	})

	t.Run("fails when user is not a verified organizer", func(t *testing.T) {
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedOrganizer: false,
			}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfig(),
		}

		_, err := uc.GetCertificateImportSignMessage(ctx, eventID, []ImportCertificateReceiversRequest{
			{Email: strPtr("a@b.com")},
		}, currentUser)

		require.Error(t, err)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("fails when event is not found", func(t *testing.T) {
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedOrganizer: true,
			}, nil)

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventID).Return(nil, errors.New("not found"))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			cfg:                        createMockConfig(),
		}

		_, err := uc.GetCertificateImportSignMessage(ctx, eventID, []ImportCertificateReceiversRequest{
			{Email: strPtr("a@b.com")},
		}, currentUser)

		require.Error(t, err)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("fails when user is not event owner", func(t *testing.T) {
		otherOwner := uuid.New()
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedOrganizer: true,
			}, nil)

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventID).Return(&entity.Event{
			Id:                eventID,
			OwnerCredentialId: otherOwner,
		}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			cfg:                        createMockConfig(),
		}

		_, err := uc.GetCertificateImportSignMessage(ctx, eventID, []ImportCertificateReceiversRequest{
			{Email: strPtr("a@b.com")},
		}, currentUser)

		require.Error(t, err)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})
}

func TestValidateReceiverRequests(t *testing.T) {
	t.Run("passes when receiver has only email", func(t *testing.T) {
		err := validateReceiverRequests([]ImportCertificateReceiversRequest{
			{Email: strPtr("a@b.com")},
		})
		assert.NoError(t, err)
	})

	t.Run("passes when receiver has only wallet address", func(t *testing.T) {
		err := validateReceiverRequests([]ImportCertificateReceiversRequest{
			{WalletAddress: strPtr("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")},
		})
		assert.NoError(t, err)
	})

	t.Run("fails when receiver has neither email nor wallet address", func(t *testing.T) {
		err := validateReceiverRequests([]ImportCertificateReceiversRequest{
			{FirstName: strPtr("John")},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must have exactly one of email or wallet_address")
	})

	t.Run("fails when receiver has both email and wallet address", func(t *testing.T) {
		err := validateReceiverRequests([]ImportCertificateReceiversRequest{
			{Email: strPtr("a@b.com"), WalletAddress: strPtr("0x123")},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must have exactly one of email or wallet_address")
	})

	t.Run("fails when wallet address is not a valid hex address", func(t *testing.T) {
		err := validateReceiverRequests([]ImportCertificateReceiversRequest{
			{WalletAddress: strPtr("not-a-valid-address")},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid wallet address format")
	})

	t.Run("reports the correct index in error message", func(t *testing.T) {
		err := validateReceiverRequests([]ImportCertificateReceiversRequest{
			{Email: strPtr("valid@example.com")},
			{FirstName: strPtr("missing-identifier")}, // index 1 is the bad one
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "index 1")
	})
}

func TestBuildReceiverHashes(t *testing.T) {
	t.Run("produces consistent hashes for same input", func(t *testing.T) {
		reqs := []ImportCertificateReceiversRequest{
			{
				FirstName:           strPtr("John"),
				LastName:            strPtr("Doe"),
				AcademicInstitution: strPtr("MIT"),
				CertificateTitle:    strPtr("Best Award"),
				CertificateSubtitle: strPtr("For Excellence"),
			},
		}
		h1 := buildReceiverHashes(reqs)
		h2 := buildReceiverHashes(reqs)
		assert.Equal(t, h1, h2)
		assert.Len(t, h1, 1)
		assert.NotEmpty(t, h1[0])
	})

	t.Run("different receivers produce different hashes", func(t *testing.T) {
		req1 := []ImportCertificateReceiversRequest{{FirstName: strPtr("Alice"), LastName: strPtr("A")}}
		req2 := []ImportCertificateReceiversRequest{{FirstName: strPtr("Bob"), LastName: strPtr("B")}}
		assert.NotEqual(t, buildReceiverHashes(req1), buildReceiverHashes(req2))
	})

	t.Run("nil optional fields produce empty string in CSV", func(t *testing.T) {
		req := []ImportCertificateReceiversRequest{{FirstName: strPtr("Alice"), LastName: strPtr("Smith")}}
		hashes := buildReceiverHashes(req)
		assert.Len(t, hashes, 1)
		assert.NotEmpty(t, hashes[0])
	})
}
