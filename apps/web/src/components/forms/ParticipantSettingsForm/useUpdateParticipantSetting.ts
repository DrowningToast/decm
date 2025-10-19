import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import type { EventconfigUpdateEventRegistrationConfigRequest } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

export function useUpdateParticipantSetting(eventId: string) {
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
                queryClient.invalidateQueries({ queryKey: ["event"] });
            },
        });

    async function updateParticipantSetting(
        participantSetting: EventconfigUpdateEventRegistrationConfigRequest,
    ) {
        try {
            await _updateParticipantSetting(participantSetting);
            toast.success("Participant setting updated successfully");
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
