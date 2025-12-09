import { useMemo, useEffect } from "react";
import { useCertificateDetailNavStore } from "@/components/BottomNav/stores/certificates";
import { useMyCertificatesListViewModel } from "@/hooks/useMyCertificatesListViewModel";
import type { Certificate as CertificateData } from "@/services/CertificateService/mapper";

interface CertificateViewModel {
    id: string;
    name: string;
    event: string;
    eventId: string;
    issuedAt: string;
    status: "completed" | "pending";
    certificateTitle?: string;
    certificateSubtitle?: string;
    academicInstitution?: string;
    certificateImageUrl?: string;
    certificateContractAddress?: string;
    eventContractAddress?: string;
    verifiableCredentialUrl?: string;
}

const mapCertificateToViewModel = (cert: CertificateData): CertificateViewModel => {
    return {
        id: cert.id,
        name: cert.name || "Untitled Certificate",
        event: cert.eventName || "Unknown Event",
        eventId: cert.eventId,
        issuedAt: cert.createdAt,
        status: cert.certificateTokenId ? "completed" : "pending",
        certificateTitle: cert.certificateTitle,
        certificateSubtitle: cert.certificateSubtitle,
        academicInstitution: cert.academicInstitution,
        certificateImageUrl: undefined, // Fetched separately via useCertificateImage hook
        certificateContractAddress: cert.eventCertificateAddress,
        eventContractAddress: cert.eventContractAddress,
        verifiableCredentialUrl: undefined, // TODO: Add when VC URL is available
    };
};

export const useCertificateDetailUsecase = (certificateId: string) => {
    const { setCertificateId, setIsClaimed } = useCertificateDetailNavStore();
    const { claimedCertificates, unclaimedCertificates, isLoading, isError } =
        useMyCertificatesListViewModel();

    const certificate = useMemo(() => {
        // Combine claimed and unclaimed certificates
        const allCertificates = [...claimedCertificates, ...unclaimedCertificates];

        // Find certificate by ID
        const foundCert = allCertificates.find((cert) => cert.id === certificateId);

        if (!foundCert) return null;

        return mapCertificateToViewModel(foundCert);
    }, [certificateId, claimedCertificates, unclaimedCertificates]);

    const formattedDate = useMemo(() => {
        if (!certificate) return "";
        return new Date(certificate.issuedAt).toLocaleDateString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
        });
    }, [certificate]);

    // Set the certificate ID and claimed status in the nav store when component mounts
    useEffect(() => {
        setCertificateId(certificateId);
        setIsClaimed(certificate?.status === "completed");
        return () => {
            setCertificateId(null);
            setIsClaimed(false);
        };
    }, [certificateId, certificate?.status, setCertificateId, setIsClaimed]);

    return {
        certificate,
        formattedDate,
        isLoading,
        isError,
    };
};
