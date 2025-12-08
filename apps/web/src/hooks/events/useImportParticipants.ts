import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useNavigate } from "@/router";
import type { EventRegistrationImportEventParticipantsRequest } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

export function useImportParticipants(eventId: string) {
    const navigate = useNavigate();

    const { mutateAsync: importParticipants, isPending: isImportingParticipants } = useMutation({
        mutationKey: ["import-participants"],
        mutationFn: (data: EventRegistrationImportEventParticipantsRequest) =>
            coreApiClient.v1.importEventParticipants({ eventId }, data),
        onSuccess: () => {
            toast.success("Participants imported successfully");
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
            navigate(`/host/events/:eventId`, {
                params: { eventId },
            });
        },
        onError: (error) => {
            toast.error("Failed to import participants", {
                description: error.message,
            });
            console.error(error);
        },
    });

    return {
        importParticipants,
        isImportingParticipants,
    };
}
