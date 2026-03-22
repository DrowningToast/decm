import { useMemo, useEffect } from "react";
import { useCertificateDetailNavStore } from "@/components/BottomNav/stores/certificates";
import { useMyCertificatesListViewModel } from "@/hooks/useMyCertificatesListViewModel";
import type {
    Certificate as CertificateData,
    CertificateShare,
} from "@/services/CertificateService/mapper";
import { useTranslation } from "react-i18next";

interface CertificateViewModel {
    id: string;
    name: string;
    event: string;
    eventId: string;
    issuedAt: string;
    status: "completed" | "queued" | "expired" | "aborted" | "pending";
    estimatedDeadline?: Date;
    certificateTitle?: string;
    certificateSubtitle?: string;
    academicInstitution?: string;
    certificateImageUrl?: string;
    certificateContractAddress?: string;
    eventContractAddress?: string;
    verifiableCredentialUrl?: string;
    shareable?: CertificateShare;
}

const mapCertificateToViewModel = (cert: CertificateData): CertificateViewModel => {
    const isAborted = cert.abortedAt != null;
    const isExpired =
        !isAborted &&
        cert.estimatedDeadline != null &&
        new Date(cert.estimatedDeadline) < new Date();
    const status = cert.certificateTokenId
        ? "completed"
        : isAborted
          ? "aborted"
          : isExpired
            ? "expired"
            : cert.estimatedDeadline
              ? "queued"
              : "pending";

    return {
        id: cert.id,
        name: cert.name || "Untitled Certificate",
        event: cert.eventName || "Unknown Event",
        eventId: cert.eventId,
        issuedAt: cert.createdAt,
        status,
        estimatedDeadline:
            status === "queued" && cert.estimatedDeadline
                ? new Date(cert.estimatedDeadline)
                : undefined,
        // Note: for "expired" status, estimatedDeadline is intentionally omitted
        certificateTitle: cert.certificateTitle,
        certificateSubtitle: cert.certificateSubtitle,
        academicInstitution: cert.academicInstitution,
        certificateImageUrl: undefined, // Fetched separately via useCertificateImage hook
        certificateContractAddress: cert.eventCertificateAddress,
        eventContractAddress: cert.eventContractAddress,
        verifiableCredentialUrl: undefined, // TODO: Add when VC URL is available
        shareable: cert.certificateShare,
    };
};

export const useCertificateDetailUsecase = (certificateId: string) => {
    const { setCertificateId, setIsClaimed } = useCertificateDetailNavStore();
    const { claimedCertificates, unclaimedCertificates, isLoading, isError } =
        useMyCertificatesListViewModel();
    const { i18n } = useTranslation();

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
        const locale = i18n.language === "th" ? "th-TH" : "en-US";
        return new Date(certificate.issuedAt).toLocaleDateString(locale, {
            year: "numeric",
            month: "short",
            day: "numeric",
        });
    }, [certificate, i18n.language]);

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
