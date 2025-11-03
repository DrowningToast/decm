import { useQuery } from "@tanstack/react-query";

export type InboxContentType = "event-invitation" | "certificate";

export interface InboxDetail {
    id: string;
    title: string;
    sender: string;
    date: string;
    status: "pending" | "available" | "expired" | "action-required";
    contentType: InboxContentType;
    eventId?: string;
    certificateId?: string;
    description?: string;
    isUserInEvent?: boolean;
}

const MOCK_INBOX_DETAILS: Record<string, InboxDetail> = {
    "1": {
        id: "1",
        title: "Event Invitation",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "pending",
        contentType: "event-invitation",
        eventId: "evt-001",
        description:
            "You have been invited to join the ToBeIT69 event. This is a great opportunity to learn and network with other professionals in the tech community.",
        isUserInEvent: false,
    },
    "2": {
        id: "2",
        title: "Event Invitation",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "available",
        contentType: "event-invitation",
        eventId: "evt-002",
        description:
            "You have accepted the invitation to join the ToBeIT69 event. We look forward to seeing you there!",
        isUserInEvent: true,
    },
    "3": {
        id: "3",
        title: "Event Invitation",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "expired",
        contentType: "event-invitation",
        eventId: "evt-003",
        description: "This event invitation has expired. You can no longer accept this invitation.",
        isUserInEvent: false,
    },
    "4": {
        id: "4",
        title: "New certificate",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "action-required",
        contentType: "certificate",
        certificateId: "cert-001",
        description:
            "Congratulations! You have earned a new certificate from the ToBeIT69 event. Click below to view your certificate.",
        isUserInEvent: true,
    },
    "5": {
        id: "5",
        title: "New certificate",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "available",
        contentType: "certificate",
        certificateId: "cert-002",
        description:
            "Your certificate is ready! You can now view and download it from your certificates page.",
        isUserInEvent: false,
    },
};

interface UseInboxDetailOptions {
    inboxId: string;
}

export const useInboxDetailUsecase = (options: UseInboxDetailOptions) => {
    // const api = useApi();
    const {
        data: inboxDetail = null,
        isLoading,
        error,
    } = useQuery({
        queryKey: ["participant-inbox-detail", options.inboxId],
        queryFn: async () => {
            try {
                // TODO: Replace with actual API call once endpoint is available
                // const response = await api.getInboxDetail(options.inboxId);
                return MOCK_INBOX_DETAILS[options.inboxId] || null;
            } catch (error) {
                console.error("Failed to fetch inbox detail:", error);
                return null;
            }
        },
    });

    return { inboxDetail, isLoading, error };
};
