// Use this hook when the user is in the authentication flow

import { useCallback, useEffect, useMemo } from "react";
import { useCheckOnboardStatus } from "../Onboard/useCheckOnboardStatus";

import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { match } from "ts-pattern";
import { useNavigate } from "react-router-dom";
import { OnboardMethods } from "@/pages/onboard/[method]";
import { useAccount } from "wagmi";
import { AxiosError } from "axios";
import { useSignout } from "@/components/useSignout";

export const useSignInPageRedirect = () => {
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
    const [expiresIn] = useLocalStorage<number | undefined>(
        LOCAL_STORAGE_KEYS.EXPIRES_IN,
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
            // Don't redirect based on onboardStatus when there's a non-401 API error
            // But still allow accessToken-based redirects (OAuth flow)
            .with(
                {
                    hasNon401Error: true,
                    hasAccessToken: false,
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
                    navigate("/auth/success");
                },
            )
            .with(
                {
                    hasAccessToken: true,
                    hasAuthenticationCredentialId: true,
                    hasProfileId: true,
                },
                () => {
                    const queryParams = new URLSearchParams();
                    if (!accessToken || !expiresIn) {
                        return;
                    }
                    queryParams.set("access_token", accessToken);
                    queryParams.set("expires_in", expiresIn.toString());
                    const searchString = queryParams.toString();
                    navigate(`/auth/success?${searchString}`);
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
        expiresIn,
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
            // Sign out when 401 Unauthorized error occurs
            .with(
                {
                    isUnauthorizedError: true,
                },
                () => {
                    signout({ showSuccessToast: false });
                },
            )
            // Don't redirect based on onboardStatus when there's a non-401 API error
            // User should stay on signin page when not authenticated
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
                    navigate(`/onboard/${OnboardMethods.WALLET}`);
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
            .with(
                {
                    hasAddress: true,
                    hasAuthSignSignature: false,
                },
                () => {
                    console.log("navigate to sign message");
                    navigate("/signin/sign-message");
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
