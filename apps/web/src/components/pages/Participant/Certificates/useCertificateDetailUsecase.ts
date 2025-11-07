import { useMemo } from "react";
import { useCertificateDetailNavStore } from "@/components/Botto/stores";

interface Certificate {
    id: string;
    name: string;
    event: string;
    eventId: string;
    issuedAt: string;
    status: "completed" | "pending";
    description?: string;
    certificateImageUrl?: string;
    certificateContractAddress?: string;
    eventContractAddress?: string;
    verifiableCredentialUrl?: string;
}

// Mock data - TODO: Replace with API call
const mockCertificates: Record<string, Certificate> = {
    "1": {
        id: "1",
        name: "Participation award",
        event: "ToBeIT69",
        eventId: "1",
        issuedAt: "2024-09-24",
        status: "completed",
        description:
            "Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum ",
        certificateImageUrl: undefined, // Placeholder for certificate image
        certificateContractAddress: "0x0000...0000",
        eventContractAddress: "0x0000...0000",
        verifiableCredentialUrl: "#",
    },
    "2": {
        id: "2",
        name: "Winning team award",
        event: "ToBeIT69",
        eventId: "1",
        issuedAt: "2024-09-24",
        status: "pending",
        description:
            "Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum ",
        certificateImageUrl: undefined,
        certificateContractAddress: "0x0000...0000",
        eventContractAddress: "0x0000...0000",
        verifiableCredentialUrl: "#",
    },
};

export const useCertificateDetailUsecase = (certificateId: string) => {
    const { setCertificateId } = useCertificateDetailNavStore();

    const certificate = useMemo(() => {
        // TODO: Fetch from API based on certificateId
        return mockCertificates[certificateId] || null;
    }, [certificateId]);

    const formattedDate = useMemo(() => {
        if (!certificate) return "";
        return new Date(certificate.issuedAt).toLocaleDateString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
        });
    }, [certificate]);

    // Set the certificate ID in the nav store when component mounts
    useMemo(() => {
        setCertificateId(certificateId);
        return () => {
            setCertificateId(null);
        };
    }, [certificateId, setCertificateId]);

    return {
        certificate,
        formattedDate,
    };
};
