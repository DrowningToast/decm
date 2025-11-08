# Event Detail Page - Usecase Hooks Documentation

## Overview

The Event Detail page uses the `useEventDetailUsecase` hook to manage event details, password-protected events, and invitation-only events with proper state management and toast notifications.

## Architecture

### Files Structure

```
Participant/Events/
├── EventDetailPage.tsx           # Main page component
├── useEventDetailUsecase.ts      # Business logic hook
├── useEventsListUsecase.ts       # Event list hook
└── README.md                     # This file
```

### Related Files

```
BottomNav/
├── stores/
│   ├── event-password.ts         # Password submission store
│   └── event-invitation.ts       # Invitation acceptance store
└── variants/
    ├── EventPasswordNav.tsx      # Password input UI
    ├── InvitedNav.tsx           # Invited status UI
    ├── InvitationRequiredNav.tsx # Not invited UI
    └── ParticipatingNav.tsx     # Participating status UI
```

## Usage

### Basic Implementation

```tsx
import { EventDetailPage } from "@/components/pages/Participant/Events/EventDetailPage";

// In your route component
<EventDetailPage eventId="1" />;
```

### Hook Usage

```tsx
import { useEventDetailUsecase } from "./useEventDetailUsecase";

const {
    event, // Event detail data
    isLoading, // Loading state
    error, // Error state
    submitPassword, // Function to submit password
    acceptInvitation, // Function to accept invitation
    bottomNavVariant, // Computed bottom nav variant
    hasJoinedPasswordEvent, // Password event joined state
    hasAcceptedInvitation, // Invitation accepted state
} = useEventDetailUsecase({ eventId: "1" });
```

## Event Types

### 1. Password-Required Events

**Access Type**: `password` or `requiresPassword: true`

**States**:

- `input-password`: User needs to enter password (shows `event-password` variant)
- `joined`: User has successfully entered the correct password (shows `participating` variant)

**Mock Data Example**:

```typescript
{
    id: "1",
    name: "Password Protected Event",
    accessType: "password",
    requiresPassword: true,
    correctPassword: "decm2024",  // For testing
    hasJoined: false,
}
```

**User Flow**:

1. User views event detail
2. Bottom nav shows password input field (`event-password` variant)
3. User enters password and clicks submit
4. System validates password:
    - ✅ **Correct**: Toast success message, updates to `participating` variant
    - ❌ **Incorrect**: Toast error message, stays on `event-password` variant
5. After joining, shows `participating` variant

**Toast Messages**:

- Success: "Password correct! You have joined the event."
- Error: "Incorrect password - Please try again with the correct password."

### 2. Invitation-Only Events

**Access Type**: `invite-only`

**States**:

- `not-invited`: User is not invited (shows `invitation-required` variant)
- `invited`: User is invited but hasn't accepted (shows `invited` variant - clickable)
- `accepted`: User has accepted the invitation (shows `participating` variant)

**Mock Data Example**:

```typescript
{
    id: "2",
    name: "Invitation Only Event",
    accessType: "invite-only",
    invitationStatus: "invited",  // or "not-invited" / "accepted"
}
```

**User Flow**:

**Scenario A: User is Invited**

1. User views event detail
2. Bottom nav shows "You're invited! Click here to continue." (`invited` variant)
3. User clicks on the message
4. System accepts invitation:
    - ✅ Toast success message, updates to `participating` variant
    - ❌ Toast error message if something fails
5. After accepting, shows `participating` variant

**Scenario B: User is Not Invited**

1. User views event detail
2. Bottom nav shows "Invitation is required." (`invitation-required` variant)
3. No action available - display only

**Toast Messages**:

- Success: "Invitation accepted! - You have successfully joined the event."
- Error: "Failed to accept invitation"

### 3. Closed Events

**Status**: `closed`

**Behavior**: No bottom nav displayed

## Bottom Nav Variants

The `bottomNavVariant` is automatically computed based on event state:

| Event Type  | State       | Variant               | Description               |
| ----------- | ----------- | --------------------- | ------------------------- |
| Password    | Not joined  | `event-password`      | Shows password input      |
| Password    | Joined      | `participating`       | Shows joined status       |
| Invite-only | Not invited | `invitation-required` | Shows message (read-only) |
| Invite-only | Invited     | `invited`             | Shows clickable message   |
| Invite-only | Accepted    | `participating`       | Shows accepted status     |
| Closed      | Any         | `undefined`           | No bottom nav             |

## Mock Data

### Testing Different States

The mock data includes different event scenarios:

```typescript
// Event ID "1" - Password-required event
{
    id: "1",
    accessType: "password",
    correctPassword: "decm2024",  // Use this to test
    hasJoined: false,
}

// Event ID "2" - Invited-only event (user is invited)
{
    id: "2",
    accessType: "invite-only",
    invitationStatus: "invited",
}

// Event ID "3" - Closed event
{
    id: "3",
    status: "closed",
}
```

### Modifying Mock Data

To test different states, modify the mock data in `useEventDetailUsecase.ts`:

```typescript
const mockEventDetails: Record<string, EventDetail> = {
    "1": {
        // ... modify states here
        hasJoined: true, // Test joined state
    },
    "2": {
        // ... modify invitation status
        invitationStatus: "not-invited", // Test not invited
    },
};
```

## State Management

### Password Store (`event-password.ts`)

```typescript
const {
    password, // Current password input
    setPassword, // Update password input
    resetPassword, // Clear password field
    onSubmitCallback, // Callback function for submission
    setOnSubmitCallback, // Set the callback
} = useEventPasswordNavStore();
```

### Invitation Store (`event-invitation.ts`)

```typescript
const {
    onAcceptCallback, // Callback function for acceptance
    setOnAcceptCallback, // Set the callback
} = useEventInvitationNavStore();
```

## Implementation Details

### Password Validation Flow

1. User types password in `EventPasswordNav` component
2. Password is stored in `useEventPasswordNavStore`
3. When user clicks submit, `onSubmitCallback` is triggered
4. `EventDetailPage` receives the password and calls `submitPassword` mutation
5. Mutation validates password against mock data (or API)
6. Toast notification is shown
7. On success, `hasJoined` is updated to `true`
8. Password field is reset
9. Bottom nav switches to `participating` variant

### Invitation Acceptance Flow

1. User sees `InvitedNav` component (clickable message)
2. When user clicks, `onAcceptCallback` is triggered
3. `EventDetailPage` calls `acceptInvitation` mutation
4. Mutation processes acceptance (mock or API)
5. Toast notification is shown
6. On success, `invitationStatus` is updated to `"accepted"`
7. Bottom nav switches to `participating` variant

## Toast Notifications

Using `sonner` for toast notifications:

```typescript
// Success toast
toast.success("Title", {
    description: "Description text",
});

// Error toast
toast.error("Title", {
    description: "Error description",
});
```

## API Integration (TODO)

When integrating with the real API, replace the mock mutations:

### Password Submit

```typescript
mutationFn: async ({ password }: { password: string }) => {
    // Replace this
    const response = await api.joinPasswordProtectedEvent(eventId, password);
    return response.data;
};
```

### Accept Invitation

```typescript
mutationFn: async () => {
    // Replace this
    const response = await api.acceptEventInvitation(eventId);
    return response.data;
};
```

## Translations

All user-facing text uses i18n translations:

### English (`en.json`)

```json
{
    "participant": {
        "events": {
            "detail": {
                "passwordRequired": "Password Required",
                "inviteOnly": "Invite only",
                "joined": "Joined",
                "accepted": "Accepted"
            },
            "invited": "You're invited! Click here to continue.",
            "invitationRequired": "Invitation is required.",
            "passwordPlaceholder": "Password here"
        }
    }
}
```

### Thai (`th.json`)

```json
{
    "participant": {
        "events": {
            "detail": {
                "passwordRequired": "ต้องใช้รหัสผ่าน",
                "inviteOnly": "เฉพาะผู้ได้รับเชิญ",
                "joined": "เข้าร่วมแล้ว",
                "accepted": "ยอมรับแล้ว"
            },
            "invited": "คุณได้รับเชิญแล้ว! คลิกที่นี่เพื่อดำเนินการต่อ",
            "invitationRequired": "ต้องการคำเชิญ",
            "passwordPlaceholder": "กรอกรหัสผ่านที่นี่"
        }
    }
}
```

## Best Practices

1. **Always use the usecase hook** - Don't manage event state directly in the component
2. **Reset password field** - Clear the password after submission for security
3. **Show loading states** - Use `isLoading` to display loading UI
4. **Handle errors gracefully** - Show error toast and keep UI interactive
5. **Use translations** - Never hard-code user-facing text
6. **Test all states** - Verify all event types and states work correctly

## Testing Checklist

- [ ] Password-required event - correct password
- [ ] Password-required event - incorrect password
- [ ] Password-required event - already joined state
- [ ] Invite-only event - user is invited
- [ ] Invite-only event - user is not invited
- [ ] Invite-only event - user has accepted
- [ ] Closed event - no bottom nav shown
- [ ] Loading state displays correctly
- [ ] Error state displays correctly
- [ ] Toast notifications appear correctly
- [ ] Translations work in both English and Thai

## Future Enhancements

1. Add support for public events
2. Add support for shortlisted status
3. Implement actual API calls
4. Add event capacity checking
5. Add event date/time validation
6. Add event cancellation flow
7. Add waitlist functionality
