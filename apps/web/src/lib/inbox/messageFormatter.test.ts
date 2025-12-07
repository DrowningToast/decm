import { describe, it, expect } from "vitest";
import { formatInboxMessage, translatableMessageSchema } from "./messageFormatter";

describe("translatableMessageSchema", () => {
    it("should validate a message with both en and th", () => {
        const result = translatableMessageSchema.safeParse({
            en: "Hello",
            th: "สวัสดี",
        });
        expect(result.success).toBe(true);
    });

    it("should validate a message with only en", () => {
        const result = translatableMessageSchema.safeParse({
            en: "Hello",
        });
        expect(result.success).toBe(true);
    });

    it("should fail validation when en is missing", () => {
        const result = translatableMessageSchema.safeParse({
            th: "สวัสดี",
        });
        expect(result.success).toBe(false);
    });

    it("should fail validation for non-object input", () => {
        const result = translatableMessageSchema.safeParse("Hello");
        expect(result.success).toBe(false);
    });

    it("should fail validation for empty object", () => {
        const result = translatableMessageSchema.safeParse({});
        expect(result.success).toBe(false);
    });
});

describe("formatInboxMessage", () => {
    describe("with valid JSON and both translations", () => {
        const validJson = '{"en": "Hello World", "th": "สวัสดีโลก"}';

        it("should return English when language is en", () => {
            const result = formatInboxMessage(validJson, "en", "Fallback");
            expect(result).toBe("Hello World");
        });

        it("should return Thai when language is th", () => {
            const result = formatInboxMessage(validJson, "th", "Fallback");
            expect(result).toBe("สวัสดีโลก");
        });
    });

    describe("with valid JSON and only English", () => {
        const englishOnlyJson = '{"en": "Hello World"}';

        it("should return English when language is en", () => {
            const result = formatInboxMessage(englishOnlyJson, "en", "Fallback");
            expect(result).toBe("Hello World");
        });

        it("should fall back to English when language is th but th is missing", () => {
            const result = formatInboxMessage(englishOnlyJson, "th", "Fallback");
            expect(result).toBe("Hello World");
        });
    });

    describe("with invalid input", () => {
        it("should return fallback for invalid JSON", () => {
            const result = formatInboxMessage("invalid json", "en", "Fallback Message");
            expect(result).toBe("Fallback Message");
        });

        it("should return fallback for null content", () => {
            const result = formatInboxMessage(null, "en", "Fallback Message");
            expect(result).toBe("Fallback Message");
        });

        it("should return fallback for undefined content", () => {
            const result = formatInboxMessage(undefined, "en", "Fallback Message");
            expect(result).toBe("Fallback Message");
        });

        it("should return fallback for empty string", () => {
            const result = formatInboxMessage("", "en", "Fallback Message");
            expect(result).toBe("Fallback Message");
        });

        it("should return fallback for JSON without en key", () => {
            const result = formatInboxMessage('{"th": "สวัสดี"}', "en", "Fallback Message");
            expect(result).toBe("Fallback Message");
        });

        it("should return fallback for JSON with non-string values", () => {
            const result = formatInboxMessage('{"en": 123}', "en", "Fallback Message");
            expect(result).toBe("Fallback Message");
        });

        it("should return fallback for JSON array", () => {
            const result = formatInboxMessage('["Hello", "World"]', "en", "Fallback Message");
            expect(result).toBe("Fallback Message");
        });
    });

    describe("edge cases", () => {
        it("should return empty string when no fallback is provided and content is null", () => {
            const result = formatInboxMessage(null, "en");
            expect(result).toBe("");
        });

        it("should handle empty strings in translations", () => {
            const result = formatInboxMessage('{"en": "", "th": "สวัสดี"}', "en", "Fallback");
            expect(result).toBe("");
        });

        it("should handle unicode characters correctly", () => {
            const unicodeJson = '{"en": "Hello 👋", "th": "สวัสดี 🙏"}';
            expect(formatInboxMessage(unicodeJson, "en", "Fallback")).toBe("Hello 👋");
            expect(formatInboxMessage(unicodeJson, "th", "Fallback")).toBe("สวัสดี 🙏");
        });

        it("should handle long text correctly", () => {
            const longText = "A".repeat(10000);
            const longJson = JSON.stringify({ en: longText, th: "Thai" });
            const result = formatInboxMessage(longJson, "en", "Fallback");
            expect(result).toBe(longText);
        });

        it("should handle special characters in JSON", () => {
            const specialJson = '{"en": "Hello \\"World\\"", "th": "สวัสดี \\"โลก\\""}';
            expect(formatInboxMessage(specialJson, "en", "Fallback")).toBe('Hello "World"');
            expect(formatInboxMessage(specialJson, "th", "Fallback")).toBe('สวัสดี "โลก"');
        });
    });
});
