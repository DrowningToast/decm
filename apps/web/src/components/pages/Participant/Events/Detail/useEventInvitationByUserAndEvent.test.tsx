import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useEventInvitationByUserAndEvent } from "./useEventInvitationByUserAndEvent";

vi.mock("@/services/services", () => ({
    eventRegistrationService: {
        getInvitationOfUserAndEventId: vi.fn(),
    },
}));

import { eventRegistrationService } from "@/services/services";

const mockInvitation = {
    registrationInvitation: { id: "inv-1", eventId: "event-1" },
    inbox: { id: "msg-1" },
};

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false } },
    });
    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("useEventInvitationByUserAndEvent", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("returns invitation data when query succeeds", async () => {
        vi.mocked(eventRegistrationService.getInvitationOfUserAndEventId).mockResolvedValue(
            mockInvitation as never,
        );

        const { result } = renderHook(() => useEventInvitationByUserAndEvent("event-1", "user-1"), {
            wrapper: createWrapper(),
        });

        await waitFor(() => expect(result.current.isLoading).toBe(false));

        expect(result.current.invitation).toEqual(mockInvitation);
        expect(result.current.error).toBeNull();
        expect(eventRegistrationService.getInvitationOfUserAndEventId).toHaveBeenCalledWith(
            "event-1",
        );
    });

    it("does not fetch when userId is undefined", () => {
        const { result } = renderHook(
            () => useEventInvitationByUserAndEvent("event-1", undefined),
            { wrapper: createWrapper() },
        );

        expect(result.current.isLoading).toBe(false);
        expect(eventRegistrationService.getInvitationOfUserAndEventId).not.toHaveBeenCalled();
    });

    it("does not fetch when eventId is empty", () => {
        const { result } = renderHook(() => useEventInvitationByUserAndEvent("", "user-1"), {
            wrapper: createWrapper(),
        });

        expect(result.current.isLoading).toBe(false);
        expect(eventRegistrationService.getInvitationOfUserAndEventId).not.toHaveBeenCalled();
    });

    it("returns error when query fails", async () => {
        vi.mocked(eventRegistrationService.getInvitationOfUserAndEventId).mockRejectedValue(
            new Error("not found"),
        );

        const { result } = renderHook(() => useEventInvitationByUserAndEvent("event-1", "user-1"), {
            wrapper: createWrapper(),
        });

        await waitFor(() => expect(result.current.isLoading).toBe(false));

        expect(result.current.error).toBeTruthy();
        expect(result.current.invitation).toBeUndefined();
    });
});
