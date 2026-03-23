import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useEventViewModelUsecase } from "./useEventViewModelUsecase";

vi.mock("@/services/services", () => ({
    eventService: {
        getEventViewModelExtended: vi.fn(),
    },
}));

vi.mock("@/context/AuthContext", () => ({
    useAuth: vi.fn(() => ({
        user: { authenticationCredentialId: "cred-1" },
    })),
}));

import { eventService } from "@/services/services";
import { useAuth } from "@/context/AuthContext";

const baseEvent = {
    eventStatus: "open",
    eventType: "private" as const,
    isJoined: false,
    isInvited: false,
    signatureExpiredAt: undefined,
};

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false } },
    });
    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("useEventViewModelUsecase", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(useAuth).mockReturnValue({
            user: { authenticationCredentialId: "cred-1" },
        } as never);
    });

    it("returns event data when query succeeds", async () => {
        vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue(baseEvent);

        const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
            wrapper: createWrapper(),
        });

        await waitFor(() => expect(result.current.isLoading).toBe(false));

        expect(result.current.event).toEqual(baseEvent);
        expect(result.current.error).toBeNull();
    });

    it("is loading while fetching", () => {
        vi.mocked(eventService.getEventViewModelExtended).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
            wrapper: createWrapper(),
        });

        expect(result.current.isLoading).toBe(true);
    });

    it("does not fetch when eventId is empty", () => {
        const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "" }), {
            wrapper: createWrapper(),
        });

        expect(result.current.isLoading).toBe(false);
        expect(eventService.getEventViewModelExtended).not.toHaveBeenCalled();
    });

    describe("bottomNavVariant", () => {
        it("returns undefined when no event data", async () => {
            vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue(undefined);

            const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.bottomNavVariant).toBeUndefined();
        });

        it("returns undefined when no user", async () => {
            vi.mocked(useAuth).mockReturnValue({ user: null } as never);
            vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue(baseEvent);

            const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.bottomNavVariant).toBeUndefined();
        });

        it("returns undefined when event is closed", async () => {
            vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue({
                ...baseEvent,
                eventStatus: "closed",
            });

            const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.bottomNavVariant).toBeUndefined();
        });

        it("returns participating when joined and no signatureExpiredAt", async () => {
            vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue({
                ...baseEvent,
                isJoined: true,
                signatureExpiredAt: undefined,
            });

            const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.bottomNavVariant).toBe("participating");
        });

        it("returns event-password for private event type", async () => {
            vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue({
                ...baseEvent,
                eventType: "private",
                isJoined: false,
            });

            const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.bottomNavVariant).toBe("event-password");
        });

        it("returns invited when invite type and user is invited", async () => {
            vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue({
                ...baseEvent,
                eventType: "invite",
                isInvited: true,
                isJoined: false,
            });

            const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.bottomNavVariant).toBe("invited");
        });

        it("returns invitation-required when invite type and user is not invited", async () => {
            vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue({
                ...baseEvent,
                eventType: "invite",
                isInvited: false,
                isJoined: false,
            });

            const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.bottomNavVariant).toBe("invitation-required");
        });

        it("returns undefined for public/default event type", async () => {
            vi.mocked(eventService.getEventViewModelExtended).mockResolvedValue({
                ...baseEvent,
                eventType: "public" as never,
                isJoined: false,
            });

            const { result } = renderHook(() => useEventViewModelUsecase({ eventId: "event-1" }), {
                wrapper: createWrapper(),
            });

            await waitFor(() => expect(result.current.isLoading).toBe(false));

            expect(result.current.bottomNavVariant).toBeUndefined();
        });
    });
});
