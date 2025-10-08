// Type definitions for i18n translations
import enTranslations from "./locales/en.json";

// Infer the translation keys from the English translation file
export type TranslationKeys = typeof enTranslations;

// Deep keys type for nested translation objects
type PathImpl<T, Key extends keyof T> = Key extends string
	? T[Key] extends Record<string, unknown>
		?
				| `${Key}.${PathImpl<T[Key], Exclude<keyof T[Key], keyof unknown[]>> & string}`
				| `${Key}.${Exclude<keyof T[Key], keyof unknown[]> & string}`
		: never
	: never;

type PathImpl2<T> = PathImpl<T, keyof T> | keyof T;

export type TranslationPath =
	PathImpl2<TranslationKeys> extends string | keyof TranslationKeys
		? PathImpl2<TranslationKeys>
		: keyof TranslationKeys;

// Extend react-i18next module with our translation keys
declare module "react-i18next" {
	interface CustomTypeOptions {
		defaultNS: "translation";
		resources: {
			translation: TranslationKeys;
		};
	}
}
