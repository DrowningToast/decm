import { renderHook, act, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import { toast } from "sonner";
import { usePasswordPrompt } from "./usePasswordPrompt";

const mockOpen = vi.fn();
const mockSetOnSuccess = vi.fn();
const mockSetOnError = vi.fn();

vi.mock("@/components/providers/SignPasswordModal/store", () => ({
    useSignPasswordModalStore: () => ({
        open: mockOpen,
        setOnSuccess: mockSetOnSuccess,
        setOnError: mockSetOnError,
    }),
}));

vi.mock("sonner", () => ({
    toast: {
        error: vi.fn(),
        success: vi.fn(),
    },
}));

const mockUseAuth = vi.fn();
vi.mock("@/context/AuthContext", () => ({
    useAuth: () => mockUseAuth(),
}));

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            mutations: { retry: false },
        },
    });
    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

const promptParams = {
    eventContractAddress: "0xcontract",
    transactionType: "Registration Confirmation",
    title: "Sign to register",
    description: "Confirm your registration",
    details: "Registering for Event A",
};

describe("usePasswordPrompt", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    describe("SYSTEM_MANAGED user", () => {
        beforeEach(() => {
            mockUseAuth.mockReturnValue({
                user: { solutionStatus: "SYSTEM_MANAGED" },
            });
        });

        it("should open the sign password modal", async () => {
            const { result } = renderHook(() => usePasswordPrompt(), {
                wrapper: createWrapper(),
            });

            act(() => {
                result.current.mutate(promptParams);
            });

            await waitFor(() => expect(mockOpen).toHaveBeenCalled());

            expect(mockOpen).toHaveBeenCalledWith(
                promptParams.title,
                promptParams.description,
                true,
                true,
                {
                    contractAddress: promptParams.eventContractAddress,
                    transactionType: promptParams.transactionType,
                    details: promptParams.details,
                },
                expect.any(Function),
            );
        });

        it("should register onSuccess and onError callbacks", async () => {
            const { result } = renderHook(() => usePasswordPrompt(), {
                wrapper: createWrapper(),
            });

            act(() => {
                result.current.mutate(promptParams);
            });

            await waitFor(() => expect(mockSetOnSuccess).toHaveBeenCalled());

            expect(mockSetOnSuccess).toHaveBeenCalledWith(expect.any(Function));
            expect(mockSetOnError).toHaveBeenCalledWith(expect.any(Function));
        });

        it("should resolve the mutation with the password when onSuccess is called", async () => {
            let capturedOnSuccess: ((params: { value: string }) => void) | undefined;
            mockSetOnSuccess.mockImplementation((cb: (params: { value: string }) => void) => {
                capturedOnSuccess = cb;
            });

            const { result } = renderHook(() => usePasswordPrompt(), {
                wrapper: createWrapper(),
            });

            act(() => {
                result.current.mutate(promptParams);
            });

            await waitFor(() => expect(capturedOnSuccess).toBeDefined());

            await act(async () => {
                capturedOnSuccess?.({ value: "user-entered-pin" });
            });

            await waitFor(() => expect(result.current.data).toBe("user-entered-pin"));
        });

        it("should not show an error toast", async () => {
            const { result } = renderHook(() => usePasswordPrompt(), {
                wrapper: createWrapper(),
            });

            act(() => {
                result.current.mutate(promptParams);
            });

            await waitFor(() => expect(mockOpen).toHaveBeenCalled());

            expect(toast.error).not.toHaveBeenCalled();
        });
    });

    describe("BYOK user", () => {
        beforeEach(() => {
            mockUseAuth.mockReturnValue({
                user: { solutionStatus: "BYOK" },
            });
        });

        it("should show an error toast", async () => {
            const { result } = renderHook(() => usePasswordPrompt(), {
                wrapper: createWrapper(),
            });

            act(() => {
                result.current.mutate(promptParams);
            });

            await waitFor(() => expect(toast.error).toHaveBeenCalled());

            expect(toast.error).toHaveBeenCalledWith("walletDebug.byokNotSupported");
        });

        it("should not open the sign password modal", async () => {
            const { result } = renderHook(() => usePasswordPrompt(), {
                wrapper: createWrapper(),
            });

            act(() => {
                result.current.mutate(promptParams);
            });

            await waitFor(() => expect(toast.error).toHaveBeenCalled());

            expect(mockOpen).not.toHaveBeenCalled();
        });
    });
});
