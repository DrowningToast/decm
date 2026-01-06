package oauth

import (
	"apps/backend/common/customerror"
	"apps/backend/common/metrics"
	"apps/backend/common/utils"
	"apps/backend/core-api/config"
	"apps/backend/services/log"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUser represents user data from Google
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// OAuthResult contains both token and user information
type OAuthResult struct {
	Token *oauth2.Token `json:"token,omitempty"`
	User  *GoogleUser   `json:"user,omitempty"`
}

type GoogleOAuthService struct {
	googleConfig *oauth2.Config
	SessionStore *session.Store
	httpClient   *http.Client
	logger       *slog.Logger
}

var _ OAuthService = &GoogleOAuthService{}

func NewGoogleOAuthService() *GoogleOAuthService {
	cfg := config.LoadConfig()

	// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
	googleConfig := &oauth2.Config{
		RedirectURL:  cfg.GoogleOAuth.RedirectURL,
		ClientID:     cfg.GoogleOAuth.ClientID,
		ClientSecret: cfg.GoogleOAuth.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	// Configure session store with security settings
	sessionStore := session.New(session.Config{
		Expiration:     time.Hour * 2, // 2 hour expiration
		CookieHTTPOnly: true,
		KeyLookup:      "cookie:google_oauth_session", // Format: <source>:<name>
		CookiePath:     "/",                           // Ensure cookie is sent on all paths
		CookieSecure:   cfg.Env == "production",       // Only secure in production
		CookieSameSite: "Lax",
		CookieDomain:   cfg.CookieDomain, // Set domain from config (empty for dev, domain for production)
		KeyGenerator: func() string {
			return utils.GenerateSecureRandomString(32)
		},
	})

	// Configure HTTP client with timeouts and retry logic
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  false,
			MaxIdleConnsPerHost: 10,
		},
	}

	return &GoogleOAuthService{
		googleConfig: googleConfig,
		SessionStore: sessionStore,
		httpClient:   httpClient,
		logger:       log.NewLogger(),
	}
}

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *GoogleOAuthService) getUserDataFromGoogle(ctx context.Context, code string) (*oauth2.Token, error) {
	start := time.Now()
	operation := "token_exchange"

	// Use code to get token and get user info from Google.
	token, err := s.googleConfig.Exchange(ctx, code)

	duration := time.Since(start).Seconds()
	status := "success"

	if err != nil {
		status = "error"
		metrics.GoogleOAuthOperationsTotal.WithLabelValues(operation, status).Inc()
		metrics.GoogleOAuthOperationDuration.WithLabelValues(operation, status).Observe(duration)

		s.logger.Error("Code exchange wrong: %s", slog.String("error", err.Error()))
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("code exchange wrong: %s"))
	}

	metrics.GoogleOAuthOperationsTotal.WithLabelValues(operation, status).Inc()
	metrics.GoogleOAuthOperationDuration.WithLabelValues(operation, status).Observe(duration)

	return token, nil
}

func (s *GoogleOAuthService) Login(session *session.Session) (*string, error) {
	state, err := generateState()
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	session.Set("state", state)
	if err := session.Save(); err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	url := s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return &url, nil
}

func (s *GoogleOAuthService) Callback(ctx context.Context, session *session.Session, code string, state string) (*oauth2.Token, error) {
	savedStateInterface := session.Get("state")
	if savedStateInterface == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("state not found in session"))
	}
	savedState, ok := savedStateInterface.(string)
	if !ok {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid state type in session"))
	}
	if savedState != state {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("state mismatch"))
	}

	token, err := s.getUserDataFromGoogle(ctx, code)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return token, nil
}

// validateToken checks if the token is valid and not expired
func (s *GoogleOAuthService) validateToken(token *oauth2.Token) error {
	if token == nil {
		return errors.New("token is nil")
	}

	if token.AccessToken == "" {
		return errors.New("access token is empty")
	}

	// Only validate expiry if it's actually set (not zero time)
	if !token.Expiry.IsZero() {
		if !token.Valid() {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("token has expired or is invalid"))
		}

		// Warn if token expires soon (within 5 minutes)
		if time.Until(token.Expiry) < 5*time.Minute {
			s.logger.Warn("Warning: Token expires soon", slog.String("expiry", token.Expiry.String()))
		}
	}

	return nil
}

// getUserInfoFromGoogle fetches user information using OAuth2 client
func (s *GoogleOAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUser, error) {
	start := time.Now()
	operation := "get_user_info"

	// Validate token first
	if err := s.validateToken(token); err != nil {
		s.logger.Warn("Invalid token: %v", slog.String("error", err.Error()))
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Create authenticated HTTP client using oauth2 library
	client := s.googleConfig.Client(ctx, token)

	// Make request to Google's userinfo endpoint
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	duration := time.Since(start).Seconds()

	if err != nil {
		metrics.GoogleOAuthOperationsTotal.WithLabelValues(operation, "error").Inc()
		metrics.GoogleOAuthOperationDuration.WithLabelValues(operation, "error").Observe(duration)

		s.logger.Warn("Failed to get user info: %v", slog.String("error", err.Error()))
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to get user info"))
	}
	defer func() { _ = resp.Body.Close() }()

	status := "success"
	if resp.StatusCode >= 400 {
		status = "error"
	}
	metrics.GoogleOAuthOperationsTotal.WithLabelValues(operation, status).Inc()
	metrics.GoogleOAuthOperationDuration.WithLabelValues(operation, status).Observe(duration)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Warn("Failed to read response body: %v", slog.String("error", err.Error()))
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to read response body"))
	}

	// Handle different status codes
	switch resp.StatusCode {
	case http.StatusOK:
		var user OAuthUser
		if err := json.Unmarshal(body, &user); err != nil {
			s.logger.Warn("Failed to unmarshal user info: %v", slog.String("error", err.Error()))
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to unmarshal user info"))
		}

		// Validate required fields
		if user.Email == "" {
			s.logger.Warn("User email is empty", slog.String("email", user.Email))
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user email is empty"))
		}

		s.logger.Info("Successfully retrieved user info for: %s", slog.String("email", user.Email))
		return &OAuthUser{
			Id:    user.Id,
			Email: user.Email,
		}, nil

	case http.StatusUnauthorized:
		s.logger.Warn("Unauthorized: token may be expired or invalid", slog.String("error", "unauthorized: token may be expired or invalid"))
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("unauthorized: token may be expired or invalid"))

	case http.StatusForbidden:
		s.logger.Warn("Forbidden: insufficient permissions for user info", slog.String("error", "forbidden: insufficient permissions for user info"))
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("forbidden: insufficient permissions for user info"))

	case http.StatusTooManyRequests:
		s.logger.Error("Rate limited: too many requests to Google API", slog.String("error", "rate limited: too many requests to Google API"))
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("rate limited: too many requests to Google API"))

	default:
		s.logger.Error("Google API error %d: %s", slog.Int("status_code", resp.StatusCode), slog.String("error", string(body)))
		return nil, customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("google API error %d: %s", resp.StatusCode, string(body)))
	}
}
