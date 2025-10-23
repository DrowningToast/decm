// Use this hook when the user is in the authentication flow

import { useCallback, useEffect } from "react";
import { useCheckOnboardStatus } from "../Onboard/useCheckOnboardStatus";
import { useNavigate } from "@/router";
import { useWalletClient } from "wagmi";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { match } from "ts-pattern";
import { OnboardMethods } from "@/pages/onboard/[method]";
import { useAuth } from "@/context/AuthContext";

// eg. Sign in, Sign up, onboard page
export const useAuthPageRedirection = () => {
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
            );
    }, [
        accessToken,
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

    // Handle if the user context is authenticated
    const { user } = useAuth();
    useEffect(() => {
        if (user) {
            navigate("/app");
        }
    }, [user, navigate]);

    useEffect(() => {
        authCheckGoogle();
        authCheckWallet();
    }, [authCheckGoogle, authCheckWallet]);
};
