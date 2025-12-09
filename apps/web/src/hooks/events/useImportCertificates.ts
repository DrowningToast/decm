import { certificateService } from "@/services/services";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useNavigate } from "@/router";
import type { EventImportCertificateReceiverRequest } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

interface UseImportCertificatesOptions {
    onError?: (error: Error) => void;
}

export function useImportCertificates(eventId: string, options?: UseImportCertificatesOptions) {
    const navigate = useNavigate();
    const { t } = useTranslation();

    const { mutateAsync: importCertificates, isPending: isImportingCertificates } = useMutation({
        mutationKey: ["import-certificates"],
        mutationFn: (data: {
            hostPin: string;
            receivers: EventImportCertificateReceiverRequest[];
        }) =>
            certificateService.importCertificates({
                eventId,
                hostPin: data.hostPin,
                receivers: data.receivers,
            }),
        onSuccess: () => {
            toast.success(t("certificateImport.importSuccess"));
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
            navigate("/host/events/:eventId", { params: { eventId } });
        },
        onError: (error: Error) => {
            toast.error(t("certificateImport.importError"), {
                description: error.message,
            });
            console.error(error);
            // Call the onError callback if provided (e.g., to set error state)
            options?.onError?.(error);
        },
    });

    return {
        importCertificates,
        isImportingCertificates,
    };
}
