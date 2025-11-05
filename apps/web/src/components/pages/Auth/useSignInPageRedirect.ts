// Use this hook when the user is in the authentication flow

import { useCallback, useEffect } from "react";
import { useCheckOnboardStatus } from "../Onboard/useCheckOnboardStatus";
import { useWalletClient } from "wagmi";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { match } from "ts-pattern";
import { useNavigate } from "react-router-dom";
import { OnboardMethods } from "@/pages/onboard/[method]";

export const useSignInPageRedirect = () => {
    const navigate = useNavigate();
    const { onboardStatus, isLoading } = useCheckOnboardStatus();
    const { data: walletClient } = useWalletClient();
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

    const authCheckGoogle = useCallback(async () => {
        match({
            isLoading,
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
        isLoading,
        navigate,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
    ]);

    const authCheckWallet = useCallback(async () => {
        match({
            isLoading,
            hasWalletClient: !!walletClient,
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
            .with(
                {
                    hasWalletClient: false,
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
                    hasWalletClient: true,
                    hasAuthSignSignature: false,
                },
                () => {
                    navigate("/signin/sign-message");
                },
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
