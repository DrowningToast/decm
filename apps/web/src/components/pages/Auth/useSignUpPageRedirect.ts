// Use this hook when the user is in the authentication flow

import { useCallback, useEffect, useMemo } from "react";
import { useCheckOnboardStatus } from "../Onboard/useCheckOnboardStatus";
import { useNavigate } from "@/router";
import { useAccount } from "wagmi";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { match } from "ts-pattern";
import { OnboardMethods } from "@/pages/onboard/[method]";

export const useSignUpPageRedirect = () => {
    const navigate = useNavigate();
    const { onboardStatus, isLoading: isOnboardStatusLoading } = useCheckOnboardStatus();
    const { address, isReconnecting, isDisconnected, isConnecting } = useAccount();
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
        isLoading,
        navigate,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
    ]);

    const authCheckWallet = useCallback(async () => {
        match({
            isLoading,
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
            );
    }, [
        authSignSignature,
        isLoading,
        navigate,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
        address,
    ]);

    useEffect(() => {
        authCheckGoogle();
        authCheckWallet();
    }, [authCheckGoogle, authCheckWallet]);

    return {
        isLoading,
    };
};
