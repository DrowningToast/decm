import { useMutation, useQueryClient } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import type { CoreApiInternalHandlerEventSignEventCertificatesResponse } from "@decm/api";

interface SignEventCertificatesParams {
    eventId: string;
    issuerPin: string;
}

export function useSignEventCertificates() {
    const queryClient = useQueryClient();
    const { t } = useTranslation();

    const {
        mutate: signEventCertificates,
        isPending: isSigning,
        error: signingError,
    } = useMutation({
        mutationFn: ({ eventId, issuerPin }: SignEventCertificatesParams) =>
            coreApiClient.v1.signEventCertificates({ eventId }, { issuer_pin: issuerPin }),
        onSuccess: (data: CoreApiInternalHandlerEventSignEventCertificatesResponse) => {
            const certificatesCount = data.certificates?.length || 0;
            toast.success(t("issuer.sign.signingSuccess", { count: certificatesCount }));
            // Invalidate related queries to refresh data
            queryClient.invalidateQueries({
                queryKey: ["issuer", "events"],
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
