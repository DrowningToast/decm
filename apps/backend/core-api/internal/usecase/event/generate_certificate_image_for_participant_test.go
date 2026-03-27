package event

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCertificateImageForParticipant_Authorization(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	configErr := errors.New("no config")

	t.Run("BYOK wallet user authorized by wallet_address", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		uc := &EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			EventCertificateConfigDg:    mockCertConfigDg,
			logger:                      logger,
		}

		certID := uuid.New()
		eventID := uuid.New()
		walletAddr := "0xabc123"
		cert := &entity.EventCertificate{
			Id:                    certID,
			EventId:               eventID,
			ReceiverWalletAddress: &walletAddr,
		}

		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).Return(nil, configErr)

		_, err := uc.GenerateCertificateImageForParticipant(ctx, certID, uuid.New(), nil, &walletAddr)

		// Authorization passed — error is from config lookup, not 403
		assert.Error(t, err)
		assert.NotContains(t, err.Error(), "does not own this certificate")
	})

	t.Run("User not matching credential, email, or wallet returns 403", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		uc := &EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			logger:                      logger,
		}

		certID := uuid.New()
		otherWallet := "0xother"
		cert := &entity.EventCertificate{
			Id:                    certID,
			EventId:               uuid.New(),
			ReceiverWalletAddress: &otherWallet,
		}

		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		myWallet := "0xmine"
		_, err := uc.GenerateCertificateImageForParticipant(ctx, certID, uuid.New(), nil, &myWallet)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not own this certificate")
	})

	t.Run("User authorized by credential ID still works", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		uc := &EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			EventCertificateConfigDg:    mockCertConfigDg,
			logger:                      logger,
		}

		userID := uuid.New()
		certID := uuid.New()
		eventID := uuid.New()
		cert := &entity.EventCertificate{
			Id:                   certID,
			EventId:              eventID,
			ReceiverCredentialId: &userID,
		}

		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).Return(nil, configErr)

		_, err := uc.GenerateCertificateImageForParticipant(ctx, certID, userID, nil, nil)

		assert.Error(t, err)
		assert.NotContains(t, err.Error(), "does not own this certificate")
	})

	t.Run("Wallet address comparison is case-insensitive", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		uc := &EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			EventCertificateConfigDg:    mockCertConfigDg,
			logger:                      logger,
		}

		certID := uuid.New()
		eventID := uuid.New()
		storedWallet := "0xabc123" // stored lowercase
		cert := &entity.EventCertificate{
			Id:                    certID,
			EventId:               eventID,
			ReceiverWalletAddress: &storedWallet,
		}

		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).Return(nil, configErr)

		upperWallet := "0xABC123" // user sends uppercase
		_, err := uc.GenerateCertificateImageForParticipant(ctx, certID, uuid.New(), nil, &upperWallet)

		assert.Error(t, err)
		assert.NotContains(t, err.Error(), "does not own this certificate")
	})
}
