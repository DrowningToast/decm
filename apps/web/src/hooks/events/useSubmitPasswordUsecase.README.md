# useSubmitPasswordUsecase Hook

## Overview

A custom React hook for handling password-protected event registration with a two-step flow: preview (password validation) and confirm (final registration with PII data).

**Automatically integrates with `EventPasswordNav` bottom navigation** - When the hook is initialized, it sets up a callback so that password submissions from the bottom nav trigger the preview function.

## Location

`apps/web/src/hooks/events/useSubmitPasswordUsecase.ts`

## Prerequisites

- The event must be password-required
- The user must be authenticated

## Store Integration

This hook automatically integrates with `useEventPasswordNavStore` to:

- Connect the `EventPasswordNav` bottom navigation with the password preview logic
- Manage password input state across the component
- Reset password field after successful operations

## Features

### 1. Preview Function

- Checks the password with the backend service
- Shows **destructive red toast** if the password is incorrect
- Shows **alert modal with registration requirements** if the password is correct
- Fetches and displays the event's registration configuration

### 2. Confirm Function

- Submits final registration with password and PII data
- Shows success toast on completion
- Handles validation errors

## API Reference

### Hook Signature

```typescript
function useSubmitPasswordUsecase(eventId: string): UseSubmitPasswordUsecaseReturn;
```

### Return Value

```typescript
interface UseSubmitPasswordUsecaseReturn {
    // Preview function - checks password
    preview: (password: string) => Promise<void>;

    // Confirm function - submits registration
    confirm: (password: string, piiData: ParticipantPIIData) => Promise<void>;

    // Loading states
    isPreviewLoading: boolean;
    isConfirmLoading: boolean;

    // Password validation state
    isPasswordValid: boolean;

    // Registration configuration to display
    registrationConfig: RegistrationConfig | null;

    // Reset function
    resetPasswordValidation: () => void;
}
```

### Types

#### ParticipantPIIData

```typescript
interface ParticipantPIIData {
    firstName?: string;
    lastName?: string;
    email?: string;
    bio?: string;
    phoneNumber?: string;
    address?: string;
    academicInstitution?: string;
    academicEmail?: string;
}
```

#### RegistrationConfig

```typescript
interface RegistrationConfig {
    requireFirstName?: boolean;
    requireLastName?: boolean;
    requireEmail?: boolean;
    requireBio?: boolean;
    requirePhoneNumber?: boolean;
    requireAddress?: boolean;
    requireAcademicInstitution?: boolean;
    requireAcademicEmail?: boolean;
}
```

## Usage Example

### Basic Usage with EventPasswordNav (Recommended)

```typescript
import { useSubmitPasswordUsecase } from "@/hooks/events/useSubmitPasswordUsecase";

function EventDetailPage({ eventId }: { eventId: string }) {
    const {
        isPasswordValid,
        registrationConfig,
        isPreviewLoading,
        confirm,
        resetPasswordValidation,
    } = useSubmitPasswordUsecase(eventId);

    // That's it! The hook automatically connects to EventPasswordNav
    // When user submits password from bottom nav, preview() is called

    const handleConfirm = async () => {
        const piiData = {
            firstName: "John",
            lastName: "Doe",
            email: "john@example.com",
        };

        // Get password from store
        const { password } = useEventPasswordNavStore.getState();
        await confirm(password, piiData);
    };

    return (
        <div>
            {/* Your event details UI */}

            {/* Alert dialog shows when isPasswordValid is true */}
            <AlertDialog open={isPasswordValid} onOpenChange={(open) => !open && resetPasswordValidation()}>
                <AlertDialogContent>
                    {/* Display requirements from registrationConfig */}
                    {/* Collect PII data */}
                    <AlertDialogAction onClick={handleConfirm}>
                        Confirm
                    </AlertDialogAction>
                </AlertDialogContent>
            </AlertDialog>

            {/* EventPasswordNav is rendered in BottomContainer based on event state */}
        </div>
    );
}
```

### Manual Usage (Without EventPasswordNav)

```typescript
import { useSubmitPasswordUsecase } from "@/hooks/events/useSubmitPasswordUsecase";
import { useState } from "react";

function EventRegistration({ eventId }: { eventId: string }) {
    const [password, setPassword] = useState("");
    const {
        preview,
        confirm,
        isPreviewLoading,
        isPasswordValid,
        registrationConfig,
    } = useSubmitPasswordUsecase(eventId);

    const handleCheckPassword = async () => {
        await preview(password);
    };

    const handleConfirm = async () => {
        const piiData = {
            firstName: "John",
            lastName: "Doe",
            email: "john@example.com",
        };
        await confirm(password, piiData);
    };

    return (
        <div>
            <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />
            <button onClick={handleCheckPassword} disabled={isPreviewLoading}>
                Check Password
            </button>

            {isPasswordValid && (
                <button onClick={handleConfirm}>
                    Confirm Registration
                </button>
            )}
        </div>
    );
}
```

### Complete Example with Alert Modal

See `useSubmitPasswordUsecase.example.tsx` for a complete implementation with:

- Password input field
- Alert dialog with registration requirements
- PII data form
- Loading states
- Error handling
- i18n translations

## Flow Diagram

### With EventPasswordNav (Automatic)

```
Hook initializes
        ↓
useEffect sets callback in store
        ↓
User types password in EventPasswordNav
        ↓
User clicks submit in bottom nav
        ↓
Store triggers callback → preview(password)
        ↓
    ┌───────────────────────┐
    │   Check Password      │
    │   with Backend        │
    └───────────────────────┘
            ↓
    ┌───────────────┐
    │  Incorrect?   │
    └───────────────┘
      ↓Yes       ↓No
      ↓          ↓
 Destructive   Fetch Config
 Red Toast          ↓
                Show Alert
                Modal with
                Requirements
                    ↓
            User fills PII
                    ↓
        Calls confirm(password, piiData)
                    ↓
            Submit Registration
                    ↓
            Success Toast
                    ↓
        Password field reset
```

## Toast Messages

### Error Toasts (Destructive - Red)

- **Incorrect Password**: Shows when password validation fails
- **Not Found**: Shows when event doesn't exist
- **Generic Error**: Shows for unexpected errors

### Success Toast

- **Registration Success**: Shows when final registration completes

## Translation Keys

### Required i18n Keys

```json
{
    "errors": {
        "invalidInput": "Invalid input",
        "unauthorized": "You are not authorized",
        "notFound": "The requested resource was not found",
        "generic": "Something went wrong"
    },
    "validation": {
        "passwordRequired": "Password is required"
    },
    "event": {
        "notFound": "Event not found",
        "registration": {
            "incorrectPassword": "Incorrect password. Please try again.",
            "invalidData": "Please check your information and try again.",
            "success": "Registration successful!",
            "successDescription": "You have successfully registered for this event."
        }
    }
}
```

## Backend Integration

### Service Used

```typescript
EventRegistrationService
- checkPassword(eventId: string, password: string)
- getConfiguration(eventId: string)
```

### API Endpoints

- `POST /api/v1/events/:eventId/check-password` - Validate password
- `GET /api/v1/events/:eventId/registration-config` - Get requirements

## Error Handling

### HTTP Status Codes

- **401/403**: Incorrect password → Shows destructive toast
- **404**: Event not found → Shows destructive toast
- **400**: Invalid data in confirm → Shows destructive toast
- **500**: Server error → Shows generic error toast

### Axios Errors

The hook automatically handles `AxiosError` instances and displays appropriate error messages based on the response status.

## State Management

### Internal State

- `isPreviewLoading`: Loading state for password check
- `isConfirmLoading`: Loading state for final submission
- `isPasswordValid`: Controls alert modal visibility
- `registrationConfig`: Stores fetched registration requirements

### Store Integration (useEventPasswordNavStore)

The hook automatically integrates with the password store:

- **On mount**: Sets `onSubmitCallback` to trigger `preview()` function
- **On success/reset**: Calls `resetPassword()` to clear the password field
- **EventPasswordNav** component reads from this store and triggers the callback

### State Flow

1. Hook initializes → Sets callback in store
2. User types password in EventPasswordNav → Store updates `password` state
3. User clicks submit → Store triggers `onSubmitCallback(password)`
4. Callback executes → `preview()` is called → `isPreviewLoading = true`
5. Password validated → `isPasswordValid = true`, `registrationConfig` set
6. Modal opens automatically when `isPasswordValid = true`
7. User fills PII and confirms → `isConfirmLoading = true`
8. Success → Modal closes, states reset, password field cleared

## Best Practices

### ✅ DO

- Always validate password input before calling `preview()`
- Check `isPreviewLoading` before allowing re-submission
- Use `isPasswordValid` to control modal visibility
- Call `resetPasswordValidation()` when modal closes
- Provide loading states in UI
- Handle all error cases with appropriate messages

### ❌ DON'T

- Don't call `preview()` with empty password
- Don't bypass password check and call `confirm()` directly
- Don't forget to reset state after successful registration
- Don't show modal without checking `isPasswordValid`
- Don't ignore loading states in UI

## Testing Considerations

### Mock Data

The `confirm()` function currently uses mock implementation:

```typescript
// TODO: Replace with actual registration API call
console.log("[Mock] Submitting registration with:", {
    eventId,
    password,
    piiData,
});
```

### Testing Checklist

- [ ] Test with correct password
- [ ] Test with incorrect password
- [ ] Test with empty password
- [ ] Test with network errors
- [ ] Test with non-existent event ID
- [ ] Test modal open/close flow
- [ ] Test form submission with PII data
- [ ] Test loading states
- [ ] Test error toast displays
- [ ] Test success toast displays

## Future Enhancements

### TODO Items

1. Replace mock `confirm()` implementation with actual API call
2. Transform `registrationConfig` based on actual API response structure
3. Add retry mechanism for network failures
4. Add analytics tracking for registration flow
5. Add validation for PII data before submission
6. Cache registration config to avoid refetching
7. Add support for partial PII data (optional fields)

## Related Files

- Service: `apps/web/src/services/EventRegistration.ts`
- Example: `apps/web/src/hooks/events/useSubmitPasswordUsecase.example.tsx`
- Translations: `apps/web/src/lib/i18n/locales/en.json`, `th.json`
- Components: `apps/web/src/components/ui/alert-dialog.tsx`

## Dependencies

```json
{
    "react": "^19.x",
    "react-i18next": "^13.x",
    "sonner": "^1.x",
    "axios": "^1.x"
}
```

## Author Notes

- The hook follows the project's usecase pattern
- Toast messages use Sonner library with destructive variant for errors
- Alert modal uses shadcn/ui AlertDialog component
- All text is internationalized with react-i18next
- Error handling follows project's error handling patterns
