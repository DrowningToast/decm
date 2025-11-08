# Event Detail Page - Usage Examples

## Quick Start

### 1. Basic Usage in Route

```tsx
// In your route file (e.g., pages/events/[id].tsx)
import { EventDetailPage } from "@/components/pages/Participant/Events/EventDetailPage";

export default function EventDetailRoute() {
    const { id } = useParams(); // or however you get the event ID

    return <EventDetailPage eventId={id} />;
}
```

## Testing Different Event Types

### Password-Required Event (Event ID: "1")

**Test Password**: `decm2024`

```tsx
// Navigate to event with ID "1"
<EventDetailPage eventId="1" />
```

**Expected Behavior**:

1. Shows event details
2. Bottom nav shows password input field
3. Enter password "decm2024" and submit
4. Toast: "Password correct! You have joined the event."
5. Bottom nav changes to "Participating"
6. Participation status shows "Joined"

**Try Wrong Password**:

1. Enter wrong password (e.g., "wrong123")
2. Toast: "Incorrect password - Please try again"
3. Stay on password input field

### Invitation-Only Event - Invited (Event ID: "2")

```tsx
<EventDetailPage eventId="2" />
```

**Expected Behavior**:

1. Shows event details
2. Bottom nav shows "You're invited! Click here to continue."
3. Click on the message
4. Toast: "Invitation accepted! You have successfully joined the event."
5. Bottom nav changes to "Participating"
6. Participation status shows "Accepted"

### Invitation-Only Event - Not Invited

To test this, modify the mock data in `useEventDetailUsecase.ts`:

```typescript
"2": {
    // ... other fields
    invitationStatus: "not-invited",  // Change from "invited"
}
```

**Expected Behavior**:

1. Shows event details
2. Bottom nav shows "Invitation is required." (read-only)
3. No action available

### Closed Event (Event ID: "3")

```tsx
<EventDetailPage eventId="3" />
```

**Expected Behavior**:

1. Shows event details
2. Status shows "No longer accepting request"
3. No bottom nav displayed

## Mock Data Structure

### Complete Event Detail Example

```typescript
{
    id: "1",
    name: "ToBelT69 - Password Event",
    description: "This is a password-protected event...",
    eventName: "ToBelT69",
    dateTime: "2024-09-24",
    finalCallDate: "2024-09-24",
    status: "accepting",  // or "closed"
    accessType: "password",  // or "invite-only" or "public"
    requiresPassword: true,
    correctPassword: "decm2024",  // Mock password
    hasJoined: false,  // Changes to true after joining
    seatsAvailable: 50,
    totalSeats: 100,
}
```

### Invitation-Only Event Example

```typescript
{
    id: "2",
    name: "ToBelT69 - Invitation Only",
    description: "This is an invitation-only event...",
    eventName: "ToBelT69",
    dateTime: "2024-09-25",
    finalCallDate: "2024-09-25",
    status: "accepting",
    accessType: "invite-only",
    invitationStatus: "invited",  // or "not-invited" or "accepted"
    seatsAvailable: 30,
    totalSeats: 50,
}
```

## Creating New Mock Events

### Add a New Password Event

In `useEventDetailUsecase.ts`:

```typescript
const mockEventDetails: Record<string, EventDetail> = {
    // ... existing events
    "4": {
        id: "4",
        name: "New Password Event",
        description: "Another password-protected event",
        eventName: "New Event",
        dateTime: "2024-10-01",
        finalCallDate: "2024-10-01",
        status: "accepting",
        accessType: "password",
        requiresPassword: true,
        correctPassword: "test123", // Your test password
        hasJoined: false,
        seatsAvailable: 100,
        totalSeats: 200,
    },
};
```

### Add a New Invitation Event

```typescript
const mockEventDetails: Record<string, EventDetail> = {
    // ... existing events
    "5": {
        id: "5",
        name: "VIP Only Event",
        description: "Exclusive invitation-only event",
        eventName: "VIP Event",
        dateTime: "2024-10-05",
        finalCallDate: "2024-10-05",
        status: "accepting",
        accessType: "invite-only",
        invitationStatus: "invited", // Change as needed
        seatsAvailable: 20,
        totalSeats: 30,
    },
};
```

## Testing All States

### Password Event States

```typescript
// State 1: Need to enter password
{
    hasJoined: false,
    // Shows: event-password variant
}

// State 2: Already joined
{
    hasJoined: true,
    // Shows: participating variant
}
```

### Invitation Event States

```typescript
// State 1: Not invited
{
    invitationStatus: "not-invited",
    // Shows: invitation-required variant
}

// State 2: Invited but not accepted
{
    invitationStatus: "invited",
    // Shows: invited variant (clickable)
}

// State 3: Invitation accepted
{
    invitationStatus: "accepted",
    // Shows: participating variant
}
```

## Toast Messages Reference

### Success Messages

```typescript
// Password correct
toast.success("Password correct! You have joined the event.", {
    description: "You can now participate in this event.",
});

// Invitation accepted
toast.success("Invitation accepted!", {
    description: "You have successfully joined the event.",
});
```

### Error Messages

```typescript
// Wrong password
toast.error("Incorrect password", {
    description: "Please try again with the correct password.",
});

// Failed to accept invitation
toast.error("Failed to accept invitation", {
    description: error.message,
});
```

## Bottom Nav Variants Reference

| Variant               | When Shown                    | User Action             |
| --------------------- | ----------------------------- | ----------------------- |
| `event-password`      | Password required, not joined | Enter password & submit |
| `invited`             | User is invited, not accepted | Click message to accept |
| `invitation-required` | User is not invited           | None (read-only)        |
| `participating`       | User has joined/accepted      | None (display status)   |
| No variant            | Event is closed               | None                    |

## Integration with Event List

The event list uses `useEventsListUsecase` which includes:

```typescript
const { events, isLoading, error } = useEventsListUsecase({
    filterType: "all", // or "my-events"
});
```

Click on an event card to navigate to detail:

```typescript
<EventCard
    event={event}
    onClick={() => navigate(`/events/${event.id}`)}
/>
```

## Common Patterns

### Loading State

```tsx
if (isLoading) {
    return (
        <Typography variant="text" tag="p" color="muted" className="animate-pulse">
            {t("common.loading")}
        </Typography>
    );
}
```

### Error State

```tsx
if (error || !event) {
    return (
        <Typography variant="text" tag="p" color="destructive">
            {t("errors.generic")}
        </Typography>
    );
}
```

### Conditional Bottom Nav

```tsx
{
    bottomNavVariant && (
        <BottomNav variant={bottomNavVariant} onBack={() => window.history.back()} />
    );
}
```

## Debugging Tips

### Check Store State

```tsx
// In EventDetailPage component, add console.logs:
const { password } = useEventPasswordNavStore();
console.log("Current password:", password);

const { event } = useEventDetailUsecase({ eventId });
console.log("Event state:", event);
console.log("Bottom nav variant:", bottomNavVariant);
```

### Test Callbacks

```tsx
// Add in useEffect to verify callbacks are set:
useEffect(() => {
    console.log("Password callback set");
    setOnSubmitCallback((password: string) => {
        console.log("Submitting password:", password);
        submitPassword({ password });
    });
}, [setOnSubmitCallback, submitPassword]);
```

### Check Toast Notifications

Make sure `Toaster` component is included in your app:

```tsx
// In your app root or layout
import { Toaster } from "sonner";

function App() {
    return (
        <>
            <YourContent />
            <Toaster />
        </>
    );
}
```

## Advanced Customization

### Custom Password Validation

Modify the mutation in `useEventDetailUsecase.ts`:

```typescript
mutationFn: async ({ password }: { password: string }) => {
    // Add custom validation logic
    if (password.length < 6) {
        throw new Error("Password must be at least 6 characters");
    }

    // Check against correct password
    if (password === event.correctPassword) {
        return { success: true, message: "Welcome!" };
    }

    throw new Error("Incorrect password");
};
```

### Custom Toast Duration

```typescript
toast.success("Message", {
    description: "Description",
    duration: 5000, // 5 seconds
});
```

### Redirect After Joining

```typescript
onSuccess: (data) => {
    toast.success(data.message);
    queryClient.invalidateQueries(["event-detail", eventId]);

    // Redirect to event page or dashboard
    setTimeout(() => {
        navigate("/events");
    }, 2000);
};
```

## Troubleshooting

### Password not submitting

- Check if `onSubmitCallback` is set in `useEventPasswordNavStore`
- Verify `setOnSubmitCallback` is called in `useEffect`
- Check if password field is not empty

### Toast not showing

- Ensure `Toaster` component is rendered
- Check browser console for errors
- Verify `sonner` package is installed

### Bottom nav not updating

- Check if `bottomNavVariant` is computing correctly
- Verify event state updates after mutations
- Check React Query cache invalidation

### Translations not working

- Verify translation keys exist in both `en.json` and `th.json`
- Check `useTranslation` hook is called
- Verify i18n is initialized in app root

## Further Reading

- [React Query Documentation](https://tanstack.com/query/latest)
- [Sonner Toast Documentation](https://sonner.emilkowal.ski/)
- [Zustand Store Documentation](https://zustand-demo.pmnd.rs/)
