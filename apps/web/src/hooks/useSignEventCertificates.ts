import { useMutation, useQueryClient } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { QUERY_KEY } from "@/lib/queryKeys";

type SignEventCertificatesParams =
    | { eventId: string; issuerPin: string }
    | { eventId: string; issuerSignMessage: string; issuerSignature: string };

export function useSignEventCertificates() {
    const queryClient = useQueryClient();
    const { t } = useTranslation();

    const {
        mutate: signEventCertificates,
        isPending: isSigning,
        error: signingError,
    } = useMutation({
        mutationFn: (params: SignEventCertificatesParams) => {
            const { eventId, ...rest } = params;
            if ("issuerPin" in rest) {
                return certificateService.signEventCertificates(eventId, {
                    issuerPin: rest.issuerPin,
                });
            }
            return certificateService.signEventCertificates(eventId, {
                issuerSignMessage: rest.issuerSignMessage,
                issuerSignature: rest.issuerSignature,
            });
        },
        onSuccess: (data) => {
            const certificatesCount = data.totalSigned;
            toast.success(t("issuer.sign.signingSuccess", { count: certificatesCount }));
            // Invalidate related queries to refresh data
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.all,
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
