---
description: "Google OAuth 2.0 authentication patterns and implementation guidelines for DECM platform"
---

# Google OAuth 2.0 Authentication Patterns

## 🔐 OAuth Flow Architecture

### Service Layer Pattern
OAuth services follow a clean architecture with interface segregation:

**Interface Definition** ([oauth.go](mdc:apps/backend/services/oauth/oauth.go)):
```go
type OAuthService interface {
    Login(session *session.Session) (*string, *customerror.Err)
    Callback(ctx context.Context, session *session.Session, code string, state string) (*oauth2.Token, *customerror.Err)
    GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUser, *customerror.Err)
}
```

**Implementation** ([google.go](mdc:apps/backend/services/oauth/google.go)):
```go
type GoogleOAuthService struct {
    googleConfig *oauth2.Config
    SessionStore *session.Store
    httpClient   *http.Client
    logger       *slog.Logger
}
```

## 🛡️ Security Patterns

### CSRF Protection with State Validation
**Always generate cryptographically secure state**:
```go
func generateState() (string, error) {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}
```

**State validation in callbacks**:
```go
savedState := session.Get("state")
if savedState != state {
    return nil, customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("state mismatch"))
}
```

### Session Security Configuration
```go
sessionStore := session.New(session.Config{
    Expiration:     time.Hour * 2,          // 2 hour expiration
    CookieHTTPOnly: true,                   // Prevent XSS
    CookieSecure:   cfg.ENV == "production", // HTTPS in prod
    CookieSameSite: "Lax",                  // CSRF protection
    KeyGenerator: func() string {
        return utils.GenerateSecureRandomString(32)
    },
})
```

## 📡 Handler Patterns

### OAuth Flow Initiation
Pattern: Redirect with state management ([handle_request_google_oauth.go](mdc:apps/backend/core-api/internal/handler/auth/handle_request_google_oauth.go)):

```go
func (h Handler) RequestGoogleOAuth(ctx *fiber.Ctx) error {
    session, err := h.GoogleOAuthService.SessionStore.Get(ctx)
    if err != nil {
        return *customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
    }

    url, err := h.GoogleOAuthService.Login(session)
    if err != nil {
        return *customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
    }

    return ctx.Redirect(*url, fiber.StatusTemporaryRedirect)
}
```

### Token Verification Handler
Pattern: Validate input, verify state, return tokens ([handle_verify_google_oauth.go](mdc:apps/backend/core-api/internal/handler/auth/handle_verify_google_oauth.go)):

```go
func (h Handler) VerifyGoogleOAuth(ctx *fiber.Ctx) error {
    requestBody := verifyGoogleOAuthRequest{}
    if err := requestBody.Parse(ctx); err != nil {
        return *err
    }
    if err := requestBody.IsValid(); err != nil {
        return *err
    }

    session, err := h.GoogleOAuthService.SessionStore.Get(ctx)
    if err != nil {
        return *customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
    }

    token, err := h.AuthUc.VerifyGoogleOAuthCode(ctx.UserContext(), session, requestBody.Code, requestBody.State)
    if err != nil {
        return *customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
    }

    response := verifyGoogleOAuthResponse{
        AccessToken:  token.AccessToken,
        ExpiresIn:    int(token.Expiry.Sub(time.Now()).Seconds()),
        RefreshToken: token.RefreshToken,
    }
    return ctx.Status(fiber.StatusOK).JSON(response)
}
```

## 🔄 Registration Integration

### OAuth User Registration
Pattern: Token parsing, user info retrieval, wallet generation ([handle_register_with_google_oauth.go](mdc:apps/backend/core-api/internal/handler/onboard/handle_register_with_google_oauth.go)):

```go
func (h Handler) RegisterWithGoogleOAuth(ctx *fiber.Ctx) error {
    // Parse and validate request
    requestBody := registerWithGoogleOAuthRequest{}
    if err := requestBody.Parse(ctx); err != nil {
        return *err
    }

    // Parse OAuth token
    token, err := oauth.ParseToken(requestBody.AccessToken, requestBody.RefreshToken)
    if err != nil {
        return *customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, err)
    }

    // Register user with Google OAuth
    jwt, mnemonic, err := h.OnboardUc.RegisterWithGoogle(ctx.UserContext(), token, requestBody.Password)
    if err != nil {
        return *err
    }

    // Set session cookie
    cookie := new(fiber.Cookie)
    cookie.Name = "session"
    cookie.Value = *jwt
    cookie.Expires = time.Now().Add(h.SessionExpiration)
    ctx.Cookie(cookie)

    return ctx.Status(fiber.StatusOK).JSON(registerWithGoogleOAuthResponse{
        Mnemonic: mnemonic,
    })
}
```

## 🗃️ Database Integration

### Encrypted OAuth References
Pattern: Store encrypted Google connector references ([authentication_credentials.sql](mdc:packages/database/queries/authentication_credentials.sql)):

```sql
-- Encrypt on INSERT/UPDATE
google_connector_ref = CASE 
    WHEN sqlc.narg(google_connector_ref) IS NOT NULL 
    THEN pgp_sym_encrypt(sqlc.narg(google_connector_ref), sqlc.arg(encryption_key)::varchar)
    ELSE NULL 
END::text

-- Decrypt on SELECT
CASE 
    WHEN google_connector_ref IS NOT NULL 
    THEN pgp_sym_decrypt(google_connector_ref::bytea, sqlc.arg(encryption_key)::varchar)
    ELSE NULL 
END::text as google_connector_ref
```

### OAuth User Lookup
Pattern: Secure lookup by encrypted Google reference:
```sql
-- name: GetAuthenticationCredentialByGoogleConnectorRef :one
SELECT * FROM authentication_credentials 
WHERE google_connector_ref = pgp_sym_encrypt(sqlc.arg(google_connector_ref), sqlc.arg(encryption_key)::varchar)
```

## ⚙️ Configuration Patterns

### Environment-Based OAuth Config
Pattern: Separate configuration struct ([google.go](mdc:apps/backend/core-api/config/google/google.go)):

```go
type GoogleOAuthConfig struct {
    RedirectURL  string `env:"REDIRECT_URL" envDefault:"http://localhost:8000/api/v1/onboard/google/callback"`
    ClientID     string `env:"CLIENT_ID"`
    ClientSecret string `env:"CLIENT_SECRET"`
}
```

### OAuth 2.0 Client Configuration
```go
googleConfig := &oauth2.Config{
    RedirectURL:  cfg.GoogleOAuth.RedirectURL,
    ClientID:     cfg.GoogleOAuth.ClientID,
    ClientSecret: cfg.GoogleOAuth.ClientSecret,
    Scopes: []string{
        "https://www.googleapis.com/auth/userinfo.email",
        "https://www.googleapis.com/auth/photoslibrary.readonly",
    },
    Endpoint: google.Endpoint,
}
```

## 🚨 Error Handling Patterns

### Token Validation
```go
func (s *GoogleOAuthService) validateToken(token *oauth2.Token) error {
    if token == nil {
        return errors.New("token is nil")
    }
    if token.AccessToken == "" {
        return errors.New("access token is empty")
    }
    if !token.Valid() {
        return errors.New("token has expired or is invalid")
    }
    // Warn if token expires soon
    if time.Until(token.Expiry) < 5*time.Minute {
        s.logger.Warn("Warning: Token expires soon", slog.String("expiry", token.Expiry.String()))
    }
    return nil
}
```

### HTTP Status Code Handling
```go
switch resp.StatusCode {
case http.StatusOK:
    // Process user data
case http.StatusUnauthorized:
    return nil, customerror.TryParseAsCustomErr(&customerror.ErrUnauthorized, errors.New("token expired"))
case http.StatusForbidden:
    return nil, customerror.TryParseAsCustomErr(&customerror.ErrInsufficientPermission, errors.New("insufficient permissions"))
case http.StatusTooManyRequests:
    return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, errors.New("rate limited"))
default:
    return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, fmt.Errorf("API error %d", resp.StatusCode))
}
```

## 🎯 Use Case Layer

### Clean Architecture Integration
Pattern: Use case orchestrates OAuth flow ([auth.go](mdc:apps/backend/core-api/internal/usecase/auth/auth.go)):

```go
type AuthUsecase struct {
    googleOAuthService *oauth_services.GoogleOAuthService
}

func (u *AuthUsecase) VerifyGoogleOAuthCode(ctx context.Context, session *session.Session, code string, state string) (*oauth2.Token, *customerror.Err) {
    token, err := u.googleOAuthService.Callback(ctx, session, code, state)
    if err != nil {
        return nil, err.Extend("failed to verify google oauth code")
    }
    return token, nil
}
```

## 📋 Development Checklist

When implementing OAuth features:

**Security Requirements:**
- [ ] CSRF state validation implemented
- [ ] Secure session configuration
- [ ] Token validation with expiry checking
- [ ] Encrypted storage of OAuth references
- [ ] HTTPS-only cookies in production

**Error Handling:**
- [ ] Comprehensive HTTP status code handling
- [ ] Structured error responses using customerror
- [ ] Logging for debugging and monitoring
- [ ] Graceful degradation on API failures

**Testing:**
- [ ] OAuth flow integration tests
- [ ] State mismatch attack prevention
- [ ] Token expiry handling
- [ ] Rate limiting scenarios

**Configuration:**
- [ ] Environment-based OAuth credentials
- [ ] Separate development and production configs
- [ ] Proper redirect URL configuration
- [ ] Scope management for minimal permissions

## 🛠️ Development Commands

```bash
# Test OAuth endpoints
curl http://localhost:8080/api/v1/auth/request-google-oauth

# Database operations with OAuth
bun db:generate     # Regenerate after OAuth schema changes
bun compose:up      # Start database with encrypted storage

# API generation
bun gen-api         # Update TypeScript client with OAuth endpoints
```

## 🔗 Related Files

- **Service Layer**: [oauth.go](mdc:apps/backend/services/oauth/oauth.go), [google.go](mdc:apps/backend/services/oauth/google.go)
- **Handlers**: [auth/](mdc:apps/backend/core-api/internal/handler/auth/), [onboard/](mdc:apps/backend/core-api/internal/handler/onboard/)
- **Database**: [authentication_credentials.sql](mdc:packages/database/queries/authentication_credentials.sql)
- **Config**: [google.go](mdc:apps/backend/core-api/config/google/google.go)
- **Use Cases**: [auth.go](mdc:apps/backend/core-api/internal/usecase/auth/auth.go), [onboard.go](mdc:apps/backend/core-api/internal/usecase/onboard/onboard.go)

## 🔑 Key Principles

1. **Security First**: Always implement CSRF protection and secure session management
2. **Clean Architecture**: Separate concerns between handlers, use cases, and services
3. **Error Resilience**: Handle all possible OAuth failure scenarios gracefully
4. **Data Protection**: Encrypt sensitive OAuth data using established PII encryption patterns
5. **Configuration Management**: Use environment variables for all OAuth credentials
