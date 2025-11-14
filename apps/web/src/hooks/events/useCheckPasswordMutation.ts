import { eventRegistrationService } from "@/services/services";
import { useMutation } from "@tanstack/react-query";

export const useCheckPasswordMutation = () => {
    const mutation = useMutation({
        mutationFn: async ({ eventId, password }: { eventId: string; password: string }) => {
            const response = await eventRegistrationService.checkPassword(eventId, password);
            return response;
        },
    });

    return mutation;
};
