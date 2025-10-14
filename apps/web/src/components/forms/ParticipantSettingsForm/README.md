# Participant Settings Form

A comprehensive form component for configuring event participant registration requirements with a live preview.

## Overview

The Participant Settings Form allows event hosts to configure:

1. **Registration Settings**
    - Event Type (Public/Private)
    - Booking Approval Requirement
    - Ticket Transferability

2. **Participant Requirements**
    - Field-level configuration (Not Required/Optional/Required)
    - Basic Information: First Name, Last Name, Email, Phone Number
    - Additional Information: Bio, Address
    - Academic Information: Academic Institution, Academic Email

3. **Live Preview**
    - Real-time preview of the registration form
    - Shows exactly how participants will see the form
    - Updates dynamically based on settings

## Components

### ParticipantSettingsForm

Main form component for configuring participant settings.

**Props:**

- `defaultValues?: Partial<ParticipantSettingsData>` - Initial form values
- `onSubmit: (data: ParticipantSettingsData) => void | Promise<void>` - Submit handler
- `isLoading?: boolean` - Loading state

**Example:**

```tsx
import { ParticipantSettingsForm } from "@/components/forms/ParticipantSettingsForm";

function MyPage() {
    const handleSubmit = async (data) => {
        // Save settings to API
        await saveParticipantSettings(data);
    };

    return (
        <ParticipantSettingsForm
            defaultValues={{
                eventType: "public",
                isBookingRequired: false,
                isTicketTransferable: true,
                // All fields default to "not_required"
                // Override specific fields as needed
                firstName: "required",
                lastName: "required",
                email: "required",
            }}
            onSubmit={handleSubmit}
        />
    );
}
```

### RegistrationFormPreview

Preview component showing how the registration form will appear to participants.

**Props:**

- `settings: ParticipantSettingsData` - Configuration to preview

**Example:**

```tsx
import { RegistrationFormPreview } from "@/components/forms/ParticipantSettingsForm";

function MyPage() {
    const [settings, setSettings] = useState(defaultParticipantSettings);

    return <RegistrationFormPreview settings={settings} />;
}
```

## Schema

### ParticipantSettingsData

```typescript
{
    // Registration Settings
    eventType: "public" | "private",
    isBookingRequired: boolean,
    isTicketTransferable: boolean,

    // Participant Requirements
    firstName: "not_required" | "optional" | "required",
    lastName: "not_required" | "optional" | "required",
    email: "not_required" | "optional" | "required",
    bio: "not_required" | "optional" | "required",
    phoneNumber: "not_required" | "optional" | "required",
    address: "not_required" | "optional" | "required",
    academicInstitution: "not_required" | "optional" | "required",
    academicEmail: "not_required" | "optional" | "required"
}
```

## Default Values

```typescript
import { defaultParticipantSettings } from "@/lib/schemas/participantSettingsSchema";

// Default configuration
{
    eventType: "public",
    isBookingRequired: false,
    isTicketTransferable: true,
    firstName: "not_required",
    lastName: "not_required",
    email: "not_required",
    bio: "not_required",
    phoneNumber: "not_required",
    address: "not_required",
    academicInstitution: "not_required",
    academicEmail: "not_required"
}
```

## Usage in Page

The main page component demonstrates the complete implementation:

**Location:** `apps/web/src/components/pages/HostPages/EventsPage/EventParticipantSettingPage.tsx`

**Features:**

- Modal preview dialog for registration form
- Preview button next to save button
- Form state management with React Hook Form
- Real-time preview updates with watch()
- API integration placeholder
- i18n translations
- Responsive design

## Translations

All text is internationalized using react-i18next with the namespace `participantSettings`:

**English:** `apps/web/src/lib/i18n/locales/en.json`
**Thai:** `apps/web/src/lib/i18n/locales/th.json`

**Translation Keys:**

- `participantSettings.pageTitle`
- `participantSettings.registrationSettings`
- `participantSettings.eventType`
- `participantSettings.fields.firstName`
- `participantSettings.preview.title`
- ... and more

## Routing

The page is accessible at:

```
/host/events/:eventId/settings/participants
```

Example: `/host/events/123/settings/participants`

## Validation

Built with React Hook Form and Zod validation:

- All enum values are validated
- Type-safe form submission
- Real-time validation feedback
- Disabled submit button when form is invalid or pristine

## Styling

Uses Tailwind CSS and Radix UI components:

- Responsive grid layouts
- Consistent spacing with Tailwind classes
- Accessible form controls
- Dark mode support via CSS variables

## Related Files

- Schema: `apps/web/src/lib/schemas/participantSettingsSchema.ts`
- Form: `apps/web/src/components/forms/ParticipantSettingsForm/ParticipantSettingsForm.tsx`
- Preview: `apps/web/src/components/forms/ParticipantSettingsForm/RegistrationFormPreview.tsx`
- Page Component: `apps/web/src/components/pages/HostPages/EventsPage/EventParticipantSettingPage.tsx`
- Route: `apps/web/src/pages/host/events/[eventId]/settings/participant/index.tsx`
- Translations: `apps/web/src/lib/i18n/locales/{en,th}.json`

## TODO: API Integration

When integrating with the backend API:

1. Create backend endpoint to save participant settings
2. Add to OpenAPI spec with Swagger annotations
3. Generate TypeScript client: `pnpm gen-api:core`
4. Replace the TODO comment in `participants.tsx` with actual API call
5. Add proper error handling and success feedback (toast notifications)

**Example API Integration:**

```typescript
import { DefaultApi } from "@decm/api";

const api = new DefaultApi({ basePath: "http://localhost:8080/api/v1" });

const handleSubmit = async (data: ParticipantSettingsData) => {
    try {
        await api.updateEventParticipantSettings({
            eventId: eventId,
            participantSettings: data,
        });
        toast.success(t("participantSettings.saveSuccess"));
    } catch (error) {
        toast.error(t("participantSettings.saveError"));
        throw error;
    }
};
```

## Best Practices

1. **Form State:** Use React Hook Form's `watch` to sync form values with preview modal
2. **Validation:** Zod schema ensures type safety and runtime validation
3. **i18n:** All user-facing text should use translation keys
4. **Loading States:** Disable form controls during submission
5. **Accessibility:** Use semantic HTML and proper labels
6. **Responsive:** Mobile-first design with grid layouts
7. **Preview:** Preview modal shows real-time changes based on form values

## Testing

When implementing tests:

1. Test form submission with valid data
2. Test validation for each field requirement
3. Test preview updates based on settings changes
4. Test tab switching between Settings and Preview
5. Test i18n translations (both English and Thai)
6. Test loading states and disabled states
