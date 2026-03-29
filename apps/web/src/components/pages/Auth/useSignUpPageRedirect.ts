// Use this hook when the user is in the authentication flow

import { useCallback, useEffect, useMemo } from "react";
import { useCheckOnboardStatus } from "../Onboard/useCheckOnboardStatus";
import { useNavigate } from "@/router";
import { useAccount } from "wagmi";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { match } from "ts-pattern";
import { OnboardMethods } from "@/pages/onboard/[method]";
import { AxiosError } from "axios";
import { useSignout } from "@/components/useSignout";

export const useSignUpPageRedirect = () => {
    const navigate = useNavigate();
    const { signout } = useSignout();
    const {
        onboardStatus,
        isLoading: isOnboardStatusLoading,
        error: onboardStatusError,
    } = useCheckOnboardStatus();
    const { address, isReconnecting, isConnecting } = useAccount();
    const [authSignSignature] = useLocalStorage<string | undefined>(
        LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
        undefined,
    );
    const [accessToken] = useLocalStorage<string | undefined>(
        LOCAL_STORAGE_KEYS.ACCESS_TOKEN,
        undefined,
    );

    const isLoading = useMemo(() => {
        return isOnboardStatusLoading || isReconnecting || isConnecting;
    }, [isOnboardStatusLoading, isReconnecting, isConnecting]);

    // Check if error is specifically a 401 Unauthorized
    const isUnauthorizedError =
        onboardStatusError instanceof AxiosError && onboardStatusError.response?.status === 401;

    // Don't redirect if there's a non-401 API error
    const hasNon401Error = !!onboardStatusError && !isUnauthorizedError;

    const authCheckGoogle = useCallback(async () => {
        match({
            isLoading,
            isUnauthorizedError,
            hasNon401Error,
            hasAuthenticationCredentialId: !!onboardStatus?.authentication_credential_id,
            hasProfileId: !!onboardStatus?.profile_id,
            hasAccessToken: !!accessToken,
        })
            .returnType<void>()
            .with(
                {
                    isLoading: true,
                },
                () => {
                    return;
                },
            )
            // Sign out when 401 Unauthorized error occurs
            .with(
                {
                    isUnauthorizedError: true,
                },
                () => {
                    signout({ showSuccessToast: false });
                },
            )
            // Don't redirect when there's a non-401 API error - user should stay on signup page
            .with(
                {
                    hasNon401Error: true,
                },
                () => {
                    return;
                },
            )
            .with(
                {
                    hasAccessToken: true,
                    hasAuthenticationCredentialId: false,
                    hasProfileId: false,
                },
                () => {
                    navigate("/onboard/:method", {
                        params: {
                            method: OnboardMethods.GOOGLE,
                        },
                    });
                },
            )
            .with(
                {
                    hasAccessToken: true,
                    hasAuthenticationCredentialId: true,
                    hasProfileId: true,
                },
                () => {
                    navigate("/onboard/:method", {
                        params: {
                            method: OnboardMethods.GOOGLE,
                        },
                    });
                },
            )
            .with(
                {
                    hasAuthenticationCredentialId: true,
                    hasProfileId: true,
                },
                () => {
                    navigate("/app");
                },
            );
    }, [
        accessToken,
        hasNon401Error,
        isLoading,
        isUnauthorizedError,
        navigate,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
        signout,
    ]);

    const authCheckWallet = useCallback(async () => {
        match({
            isLoading,
            isUnauthorizedError,
            hasNon401Error,
            hasAddress: !!address,
            hasAuthenticationCredentialId: !!onboardStatus?.authentication_credential_id,
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
                },
            )
            // Sign out when 401 Unauthorized error occurs (expired/invalid session)
            .with(
                {
                    isUnauthorizedError: true,
                },
                () => {
                    signout({ showSuccessToast: false });
                },
            )
            // No wallet connected, nothing to do
            .with(
                {
                    hasAddress: false,
                },
                () => {
                    return;
                },
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
                },
            )
            .with(
                {
                    hasAuthenticationCredentialId: true,
                    hasProfileId: true,
                },
                () => {
                    navigate("/app");
                },
            )
            // Wallet connected with no account yet — proceed to onboard
            // (hasNon401Error here is expected: unauthenticated users get 400 from checkOnboardStatus)
            .with(
                {
                    hasAddress: true,
                },
                () => {
                    navigate("/onboard/:method", {
                        params: {
                            method: OnboardMethods.WALLET,
                        },
                    });
                },
            )
            // Don't redirect when there's a non-401 API error - user should stay on signup page
            .with(
                {
                    hasNon401Error: true,
                },
                () => {
                    return;
                },
            );
    }, [
        authSignSignature,
        hasNon401Error,
        isLoading,
        isUnauthorizedError,
        navigate,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
        address,
        signout,
    ]);

    useEffect(() => {
        authCheckGoogle();
        authCheckWallet();
    }, [authCheckGoogle, authCheckWallet]);

    return {
        isLoading,
    };
};
