package event_registration

import (
	"apps/backend/services/auth"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	customerror "apps/backend/common/customerror"
	eventRegistrationUc "apps/backend/core-api/internal/usecase/event_registration"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func sp(s string) *string { return &s }

func newTestApp(h *Handler) *fiber.App {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	app := fiber.New(fiber.Config{
		ErrorHandler: customerror.GetErrFiberHandler(discardLogger),
	})
	claims := &auth.JwtClaims{UserId: uuid.New()}
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})
	app.Post("/event-registration/invitation/:event_id/import", h.ImportEventParticipants)
	return app
}

func postImport(app *fiber.App, eventID uuid.UUID, body interface{}) *http.Response {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/event-registration/invitation/"+eventID.String()+"/import", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	return resp
}

// ---------------------------------------------------------------------------
// IsValid — unit tests (no fiber, no usecase)
// ---------------------------------------------------------------------------

func TestParticipantXORValidation_IsValid(t *testing.T) {
	validEventID := uuid.New().String()

	t.Run("fails when both email and wallet_address are provided", func(t *testing.T) {
		req := ImportEventParticipantsRequest{
			EventID: validEventID,
			Participants: []ParticipantRequestItem{
				{FirstName: "John", LastName: "Doe", Email: "john@example.com", WalletAddress: sp("0xabc")},
			},
		}
		err := req.IsValid()
		require.Error(t, err)
	})

	t.Run("fails when neither email nor wallet_address are provided", func(t *testing.T) {
		req := ImportEventParticipantsRequest{
			EventID: validEventID,
			Participants: []ParticipantRequestItem{
				{FirstName: "John", LastName: "Doe"},
			},
		}
		err := req.IsValid()
		require.Error(t, err)
	})

	t.Run("passes when only email is provided", func(t *testing.T) {
		req := ImportEventParticipantsRequest{
			EventID: validEventID,
			Participants: []ParticipantRequestItem{
				{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
			},
		}
		err := req.IsValid()
		require.NoError(t, err)
	})

	t.Run("passes when only wallet_address is provided", func(t *testing.T) {
		req := ImportEventParticipantsRequest{
			EventID: validEventID,
			Participants: []ParticipantRequestItem{
				{FirstName: "John", LastName: "Doe", WalletAddress: sp("0xabc")},
			},
		}
		err := req.IsValid()
		require.NoError(t, err)
	})

	t.Run("fails when any entry in a mixed list violates XOR", func(t *testing.T) {
		req := ImportEventParticipantsRequest{
			EventID: validEventID,
			Participants: []ParticipantRequestItem{
				{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
				{FirstName: "Jane", LastName: "Smith"}, // violates XOR
			},
		}
		err := req.IsValid()
		require.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// HTTP handler — integration tests (XOR errors must return 400)
// ---------------------------------------------------------------------------

func TestImportEventParticipantsHandler_XOR(t *testing.T) {
	eventID := uuid.New()
	h := &Handler{
		AuthenticationService: auth.AuthService{},
		EventRegistrationUc:   &eventRegistrationUc.EventRegistrationUsecase{},
	}
	app := newTestApp(h)

	t.Run("returns 400 when participant has both email and wallet_address", func(t *testing.T) {
		body := map[string]interface{}{
			"event_id": eventID.String(),
			"participants": []map[string]interface{}{
				{"first_name": "John", "last_name": "Doe", "email": "john@example.com", "wallet_address": "0xabc"},
			},
		}
		resp := postImport(app, eventID, body)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 400 when participant has neither email nor wallet_address", func(t *testing.T) {
		body := map[string]interface{}{
			"event_id": eventID.String(),
			"participants": []map[string]interface{}{
				{"first_name": "John", "last_name": "Doe"},
			},
		}
		resp := postImport(app, eventID, body)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
