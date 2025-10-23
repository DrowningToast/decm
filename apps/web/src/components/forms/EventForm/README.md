# EventForm Component

## Overview

A fully functional event creation/editing form built with React Hook Form and Zod validation.

## Features

- ✅ **Type-safe validation** with Zod schema
- ✅ **i18n support** for English and Thai
- ✅ **Accessible form controls** with ARIA attributes
- ✅ **Loading states** during submission
- ✅ **Error messages** with proper translations
- ✅ **Create/Edit modes** for different use cases

## Form Fields

| Field         | Type   | Required | Validation                          |
| ------------- | ------ | -------- | ----------------------------------- |
| `name`        | string | Yes      | Min 3 characters                    |
| `description` | string | No       | Optional text area                  |
| `eventBanner` | File   | No       | Image file (JPEG/PNG/WebP), max 5MB |
| `eventIcon`   | File   | No       | Image file (JPEG/PNG/WebP), max 5MB |
| `startDate`   | Date   | Yes      | Valid date (Date Picker)            |
| `endDate`     | Date   | Yes      | Must be after start date            |
| `seatsCount`  | number | Yes      | Integer, minimum 1                  |

## Usage

### Basic Usage (Create Mode)

```tsx
import { EventForm } from "@/components/forms/EventForm";
import type { EventFormData } from "@/lib/schemas/eventFormSchema";

const MyPage = () => {
    const handleSubmit = async (data: EventFormData) => {
        console.log("Event data:", data);
        // Handle form submission
    };

    return <EventForm onSubmit={handleSubmit} mode="create" />;
};
```

### Edit Mode with Default Values

```tsx
const existingEvent = {
    name: "Tech Conference 2024",
    description: "Annual tech conference",
    eventBanner: new File([""], "banner.jpg", { type: "image/jpeg" }),
    eventIcon: new File([""], "icon.png", { type: "image/png" }),
    startDate: new Date("2024-03-20"),
    endDate: new Date("2024-03-20"),
    seatsCount: 100,
};

<EventForm defaultValues={existingEvent} onSubmit={handleUpdate} mode="edit" />;
```

### With Loading State

```tsx
const [isLoading, setIsLoading] = useState(false);

const handleSubmit = async (data: EventFormData) => {
    setIsLoading(true);
    try {
        await api.createEvent(data);
    } finally {
        setIsLoading(false);
    }
};

<EventForm onSubmit={handleSubmit} isLoading={isLoading} mode="create" />;
```

## Props

### `EventFormProps`

| Prop            | Type                                             | Default      | Description                          |
| --------------- | ------------------------------------------------ | ------------ | ------------------------------------ |
| `defaultValues` | `Partial<EventFormData>`                         | `undefined`  | Default form values for editing      |
| `onSubmit`      | `(data: EventFormData) => void \| Promise<void>` | **Required** | Callback when form is submitted      |
| `isLoading`     | `boolean`                                        | `false`      | Shows loading state on submit button |
| `mode`          | `'create' \| 'edit'`                             | `'create'`   | Form mode affects submit button text |

## Validation Schema

Located at: `apps/web/src/lib/schemas/eventFormSchema.ts`

### Rules:

- **name**: Required, minimum 3 characters
- **description**: Optional string
- **eventBanner**: Optional File object, must be JPEG/PNG/WebP format, max 5MB
- **eventIcon**: Optional File object, must be JPEG/PNG/WebP format, max 5MB
- **startDate**: Required, valid Date object (selected via Date Picker)
- **endDate**: Required, valid Date object (selected via Date Picker), must be after startDate
- **seatsCount**: Required, integer, minimum 1

### Error Messages

All error messages are translatable via i18n:

- `events.validation.nameRequired`
- `events.validation.nameMinLength`
- `events.validation.eventBannerRequired`
- `events.validation.eventBannerSize`
- `events.validation.eventBannerType`
- `events.validation.eventIconRequired`
- `events.validation.eventIconSize`
- `events.validation.eventIconType`
- `events.validation.startDateRequired`
- `events.validation.endDateRequired`
- `events.validation.endDateAfterStart`
- `events.validation.seatsCountRequired`
- `events.validation.seatsCountMin`
- `events.validation.seatsCountInteger`

## Example: CreateEventPage Implementation

```tsx
import { useTranslation } from 'react-i18next';
import { Typography } from '@/components/typography/typography';
import { EventForm } from '@/components/forms/EventForm';
import type { EventFormData } from '@/lib/schemas/eventFormSchema';
import { toast } from 'sonner';

export const CreateEventPage = () => {
  const { t } = useTranslation();

  const handleCreateEvent = async (data: EventFormData) => {
    try {
      // API call to create event
      await api.createEvent(data);

      toast.success(t('common.success'), {
        description: \`Event "\${data.name}" created successfully\`,
      });
    } catch (error) {
      toast.error(t('common.error'), {
        description: t('errors.generic'),
      });
    }
  };

  return (
    <div className="space-y-6">
      <Typography variant="header" tag="h1">
        {t('events.createEvent')}
      </Typography>
      <EventForm onSubmit={handleCreateEvent} mode="create" />
    </div>
  );
};
```

## Styling

The form uses:

- **Radix UI** components (Label, Popover)
- **Custom UI components** (Input, Textarea, Button, Calendar)
- **shadcn/ui Date Picker** pattern with React DayPicker
- **date-fns** for date formatting
- **Tailwind CSS** for styling
- **Typography component** for consistent text rendering

## Accessibility

- Proper ARIA labels and attributes
- Error messages with `role="alert"`
- Required field indicators (\*)
- Disabled state during loading
- Keyboard navigation support

## i18n Keys

### Form Labels

- `events.form.name`
- `events.form.description`
- `events.form.eventBanner`
- `events.form.eventBannerPlaceholder`
- `events.form.eventBannerChange`
- `events.form.eventBannerRemove`
- `events.form.eventIcon`
- `events.form.eventIconPlaceholder`
- `events.form.eventIconChange`
- `events.form.eventIconRemove`
- `events.form.startDate`
- `events.form.endDate`
- `events.form.seatsCount`

### Form Placeholders

- `events.form.namePlaceholder`
- `events.form.descriptionPlaceholder`
- `events.form.startDatePlaceholder`
- `events.form.endDatePlaceholder`
- `events.form.seatsCountPlaceholder`

### Submit Buttons

- `events.form.submitCreate`
- `events.form.submitUpdate`

## Type Definitions

```typescript
import { z } from "zod";

export type EventFormData = {
    name: string;
    description?: string;
    eventBanner?: File; // Image file (JPEG/PNG/WebP, max 5MB)
    eventIcon?: File; // Image file (JPEG/PNG/WebP, max 5MB)
    startDate: Date; // Date object from Date Picker
    endDate: Date; // Date object from Date Picker
    seatsCount: number;
};
```

## Related Files

- Schema: `apps/web/src/lib/schemas/eventFormSchema.ts`
- Component: `apps/web/src/components/forms/EventForm/EventForm.tsx`
- Translations (EN): `apps/web/src/lib/i18n/locales/en.json`
- Translations (TH): `apps/web/src/lib/i18n/locales/th.json`
- Usage Example: `apps/web/src/components/pages/HostPages/CreateEventPage/CreateEventPage.tsx`
- Route Page: `apps/web/src/pages/host/events/create/index.tsx`
