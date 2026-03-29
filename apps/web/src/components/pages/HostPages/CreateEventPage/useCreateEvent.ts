import { useNavigate } from "@/router";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { eventService } from "@/services/services";
import type {
    CreateEventWithPasswordInput,
    CreateEventWithSignatureInput,
} from "@/services/EventService/EventService";

export function useCreateEvent() {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const { mutateAsync: _createEventWithPassword, isPending: isCreatingWithPassword } =
        useMutation({
            mutationKey: ["createEvent", "password"],
            mutationFn: (input: CreateEventWithPasswordInput) =>
                eventService.createEventWithPassword(input),
        });

    const { mutateAsync: _createEventWithSignature, isPending: isCreatingWithSignature } =
        useMutation({
            mutationKey: ["createEvent", "signature"],
            mutationFn: (input: CreateEventWithSignatureInput) =>
                eventService.createEventWithSignature(input),
        });

    async function createEventWithPassword(input: CreateEventWithPasswordInput) {
        try {
            const eventId = await _createEventWithPassword(input);
            toast.success(t("createEvent.success"));
            navigate("/host/events/:eventId/settings/participant", {
                params: { eventId },
            });
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    async function createEventWithSignature(input: CreateEventWithSignatureInput) {
        try {
            const eventId = await _createEventWithSignature(input);
            toast.success(t("createEvent.success"));
            navigate("/host/events/:eventId/settings/participant", {
                params: { eventId },
            });
        } catch (error) {
            console.error(error);
            toast.error(t("errors.generic"));
            throw error;
        }
    }

    return {
        createEventWithPassword,
        createEventWithSignature,
        isCreatingEvent: isCreatingWithPassword || isCreatingWithSignature,
    };
}
