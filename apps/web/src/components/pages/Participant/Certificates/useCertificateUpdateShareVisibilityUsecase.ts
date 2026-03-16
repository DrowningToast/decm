import { useCallback } from "react";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { certificateService } from "@/services/services";

export const useCertificateUpdateShareVisibilityUsecase = () => {
    const { t } = useTranslation();
    const mutation = useMutation({
        mutationFn: ({ shareId, active }: { shareId: string; active: boolean }) =>
            certificateService.updateCertificateShare(shareId, { active }),
        onSuccess: async () => {
            await queryClient.invalidateQueries({
                queryKey: QUERY_KEY.certificate.myCertificatesListViewModel,
            });
            toast.success(t("participant.certificates.detail.saveSuccess", "Changes saved"));
        },
        onError: (error) => {
            toast.error(
                t(
                    "participant.certificates.detail.saveError",
                    "Failed to save changes: {{error}}",
                    {
                        error: error instanceof Error ? error.message : String(error),
                    },
                ),
            );
        },
    });

    const updateShareVisibility = useCallback(
        (shareId: string, active: boolean) => {
            mutation.mutate({ shareId, active });
        },
        // mutation.mutate is stable across renders (React Query guarantee)
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [mutation.mutate],
    );

    return { updateShareVisibility };
};
