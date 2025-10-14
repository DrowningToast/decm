# RequirementItem Component

## 📋 Overview

A reusable UI component for displaying requirement status of form fields or data collection points. Supports three states: **Required**, **Optional**, and **Not Required**, each with distinct visual styling and icons.

---

## 🎯 Features

- ✅ **Three Status States**: Required, Optional, Not Required
- ✅ **Visual Indicators**: Icons and color coding for each state
- ✅ **Fully Internationalized**: Supports English, Thai, and more
- ✅ **TypeScript**: Type-safe with exported types
- ✅ **Dark Mode Support**: Adapts to light and dark themes
- ✅ **Reusable**: Can be used anywhere in the application
- ✅ **Accessible**: Semantic HTML with proper ARIA support

---

## 📖 Usage

### Basic Import

```tsx
import { RequirementItem, type RequirementStatus } from "@/components/ui/requirement-item";
```

### Simple Example

```tsx
<RequirementItem label="First Name" status="required" />
```

### All Three States

```tsx
{
    /* Required - Green with checkmark */
}
<RequirementItem label="First Name" status="required" />;

{
    /* Optional - Gray with minus icon */
}
<RequirementItem label="Bio" status="optional" />;

{
    /* Not Required - Red with X icon */
}
<RequirementItem label="Address" status="not_required" />;
```

### In a Grid Layout

```tsx
<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
    <RequirementItem label="First Name" status="required" />
    <RequirementItem label="Last Name" status="required" />
    <RequirementItem label="Email" status="required" />
    <RequirementItem label="Phone Number" status="optional" />
    <RequirementItem label="Bio" status="optional" />
    <RequirementItem label="Address" status="not_required" />
</div>
```

### With Custom ClassName

```tsx
<RequirementItem
    label="Custom Field"
    status="required"
    className="shadow-md hover:shadow-lg transition-shadow"
/>
```

### Dynamic Data-Driven

```tsx
const requirements: Record<string, RequirementStatus> = {
    firstName: "required",
    lastName: "required",
    email: "required",
    bio: "optional",
    phoneNumber: "optional",
    address: "not_required",
};

<div className="grid grid-cols-2 gap-4">
    {Object.entries(requirements).map(([field, status]) => (
        <RequirementItem key={field} label={t(`fields.${field}`)} status={status} />
    ))}
</div>;
```

---

## 🎨 Props

```typescript
interface RequirementItemProps {
    /**
     * Label text for the requirement field
     */
    label: string;

    /**
     * Status of the requirement: "required", "optional", or "not_required"
     */
    status: RequirementStatus;

    /**
     * Optional custom className for the container
     */
    className?: string;
}

export type RequirementStatus = "required" | "optional" | "not_required";
```

| Prop        | Type                | Required | Default | Description                                       |
| ----------- | ------------------- | -------- | ------- | ------------------------------------------------- |
| `label`     | `string`            | ✅       | -       | The text label for the requirement                |
| `status`    | `RequirementStatus` | ✅       | -       | Status: "required", "optional", or "not_required" |
| `className` | `string`            | ❌       | `""`    | Additional CSS classes for styling                |

---

## 🎨 Visual States

### 1. Required

- **Icon**: ✓ Green checkmark (`CheckCircle2Icon`)
- **Text**: "Required" / "จำเป็น"
- **Color**: Green (`text-green-600`)
- **Background**: Light green (`bg-green-50`)
- **Border**: Green (`border-green-200`)
- **Use Case**: Fields that must be filled

```tsx
<RequirementItem label="Email" status="required" />
```

**Appearance:**

```
┌────────────────────────────────────────┐
│ Email                    ✓ Required    │  <- Green background
└────────────────────────────────────────┘
```

### 2. Optional

- **Icon**: ⊖ Gray minus circle (`MinusCircleIcon`)
- **Text**: "Optional" / "ไม่บังคับ"
- **Color**: Muted gray (`text-muted-foreground`)
- **Background**: Light gray (`bg-muted/10`)
- **Border**: Default border
- **Use Case**: Fields that users can skip

```tsx
<RequirementItem label="Bio" status="optional" />
```

**Appearance:**

```
┌────────────────────────────────────────┐
│ Bio                      ⊖ Optional    │  <- Gray background
└────────────────────────────────────────┘
```

### 3. Not Required

- **Icon**: ✕ Red X (`XCircleIcon`)
- **Text**: "Not Required" / "ไม่ต้องการ"
- **Color**: Red (`text-red-600`)
- **Background**: Light red (`bg-red-50`)
- **Border**: Red (`border-red-200`)
- **Use Case**: Fields explicitly not needed

```tsx
<RequirementItem label="Address" status="not_required" />
```

**Appearance:**

```
┌────────────────────────────────────────┐
│ Address                  ✕ Not Required│  <- Red background
└────────────────────────────────────────┘
```

---

## 🌐 Internationalization

The component uses i18n for all status text:

### Translation Keys

```json
{
    "common": {
        "required": "Required",
        "optional": "Optional",
        "notRequired": "Not Required"
    }
}
```

### Supported Languages

| Status       | English      | Thai       |
| ------------ | ------------ | ---------- |
| Required     | Required     | จำเป็น     |
| Optional     | Optional     | ไม่บังคับ  |
| Not Required | Not Required | ไม่ต้องการ |

### Adding More Languages

Add translations to your locale files:

```json
// fr.json (French)
{
    "common": {
        "required": "Obligatoire",
        "optional": "Facultatif",
        "notRequired": "Non requis"
    }
}
```

---

## 🎯 Common Use Cases

### 1. Event Participant Requirements

```tsx
const participantRequirements = {
    firstName: "required",
    lastName: "required",
    email: "required",
    bio: "optional",
    phoneNumber: "optional",
    address: "not_required",
};

<Accordion>
    <AccordionItem value="requirements">
        <AccordionTrigger>Participant Requirements</AccordionTrigger>
        <AccordionContent>
            <div className="grid grid-cols-2 gap-4">
                {Object.entries(participantRequirements).map(([field, status]) => (
                    <RequirementItem
                        key={field}
                        label={t(`participants.fields.${field}`)}
                        status={status}
                    />
                ))}
            </div>
        </AccordionContent>
    </AccordionItem>
</Accordion>;
```

### 2. Form Field Configuration

```tsx
const formConfig = [
    { label: "Username", status: "required" },
    { label: "Display Name", status: "optional" },
    { label: "Middle Name", status: "not_required" },
];

<div className="space-y-2">
    {formConfig.map((field) => (
        <RequirementItem key={field.label} label={field.label} status={field.status} />
    ))}
</div>;
```

### 3. API Response Validation

```tsx
interface APIFieldRequirement {
    name: string;
    required: boolean;
    optional: boolean;
}

function mapToStatus(field: APIFieldRequirement): RequirementStatus {
    if (field.required) return "required";
    if (field.optional) return "optional";
    return "not_required";
}

<div className="grid grid-cols-3 gap-4">
    {apiFields.map((field) => (
        <RequirementItem key={field.name} label={field.name} status={mapToStatus(field)} />
    ))}
</div>;
```

---

## 🎨 Styling & Customization

### Default Styling

The component comes with pre-defined styles that adapt to dark mode:

```tsx
// Light mode
bg-green-50 border-green-200 text-green-600  // Required
bg-muted/10 border-[#D9D9D91A] text-muted-foreground  // Optional
bg-red-50 border-red-200 text-red-600  // Not Required

// Dark mode
bg-green-950/20 border-green-800 text-green-600  // Required
bg-muted/10 border-[#D9D9D91A] text-muted-foreground  // Optional
bg-red-950/20 border-red-800 text-red-600  // Not Required
```

### Custom Styling

Add custom classes via the `className` prop:

```tsx
<RequirementItem label="Custom Field" status="required" className="shadow-lg rounded-xl border-2" />
```

### Override Styles with Tailwind

```tsx
<RequirementItem
    label="Special Field"
    status="optional"
    className="bg-blue-50 border-blue-200 hover:bg-blue-100 transition-colors"
/>
```

---

## 💡 Best Practices

### ✅ DO:

```tsx
// Use clear, descriptive labels
<RequirementItem label="Email Address" status="required" />

// Use consistent status logic across your app
const getStatus = (field) => field.mandatory ? "required" : "optional";

// Group related requirements together
<div className="grid grid-cols-2 gap-4">
    <RequirementItem label="First Name" status="required" />
    <RequirementItem label="Last Name" status="required" />
</div>

// Use i18n for labels
<RequirementItem label={t("fields.email")} status="required" />
```

### ❌ DON'T:

```tsx
// Don't use vague labels
<RequirementItem label="Field 1" status="required" />

// Don't hardcode text
<RequirementItem label="Email" status="required" />  // Missing i18n

// Don't mix different UI patterns
// Either use RequirementItem OR custom badges, not both

// Don't use inconsistent status values
status="REQUIRED"  // ❌ Wrong (uppercase)
status={true}      // ❌ Wrong (boolean)
status="required"  // ✅ Correct
```

---

## 🧪 Testing

### Unit Test Example

```tsx
import { render, screen } from "@testing-library/react";
import { RequirementItem } from "@/components/ui/requirement-item";

describe("RequirementItem", () => {
    it("renders required status correctly", () => {
        render(<RequirementItem label="Email" status="required" />);
        expect(screen.getByText("Email")).toBeInTheDocument();
        expect(screen.getByText("Required")).toBeInTheDocument();
    });

    it("renders optional status correctly", () => {
        render(<RequirementItem label="Bio" status="optional" />);
        expect(screen.getByText("Optional")).toBeInTheDocument();
    });

    it("renders not_required status correctly", () => {
        render(<RequirementItem label="Address" status="not_required" />);
        expect(screen.getByText("Not Required")).toBeInTheDocument();
    });
});
```

---

## 🔧 TypeScript Support

### Type Definitions

```typescript
export type RequirementStatus = "required" | "optional" | "not_required";

interface RequirementItemProps {
    label: string;
    status: RequirementStatus;
    className?: string;
}
```

### Usage with Type Safety

```typescript
// Type-safe data structure
const requirements: Record<string, RequirementStatus> = {
    firstName: "required",
    lastName: "required",
    email: "required",
    bio: "optional",
    address: "not_required",
};

// Type-safe helper function
function getRequirementStatus(isRequired: boolean, isOptional: boolean): RequirementStatus {
    if (isRequired) return "required";
    if (isOptional) return "optional";
    return "not_required";
}
```

---

## 📊 Performance

- **Lightweight**: Minimal re-renders
- **No Heavy Dependencies**: Uses only icons from `lucide-react`
- **Optimized**: Conditional rendering for icons and styles
- **Memoizable**: Can be wrapped with `React.memo` if needed

---

## 🚀 Future Enhancements

Potential improvements:

1. **Tooltip Support**: Add hover tooltips for additional info
2. **Animation**: Smooth transitions between states
3. **Custom Icons**: Allow custom icons per status
4. **Badge Variant**: Add compact badge mode
5. **Interactive**: Click to toggle status (for admin panels)

---

## 📚 Related Components

- `Accordion` - For grouping requirements
- `Typography` - For consistent text styling
- `Card` - Alternative container for requirements

---

**Component Location**: `/components/ui/requirement-item.tsx`  
**Created**: October 14, 2025  
**Status**: ✅ Production Ready
