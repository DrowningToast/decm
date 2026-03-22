import { describe, it, expect } from "vitest";
import { keccak256, toBytes } from "viem";
import { computeCertificateHash, verifyCertificateHash } from "./verifyHash";

// Helper: the exact CSV string the backend produces, so tests document the contract.
function backendCsv(
    firstName: string,
    lastName: string,
    academicInstitution: string,
    certificateTitle: string,
    certificateSubtitle: string,
): string {
    const name = `${firstName} ${lastName}`;
    return `${name},${academicInstitution},${certificateTitle},${certificateSubtitle}`;
}

describe("computeCertificateHash", () => {
    it("matches backend output for full inputs", () => {
        const input = {
            firstName: "Alice",
            lastName: "Smith",
            academicInstitution: "MIT",
            certificateTitle: "Certificate of Completion",
            certificateSubtitle: "Blockchain Development",
        };
        const expected = keccak256(
            toBytes(
                backendCsv(
                    "Alice",
                    "Smith",
                    "MIT",
                    "Certificate of Completion",
                    "Blockchain Development",
                ),
            ),
        );
        expect(computeCertificateHash(input)).toBe(expected);
    });

    it("joins name as 'firstName lastName' with a single space", () => {
        // Backend: fmt.Sprintf("%s %s", firstName, lastName) — always a space between
        const result = computeCertificateHash({
            firstName: "John",
            lastName: "Doe",
            certificateTitle: "Award",
        });
        const expected = keccak256(toBytes(backendCsv("John", "Doe", "", "Award", "")));
        expect(result).toBe(expected);
    });

    it("uses a space for name when both firstName and lastName are undefined", () => {
        // Go: fmt.Sprintf("%s %s", "", "") == " "
        const result = computeCertificateHash({ certificateTitle: "T" });
        const expected = keccak256(toBytes(backendCsv("", "", "", "T", "")));
        expect(result).toBe(expected);
    });

    it("treats undefined fields the same as empty strings", () => {
        const withUndefined = computeCertificateHash({
            firstName: "Bob",
            lastName: "Lee",
        });
        const withEmpty = computeCertificateHash({
            firstName: "Bob",
            lastName: "Lee",
            academicInstitution: "",
            certificateTitle: "",
            certificateSubtitle: "",
        });
        expect(withUndefined).toBe(withEmpty);
    });

    it("produces different hashes for different inputs", () => {
        const a = computeCertificateHash({
            firstName: "Alice",
            lastName: "Smith",
            certificateTitle: "A",
        });
        const b = computeCertificateHash({
            firstName: "Bob",
            lastName: "Jones",
            certificateTitle: "A",
        });
        expect(a).not.toBe(b);
    });

    it("is sensitive to field order — swapping title and subtitle changes the hash", () => {
        const a = computeCertificateHash({ certificateTitle: "Title", certificateSubtitle: "Sub" });
        const b = computeCertificateHash({ certificateTitle: "Sub", certificateSubtitle: "Title" });
        expect(a).not.toBe(b);
    });

    it("handles Thai (multibyte UTF-8) characters correctly", () => {
        const title = "ได้รับรางวัลแข่งขัน Application Development Workshop";
        const result = computeCertificateHash({
            firstName: "สมชาย",
            lastName: "ใจดี",
            certificateTitle: title,
        });
        const expected = keccak256(toBytes(backendCsv("สมชาย", "ใจดี", "", title, "")));
        expect(result).toBe(expected);
    });

    it("always returns a 0x-prefixed 32-byte hex string", () => {
        const result = computeCertificateHash({ firstName: "X" });
        expect(result).toMatch(/^0x[0-9a-f]{64}$/i);
    });
});

describe("verifyCertificateHash", () => {
    it("returns true when computed hash matches the proof hash", () => {
        const input = {
            firstName: "Alice",
            lastName: "Smith",
            academicInstitution: "MIT",
            certificateTitle: "Certificate of Completion",
            certificateSubtitle: "",
        };
        const proofHash = computeCertificateHash(input);
        expect(verifyCertificateHash(input, proofHash)).toBe(true);
    });

    it("returns false when the proof hash does not match", () => {
        const input = { firstName: "Alice", lastName: "Smith", certificateTitle: "Cert" };
        const wrongHash = "0x" + "ab".repeat(32);
        expect(verifyCertificateHash(input, wrongHash)).toBe(false);
    });

    it("is case-insensitive for the proof hash", () => {
        const input = { firstName: "Alice", lastName: "Smith", certificateTitle: "Cert" };
        const proofHash = computeCertificateHash(input);
        expect(verifyCertificateHash(input, proofHash.toUpperCase())).toBe(true);
        expect(verifyCertificateHash(input, proofHash.toLowerCase())).toBe(true);
    });
});
