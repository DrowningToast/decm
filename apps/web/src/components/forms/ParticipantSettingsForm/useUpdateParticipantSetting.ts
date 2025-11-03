import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { useNavigate } from "@/router";
import type { EventconfigUpdateEventRegistrationConfigRequest } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { queryKeys } from "@/lib/queryKeys";

export function useUpdateParticipantSetting(eventId: string) {
    const navigate = useNavigate();

    const { mutateAsync: _updateParticipantSetting, isPending: isUpdatingParticipantSetting } =
        useMutation({
            mutationKey: ["update-participant-setting"],
            mutationFn: async (
                participantSetting: EventconfigUpdateEventRegistrationConfigRequest,
            ) => {
                try {
                    const res = await coreApiClient.v1.updateEventRegistrationConfig(
                        {
                            eventId,
                        },
                        participantSetting,
                    );

                    return res;
                } catch (error) {
                    console.error(error);
                    throw error;
                }
            },
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: queryKeys.event.all });
            },
        });

    async function updateParticipantSetting(
        participantSetting: EventconfigUpdateEventRegistrationConfigRequest,
    ) {
        try {
            await _updateParticipantSetting(participantSetting);
            toast.success("Participant setting updated successfully");
            navigate("/host/events/:eventId", {
                params: {
                    eventId,
                },
            });
        } catch (error) {
            console.error(error);
            toast.error("Failed to update participant setting");
        }
    }

    return {
        updateParticipantSetting,
        isUpdatingParticipantSetting,
    };
}
