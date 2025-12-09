import { certificateService } from "@/services/services";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

export function useRevokeEventCertificate() {
    const { t } = useTranslation();

    const {
        mutate: revokeEventCertificate,
        isPending: isRevoking,
        error: revokeError,
    } = useMutation({
        mutationFn: ({ certificateIds, eventId }: { certificateIds: string[]; eventId: string }) =>
            certificateService.revokeCertificates({ eventId, certificateIds }),
        onSuccess: () => {
            toast.success(t("event.certificates.revokeSuccess"));
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
        },
        onError: (error) => {
            toast.error(t("event.certificates.revokeError"));
            console.error(error);
        },
    });

    return {
        revokeEventCertificate,
        isRevoking,
        revokeError,
    };
}
