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

describe('AuthService - Additional Edge Cases', () => {
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

    describe('createAccount - Edge Cases', () => {
        it('should handle empty password for Google OAuth', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'token',
                    password: '',
                    expiresIn: 3600,
                }),
            ).rejects.toThrow();
        });

        it('should handle very short access tokens', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'x',
                    password: 'pass',
                    expiresIn: 3600,
                }),
            ).resolves.toBeDefined();
        });

        it('should handle whitespace-only access token', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: '   ',
                    password: 'pass',
                    expiresIn: 3600,
                }),
            ).rejects.toThrow();
        });

        it('should handle zero expiresIn', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);
            (mockCoreApi.v1.registerWithGoogleOauth as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            const result = await authService.createAccount({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: 'token',
                password: 'password',
                expiresIn: 0,
            });

            expect(result.authentication_credential_id).toBe('123');
        });

        it('should handle negative expiresIn', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);
            (mockCoreApi.v1.registerWithGoogleOauth as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            const result = await authService.createAccount({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: 'token',
                password: 'password',
                expiresIn: -1,
            });

            expect(result.authentication_credential_id).toBe('123');
        });

        it('should handle whitespace-only signature', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodWallet,
                    signSignature: '   ',
                }),
            ).rejects.toThrow();
        });

        it('should handle API errors from registerWithGoogleOauth', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);
            (mockCoreApi.v1.registerWithGoogleOauth as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('API Error'),
            );

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'token',
                    password: 'password',
                    expiresIn: 3600,
                }),
            ).rejects.toThrow('API Error');
        });

        it('should handle API errors from registerWithWallet', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);
            (mockCoreApi.v1.registerWithWallet as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Wallet Error'),
            );

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodWallet,
                    signSignature: 'signature',
                }),
            ).rejects.toThrow('Wallet Error');
        });

        it('should handle network timeout', async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue(null);
            (mockCoreApi.v1.registerWithGoogleOauth as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Network timeout'),
            );

            await expect(
                authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: 'token',
                    password: 'password',
                    expiresIn: 3600,
                }),
            ).rejects.toThrow('Network timeout');
        });
    });

    describe('createProfile - Edge Cases', () => {
        it('should handle minimal profile data', async () => {
            const minimalProfile = {
                firstName: 'John',
                lastName: 'Doe',
                email: 'john@example.com',
            };

            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
                ...minimalProfile,
            });

            const result = await authService.createProfile('123', minimalProfile);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith(
                expect.objectContaining({
                    authentication_credential_id: '123',
                    first_name: 'John',
                    last_name: 'Doe',
                    email: 'john@example.com',
                }),
            );
            expect(result).toBeDefined();
        });

        it('should handle profile with only optional fields', async () => {
            const optionalOnlyProfile = {
                bio: 'Just a bio',
                phoneNumber: '123-456-7890',
            };

            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await authService.createProfile('123', optionalOnlyProfile);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith(
                expect.objectContaining({
                    authentication_credential_id: '123',
                    bio: 'Just a bio',
                    phone_number: '123-456-7890',
                }),
            );
        });

        it('should handle profile with all privacy flags false', async () => {
            const profileWithAllPrivate = {
                firstName: 'John',
                lastName: 'Doe',
                email: 'john@example.com',
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isBioPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isAcademicEmailPublic: false,
                isAcademicInstitutionPublic: false,
                isProfilePicturePublic: false,
            };

            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await authService.createProfile('123', profileWithAllPrivate);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith(
                expect.objectContaining({
                    is_first_name_public: false,
                    is_last_name_public: false,
                    is_email_public: false,
                }),
            );
        });

        it('should handle profile with special characters', async () => {
            const profileWithSpecialChars = {
                firstName: "O'Brien",
                lastName: 'Müller-Schmidt',
                email: 'test+alias@example.com',
                address: '123 Main St., Apt #4',
            };

            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await authService.createProfile('123', profileWithSpecialChars);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith(
                expect.objectContaining({
                    first_name: "O'Brien",
                    last_name: 'Müller-Schmidt',
                    email: 'test+alias@example.com',
                    address: '123 Main St., Apt #4',
                }),
            );
        });

        it('should handle very long profile URLs', async () => {
            const longUrl = 'https://example.com/' + 'a'.repeat(500) + '.jpg';
            const profileWithLongUrl = {
                firstName: 'John',
                profilePictureUrl: longUrl,
            };

            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await authService.createProfile('123', profileWithLongUrl);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith(
                expect.objectContaining({
                    profile_picture_url: longUrl,
                }),
            );
        });

        it('should handle API errors during profile creation', async () => {
            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Profile creation failed'),
            );

            await expect(
                authService.createProfile('123', { firstName: 'John' }),
            ).rejects.toThrow('Profile creation failed');
        });

        it('should handle empty credential ID', async () => {
            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '',
            });

            await authService.createProfile('', { firstName: 'John' });

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith(
                expect.objectContaining({
                    authentication_credential_id: '',
                }),
            );
        });
    });

    describe('updateProfile - Edge Cases', () => {
        it('should not allow email updates', async () => {
            const profileWithEmail = {
                firstName: 'Jane',
                email: 'should-not-update@example.com',
            };

            (mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await authService.updateProfile('123', profileWithEmail);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: '123' },
                expect.not.objectContaining({
                    email: expect.anything(),
                }),
            );
        });

        it('should handle partial profile updates', async () => {
            const partialUpdate = {
                bio: 'New bio',
            };

            (mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await authService.updateProfile('123', partialUpdate);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: '123' },
                expect.objectContaining({
                    bio: 'New bio',
                }),
            );
        });

        it('should handle update with no changes', async () => {
            (mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await authService.updateProfile('123', {});

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalled();
        });

        it('should handle API errors during update', async () => {
            (mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Update failed'),
            );

            await expect(
                authService.updateProfile('123', { firstName: 'Jane' }),
            ).rejects.toThrow('Update failed');
        });

        it('should handle updating all privacy flags', async () => {
            const privacyUpdate = {
                isFirstNamePublic: true,
                isLastNamePublic: true,
                isEmailPublic: true,
                isBioPublic: false,
                isPhoneNumberPublic: false,
            };

            (mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await authService.updateProfile('123', privacyUpdate);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: '123' },
                expect.objectContaining({
                    is_first_name_public: true,
                    is_last_name_public: true,
                    is_email_public: true,
                    is_bio_public: false,
                    is_phone_number_public: false,
                }),
            );
        });

        it('should handle concurrent updates', async () => {
            (mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: '123',
            });

            await Promise.all([
                authService.updateProfile('123', { firstName: 'Update1' }),
                authService.updateProfile('123', { lastName: 'Update2' }),
            ]);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledTimes(2);
        });
    });

    describe('signOut - Edge Cases', () => {
        it('should handle successful sign out', async () => {
            (mockCoreApi.v1.logout as ReturnType<typeof vi.fn>).mockResolvedValue({});

            await authService.signOut();

            expect(mockCoreApi.v1.logout).toHaveBeenCalledTimes(1);
        });

        it('should handle sign out errors', async () => {
            (mockCoreApi.v1.logout as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Logout failed'),
            );

            await expect(authService.signOut()).rejects.toThrow('Logout failed');
        });

        it('should handle multiple sign out calls', async () => {
            (mockCoreApi.v1.logout as ReturnType<typeof vi.fn>).mockResolvedValue({});

            await authService.signOut();
            await authService.signOut();

            expect(mockCoreApi.v1.logout).toHaveBeenCalledTimes(2);
        });

        it('should handle network errors during sign out', async () => {
            (mockCoreApi.v1.logout as ReturnType<typeof vi.fn>).mockRejectedValue(
                new Error('Network error'),
            );

            await expect(authService.signOut()).rejects.toThrow('Network error');
        });
    });
});