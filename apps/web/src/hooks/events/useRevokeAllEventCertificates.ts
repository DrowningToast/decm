import { certificateService } from "@/services/services";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

export function useRevokeAllEventCertificates() {
    const { t } = useTranslation();

    const {
        mutate: revokeAllEventCertificates,
        isPending: isRevokingAll,
        error: revokeAllError,
    } = useMutation({
        mutationFn: ({ eventId }: { eventId: string }) =>
            certificateService.revokeAllCertificates(eventId),
        onSuccess: () => {
            toast.success(t("event.certificates.revokeAllSuccess"));
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
        },
        onError: (error) => {
            toast.error(t("event.certificates.revokeAllError"));
            console.error(error);
        },
    });

    return {
        revokeAllEventCertificates,
        isRevokingAll,
        revokeAllError,
    };
}
