import { describe, it, expect } from "vitest";
import { verifySignature, verifyProof } from "./verifySignature";

// Real values from the certificate share API
const SIGN_MESSAGE =
    '{"eventContractAddress":"0x428e6c4258cA560bbb56D83F0cfbAD8dc33f6696","receivers":["0x31c1ba958f8f25fafbf8686d7107c700d11575f9dbe5137c2c6747cf584db7ab"]}';
const REAL_SIGNATURE =
    "0x40f70443932b022dfaff2890498d750410b63bb6ae52f327ac8a8c5eef316c9d13c3647130efdf47a0169f5c38d49d7c59089bc46e3291746b907e203790240b1b";
const REAL_PUBLIC_KEY = "0x46494F89533057aD6865B86D9619aCd9A3cf7687";

// A signature made by a different key (wrong signer)
const OTHER_SIGNATURE =
    "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1b";

describe("verifySignature", () => {
    it("returns true when signature was made by the claimed address", async () => {
        const result = await verifySignature({
            message: SIGN_MESSAGE,
            signature: REAL_SIGNATURE,
            claimedAddress: REAL_PUBLIC_KEY,
        });
        expect(result).toBe(true);
    });

    it("returns true regardless of address checksum casing", async () => {
        const result = await verifySignature({
            message: SIGN_MESSAGE,
            signature: REAL_SIGNATURE,
            claimedAddress: REAL_PUBLIC_KEY.toLowerCase(),
        });
        expect(result).toBe(true);
    });

    it("returns false when the claimed address does not match the signer", async () => {
        const result = await verifySignature({
            message: SIGN_MESSAGE,
            signature: REAL_SIGNATURE,
            claimedAddress: "0x0000000000000000000000000000000000000001",
        });
        expect(result).toBe(false);
    });

    it("returns false for a malformed signature", async () => {
        const result = await verifySignature({
            message: SIGN_MESSAGE,
            signature: "0xdeadbeef",
            claimedAddress: REAL_PUBLIC_KEY,
        });
        expect(result).toBe(false);
    });

    it("returns false for an empty signature", async () => {
        const result = await verifySignature({
            message: SIGN_MESSAGE,
            signature: "",
            claimedAddress: REAL_PUBLIC_KEY,
        });
        expect(result).toBe(false);
    });
});

describe("verifyProof", () => {
    const proof = {
        signMessage: SIGN_MESSAGE,
        host: {
            signature: REAL_SIGNATURE,
            publicKey: REAL_PUBLIC_KEY,
        },
        issuers: [
            {
                issuerSignature: REAL_SIGNATURE,
                issuerPublicKey: REAL_PUBLIC_KEY,
            },
        ],
    };

    it("returns verified=true for host and all issuers when signatures are valid", async () => {
        const result = await verifyProof(proof);
        expect(result.host).toBe(true);
        expect(result.issuers).toEqual([true]);
    });

    it("returns verified=false for host when host signature is wrong", async () => {
        const result = await verifyProof({
            ...proof,
            host: { signature: OTHER_SIGNATURE, publicKey: REAL_PUBLIC_KEY },
        });
        expect(result.host).toBe(false);
        expect(result.issuers).toEqual([true]);
    });

    it("returns verified=false for the issuer whose signature is wrong", async () => {
        const result = await verifyProof({
            ...proof,
            issuers: [{ issuerSignature: OTHER_SIGNATURE, issuerPublicKey: REAL_PUBLIC_KEY }],
        });
        expect(result.host).toBe(true);
        expect(result.issuers).toEqual([false]);
    });

    it("handles multiple issuers independently", async () => {
        const result = await verifyProof({
            ...proof,
            issuers: [
                { issuerSignature: REAL_SIGNATURE, issuerPublicKey: REAL_PUBLIC_KEY },
                { issuerSignature: OTHER_SIGNATURE, issuerPublicKey: REAL_PUBLIC_KEY },
            ],
        });
        expect(result.issuers).toEqual([true, false]);
    });
});
