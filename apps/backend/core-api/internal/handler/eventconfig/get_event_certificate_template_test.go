package eventconfig

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"decm-database/go/generated"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	event_usecase "apps/backend/core-api/internal/usecase/event"
	eventconfig_usecase "apps/backend/core-api/internal/usecase/eventconfig"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

// mockEventDataGateway implements eventdatagateway.EventDataGateway.
// Only GetEventById is configurable; all other methods are stubs.
type mockEventDataGateway struct{ mock.Mock }

func (m *mockEventDataGateway) GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *mockEventDataGateway) CreateEvent(_ context.Context, _ eventdatagateway.CreateEventParameters) (*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventDataGateway) GetViewModelById(_ context.Context, _ uuid.UUID) (*entity.Event, *entity.EventRegistrationConfig, *entity.EventContract, error) {
	return nil, nil, nil, errors.New("not implemented")
}

func (m *mockEventDataGateway) ListEventsByOwnerCredentialID(_ context.Context, _ uuid.UUID, _ int32, _ int32) ([]*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventDataGateway) UpdateEvent(_ context.Context, _ uuid.UUID, _ eventdatagateway.UpdateEventParameters) (*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventDataGateway) DeleteEvent(_ context.Context, _ uuid.UUID) (*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventDataGateway) ListEvents(_ context.Context, _ *int32, _ *int32) ([]*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventDataGateway) ListEventsByEventAttendeeCredentialID(_ context.Context, _ uuid.UUID, _ *int32, _ *int32) ([]*entity.Event, error) {
	return nil, errors.New("not implemented")
}

// mockEventIssuerDataGateway implements eventdatagateway.EventIssuerDataGateway.
// Only GetEventIssuerByEventIDAndIssuerCredentialID is configurable.
type mockEventIssuerDataGateway struct{ mock.Mock }

func (m *mockEventIssuerDataGateway) GetEventIssuerByEventIDAndIssuerCredentialID(ctx context.Context, eventID uuid.UUID, issuerCredentialID uuid.UUID) (generated.EventIssuer, error) {
	args := m.Called(ctx, eventID, issuerCredentialID)
	return args.Get(0).(generated.EventIssuer), args.Error(1)
}

func (m *mockEventIssuerDataGateway) CreateEventIssuer(_ context.Context, _ generated.CreateEventIssuerParams) (*generated.EventIssuer, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) GetEventIssuerByID(_ context.Context, _ uuid.UUID) (generated.EventIssuer, error) {
	return generated.EventIssuer{}, errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) GetEventIssuersByEventID(_ context.Context, _ uuid.UUID) ([]generated.EventIssuer, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) UpdateEventIssuer(_ context.Context, _ generated.UpdateEventIssuerParams) (*generated.EventIssuer, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) DeleteEventIssuer(_ context.Context, _ uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) UpdateEventIssuerSigningStatus(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ int32) error {
	return errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) ResetAllEventIssuersSigningStatus(_ context.Context, _ uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) HasSignedIssuers(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) GetSignedIssuersCount(_ context.Context, _ uuid.UUID) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) GetTotalIssuersCount(_ context.Context, _ uuid.UUID) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockEventIssuerDataGateway) AllIssuersHaveSigned(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, errors.New("not implemented")
}

// mockEventCertificateConfigDataGateway implements eventdatagateway.EventCertificateConfigDataGateway.
// Only GetEventCertificateConfigByEventID is configurable.
type mockEventCertConfigDataGateway struct{ mock.Mock }

func (m *mockEventCertConfigDataGateway) GetEventCertificateConfigByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventCertificateConfig, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateConfig), args.Error(1)
}

func (m *mockEventCertConfigDataGateway) CreateEventCertificateConfig(_ context.Context, _ generated.CreateEventCertificateConfigParams) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventCertConfigDataGateway) GetEventCertificateConfigByID(_ context.Context, _ uuid.UUID) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventCertConfigDataGateway) UpdateEventCertificateConfig(_ context.Context, _ generated.UpdateEventCertificateConfigParams) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventCertConfigDataGateway) UpdateEventCertificateTextConfig(_ context.Context, _ generated.UpdateEventCertificateTextConfigParams) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventCertConfigDataGateway) ToggleEventCertificateConfigPublished(_ context.Context, _ generated.ToggleEventCertificateConfigPublishedParams) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEventCertConfigDataGateway) DeleteEventCertificateConfig(_ context.Context, _ uuid.UUID) error {
	return errors.New("not implemented")
}

// ---------------------------------------------------------------------------
// Test helper
// ---------------------------------------------------------------------------

// buildTestApp wires a minimal Handler into a Fiber app and returns it along
// with the event usecase (so callers can set expectations on its mocks).
func buildTestApp(
	mockEventDg *mockEventDataGateway,
	mockIssuerDg *mockEventIssuerDataGateway,
	mockCertCfgDg *mockEventCertConfigDataGateway,
) *fiber.App {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Build EventUsecase with only the DGs the handler actually uses.
	// Remaining parameters are nil; they are not invoked in these tests.
	eventUc := event_usecase.NewEventUsecase(
		mockEventDg,  // EventDataGateway
		nil,          // EventContractDataGateway
		nil,          // EventContractFactoryDg
		mockIssuerDg, // EventIssuerDataGateway
		nil,          // EventCertificateDataGateway
		nil,          // CertificateShareDataGateway
		nil,          // EventCertificateSignatureDataGateway
		nil,          // EventCertificateConfigDg
		nil,          // EventCertificateFontFamilyDg
		nil,          // AuthenticationCredentialDg
		nil,          // EventRegistrationInvitationDg
		nil,          // EventAttendeeDg
		nil,          // UserSignatureDg
		nil,          // InboxMessageDg
		nil,          // BlockchainClientDg
		nil,          // ethClient
		nil,          // S3DataGateway
		discardLogger,
		nil, // authService
		nil, // cfg
	)

	eventConfigUc := &eventconfig_usecase.EventConfigUsecase{
		EventCertificateDg: mockCertCfgDg,
		// S3Service zero-value is fine: tests that reach this point make
		// GetEventCertificateConfigByEventID return an error, so GetFile is
		// never called.
	}

	h := &Handler{
		EventUc:               eventUc,
		EventConfigUc:         eventConfigUc,
		AuthenticationService: &auth.AuthService{},
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: customerror.GetErrFiberHandler(discardLogger),
	})
	app.Get("/events/:event_id/config/certificate/template", h.GetEventCertificateTemplate)
	return app
}

// injectUser is a middleware that sets the authenticated user in Fiber Locals.
func injectUser(claims *auth.JwtClaims) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	}
}

func doRequest(app *fiber.App, eventID string, claims *auth.JwtClaims) *http.Response {
	req := httptest.NewRequest(http.MethodGet, "/events/"+eventID+"/config/certificate/template", nil)
	if claims != nil {
		// We must inject the user before the handler runs. Register a
		// per-request middleware on a fresh sub-app that wraps the route.
		subApp := fiber.New(fiber.Config{ErrorHandler: app.Config().ErrorHandler})
		subApp.Use(injectUser(claims))
		// Re-register the route on the sub-app with the same handler by
		// proxying through the original app's router is tricky, so instead
		// we build the full test via buildTestApp below.  This helper is
		// only used when the app was already built with the user middleware.
		_ = subApp
	}
	resp, _ := app.Test(req, 3000)
	return resp
}

// buildTestAppWithUser is the primary helper: builds the app and registers the
// user-injection middleware so every request is authenticated as `claims`.
func buildTestAppWithUser(
	claims *auth.JwtClaims,
	mockEventDg *mockEventDataGateway,
	mockIssuerDg *mockEventIssuerDataGateway,
	mockCertCfgDg *mockEventCertConfigDataGateway,
) *fiber.App {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	eventUc := event_usecase.NewEventUsecase(
		mockEventDg,
		nil, nil,
		mockIssuerDg,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		discardLogger,
		nil, nil,
	)

	eventConfigUc := &eventconfig_usecase.EventConfigUsecase{
		EventCertificateDg: mockCertCfgDg,
	}

	h := &Handler{
		EventUc:               eventUc,
		EventConfigUc:         eventConfigUc,
		AuthenticationService: &auth.AuthService{},
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: customerror.GetErrFiberHandler(discardLogger),
	})
	if claims != nil {
		app.Use(injectUser(claims))
	}
	app.Get("/events/:event_id/config/certificate/template", h.GetEventCertificateTemplate)
	return app
}

func get(app *fiber.App, eventID string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, "/events/"+eventID+"/config/certificate/template", nil)
	resp, _ := app.Test(req, 3000)
	return resp
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

const validEventID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

func TestGetEventCertificateTemplate_InvalidEventID(t *testing.T) {
	app := buildTestAppWithUser(
		&auth.JwtClaims{UserId: uuid.New()},
		new(mockEventDataGateway),
		new(mockEventIssuerDataGateway),
		new(mockEventCertConfigDataGateway),
	)

	resp := get(app, "not-a-uuid")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetEventCertificateTemplate_Unauthenticated(t *testing.T) {
	// No claims injected — RequireUserContext returns 401.
	app := buildTestAppWithUser(
		nil, // no user
		new(mockEventDataGateway),
		new(mockEventIssuerDataGateway),
		new(mockEventCertConfigDataGateway),
	)

	resp := get(app, validEventID)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetEventCertificateTemplate_EventNotFound(t *testing.T) {
	eventID := uuid.MustParse(validEventID)
	userID := uuid.New()

	mockEventDg := new(mockEventDataGateway)
	mockEventDg.On("GetEventById", mock.Anything, eventID).
		Return(nil, nil) // nil event, nil error → handler returns 404

	app := buildTestAppWithUser(
		&auth.JwtClaims{UserId: userID},
		mockEventDg,
		new(mockEventIssuerDataGateway),
		new(mockEventCertConfigDataGateway),
	)

	resp := get(app, validEventID)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockEventDg.AssertExpectations(t)
}

func TestGetEventCertificateTemplate_ForbiddenNotOwnerNotIssuer(t *testing.T) {
	eventID := uuid.MustParse(validEventID)
	ownerID := uuid.New()
	otherUserID := uuid.New()

	event := &entity.Event{OwnerCredentialId: ownerID}

	mockEventDg := new(mockEventDataGateway)
	mockEventDg.On("GetEventById", mock.Anything, eventID).Return(event, nil)

	// User is not the owner — IsUserIssuerForEvent is called.
	// pgx.ErrNoRows signals "no issuer row found" → returns (false, nil).
	mockIssuerDg := new(mockEventIssuerDataGateway)
	mockIssuerDg.On("GetEventIssuerByEventIDAndIssuerCredentialID", mock.Anything, eventID, otherUserID).
		Return(generated.EventIssuer{}, pgx.ErrNoRows)

	app := buildTestAppWithUser(
		&auth.JwtClaims{UserId: otherUserID},
		mockEventDg,
		mockIssuerDg,
		new(mockEventCertConfigDataGateway),
	)

	resp := get(app, validEventID)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	mockEventDg.AssertExpectations(t)
	mockIssuerDg.AssertExpectations(t)
}

func TestGetEventCertificateTemplate_AllowedOwner(t *testing.T) {
	eventID := uuid.MustParse(validEventID)
	ownerID := uuid.New()

	event := &entity.Event{OwnerCredentialId: ownerID}

	mockEventDg := new(mockEventDataGateway)
	mockEventDg.On("GetEventById", mock.Anything, eventID).Return(event, nil)

	// IsUserIssuerForEvent is NOT called when the caller is the owner.

	// After auth passes, GetEventCertificateTemplateSVG is called.
	// Return an error so we never touch the zero-value S3Service.
	mockCertCfgDg := new(mockEventCertConfigDataGateway)
	mockCertCfgDg.On("GetEventCertificateConfigByEventID", mock.Anything, eventID).
		Return(nil, errors.New("db error"))

	app := buildTestAppWithUser(
		&auth.JwtClaims{UserId: ownerID},
		mockEventDg,
		new(mockEventIssuerDataGateway),
		mockCertCfgDg,
	)

	resp := get(app, validEventID)

	// Auth passed (not 401/403). The 500 is from the downstream DB error,
	// not from the access-control guard.
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotEqual(t, http.StatusForbidden, resp.StatusCode)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockEventDg.AssertExpectations(t)
	mockCertCfgDg.AssertExpectations(t)
}

func TestGetEventCertificateTemplate_AllowedIssuer(t *testing.T) {
	eventID := uuid.MustParse(validEventID)
	ownerID := uuid.New()
	issuerID := uuid.New()

	event := &entity.Event{OwnerCredentialId: ownerID}

	mockEventDg := new(mockEventDataGateway)
	mockEventDg.On("GetEventById", mock.Anything, eventID).Return(event, nil)

	// Issuer is found in the event_issuers table → (true, nil).
	mockIssuerDg := new(mockEventIssuerDataGateway)
	mockIssuerDg.On("GetEventIssuerByEventIDAndIssuerCredentialID", mock.Anything, eventID, issuerID).
		Return(generated.EventIssuer{}, nil)

	// After auth passes, GetEventCertificateTemplateSVG is called.
	mockCertCfgDg := new(mockEventCertConfigDataGateway)
	mockCertCfgDg.On("GetEventCertificateConfigByEventID", mock.Anything, eventID).
		Return(nil, errors.New("db error"))

	app := buildTestAppWithUser(
		&auth.JwtClaims{UserId: issuerID},
		mockEventDg,
		mockIssuerDg,
		mockCertCfgDg,
	)

	resp := get(app, validEventID)

	// Auth passed (not 401/403).
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotEqual(t, http.StatusForbidden, resp.StatusCode)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockEventDg.AssertExpectations(t)
	mockIssuerDg.AssertExpectations(t)
	mockCertCfgDg.AssertExpectations(t)
}

func TestGetEventCertificateTemplate_IssuersCheckSkippedForOwner(t *testing.T) {
	// Verifies that IsUserIssuerForEvent is never called when the caller IS the owner.
	eventID := uuid.MustParse(validEventID)
	ownerID := uuid.New()

	event := &entity.Event{OwnerCredentialId: ownerID}

	mockEventDg := new(mockEventDataGateway)
	mockEventDg.On("GetEventById", mock.Anything, eventID).Return(event, nil)

	mockIssuerDg := new(mockEventIssuerDataGateway) // No expectations set.

	mockCertCfgDg := new(mockEventCertConfigDataGateway)
	mockCertCfgDg.On("GetEventCertificateConfigByEventID", mock.Anything, eventID).
		Return(nil, errors.New("db error"))

	app := buildTestAppWithUser(
		&auth.JwtClaims{UserId: ownerID},
		mockEventDg,
		mockIssuerDg,
		mockCertCfgDg,
	)

	get(app, validEventID)

	// AssertExpectations on mockIssuerDg with no expectations means no call was made.
	mockIssuerDg.AssertExpectations(t)
}
