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

type ImportCertificatesInput =
    | { hostPin: string; receivers: EventImportCertificateReceiverRequest[] }
    | {
          hostSignMessage: string;
          hostSignature: string;
          receivers: EventImportCertificateReceiverRequest[];
      };

export function useImportCertificates(eventId: string, options?: UseImportCertificatesOptions) {
    const navigate = useNavigate();
    const { t } = useTranslation();

    const { mutateAsync: importCertificates, isPending: isImportingCertificates } = useMutation({
        mutationKey: ["import-certificates"],
        mutationFn: (data: ImportCertificatesInput) => {
            if ("hostPin" in data) {
                return certificateService.importCertificates({
                    eventId,
                    hostPin: data.hostPin,
                    receivers: data.receivers,
                });
            }
            return certificateService.importCertificates({
                eventId,
                hostSignMessage: data.hostSignMessage,
                hostSignature: data.hostSignature,
                receivers: data.receivers,
            });
        },
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
            options?.onError?.(error);
        },
    });

    const { mutateAsync: getSignMessage, isPending: isGettingSignMessage } = useMutation({
        mutationKey: ["certificate-import-sign-message", eventId],
        mutationFn: (receivers: EventImportCertificateReceiverRequest[]) =>
            certificateService.getCertificateImportSignMessage({ eventId, receivers }),
    });

    return {
        importCertificates,
        isImportingCertificates,
        getSignMessage,
        isGettingSignMessage,
    };
}
