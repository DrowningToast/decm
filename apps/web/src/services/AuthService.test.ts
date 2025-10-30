import { describe, it, expect, beforeEach, vi } from 'vitest';
import { AuthService } from './AuthService';
import { OnboardRegistrationMethod } from '@decm/api';
import type { CoreApiType } from '@/lib/api/api';
import type { OnboardService } from './OnboardService';

describe('AuthService', () => {
    let mockCoreApi: CoreApiType;
    let mockOnboardService: OnboardService;
    let authService: AuthService;

    beforeEach(() => {
        mockCoreApi = {
            v1: {
                registerWithGoogleOauth: vi.fn(),
                registerWithWallet: vi.fn(),
                createProfile: vi.fn(),
                updateProfileByCredentialId: vi.fn(),
                logout: vi.fn(),
            },
        } as unknown as CoreApiType;

        mockOnboardService = {
            checkOnboardStatus: vi.fn(),
        } as unknown as OnboardService;

        authService = new AuthService(mockCoreApi, mockOnboardService);
    });

    describe('createAccount', () => {
        it('should create account with Google OAuth', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);
            (mockCoreApi.v1.registerWithGoogleOauth as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            const result = await authService.createAccount({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: 'test-token',
                password: 'test-password',
                expiresIn: 3600,
            });

            expect(mockOnboardService.checkOnboardStatus).toHaveBeenCalled();
            expect(mockCoreApi.v1.registerWithGoogleOauth).toHaveBeenCalledWith({
                access_token: 'test-token',
                password: 'test-password',
            });
            expect(result.authentication_credential_id).toBe('123');
        });

        it('should create account with Wallet', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);
            (mockCoreApi.v1.registerWithWallet as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '456',
            });

            const result = await authService.createAccount({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: 'test-signature',
            });

            expect(mockOnboardService.checkOnboardStatus).toHaveBeenCalled();
            expect(mockCoreApi.v1.registerWithWallet).toHaveBeenCalledWith({
                signed_message: 'test-signature',
            });
            expect(result.authentication_credential_id).toBe('456');
        });

        it('should throw error if authentication credential already exists', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: 'existing-id',
            });

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'test-token',
                    password: 'test-password',
                    expiresIn: 3600,
                }),
            ).rejects.toThrow();
        });

        it('should throw error for invalid Google OAuth parameters', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: '',
                    password: '',
                    expiresIn: 3600,
                }),
            ).rejects.toThrow();
        });

        it('should throw error for invalid Wallet signature', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodWallet,
                    signSignature: '',
                }),
            ).rejects.toThrow();
        });

        it('should throw error for invalid method', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);

            await expect(
                authService.createAccount({
                    method: 'invalid' as OnboardRegistrationMethod,
                } as any),
            ).rejects.toThrow();
        });
    });

    describe('createProfile', () => {
        it('should create profile successfully', async () => {
            const profileData = {
                firstName: 'John',
                lastName: 'Doe',
                email: 'john.doe@example.com',
            };

            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
                ...profileData,
            });

            const result = await authService.createProfile('123', profileData);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith({
                authentication_credential_id: '123',
                first_name: 'John',
                last_name: 'Doe',
                email: 'john.doe@example.com',
            });
            expect(result).toBeDefined();
        });

        it('should handle all profile fields', async () => {
            const profileData = {
                firstName: 'John',
                lastName: 'Doe',
                email: 'john.doe@example.com',
                academicEmail: 'john@university.edu',
                academicInstitution: 'University',
                address: '123 Main St',
                bio: 'Test bio',
                phoneNumber: '123-456-7890',
                profilePictureUrl: 'https://example.com/pic.jpg',
                isAcademicEmailPublic: true,
                isAcademicInstitutionPublic: true,
                isAddressPublic: false,
                isBioPublic: true,
                isEmailPublic: true,
                isFirstNamePublic: true,
                isLastNamePublic: true,
                isPhoneNumberPublic: false,
                isProfilePicturePublic: true,
            };

            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
                ...profileData,
            });

            await authService.createProfile('123', profileData);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith(
                expect.objectContaining({
                    authentication_credential_id: '123',
                    first_name: 'John',
                    last_name: 'Doe',
                    email: 'john.doe@example.com',
                    academic_email: 'john@university.edu',
                    academic_institution: 'University',
                    address: '123 Main St',
                    bio: 'Test bio',
                    phone_number: '123-456-7890',
                    profile_picture_url: 'https://example.com/pic.jpg',
                    is_academic_email_public: true,
                    is_academic_institution_public: true,
                    is_address_public: false,
                    is_bio_public: true,
                    is_email_public: true,
                    is_first_name_public: true,
                    is_last_name_public: true,
                    is_phone_number_public: false,
                    is_profile_picture_public: true,
                }),
            );
        });
    });

    describe('updateProfile', () => {
        it('should update profile successfully', async () => {
            const profileData = {
                firstName: 'Jane',
                lastName: 'Smith',
            };

            (mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
                ...profileData,
            });

            const result = await authService.updateProfile('123', profileData);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: '123' },
                expect.objectContaining({
                    first_name: 'Jane',
                    last_name: 'Smith',
                }),
            );
            expect(result).toBeDefined();
        });

        it('should not include email in update', async () => {
            const profileData = {
                firstName: 'Jane',
                email: 'should-not-be-included@example.com',
            };

            await authService.updateProfile('123', profileData);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: '123' },
                expect.not.objectContaining({
                    email: expect.anything(),
                }),
            );
        });
    });

    describe('signOut', () => {
        it('should sign out successfully', async () => {
            (mockCoreApi.v1.logout as ReturnType<typeof vi.fn>).mockResolvedValue({});

            await authService.signOut();

            expect(mockCoreApi.v1.logout).toHaveBeenCalled();
        });
    });
});
