import { describe, it, expect, vi, beforeEach } from "vitest";
import { CheckOnboardParams, OnboardService } from "./OnboardService";
import { OnboardRegistrationMethod } from "@decm/api";
import type { CoreApiType } from "@/lib/api/api";

describe("OnboardService", () => {
    let mockCoreApi: CoreApiType;
    let onboardService: OnboardService;

    beforeEach(() => {
        mockCoreApi = {
            v1: {
                checkOnboardStatus: vi.fn(),
                getSignMessage: vi.fn(),
            },
        } as unknown as CoreApiType;

        onboardService = new OnboardService(mockCoreApi);
    });

    describe("checkOnboardStatus", () => {
        it("should check onboard status via JWT cookie when no params provided", async () => {
            const mockResponse = {
                authentication_credential_id: "test-id",
                onboard_completed: true,
            };

            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(
                mockResponse,
            );

            const result = await onboardService.checkOnboardStatus();

            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledWith({});
            expect(result).toEqual(mockResponse);
        });

        it("should check onboard status with Google OAuth method", async () => {
            const mockResponse = {
                authentication_credential_id: "google-id",
                onboard_completed: false,
            };

            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(
                mockResponse,
            );

            const params = {
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: "google-access-token",
                expiresIn: 3600,
            } satisfies CheckOnboardParams;

            const result = await onboardService.checkOnboardStatus(params);

            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledWith({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                access_token: "google-access-token",
                expires_in: 3600,
            });
            expect(result).toEqual(mockResponse);
        });

        it("should check onboard status with Wallet method", async () => {
            const mockResponse = {
                authentication_credential_id: "wallet-id",
                onboard_completed: false,
            };

            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(
                mockResponse,
            );

            const params = {
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: "0x123abc",
            } satisfies CheckOnboardParams;

            const result = await onboardService.checkOnboardStatus(params);

            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledWith({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                message_signature: "0x123abc",
            });
            expect(result).toEqual(mockResponse);
        });

        it("should return null when Google OAuth params are invalid", async () => {
            const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

            const params = {
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: "",
                expiresIn: 0,
            } satisfies CheckOnboardParams;

            const result = await onboardService.checkOnboardStatus(params);

            expect(result).toBeNull();
            expect(consoleSpy).toHaveBeenCalledWith("Invalid access token or expires in");
            consoleSpy.mockRestore();
        });

        it("should return null when Wallet signature is invalid", async () => {
            const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

            const params = {
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: "",
            };

            const result = await onboardService.checkOnboardStatus(params);

            expect(result).toBeNull();
            expect(consoleSpy).toHaveBeenCalledWith("Invalid sign signature");
            consoleSpy.mockRestore();
        });

        it("should throw error for invalid method", async () => {
            const params = {
                method: "invalid-method" as OnboardRegistrationMethod,
            } as never;

            await expect(onboardService.checkOnboardStatus(params)).rejects.toThrow(
                "Invalid method",
            );
        });
    });

    describe("getSignMessage", () => {
        it("should return sign message successfully", async () => {
            const mockMessage = "Please sign this message to authenticate";

            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockResolvedValue({
                message: mockMessage,
            });

            const result = await onboardService.getSignMessage();

            expect(mockCoreApi.v1.getSignMessage).toHaveBeenCalled();
            expect(result).toBe(mockMessage);
        });

        it("should handle API errors and throw with descriptive message", async () => {
            const error = new Error("Network error");

            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockRejectedValue(error);

            const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

            await expect(onboardService.getSignMessage()).rejects.toThrow(
                "Failed to retrieve signing message: Network error",
            );

            expect(consoleSpy).toHaveBeenCalledWith(
                "Failed to fetch sign message: Network error",
                error,
            );

            consoleSpy.mockRestore();
        });

        it("should handle unknown errors", async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockRejectedValue(
                "unknown error",
            );

            const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

            await expect(onboardService.getSignMessage()).rejects.toThrow(
                "Failed to retrieve signing message: Unknown error",
            );

            consoleSpy.mockRestore();
        });
    });
});
