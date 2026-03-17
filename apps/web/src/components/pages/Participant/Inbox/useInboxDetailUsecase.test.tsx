import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useInboxDetailUsecase } from "./useInboxDetailUsecase";
import { EntityInboxMessageType } from "@decm/api";
import type { InboxMessageDetail } from "@/services/InboxService/InboxService";

// Mock dependencies
const mockMarkAsRead = vi.fn().mockResolvedValue(undefined);

vi.mock("@/hooks/inbox/useInboxMessage", () => ({
    useInboxMessage: vi.fn(() => ({
        data: null,
        isLoading: true,
        error: null,
        refetch: vi.fn(),
    })),
}));

vi.mock("@/hooks/inbox/useMarkInboxMessageAsRead", () => ({
    useMarkInboxMessageAsRead: vi.fn(() => ({
        mutateAsync: mockMarkAsRead,
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, defaultValue: string) => defaultValue,
    }),
}));

import { useInboxMessage } from "@/hooks/inbox/useInboxMessage";

describe("useInboxDetailUsecase", () => {
    let queryClient: QueryClient;

    beforeEach(() => {
        queryClient = new QueryClient({
            defaultOptions: { queries: { retry: false } },
        });
        vi.clearAllMocks();
    });

    const wrapper = ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );

    it("returns loading state initially", () => {
        const { result } = renderHook(() => useInboxDetailUsecase({ inboxId: "msg-1" }), {
            wrapper,
        });

        expect(result.current.isLoading).toBe(true);
        expect(result.current.inboxDetail).toBeNull();
    });

    it("returns mapped inbox detail after data loads", async () => {
        const mockMessage: InboxMessageDetail = {
            id: "msg-1",
            entityInboxMessageType: EntityInboxMessageType.InboxMessageTypeGeneral,
            messageType: "general",
            messageContent: "Hello world",
            isRead: true,
            createdAt: new Date("2024-01-01T00:00:00Z"),
            updatedAt: new Date("2024-01-02T00:00:00Z"),
            senderCredentialEmail: "sender@test.com",
        };

        vi.mocked(useInboxMessage).mockReturnValue({
            data: mockMessage,
            isLoading: false,
            error: null,
            refetch: vi.fn(),
        });

        const { result } = renderHook(() => useInboxDetailUsecase({ inboxId: "msg-1" }), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxDetail).not.toBeNull();
        expect(result.current.inboxDetail!.id).toBe("msg-1");
        expect(result.current.inboxDetail!.contentType).toBe("general");
        expect(result.current.inboxDetail!.sender).toBe("sender@test.com");
        expect(result.current.inboxDetail!.title).toBe("Message");
    });

    it("marks message as read on mount when unread", async () => {
        const mockMessage: InboxMessageDetail = {
            id: "msg-1",
            entityInboxMessageType: EntityInboxMessageType.InboxMessageTypeGeneral,
            messageType: "general",
            messageContent: "Hello",
            isRead: false,
            createdAt: new Date("2024-01-01T00:00:00Z"),
            updatedAt: new Date("2024-01-02T00:00:00Z"),
        };

        vi.mocked(useInboxMessage).mockReturnValue({
            data: mockMessage,
            isLoading: false,
            error: null,
            refetch: vi.fn(),
        });

        renderHook(() => useInboxDetailUsecase({ inboxId: "msg-1" }), { wrapper });

        await waitFor(() => {
            expect(mockMarkAsRead).toHaveBeenCalledWith("msg-1");
        });
    });

    it("does not mark message as read when already read", async () => {
        const mockMessage: InboxMessageDetail = {
            id: "msg-1",
            entityInboxMessageType: EntityInboxMessageType.InboxMessageTypeGeneral,
            messageType: "general",
            messageContent: "Hello",
            isRead: true,
            createdAt: new Date("2024-01-01T00:00:00Z"),
            updatedAt: new Date("2024-01-02T00:00:00Z"),
        };

        vi.mocked(useInboxMessage).mockReturnValue({
            data: mockMessage,
            isLoading: false,
            error: null,
            refetch: vi.fn(),
        });

        renderHook(() => useInboxDetailUsecase({ inboxId: "msg-1" }), { wrapper });

        // Give effects a chance to run
        await waitFor(() => {
            expect(mockMarkAsRead).not.toHaveBeenCalled();
        });
    });

    it("maps registration invitation type correctly", async () => {
        const mockMessage: InboxMessageDetail = {
            id: "msg-2",
            entityInboxMessageType:
                EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
            messageType: "event_registration_invitation",
            messageContent: "You are invited",
            isRead: true,
            createdAt: new Date("2024-01-01T00:00:00Z"),
            updatedAt: new Date("2024-01-02T00:00:00Z"),
            senderCredentialEmail: "host@test.com",
            eventId: "evt-1",
            acceptedAt: new Date("2024-01-15T00:00:00Z"),
        };

        vi.mocked(useInboxMessage).mockReturnValue({
            data: mockMessage,
            isLoading: false,
            error: null,
            refetch: vi.fn(),
        });

        const { result } = renderHook(() => useInboxDetailUsecase({ inboxId: "msg-2" }), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxDetail!.contentType).toBe("event-invitation");
        expect(result.current.inboxDetail!.title).toBe("Event Registration Invitation");
        expect(result.current.inboxDetail!.eventId).toBe("evt-1");
    });

    it("maps certificate invitation type correctly", async () => {
        const mockMessage: InboxMessageDetail = {
            id: "msg-3",
            entityInboxMessageType:
                EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
            messageType: "event_certificate_invitation",
            messageContent: "You have a certificate",
            isRead: true,
            createdAt: new Date("2024-01-01T00:00:00Z"),
            updatedAt: new Date("2024-01-02T00:00:00Z"),
            senderCredentialEmail: "host@test.com",
            eventName: "Test Event",
            certificateId: "cert-1",
            certificateTitle: "Completion",
            tokenId: "tok-1",
            hasParticipantJoinedEvent: true,
        };

        vi.mocked(useInboxMessage).mockReturnValue({
            data: mockMessage,
            isLoading: false,
            error: null,
            refetch: vi.fn(),
        });

        const { result } = renderHook(() => useInboxDetailUsecase({ inboxId: "msg-3" }), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxDetail!.contentType).toBe("certificate");
        expect(result.current.inboxDetail!.title).toBe("Certificate Invitation");
        expect(result.current.inboxDetail!.certificateId).toBe("cert-1");
        expect(result.current.inboxDetail!.isUserInEvent).toBe(true);
    });

    it("handles error state", async () => {
        const mockError = new Error("Network error");
        vi.mocked(useInboxMessage).mockReturnValue({
            data: null,
            isLoading: false,
            error: mockError,
            refetch: vi.fn(),
        });

        const { result } = renderHook(() => useInboxDetailUsecase({ inboxId: "msg-1" }), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.error).toBe(mockError);
        expect(result.current.inboxDetail).toBeNull();
    });

    it("handles null message data", async () => {
        vi.mocked(useInboxMessage).mockReturnValue({
            data: null,
            isLoading: false,
            error: null,
            refetch: vi.fn(),
        });

        const { result } = renderHook(() => useInboxDetailUsecase({ inboxId: "msg-1" }), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxDetail).toBeNull();
    });
});
