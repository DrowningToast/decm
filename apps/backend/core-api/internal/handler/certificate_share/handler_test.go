package certificate_share_handler

import (
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
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

	goccyjson "github.com/goccy/go-json"

	customerror "apps/backend/common/customerror"

	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	certificatecontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/certificate_contract"

	certificate_share_usecase "apps/backend/core-api/internal/usecase/certificate_share"

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

func (m *mockCertificateShareDataGateway) GetCertificateShareByID(ctx context.Context, id uuid.UUID) (*entity.CertificateShare, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificateShare), args.Error(1)
}

func (m *mockCertificateShareDataGateway) GetCertificateShareByEventCertificateID(ctx context.Context, eventCertificateID uuid.UUID) (*entity.CertificateShare, error) {
	args := m.Called(ctx, eventCertificateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CertificateShare), args.Error(1)
}

func (m *mockCertificateShareDataGateway) UpdateCertificateShare(ctx context.Context, id uuid.UUID, params event_datagateway.UpdateCertificateShareParameters) (*entity.CertificateShare, error) {
	args := m.Called(ctx, id, params)
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

type mockCertificateImageGenerator struct{ mock.Mock }

var _ certificate_share_usecase.CertificateImageGenerator = (*mockCertificateImageGenerator)(nil)

func (m *mockCertificateImageGenerator) GenerateCertificateImage(ctx context.Context, certificateID uuid.UUID) ([]byte, error) {
	args := m.Called(ctx, certificateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
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

func (m *mockCertContractDg) UsedSignatures(_ context.Context, _ []byte) (bool, error) {
	return false, nil
}

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
	imageGen ...certificate_share_usecase.CertificateImageGenerator,
) *fiber.App {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	var gen certificate_share_usecase.CertificateImageGenerator
	if len(imageGen) > 0 {
		gen = imageGen[0]
	}

	uc := &certificate_share_usecase.CertificateShareUsecase{
		EventCertificateDataGateway:  mockCertDg,
		CertificateShareDg:           mockShareDg,
		CertificateContractFactoryDg: mockFactoryDg,
		CertificateImageGenerator:    gen,
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
	app.Post("/certificate-shares/config/:certificate_id", h.CreateCertificateShare)
	app.Patch("/certificate-shares/config/:share_id", h.UpdateCertificateShare)
	app.Post("/certificate-shares/:handle", h.GetCertificateShareData)
	app.Get("/certificate-shares/:handle/image", h.GetCertificateShareImage)

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

	resp := doRequest(app, http.MethodPost, "/certificate-shares/config/"+certID.String(), nil)

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

	resp := doRequest(app, http.MethodPost, "/certificate-shares/config/not-a-uuid", nil)

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

	resp := doRequest(app, http.MethodPost, "/certificate-shares/config/"+certID.String(), nil)

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

	resp := doRequest(app, http.MethodPost, "/certificate-shares/config/"+certID.String(), nil)

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
	mockShareDg.On("GetCertificateShareByEventCertificateID", mock.Anything, certID).Return(nil, nil)
	mockShareDg.On("CreateCertificateShare", mock.Anything, mock.MatchedBy(func(p event_datagateway.CreateCertificateShareParameters) bool {
		return p.EventCertificateId == certID && p.Active && p.Password == nil
	})).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
		Handle:             handle,
		Active:             true,
		Password:           nil,
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		mockCertDg,
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPost, "/certificate-shares/config/"+certID.String(), nil)

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
	mockShareDg.On("GetCertificateShareByEventCertificateID", mock.Anything, certID).Return(nil, nil)
	mockShareDg.On("CreateCertificateShare", mock.Anything, mock.MatchedBy(func(p event_datagateway.CreateCertificateShareParameters) bool {
		// Password should be a non-nil hashed value (Argon2id format)
		return p.EventCertificateId == certID && p.Active && p.Password != nil
	})).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
		Handle:             handle,
		Active:             true,
		Password:           strPtr("$argon2id$..."),
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		mockCertDg,
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPost, "/certificate-shares/config/"+certID.String(),
		map[string]string{"password": "mysecret"})

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockCertDg.AssertExpectations(t)
	mockShareDg.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: GetCertificateShareData
// ---------------------------------------------------------------------------

func TestGetCertificateShareData_NotFound(t *testing.T) {
	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "gone").Return(nil, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/gone", nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareData_PasswordProtected_NoPasswordGiven(t *testing.T) {
	pw := "secret"
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "locked-handle",
		Active:             true,
		Password:           &pw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "locked-handle").Return(share, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/locked-handle", nil)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareData_PasswordProtected_WrongPassword(t *testing.T) {
	rawPw := "right"
	hashedPw, err := hashutils.HashPassword(rawPw)
	require.NoError(t, err)
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "locked-handle",
		Active:             true,
		Password:           &hashedPw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "locked-handle").Return(share, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPost, "/certificate-shares/locked-handle",
		map[string]string{"password": "wrong"})

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareData_NotClaimed(t *testing.T) {
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "unclaimed-handle",
		Active:             true,
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

	resp := doRequest(app, http.MethodPost, "/certificate-shares/unclaimed-handle", nil)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

func TestGetCertificateShareData_Success_Public(t *testing.T) {
	certID := uuid.New()
	tokenID := "99"
	contractAddr := "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "claimed-handle",
		Active:             true,
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

	resp := doRequest(app, http.MethodPost, "/certificate-shares/claimed-handle", nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
	mockFactory.AssertExpectations(t)
	mockContract.AssertExpectations(t)
}

func TestGetCertificateShareData_Success_WithPassword(t *testing.T) {
	rawPw := "right"
	hashedPw, err := hashutils.HashPassword(rawPw)
	require.NoError(t, err)
	certID := uuid.New()
	tokenID := "7"
	contractAddr := "0xcafebabecafebabecafebabecafebabecafebabe"
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "pw-handle",
		Active:             true,
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
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "pw-handle").Return(share, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(cert, nil)

	mockFactory := new(mockCertContractFactoryDg)
	mockContract := new(mockCertContractDg)
	mockFactory.On("GetContract", common.HexToAddress(contractAddr)).Return(mockContract, nil)
	mockContract.On("GetTokenData", mock.Anything, big.NewInt(7)).Return(payload, nil)

	app := buildTestApp(nil, mockCertDg, mockShareDg, mockFactory)

	resp := doRequest(app, http.MethodPost, "/certificate-shares/pw-handle",
		map[string]string{"password": rawPw})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
	mockFactory.AssertExpectations(t)
	mockContract.AssertExpectations(t)
}

// TestGetCertificateShareData_AtContextField_NoSegfault is a regression test for
// the SIGSEGV that occurred when goccy/go-json (Fiber's global JSON encoder)
// tried to encode a CertificatePayload whose Header.Context []string field has
// json:"@context". The '@' prefix caused goccy/go-json to compute the wrong
// field offset, reading raw string bytes as a pointer and segfaulting.
// The fix marshals the response with encoding/json directly, bypassing goccy.
// This test configures Fiber with goccy/go-json as encoder (exactly like
// production) to guarantee the workaround remains in place.
func TestGetCertificateShareData_AtContextField_NoSegfault(t *testing.T) {
	certID := uuid.New()
	tokenID := "0"
	contractAddr := "0x46494f89533057ad6865b86d9619acd9a3cf7687"
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "vc-handle",
		Active:             true,
		Password:           nil,
	}
	cert := &entity.EventCertificate{
		Id:                      certID,
		CertificateTokenId:      &tokenID,
		EventCertificateAddress: &contractAddr,
	}
	// Full W3C VC payload — the @context []string is what triggered the segfault.
	payload := &entity.CertificatePayload{
		Header: entity.CertificatePayloadHeader{
			Context:      []string{"https://www.w3.org/2018/credentials/v1"},
			Type:         []string{"VerifiableCredential", "EventCertificate"},
			Id:           "cert-1",
			Issuer:       "issuer-1",
			IssuanceDate: "1772294796",
		},
		Data: entity.CertificatePayloadData{
			EventName:        "Test Event",
			CertificateTitle: "Certificate of Completion",
			Status:           "VALID",
		},
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "vc-handle").Return(share, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(cert, nil)

	mockFactory := new(mockCertContractFactoryDg)
	mockContract := new(mockCertContractDg)
	mockFactory.On("GetContract", common.HexToAddress(contractAddr)).Return(mockContract, nil)
	mockContract.On("GetTokenData", mock.Anything, big.NewInt(0)).Return(payload, nil)

	// Mirror the production Fiber config: use goccy/go-json as the encoder.
	// Without the stdjson.Marshal workaround in the handler this would segfault.
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := &certificate_share_usecase.CertificateShareUsecase{
		EventCertificateDataGateway:  mockCertDg,
		CertificateShareDg:           mockShareDg,
		CertificateContractFactoryDg: mockFactory,
	}
	h := &Handler{
		CertificateShareUc:    uc,
		AuthenticationService: &auth.AuthService{},
		Logger:                logger,
	}
	app := fiber.New(fiber.Config{
		JSONEncoder:  goccyjson.Marshal,
		JSONDecoder:  goccyjson.Unmarshal,
		ErrorHandler: customerror.GetErrFiberHandler(logger),
	})
	app.Post("/certificate-shares/:handle", h.GetCertificateShareData)

	resp := doRequest(app, http.MethodPost, "/certificate-shares/vc-handle", nil)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var got CertificateShareDataResponse
	body, _ := io.ReadAll(resp.Body)
	require.NoError(t, json.Unmarshal(body, &got))
	require.NotNil(t, got.Data)
	assert.Equal(t, []string{"https://www.w3.org/2018/credentials/v1"}, got.Data.Header.Context)
	assert.Equal(t, "Test Event", got.Data.Data.EventName)
	require.NotNil(t, got.Contract)
	assert.Equal(t, contractAddr, got.Contract.EventCertificateContractAddress)
	assert.Equal(t, tokenID, got.Contract.CertificateTokenId)
}

// ---------------------------------------------------------------------------
// Tests: UpdateCertificateShare
// ---------------------------------------------------------------------------

func TestUpdateCertificateShare_Unauthenticated(t *testing.T) {
	shareID := uuid.New()
	app := buildTestApp(nil, new(mockEventCertificateDataGateway), new(mockCertificateShareDataGateway), new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodPatch, "/certificate-shares/config/"+shareID.String(), nil)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUpdateCertificateShare_InvalidShareUUID(t *testing.T) {
	userID := uuid.New()
	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		new(mockEventCertificateDataGateway),
		new(mockCertificateShareDataGateway),
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPatch, "/certificate-shares/config/not-a-uuid", nil)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateCertificateShare_ShareNotFound(t *testing.T) {
	userID := uuid.New()
	shareID := uuid.New()

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByID", mock.Anything, shareID).Return(nil, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		new(mockEventCertificateDataGateway),
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPatch, "/certificate-shares/config/"+shareID.String(), nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestUpdateCertificateShare_Forbidden_NotOwner(t *testing.T) {
	ownerID := uuid.New()
	callerID := uuid.New()
	certID := uuid.New()
	shareID := uuid.New()

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByID", mock.Anything, shareID).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
	}, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(&entity.EventCertificate{
		Id:                   certID,
		ReceiverCredentialId: &ownerID,
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: callerID},
		mockCertDg,
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPatch, "/certificate-shares/config/"+shareID.String(), nil)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

func TestUpdateCertificateShare_Success_WithPassword(t *testing.T) {
	userID := uuid.New()
	certID := uuid.New()
	shareID := uuid.New()

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByID", mock.Anything, shareID).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
	}, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(&entity.EventCertificate{
		Id:                   certID,
		ReceiverCredentialId: &userID,
	}, nil)

	mockShareDg.On("UpdateCertificateShare", mock.Anything, shareID, mock.MatchedBy(func(p event_datagateway.UpdateCertificateShareParameters) bool {
		return p.Password.Changed && p.Password.Value != nil
	})).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
		Password:           strPtr("$argon2id$..."),
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		mockCertDg,
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	resp := doRequest(app, http.MethodPatch, "/certificate-shares/config/"+shareID.String(),
		map[string]string{"password": "newpassword"})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

func TestUpdateCertificateShare_Success_RemovePassword(t *testing.T) {
	userID := uuid.New()
	certID := uuid.New()
	shareID := uuid.New()

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByID", mock.Anything, shareID).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
	}, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(&entity.EventCertificate{
		Id:                   certID,
		ReceiverCredentialId: &userID,
	}, nil)

	mockShareDg.On("UpdateCertificateShare", mock.Anything, shareID, mock.MatchedBy(func(p event_datagateway.UpdateCertificateShareParameters) bool {
		return p.Password.Changed && p.Password.Value == nil
	})).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
		Password:           nil,
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		mockCertDg,
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	// Send empty string password to explicitly remove protection
	resp := doRequest(app, http.MethodPatch, "/certificate-shares/config/"+shareID.String(), map[string]string{"password": ""})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

func TestUpdateCertificateShare_Success_OmitPasswordField_KeepsExisting(t *testing.T) {
	userID := uuid.New()
	certID := uuid.New()
	shareID := uuid.New()
	existingPw := strPtr("$argon2id$existing")

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByID", mock.Anything, shareID).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
		Password:           existingPw,
	}, nil)

	mockCertDg := new(mockEventCertificateDataGateway)
	mockCertDg.On("GetEventCertificateByID", mock.Anything, certID).Return(&entity.EventCertificate{
		Id:                   certID,
		ReceiverCredentialId: &userID,
	}, nil)

	mockShareDg.On("UpdateCertificateShare", mock.Anything, shareID, mock.MatchedBy(func(p event_datagateway.UpdateCertificateShareParameters) bool {
		return !p.Password.Changed
	})).Return(&entity.CertificateShare{
		Id:                 shareID,
		EventCertificateId: certID,
		Password:           existingPw,
	}, nil)

	app := buildTestApp(
		&auth.JwtClaims{UserId: userID},
		mockCertDg,
		mockShareDg,
		new(mockCertContractFactoryDg),
	)

	// Omit the password field entirely — existing protection must be preserved
	resp := doRequest(app, http.MethodPatch, "/certificate-shares/config/"+shareID.String(), map[string]bool{"active": true})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
	mockCertDg.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: GetCertificateShareImage
// ---------------------------------------------------------------------------

func TestGetCertificateShareImage_NotFound(t *testing.T) {
	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "gone").Return(nil, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/gone/image", nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareImage_PasswordProtected_NoPasswordGiven(t *testing.T) {
	pw := "secret"
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "locked",
		Active:             true,
		Password:           &pw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "locked").Return(share, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	resp := doRequest(app, http.MethodGet, "/certificate-shares/locked/image", nil)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareImage_PasswordProtected_WrongPassword(t *testing.T) {
	rawPw := "right"
	hashedPw, err := hashutils.HashPassword(rawPw)
	require.NoError(t, err)
	certID := uuid.New()
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "locked",
		Active:             true,
		Password:           &hashedPw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "locked").Return(share, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg))

	req := httptest.NewRequest(http.MethodGet, "/certificate-shares/locked/image?password=wrong", nil)
	resp, _ := app.Test(req, 5000)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockShareDg.AssertExpectations(t)
}

func TestGetCertificateShareImage_Success_Public(t *testing.T) {
	certID := uuid.New()
	pngBytes := []byte{0x89, 0x50, 0x4E, 0x47}
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "public-handle",
		Active:             true,
		Password:           nil,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "public-handle").Return(share, nil)

	mockImageGen := new(mockCertificateImageGenerator)
	mockImageGen.On("GenerateCertificateImage", mock.Anything, certID).Return(pngBytes, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg), mockImageGen)

	resp := doRequest(app, http.MethodGet, "/certificate-shares/public-handle/image", nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, pngBytes, body)
	mockShareDg.AssertExpectations(t)
	mockImageGen.AssertExpectations(t)
}

func TestGetCertificateShareImage_Success_WithPassword(t *testing.T) {
	rawPw := "right"
	hashedPw, err := hashutils.HashPassword(rawPw)
	require.NoError(t, err)
	certID := uuid.New()
	pngBytes := []byte{0x89, 0x50, 0x4E, 0x47}
	share := &entity.CertificateShare{
		Id:                 uuid.New(),
		EventCertificateId: certID,
		Handle:             "pw-handle",
		Active:             true,
		Password:           &hashedPw,
	}

	mockShareDg := new(mockCertificateShareDataGateway)
	mockShareDg.On("GetCertificateShareByHandle", mock.Anything, "pw-handle").Return(share, nil)

	mockImageGen := new(mockCertificateImageGenerator)
	mockImageGen.On("GenerateCertificateImage", mock.Anything, certID).Return(pngBytes, nil)

	app := buildTestApp(nil, new(mockEventCertificateDataGateway), mockShareDg, new(mockCertContractFactoryDg), mockImageGen)

	req := httptest.NewRequest(http.MethodGet, "/certificate-shares/pw-handle/image?password="+rawPw, nil)
	resp, _ := app.Test(req, 5000)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, pngBytes, body)
	mockShareDg.AssertExpectations(t)
	mockImageGen.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func strPtr(s string) *string { return &s }
