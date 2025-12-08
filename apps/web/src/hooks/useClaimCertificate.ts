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
        onSuccess: () => {
            toast.success(
                t("participant.certificates.claimSuccess", "Certificate claimed successfully!"),
            );

            // Invalidate related queries to refresh certificate data
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.certificate.all,
            });
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.inbox.all,
            });
        },
        onError: (error: unknown) => {
            console.error("Error claiming certificate:", error);

            // Parse error message from backend
            const errorMessage =
                error?.error?.message ||
                error?.message ||
                t(
                    "participant.certificates.claimError",
                    "Failed to claim certificate. Please try again.",
                );

            toast.error(errorMessage);
        },
    });

    return {
        claimCertificate,
        isClaiming,
        claimError,
    };
}
