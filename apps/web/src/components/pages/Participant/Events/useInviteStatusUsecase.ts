import { useQuery } from "@tanstack/react-query";

export interface InviteStatus {
    isInvited: boolean;
    invitedAt?: string;
    invitedBy?: string;
    status?: "pending" | "accepted" | "rejected";
}

// Mock data for user invite status - in real app, would fetch from API
// This checks if current user is invited to an invite-only event
const mockInviteStatuses: Record<string, InviteStatus> = {
    "1": {
        isInvited: false,
        status: "pending",
    },
    "2": {
        isInvited: true,
        invitedAt: "2024-09-15",
        invitedBy: "admin@example.com",
        status: "pending",
    },
    "3": {
        isInvited: false,
        status: "rejected",
    },
};

/**
 * Hook to check if the current user is invited to an event
 * Used for invite-only events to determine if user can access/join
 */
export const useInviteStatusUsecase = (eventId: string) => {
    const {
        data: inviteStatus,
        isLoading,
        error,
    } = useQuery({
        queryKey: ["event-invite-status", eventId],
        queryFn: async () => {
            try {
                // TODO: Replace with actual API call once endpoint is available
                // const response = await api.checkEventInviteStatus(eventId);
                const status = mockInviteStatuses[eventId];
                if (!status) {
                    return {
                        isInvited: false,
                        status: "pending",
                    } as InviteStatus;
                }
                return status;
            } catch (error) {
                console.error("Failed to fetch invite status:", error);
                throw error;
            }
        },
    });

    return { inviteStatus, isLoading, error };
};
