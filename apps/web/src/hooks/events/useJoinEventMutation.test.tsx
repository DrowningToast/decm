import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import { toast } from "sonner";
import { useJoinEventMutation } from "./useJoinEventMutation";

vi.mock("@/services/services", () => ({
    eventRegistrationService: {
        joinEventWithAccountPassword: vi.fn(),
    },
}));

vi.mock("@/context/AuthContext", () => ({
    useAuth: vi.fn(() => ({
        user: {
            authenticationCredentialId: "cred-123",
            walletAddress: "0xabc",
        },
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

vi.mock("sonner", () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    },
}));

import { eventRegistrationService } from "@/services/services";

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: { retry: false },
            mutations: { retry: false },
        },
    });

    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

const registrationData = {
    firstName: "John",
    lastName: "Doe",
    bio: undefined,
    email: "john@example.com",
    phoneNumber: undefined,
    address: undefined,
    academicEmail: undefined,
    academicInstitution: undefined,
};

describe("useJoinEventMutation", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    describe("joinWithAccountPassword", () => {
        it("should call joinEventWithAccountPassword with correct params", async () => {
            vi.mocked(eventRegistrationService.joinEventWithAccountPassword).mockResolvedValue(
                undefined,
            );

            const { result } = renderHook(() => useJoinEventMutation(), {
                wrapper: createWrapper(),
            });

            result.current.joinWithAccountPassword.mutate({
                eventId: "event-1",
                accountPassword: "my-pin",
                registrationData,
            });

            await waitFor(() =>
                expect(result.current.joinWithAccountPassword.isSuccess).toBe(true),
            );

            expect(eventRegistrationService.joinEventWithAccountPassword).toHaveBeenCalledWith({
                eventId: "event-1",
                accountPassword: "my-pin",
                eventPassword: undefined,
                registrationData,
            });
        });

        it("should show success toast on success", async () => {
            vi.mocked(eventRegistrationService.joinEventWithAccountPassword).mockResolvedValue(
                undefined,
            );

            const { result } = renderHook(() => useJoinEventMutation(), {
                wrapper: createWrapper(),
            });

            result.current.joinWithAccountPassword.mutate({
                eventId: "event-1",
                accountPassword: "my-pin",
                registrationData,
            });

            await waitFor(() =>
                expect(result.current.joinWithAccountPassword.isSuccess).toBe(true),
            );

            expect(toast.success).toHaveBeenCalledWith("events.registration.piiForm.submitSuccess");
        });

        it("should show error toast on failure", async () => {
            vi.mocked(eventRegistrationService.joinEventWithAccountPassword).mockRejectedValue(
                new Error("Wrong password"),
            );

            const { result } = renderHook(() => useJoinEventMutation(), {
                wrapper: createWrapper(),
            });

            result.current.joinWithAccountPassword.mutate({
                eventId: "event-1",
                accountPassword: "wrong",
                registrationData,
            });

            await waitFor(() => expect(result.current.joinWithAccountPassword.isError).toBe(true));

            expect(toast.error).toHaveBeenCalledWith("events.registration.piiForm.submitError");
        });

        it("should start with idle status", () => {
            const { result } = renderHook(() => useJoinEventMutation(), {
                wrapper: createWrapper(),
            });

            expect(result.current.joinWithAccountPassword.isPending).toBe(false);
            expect(result.current.joinWithAccountPassword.isSuccess).toBe(false);
            expect(result.current.joinWithAccountPassword.isError).toBe(false);
        });
    });

    describe("joinWithEventPassword", () => {
        it("should call joinEventWithAccountPassword with event password and empty account password", async () => {
            vi.mocked(eventRegistrationService.joinEventWithAccountPassword).mockResolvedValue(
                undefined,
            );

            const { result } = renderHook(() => useJoinEventMutation(), {
                wrapper: createWrapper(),
            });

            result.current.joinWithEventPassword.mutate({
                eventId: "event-1",
                eventPassword: "event-secret",
                registrationData,
            });

            await waitFor(() => expect(result.current.joinWithEventPassword.isSuccess).toBe(true));

            expect(eventRegistrationService.joinEventWithAccountPassword).toHaveBeenCalledWith({
                eventId: "event-1",
                accountPassword: "",
                eventPassword: "event-secret",
                registrationData,
            });
        });

        it("should show success toast on success", async () => {
            vi.mocked(eventRegistrationService.joinEventWithAccountPassword).mockResolvedValue(
                undefined,
            );

            const { result } = renderHook(() => useJoinEventMutation(), {
                wrapper: createWrapper(),
            });

            result.current.joinWithEventPassword.mutate({
                eventId: "event-1",
                eventPassword: "event-secret",
                registrationData,
            });

            await waitFor(() => expect(result.current.joinWithEventPassword.isSuccess).toBe(true));

            expect(toast.success).toHaveBeenCalledWith("events.registration.piiForm.submitSuccess");
        });

        it("should show error toast on failure", async () => {
            vi.mocked(eventRegistrationService.joinEventWithAccountPassword).mockRejectedValue(
                new Error("Wrong event password"),
            );

            const { result } = renderHook(() => useJoinEventMutation(), {
                wrapper: createWrapper(),
            });

            result.current.joinWithEventPassword.mutate({
                eventId: "event-1",
                eventPassword: "wrong",
                registrationData,
            });

            await waitFor(() => expect(result.current.joinWithEventPassword.isError).toBe(true));

            expect(toast.error).toHaveBeenCalledWith("events.registration.piiForm.submitError");
        });
    });
});
