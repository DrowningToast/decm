# Internationalization (i18n) Guide

This guide explains how to use the i18n system in the DECM web application.

## Overview

The application uses [react-i18next](https://react.i18next.com/) for internationalization, providing:

- Multiple language support (English, Thai)
- Automatic language detection
- Type-safe translations
- Easy language switching

## Supported Languages

- **English (en)** - Default language
- **Thai (th)** - Thai language

## Quick Start

### Using Translations in Components

```tsx
import { useTranslation } from "react-i18next";

function MyComponent() {
	const { t } = useTranslation();

	return (
		<div>
			<h1>{t("common.welcome")}</h1>
			<p>{t("home.hero.subtitle")}</p>
		</div>
	);
}
```

### Adding Language Switcher

```tsx
import { LanguageSwitcher } from "@/components/LanguageSwitcher";

function Header() {
	return (
		<header>
			<nav>{/* Your navigation items */}</nav>
			<LanguageSwitcher />
		</header>
	);
}
```

## Translation Files

Translation files are located in `apps/web/src/lib/i18n/locales/`:

```
locales/
├── en.json  # English translations
└── th.json  # Thai translations
```

### Structure

All translation files follow the same structure:

```json
{
	"common": {
		"welcome": "Welcome",
		"loading": "Loading..."
	},
	"auth": {
		"signIn": "Sign In",
		"signUp": "Sign Up"
	},
	"home": {
		"hero": {
			"title": "Your Title",
			"subtitle": "Your Subtitle"
		}
	}
}
```

## Adding New Translations

### 1. Add to Translation Files

Add the same key to both `en.json` and `th.json`:

**en.json:**

```json
{
	"myFeature": {
		"title": "My Feature",
		"description": "This is a description"
	}
}
```

**th.json:**

```json
{
	"myFeature": {
		"title": "ฟีเจอร์ของฉัน",
		"description": "นี่คือคำอธิบาย"
	}
}
```

### 2. Use in Component

```tsx
function MyFeature() {
	const { t } = useTranslation();

	return (
		<div>
			<h2>{t("myFeature.title")}</h2>
			<p>{t("myFeature.description")}</p>
		</div>
	);
}
```

## Advanced Usage

### Interpolation

Use variables in translations:

**Translation file:**

```json
{
	"welcome": "Welcome, {{name}}!"
}
```

**Component:**

```tsx
const { t } = useTranslation();
t("welcome", { name: "John" }); // "Welcome, John!"
```

### Pluralization

Handle singular/plural forms:

**Translation file:**

```json
{
	"itemCount": "{{count}} item",
	"itemCount_plural": "{{count}} items"
}
```

**Component:**

```tsx
const { t } = useTranslation();
t("itemCount", { count: 1 }); // "1 item"
t("itemCount", { count: 5 }); // "5 items"
```

### Accessing i18n Instance

For more control, access the i18n instance directly:

```tsx
import { useTranslation } from "react-i18next";

function LanguageInfo() {
	const { i18n } = useTranslation();

	const changeLanguage = (lng: string) => {
		i18n.changeLanguage(lng);
	};

	return (
		<div>
			<p>Current language: {i18n.language}</p>
			<button onClick={() => changeLanguage("th")}>Switch to Thai</button>
		</div>
	);
}
```

## Language Detection

The application automatically detects the user's language in this order:

1. **localStorage** - Saved preference (`decm-language` key)
2. **Browser navigator** - Browser language settings
3. **HTML tag** - Document language attribute
4. **Fallback** - English (en)

## Type Safety

The i18n system is configured with TypeScript for type-safe translation keys:

```tsx
// ✅ Valid - TypeScript will autocomplete
t("common.welcome");
t("home.hero.title");

// ❌ Invalid - TypeScript will show error
t("nonexistent.key");
```

## Configuration

The i18n configuration is in `apps/web/src/lib/i18n/config.ts`:

```typescript
i18n
	.use(LanguageDetector)
	.use(initReactI18next)
	.init({
		resources: {
			/* translations */
		},
		fallbackLng: "en",
		supportedLngs: ["en", "th"],
		debug: process.env.NODE_ENV === "development",
		interpolation: {
			escapeValue: false, // React already escapes
		},
		detection: {
			order: ["localStorage", "navigator", "htmlTag"],
			caches: ["localStorage"],
			lookupLocalStorage: "decm-language",
		},
	});
```

## Adding a New Language

### 1. Create Translation File

Create a new file in `locales/` (e.g., `ja.json` for Japanese):

```json
{
	"common": {
		"welcome": "ようこそ"
	}
}
```

### 2. Update Configuration

**config.ts:**

```typescript
import jaTranslations from "./locales/ja.json";

export const languages = {
	en: { label: "English", flag: "🇬🇧" },
	th: { label: "ไทย", flag: "🇹🇭" },
	ja: { label: "日本語", flag: "🇯🇵" }, // Add new language
} as const;

i18n.init({
	resources: {
		en: { translation: enTranslations },
		th: { translation: thTranslations },
		ja: { translation: jaTranslations }, // Add new resource
	},
	supportedLngs: ["en", "th", "ja"], // Add to supported languages
	// ... rest of config
});
```

### 3. Update Language Switcher

The `LanguageSwitcher` component automatically picks up new languages from the `languages` object.

## Best Practices

### 1. Use Namespaces

Group related translations together:

```json
{
	"auth": {
		/* authentication related */
	},
	"profile": {
		/* profile related */
	},
	"events": {
		/* events related */
	}
}
```

### 2. Consistent Keys

Use consistent naming across languages:

```json
// ✅ Good - same structure
// en.json
{ "auth": { "signIn": "Sign In" } }

// th.json
{ "auth": { "signIn": "เข้าสู่ระบบ" } }

// ❌ Bad - different structure
// en.json
{ "auth": { "signIn": "Sign In" } }

// th.json
{ "authentication": { "login": "เข้าสู่ระบบ" } }
```

### 3. Avoid Hard-coded Strings

```tsx
// ❌ Bad - hard-coded
<button>Sign In</button>

// ✅ Good - translated
<button>{t('auth.signIn')}</button>
```

### 4. Provide Context

Use descriptive keys that provide context:

```json
{
	"button": {
		"save": "Save",
		"cancel": "Cancel"
	},
	"validation": {
		"required": "This field is required"
	}
}
```

### 5. Keep Translations Short

For UI elements, keep translations concise:

```json
{
	"nav": {
		"home": "Home",
		"events": "Events",
		"profile": "Profile"
	}
}
```

## Testing Translations

### Check Missing Keys

Run in development mode to see debug logs:

```bash
pnpm dev
```

Missing translation keys will be logged to the console.

### Test Language Switching

1. Start the dev server
2. Open the application
3. Use the LanguageSwitcher component
4. Verify all translations appear correctly

## Common Translation Keys

The application includes pre-configured translation keys for:

- **common** - UI elements, actions
- **nav** - Navigation items
- **auth** - Authentication flow
- **signup** - Sign up page
- **home** - Home page content
- **profile** - User profile
- **events** - Event management
- **credentials** - Digital credentials
- **portfolio** - User portfolio
- **validation** - Form validation messages
- **errors** - Error messages

## Resources

- [react-i18next Documentation](https://react.i18next.com/)
- [i18next Documentation](https://www.i18next.com/)
- [Language Detection Plugin](https://github.com/i18next/i18next-browser-languageDetector)

## Support

For questions or issues with translations:

1. Check this documentation
2. Review the translation files
3. Check the i18n configuration
4. Consult the react-i18next documentation
