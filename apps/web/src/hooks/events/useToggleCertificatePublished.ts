import { useMutation, useQueryClient } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { QUERY_KEY } from "@/lib/queryKeys";

interface ToggleCertificatePublishedParams {
    eventId: string;
    isPublished: boolean;
}

export const useToggleCertificatePublished = () => {
    const queryClient = useQueryClient();
    const { t } = useTranslation();

    return useMutation({
        mutationFn: ({ eventId, isPublished }: ToggleCertificatePublishedParams) =>
            certificateService.toggleCertificatePublished(eventId, isPublished),
        onSuccess: (_data, variables) => {
            // Invalidate and refetch certificate config
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.certificate.config(variables.eventId),
            });

            if (variables.isPublished) {
                toast.success(t("event.certificates.publishConfigSuccess"));
            }
        },
        onError: (error: unknown) => {
            const errorMessage =
                error && typeof error === "object" && "response" in error
                    ? (error.response as { data?: { message?: string } })?.data?.message
                    : undefined;
            toast.error(errorMessage || t("event.certificates.updatePublishedStatusError"));
        },
    });
};
