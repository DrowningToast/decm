import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { AxiosError } from "axios";

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

const mockNavigate = vi.fn();
vi.mock("react-router-dom", () => ({
    useNavigate: () => mockNavigate,
}));

const mockSignout = vi.fn();
vi.mock("@/components/useSignout", () => ({
    useSignout: () => ({ signout: mockSignout }),
}));

const mockCheckOnboardStatus = vi.fn();
vi.mock("../Onboard/useCheckOnboardStatus", () => ({
    useCheckOnboardStatus: () => mockCheckOnboardStatus(),
}));

const mockUseAccount = vi.fn();
vi.mock("wagmi", () => ({
    useAccount: () => mockUseAccount(),
}));

const mockUseLocalStorage = vi.fn();
vi.mock("@/hooks/use-local-storage", () => ({
    useLocalStorage: () => mockUseLocalStorage(),
}));

vi.mock("@/pages/onboard/[method]", () => ({
    OnboardMethods: { WALLET: "wallet", GOOGLE: "google" },
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
    expiresIn = undefined as number | undefined,
} = {}) {
    mockUseAccount.mockReturnValue({ address, isConnecting, isReconnecting });
    mockCheckOnboardStatus.mockReturnValue({ onboardStatus, isLoading, error });
    // useLocalStorage is called three times: AUTH_SIGN_SIGNATURE, ACCESS_TOKEN, EXPIRES_IN
    mockUseLocalStorage
        .mockReturnValueOnce([authSignSignature, vi.fn()])
        .mockReturnValueOnce([accessToken, vi.fn()])
        .mockReturnValueOnce([expiresIn, vi.fn()]);
}

// ---------------------------------------------------------------------------
// Import after mocks
// ---------------------------------------------------------------------------

import { useSignInPageRedirect } from "./useSignInPageRedirect";

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("useSignInPageRedirect – authCheckWallet", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("does nothing while loading", async () => {
        setupDefaults({ isLoading: true, address: "0xABC" });

        renderHook(() => useSignInPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).not.toHaveBeenCalled();
            expect(mockSignout).not.toHaveBeenCalled();
        });
    });

    it("does nothing when wallet is not connected", async () => {
        setupDefaults({ address: undefined });

        renderHook(() => useSignInPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).not.toHaveBeenCalled();
            expect(mockSignout).not.toHaveBeenCalled();
        });
    });

    it("navigates to /signin/sign-message when wallet connects and user has no signature (400 error from status check)", async () => {
        // Regression case: checkOnboardStatus({}) returns 400 for unauthenticated users.
        // Wallet navigation to /signin/sign-message must still fire despite hasNon401Error=true.
        setupDefaults({
            address: "0xABC123",
            authSignSignature: undefined,
            error: makeAxiosError(400),
        });

        renderHook(() => useSignInPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/signin/sign-message");
        });
        expect(mockSignout).not.toHaveBeenCalled();
    });

    it("calls signout when session is expired (401 error)", async () => {
        setupDefaults({
            address: "0xABC123",
            error: makeAxiosError(401),
        });

        renderHook(() => useSignInPageRedirect());

        await waitFor(() => {
            expect(mockSignout).toHaveBeenCalledWith({ showSuccessToast: false });
        });
        expect(mockNavigate).not.toHaveBeenCalled();
    });

    it("navigates to /app when user is fully authenticated", async () => {
        setupDefaults({
            address: "0xABC123",
            onboardStatus: {
                authentication_credential_id: "cred-1",
                profile_id: "prof-1",
            },
        });

        renderHook(() => useSignInPageRedirect());

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

        renderHook(() => useSignInPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/onboard/wallet");
        });
    });

    it("does nothing when wallet is not connected even with a non-401 error", async () => {
        setupDefaults({
            address: undefined,
            error: makeAxiosError(400),
        });

        renderHook(() => useSignInPageRedirect());

        await waitFor(() => {
            expect(mockNavigate).not.toHaveBeenCalled();
            expect(mockSignout).not.toHaveBeenCalled();
        });
    });
});
