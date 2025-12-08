import { QUERY_KEY } from "@/lib/queryKeys";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { authService } from "@/services/services";

// Mock type for ProfileUpdateProfileRequest until API is regenerated
export interface ProfileUpdateProfileRequest {
    academic_email?: string;
    academic_institution?: string;
    address?: string;
    bio?: string;
    email?: string;
    first_name?: string;
    is_academic_email_public?: boolean;
    is_academic_institution_public?: boolean;
    is_address_public?: boolean;
    is_bio_public?: boolean;
    is_email_public?: boolean;
    is_first_name_public?: boolean;
    is_last_name_public?: boolean;
    is_phone_number_public?: boolean;
    is_profile_picture_public?: boolean;
    last_name?: string;
    phone_number?: string;
    profile_picture_url?: string;
}

export const useUpdateProfile = () => {
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: async (profile: ProfileUpdateProfileRequest) => {
            // Get the current profile to extract credential_id
            const currentProfile = await authService.getMyProfile();

            if (!currentProfile?.authenticationCredentialId) {
                throw new Error("No credential ID found. Please ensure you're logged in.");
            }

            // Transform snake_case to camelCase for service layer
            return authService.updateProfile(currentProfile.authenticationCredentialId, {
                academicEmail: profile.academic_email,
                academicInstitution: profile.academic_institution,
                address: profile.address,
                bio: profile.bio,
                email: profile.email,
                firstName: profile.first_name,
                lastName: profile.last_name,
                phoneNumber: profile.phone_number,
                profilePictureUrl: profile.profile_picture_url,
                isAcademicEmailPublic: profile.is_academic_email_public,
                isAcademicInstitutionPublic: profile.is_academic_institution_public,
                isAddressPublic: profile.is_address_public,
                isBioPublic: profile.is_bio_public,
                isEmailPublic: profile.is_email_public,
                isFirstNamePublic: profile.is_first_name_public,
                isLastNamePublic: profile.is_last_name_public,
                isPhoneNumberPublic: profile.is_phone_number_public,
                isProfilePicturePublic: profile.is_profile_picture_public,
            });
        },
        onSuccess: () => {
            // Invalidate and refetch profile
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
        },
        onError: (error) => {
            console.error("Profile update error:", error);
        },
    });

    return {
        updateProfile: mutation.mutateAsync,
        isLoading: mutation.isPending,
        error: mutation.error,
    };
};
