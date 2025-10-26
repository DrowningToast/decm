# Certificate Settings Page

## Overview

The Certificate Settings Page allows event hosts to configure certificate issuers and upload SVG certificate templates with dynamic keywords that will be replaced with actual event and participant data.

## Location

- **Component**: `apps/web/src/components/pages/HostPages/EventPages/CertificateSettingsPage.tsx`
- **Route**: `/host/events/[eventId]/settings/certificate`

## Features

### Step 1: Issuer Settings

1. **Issuer Search**
    - Search issuers by name, email, or organization
    - Opens a modal dialog with search results (4 mock results)
    - Select multiple issuers with checkboxes
    - Shows count of selected issuers in modal
    - "Choose" button to confirm selection and add to selected list

2. **Selected Issuers Table**
    - Display issuer name, email, and organization
    - Remove issuers with delete button
    - Shows count of selected issuers

3. **Important Notice Alert**
    - Warning that issuer settings can only be configured once
    - Can be changed after confirmation, but should be done carefully

### Step 2: Certificate Template Settings

1. **Figma Instructions**
    - Step-by-step guide on creating certificate templates in Figma
    - Three instruction steps with placeholder images
    - Visual guidance for designers

2. **Available Keywords**
    - Display of all supported dynamic keywords:
        - `{{ eventName }}`
        - `{{ name }}`
        - `{{ academicInstitutionName }}`
        - `{{ startDate }}`
        - `{{ endDate }}`

3. **SVG Upload**
    - File input for SVG certificate templates
    - Only accepts `.svg` or `image/svg+xml` files
    - Shows selected file name

4. **Template Preview**
    - Live preview of uploaded SVG template
    - Rendered directly in the browser
    - Maximum height: 400px with scroll

5. **Detected Keywords Table**
    - Automatically parses SVG content for keywords
    - Shows keyword name, X position, Y position
    - Count column shows how many times keyword appears
    - Helps verify template correctness

## Components Used

### UI Components

- `Table`, `TableHeader`, `TableBody`, `TableRow`, `TableHead`, `TableCell` - Display data in tables
- `Alert`, `AlertTitle`, `AlertDescription` - Show warnings and information
- `Input` - Text search and file upload
- `Button` - Actions and navigation
- `Label` - Form labels
- `Typography` - Consistent text rendering

### Layout Components

- `PageContainer` - Page wrapper with title and description
- `SectionContainer` - Content section wrapper

### Icons (lucide-react)

- `Search` - Search button icon
- `Trash2` - Delete issuer icon
- `Upload` - Upload file icon
- `AlertCircle` - Warning alert icon
- `Info` - Information alert icon
- `ImageIcon` - Placeholder for instruction images

## Data Structures

### Issuer

```typescript
interface Issuer {
    id: string;
    name: string;
    email: string;
    organization?: string;
}
```

### DetectedKeyword

```typescript
interface DetectedKeyword {
    keyword: string;
    x: number;
    y: number;
    count: number;
}
```

## State Management

- `searchQuery` - Current search input value
- `selectedIssuers` - Array of selected issuer objects
- `isSearching` - Loading state for issuer search
- `searchResults` - Array of search result issuers (shown in modal)
- `isSearchModalOpen` - Modal dialog open/close state
- `tempSelectedIds` - Set of temporarily selected issuer IDs in modal
- `svgFile` - Uploaded SVG file object
- `svgPreview` - SVG content as string for preview
- `detectedKeywords` - Parsed keywords from SVG
- `isLoading` - Loading state for form submission

## Key Functions

### `handleSearchIssuers()`

- Searches for issuers based on search query
- Opens modal dialog with 4 mock search results
- Pre-selects already selected issuers in modal
- TODO: Implement actual API call

### `handleToggleIssuerSelection(issuerId)`

- Toggles issuer selection in modal
- Updates temporary selection state (Set)

### `handleConfirmIssuerSelection()`

- Merges temporary selections with current selections
- Avoids duplicates by checking IDs
- Closes modal and resets temporary state

### `handleCancelIssuerSelection()`

- Closes modal without saving changes
- Resets temporary selection state

### `handleRemoveIssuer(issuerId)`

- Removes issuer from selected list by ID

### `handleFileSelect(event)`

- Handles SVG file selection
- Validates file type (must be SVG)
- Reads file content for preview
- Triggers keyword parsing

### `parseKeywordsFromSVG(svgContent)`

- Parses SVG XML to find text elements
- Extracts keywords matching pattern: `{{ keywordName }}`
- Gets X and Y coordinates from text elements
- Counts keyword occurrences
- Updates detected keywords state

### `handleSubmit()`

- Validates form data
- Submits certificate settings
- TODO: Implement actual API call

### `handleCancel()`

- Resets form to initial state
- Clears all selected data

## Form Validation

The submit button is disabled when:

- No issuers are selected
- No SVG file is uploaded
- No keywords are detected in the template
- Form is currently loading/submitting

## Internationalization (i18n)

### Translation Keys

All text is internationalized using react-i18next:

```typescript
// Page
certificateSettings.pageTitle
certificateSettings.pageDescription

// Step 1
certificateSettings.step1.title
certificateSettings.step1.description
certificateSettings.step1.searchLabel
certificateSettings.step1.searchPlaceholder
certificateSettings.step1.searchButton
certificateSettings.step1.selectedIssuers
certificateSettings.step1.table.*
certificateSettings.step1.alert.*
certificateSettings.step1.modal.title
certificateSettings.step1.modal.description
certificateSettings.step1.modal.selected
certificateSettings.step1.modal.noResults
certificateSettings.step1.modal.chooseButton

// Step 2
certificateSettings.step2.title
certificateSettings.step2.description
certificateSettings.step2.instructions.*
certificateSettings.step2.keywords.*
certificateSettings.step2.upload.*
certificateSettings.step2.preview.*
certificateSettings.step2.detectedKeywords.*

// Actions
certificateSettings.confirmButton
certificateSettings.saveSuccess
certificateSettings.saveError

// Common
common.cancel
common.loading
common.searching
```

### Supported Languages

- **English (en)**: Full translations
- **Thai (th)**: Full translations (ไทย)

## Usage Example

### In Route File

```tsx
import { CertificateSettingsPage } from "@/components/pages/HostPages/EventPages/CertificateSettingsPage";

export default function Page() {
    return <CertificateSettingsPage />;
}
```

### With useParams

The component automatically gets `eventId` from URL params:

```
/host/events/123/settings/certificate
```

## TODO / Future Improvements

1. **API Integration**
    - Implement actual issuer search API call
    - Implement certificate settings save API call
    - Handle loading and error states properly

2. **Issuer Management**
    - Add issuer profile view/modal
    - Implement issuer invitation system
    - Add issuer role/permission management

3. **Template Validation**
    - Validate that required keywords are present
    - Warn if unknown keywords are detected
    - Check SVG file size (implement 5MB limit)
    - Validate SVG structure and format

4. **Enhanced Preview**
    - Show preview with sample data filled in
    - Allow zooming and panning of preview
    - Side-by-side comparison of template and preview

5. **Template Library**
    - Provide pre-made certificate templates
    - Template gallery with categories
    - Allow saving custom templates

6. **Keyword Editor**
    - Visual editor to position keywords
    - Drag-and-drop keyword placement
    - Live preview while editing

7. **Accessibility**
    - Add proper ARIA labels
    - Keyboard navigation support
    - Screen reader announcements

8. **Testing**
    - Unit tests for keyword parsing
    - Integration tests for form submission
    - E2E tests for complete workflow

## Related Files

- Translation files:
    - `apps/web/src/lib/i18n/locales/en.json`
    - `apps/web/src/lib/i18n/locales/th.json`
- UI components:
    - `apps/web/src/components/ui/table.tsx` (new)
    - `apps/web/src/components/ui/alert.tsx` (new)
- Route file:
    - `apps/web/src/pages/host/events/[eventId]/settings/certificate/index.tsx`

## Design Decisions

### Why SVG?

- Vector graphics scale perfectly for different sizes
- Easy to parse and manipulate text elements
- Can extract coordinate information for keyword placement
- Lightweight and web-friendly format

### Why Keyword Pattern `{{ keywordName }}`?

- Common templating pattern (similar to Handlebars, Jinja)
- Easy to type and remember
- Visually distinct in design tools
- Simple to parse with regex

### Why DOMParser for SVG Parsing?

- Native browser API, no extra dependencies
- Reliable XML parsing
- Easy to query elements with querySelector
- Extract attributes (x, y coordinates) easily

### Why One-Time Issuer Setting?

- Maintains certificate integrity and authenticity
- Prevents unauthorized changes to issued certificates
- Audit trail for accountability
- Can still be changed but with careful consideration

## Development Notes

### SVG File Requirements

For the SVG template to work correctly:

1. Text elements should contain keywords in the format `{{ keywordName }}`
2. Text elements must have `x` and `y` attributes for position detection
3. Keywords should match the supported keywords list
4. File should be valid SVG XML

### Keyword Detection Algorithm

1. Parse SVG content with DOMParser
2. Query all `text` and `tspan` elements
3. Search text content with regex: `/\{\{\s*(\w+)\s*\}\}/g`
4. Extract `x` and `y` attributes from elements
5. Count occurrences of each keyword
6. Store results in state

## Troubleshooting

### Keywords Not Detected

- Check that SVG file is valid XML
- Verify keywords are in text elements (not paths or shapes)
- Ensure keyword format matches: `{{ keywordName }}`
- Check that text elements have x and y attributes

### Preview Not Showing

- Verify SVG file uploaded successfully
- Check browser console for errors
- Ensure SVG content is valid XML
- Check that SVG doesn't have external dependencies

### Search Not Working

- Implement the API endpoint for issuer search
- Check network tab for API errors
- Verify authentication tokens
- Check search query validation

## Browser Compatibility

- Chrome/Edge: Full support
- Firefox: Full support
- Safari: Full support
- Mobile browsers: Responsive design supported

## Performance Considerations

- SVG parsing is done client-side for immediate feedback
- Large SVG files may cause slow rendering (implement size limit)
- Table rendering optimized for up to ~100 issuers
- Keyword detection is O(n) where n = number of text elements

---

**Last Updated**: October 14, 2025
**Created By**: AI Assistant
**Status**: ✅ Completed - Ready for API integration
