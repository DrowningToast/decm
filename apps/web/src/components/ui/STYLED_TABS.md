# Styled Tabs Component

A reusable, consistently styled tabs component for the DECM platform with the signature Cormorant Garamond font and primary color scheme.

## Overview

The `styled-tabs` components provide a pre-styled wrapper around the base Radix UI tabs, ensuring visual consistency across all tabbed interfaces in the application.

## Components

### StyledTabs

Main container for the tabs component.

**Props:**

- `defaultValue: string` - The default active tab value (required)
- `children: React.ReactNode` - Tab content (required)
- `className?: string` - Additional CSS classes (optional)

### StyledTabsList

Container for tab triggers with consistent background styling.

**Props:**

- `children: React.ReactNode` - TabsTrigger components (required)

**Styling:**

- Width: Full width (`w-full`)
- Height: 40px (`h-10`)
- Background: Light pink (`bg-[#E9DEDE]`)

### StyledTabsTrigger

Individual tab trigger/button with Cormorant Garamond font.

**Props:**

- `value: string` - Unique identifier for the tab (required)
- `children: React.ReactNode` - Tab label content (required)

**Styling:**

- Font: Cormorant Garamond, 16px
- Default: Gray text (`text-gray-900`)
- Active: Primary background with white text (`bg-primary text-white`)
- Smooth state transitions

### StyledTabsContent

Content container for each tab panel.

**Props:**

- `value: string` - Matches corresponding TabsTrigger value (required)
- `children: React.ReactNode` - Tab content (required)
- `className?: string` - Additional CSS classes, defaults to `"mt-6"` (optional)

## Basic Usage

```tsx
import {
    StyledTabs,
    StyledTabsList,
    StyledTabsTrigger,
    StyledTabsContent,
} from "@/components/ui/styled-tabs";

function MyComponent() {
    return (
        <StyledTabs defaultValue="tab1">
            <StyledTabsList>
                <StyledTabsTrigger value="tab1">First Tab</StyledTabsTrigger>
                <StyledTabsTrigger value="tab2">Second Tab</StyledTabsTrigger>
                <StyledTabsTrigger value="tab3">Third Tab</StyledTabsTrigger>
            </StyledTabsList>

            <StyledTabsContent value="tab1">
                <div>Content for first tab</div>
            </StyledTabsContent>

            <StyledTabsContent value="tab2">
                <div>Content for second tab</div>
            </StyledTabsContent>

            <StyledTabsContent value="tab3">
                <div>Content for third tab</div>
            </StyledTabsContent>
        </StyledTabs>
    );
}
```

## Real-World Examples

### Example 1: Event Details Page

From `HostEventDetailsPage.tsx`:

```tsx
<StyledTabs defaultValue="event-info">
    <StyledTabsList>
        <StyledTabsTrigger value="event-info">Event Info</StyledTabsTrigger>
        <StyledTabsTrigger value="participants">Participants</StyledTabsTrigger>
        <StyledTabsTrigger value="certificates">Certificates</StyledTabsTrigger>
    </StyledTabsList>

    <StyledTabsContent value="event-info">
        <div className="flex flex-col gap-4 lg:flex-row">{/* Event information content */}</div>
    </StyledTabsContent>

    <StyledTabsContent value="participants">
        <div className="space-y-4">{/* Participant list and settings */}</div>
    </StyledTabsContent>

    <StyledTabsContent value="certificates">
        <div className="w-full bg-primary/10 border border-primary/20 rounded-lg p-6">
            {/* Certificate configuration */}
        </div>
    </StyledTabsContent>
</StyledTabs>
```

### Example 2: Settings Page with i18n

From `EventParticipantSettingPage.tsx`:

```tsx
import { useTranslation } from "react-i18next";

function SettingsPage() {
    const { t } = useTranslation();

    return (
        <StyledTabs defaultValue="settings">
            <StyledTabsList>
                <StyledTabsTrigger value="settings">
                    {t("participantSettings.tabs.settings")}
                </StyledTabsTrigger>
                <StyledTabsTrigger value="preview">
                    {t("participantSettings.tabs.preview")}
                </StyledTabsTrigger>
            </StyledTabsList>

            <StyledTabsContent value="settings">
                <ParticipantSettingsForm defaultValues={formValues} onSubmit={handleSubmit} />
            </StyledTabsContent>

            <StyledTabsContent value="preview">
                <RegistrationFormPreview settings={formValues} />
            </StyledTabsContent>
        </StyledTabs>
    );
}
```

## Customization

### Custom Content Margin

Override the default `mt-6` margin on content:

```tsx
<StyledTabsContent value="tab1" className="mt-0">
    <div>Content without top margin</div>
</StyledTabsContent>
```

### Additional Container Classes

Add custom classes to the main container:

```tsx
<StyledTabs defaultValue="tab1" className="max-w-4xl mx-auto">
    {/* tabs content */}
</StyledTabs>
```

## Styling Details

### Colors

- **Background**: `#E9DEDE` (light pink)
- **Active Tab**: `primary` (from theme)
- **Inactive Tab**: `text-gray-900`
- **Active Text**: White

### Typography

- **Font Family**: Cormorant Garamond (serif)
- **Font Size**: 16px
- **Font Weight**: Normal (inherited)

### Dimensions

- **TabsList Height**: 40px (`h-10`)
- **TabsList Width**: Full width (`w-full`)
- **Content Margin**: 24px top (`mt-6`) - customizable

## Accessibility

The styled tabs inherit all accessibility features from Radix UI:

- ✅ Keyboard navigation (Arrow keys, Home, End)
- ✅ ARIA attributes automatically applied
- ✅ Focus management
- ✅ Screen reader support

## Migration from Base Tabs

**Before:**

```tsx
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

<Tabs defaultValue="tab1">
    <TabsList className="w-full h-10 bg-[#E9DEDE]">
        <TabsTrigger
            value="tab1"
            className="data-[state=active]:bg-primary text-gray-900 data-[state=active]:text-white"
            style={{
                fontFamily: "Cormorant Garamond",
                fontSize: "16px",
            }}
        >
            Tab 1
        </TabsTrigger>
    </TabsList>
    <TabsContent value="tab1" className="mt-6">
        Content
    </TabsContent>
</Tabs>;
```

**After:**

```tsx
import {
    StyledTabs,
    StyledTabsList,
    StyledTabsTrigger,
    StyledTabsContent,
} from "@/components/ui/styled-tabs";

<StyledTabs defaultValue="tab1">
    <StyledTabsList>
        <StyledTabsTrigger value="tab1">Tab 1</StyledTabsTrigger>
    </StyledTabsList>
    <StyledTabsContent value="tab1">Content</StyledTabsContent>
</StyledTabs>;
```

## When to Use

### ✅ DO Use StyledTabs:

- Event management pages
- Settings interfaces
- Multi-section detail views
- Dashboard layouts
- Any tabbed interface requiring DECM branding

### ❌ DON'T Use StyledTabs:

- When tabs need completely different styling
- Inside components that already have their own design system
- For navigation menus (use router-based navigation instead)
- When base Radix tabs with custom styling is more appropriate

## Component Files

- **Component**: `apps/web/src/components/ui/styled-tabs.tsx`
- **Base Component**: `apps/web/src/components/ui/tabs.tsx` (Radix UI)
- **Documentation**: `apps/web/src/components/ui/STYLED_TABS.md` (this file)

## Related Components

- Base Tabs: `@/components/ui/tabs`
- Typography: `@/components/typography/typography`
- PageContainer: `@/components/container/PageContainer`
- SectionContainer: `@/components/container/SectionContainer`

## Examples in Codebase

1. **HostEventDetailsPage**: Three-tab interface (Event Info, Participants, Certificates)
    - File: `apps/web/src/components/pages/HostPages/EventsPage/HostEventDetailsPage.tsx`
    - Tabs: Event information, participant management, certificate settings

2. **EventParticipantSettingPage**: Two-tab interface (Settings, Preview)
    - File: `apps/web/src/components/pages/HostPages/EventsPage/EventParticipantSettingPage.tsx`
    - Tabs: Form configuration, live preview

## Performance Considerations

- Components are lightweight wrappers
- No additional rendering overhead
- Styling is applied via CSS classes (no inline style recalculation except font-family)
- Tab content is conditionally rendered based on active state

## Browser Support

Inherits support from Radix UI tabs:

- ✅ Chrome/Edge (latest)
- ✅ Firefox (latest)
- ✅ Safari (latest)
- ✅ Mobile browsers (iOS Safari, Chrome Mobile)

## Future Enhancements

Potential improvements for future versions:

- [ ] Support for icon + text tab triggers
- [ ] Variant prop for different color schemes
- [ ] Size variants (small, medium, large)
- [ ] Optional loading state for tab content
- [ ] Orientation support (vertical tabs)
- [ ] Custom animation options
