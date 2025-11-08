# EventPasswordNavStore Integration Guide

## Overview

The `useSubmitPasswordUsecase` hook is now **fully integrated** with the `useEventPasswordNavStore` for seamless password submission from the bottom navigation.

## How It Works

### Automatic Connection

When you initialize the hook, it automatically:

1. **Sets up a callback** in the store via `useEffect`
2. **Listens for password submissions** from `EventPasswordNav`
3. **Triggers the preview function** when user submits
4. **Resets the password field** after success or manual reset

### Architecture

```
┌─────────────────────────────────────────┐
│   ParticipantEventDetailPage.tsx       │
│                                         │
│   useSubmitPasswordUsecase(eventId)   │
│   └─ useEffect → setOnSubmitCallback  │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│   useEventPasswordNavStore (Zustand)   │
│                                         │
│   - password: string                   │
│   - setPassword()                      │
│   - resetPassword()                    │
│   - onSubmitCallback()  ← SET HERE    │
└─────────────────────────────────────────┘
                    ↑          ↓
                    │          │
         User types │          │ User clicks submit
                    │          │
┌─────────────────────────────────────────┐
│   EventPasswordNav.tsx                 │
│   (Bottom Navigation Component)        │
│                                         │
│   <input onChange={setPassword} />     │
│   <button onClick={onSubmitCallback} />│
└─────────────────────────────────────────┘
```

## Usage Example

### In ParticipantEventDetailPage.tsx

```typescript
import { useSubmitPasswordUsecase } from "@/hooks/events/useSubmitPasswordUsecase";
import { useEventPasswordNavStore } from "@/components/BottomNav/stores/event-password";

export function ParticipantEventDetailPage() {
    const { id: eventId } = useParams();

    // Initialize the hook - that's all you need!
    const {
        isPasswordValid,
        registrationConfig,
        isPreviewLoading,
        confirm,
        resetPasswordValidation,
    } = useSubmitPasswordUsecase(eventId);

    // Handle final registration confirmation
    const handleConfirm = async () => {
        // Get the password from the store
        const { password } = useEventPasswordNavStore.getState();

        // Mock PII data (replace with actual form data)
        const piiData = {
            firstName: "John",
            lastName: "Doe",
            email: "john@example.com",
        };

        await confirm(password, piiData);
    };

    return (
        <div>
            {/* Your event details UI */}

            {/* Alert modal - opens automatically when password is valid */}
            <AlertDialog
                open={isPasswordValid}
                onOpenChange={(open) => !open && resetPasswordValidation()}
            >
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Registration Requirements</AlertDialogTitle>
                        <AlertDialogDescription>
                            Please provide the following information
                        </AlertDialogDescription>
                    </AlertDialogHeader>

                    {/* Display requirements */}
                    {registrationConfig && (
                        <div className="space-y-2">
                            {registrationConfig.requireFirstName && (
                                <p>• First Name (Required)</p>
                            )}
                            {registrationConfig.requireLastName && (
                                <p>• Last Name (Required)</p>
                            )}
                            {registrationConfig.requireEmail && (
                                <p>• Email (Required)</p>
                            )}
                        </div>
                    )}

                    {/* PII Data Form */}
                    <div className="space-y-2">
                        {/* Add your form fields here */}
                    </div>

                    <AlertDialogFooter>
                        <AlertDialogCancel onClick={resetPasswordValidation}>
                            Cancel
                        </AlertDialogCancel>
                        <AlertDialogAction onClick={handleConfirm}>
                            Confirm Registration
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>

            {/* EventPasswordNav renders automatically in BottomContainer */}
        </div>
    );
}
```

## Flow Diagram

```
┌────────────────────────────────────────────────┐
│ 1. Component mounts                            │
│    useSubmitPasswordUsecase(eventId)          │
└────────────────────────────────────────────────┘
                    ↓
┌────────────────────────────────────────────────┐
│ 2. useEffect runs                              │
│    setOnSubmitCallback((password) => {        │
│        preview(password)                       │
│    })                                          │
└────────────────────────────────────────────────┘
                    ↓
┌────────────────────────────────────────────────┐
│ 3. User types password in EventPasswordNav    │
│    Store updates: password = "user-input"     │
└────────────────────────────────────────────────┘
                    ↓
┌────────────────────────────────────────────────┐
│ 4. User clicks submit button in bottom nav    │
│    EventPasswordNav calls onSubmitCallback()  │
└────────────────────────────────────────────────┘
                    ↓
┌────────────────────────────────────────────────┐
│ 5. Callback triggers preview(password)        │
│    - Shows loading state                       │
│    - Calls EventRegistrationService           │
│    - Validates password with backend          │
└────────────────────────────────────────────────┘
                    ↓
         ┌──────────┴──────────┐
         ↓                     ↓
┌─────────────────┐   ┌─────────────────┐
│ Incorrect       │   │ Correct         │
│ Password        │   │ Password        │
│                 │   │                 │
│ • Destructive   │   │ • Fetch config  │
│   toast error   │   │ • Set modal     │
│ • Stay on page  │   │   open state    │
└─────────────────┘   └─────────────────┘
                              ↓
                ┌─────────────────────────┐
                │ 6. Modal opens          │
                │    (isPasswordValid=true)│
                │    Display requirements │
                └─────────────────────────┘
                              ↓
                ┌─────────────────────────┐
                │ 7. User fills PII data  │
                │    Clicks confirm       │
                └─────────────────────────┘
                              ↓
                ┌─────────────────────────┐
                │ 8. confirm() called     │
                │    - Submit registration│
                │    - Success toast      │
                │    - Reset password     │
                │    - Close modal        │
                └─────────────────────────┘
```

## Key Benefits

### 🎯 Zero Configuration

- Just call `useSubmitPasswordUsecase(eventId)`
- No manual callback setup needed
- Works automatically with EventPasswordNav

### 🔄 Automatic State Management

- Password field resets after operations
- Modal state managed by the hook
- Loading states handled internally

### 🎨 Clean Component Code

- No password state in your component
- No manual form handling
- Focus on UI and business logic

### 🐛 Error Handling Built-in

- Destructive toasts for errors
- Network error handling
- Validation error handling

## Store Methods Used

```typescript
// From useEventPasswordNavStore:

setOnSubmitCallback: (callback: (password: string) => void) => void
// Sets the callback that runs when user submits from bottom nav
// Called automatically in useEffect

resetPassword: () => void
// Clears the password field
// Called after successful registration or reset
```

## When Password Field Gets Cleared

The password is automatically cleared in these scenarios:

1. ✅ **After successful registration** - `confirm()` success
2. ✅ **Manual reset** - `resetPasswordValidation()` called
3. ✅ **Modal close with reset** - User cancels and resets

## Best Practices

### ✅ DO

- Initialize the hook at component mount
- Use `isPasswordValid` to control modal visibility
- Call `resetPasswordValidation()` when modal closes
- Get password from store for `confirm()` call

### ❌ DON'T

- Don't manually set `onSubmitCallback` - it's automatic
- Don't manage password state separately
- Don't forget to reset when modal closes
- Don't call `preview()` manually (unless not using EventPasswordNav)

## Troubleshooting

### Password submission not working

- ✅ Check that hook is initialized: `useSubmitPasswordUsecase(eventId)`
- ✅ Verify EventPasswordNav is rendered in BottomContainer
- ✅ Check browser console for errors

### Modal not opening

- ✅ Verify `isPasswordValid` state is used for modal open prop
- ✅ Check that password check is successful (not returning error)
- ✅ Ensure AlertDialog component is imported correctly

### Password not clearing

- ✅ Call `resetPasswordValidation()` when modal closes
- ✅ Check that `resetPassword` is being called in hook
- ✅ Verify store state updates

## Related Files

- Hook: `/apps/web/src/hooks/events/useSubmitPasswordUsecase.ts`
- Store: `/apps/web/src/components/BottomNav/stores/event-password.ts`
- Bottom Nav: `/apps/web/src/components/BottomNav/variants/EventPasswordNav.tsx`
- Documentation: `/apps/web/src/hooks/events/useSubmitPasswordUsecase.README.md`
