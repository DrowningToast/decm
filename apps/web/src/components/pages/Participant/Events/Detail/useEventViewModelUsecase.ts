import { useQuery } from "@tanstack/react-query";
import { eventService } from "@/services/services";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useAuth } from "@/context/AuthContext";
import { useMemo } from "react";
import type { BottomNavVariant } from "@/components/BottomNav/BottomNav";

interface UseEventDetailUsecaseOptions {
    eventId: string;
}

export const useEventViewModelUsecase = ({ eventId }: UseEventDetailUsecaseOptions) => {
    const { user } = useAuth();
    const {
        data: event,
        isLoading,
        error,
    } = useQuery({
        queryKey: QUERY_KEY.event.viewmodel(eventId, user?.authenticationCredentialId),
        queryFn: async () => {
            try {
                const response = await eventService.getEventViewModelExtended(eventId);
                return response;
            } catch (error) {
                console.error("Failed to fetch event detail:", error);
                throw error;
            }
        },
        enabled: !!eventId,
    });

    // Determine bottom nav variant
    const bottomNavVariant = useMemo<BottomNavVariant | undefined>(() => {
        if (!event || !user) {
            return undefined;
        }
        if (event.eventStatus === "closed") {
            return undefined;
        }
        if (event.isJoined && !event.signatureExpiredAt) {
            return "participating";
        }
        switch (event.eventType) {
            case "invite": {
                if (event.isInvited) {
                    return "invited";
                }
                return "invitation-required";
            }
            case "private": {
                return "event-password";
            }
            default: {
                return undefined;
            }
        }
    }, [event, user]);

    return {
        event,
        isLoading,
        error,
        // States
        bottomNavVariant,
    };
};
