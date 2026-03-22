import { recoverMessageAddress } from "viem";

export interface VerifySignatureParams {
    message: string;
    signature: string;
    claimedAddress: string;
}

export interface ProofVerificationResult {
    host: boolean;
    issuers: boolean[];
}

/**
 * Verifies that an Ethereum personal_sign signature was produced by the claimed address.
 * Returns false (instead of throwing) for malformed/invalid inputs.
 */
export async function verifySignature({
    message,
    signature,
    claimedAddress,
}: VerifySignatureParams): Promise<boolean> {
    try {
        const recovered = await recoverMessageAddress({
            message,
            signature: signature as `0x${string}`,
        });
        return recovered.toLowerCase() === claimedAddress.toLowerCase();
    } catch {
        return false;
    }
}

/**
 * Verifies the host and all issuer signatures in a certificate proof against
 * the shared signMessage.
 */
export async function verifyProof(proof: {
    signMessage: string;
    host: { signature: string; publicKey: string };
    issuers: { issuerSignature: string; issuerPublicKey: string }[];
}): Promise<ProofVerificationResult> {
    const [host, ...issuers] = await Promise.all([
        verifySignature({
            message: proof.signMessage,
            signature: proof.host.signature,
            claimedAddress: proof.host.publicKey,
        }),
        ...proof.issuers.map((issuer) =>
            verifySignature({
                message: proof.signMessage,
                signature: issuer.issuerSignature,
                claimedAddress: issuer.issuerPublicKey,
            }),
        ),
    ]);

    return { host, issuers };
}
