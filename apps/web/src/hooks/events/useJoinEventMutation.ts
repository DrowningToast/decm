import { useMutation, useQueryClient } from "@tanstack/react-query";
import { eventRegistrationService } from "@/services/services";
import { QUERY_KEY } from "@/lib/queryKeys";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

interface JoinEventWithPasswordParams {
    eventId: string;
    accountPassword: string;
    registrationData: {
        firstName?: string;
        lastName?: string;
        bio?: string;
        email?: string;
        phoneNumber?: string;
        address?: string;
        academicEmail?: string;
        academicInstitution?: string;
    };
}

interface JoinEventWithEventPasswordParams {
    eventId: string;
    eventPassword: string;
    registrationData: {
        firstName?: string;
        lastName?: string;
        bio?: string;
        email?: string;
        phoneNumber?: string;
        address?: string;
        academicEmail?: string;
        academicInstitution?: string;
    };
}

/**
 * Hook for joining events with account password or event password
 * Wraps eventRegistrationService calls in React Query mutations
 */
export const useJoinEventMutation = () => {
    const { t } = useTranslation();
    const queryClient = useQueryClient();

    const joinWithAccountPassword = useMutation({
        mutationFn: async (params: JoinEventWithPasswordParams) => {
            return await eventRegistrationService.joinEventWithAccountPassword({
                eventId: params.eventId,
                accountPassword: params.accountPassword,
                registrationData: params.registrationData,
            });
        },
        onSuccess: (_, variables) => {
            toast.success(t("events.registration.piiForm.submitSuccess"));
            // Invalidate relevant queries
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.detail(variables.eventId),
            });
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.eventRegistration.config(variables.eventId),
            });
        },
        onError: (error) => {
            console.error("Failed to join event with account password:", error);
            toast.error(t("events.registration.piiForm.submitError"));
        },
    });

    const joinWithEventPassword = useMutation({
        mutationFn: async (params: JoinEventWithEventPasswordParams) => {
            return await eventRegistrationService.joinEventWithEventPassword({
                eventId: params.eventId,
                eventPassword: params.eventPassword,
                registrationData: params.registrationData,
            });
        },
        onSuccess: (_, variables) => {
            toast.success(t("events.registration.piiForm.submitSuccess"));
            // Invalidate relevant queries
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.detail(variables.eventId),
            });
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.eventRegistration.config(variables.eventId),
            });
        },
        onError: (error) => {
            console.error("Failed to join event with event password:", error);
            toast.error(t("events.registration.piiForm.submitError"));
        },
    });

    return {
        joinWithAccountPassword,
        joinWithEventPassword,
    };
};
