import { useMutation, useQueryClient } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { QUERY_KEY } from "@/lib/queryKeys";
import type {
    ClaimCertificateWithPinParams,
    ClaimCertificateWithSignatureParams,
} from "@/services/CertificateService/mapper";

type ClaimCertificateParams = ClaimCertificateWithPinParams | ClaimCertificateWithSignatureParams;

/**
 * Hook for claiming/minting a certificate as a participant
 * Supports two methods:
 * 1. PIN/Password flow: Provide accountPassword
 * 2. Wallet signature flow: Provide signature + signMessage
 */
export function useClaimCertificate() {
    const queryClient = useQueryClient();
    const { t } = useTranslation();

    const {
        mutateAsync: claimCertificate,
        isPending: isClaiming,
        error: claimError,
    } = useMutation({
        mutationFn: (params: ClaimCertificateParams) => certificateService.claimCertificate(params),
        onSuccess: (result) => {
            if (result.status === "queued") {
                toast.success(
                    t(
                        "participant.certificates.claimQueued",
                        "Certificate claim submitted — minting in progress!",
                    ),
                );
            } else {
                toast.success(
                    t("participant.certificates.claimSuccess", "Certificate claimed successfully!"),
                );
            }

            // Invalidate related queries to refresh certificate data
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.certificate.all,
            });
            // Invalidate the my certificates list view model query
            // This is used by useCertificateDetailUsecase to determine if certificate is claimed
            queryClient.invalidateQueries({
                queryKey: ["my-certificates-list-viewmodel"],
            });
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.inbox.all,
            });
        },
        onError: (error: unknown) => {
            console.error("Error claiming certificate:", error);

            // Parse error message from backend
            let errorMessage: string;
            if (error && typeof error === "object") {
                if (
                    "error" in error &&
                    error.error &&
                    typeof error.error === "object" &&
                    "message" in error.error
                ) {
                    errorMessage = String(error.error.message);
                } else if ("message" in error) {
                    errorMessage = String(error.message);
                } else {
                    errorMessage = t(
                        "participant.certificates.claimError",
                        "Failed to claim certificate. Please try again.",
                    );
                }
            } else {
                errorMessage = t(
                    "participant.certificates.claimError",
                    "Failed to claim certificate. Please try again.",
                );
            }

            toast.error(errorMessage);
        },
    });

    return {
        claimCertificate,
        isClaiming,
        claimError,
    };
}
