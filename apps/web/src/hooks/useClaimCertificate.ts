import { useMutation, useQueryClient } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { QUERY_KEY } from "@/lib/queryKeys";

interface ClaimCertificateParams {
    certificateId: string;
    eventId: string;
    accountPassword: string;
}

/**
 * Hook for claiming/minting a certificate as a participant
 * This will sign the certificate on the blockchain and mint the NFT
 */
export function useClaimCertificate() {
    const queryClient = useQueryClient();
    const { t } = useTranslation();

    const {
        mutateAsync: claimCertificate,
        isPending: isClaiming,
        error: claimError,
    } = useMutation({
        mutationFn: async ({ certificateId, eventId, accountPassword }: ClaimCertificateParams) => {
            // TODO: Backend API Integration Required
            // =====================================
            //
            // Expected API Endpoint:
            //   POST /api/v1/events/{event_id}/certificates/{certificate_id}/claim
            //   or POST /api/v1/certificates/{certificate_id}/claim
            //
            // Request Body:
            //   {
            //     "account_password": string  // User's verified account password/PIN
            //   }
            //
            // Response Body:
            //   {
            //     "certificate_id": string,
            //     "certificate_token_id": string,  // NFT token ID after minting
            //     "transaction_hash": string,      // Blockchain transaction hash
            //     "claimed_at": string            // ISO timestamp
            //   }
            //
            // Backend Implementation Notes:
            //   1. Verify the account_password matches the user's stored password
            //   2. Check if certificate belongs to the authenticated user
            //   3. Verify certificate hasn't already been claimed (token_id is null)
            //   4. Verify user is a participant in the event
            //   5. Mint certificate NFT on blockchain using system wallet
            //   6. Update certificate record with token_id and claimed_at timestamp
            //   7. Update inbox message status to "claimed"
            //
            // OpenAPI Annotation Example:
            //   // @Summary Claim certificate as participant
            //   // @Description Mint certificate NFT on blockchain and mark as claimed
            //   // @Tags Certificates, Participant
            //   // @Accept json
            //   // @Produce json
            //   // @Param event_id path string true "Event ID"
            //   // @Param certificate_id path string true "Certificate ID"
            //   // @Param body body ClaimCertificateRequest true "Account password"
            //   // @Success 200 {object} ClaimCertificateResponse
            //   // @Failure 400 {object} customerror.ErrResponse "Invalid request"
            //   // @Failure 401 {object} customerror.ErrResponse "Invalid password"
            //   // @Failure 404 {object} customerror.ErrResponse "Certificate not found"
            //   // @Failure 409 {object} customerror.ErrResponse "Certificate already claimed"
            //   // @Router /api/v1/events/{event_id}/certificates/{certificate_id}/claim [post]
            //
            // Once backend is ready, replace the mock below with:
            // const response = await coreApiClient.v1.claimCertificate(
            //     { certificateId, eventId },
            //     { account_password: accountPassword }
            // );
            // return response;

            // Temporary mock response for development
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({
                        certificate_id: certificateId,
                        certificate_token_id: "mock-token-id-" + Date.now(),
                        transaction_hash: "0xmock-" + Math.random().toString(36).substring(2, 15),
                        claimed_at: new Date().toISOString(),
                    });
                }, 1500);
            });
        },
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
        onError: (error) => {
            console.error("Error claiming certificate:", error);
            toast.error(
                t(
                    "participant.certificates.claimError",
                    "Failed to claim certificate. Please try again.",
                ),
            );
        },
    });

    return {
        claimCertificate,
        isClaiming,
        claimError,
    };
}
