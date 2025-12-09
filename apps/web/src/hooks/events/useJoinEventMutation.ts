import { useMutation, useQueryClient } from "@tanstack/react-query";
import { eventRegistrationService } from "@/services/services";
import { QUERY_KEY } from "@/lib/queryKeys";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/context/AuthContext";

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
    const { user } = useAuth();

    const joinWithAccountPassword = useMutation({
        mutationFn: async (params: JoinEventWithPasswordParams) => {
            return await eventRegistrationService.joinEventWithAccountPassword({
                eventId: params.eventId,
                accountPassword: params.accountPassword,
                eventPassword: undefined,
                registrationData: params.registrationData as {
                    firstName: string | undefined;
                    lastName: string | undefined;
                    bio: string | undefined;
                    email: string | undefined;
                    phoneNumber: string | undefined;
                    address: string | undefined;
                    academicEmail: string | undefined;
                    academicInstitution: string | undefined;
                },
            });
        },
        onSuccess: (_, variables) => {
            toast.success(t("events.registration.piiForm.submitSuccess"));
            // Invalidate relevant queries
            // Invalidate viewmodel query (used by event detail page to show participation status)
            // Use the exact query key structure to ensure proper invalidation
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.viewmodel(
                    variables.eventId,
                    user?.authenticationCredentialId,
                ),
            });
            // Also invalidate all viewmodel queries for this event (in case userId is different)
            queryClient.invalidateQueries({
                queryKey: ["event", variables.eventId, "viewmodel"],
            });
            // Invalidate event detail queries
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.byId(variables.eventId),
            });
            // Invalidate registration config
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.registrationConfig(variables.eventId),
            });
            // Invalidate invitations (user might have used an invitation)
            queryClient.invalidateQueries({
                queryKey: ["event", variables.eventId, "invitations"],
            });
            // Invalidate specific user invitation query
            if (user?.walletAddress) {
                queryClient.invalidateQueries({
                    queryKey: QUERY_KEY.event.invitations.ofUserAndEventId(
                        variables.eventId,
                        user.walletAddress,
                    ),
                });
            }
            // Invalidate attendees (count might have changed)
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.attendees.byEventId(variables.eventId),
            });
        },
        onError: (error) => {
            console.error("Failed to join event with account password:", error);
            toast.error(t("events.registration.piiForm.submitError"));
        },
    });

    const joinWithEventPassword = useMutation({
        mutationFn: async (params: JoinEventWithEventPasswordParams) => {
            // Note: The service method accepts both accountPassword and eventPassword
            // For event password only, we pass an empty string for accountPassword
            // This should be adjusted based on actual API requirements
            return await eventRegistrationService.joinEventWithAccountPassword({
                eventId: params.eventId,
                accountPassword: "", // Empty for event password only flow
                eventPassword: params.eventPassword,
                registrationData: params.registrationData as {
                    firstName: string | undefined;
                    lastName: string | undefined;
                    bio: string | undefined;
                    email: string | undefined;
                    phoneNumber: string | undefined;
                    address: string | undefined;
                    academicEmail: string | undefined;
                    academicInstitution: string | undefined;
                },
            });
        },
        onSuccess: (_, variables) => {
            toast.success(t("events.registration.piiForm.submitSuccess"));
            // Invalidate relevant queries
            // Invalidate viewmodel query (used by event detail page to show participation status)
            // Use the exact query key structure to ensure proper invalidation
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.viewmodel(
                    variables.eventId,
                    user?.authenticationCredentialId,
                ),
            });
            // Also invalidate all viewmodel queries for this event (in case userId is different)
            queryClient.invalidateQueries({
                queryKey: ["event", variables.eventId, "viewmodel"],
            });
            // Invalidate event detail queries
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.byId(variables.eventId),
            });
            // Invalidate registration config
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.registrationConfig(variables.eventId),
            });
            // Invalidate invitations (user might have used an invitation)
            queryClient.invalidateQueries({
                queryKey: ["event", variables.eventId, "invitations"],
            });
            // Invalidate specific user invitation query
            if (user?.walletAddress) {
                queryClient.invalidateQueries({
                    queryKey: QUERY_KEY.event.invitations.ofUserAndEventId(
                        variables.eventId,
                        user.walletAddress,
                    ),
                });
            }
            // Invalidate attendees (count might have changed)
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.attendees.byEventId(variables.eventId),
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
