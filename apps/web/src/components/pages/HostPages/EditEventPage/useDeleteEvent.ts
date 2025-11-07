import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { useNavigate } from "@/router";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { QUERY_KEY } from "@/lib/queryKeys";

export function useDeleteEvent(eventId: string) {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const { mutateAsync: _deleteEvent, isPending: isDeletingEvent } = useMutation({
        mutationKey: ["deleteEvent", eventId],
        mutationFn: async (hostPassword: string) => {
            return await coreApiClient.v1.deleteEventById(
                {
                    eventId,
                },
                {
                    host_password: hostPassword,
                },
            );
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.hostEvents.all });
        },
    });

    async function deleteEvent(hostPassword: string) {
        try {
            await _deleteEvent(hostPassword);
            toast.success(t("deleteEvent.success"));

            navigate("/host/events");
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    return {
        deleteEvent,
        isDeletingEvent,
    };
}
