import { useCallback } from "react";
import { useMutation } from "@tanstack/react-query";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { certificateService } from "@/services/services";

export const useCertificateUpdateShareVisibilityUsecase = () => {
    const mutation = useMutation({
        mutationFn: ({ shareId, active }: { shareId: string; active: boolean }) =>
            certificateService.updateCertificateShare(shareId, { active }),
        onSuccess: async () => {
            await queryClient.invalidateQueries({
                queryKey: QUERY_KEY.certificate.myCertificatesListViewModel,
            });
        },
    });

    const updateShareVisibility = useCallback(
        async (shareId: string, active: boolean): Promise<void> => {
            await mutation.mutateAsync({ shareId, active });
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [mutation.mutateAsync],
    );

    return { updateShareVisibility };
};
