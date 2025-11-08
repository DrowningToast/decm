# Role Guard Middleware - Implementation Summary

## Overview

A flexible, generic role-based access control middleware for protecting routes based on user roles.

## Key Features

✅ **Generic `RequireRole()` function** - Accepts variadic roles, user must have at least ONE  
✅ **Role constants** - Type-safe role definitions  
✅ **Convenience methods** - Pre-configured shortcuts for common use cases  
✅ **Flexible & Extensible** - Easy to add new roles in the future

## Implementation

### Core Function

```go
// User must have ONE OR MORE of the specified roles
func (m *RoleGuardMiddleware) RequireRole(roles ...Role) fiber.Handler
```

### Available Roles

```go
const (
    RoleVerifiedOrganizer Role = "verified_organizer"
    RoleVerifiedIssuer    Role = "verified_issuer"
)
```

### Convenience Methods

All convenience methods use the generic `RequireRole()` internally:

```go
// Single role shortcuts
RequireVerifiedOrganizer()          // = RequireRole(RoleVerifiedOrganizer)
RequireVerifiedIssuer()             // = RequireRole(RoleVerifiedIssuer)

// Multiple roles shortcut
RequireVerifiedOrganizerOrIssuer()  // = RequireRole(RoleVerifiedOrganizer, RoleVerifiedIssuer)
```

## Usage Examples

### 1. Generic RequireRole() - Recommended Approach

```go
import roleguard "apps/backend/core-api/internal/middleware/role_guard"

func (h *Handler) Mount(r fiber.Router) {
    authGuard := h.AuthenticationGuardMiddleware
    roleGuard := h.RoleGuardMiddleware

    eventGroup := r.Group("/events").Use(authGuard.Middleware)

    // Single role
    eventGroup.Post("/",
        roleGuard.RequireRole(roleguard.RoleVerifiedOrganizer),
        h.CreateEvent,
    )

    // Multiple roles (OR logic - user needs at least one)
    eventGroup.Get("/admin/dashboard",
        roleGuard.RequireRole(
            roleguard.RoleVerifiedOrganizer,
            roleguard.RoleVerifiedIssuer,
        ),
        h.GetAdminDashboard,
    )
}
```

### 2. Convenience Methods - Quick Shortcuts

```go
// Organizer-only route
eventGroup.Post("/",
    h.RoleGuardMiddleware.RequireVerifiedOrganizer(),
    h.CreateEvent,
)

// Issuer-only route
certGroup.Post("/issue",
    h.RoleGuardMiddleware.RequireVerifiedIssuer(),
    h.IssueCertificate,
)

// Organizer OR Issuer route
adminGroup.Get("/dashboard",
    h.RoleGuardMiddleware.RequireVerifiedOrganizerOrIssuer(),
    h.GetDashboard,
)
```

## Integration Checklist

### 1. Initialize in main.go

```go
roleGuardMiddleware := roleguard.New(authService)
```

### 2. Add to Handler Struct

```go
type Handler struct {
    RoleGuardMiddleware *roleguard.RoleGuardMiddleware
    // ... other fields
}
```

### 3. Pass to Handler Constructor

```go
eventHandler := event.NewHandler(
    eventUc,
    authService,
    authenticationGuardMiddleware,
    roleGuardMiddleware,  // Add this
    logger,
)
```

### 4. Use in Routes

```go
eventGroup.Post("/",
    roleGuard.RequireRole(roleguard.RoleVerifiedOrganizer),
    h.CreateEvent,
)
```

## Error Responses

### Example Error Response

```json
{
    "error": {
        "code": 403,
        "message": "Forbidden",
        "details": "user does not have required role. Required one of: verified_organizer, verified_issuer"
    }
}
```

## Advantages Over Previous Approach

### Before (Usecase-level checks)

```go
// Repeated in every usecase method
func (u *EventUsecase) CreateEvent(...) {
    if user.IsVerifiedOrganizer == nil || !*user.IsVerifiedOrganizer {
        return customerror.Parse(&customerror.ErrForbidden, ...)
    }
    // Business logic...
}
```

### After (Middleware-level checks)

```go
// Declared once in routes
eventGroup.Post("/",
    roleGuard.RequireRole(roleguard.RoleVerifiedOrganizer),
    h.CreateEvent,
)

// Usecase focuses on business logic only
func (u *EventUsecase) CreateEvent(...) {
    // Business logic...
}
```

## Benefits

1. **DRY Principle** - No repeated role checks across usecases
2. **Separation of Concerns** - Authorization separated from business logic
3. **Clear Routes** - Role requirements visible in route definitions
4. **Type Safety** - Role constants prevent typos
5. **Flexible** - Easy to require multiple roles with OR logic
6. **Extensible** - Simple to add new roles (just add to the switch statement)
7. **Testable** - Middleware can be tested independently

## Future Extensions

To add a new role (e.g., `RoleVerifiedStudent`):

1. **Add constant:**

```go
const (
    RoleVerifiedOrganizer Role = "verified_organizer"
    RoleVerifiedIssuer    Role = "verified_issuer"
    RoleVerifiedStudent   Role = "verified_student"  // New role
)
```

2. **Add case in RequireRole:**

```go
case RoleVerifiedStudent:
    if claims.IsVerifiedStudent != nil && *claims.IsVerifiedStudent {
        hasRole = true
    }
```

3. **Use it:**

```go
studentGroup.Post("/submit",
    roleGuard.RequireRole(roleguard.RoleVerifiedStudent),
    h.SubmitAssignment,
)
```

## Testing

The middleware includes comprehensive tests:

- ✅ Single role requirements
- ✅ Multiple role requirements (OR logic)
- ✅ Missing roles (403 Forbidden)
- ✅ No user context (401 Unauthorized)
- ✅ Convenience methods
- ✅ Error message formatting

Run tests:

```bash
cd apps/backend
go test ./core-api/internal/middleware/role_guard/... -v
```

## Files Created

- `role_guard.go` - Main middleware implementation
- `role_guard_test.go` - Comprehensive test suite
- `README.md` - Detailed usage documentation
- `INTEGRATION_EXAMPLE.md` - Step-by-step integration guide
- `SUMMARY.md` - This summary document

## Conclusion

The role guard middleware provides a clean, flexible, and maintainable way to protect routes based on user roles. The generic `RequireRole()` function with variadic parameters makes it easy to handle both single and multiple role requirements with a simple, consistent API.
