# Wrapped Form Input Components

Reusable, all-in-one form input components that integrate seamlessly with React Hook Form. These components include labels, error messages, and proper validation handling out of the box.

## 📍 Location

`apps/web/src/components/forms/wrapped-inputs/`

## 🎯 Features

- ✅ **All-in-one**: Includes Label, Input, and Error Message
- ✅ **React Hook Form Integration**: Works with `control` prop
- ✅ **TypeScript Support**: Fully typed with generics
- ✅ **i18n Ready**: All text supports translation
- ✅ **Accessible**: Proper ARIA attributes and semantic HTML
- ✅ **Consistent Styling**: Uses project's design system
- ✅ **Required Field Indicator**: Automatic `*` for required fields

## 📦 Available Components

### 1. `<WrappedInput />`

Text, email, number, password, tel, and URL inputs.

### 2. `<WrappedTextarea />`

Multi-line text input for longer content.

### 3. `<WrappedDateSelect />`

Date picker with calendar popover.

### 4. `<WrappedInputFile />`

File upload with preview (image preview or file info display).

## 🚀 Basic Usage

### Import

```tsx
import {
    WrappedInput,
    WrappedTextarea,
    WrappedDateSelect,
    WrappedInputFile,
} from "@/components/forms/wrapped-inputs";
```

### Setup Form

```tsx
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { myFormSchema, type MyFormData } from "@/lib/schemas/myFormSchema";

const MyForm = () => {
    const { handleSubmit, control } = useForm<MyFormData>({
        resolver: zodResolver(myFormSchema),
        defaultValues: {
            name: "",
            description: "",
            startDate: undefined,
            file: undefined,
        },
    });

    const onSubmit = (data: MyFormData) => {
        console.log(data);
    };

    return <form onSubmit={handleSubmit(onSubmit)}>{/* Your wrapped components here */}</form>;
};
```

## 📖 Component Documentation

### WrappedInput

#### Props

```tsx
interface WrappedInputProps<T extends FieldValues> {
    name: Path<T>; // Field name (from form schema)
    control: Control<T>; // React Hook Form control
    label: string; // Label text (use t() for i18n)
    placeholder?: string; // Placeholder text
    type?: "text" | "email" | "number" | "password" | "tel" | "url";
    required?: boolean; // Shows * indicator
    disabled?: boolean; // Disables input
    className?: string; // Additional CSS classes
    min?: number; // Min value (for number type)
    max?: number; // Max value (for number type)
    step?: number; // Step value (for number type)
}
```

#### Examples

```tsx
// Text input
<WrappedInput
    name="name"
    control={control}
    label={t("events.form.name")}
    placeholder={t("events.form.namePlaceholder")}
    required
    disabled={isLoading}
/>

// Email input
<WrappedInput
    name="email"
    control={control}
    label={t("profile.email")}
    placeholder="user@example.com"
    type="email"
    required
/>

// Number input
<WrappedInput
    name="seatsCount"
    control={control}
    label={t("events.form.seatsCount")}
    placeholder={t("events.form.seatsCountPlaceholder")}
    type="number"
    required
    min={1}
    step={1}
/>
```

### WrappedTextarea

#### Props

```tsx
interface WrappedTextareaProps<T extends FieldValues> {
    name: Path<T>;
    control: Control<T>;
    label: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    className?: string;
    rows?: number; // Number of visible text lines
    maxLength?: number; // Maximum character length
}
```

#### Example

```tsx
<WrappedTextarea
    name="description"
    control={control}
    label={t("events.form.description")}
    placeholder={t("events.form.descriptionPlaceholder")}
    disabled={isLoading}
    rows={4}
    maxLength={500}
/>
```

### WrappedDateSelect

#### Props

```tsx
interface WrappedDateSelectProps<T extends FieldValues> {
    name: Path<T>;
    control: Control<T>;
    label: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    className?: string;
    minDate?: Date; // Minimum selectable date
    maxDate?: Date; // Maximum selectable date
    disablePastDates?: boolean; // Disable dates before today
}
```

#### Examples

```tsx
// Basic date picker
<WrappedDateSelect
    name="startDate"
    control={control}
    label={t("events.form.startDate")}
    placeholder={t("events.form.startDatePlaceholder")}
    required
    disablePastDates
/>

// Date range (end date must be after start date)
<div className="grid grid-cols-2 gap-4">
    <WrappedDateSelect
        name="startDate"
        control={control}
        label={t("events.form.startDate")}
        required
        disablePastDates
    />
    <WrappedDateSelect
        name="endDate"
        control={control}
        label={t("events.form.endDate")}
        required
        disablePastDates
        minDate={watch("startDate")} // End date after start date
    />
</div>
```

### WrappedInputFile

#### Props

```tsx
interface WrappedInputFileProps<T extends FieldValues> {
    name: Path<T>;
    control: Control<T>;
    label: string;
    accept?: string; // File types (default: images)
    required?: boolean;
    disabled?: boolean;
    className?: string;
    maxSize?: number; // Max file size in bytes (default: 5MB)
    previewClassName?: string; // Additional classes for preview
}
```

#### Features

- **Image Files**: Shows full image preview
- **Non-Image Files**: Shows file info (name, size, icon)
- **Single File**: Only one file at a time
- **Actions**: Change file or remove file buttons

#### Examples

```tsx
// Image upload with preview
<WrappedInputFile
    name="eventBanner"
    control={control}
    label={t("events.form.eventBanner")}
    accept="image/jpeg,image/jpg,image/png,image/webp"
    required
    maxSize={5 * 1024 * 1024} // 5MB
    previewClassName="object-cover"
/>

// Any file type
<WrappedInputFile
    name="document"
    control={control}
    label={t("documents.upload")}
    accept=".pdf,.doc,.docx"
    required
    maxSize={10 * 1024 * 1024} // 10MB
/>
```

## 🎨 Styling

All components follow the project's design system:

- Consistent border styling: `border border-[#D9D9D91A]`
- Proper spacing with Tailwind classes
- Typography component integration
- Required field indicator (`*`)
- Error message styling with `text-destructive`

## 🌐 i18n Integration

All wrapped components work seamlessly with react-i18next:

```tsx
import { useTranslation } from "react-i18next";

const MyForm = () => {
    const { t } = useTranslation();
    const { control, handleSubmit } = useForm();

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <WrappedInput
                name="name"
                control={control}
                label={t("form.name")} // Translated label
                placeholder={t("form.namePlaceholder")} // Translated placeholder
                required
            />
        </form>
    );
};
```

### Required Translation Keys

Add these to your `locales/en.json` and `locales/th.json`:

```json
{
    "common": {
        "selectDate": "Select date",
        "clickToUpload": "Click to upload",
        "change": "Change",
        "remove": "Remove"
    }
}
```

## ✅ Validation

Error messages are automatically displayed when validation fails:

```tsx
// Schema with validation
import { z } from "zod";

const eventFormSchema = z.object({
    name: z
        .string()
        .min(3, "events.validation.nameMinLength")
        .max(100, "events.validation.nameMaxLength"),
    email: z.string().email("validation.invalidEmail"),
    seatsCount: z
        .number()
        .min(1, "events.validation.seatsCountMin")
        .int("events.validation.seatsCountInteger"),
});
```

Error messages are translated and displayed below each field automatically.

## 📊 Real-World Example

See `apps/web/src/components/forms/EventForm/EventForm.tsx` for a complete implementation:

**Before (489 lines)** → **After (160 lines)** 🎉

```tsx
export const EventForm = ({ onSubmit, isLoading }: EventFormProps) => {
    const { t } = useTranslation();
    const { handleSubmit, control } = useForm<EventFormData>({
        resolver: zodResolver(eventFormSchema),
    });

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Event Banner & Icon */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="lg:col-span-3">
                    <WrappedInputFile
                        name="eventBanner"
                        control={control}
                        label={t("events.form.eventBanner")}
                        accept="image/*"
                        required
                    />
                </div>
                <div>
                    <WrappedInputFile
                        name="eventIcon"
                        control={control}
                        label={t("events.form.eventIcon")}
                        accept="image/*"
                        required
                    />
                </div>
            </div>

            {/* Event Name */}
            <WrappedInput
                name="name"
                control={control}
                label={t("events.form.name")}
                placeholder={t("events.form.namePlaceholder")}
                required
                disabled={isLoading}
            />

            {/* Description */}
            <WrappedTextarea
                name="description"
                control={control}
                label={t("events.form.description")}
                placeholder={t("events.form.descriptionPlaceholder")}
            />

            {/* Dates */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <WrappedDateSelect
                    name="startDate"
                    control={control}
                    label={t("events.form.startDate")}
                    required
                    disablePastDates
                />
                <WrappedDateSelect
                    name="endDate"
                    control={control}
                    label={t("events.form.endDate")}
                    required
                    disablePastDates
                />
            </div>

            {/* Seats Count */}
            <WrappedInput
                name="seatsCount"
                control={control}
                label={t("events.form.seatsCount")}
                type="number"
                required
                min={1}
                step={1}
            />

            {/* Submit */}
            <Button type="submit" disabled={isLoading}>
                {t("common.submit")}
            </Button>
        </form>
    );
};
```

## 🔧 Customization

### Adding Custom Styling

```tsx
<WrappedInput
    name="name"
    control={control}
    label="Name"
    className="custom-input-class" // Applied to input element
/>
```

### Extending Components

Create your own wrapped components following the same pattern:

```tsx
// WrappedSelect.tsx
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";

interface WrappedSelectProps<T extends FieldValues> {
    name: Path<T>;
    control: Control<T>;
    label: string;
    options: Array<{ value: string; label: string }>;
    required?: boolean;
}

export const WrappedSelect = <T extends FieldValues>({
    name,
    control,
    label,
    options,
    required,
}: WrappedSelectProps<T>) => {
    const { t } = useTranslation();

    return (
        <Controller
            name={name}
            control={control}
            render={({ field, fieldState: { error } }) => (
                <div className="space-y-2">
                    <Label>
                        <Typography variant="text" tag="span">
                            {label}
                        </Typography>
                        {required && <span className="text-destructive ml-1">*</span>}
                    </Label>
                    <Select value={field.value} onValueChange={field.onChange}>
                        {/* Select implementation */}
                    </Select>
                    {error && (
                        <Typography variant="text" tag="p" className="text-sm text-destructive">
                            {t(error.message as string)}
                        </Typography>
                    )}
                </div>
            )}
        />
    );
};
```

## 🚨 Common Issues

### Issue: TypeScript errors with generics

**Solution**: Ensure your form data type extends `FieldValues` from React Hook Form.

### Issue: Error messages not translating

**Solution**: Make sure your error messages in Zod schema are translation keys (strings), not direct error messages.

### Issue: File preview not updating

**Solution**: The `WrappedInputFile` component handles this automatically. Ensure you're not manually managing file state.

## 📚 Related Documentation

- [Typography Usage](../../../README.md#typography)
- [i18n Translations](../../../../lib/i18n/README.md)
- [React Hook Form Docs](https://react-hook-form.com/)
- [Zod Validation](https://zod.dev/)

## 🎓 Best Practices

1. **Always use `control` prop** from `useForm()` hook
2. **Use translation keys** for labels and placeholders
3. **Define validation** in Zod schema, not in components
4. **Keep forms simple** - one wrapped component per field
5. **Test with different languages** to ensure proper text wrapping
6. **Use semantic HTML** - components already handle this
7. **Leverage TypeScript** - full type safety out of the box

## 💡 Tips

- Use `required` prop instead of adding `*` manually
- Combine with `formState.isSubmitting` for disabled state
- Use grid layouts for responsive form fields
- Group related fields with semantic `<div>` wrappers
- Leverage `min`, `max`, `step` for number inputs
- Use `disablePastDates` for event/booking dates

---

Created for DECM Platform | React 19 + React Hook Form + Zod + TypeScript
