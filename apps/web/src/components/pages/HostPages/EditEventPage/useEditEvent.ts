import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { useNavigate } from "@/router";
import type { UpdateEventPayload } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

export function useEditEvent(eventId: string) {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const { mutateAsync: _editEvent, isPending: isEditingEvent } = useMutation({
        mutationKey: ["editEvent", eventId],
        mutationFn: async (event: UpdateEventPayload) => {
            return await coreApiClient.v1.updateEvent(
                {
                    eventId,
                },
                event,
            );
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["event"] });
        },
    });

    async function editEvent(event: UpdateEventPayload) {
        try {
            const response = await _editEvent(event);
            toast.success(t("editEvent.success"));

            navigate("/host/events/:eventId", {
                params: {
                    eventId: response.id ?? "",
                },
            });
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    return {
        editEvent,
        isEditingEvent,
    };
}
