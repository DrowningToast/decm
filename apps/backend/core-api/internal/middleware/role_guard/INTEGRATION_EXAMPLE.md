# Role Guard Middleware - Complete Integration Example

This document shows a complete example of integrating the role guard middleware into the DECM backend.

## Step 1: Update main.go

Add the role guard middleware initialization:

```go
package main

import (
    // ... existing imports ...
    roleguard "apps/backend/core-api/internal/middleware/role_guard"
)

func main() {
    // ... existing setup code ...

    // Initialize middlewares
    authenticationGuardMiddleware := authenticationguard.New(authService)
    verifyJwtMiddleware := verifyjwt.New(authService)
    roleGuardMiddleware := roleguard.New(authService)  // 👈 Add this line

    // Update handler initializations to include roleGuardMiddleware

    // Event Handler
    eventHandler := event.NewHandler(
        eventUc,
        eventConfigUc,
        profileUc,
        eventRegistrationInvitationUc,
        authService,
        authenticationGuardMiddleware,
        roleGuardMiddleware,  // 👈 Add this parameter
        logger,
    )
    eventHandler.Mount(apiV1)

    // Issuer Handler
    issuerHandler := issuer.NewHandler(
        issuerUc,
        authService,
        authenticationGuardMiddleware,
        roleGuardMiddleware,  // 👈 Add this parameter
    )
    issuerHandler.Mount(apiV1)

    // ... rest of the code ...
}
```

## Step 2: Update Event Handler Structure

Update `apps/backend/core-api/internal/handler/event/handler.go`:

```go
package event

import (
    "log/slog"

    "apps/backend/services/auth"
    authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
    roleguard "apps/backend/core-api/internal/middleware/role_guard"  // 👈 Add import
    event_registration_invitation_usecase "apps/backend/core-api/internal/usecase/event_registration_invitation"
    eventconfig "apps/backend/core-api/internal/usecase/eventconfig"
    "apps/backend/core-api/internal/usecase/event"
    profile "apps/backend/core-api/internal/usecase/profile"
)

type Handler struct {
    EventUc                                *event.EventUseCase
    EventConfigUc                          *eventconfig.EventConfigUseCase
    ProfileUc                              *profile.ProfileUsecase
    EventRegistrationInvitationUc          *event_registration_invitation_usecase.EventRegistrationInvitationUseCase
    AuthenticationService                  *auth.AuthService
    AuthenticationGuardMiddleware          *authenticationguard.AuthenticationGuardMiddleware
    RoleGuardMiddleware                    *roleguard.RoleGuardMiddleware  // 👈 Add this field
    Logger                                 *slog.Logger
}

func NewHandler(
    eventUc *event.EventUseCase,
    eventConfigUc *eventconfig.EventConfigUseCase,
    profileUc *profile.ProfileUsecase,
    eventRegistrationInvitationUc *event_registration_invitation_usecase.EventRegistrationInvitationUseCase,
    authenticationService *auth.AuthService,
    authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware,
    roleGuardMiddleware *roleguard.RoleGuardMiddleware,  // 👈 Add this parameter
    logger *slog.Logger,
) *Handler {
    return &Handler{
        EventUc:                                eventUc,
        EventConfigUc:                          eventConfigUc,
        ProfileUc:                              profileUc,
        EventRegistrationInvitationUc:          eventRegistrationInvitationUc,
        AuthenticationService:                  authenticationService,
        AuthenticationGuardMiddleware:          authenticationGuardMiddleware,
        RoleGuardMiddleware:                    roleGuardMiddleware,  // 👈 Add this field
        Logger:                                 logger,
    }
}
```

## Step 3: Update Event Routes

Update `apps/backend/core-api/internal/handler/event/routes.go`:

```go
package event

import (
    "apps/backend/common/log"
    "apps/backend/core-api/internal/entity"
    "github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
    logger := log.LoadLogger()
    defer logger.Info("Mounted event routes")

    // Base event group - all routes require authentication
    eventGroup := r.Group("/events").Use(
        h.AuthenticationGuardMiddleware.Middleware,
    )

    // ====== PUBLIC ROUTES (Authenticated users) ======
    // View event details - any authenticated user
    eventGroup.Get("/:event_id", h.GetEventById)
    eventGroup.Get("/:event_id/contracts", h.GetEventContractByEventID)
    eventGroup.Get("/:event_id/issuers", h.GetEventIssuersByEventID)
    eventGroup.Get("/:event_id/registration/invitations", h.GetEventRegistrationInvitationsByEventID)

    // View own events - any authenticated user
    eventGroup.Get("/owner-credentials/:owner_credential_id", h.GetEventsByOwnerCredentialsId)

    // ====== ORGANIZER-ONLY ROUTES ======
    // Create event - requires verified organizer
    eventGroup.Post("/",
        h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // 👈 Add role guard
        h.CreateEvent,
    )

    // Update event - requires verified organizer
    eventGroup.Put("/:event_id",
        h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // 👈 Add role guard
        h.UpdateEvent,
    )

    // Delete event - requires verified organizer
    eventGroup.Delete("/:event_id",
        h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // 👈 Add role guard
        h.DeleteEvent,
    )

    // Manage event contracts - requires verified organizer
    eventGroup.Put("/:event_id/contracts",
        h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // 👈 Add role guard
        h.UpdateEventContract,
    )
    eventGroup.Delete("/:event_id/contracts",
        h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // 👈 Add role guard
        h.DeleteEventContract,
    )

    // Manage event issuers - requires verified organizer
    eventGroup.Post("/:event_id/issuers",
        h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // 👈 Add role guard
        h.CreateEventIssuer,
    )
    eventGroup.Put("/:event_id/issuers",
        h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // 👈 Add role guard
        h.UpdateEventIssuer,
    )
    eventGroup.Delete("/:event_id/issuers/:issuer_id",
        h.RoleGuardMiddleware.RequireVerifiedOrganizer(),  // 👈 Add role guard
        h.DeleteEventIssuer,
    )
}
```

## Step 4: Simplify Usecase Logic (Optional)

Now that role checks are at the middleware level, you can simplify your usecases:

### Before (with role check in usecase):

```go
// apps/backend/core-api/internal/usecase/event/create_event.go
func (u *EventUsecase) CreateEvent(ctx context.Context, params CreateEventParams, user *auth.JwtClaims) (*entity.Event, error) {
    // Role validation
    credential, err := u.AuthenticationCredentialsDg.GetAuthenticationCredentialById(ctx, user.UserId)
    if err != nil {
        return nil, err
    }

    isVerifiedOrganizer := credential.IsVerifiedOrganizer
    if !isVerifiedOrganizer {
        return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("user is not a verified organizer"))
    }

    // Business logic continues...
}
```

### After (role check handled by middleware):

```go
// apps/backend/core-api/internal/usecase/event/create_event.go
func (u *EventUsecase) CreateEvent(ctx context.Context, params CreateEventParams, user *auth.JwtClaims) (*entity.Event, error) {
    // Role already verified by middleware - focus on business logic
    credential, err := u.AuthenticationCredentialsDg.GetAuthenticationCredentialById(ctx, user.UserId)
    if err != nil {
        return nil, err
    }

    // Business logic continues...
}
```

**Note**: You can keep the usecase-level checks for additional validation or remove them for cleaner code. The middleware provides the primary authorization layer.

## Step 5: Update Issuer Handler (Similar Pattern)

Update `apps/backend/core-api/internal/handler/issuer/handler.go`:

```go
package issuer

import (
    "apps/backend/services/auth"
    authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
    roleguard "apps/backend/core-api/internal/middleware/role_guard"  // 👈 Add import
    "apps/backend/core-api/internal/usecase/issuer"
)

type Handler struct {
    IssuerUc                       *issuer.IssuerUseCase
    AuthenticationService          *auth.AuthService
    AuthenticationGuardMiddleware  *authenticationguard.AuthenticationGuardMiddleware
    RoleGuardMiddleware            *roleguard.RoleGuardMiddleware  // 👈 Add field
}

func NewHandler(
    issuerUc *issuer.IssuerUseCase,
    authenticationService *auth.AuthService,
    authenticationGuardMiddleware *authenticationguard.AuthenticationGuardMiddleware,
    roleGuardMiddleware *roleguard.RoleGuardMiddleware,  // 👈 Add parameter
) *Handler {
    return &Handler{
        IssuerUc:                      issuerUc,
        AuthenticationService:         authenticationService,
        AuthenticationGuardMiddleware: authenticationGuardMiddleware,
        RoleGuardMiddleware:           roleGuardMiddleware,  // 👈 Initialize
    }
}
```

Update `apps/backend/core-api/internal/handler/issuer/routes.go`:

```go
func (h *Handler) Mount(r fiber.Router) {
    issuerGroup := r.Group("/issuers").Use(
        h.AuthenticationGuardMiddleware.Middleware,
    )

    // List verified issuers - any authenticated user can view
    issuerGroup.Get("/", h.GetVerifiedIssuers)

    // Issue certificates - requires verified issuer role
    issuerGroup.Post("/certificates",
        h.RoleGuardMiddleware.RequireVerifiedIssuer(),  // 👈 Add role guard
        h.IssueCertificate,
    )
}
```

## Testing the Integration

### Test with curl:

```bash
# 1. Login as a verified organizer
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"wallet_address": "0x...", "signature": "..."}'

# 2. Try to create an event (should succeed for verified organizer)
curl -X POST http://localhost:8080/api/v1/events \
  -H "Cookie: session=<jwt_token>" \
  -F "name=Test Event" \
  -F "start_date=2025-12-01T10:00:00Z" \
  ...

# 3. Try with non-organizer (should return 403 Forbidden)
# Login as regular user, then:
curl -X POST http://localhost:8080/api/v1/events \
  -H "Cookie: session=<non_organizer_token>" \
  -F "name=Test Event" \
  ...

# Expected response:
# {
#   "error": {
#     "code": 403,
#     "message": "Forbidden",
#     "details": "user is not a verified organizer"
#   }
# }
```

## Benefits of This Integration

1. **Clear Route Protection**: Looking at `routes.go` immediately shows which endpoints require special roles
2. **Reduced Boilerplate**: No need to repeat role checks in every usecase
3. **Early Rejection**: Unauthorized requests are rejected at the handler level before reaching business logic
4. **Consistent Errors**: All role-based authorization errors follow the same format
5. **Easy to Test**: Middleware can be tested independently of business logic

## Rollout Strategy

### Phase 1: Add Middleware (Non-Breaking)

- Add role guard middleware to project
- Keep existing usecase-level checks
- Test in development

### Phase 2: Apply to New Routes

- Use middleware for all new endpoints
- Gradually migrate existing routes

### Phase 3: Cleanup (Optional)

- Remove redundant usecase-level role checks
- Keep usecase checks only for complex business rules

This approach ensures backward compatibility while improving code organization.
