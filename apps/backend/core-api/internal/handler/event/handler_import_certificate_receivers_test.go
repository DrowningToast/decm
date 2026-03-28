package event

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

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func strPtr(s string) *string { return &s }

func newCertificateTestApp(h *Handler) *fiber.App {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	app := fiber.New(fiber.Config{
		ErrorHandler: customerror.GetErrFiberHandler(discardLogger),
	})
	claims := &auth.JwtClaims{UserId: uuid.New()}
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})
	app.Post("/events/certificates/import", h.ImportCertificateReceivers)
	return app
}

func postCertificateImport(app *fiber.App, body interface{}) *http.Response {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/events/certificates/import", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	return resp
}

// ---------------------------------------------------------------------------
// IsValid — unit tests (no fiber, no usecase)
// ---------------------------------------------------------------------------

func TestImportCertificateReceiversRequest_IsValid(t *testing.T) {
	validEventID := uuid.New()
	validPin := "123456"

	t.Run("passes when only email is provided (PIN path)", func(t *testing.T) {
		req := ImportCertificateReceiversRequest{
			EventID: validEventID,
			HostPin: &validPin,
			Receivers: []ImportCertificateReceiverRequest{
				{Email: strPtr("alice@example.com")},
			},
		}
		assert.NoError(t, req.IsValid())
	})

	t.Run("passes when only wallet_address is provided (PIN path)", func(t *testing.T) {
		req := ImportCertificateReceiversRequest{
			EventID: validEventID,
			HostPin: &validPin,
			Receivers: []ImportCertificateReceiverRequest{
				{WalletAddress: strPtr("0x1234567890abcdef1234567890abcdef12345678")},
			},
		}
		assert.NoError(t, req.IsValid())
	})

	t.Run("passes with wallet-based auth path", func(t *testing.T) {
		signMsg := `{"eventContractAddress":"0xabc","receivers":["0xhash"]}`
		sig := "0xdeadbeef"
		req := ImportCertificateReceiversRequest{
			EventID:         validEventID,
			HostSignMessage: &signMsg,
			HostSignature:   &sig,
			Receivers: []ImportCertificateReceiverRequest{
				{Email: strPtr("alice@example.com")},
			},
		}
		assert.NoError(t, req.IsValid())
	})

	t.Run("fails when both email and wallet_address are provided", func(t *testing.T) {
		req := ImportCertificateReceiversRequest{
			EventID: validEventID,
			HostPin: &validPin,
			Receivers: []ImportCertificateReceiverRequest{
				{Email: strPtr("alice@example.com"), WalletAddress: strPtr("0xabc")},
			},
		}
		assert.Error(t, req.IsValid())
	})

	t.Run("fails when neither email nor wallet_address is provided", func(t *testing.T) {
		req := ImportCertificateReceiversRequest{
			EventID: validEventID,
			HostPin: &validPin,
			Receivers: []ImportCertificateReceiverRequest{
				{},
			},
		}
		assert.Error(t, req.IsValid())
	})

	t.Run("fails when receivers list is empty", func(t *testing.T) {
		req := ImportCertificateReceiversRequest{
			EventID:   validEventID,
			HostPin:   &validPin,
			Receivers: []ImportCertificateReceiverRequest{},
		}
		assert.Error(t, req.IsValid())
	})

	t.Run("fails when any entry in a mixed list violates XOR", func(t *testing.T) {
		req := ImportCertificateReceiversRequest{
			EventID: validEventID,
			HostPin: &validPin,
			Receivers: []ImportCertificateReceiverRequest{
				{Email: strPtr("alice@example.com")},
				{}, // violates XOR — neither
			},
		}
		assert.Error(t, req.IsValid())
	})

	t.Run("fails when neither host_pin nor wallet auth is provided", func(t *testing.T) {
		req := ImportCertificateReceiversRequest{
			EventID: validEventID,
			Receivers: []ImportCertificateReceiverRequest{
				{Email: strPtr("alice@example.com")},
			},
		}
		assert.Error(t, req.IsValid())
	})

	t.Run("fails when both host_pin and wallet signature are provided", func(t *testing.T) {
		signMsg := `{"eventContractAddress":"0xabc","receivers":["0xhash"]}`
		sig := "0xdeadbeef"
		req := ImportCertificateReceiversRequest{
			EventID:         validEventID,
			HostPin:         &validPin,
			HostSignMessage: &signMsg,
			HostSignature:   &sig,
			Receivers: []ImportCertificateReceiverRequest{
				{Email: strPtr("alice@example.com")},
			},
		}
		assert.Error(t, req.IsValid())
	})
}
