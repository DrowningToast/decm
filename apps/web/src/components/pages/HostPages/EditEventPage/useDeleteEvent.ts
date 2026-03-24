import { queryClient } from "@/lib/api/queryClient";
import { useNavigate } from "@/router";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { QUERY_KEY } from "@/lib/queryKeys";
import { eventService } from "@/services/services";
import type { WalletSignatureAuth } from "@/services/EventService/EventService";

export function useDeleteEvent(eventId: string) {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const { mutateAsync: _deleteWithPassword, isPending: isDeletingWithPassword } = useMutation({
        mutationKey: ["deleteEvent", eventId, "password"],
        mutationFn: (hostPassword: string) =>
            eventService.deleteEventWithPassword(eventId, hostPassword),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.hostEvents.all });
        },
    });

    const { mutateAsync: _deleteWithSignature, isPending: isDeletingWithSignature } = useMutation({
        mutationKey: ["deleteEvent", eventId, "signature"],
        mutationFn: (auth: WalletSignatureAuth) =>
            eventService.deleteEventWithSignature(eventId, auth),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.hostEvents.all });
        },
    });

    async function deleteEventWithPassword(hostPassword: string) {
        try {
            await _deleteWithPassword(hostPassword);
            toast.success(t("deleteEvent.success"));
            navigate("/host/events");
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    async function deleteEventWithSignature(auth: WalletSignatureAuth) {
        try {
            await _deleteWithSignature(auth);
            toast.success(t("deleteEvent.success"));
            navigate("/host/events");
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    return {
        deleteEventWithPassword,
        deleteEventWithSignature,
        isDeletingEvent: isDeletingWithPassword || isDeletingWithSignature,
    };
}
