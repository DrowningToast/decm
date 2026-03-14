import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { certificateService } from "@/services/services";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

export const useCertificateCreateShareLinkUsecase = () => {
    const { t } = useTranslation();
    const mutation = useMutation({
        mutationFn: (certificateId: string, password?: string) =>
            certificateService.createCertificateShareableLink(certificateId, password),
        onSuccess: async () => {
            await queryClient.invalidateQueries({
                queryKey: QUERY_KEY.certificate.myCertificatesListViewModel,
            });
            toast.success(
                t(
                    "participant.certificates.shareSuccess",
                    "Certificate share link created successfully!",
                ),
            );
        },
        onError: (error) => {
            toast.error(t("participant.certificates.shareError", { error: error.message }));
        },
    });

    return {
        createCertificateShareLink: mutation.mutateAsync,
        isPending: mutation.isPending,
        error: mutation.error,
    };
};
