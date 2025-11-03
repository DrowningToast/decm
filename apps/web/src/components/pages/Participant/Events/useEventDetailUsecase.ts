import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type { Event } from "./useEventsListUsecase";
import { AxiosError } from "axios";
import { useTranslation } from "react-i18next";
import { handleAxiosError } from "@/common/Err";

export interface EventDetail extends Event {
    hasJoined?: boolean;
    invitationStatus?: "not-invited" | "invited" | "accepted";
    correctPassword?: string; // Mock password for testing
}

// Mock event detail data
const mockEventDetails: Record<string, EventDetail> = {
    "1": {
        id: "1",
        name: "ToBelT69 - Password Event",
        description:
            "This is a password-protected event. You need to enter the correct password to join.",
        eventName: "ToBelT69",
        dateTime: "2024-09-24",
        finalCallDate: "2024-09-24",
        status: "accepting",
        accessType: "password",
        requiresPassword: true,
        correctPassword: "decm2024",
        hasJoined: false,
        seatsAvailable: 50,
        totalSeats: 100,
    },
    "2": {
        id: "2",
        name: "ToBelT69 - Invitation Only",
        description:
            "This is an invitation-only event. You can only join if you have been invited.",
        eventName: "ToBelT69",
        dateTime: "2024-09-25",
        finalCallDate: "2024-09-25",
        status: "accepting",
        accessType: "invite-only",
        invitationStatus: "invited",
        seatsAvailable: 30,
        totalSeats: 50,
    },
    "3": {
        id: "3",
        name: "ToBelT69 - Closed Event",
        description: "This event is now closed and no longer accepting participants.",
        eventName: "ToBelT69",
        dateTime: "2024-09-20",
        finalCallDate: "2024-09-20",
        status: "closed",
        accessType: "public",
    },
};

interface UseEventDetailUsecaseOptions {
    eventId: string;
}

export const useEventDetailUsecase = ({ eventId }: UseEventDetailUsecaseOptions) => {
    const queryClient = useQueryClient();
    const { t } = useTranslation();

    // Fetch event detail
    const {
        data: event,
        isLoading,
        error,
    } = useQuery({
        queryKey: ["event-detail", eventId],
        queryFn: async () => {
            try {
                // TODO: Replace with actual API call
                // const response = await api.getEventDetail(eventId);
                const eventDetail = mockEventDetails[eventId];
                if (!eventDetail) {
                    throw new Error("Event not found");
                }
                return eventDetail;
            } catch (error) {
                console.error("Failed to fetch event detail:", error);
                throw error;
            }
        },
    });

    // Submit password mutation
    const submitPasswordMutation = useMutation({
        mutationFn: async ({ password }: { password: string }) => {
            // TODO: Replace with actual API call
            // const response = await api.joinPasswordProtectedEvent(eventId, password);

            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 500));

            if (!event) {
                throw new Error("Event not found");
            }

            // Mock password validation
            if (password === event.correctPassword) {
                return { success: true, message: "Password correct! You have joined the event." };
            } else {
                throw new Error("Incorrect password");
            }
        },
        onSuccess: (data) => {
            toast.success(data.message, {
                description: "You can now participate in this event.",
            });

            // Update the event detail to mark as joined
            queryClient.setQueryData(
                ["event-detail", eventId],
                (oldData: EventDetail | undefined) => {
                    if (!oldData) return oldData;
                    return {
                        ...oldData,
                        hasJoined: true,
                    };
                },
            );
        },
        onError: (error: Error) => {
            if (error instanceof AxiosError) {
                handleAxiosError(t, error);
            }
        },
    });

    // Accept invitation mutation
    const acceptInvitationMutation = useMutation({
        mutationFn: async () => {
            // TODO: Replace with actual API call
            // const response = await api.acceptEventInvitation(eventId);

            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 500));

            if (!event) {
                throw new Error("Event not found");
            }

            if (event.invitationStatus !== "invited") {
                throw new Error("You are not invited to this event");
            }

            return { success: true, message: "Invitation accepted!" };
        },
        onSuccess: (data) => {
            toast.success(data.message, {
                description: "You have successfully joined the event.",
            });

            // Update the event detail to mark as accepted
            queryClient.setQueryData(
                ["event-detail", eventId],
                (oldData: EventDetail | undefined) => {
                    if (!oldData) return oldData;
                    return {
                        ...oldData,
                        invitationStatus: "accepted",
                    };
                },
            );
        },
        onError: (error: Error) => {
            toast.error("Failed to accept invitation", {
                description: error.message,
            });
        },
    });

    // Computed states
    const isPasswordRequired = event?.accessType === "password" || event?.requiresPassword;
    const isInviteOnly = event?.accessType === "invite-only";
    const isClosed = event?.status === "closed";

    // Password event states
    const hasJoinedPasswordEvent = isPasswordRequired && event?.hasJoined;
    const needsPasswordInput = isPasswordRequired && !event?.hasJoined;

    // Invite-only event states
    const isNotInvited = isInviteOnly && event?.invitationStatus === "not-invited";
    const isInvited = isInviteOnly && event?.invitationStatus === "invited";
    const hasAcceptedInvitation = isInviteOnly && event?.invitationStatus === "accepted";

    // Determine bottom nav variant
    const getBottomNavVariant = ():
        | "event-password"
        | "invitation-required"
        | "invited"
        | "participating"
        | undefined => {
        if (isClosed) return undefined;

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
        submitPassword: submitPasswordMutation.mutate,
        isSubmittingPassword: submitPasswordMutation.isPending,
        acceptInvitation: acceptInvitationMutation.mutate,
        isAcceptingInvitation: acceptInvitationMutation.isPending,
        // States
        isPasswordRequired,
        isInviteOnly,
        isClosed,
        hasJoinedPasswordEvent,
        needsPasswordInput,
        isNotInvited,
        isInvited,
        hasAcceptedInvitation,
        bottomNavVariant: getBottomNavVariant(),
    };
};
