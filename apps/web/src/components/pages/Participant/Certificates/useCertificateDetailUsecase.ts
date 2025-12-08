import { useMemo, useEffect } from "react";
import { useCertificateDetailNavStore } from "@/components/BottomNav/stores/certificates";
import { useMyCertificatesListViewModel } from "@/hooks/useMyCertificatesListViewModel";
import type { EntityEventCertificate } from "@decm/api";

interface Certificate {
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

const mapEventCertificateToViewModel = (cert: EntityEventCertificate): Certificate => {
    return {
        id: cert.id || "",
        name: cert.name || "Untitled Certificate",
        event: cert.event_name || "Unknown Event",
        eventId: cert.event_id || "",
        issuedAt: cert.created_at || new Date().toISOString(),
        status: cert.certificate_token_id ? "completed" : "pending",
        certificateTitle: cert.certificate_title || undefined,
        certificateSubtitle: cert.certificate_subtitle || undefined,
        academicInstitution: cert.academic_institution || undefined,
        certificateImageUrl: undefined, // Fetched separately via useCertificateImage hook
        certificateContractAddress: cert.event_certificate_address || undefined,
        eventContractAddress: cert.event_contract_address || undefined,
        verifiableCredentialUrl: undefined, // TODO: Add when VC URL is available
    };
};

export const useCertificateDetailUsecase = (certificateId: string) => {
    const { setCertificateId, setIsClaimed } = useCertificateDetailNavStore();
    const { claimedCertificates, unclaimedCertificates, isLoading, isError } =
        useMyCertificatesListViewModel();

    const certificate = useMemo(() => {
        // Combine claimed and unclaimed certificates
        const allCertificates = [...(claimedCertificates || []), ...(unclaimedCertificates || [])];

        // Find certificate by ID
        const foundCert = allCertificates.find((cert) => cert.id === certificateId);

        if (!foundCert) return null;

        return mapEventCertificateToViewModel(foundCert);
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
