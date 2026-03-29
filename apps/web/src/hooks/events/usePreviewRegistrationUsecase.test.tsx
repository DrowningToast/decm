import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { AxiosError } from "axios";
import { usePreviewRegistrationUsecase } from "./usePreviewRegistrationUsecase";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({ t: (key: string) => key }),
}));

vi.mock("sonner", () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    },
}));

const mockSetOnSubmitCallback = vi.fn();
const mockSetOnAcceptCallback = vi.fn();

vi.mock("@/components/BottomNav/stores/event-password", () => ({
    useEventPasswordNavStore: () => ({
        setOnSubmitCallback: mockSetOnSubmitCallback,
    }),
}));

vi.mock("@/components/BottomNav/stores/event-invitation", () => ({
    useEventInvitationNavStore: () => ({
        setOnAcceptCallback: mockSetOnAcceptCallback,
    }),
}));

const mockCheckPassword = vi.fn();

vi.mock("./useCheckPasswordMutation", () => ({
    useCheckPasswordMutation: () => ({
        mutateAsync: mockCheckPassword,
    }),
}));

const mockEventViewModel = vi.fn();

vi.mock("@/components/pages/Participant/Events/Detail/useEventViewModelUsecase", () => ({
    useEventViewModelUsecase: (opts: { eventId: string }) => mockEventViewModel(opts),
}));

import { toast } from "sonner";

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
    });
    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("usePreviewRegistrationUsecase", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockEventViewModel.mockReturnValue({ event: null });
    });

    it("initialises with showPreviewModal false", () => {
        const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
            wrapper: createWrapper(),
        });
        expect(result.current.showPreviewModal).toBe(false);
        expect(result.current.isPasswordValid).toBe(false);
    });

    describe("previewWithPassword", () => {
        it("shows eventFull toast when event is full", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: true, eventType: "private" },
            });

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithPassword("pass");
            });

            expect(toast.error).toHaveBeenCalledWith(
                "event.registration.eventFull",
                expect.any(Object),
            );
            expect(mockCheckPassword).not.toHaveBeenCalled();
        });

        it("shows privateOnly toast when event type is not private", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "invite" },
            });

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithPassword("pass");
            });

            expect(toast.error).toHaveBeenCalledWith(
                "event.registration.privateOnly",
                expect.any(Object),
            );
            expect(mockCheckPassword).not.toHaveBeenCalled();
        });

        it("sets isPasswordValid to true and shows success toast on correct password", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "private" },
            });
            mockCheckPassword.mockResolvedValue(true);

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithPassword("correct");
            });

            expect(result.current.isPasswordValid).toBe(true);
            expect(result.current.showPreviewModal).toBe(true);
            expect(toast.success).toHaveBeenCalledWith(
                "event.registration.success",
                expect.any(Object),
            );
        });

        it("keeps isPasswordValid false and shows error toast on wrong password", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "private" },
            });
            mockCheckPassword.mockResolvedValue(false);

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithPassword("wrong");
            });

            expect(result.current.isPasswordValid).toBe(false);
            expect(toast.error).toHaveBeenCalledWith(
                "event.registration.incorrectPassword",
                expect.any(Object),
            );
        });

        it("shows unavailablePassword toast on 400 error", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "private" },
            });
            const axiosError = new AxiosError("bad request", "400", undefined, undefined, {
                status: 400,
                data: {},
                headers: {},
                config: {} as never,
                statusText: "",
            });
            mockCheckPassword.mockRejectedValue(axiosError);

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithPassword("pass");
            });

            expect(toast.error).toHaveBeenCalledWith(
                "event.registration.unavilablePassword",
                expect.any(Object),
            );
        });
    });

    describe("previewWithInvitation", () => {
        it("shows eventFull toast when event is full", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: true, eventType: "invite", isInvited: true },
            });

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithInvitation();
            });

            expect(toast.error).toHaveBeenCalledWith(
                "event.registration.eventFull",
                expect.any(Object),
            );
            expect(result.current.showPreviewModal).toBe(false);
        });

        it("shows inviteOnly toast when event type is not invite", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "private", isInvited: true },
            });

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithInvitation();
            });

            expect(toast.error).toHaveBeenCalledWith(
                "event.registration.inviteOnly",
                expect.any(Object),
            );
        });

        it("shows inviteOnly toast when user is not invited", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "invite", isInvited: false },
            });

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithInvitation();
            });

            expect(toast.error).toHaveBeenCalledWith(
                "event.registration.inviteOnly",
                expect.any(Object),
            );
        });

        it("sets showPreviewModal true when conditions are met", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "invite", isInvited: true },
            });

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithInvitation();
            });

            expect(result.current.showPreviewModal).toBe(true);
        });
    });

    describe("closePreviewModal", () => {
        it("resets isPasswordValid and closes modal", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "private" },
            });
            mockCheckPassword.mockResolvedValue(true);

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithPassword("correct");
            });

            expect(result.current.showPreviewModal).toBe(true);

            act(() => {
                result.current.closePreviewModal();
            });

            expect(result.current.showPreviewModal).toBe(false);
            expect(result.current.isPasswordValid).toBe(false);
        });

        it("closes invitation modal", async () => {
            mockEventViewModel.mockReturnValue({
                event: { isFull: false, eventType: "invite", isInvited: true },
            });

            const { result } = renderHook(() => usePreviewRegistrationUsecase("event-1"), {
                wrapper: createWrapper(),
            });

            await act(async () => {
                await result.current.previewWithInvitation();
            });

            expect(result.current.showPreviewModal).toBe(true);

            act(() => {
                result.current.closePreviewModal();
            });

            expect(result.current.showPreviewModal).toBe(false);
        });
    });

    it("registers callbacks with stores on mount", async () => {
        mockEventViewModel.mockReturnValue({ event: null });

        renderHook(() => usePreviewRegistrationUsecase("event-1"), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(mockSetOnSubmitCallback).toHaveBeenCalledWith(expect.any(Function));
            expect(mockSetOnAcceptCallback).toHaveBeenCalledWith(expect.any(Function));
        });
    });
});
