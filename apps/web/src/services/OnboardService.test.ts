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

describe('OnboardService - Additional Edge Cases', () => {
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

    describe('checkOnboardStatus - Edge Cases', () => {
        it('should handle zero expiresIn', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: 'token',
                expiresIn: 0,
            });

            expect(result).toBeNull();
        });

        it('should handle negative expiresIn', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: 'token',
                expiresIn: -1,
            });

            expect(result).toBeNull();
            consoleSpy.mockRestore();
        });

        it('should handle whitespace-only access token', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: '   ',
                expiresIn: 3600,
            });

            expect(result).toBeNull();
            consoleSpy.mockRestore();
        });

        it('should handle whitespace-only signature', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: '   ',
            });

            expect(result).toBeNull();
            consoleSpy.mockRestore();
        });

        it('should handle API errors gracefully', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('API Error'),
            );

            await expect(
                onboardService.checkOnboardStatus({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'token',
                    expiresIn: 3600,
                }),
            ).rejects.toThrow('API Error');
        });

        it('should handle undefined params correctly', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            const result = await onboardService.checkOnboardStatus();

            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledWith({});
            expect(result?.authentication_credential_id).toBe('123');
        });

        it('should handle network timeout', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Network timeout'),
            );

            await expect(
                onboardService.checkOnboardStatus({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'token',
                    expiresIn: 3600,
                }),
            ).rejects.toThrow('Network timeout');
        });

        it('should handle successful Google OAuth status check', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: 'google-123',
                profile: { email: 'user@example.com' },
            });

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: 'valid-token',
                expiresIn: 3600,
            });

            expect(result?.authentication_credential_id).toBe('google-123');
        });

        it('should handle successful Wallet status check', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: 'wallet-456',
            });

            const result = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: '0x123abc',
            });

            expect(result?.authentication_credential_id).toBe('wallet-456');
        });

        it('should handle JWT cookie authentication', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: 'jwt-789',
            });

            const result = await onboardService.checkOnboardStatus();

            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledWith({});
            expect(result?.authentication_credential_id).toBe('jwt-789');
        });

        it('should handle concurrent status checks', async () => {
            (mockCoreApi.v1.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            const results = await Promise.all([
                onboardService.checkOnboardStatus({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'token1',
                    expiresIn: 3600,
                }),
                onboardService.checkOnboardStatus({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'token2',
                    expiresIn: 3600,
                }),
            ]);

            expect(results).toHaveLength(2);
            expect(mockCoreApi.v1.checkOnboardStatus).toHaveBeenCalledTimes(2);
        });
    });

    describe('getSignMessage - Edge Cases', () => {
        it('should handle successful message retrieval', async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockResolvedValue({
                message: 'Sign this to authenticate',
            });

            const result = await onboardService.getSignMessage();

            expect(result).toBe('Sign this to authenticate');
            expect(mockCoreApi.v1.getSignMessage).toHaveBeenCalledTimes(1);
        });

        it('should handle empty message', async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockResolvedValue({
                message: '',
            });

            const result = await onboardService.getSignMessage();

            expect(result).toBe('');
        });

        it('should handle very long messages', async () => {
            const longMessage = 'A'.repeat(10000);
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockResolvedValue({
                message: longMessage,
            });

            const result = await onboardService.getSignMessage();

            expect(result).toBe(longMessage);
            expect(result.length).toBe(10000);
        });

        it('should handle messages with special characters', async () => {
            const specialMessage = 'Sign: €£¥ 中文 🎉';
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockResolvedValue({
                message: specialMessage,
            });

            const result = await onboardService.getSignMessage();

            expect(result).toBe(specialMessage);
        });

        it('should handle API rate limiting', async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Rate limit exceeded'),
            );

            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            await expect(onboardService.getSignMessage()).rejects.toThrow(
                'Failed to retrieve signing message',
            );

            consoleSpy.mockRestore();
        });

        it('should handle server errors', async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Internal server error'),
            );

            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            await expect(onboardService.getSignMessage()).rejects.toThrow(
                'Failed to retrieve signing message',
            );

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });

        it('should handle multiple concurrent message requests', async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockResolvedValue({
                message: 'Sign this message',
            });

            const results = await Promise.all([
                onboardService.getSignMessage(),
                onboardService.getSignMessage(),
                onboardService.getSignMessage(),
            ]);

            expect(results).toEqual([
                'Sign this message',
                'Sign this message',
                'Sign this message',
            ]);
            expect(mockCoreApi.v1.getSignMessage).toHaveBeenCalledTimes(3);
        });

        it('should preserve error details in thrown error', async () => {
            const detailedError = new Error('Detailed API error');
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockRejectedValue(detailedError);

            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            try {
                await onboardService.getSignMessage();
                fail('Should have thrown an error');
            } catch (error) {
                expect(error).toBeInstanceOf(Error);
                expect((error as Error).message).toContain('Failed to retrieve signing message');
            }

            consoleSpy.mockRestore();
        });

        it('should handle timeout errors', async () => {
            (mockCoreApi.v1.getSignMessage as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Request timeout'),
            );

            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            await expect(onboardService.getSignMessage()).rejects.toThrow(
                'Failed to retrieve signing message: Request timeout',
            );

            consoleSpy.mockRestore();
        });
    });
});