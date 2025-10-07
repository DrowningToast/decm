import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

// Import translation files
import enTranslations from "./locales/en.json";
import thTranslations from "./locales/th.json";

// Import type definitions for TypeScript support
import "./types";

// Define available languages
export const languages = {
	en: { label: "English", flag: "🇬🇧" },
	th: { label: "ไทย", flag: "🇹🇭" },
} as const;

export type Language = keyof typeof languages;

// Initialize i18next
i18n
	.use(LanguageDetector) // Detect user language
	.use(initReactI18next) // Pass i18n instance to react-i18next
	.init({
		resources: {
			en: { translation: enTranslations },
			th: { translation: thTranslations },
		},
		fallbackLng: "en", // Fallback language
		supportedLngs: Object.keys(languages), // Supported languages
		debug: process.env.NODE_ENV === "development", // Enable debug in development
		interpolation: {
			escapeValue: false, // React already escapes values
		},
		detection: {
			// Order of language detection
			order: ["localStorage", "navigator", "htmlTag"],
			// Cache user language preference
			caches: ["localStorage"],
			lookupLocalStorage: "decm-language",
		},
	});

export default i18n;
