// Use this hook when the user is in the authentication flow

import { useCallback, useEffect } from "react";
import { useCheckOnboardStatus } from "../Onboard/useCheckOnboardStatus";
import { useNavigate } from "@/router";
import { useWalletClient } from "wagmi";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { match } from "ts-pattern";

// eg. Sign in, Sign up, onboard page
export const useAuthRedirect = () => {
	const navigate = useNavigate();
	const { onboardStatus, isLoading } = useCheckOnboardStatus();
	const { data: walletClient } = useWalletClient();
	const [authSignSignature] = useLocalStorage<string | undefined>(
		LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
		undefined
	);

	const authCheckGoogle = useCallback(async () => {
		match({
			isLoading,
			hasAuthenticationCredentialId:
				!!onboardStatus?.authentication_credential_id,
			hasProfileId: !!onboardStatus?.profile_id,
		})
			.returnType<void>()
			.with(
				{
					isLoading: true,
				},
				() => {
					return;
				}
			)
			.with(
				{
					hasAuthenticationCredentialId: true,
					hasProfileId: false,
				},
				() => {
					navigate("/onboard/:method", {
						params: {
							method: "wallet",
						},
					});
				}
			)
			.with(
				{
					hasAuthenticationCredentialId: true,
					hasProfileId: true,
				},
				() => {
					navigate("/app");
				}
			);
	}, [isLoading, navigate, onboardStatus]);

	const authCheckWallet = useCallback(async () => {
		console.log({
			isLoading,
			hasWalletClient: !!walletClient,
			hasAuthenticationCredentialId:
				!!onboardStatus?.authentication_credential_id,
			hasProfileId: !!onboardStatus?.profile_id,
			hasAuthSignSignature: !!authSignSignature,
		});
		match({
			isLoading,
			hasWalletClient: !!walletClient,
			hasAuthenticationCredentialId:
				!!onboardStatus?.authentication_credential_id,
			hasProfileId: !!onboardStatus?.profile_id,
			hasAuthSignSignature: !!authSignSignature,
		})
			.returnType<void>()
			.with(
				{
					isLoading: true,
				},
				() => {
					return;
				}
			)
			.with(
				{
					hasWalletClient: false,
				},
				() => {
					return;
				}
			)
			.with(
				{
					hasAuthenticationCredentialId: true,
					hasProfileId: false,
				},
				() => {
					navigate("/onboard/:method", {
						params: {
							method: "wallet",
						},
					});
				}
			)
			.with(
				{
					hasAuthenticationCredentialId: true,
					hasProfileId: true,
				},
				() => {
					navigate("/app");
				}
			)
			.with(
				{
					hasWalletClient: true,
				},
				() => {
					navigate("/onboard/:method", {
						params: {
							method: "wallet",
						},
					});
				}
			);
	}, [
		authSignSignature,
		isLoading,
		navigate,
		onboardStatus?.authentication_credential_id,
		onboardStatus?.profile_id,
		walletClient,
	]);

	useEffect(() => {
		authCheckGoogle();
		authCheckWallet();
	}, [authCheckGoogle, authCheckWallet]);
};
