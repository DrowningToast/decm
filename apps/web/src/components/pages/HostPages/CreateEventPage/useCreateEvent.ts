import { coreApiClient } from "@/lib/api/api";
import { useNavigate } from "@/router";
import type { CreateEventPayload } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

export function useCreateEvent() {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const { mutateAsync: _createEvent, isPending: isCreatingEvent } = useMutation({
        mutationKey: ["createEvent"],
        mutationFn: async (event: CreateEventPayload) => {
            return await coreApiClient.v1.createEvent(event);
        },
    });

    async function createEvent(event: CreateEventPayload) {
        try {
            const response = await _createEvent(event);
            toast.success(t("createEvent.success"));
            navigate("/host/events");
            return response;
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    return {
        createEvent,
        isCreatingEvent,
    };
}
