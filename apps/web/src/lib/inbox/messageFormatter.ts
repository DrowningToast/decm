import { z } from "zod";

/**
 * Schema for translatable message content
 * Expected format: { "en": "English text", "th": "Thai text" }
 * Thai is optional - falls back to English if missing
 */
export const translatableMessageSchema = z.object({
    en: z.string(),
    th: z.string().optional(),
});

export type TranslatableMessage = z.infer<typeof translatableMessageSchema>;

export type SupportedLanguage = "en" | "th";

/**
 * Formats an inbox message content string that contains JSON translation keys
 *
 * @param content - JSON string with translation keys (e.g., '{"en": "Hello", "th": "สวัสดี"}')
 * @param language - Current language code ("en" | "th")
 * @param fallback - Fallback content to return if parsing/validation fails
 * @returns The localized message content based on the provided language
 *
 * @example
 * // With valid JSON
 * formatInboxMessage('{"en": "Hello", "th": "สวัสดี"}', "th", "Message")
 * // Returns "สวัสดี"
 *
 * @example
 * // With missing Thai translation
 * formatInboxMessage('{"en": "Hello"}', "th", "Message")
 * // Returns "Hello" (fallback to English)
 *
 * @example
 * // With invalid JSON
 * formatInboxMessage('invalid json', "en", "Fallback Message")
 * // Returns "Fallback Message"
 */
export function formatInboxMessage(
    content: string | undefined | null,
    language: SupportedLanguage,
    fallback: string = "",
): string {
    if (!content) {
        return fallback;
    }

    try {
        const parsed: unknown = JSON.parse(content);
        const result = translatableMessageSchema.safeParse(parsed);

        if (!result.success) {
            return fallback;
        }

        const { en, th } = result.data;

        // Return Thai if available and current language is Thai, otherwise fallback to English
        if (language === "th" && th) {
            return th;
        }

        return en;
    } catch {
        // JSON parsing failed - return fallback
        return fallback;
    }
}
