# Role Guard Middleware

Role guard middleware for protecting routes based on user roles (verified organizer, verified issuer).

## Overview

The `role_guard` middleware checks if an authenticated user has the required role(s) to access specific endpoints. This middleware should be used **after** the `authentication_guard` middleware, as it requires the user claims to be injected into the context.

## Features

- ✅ `RequireRole(roles ...Role)` - Generic function that accepts multiple roles (user needs ONE OR MORE)
- ✅ `RequireVerifiedOrganizer()` - Convenience method for organizer-only routes
- ✅ `RequireVerifiedIssuer()` - Convenience method for issuer-only routes
- ✅ `RequireVerifiedOrganizerOrIssuer()` - Convenience method for organizer OR issuer routes

## Available Roles

```go
const (
    RoleVerifiedOrganizer Role = "verified_organizer"
    RoleVerifiedIssuer    Role = "verified_issuer"
)
```

## Installation

```go
import (
    roleguard "apps/backend/core-api/internal/middleware/role_guard"
    "apps/backend/services/auth"
)

// Initialize the middleware
authService := auth.NewAuthService(issuer, secretKey, expiration)
roleGuardMiddleware := roleguard.New(authService)
```

## Usage Examples

### Example 1: Using Generic RequireRole() with Multiple Roles

```go
import (
    authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
    roleguard "apps/backend/core-api/internal/middleware/role_guard"
)

func (h *Handler) Mount(r fiber.Router) {
    authGuard := h.AuthenticationGuardMiddleware
    roleGuard := h.RoleGuardMiddleware

    eventGroup := r.Group("/events").Use(authGuard.Middleware)

    // Single role - organizer only
    eventGroup.Post("/",
        roleGuard.RequireRole(roleguard.RoleVerifiedOrganizer),
        h.CreateEvent,
    )

    // Multiple roles - organizer OR issuer (user needs at least one)
    eventGroup.Get("/admin/dashboard",
        roleGuard.RequireRole(
            roleguard.RoleVerifiedOrganizer,
            roleguard.RoleVerifiedIssuer,
        ),
        h.GetAdminDashboard,
    )

    // Any authenticated user
    eventGroup.Get("/:event_id", h.GetEventById)
}
```

### Example 2: Using Convenience Methods

```go
func (h *Handler) Mount(r fiber.Router) {
    authGuard := h.AuthenticationGuardMiddleware
    roleGuard := h.RoleGuardMiddleware

    eventGroup := r.Group("/events").Use(authGuard.Middleware)

    // Organizer-only (convenience method)
    eventGroup.Post("/",
        roleGuard.RequireVerifiedOrganizer(),
        h.CreateEvent,
    )

    // Organizer or Issuer (convenience method)
    eventGroup.Get("/admin",
        roleGuard.RequireVerifiedOrganizerOrIssuer(),
        h.GetAdmin,
    )
}
```

### Example 3: Protecting Certificate Issuance (Issuers Only)

```go
func (h *Handler) Mount(r fiber.Router) {
    certGroup := r.Group("/certificates").Use(
        h.AuthenticationGuardMiddleware.Middleware,
    )

    // Issue certificate - requires verified issuer only
    certGroup.Post("/issue",
        h.RoleGuardMiddleware.RequireRole(roleguard.RoleVerifiedIssuer),
        h.IssueCertificate,
    )
}
```

## Integration with Existing Code

### Update Handler Struct

Add the role guard middleware to your handler:

```go
type Handler struct {
    AuthenticationService          *auth.AuthService
    AuthenticationGuardMiddleware  *authenticationguard.AuthenticationGuardMiddleware
    RoleGuardMiddleware            *roleguard.RoleGuardMiddleware  // Add this
    EventUc                        *event_usecase.EventUseCase
}

func NewHandler(
    eventUc *event_usecase.EventUseCase,
    authService *auth.AuthService,
    authGuard *authenticationguard.AuthenticationGuardMiddleware,
    roleGuard *roleguard.RoleGuardMiddleware,  // Add this parameter
) *Handler {
    return &Handler{
        AuthenticationService:         authService,
        AuthenticationGuardMiddleware: authGuard,
        RoleGuardMiddleware:           roleGuard,  // Add this
        EventUc:                       eventUc,
    }
}
```

### Update main.go

Initialize and inject the role guard middleware:

```go
import (
    roleguard "apps/backend/core-api/internal/middleware/role_guard"
)

func main() {
    // ... existing setup ...

    authenticationGuardMiddleware := authenticationguard.New(authService)
    roleGuardMiddleware := roleguard.New(authService)  // Add this

    // Update handler initialization
    eventHandler := event.NewHandler(
        eventUc,
        authService,
        authenticationGuardMiddleware,
        roleGuardMiddleware,  // Add this
    )
    eventHandler.Mount(apiV1)
}
```

## Error Responses

When a user doesn't have the required role, the middleware returns:

```json
{
    "error": {
        "code": 403,
        "message": "Forbidden",
        "details": "user is not a verified organizer"
    }
}
```

Possible error messages:

- `"user is not a verified organizer"`
- `"user is not a verified issuer"`
- `"user is not a verified organizer or issuer"`

## Middleware Order

⚠️ **Important**: Always use role guard middleware **after** authentication guard:

```go
// ✅ Correct order
eventGroup := r.Group("/events").Use(
    authenticationGuard.Middleware,  // 1. Authenticate first
)
eventGroup.Post("/",
    roleGuard.RequireVerifiedOrganizer(),  // 2. Then check role
    h.CreateEvent,
)

// ❌ Wrong - role guard before authentication
eventGroup.Post("/",
    roleGuard.RequireVerifiedOrganizer(),  // This will fail!
    authenticationGuard.Middleware,
    h.CreateEvent,
)
```

## Migration from Usecase-Level Checks

If you currently have role checks in your usecase layer, you can migrate them to middleware:

### Before (Usecase-level check):

```go
// In usecase/event/create_event.go
func (u *EventUsecase) CreateEvent(ctx context.Context, params CreateEventParams, user *auth.JwtClaims) (*entity.Event, error) {
    // Role check in usecase
    if user.IsVerifiedOrganizer == nil || !*user.IsVerifiedOrganizer {
        return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("user is not a verified organizer"))
    }

    // Business logic...
}
```

### After (Middleware-level check):

```go
// In handler/event/routes.go
eventGroup.Post("/",
    h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // Check at route level
    h.CreateEvent,
)

// In usecase/event/create_event.go
func (u *EventUsecase) CreateEvent(ctx context.Context, params CreateEventParams, user *auth.JwtClaims) (*entity.Event, error) {
    // Role already verified by middleware, focus on business logic
    // Business logic...
}
```

## Benefits

1. **Separation of Concerns**: Authentication/authorization logic separated from business logic
2. **Reusability**: Same middleware can be applied to multiple routes
3. **Clarity**: Route protection is explicit and visible in route definitions
4. **Consistency**: Standardized error responses for authorization failures
5. **Performance**: Early rejection of unauthorized requests before reaching usecase layer

## Testing

Example test for role guard middleware:

```go
func TestRequireVerifiedOrganizer(t *testing.T) {
    authService := auth.NewAuthService("test", "secret", time.Hour)
    roleGuard := roleguard.New(authService)

    app := fiber.New()
    app.Get("/test",
        func(c *fiber.Ctx) error {
            // Mock user in context
            isVerified := true
            c.Locals("user", &auth.JwtClaims{
                IsVerifiedOrganizer: &isVerified,
            })
            return c.Next()
        },
        roleGuard.RequireVerifiedOrganizer(),
        func(c *fiber.Ctx) error {
            return c.SendStatus(200)
        },
    )

    req := httptest.NewRequest("GET", "/test", nil)
    resp, _ := app.Test(req)

    assert.Equal(t, 200, resp.StatusCode)
}
```

## Related Files

- `apps/backend/core-api/internal/middleware/authentication_guard/` - Authentication middleware
- `apps/backend/services/auth/jwt.go` - JWT claims structure
- `apps/backend/core-api/internal/entity/authentication_credentials.go` - User roles entity
