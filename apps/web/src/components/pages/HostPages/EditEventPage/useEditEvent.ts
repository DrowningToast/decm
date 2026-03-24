import { queryClient } from "@/lib/api/queryClient";
import { useNavigate } from "@/router";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { QUERY_KEY } from "@/lib/queryKeys";
import { eventService } from "@/services/services";
import type {
    EditEventWithPasswordInput,
    EditEventWithSignatureInput,
} from "@/services/EventService/EventService";

export function useEditEvent(eventId: string) {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const { mutateAsync: _editWithPassword, isPending: isEditingWithPassword } = useMutation({
        mutationKey: ["editEvent", eventId, "password"],
        mutationFn: (input: EditEventWithPasswordInput) =>
            eventService.editEventWithPassword(eventId, input),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.byId(eventId) });
        },
    });

    const { mutateAsync: _editWithSignature, isPending: isEditingWithSignature } = useMutation({
        mutationKey: ["editEvent", eventId, "signature"],
        mutationFn: (input: EditEventWithSignatureInput) =>
            eventService.editEventWithSignature(eventId, input),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.byId(eventId) });
        },
    });

    async function editEventWithPassword(input: EditEventWithPasswordInput) {
        try {
            await _editWithPassword(input);
            toast.success(t("editEvent.success"));
            navigate("/host/events/:eventId", { params: { eventId } });
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    async function editEventWithSignature(input: EditEventWithSignatureInput) {
        try {
            await _editWithSignature(input);
            toast.success(t("editEvent.success"));
            navigate("/host/events/:eventId", { params: { eventId } });
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    return {
        editEventWithPassword,
        editEventWithSignature,
        isEditingEvent: isEditingWithPassword || isEditingWithSignature,
    };
}
