import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useMyInvitation } from "./useMyInvitation";

vi.mock("@/services/services", () => ({
    eventRegistrationService: {
        getInvitationOfUserAndEventId: vi.fn(),
    },
}));

vi.mock("@/context/AuthContext", () => ({
    useAuth: vi.fn(() => ({
        user: { authenticationCredentialId: "cred-1" },
    })),
}));

import { eventRegistrationService } from "@/services/services";
import { useAuth } from "@/context/AuthContext";

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

describe("useMyInvitation", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(useAuth).mockReturnValue({
            user: { authenticationCredentialId: "cred-1" },
        } as never);
    });

    it("returns invitation when user and eventId are present", async () => {
        vi.mocked(eventRegistrationService.getInvitationOfUserAndEventId).mockResolvedValue(
            mockInvitation as never,
        );

        const { result } = renderHook(() => useMyInvitation("event-1"), {
            wrapper: createWrapper(),
        });

        await waitFor(() => expect(result.current.isLoading).toBe(false));

        expect(result.current.data).toEqual(mockInvitation);
        expect(result.current.isSuccess).toBe(true);
        expect(eventRegistrationService.getInvitationOfUserAndEventId).toHaveBeenCalledWith(
            "event-1",
        );
    });

    it("does not fetch when user is null", () => {
        vi.mocked(useAuth).mockReturnValue({ user: null } as never);

        const { result } = renderHook(() => useMyInvitation("event-1"), {
            wrapper: createWrapper(),
        });

        expect(result.current.isLoading).toBe(false);
        expect(eventRegistrationService.getInvitationOfUserAndEventId).not.toHaveBeenCalled();
    });

    it("does not fetch when eventId is empty", () => {
        const { result } = renderHook(() => useMyInvitation(""), {
            wrapper: createWrapper(),
        });

        expect(result.current.isLoading).toBe(false);
        expect(eventRegistrationService.getInvitationOfUserAndEventId).not.toHaveBeenCalled();
    });

    it("returns error when query fails", async () => {
        vi.mocked(eventRegistrationService.getInvitationOfUserAndEventId).mockRejectedValue(
            new Error("unauthorized"),
        );

        const { result } = renderHook(() => useMyInvitation("event-1"), {
            wrapper: createWrapper(),
        });

        await waitFor(() => expect(result.current.isLoading).toBe(false));

        expect(result.current.isError).toBe(true);
        expect(result.current.data).toBeUndefined();
    });
});
