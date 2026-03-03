package certificate_share

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	certificatecontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/certificate_contract"
	"apps/backend/core-api/internal/entity"
	certificate_share_usecase "apps/backend/core-api/internal/usecase/certificate_share"
	"apps/backend/services/auth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type mockCertificateShareDataGateway struct{ mock.Mock }

var _ event_datagateway.CertificateShareDataGateway = (*mockCertificateShareDataGateway)(nil)

func (m *mockCertificateShareDataGateway) CreateCertificateShare(ctx context.Context, params event_datagateway.CreateCertificateShareParameters) (*entity.CertificateShare, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificateShare), args.Error(1)
}

func (m *mockCertificateShareDataGateway) GetCertificateShareByHandle(ctx context.Context, handle string) (*entity.CertificateShare, error) {
	args := m.Called(ctx, handle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificateShare), args.Error(1)
}

type mockEventCertificateDataGateway struct{ mock.Mock }

var _ event_datagateway.EventCertificateDataGateway = (*mockEventCertificateDataGateway)(nil)

func (m *mockEventCertificateDataGateway) GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *mockEventCertificateDataGateway) CreateEventCertificate(_ context.Context, _ event_datagateway.CreateEventCertificateParameters) (*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) GetEventCertificateWithSignature(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*event_datagateway.EventCertificateWithSignature, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) GetEventCertificateByInboxMessageID(_ context.Context, _ uuid.UUID) (*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) GetEventCertificatesByEventID(_ context.Context, _ uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) GetAllEventCertificateIDsByEventID(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) GetClaimedCertificatesByEventID(_ context.Context, _ uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) GetUnclaimedReadyCertificatesByEventID(_ context.Context, _ uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) GetClaimedCertificatesByCredentialID(_ context.Context, _ uuid.UUID, _ *string) ([]*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) GetUnclaimedReadyCertificatesByCredentialID(_ context.Context, _ uuid.UUID, _ *string) ([]*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) UpdateEventCertificate(_ context.Context, _ uuid.UUID, _ event_datagateway.UpdateEventCertificateParameters) (*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) UpdateEventCertificateInboxMessageID(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*entity.EventCertificate, error) {
	return nil, nil
}
func (m *mockEventCertificateDataGateway) DeleteEventCertificate(_ context.Context, _ uuid.UUID) error {
	return nil
}

type mockCertContractFactoryDg struct{ mock.Mock }

var _ certificatecontract_datagateway.CertificateContractFactoryDataGateway = (*mockCertContractFactoryDg)(nil)

func (m *mockCertContractFactoryDg) GetContract(addr common.Address) (certificatecontract_datagateway.CertificateContractDataGateway, error) {
	args := m.Called(addr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(certificatecontract_datagateway.CertificateContractDataGateway), args.Error(1)
}

type mockCertContractDg struct{ mock.Mock }

var _ certificatecontract_datagateway.CertificateContractDataGateway = (*mockCertContractDg)(nil)

func (m *mockCertContractDg) GetTokenData(ctx context.Context, tokenId *big.Int) (*entity.CertificatePayload, error) {
	args := m.Called(ctx, tokenId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificatePayload), args.Error(1)
}
func (m *mockCertContractDg) UsedSignatures(_ context.Context, _ []byte) (bool, error) { return false, nil }
func (m *mockCertContractDg) MintNft(_ context.Context, _ certificatecontract_datagateway.MintNftParams) (*big.Int, error) {
	return nil, nil
}
func (m *mockCertContractDg) RevokeCertificate(_ context.Context, _ *big.Int, _ string, _ []byte) error {
	return nil
}
func (m *mockCertContractDg) FilterCertificateMinted(_ context.Context, _ uint64, _ uint64, _ []common.Address) ([]*certificatecontract_datagateway.CertificateMintedEvent, error) {
	return nil, nil
}
func (m *mockCertContractDg) ParseCertificateMinted(_ types.Log) (*certificatecontract_datagateway.CertificateMintedEvent, error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func injectUser(claims *auth.JwtClaims) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if claims != nil {
			c.Locals("user", claims)
		}
		return c.Next()
	}
}

func buildTestApp(
	claims *auth.JwtClaims,
	mockCertDg *mockEventCertificateDataGateway,
	mockShareDg *mockCertificateShareDataGateway,
	mockFactoryDg *mockCertContractFactoryDg,
) *fiber.App {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := &certificate_share_usecase.CertificateShareUsecase{
		EventCertificateDataGateway:  mockCertDg,
		CertificateShareDg:           mockShareDg,
		CertificateContractFactoryDg: mockFactoryDg,
	}

	h := &Handler{
		CertificateShareUc:    uc,
		AuthenticationService: &auth.AuthService{},
		Logger:                discardLogger,
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: customerror.GetErrFiberHandler(discardLogger),
	})
	if claims != nil {
		app.Use(injectUser(claims))
	}

	// Register routes directly (mirrors routes.go without the auth middleware guard,
	// since we're injecting the user via middleware in tests)
	app.Post("/certificates/:certificate_id/shares", h.CreateCertificateShare)
	app.Get("/certificate-shares/:handle", h.GetCertificateShare)
	app.Post("/certificate-shares/:handle/unlock", h.GetCertificateShareWithPassword)
	app.Get("/certificate-shares/:handle/data", h.GetCertificateShareData)
	app.Post("/certificate-shares/:handle/data/unlock", h.GetCertificateShareDataWithPassword)

	return app
}

func doRequest(app *fiber.App, method, path string, body interface{}) *http.Response {
	var buf *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	return resp
}

// ---------------------------------------------------------------------------
// Tests: CreateCertificateShare
// ---------------------------------------------------------------------------

func TestCreateCertificateShare_Unauthenticated(t *testing.T) {
	certID := uuid.New()
	app := buildTestApp(nil, new(mockEventCertificateDataGateway), new(mockCertificateShareDataGateway), new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificates/"+certID.String()+"/shares", nil)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestCreateCertificateShare_InvalidUUID(t *testing.T) {
	userID := uuid.New()
	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		new(mockEventCertificateDataGateway),
		new(mockCertificateShareDataGateway),
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPost, "/certificates/not-a-uuid/shares", nil)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateCertificateShare_CertificateNotFound(t *testing.T) {
	userID := uuid.New()
	certID := uuid.New()

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(nil, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		mockCertDg,
		new(mockCertificateShareDataGateway),
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPost, "/certificates/"+certID.String()+"/shares", nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockCertDg.AssertExpectations(t)
}

func TestCreateCertificateShare_Forbidden_NotOwner(t *testing.T) {
	ownerID := uuid.New()
	callerID := uuid.New() // different user
	certID := uuid.New()

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(&entity.EventCertificate{
		Id:                   certID,
		ReceiverCredentialId: &ownerID,
		CertificateTokenId:   strPtr("1"),
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: callerID},
		mockCertDg,
		new(mockCertificateShareDataGateway),
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPost, "/certificates/"+certID.String()+"/shares", nil)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockCertDg.AssertExpectations(t)
}

func TestCreateCertificateShare_Success_NoPassword(t *testing.T) {
	userID := uuid.New()
	certID := uuid.New()
	broadcastedAt := time.Now()
	shareID := uuid.New()
	handle := "abc123handle"

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(&entity.EventCertificate{
		Id:                      certID,
		ReceiverCredentialId:    &userID,
		CertificateTokenId:      strPtr("1"),
		EventCertificateAddress: strPtr("0xaddr"),
		BroadcastedAt:           &broadcastedAt,
	}, nil)

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("CreateCertificateShare", mock.Anything, mock.MatchedBy(func(p event_datagateway.CreateCertificateShareParameters) bool {
		return p.EventCertificateId == certID && !p.Active && p.Password == nil
	})).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
		Handle:             handle,
		Active:             false,
		Password:           nil,
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		mockCertDg,
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPost, "/certificates/"+certID.String()+"/shares", nil)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockCertDg.AssertExpectations(t)
	mockShareDg.AssertExpectations(t)
}

func TestCreateCertificateShare_Success_WithPassword(t *testing.T) {
	userID := uuid.New()
	certID := uuid.New()
	broadcastedAt := time.Now()
	shareID := uuid.New()
	handle := "pwprotectedhandle"

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(&entity.EventCertificate{
		Id:                      certID,
		ReceiverCredentialId:    &userID,
		CertificateTokenId:      strPtr("2"),
		EventCertificateAddress: strPtr("0xaddr"),
		BroadcastedAt:           &broadcastedAt,
	}, nil)

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("CreateCertificateShare", mock.Anything, mock.MatchedBy(func(p event_datagateway.CreateCertificateShareParameters) bool {
		// Password should be a non-nil hashed value (Argon2id format)
		return p.EventCertificateId == certID && !p.Active && p.Password != nil
	})).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
		Handle:             handle,
		Active:             false,
		Password:           strPtr("$argon2id$..."),
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		mockCertDg,
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPost, "/certificates/"+certID.String()+"/shares",
		map[string]string{"password": "mysecret"})

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockCertDg.AssertExpectations(t)
	mockShareDg.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: GetCertificateShare
// ---------------------------------------------------------------------------

func TestGetCertificateShare_NotFound(t *testing.T) {
	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "missing-handle").Return(nil, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/missing-handle", nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShare_PasswordLocked(t *testing.T) {
	pw := "secret"
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "locked-handle",
		Password:           &pw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "locked-handle").Return(share, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/locked-handle", nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, string(certificate_share_usecase.CertificateShareViewStatusPasswordLocked), body["status"])
	assert.Nil(t, body["certificate"])
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShare_ValidButPending(t *testing.T) {
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "pending-handle",
		Password:           nil,
	}
	cert := &entity.EventCertificate{
		Id:                      certID,
		CertificateTokenId:      nil,
		EventCertificateAddress: nil,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "pending-handle").Return(share, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(cert, nil)

	app := buildTestApp(nil, mockCertDg, mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/pending-handle", nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, string(certificate_share_usecase.CertificateShareViewStatusValidButPending), body["status"])
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

func TestGetCertificateShare_Ready(t *testing.T) {
	certID := uuid.New()
	broadcastedAt := time.Now()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "ready-handle",
		Password:           nil,
	}
	cert := &entity.EventCertificate{
		Id:                      certID,
		CertificateTokenId:      strPtr("42"),
		EventCertificateAddress: strPtr("0xaddr"),
		BroadcastedAt:           &broadcastedAt,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "ready-handle").Return(share, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(cert, nil)

	app := buildTestApp(nil, mockCertDg, mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/ready-handle", nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, string(certificate_share_usecase.CertificateShareViewStatusReady), body["status"])
	assert.NotNil(t, body["certificate"])
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: GetCertificateShareWithPassword
// ---------------------------------------------------------------------------

func TestGetCertificateShareWithPassword_MissingPassword(t *testing.T) {
	app := buildTestApp(nil, new(mockEventCertificateDataGateway), new(mockCertificateShareDataGateway), new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/some-handle/unlock",
		map[string]string{"password": ""})

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetCertificateShareWithPassword_NotFound(t *testing.T) {
	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "missing").Return(nil, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/missing/unlock",
		map[string]string{"password": "guess"})

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareWithPassword_WrongPassword(t *testing.T) {
	rawPw := "correct"
	hashedPw, err := hashutils.HashPassword(rawPw)
	require.NoError(t, err)
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "pw-handle",
		Password:           &hashedPw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "pw-handle").Return(share, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/pw-handle/unlock",
		map[string]string{"password": "wrong"})

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, string(certificate_share_usecase.CertificateShareViewStatusPasswordLocked), body["status"])
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareWithPassword_CorrectPassword(t *testing.T) {
	rawPw := "correct"
	hashedPw, err := hashutils.HashPassword(rawPw)
	require.NoError(t, err)
	certID := uuid.New()
	broadcastedAt := time.Now()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "pw-handle",
		Password:           &hashedPw,
	}
	cert := &entity.EventCertificate{
		Id:                      certID,
		CertificateTokenId:      strPtr("5"),
		EventCertificateAddress: strPtr("0xaddr"),
		BroadcastedAt:           &broadcastedAt,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "pw-handle").Return(share, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(cert, nil)

	app := buildTestApp(nil, mockCertDg, mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/pw-handle/unlock",
		map[string]string{"password": rawPw})

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, string(certificate_share_usecase.CertificateShareViewStatusReady), body["status"])
	assert.NotNil(t, body["certificate"])
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: GetCertificateShareData
// ---------------------------------------------------------------------------

func TestGetCertificateShareData_NotFound(t *testing.T) {
	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "gone").Return(nil, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/gone/data", nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareData_PasswordProtected(t *testing.T) {
	pw := "secret"
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "locked-data-handle",
		Password:           &pw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "locked-data-handle").Return(share, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/locked-data-handle/data", nil)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareData_NotClaimed(t *testing.T) {
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "unclaimed-handle",
		Password:           nil,
	}
	cert := &entity.EventCertificate{
		Id:                      certID,
		CertificateTokenId:      nil,
		EventCertificateAddress: nil,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "unclaimed-handle").Return(share, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(cert, nil)

	app := buildTestApp(nil, mockCertDg, mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/unclaimed-handle/data", nil)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

func TestGetCertificateShareData_Success(t *testing.T) {
	certID := uuid.New()
	tokenID := "99"
	contractAddr := "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "claimed-handle",
		Password:           nil,
	}
	cert := &entity.EventCertificate{
		Id:                      certID,
		CertificateTokenId:      &tokenID,
		EventCertificateAddress: &contractAddr,
	}
	payload := &entity.CertificatePayload{
		Header: entity.CertificatePayloadHeader{Id: "vc-id"},
		Data:   entity.CertificatePayloadData{EventName: "Test Event"},
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "claimed-handle").Return(share, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(cert, nil)

	mockFactory := new(mockCertContractFactoryDg)
	mockContract := new(mockCertContractDg)
	mockFactory.On("GetContract", common.HexToAddress(contractAddr)).Return(mockContract, nil)
	mockContract.On("GetTokenData", mock.Anything, big.NewInt(99)).Return(payload, nil)

	app := buildTestApp(nil, mockCertDg, mockShareDg, mockFactory)

	resp := doRequest(app, http.MethodGet, "/certificate-shares/claimed-handle/data", nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
	mockFactory.AssertExpectations(t)
	mockContract.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: GetCertificateShareDataWithPassword
// ---------------------------------------------------------------------------

func TestGetCertificateShareDataWithPassword_MissingPassword(t *testing.T) {
	app := buildTestApp(nil, new(mockEventCertificateDataGateway), new(mockCertificateShareDataGateway), new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/handle/data/unlock",
		map[string]string{"password": ""})

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetCertificateShareDataWithPassword_WrongPassword(t *testing.T) {
	rawPw := "right"
	hashedPw, err := hashutils.HashPassword(rawPw)
	require.NoError(t, err)
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "pw-data-handle",
		Password:           &hashedPw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "pw-data-handle").Return(share, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/pw-data-handle/data/unlock",
		map[string]string{"password": "wrong"})

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareDataWithPassword_Success(t *testing.T) {
	rawPw := "right"
	hashedPw, err := hashutils.HashPassword(rawPw)
	require.NoError(t, err)
	certID := uuid.New()
	tokenID := "7"
	contractAddr := "0xcafebabecafebabecafebabecafebabecafebabe"
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "pw-data-handle",
		Password:           &hashedPw,
	}
	cert := &entity.EventCertificate{
		Id:                      certID,
		CertificateTokenId:      &tokenID,
		EventCertificateAddress: &contractAddr,
	}
	payload := &entity.CertificatePayload{
		Header: entity.CertificatePayloadHeader{Id: "vc-2"},
		Data:   entity.CertificatePayloadData{EventName: "Protected Event"},
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "pw-data-handle").Return(share, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(cert, nil)

	mockFactory := new(mockCertContractFactoryDg)
	mockContract := new(mockCertContractDg)
	mockFactory.On("GetContract", common.HexToAddress(contractAddr)).Return(mockContract, nil)
	mockContract.On("GetTokenData", mock.Anything, big.NewInt(7)).Return(payload, nil)

	app := buildTestApp(nil, mockCertDg, mockShareDg, mockFactory)

	resp := doRequest(app, http.MethodPost, "/certificate-shares/pw-data-handle/data/unlock",
		map[string]string{"password": rawPw})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
	mockFactory.AssertExpectations(t)
	mockContract.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func strPtr(s string) *string { return &s }
