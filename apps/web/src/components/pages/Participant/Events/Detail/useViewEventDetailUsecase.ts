import { useQuery } from "@tanstack/react-query";
import { eventService } from "@/services/services";

interface UseEventDetailUsecaseOptions {
    eventId: string;
}

export const useViewEventDetailUsecase = ({ eventId }: UseEventDetailUsecaseOptions) => {
    // Fetch event detail
    const {
        data: event,
        isLoading,
        error,
    } = useQuery({
        queryKey: ["event-detail", eventId],
        queryFn: async () => {
            try {
                const response = await eventService.getEventById(eventId);
                return response;
            } catch (error) {
                console.error("Failed to fetch event detail:", error);
                throw error;
            }
        },
    });
    // Determine bottom nav variant
    const getBottomNavVariant = ():
        | "event-password"
        | "invitation-required"
        | "invited"
        | "participating"
        | undefined => {
        if (isClosed) {
            return undefined;
        }

        // Password-required event
        if (isPasswordRequired) {
            if (hasJoinedPasswordEvent) {
                return "participating";
            }
            return "event-password";
        }

        // Invite-only event
        if (isInviteOnly) {
            if (hasAcceptedInvitation) {
                return "participating";
            }
            if (isInvited) {
                return "invited";
            }
            return "invitation-required";
        }

        return undefined;
    };

    return {
        event,
        isLoading,
        error,
        // States
        bottomNavVariant: getBottomNavVariant(),
    };
};
