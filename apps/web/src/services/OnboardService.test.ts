import { describe, it, expect, beforeEach, vi } from 'vitest';
import { OnboardService } from './OnboardService';
import { OnboardRegistrationMethod } from '@decm/api';
import type { CoreApiType } from '@/lib/api/api';

describe('OnboardService', () => {
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

    describe('checkOnboardStatus', () => {
        it('should check status with Google OAuth method', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: 'test-token',
                expiresIn: 3600,
            });

            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledWith({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                access_token: 'test-token',
                expires_in: 3600,
            });
            expect(result?.authentication_credential_id).toBe('123');
        });

        it('should check status with Wallet method', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '456',
            });

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: 'test-signature',
            });

            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledWith({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                message_signature: 'test-signature',
            });
            expect(result?.authentication_credential_id).toBe('456');
        });

        it('should check status without params (uses JWT cookie)', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '789',
            });

            const result = await onboardService.checkOnboardStatus();

            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledWith({});
            expect(result?.authentication_credential_id).toBe('789');
        });

        it('should return null for invalid Google OAuth params', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: '',
                expiresIn: 0,
            });

            expect(result).toBeNull();
            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });

        it('should return null for invalid Wallet signature', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: '',
            });

            expect(result).toBeNull();
            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });

        it('should throw error for invalid method', async () => {
            await expect(
                onboardService.checkOnboardStatus({
                    method: 'invalid' as OnboardRegistrationMethod,
                } as any),
            ).rejects.toThrow('Invalid method');
        });
    });

    describe('getSignMessage', () => {
        it('should get sign message successfully', async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockResolvedValue({
                message: 'Sign this message',
            });

            const result = await onboardService.getSignMessage();

            expect(mockCoreApi.v1.getSignMessage).toHaveBeenCalled();
            expect(result).toBe('Sign this message');
        });

        it('should handle errors and throw with error message', async () => {
            const error = new Error('API Error');
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockRejectedValue(error);
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            await expect(onboardService.getSignMessage()).rejects.toThrow('Failed to retrieve signing message');

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });

        it('should handle non-Error rejection', async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockRejectedValue('String error');
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            await expect(onboardService.getSignMessage()).rejects.toThrow('Failed to retrieve signing message');

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });
});
