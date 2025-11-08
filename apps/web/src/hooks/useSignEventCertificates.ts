import { useMutation, useQueryClient } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

interface SignEventCertificatesParams {
    eventId: string;
    issuerPin: string;
}

interface SignEventCertificatesResponse {
    certificates: Array<{
        certificate: any;
        signature: string;
    }>;
}

export function useSignEventCertificates() {
    const queryClient = useQueryClient();
    const { t } = useTranslation();

    const {
        mutate: signEventCertificates,
        isPending: isSigning,
        error: signingError,
    } = useMutation({
        mutationFn: async ({ eventId, issuerPin }: SignEventCertificatesParams) => {
            const response = await coreApiClient.v1.signEventCertificates(
                { eventId },
                { issuer_pin: issuerPin },
            );
            return response;
        },
        onSuccess: (data: SignEventCertificatesResponse) => {
            toast.success(t("issuer.sign.signingSuccess", { count: data.certificates.length }));
            // Invalidate related queries to refresh data
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.certificates(eventId),
            });
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.issuers.byEventId(eventId),
            });
        },
        onError: (error) => {
            console.error("Error signing event certificates:", error);
            toast.error(t("issuer.sign.signingError"));
        },
    });

    return {
        signEventCertificates,
        isSigning,
        signingError,
    };
}
