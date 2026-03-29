package postgres

import (
	"apps/backend/common/pgmapper"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"context"
	"decm-database/go/generated"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubQuerier embeds the Querier interface so only the methods under test need
// to be implemented; any unimplemented call panics with a clear message.
type stubQuerier struct {
	generated.Querier
	createEventCertificateFn                      func(context.Context, generated.CreateEventCertificateParams) (generated.EventCertificate, error)
	getClaimedCertificatesByCredentialIDFn        func(context.Context, generated.GetClaimedCertificatesByCredentialIDParams) ([]generated.GetClaimedCertificatesByCredentialIDRow, error)
	getUnclaimedReadyCertificatesByCredentialIDFn func(context.Context, generated.GetUnclaimedReadyCertificatesByCredentialIDParams) ([]generated.GetUnclaimedReadyCertificatesByCredentialIDRow, error)
}

func (s *stubQuerier) CreateEventCertificate(ctx context.Context, arg generated.CreateEventCertificateParams) (generated.EventCertificate, error) {
	return s.createEventCertificateFn(ctx, arg)
}

func (s *stubQuerier) GetClaimedCertificatesByCredentialID(ctx context.Context, arg generated.GetClaimedCertificatesByCredentialIDParams) ([]generated.GetClaimedCertificatesByCredentialIDRow, error) {
	return s.getClaimedCertificatesByCredentialIDFn(ctx, arg)
}

func (s *stubQuerier) GetUnclaimedReadyCertificatesByCredentialID(ctx context.Context, arg generated.GetUnclaimedReadyCertificatesByCredentialIDParams) ([]generated.GetUnclaimedReadyCertificatesByCredentialIDRow, error) {
	return s.getUnclaimedReadyCertificatesByCredentialIDFn(ctx, arg)
}

const testEncryptionKey = "test-key-for-unit-tests"

func newTestRepo(q generated.Querier) *Repository {
	return &Repository{queries: q, piiEncryptionKey: testEncryptionKey}
}

// --- CreateEventCertificate normalization tests ---

func TestCreateEventCertificate_NormalizesWalletAddress(t *testing.T) {
	tests := []struct {
		name             string
		inputWallet      *string
		expectedPgWallet pgtype.Text
	}{
		{
			name:             "uppercase wallet is lowercased",
			inputWallet:      strPtrRepo("0xABCDEF"),
			expectedPgWallet: pgtype.Text{String: "0xabcdef", Valid: true},
		},
		{
			name:             "nil wallet stays nil (invalid pgText)",
			inputWallet:      nil,
			expectedPgWallet: pgtype.Text{},
		},
		{
			name:             "empty string wallet becomes nil (invalid pgText)",
			inputWallet:      strPtrRepo(""),
			expectedPgWallet: pgtype.Text{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedParams generated.CreateEventCertificateParams

			stub := &stubQuerier{
				createEventCertificateFn: func(_ context.Context, p generated.CreateEventCertificateParams) (generated.EventCertificate, error) {
					capturedParams = p
					return generated.EventCertificate{
						EventContractAddress: pgtype.Text{String: "0xcontract", Valid: true},
					}, nil
				},
			}

			repo := newTestRepo(stub)
			eventID := uuid.New()

			_, err := repo.CreateEventCertificate(context.Background(), event_datagateway.CreateEventCertificateParameters{
				EventID:               eventID,
				ReceiverWalletAddress: tt.inputWallet,
				EventContractAddress:  "0xcontract",
			})
			require.NoError(t, err)

			assert.Equal(t, tt.expectedPgWallet, capturedParams.ReceiverWalletAddress)
		})
	}
}

func TestCreateEventCertificate_NormalizesEmail(t *testing.T) {
	tests := []struct {
		name          string
		inputEmail    *string
		expectValid   bool
		expectedLower string
	}{
		{
			name:          "uppercase email is encrypted in lowercase",
			inputEmail:    strPtrRepo("USER@EXAMPLE.COM"),
			expectValid:   true,
			expectedLower: "user@example.com",
		},
		{
			name:        "nil email results in invalid pgText",
			inputEmail:  nil,
			expectValid: false,
		},
		{
			name:        "empty email results in invalid pgText",
			inputEmail:  strPtrRepo(""),
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedParams generated.CreateEventCertificateParams

			stub := &stubQuerier{
				createEventCertificateFn: func(_ context.Context, p generated.CreateEventCertificateParams) (generated.EventCertificate, error) {
					capturedParams = p
					return generated.EventCertificate{
						EventContractAddress: pgtype.Text{String: "0xcontract", Valid: true},
					}, nil
				},
			}

			repo := newTestRepo(stub)
			eventID := uuid.New()

			_, err := repo.CreateEventCertificate(context.Background(), event_datagateway.CreateEventCertificateParameters{
				EventID:              eventID,
				ReceiverEmail:        tt.inputEmail,
				EventContractAddress: "0xcontract",
			})
			require.NoError(t, err)

			assert.Equal(t, tt.expectValid, capturedParams.ReceiverEmail.Valid)
			if tt.expectValid {
				decrypted, decErr := pgmapper.DecryptPII(capturedParams.ReceiverEmail.String, testEncryptionKey)
				require.NoError(t, decErr)
				assert.Equal(t, tt.expectedLower, decrypted)
			}
		})
	}
}

// --- GetClaimedCertificatesByCredentialID normalization tests ---

func TestGetClaimedCertificatesByCredentialID_NormalizesWalletAddress(t *testing.T) {
	tests := []struct {
		name             string
		inputWallet      *string
		expectedPgWallet pgtype.Text
	}{
		{
			name:             "uppercase wallet is lowercased",
			inputWallet:      strPtrRepo("0xABCDEF"),
			expectedPgWallet: pgtype.Text{String: "0xabcdef", Valid: true},
		},
		{
			name:             "nil wallet stays nil",
			inputWallet:      nil,
			expectedPgWallet: pgtype.Text{},
		},
		{
			name:             "empty wallet becomes nil",
			inputWallet:      strPtrRepo(""),
			expectedPgWallet: pgtype.Text{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedParams generated.GetClaimedCertificatesByCredentialIDParams

			stub := &stubQuerier{
				getClaimedCertificatesByCredentialIDFn: func(_ context.Context, p generated.GetClaimedCertificatesByCredentialIDParams) ([]generated.GetClaimedCertificatesByCredentialIDRow, error) {
					capturedParams = p
					return nil, nil
				},
			}

			repo := newTestRepo(stub)
			_, err := repo.GetClaimedCertificatesByCredentialID(context.Background(), uuid.New(), nil, tt.inputWallet)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedPgWallet, capturedParams.ReceiverWalletAddress)
		})
	}
}

// --- GetUnclaimedReadyCertificatesByCredentialID normalization tests ---

func TestGetUnclaimedReadyCertificatesByCredentialID_NormalizesWalletAddress(t *testing.T) {
	tests := []struct {
		name             string
		inputWallet      *string
		expectedPgWallet pgtype.Text
	}{
		{
			name:             "uppercase wallet is lowercased",
			inputWallet:      strPtrRepo("0xABCDEF"),
			expectedPgWallet: pgtype.Text{String: "0xabcdef", Valid: true},
		},
		{
			name:             "nil wallet stays nil",
			inputWallet:      nil,
			expectedPgWallet: pgtype.Text{},
		},
		{
			name:             "empty wallet becomes nil",
			inputWallet:      strPtrRepo(""),
			expectedPgWallet: pgtype.Text{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedParams generated.GetUnclaimedReadyCertificatesByCredentialIDParams

			stub := &stubQuerier{
				getUnclaimedReadyCertificatesByCredentialIDFn: func(_ context.Context, p generated.GetUnclaimedReadyCertificatesByCredentialIDParams) ([]generated.GetUnclaimedReadyCertificatesByCredentialIDRow, error) {
					capturedParams = p
					return nil, nil
				},
			}

			repo := newTestRepo(stub)
			_, err := repo.GetUnclaimedReadyCertificatesByCredentialID(context.Background(), uuid.New(), nil, tt.inputWallet)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedPgWallet, capturedParams.ReceiverWalletAddress)
		})
	}
}

func strPtrRepo(s string) *string { return &s }
