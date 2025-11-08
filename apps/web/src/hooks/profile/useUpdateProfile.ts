import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import type { EntityProfile } from "@decm/api";

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
    const { t } = useTranslation();
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: async (profile: ProfileUpdateProfileRequest) => {
            // Get the current profile from cache to extract credential_id
            const currentProfile = queryClient.getQueryData<EntityProfile>(QUERY_KEY.user.profile);
            const credentialId = currentProfile?.authentication_credential_id;

            if (!credentialId) {
                throw new Error("No credential ID found. Please ensure you're logged in.");
            }

            const response = await coreApiClient.v1.updateProfileByCredentialId(
                { credentialId },
                profile,
            );
            return response;
        },
        onSuccess: () => {
            // Invalidate and refetch profile
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
            toast.success(t("profile.updateSuccess"));
        },
        onError: (error) => {
            console.error("Profile update error:", error);
            toast.error(t("profile.updateError"));
        },
    });

    return {
        updateProfile: mutation.mutateAsync,
        isLoading: mutation.isPending,
        error: mutation.error,
    };
};
