import { describe, it, expect } from "vitest";
import {
    validateCertificateRow,
    buildCertificate,
    CERTIFICATE_EITHER_OR_COLUMNS,
    CERTIFICATE_OPTIONAL_COLUMNS,
} from "./ExcelPreviewUtils";

describe("validateCertificateRow", () => {
    it("should be valid with email only", () => {
        const row = { email: "alice@example.com" };
        const result = validateCertificateRow(row);
        expect(result.isValid).toBe(true);
        expect(result.missingFields).toHaveLength(0);
    });

    it("should be valid with wallet_address only", () => {
        const row = { wallet_address: "0x1234567890abcdef1234567890abcdef12345678" };
        const result = validateCertificateRow(row);
        expect(result.isValid).toBe(true);
        expect(result.missingFields).toHaveLength(0);
    });

    it("should be invalid when both email and wallet_address are provided", () => {
        const row = { email: "alice@example.com", wallet_address: "0xabc" };
        const result = validateCertificateRow(row);
        expect(result.isValid).toBe(false);
        expect(result.missingFields).toContain("email_or_wallet_address");
    });

    it("should be invalid when neither email nor wallet_address is provided", () => {
        const row = { first_name: "Alice", last_name: "Smith" };
        const result = validateCertificateRow(row);
        expect(result.isValid).toBe(false);
        expect(result.missingFields).toContain("email_or_wallet_address");
    });

    it("should be valid with optional fields included alongside email", () => {
        const row = {
            email: "alice@example.com",
            first_name: "Alice",
            last_name: "Smith",
            academic_institution: "MIT",
            certificate_title: "Best Award",
            certificate_subtitle: "For Excellence",
        };
        const result = validateCertificateRow(row);
        expect(result.isValid).toBe(true);
    });
});

describe("buildCertificate", () => {
    it("should build with email", () => {
        const row = { email: "alice@example.com", first_name: "Alice", last_name: "Smith" };
        const cert = buildCertificate(row);
        expect(cert.email).toBe("alice@example.com");
        expect(cert.wallet_address).toBeUndefined();
        expect(cert.first_name).toBe("Alice");
        expect(cert.last_name).toBe("Smith");
    });

    it("should build with wallet_address", () => {
        const row = { wallet_address: "0xabc", first_name: "Alice", last_name: "Smith" };
        const cert = buildCertificate(row);
        expect(cert.wallet_address).toBe("0xabc");
        expect(cert.email).toBeUndefined();
    });

    it("should include optional fields when provided", () => {
        const row = {
            email: "alice@example.com",
            academic_institution: "MIT",
            certificate_title: "Best Award",
            certificate_subtitle: "For Excellence",
        };
        const cert = buildCertificate(row);
        expect(cert.academic_institution).toBe("MIT");
        expect(cert.certificate_title).toBe("Best Award");
        expect(cert.certificate_subtitle).toBe("For Excellence");
    });

    it("should leave optional string fields empty when not provided", () => {
        const row = { email: "alice@example.com", academic_institution: "" };
        const cert = buildCertificate(row);
        // academic_institution is a required string in the API type — empty string when not filled
        expect(cert.academic_institution).toBe("");
    });
});

describe("column definitions", () => {
    it("CERTIFICATE_EITHER_OR_COLUMNS should contain email and wallet_address", () => {
        expect(Object.values(CERTIFICATE_EITHER_OR_COLUMNS)).toContain("email");
        expect(Object.values(CERTIFICATE_EITHER_OR_COLUMNS)).toContain("wallet_address");
    });

    it("CERTIFICATE_OPTIONAL_COLUMNS should contain first_name, last_name, academic_institution, certificate_title, certificate_subtitle", () => {
        expect(Object.values(CERTIFICATE_OPTIONAL_COLUMNS)).toContain("first_name");
        expect(Object.values(CERTIFICATE_OPTIONAL_COLUMNS)).toContain("last_name");
        expect(Object.values(CERTIFICATE_OPTIONAL_COLUMNS)).toContain("academic_institution");
        expect(Object.values(CERTIFICATE_OPTIONAL_COLUMNS)).toContain("certificate_title");
        expect(Object.values(CERTIFICATE_OPTIONAL_COLUMNS)).toContain("certificate_subtitle");
    });
});
