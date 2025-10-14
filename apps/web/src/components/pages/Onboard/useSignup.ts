import { type ProfileCreateProfileRequest } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { onboardService } from "../../../services/OnboardService";
import {
	authService,
	type CreateAccountParams,
} from "../../../services/AuthService";

export const useSignup = () => {
	const { t } = useTranslation();

	const { mutateAsync: createAccount, isPending: isAccountPending } =
		useMutation({
			mutationKey: ["signup"],
			mutationFn: async (params: CreateAccountParams) => {
				return await authService.createAccount(params);
			},
		});

	const { mutateAsync: upsertProfile, isPending: isUpsertProfilePending } =
		useMutation({
			mutationKey: ["createProfile"],
			mutationFn: async ({
				authenticationCredentialId,
				profile,
				...params
			}: CreateAccountParams & {
				authenticationCredentialId: string;
				profile: ProfileCreateProfileRequest;
			}) => {
				const status = await onboardService.checkOnboardStatus(params);
				if (status.profile_id) {
					if (!status.authentication_credential_id) {
						throw new Error(t("errors.generic"));
					}
					return await authService.updateProfile(
						status.authentication_credential_id,
						profile
					);
				}
				return await authService.createProfile(
					authenticationCredentialId,
					profile
				);
			},
		});

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
