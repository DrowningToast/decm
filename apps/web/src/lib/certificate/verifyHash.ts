import { keccak256, toBytes } from "viem";

export interface CertificateHashInput {
    firstName?: string;
    lastName?: string;
    academicInstitution?: string;
    certificateTitle?: string;
    certificateSubtitle?: string;
}

export function computeCertificateHash(input: CertificateHashInput): `0x${string}` {
    const name = `${input.firstName ?? ""} ${input.lastName ?? ""}`;
    const csv = [
        name,
        input.academicInstitution ?? "",
        input.certificateTitle ?? "",
        input.certificateSubtitle ?? "",
    ].join(",");
    return keccak256(toBytes(csv));
}

export function verifyCertificateHash(input: CertificateHashInput, proofHash: string): boolean {
    return computeCertificateHash(input).toLowerCase() === proofHash.toLowerCase();
}
