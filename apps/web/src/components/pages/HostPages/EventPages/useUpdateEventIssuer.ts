import { coreApiClient } from "@/lib/api/api";
import type { UpdateEventIssuerPayload } from "@decm/api";
import { useMutation } from "@tanstack/react-query";

export function useUpdateEventIssuer(eventID: string) {
    const { mutateAsync: updateEventIssuer, isPending: isUpdatingEventIssuer } = useMutation({
        mutationKey: ["updateEventIssuer"],
        mutationFn: async (data: UpdateEventIssuerPayload) =>
            await coreApiClient.v1.updateEventIssuer(
                {
                    eventId: eventID,
                },
                data,
            ),
    });

    return {
        updateEventIssuer,
        isUpdatingEventIssuer,
    };
}
