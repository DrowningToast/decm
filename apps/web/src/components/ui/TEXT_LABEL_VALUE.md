# TextLabelValue Component

## 📋 Overview

A reusable UI component for displaying labeled values in a consistent format. Perfect for detail pages, data displays, and information summaries where you need to show label-value pairs.

---

## 🎯 Features

- ✅ **Clean Layout**: Label above value for consistent vertical spacing
- ✅ **Link Support**: Optional href to make values clickable
- ✅ **Icon Support**: Add icons at the end of values
- ✅ **Custom Styling**: Optional className for value customization
- ✅ **TypeScript**: Fully typed with exported interfaces
- ✅ **Accessibility**: Semantic HTML with proper structure
- ✅ **Responsive**: Works on all screen sizes

---

## 📖 Usage

### Basic Import

```tsx
import { TextLabelValue } from "@/components/ui/text-label-value";
```

### Simple Example

```tsx
<TextLabelValue label="Email" value="user@example.com" />
```

### With Link

```tsx
<TextLabelValue label="Website" value="example.com" href="https://example.com" />
```

### With End Icon

```tsx
import { ExternalLinkIcon } from "lucide-react";

<TextLabelValue
    label="Documentation"
    value="View Docs"
    href="https://docs.example.com"
    endIcon={<ExternalLinkIcon className="h-4 w-4" />}
/>;
```

### With Custom Styling

```tsx
<TextLabelValue label="Status" value="Active" valueClassName="text-green-600 font-semibold" />
```

### In a Grid Layout

```tsx
<div className="grid grid-cols-1 md:grid-cols-4 gap-4">
    <TextLabelValue label="Name" value="John Doe" />
    <TextLabelValue label="Email" value="john@example.com" />
    <TextLabelValue label="Phone" value="+1 234 567 890" />
    <TextLabelValue label="Status" value="Active" />
</div>
```

---

## 🎨 Props

```typescript
interface TextLabelValueProps {
    /**
     * Label text displayed above the value
     */
    label: string;

    /**
     * Value text to display
     */
    value: string;

    /**
     * Optional icon displayed at the end of the value
     */
    endIcon?: React.ReactNode;

    /**
     * Optional custom className for the value text
     */
    valueClassName?: string;

    /**
     * Optional href to make the value a clickable link
     */
    href?: string;
}
```

| Prop             | Type              | Required | Default | Description                          |
| ---------------- | ----------------- | -------- | ------- | ------------------------------------ |
| `label`          | `string`          | ✅       | -       | Label text shown above the value     |
| `value`          | `string`          | ✅       | -       | Value text to display                |
| `endIcon`        | `React.ReactNode` | ❌       | -       | Icon shown at the end of the value   |
| `valueClassName` | `string`          | ❌       | -       | Custom CSS classes for value styling |
| `href`           | `string`          | ❌       | -       | URL to make the value clickable      |

---

## 🎨 Visual Examples

### Default (No Link)

```
┌──────────────────┐
│ Email            │  <- Label (muted, small)
│ user@example.com │  <- Value (default text)
└──────────────────┘
```

### With Link

```
┌──────────────────┐
│ Website          │  <- Label (muted, small)
│ example.com →    │  <- Value (clickable, with icon)
└──────────────────┘
```

### With Custom Color

```
┌──────────────────┐
│ Status           │  <- Label (muted, small)
│ Active           │  <- Value (green, bold)
└──────────────────┘
```

---

## 🔧 Common Use Cases

### 1. User Profile Display

```tsx
<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
    <TextLabelValue label="Full Name" value="John Doe" />
    <TextLabelValue label="Email" value="john@example.com" />
    <TextLabelValue label="Phone" value="+1 234 567 890" />
    <TextLabelValue label="Role" value="Administrator" />
</div>
```

### 2. Event Details Page

```tsx
<div className="grid grid-cols-1 md:grid-cols-4 gap-4">
    <TextLabelValue label="Event Name" value="Tech Conference 2025" />
    <TextLabelValue label="Date" value="2025-06-15" />
    <TextLabelValue label="Location" value="Bangkok, Thailand" />
    <TextLabelValue label="Seats" value="150/200" />
</div>
```

### 3. External Links with Icons

```tsx
import { ExternalLinkIcon } from "lucide-react";

<TextLabelValue
    label="View on Blockchain"
    value="View Transaction"
    href="https://etherscan.io/tx/0x123..."
    endIcon={<ExternalLinkIcon className="h-4 w-4" />}
    valueClassName="text-blue-600 hover:underline"
/>;
```

### 4. Status Indicators

```tsx
import { CheckCircle2Icon } from "lucide-react";

<TextLabelValue
    label="Verification Status"
    value="Verified"
    endIcon={<CheckCircle2Icon className="h-4 w-4 text-green-600" />}
    valueClassName="text-green-600 font-semibold"
/>;
```

### 5. Contact Information

```tsx
<div className="space-y-4">
    <TextLabelValue label="Email" value="contact@example.com" href="mailto:contact@example.com" />
    <TextLabelValue label="Phone" value="+1 234 567 890" href="tel:+1234567890" />
    <TextLabelValue
        label="Website"
        value="www.example.com"
        href="https://www.example.com"
        endIcon={<ExternalLinkIcon className="h-4 w-4" />}
    />
</div>
```

---

## 🎨 Styling

### Default Typography Styles

- **Label**: Small, muted text
- **Value**: Base size, default text color
- **Link**: Opens in new tab with `target="_blank"` and `rel="noopener noreferrer"`

### Custom Styling Examples

#### Success Status

```tsx
<TextLabelValue label="Status" value="Active" valueClassName="text-green-600 font-semibold" />
```

#### Warning Status

```tsx
<TextLabelValue label="Status" value="Pending" valueClassName="text-yellow-600 font-medium" />
```

#### Error Status

```tsx
<TextLabelValue label="Status" value="Failed" valueClassName="text-red-600 font-bold" />
```

#### With Background

```tsx
<TextLabelValue
    label="Badge"
    value="Premium"
    valueClassName="inline-block px-3 py-1 bg-purple-100 text-purple-700 rounded-full font-medium"
/>
```

---

## 💡 Best Practices

### ✅ DO:

```tsx
// Use descriptive labels
<TextLabelValue label="Email Address" value="user@example.com" />

// Add rel="noopener noreferrer" for external links (done automatically)
<TextLabelValue label="Website" value="example.com" href="https://example.com" />

// Use consistent spacing in grids
<div className="grid grid-cols-2 gap-4">
    <TextLabelValue label="First Name" value="John" />
    <TextLabelValue label="Last Name" value="Doe" />
</div>

// Use icons to indicate link behavior
<TextLabelValue
    label="External Link"
    value="View"
    href="https://example.com"
    endIcon={<ExternalLinkIcon className="h-4 w-4" />}
/>
```

### ❌ DON'T:

```tsx
// Don't use vague labels
<TextLabelValue label="Info" value="..." />  // ❌ Too vague

// Don't mix different icon sizes
<TextLabelValue
    value="View"
    endIcon={<Icon className="h-8 w-8" />}  // ❌ Too large
/>

// Don't forget to handle empty values
<TextLabelValue label="Phone" value={phone || "N/A"} />  // ✅ Handle empty

// Don't use for complex content (use custom components instead)
<TextLabelValue label="Items" value="Item 1, Item 2, Item 3..." />  // ❌ Use list
```

---

## 🧪 Testing

### Unit Test Example

```tsx
import { render, screen } from "@testing-library/react";
import { TextLabelValue } from "@/components/ui/text-label-value";

describe("TextLabelValue", () => {
    it("renders label and value correctly", () => {
        render(<TextLabelValue label="Email" value="test@example.com" />);
        expect(screen.getByText("Email")).toBeInTheDocument();
        expect(screen.getByText("test@example.com")).toBeInTheDocument();
    });

    it("renders link when href is provided", () => {
        render(<TextLabelValue label="Website" value="example.com" href="https://example.com" />);
        const link = screen.getByRole("link");
        expect(link).toHaveAttribute("href", "https://example.com");
        expect(link).toHaveAttribute("target", "_blank");
    });

    it("renders end icon when provided", () => {
        render(
            <TextLabelValue
                label="Link"
                value="View"
                endIcon={<span data-testid="icon">→</span>}
            />,
        );
        expect(screen.getByTestId("icon")).toBeInTheDocument();
    });
});
```

---

## 🚀 Performance

- **Lightweight**: Minimal component, fast render
- **No State**: Pure presentational component
- **Optimized**: Can be memoized with `React.memo` if needed
- **Accessible**: Semantic HTML structure

---

## 🌐 Accessibility

### Built-in Features

- ✅ Semantic HTML structure
- ✅ Proper link attributes (`target`, `rel`)
- ✅ Clear label-value relationship
- ✅ Screen reader friendly

### Improvements You Can Add

```tsx
// Add ARIA labels for better screen reader support
<TextLabelValue
    label="Status"
    value="Active"
    valueClassName="text-green-600"
    aria-label="Account status is active"
/>
```

---

## 🎯 Real-World Example

### Event Details Page

```tsx
export function EventDetailsPage({ event }) {
    return (
        <div className="space-y-6">
            {/* Header Section */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <TextLabelValue label="Status" value={event.status} />
                <TextLabelValue label="Final Call" value={event.finalCallDate} />
                <TextLabelValue label="Request Type" value={event.requestType} />
                <TextLabelValue label="Seats" value={`${event.filledSeats}/${event.totalSeats}`} />
            </div>

            {/* Contact Section */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <TextLabelValue label="Organizer" value={event.organizerName} />
                <TextLabelValue
                    label="Email"
                    value={event.organizerEmail}
                    href={`mailto:${event.organizerEmail}`}
                />
                <TextLabelValue
                    label="Phone"
                    value={event.organizerPhone}
                    href={`tel:${event.organizerPhone}`}
                />
                <TextLabelValue
                    label="Website"
                    value={event.website}
                    href={event.websiteUrl}
                    endIcon={<ExternalLinkIcon className="h-4 w-4" />}
                />
            </div>
        </div>
    );
}
```

---

## 📚 Related Components

- `Typography` - Used internally for consistent text styling
- `RequirementItem` - Similar component for showing requirement status
- `Card` - Container component often used with TextLabelValue

---

## 🔧 TypeScript Support

### Type Definitions

```typescript
interface TextLabelValueProps {
    label: string;
    value: string;
    endIcon?: React.ReactNode;
    valueClassName?: string;
    href?: string;
}
```

### Usage with Type Safety

```typescript
// Type-safe props
const props: TextLabelValueProps = {
    label: "Email",
    value: "user@example.com",
    href: "mailto:user@example.com",
};

<TextLabelValue {...props} />
```

---

**Component Location**: `/components/ui/text-label-value.tsx`  
**Created**: October 14, 2025  
**Status**: ✅ Production Ready
