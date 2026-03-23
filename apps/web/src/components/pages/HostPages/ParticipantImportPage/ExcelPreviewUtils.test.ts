import { describe, it, expect } from "vitest";
import {
    validateRow,
    buildParticipant,
    ALWAYS_REQUIRED_COLUMNS,
    EITHER_OR_COLUMNS,
    OPTIONAL_COLUMNS,
} from "./ExcelPreviewUtils";

describe("validateRow", () => {
    it("should be valid with first_name, last_name, and email only", () => {
        const row = { first_name: "John", last_name: "Doe", email: "john@example.com" };
        const result = validateRow(row);
        expect(result.isValid).toBe(true);
        expect(result.missingFields).toHaveLength(0);
    });

    it("should be valid with first_name, last_name, and wallet_address only", () => {
        const row = { first_name: "John", last_name: "Doe", wallet_address: "0xabc123" };
        const result = validateRow(row);
        expect(result.isValid).toBe(true);
        expect(result.missingFields).toHaveLength(0);
    });

    it("should be invalid when both email and wallet_address are provided", () => {
        const row = {
            first_name: "John",
            last_name: "Doe",
            email: "john@example.com",
            wallet_address: "0xabc123",
        };
        const result = validateRow(row);
        expect(result.isValid).toBe(false);
        expect(result.missingFields).toContain("email_or_wallet_address");
    });

    it("should be invalid when neither email nor wallet_address is provided", () => {
        const row = { first_name: "John", last_name: "Doe" };
        const result = validateRow(row);
        expect(result.isValid).toBe(false);
        expect(result.missingFields).toContain("email_or_wallet_address");
    });

    it("should be invalid when first_name is missing", () => {
        const row = { last_name: "Doe", email: "john@example.com" };
        const result = validateRow(row);
        expect(result.isValid).toBe(false);
        expect(result.missingFields).toContain(ALWAYS_REQUIRED_COLUMNS.firstName);
    });

    it("should be invalid when last_name is missing", () => {
        const row = { first_name: "John", email: "john@example.com" };
        const result = validateRow(row);
        expect(result.isValid).toBe(false);
        expect(result.missingFields).toContain(ALWAYS_REQUIRED_COLUMNS.lastName);
    });

    it("should be invalid when first_name is whitespace only", () => {
        const row = { first_name: "   ", last_name: "Doe", email: "john@example.com" };
        const result = validateRow(row);
        expect(result.isValid).toBe(false);
        expect(result.missingFields).toContain(ALWAYS_REQUIRED_COLUMNS.firstName);
    });

    it("should be valid with optional fields included", () => {
        const row = {
            first_name: "John",
            last_name: "Doe",
            email: "john@example.com",
            phone_number: "+123",
            academic_institution: "MIT",
        };
        const result = validateRow(row);
        expect(result.isValid).toBe(true);
    });
});

describe("buildParticipant", () => {
    it("should build participant with email", () => {
        const row = { first_name: "John", last_name: "Doe", email: "john@example.com" };
        const participant = buildParticipant(row);
        expect(participant.first_name).toBe("John");
        expect(participant.last_name).toBe("Doe");
        expect(participant.email).toBe("john@example.com");
        expect(participant.wallet_address).toBeUndefined();
    });

    it("should build participant with wallet_address", () => {
        const row = { first_name: "John", last_name: "Doe", wallet_address: "0xabc123" };
        const participant = buildParticipant(row);
        expect(participant.first_name).toBe("John");
        expect(participant.last_name).toBe("Doe");
        expect(participant.wallet_address).toBe("0xabc123");
        expect(participant.email).toBeUndefined();
    });

    it("should include phone_number when provided", () => {
        const row = {
            first_name: "John",
            last_name: "Doe",
            email: "john@example.com",
            phone_number: "+123",
        };
        const participant = buildParticipant(row);
        expect(participant.phone_number).toBe("+123");
    });

    it("should omit phone_number when empty", () => {
        const row = {
            first_name: "John",
            last_name: "Doe",
            email: "john@example.com",
            phone_number: "",
        };
        const participant = buildParticipant(row);
        expect(participant.phone_number).toBeUndefined();
    });

    it("should include academic_institution when provided", () => {
        const row = {
            first_name: "John",
            last_name: "Doe",
            email: "john@example.com",
            academic_institution: "MIT",
        };
        const participant = buildParticipant(row);
        expect(participant.academic_institution).toBe("MIT");
    });
});

describe("column definitions", () => {
    it("ALWAYS_REQUIRED_COLUMNS should contain first_name and last_name", () => {
        expect(Object.values(ALWAYS_REQUIRED_COLUMNS)).toContain("first_name");
        expect(Object.values(ALWAYS_REQUIRED_COLUMNS)).toContain("last_name");
    });

    it("EITHER_OR_COLUMNS should contain email and wallet_address", () => {
        expect(Object.values(EITHER_OR_COLUMNS)).toContain("email");
        expect(Object.values(EITHER_OR_COLUMNS)).toContain("wallet_address");
    });

    it("OPTIONAL_COLUMNS should contain phone_number and academic_institution", () => {
        expect(Object.values(OPTIONAL_COLUMNS)).toContain("phone_number");
        expect(Object.values(OPTIONAL_COLUMNS)).toContain("academic_institution");
    });
});
