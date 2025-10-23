import { type ProfileCreateProfileRequest } from "@decm/api";
import { useMutation, useMutationState } from "@tanstack/react-query";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { onboardService } from "../../../services/OnboardService";
import { authService, type CreateAccountParams } from "../../../services/AuthService";

export const useSignup = () => {
    const { t } = useTranslation();

    const createAccountMutation = useMutationState({
        filters: {
            mutationKey: ["signup"],
        },
        select: (mutation) => mutation.state.status === "pending",
    });
    const { mutateAsync: createAccount, isPending: _isAccountPending } = useMutation({
        mutationKey: ["signup"],
        mutationFn: async (params: CreateAccountParams) => {
            return await authService.createAccount(params);
        },
    });
    const isAccountPending = useMemo(() => {
        return createAccountMutation.some((v) => v) || _isAccountPending;
    }, [_isAccountPending, createAccountMutation]);

    const upsertProfileMutation = useMutationState({
        filters: {
            mutationKey: ["createProfile"],
        },
        select: (mutation) => mutation.state.status === "pending",
    });
    const { mutateAsync: upsertProfile, isPending: _isUpsertProfilePending } = useMutation({
        mutationKey: ["createProfile"],
        mutationFn: async ({
            profile,
            ...params
        }: CreateAccountParams & {
            profile: ProfileCreateProfileRequest;
        }) => {
            const status = await onboardService.checkOnboardStatus(params);
            if (!status?.authentication_credential_id) {
                throw new Error(t("errors.generic"));
            }

            if (status?.profile_id) {
                return await authService.updateProfile(
                    status.authentication_credential_id,
                    profile,
                );
            }

            return await authService.createProfile(status.authentication_credential_id, profile);
        },
    });
    const isUpsertProfilePending = useMemo(() => {
        return upsertProfileMutation.some((v) => v) || _isUpsertProfilePending;
    }, [_isUpsertProfilePending, upsertProfileMutation]);

    const isLoading = useMemo(() => {
        return isAccountPending || isUpsertProfilePending;
    }, [isAccountPending, isUpsertProfilePending]);

    return {
        createAccount,
        upsertProfile,
        isAccountPending,
        isUpsertProfilePending,
        isLoading,
    };
};
