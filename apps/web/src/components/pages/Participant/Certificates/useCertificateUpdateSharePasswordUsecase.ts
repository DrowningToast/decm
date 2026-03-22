import { useCallback } from "react";
import { useMutation } from "@tanstack/react-query";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { certificateService } from "@/services/services";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

export const useCertificateUpdateSharePasswordUsecase = () => {
    const { t } = useTranslation();

    const mutation = useMutation({
        mutationFn: ({ shareId, password }: { shareId: string; password: string }) =>
            certificateService.updateCertificateShare(shareId, { password }),
        onSuccess: async () => {
            await queryClient.invalidateQueries({
                queryKey: QUERY_KEY.certificate.myCertificatesListViewModel,
            });
            toast.success(
                t(
                    "participant.certificates.passwordDialog.updateSuccess",
                    "Password protection updated successfully",
                ),
            );
        },
        onError: (error) => {
            toast.error(
                t("participant.certificates.passwordDialog.updateError", {
                    error: error.message,
                }),
            );
        },
    });

    const updateSharePassword = useCallback(
        async (shareId: string, isPasswordProtected: boolean, password: string): Promise<void> => {
            await mutation.mutateAsync({ shareId, password: isPasswordProtected ? password : "" });
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [mutation.mutateAsync],
    );

    return { updateSharePassword, isPending: mutation.isPending };
};
