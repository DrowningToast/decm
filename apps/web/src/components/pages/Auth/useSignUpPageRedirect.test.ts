import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { AxiosError } from "axios";

// ---------------------------------------------------------------------------
// Mocks — must be declared before imports of the module under test
// ---------------------------------------------------------------------------

const mockNavigate = vi.fn();

vi.mock("@/router", () => ({
    useNavigate: () => mockNavigate,
}));

const mockSignout = vi.fn();
vi.mock("@/components/useSignout", () => ({
    useSignout: () => ({ signout: mockSignout }),
}));

// Default mock state — overridden per test via mockCheckOnboardStatus.mockReturnValue
const mockCheckOnboardStatus = vi.fn();
vi.mock("../Onboard/useCheckOnboardStatus", () => ({
    useCheckOnboardStatus: () => mockCheckOnboardStatus(),
}));

// Default wagmi account state — overridden per test via mockUseAccount.mockReturnValue
const mockUseAccount = vi.fn();
vi.mock("wagmi", () => ({
    useAccount: () => mockUseAccount(),
}));

// Default localStorage state — overridden per test
const mockUseLocalStorage = vi.fn();
vi.mock("@/hooks/use-local-storage", () => ({
    useLocalStorage: () => mockUseLocalStorage(),
}));

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function makeAxiosError(status: number): AxiosError {
    const err = new AxiosError("error", undefined, undefined, undefined, {
        status,
        data: {},
        statusText: "",
        headers: {},
        // biome-ignore lint: needed for AxiosError shape
        config: {} as never,
    });
    return err;
}

function setupDefaults({
    address = undefined as string | undefined,
    isConnecting = false,
    isReconnecting = false,
    onboardStatus = undefined as Record<string, string | null | undefined> | undefined,
    isLoading = false,
    error = null as Error | null,
    authSignSignature = undefined as string | undefined,
    accessToken = undefined as string | undefined,
} = {}) {
    mockUseAccount.mockReturnValue({ address, isConnecting, isReconnecting });
    mockCheckOnboardStatus.mockReturnValue({ onboardStatus, isLoading, error });
    // useLocalStorage is called twice: first for AUTH_SIGN_SIGNATURE, second for ACCESS_TOKEN
    mockUseLocalStorage
        .mockReturnValueOnce([authSignSignature, vi.fn()])
        .mockReturnValueOnce([accessToken, vi.fn()]);
}

// ---------------------------------------------------------------------------
// Import after mocks
// ---------------------------------------------------------------------------

import { useSignUpPageRedirect } from "./useSignUpPageRedirect";

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("useSignUpPageRedirect – authCheckWallet", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("does nothing while loading", async () => {
        setupDefaults({ isLoading: true, address: "0xABC" });

        renderHook(() => useSignUpPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).not.toHaveBeenCalled();
            expect(mockSignout).not.toHaveBeenCalled();
        });
    });

    it("does nothing when wallet is not connected and there is no error", async () => {
        setupDefaults({ address: undefined });

        renderHook(() => useSignUpPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).not.toHaveBeenCalled();
            expect(mockSignout).not.toHaveBeenCalled();
        });
    });

    it("navigates to /onboard/wallet when wallet connects and user has no account (400 error from status check)", async () => {
        // This is the regression case:
        // checkOnboardStatus({}) returns 400 for unauthenticated users, setting hasNon401Error=true.
        // Wallet navigation must still fire despite the 400 error.
        setupDefaults({
            address: "0xABC123",
            error: makeAxiosError(400),
        });

        renderHook(() => useSignUpPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/onboard/:method", {
                params: { method: "wallet" },
            });
        });
        expect(mockSignout).not.toHaveBeenCalled();
    });

    it("calls signout (not navigate) when session is expired (401 error)", async () => {
        setupDefaults({
            address: "0xABC123",
            error: makeAxiosError(401),
        });

        renderHook(() => useSignUpPageRedirect());

        await waitFor(() => {
            expect(mockSignout).toHaveBeenCalledWith({ showSuccessToast: false });
        });
        expect(mockNavigate).not.toHaveBeenCalled();
    });

    it("navigates to /app when user is fully registered (has credential and profile)", async () => {
        setupDefaults({
            address: "0xABC123",
            onboardStatus: {
                authentication_credential_id: "cred-1",
                profile_id: "prof-1",
            },
        });

        renderHook(() => useSignUpPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/app");
        });
    });

    it("navigates to /onboard/wallet when user has credential but no profile", async () => {
        setupDefaults({
            address: "0xABC123",
            onboardStatus: {
                authentication_credential_id: "cred-1",
                profile_id: undefined,
            },
        });

        renderHook(() => useSignUpPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/onboard/:method", {
                params: { method: "wallet" },
            });
        });
    });

    it("does nothing when wallet is not connected even with a non-401 error", async () => {
        setupDefaults({
            address: undefined,
            error: makeAxiosError(400),
        });

        renderHook(() => useSignUpPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).not.toHaveBeenCalled();
            expect(mockSignout).not.toHaveBeenCalled();
        });
    });
});
