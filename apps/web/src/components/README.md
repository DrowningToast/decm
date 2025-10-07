# DECM Components - Figma Design System

This directory contains UI components based on the Figma design system for the DECM platform.

## Components

### Button Component

Located at: `src/components/ui/button.tsx`

The button component includes both standard variants and Figma design system variants.

#### Figma Design Variants

```tsx
import { Button } from '@/components/ui/button';

// Primary Button (Red background, light text)
<Button variant="primary" size="xl">
  Sign in with Web3 Wallet Provider
</Button>

// Secondary Dark Button (Dark background, light text)
<Button variant="secondary-dark" size="xl">
  Sign up with Google Account
</Button>

// Secondary Light Button (Light background, red text)
<Button variant="secondary-light" size="xl">
  Start building your portfolio
</Button>
```

#### Standard Variants

```tsx
<Button variant="default">Default</Button>
<Button variant="destructive">Destructive</Button>
<Button variant="outline">Outline</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="link">Link</Button>
```

#### Sizes

```tsx
<Button size="sm">Small</Button>
<Button size="default">Default</Button>
<Button size="lg">Large</Button>
<Button size="xl">Extra Large</Button> {/* Best for Figma variants */}
<Button size="icon">Icon Only</Button>
```

### Navbar Component

Located at: `src/components/layouts/navigations/PublicNavbar.tsx`

A responsive navbar with three color variants matching the Figma design system.

#### Usage

```tsx
import { PublicNavbar } from '@/components/layouts/navigations/PublicNavbar';

// Primary variant (red background)
<PublicNavbar variant="primary" />

// Secondary dark variant (dark background)
<PublicNavbar variant="secondary-dark" />

// Secondary light variant (light background)
<PublicNavbar variant="secondary-light" />
```

#### Features

- Responsive design (mobile & desktop)
- Hamburger menu for mobile devices
- Three navigation links: Features, About, Sign In
- Smooth transitions and hover effects
- Matches Figma color palette

#### Navigation Links

The navbar includes these default links:
- Features (`#features`)
- About (`#about`)
- Sign In (`/signin`)

You can modify these in the component file as needed.

## Design System Colors

The components use these colors from the Figma design:

- **Primary**: `#eb5331` (Red/Orange)
- **Secondary Dark**: `#362927` (Dark Brown)
- **Secondary Light**: `#e9dede` (Light Beige)
- **Foreground Alt**: `#e9dede` (Light text)
- **Foreground**: `#fcfcfc` (White text)

## Example Page

See `src/components/examples/ComponentExamples.tsx` for complete usage examples.

## Notes

- All components use Tailwind CSS for styling
- Components support the `className` prop for custom styling
- Button variants include text shadow effects as per Figma design
- Navbar is fixed to the top by default (use `className="relative"` to override)
- Components are fully typed with TypeScript

