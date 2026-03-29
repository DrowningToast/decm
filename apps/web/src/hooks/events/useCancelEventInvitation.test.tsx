import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useCancelEventInvitation } from "./useCancelEventInvitation";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({ t: (key: string) => key }),
}));

vi.mock("sonner", () => ({
    toast: { success: vi.fn(), error: vi.fn() },
}));

vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            cancelEventRegistrationInvitation: vi.fn(),
        },
    },
}));

const mockInvalidateQueries = vi.hoisted(() => vi.fn());
vi.mock("@/lib/api/queryClient", () => ({
    queryClient: {
        invalidateQueries: mockInvalidateQueries,
    },
}));

import { coreApiClient } from "@/lib/api/api";
import { toast } from "sonner";

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: { mutations: { retry: false } },
    });
    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("useCancelEventInvitation", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("starts with idle state", () => {
        const { result } = renderHook(() => useCancelEventInvitation(), {
            wrapper: createWrapper(),
        });

        expect(result.current.isCancelling).toBe(false);
        expect(result.current.cancelError).toBeNull();
    });

    it("calls cancelEventRegistrationInvitation with correct params on success", async () => {
        vi.mocked(coreApiClient.v1.cancelEventRegistrationInvitation).mockResolvedValue(undefined);

        const { result } = renderHook(() => useCancelEventInvitation(), {
            wrapper: createWrapper(),
        });

        act(() => {
            result.current.cancelEventInvitation({
                eventInvitationId: "inv-1",
                eventId: "event-1",
            });
        });

        await waitFor(() => expect(result.current.isCancelling).toBe(false));

        expect(coreApiClient.v1.cancelEventRegistrationInvitation).toHaveBeenCalledWith({
            event_registration_invitation_id: "inv-1",
        });
    });

    it("invalidates event invitations query on success", async () => {
        vi.mocked(coreApiClient.v1.cancelEventRegistrationInvitation).mockResolvedValue(undefined);

        const { result } = renderHook(() => useCancelEventInvitation(), {
            wrapper: createWrapper(),
        });

        act(() => {
            result.current.cancelEventInvitation({
                eventInvitationId: "inv-1",
                eventId: "event-1",
            });
        });

        await waitFor(() =>
            expect(mockInvalidateQueries).toHaveBeenCalledWith(
                expect.objectContaining({
                    queryKey: expect.arrayContaining(["event", "event-1", "invitations"]),
                }),
            ),
        );
    });

    it("shows success toast on success", async () => {
        vi.mocked(coreApiClient.v1.cancelEventRegistrationInvitation).mockResolvedValue(undefined);

        const { result } = renderHook(() => useCancelEventInvitation(), {
            wrapper: createWrapper(),
        });

        act(() => {
            result.current.cancelEventInvitation({
                eventInvitationId: "inv-1",
                eventId: "event-1",
            });
        });

        await waitFor(() =>
            expect(toast.success).toHaveBeenCalledWith("event.participants.cancelSuccess"),
        );
    });

    it("exposes error when mutation fails", async () => {
        vi.mocked(coreApiClient.v1.cancelEventRegistrationInvitation).mockRejectedValue(
            new Error("server error"),
        );

        const { result } = renderHook(() => useCancelEventInvitation(), {
            wrapper: createWrapper(),
        });

        act(() => {
            result.current.cancelEventInvitation({
                eventInvitationId: "inv-1",
                eventId: "event-1",
            });
        });

        await waitFor(() => expect(result.current.cancelError).toBeTruthy());
        expect(result.current.isCancelling).toBe(false);
    });
});
