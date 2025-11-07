import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useNavigate } from "@/router";
import type { EventCertificateImportRequest } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

export function useImportCertificates(eventId: string) {
    const navigate = useNavigate();

    const { mutateAsync: importCertificates, isPending: isImportingCertificates } = useMutation({
        mutationKey: ["import-certificates"],
        mutationFn: (data: EventCertificateImportRequest) =>
            coreApiClient.v1.importEventCertificates({ eventId }, data),
        onSuccess: () => {
            toast.success("Certificates imported successfully");
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
            navigate(`/host/events/:eventId`, {
                params: { eventId },
            });
        },
        onError: (error) => {
            toast.error("Failed to import certificates", {
                description: error.message,
            });
            console.error(error);
        },
    });

    return {
        importCertificates,
        isImportingCertificates,
    };
}
