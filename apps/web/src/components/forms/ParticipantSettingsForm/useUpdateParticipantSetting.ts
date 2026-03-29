import { queryClient } from "@/lib/api/queryClient";
import { useNavigate } from "@/router";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { QUERY_KEY } from "@/lib/queryKeys";
import { eventRegistrationService } from "@/services/services";
import type { EventRegistrationConfiguration } from "@/services/EventRegistration/EventRegistration";
import { useTranslation } from "react-i18next";

export function useUpdateParticipantSetting(eventId: string) {
    const navigate = useNavigate();
    const { t } = useTranslation();

    const { mutateAsync: _updateParticipantSetting, isPending: isUpdatingParticipantSetting } =
        useMutation({
            mutationKey: ["update-participant-setting"],
            mutationFn: async (configuration: EventRegistrationConfiguration) => {
                try {
                    return await eventRegistrationService.updateConfiguration(
                        eventId,
                        configuration,
                    );
                } catch (error) {
                    console.error(error);
                    throw error;
                }
            },
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
            },
        });

    async function updateParticipantSetting(configuration: EventRegistrationConfiguration) {
        try {
            await _updateParticipantSetting(configuration);
            toast.success(t("participantSettings.updateSuccess"));
            navigate("/host/events/:eventId", {
                params: {
                    eventId,
                },
            });
        } catch (error) {
            console.error(error);
            toast.error(t("participantSettings.updateError"));
        }
    }

    return {
        updateParticipantSetting,
        isUpdatingParticipantSetting,
    };
}
