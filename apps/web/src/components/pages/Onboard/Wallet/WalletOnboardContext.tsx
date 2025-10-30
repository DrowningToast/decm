import type React from "react";
import { createContext, useCallback, useContext, useEffect, useMemo } from "react";
import { useGetSignMessage } from "../useGetSignMessage";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { handleUniversalError } from "@/common/Err";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { useSignout } from "@/components/useSignout";
import { onboardService } from "@/services/OnboardService";
import {
    OnboardRegistrationMethod,
    type OnboardCheckOnboardStatusResponse,
    type OnboardRegisterResponse,
} from "@decm/api";
import { authService } from "@/services/AuthService";
import { useSignup } from "../useSignup";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { toast } from "sonner";

type SignMessageState =
    | {
          signMessage?: string;
          isPending: true;
      }
    | {
          signMessage: string;
          isPending: false;
      };

type WalletOnboardContextType = {
    handleSubmit: () => Promise<void>;
} & SignMessageState;

const WalletOnboardContext = createContext<WalletOnboardContextType>(
    {} as WalletOnboardContextType,
);

const WalletOnboardProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { signout } = useSignout();
    const { onboardStatusError, setStep, onboardStatus, isStatusLoading } =
        useContext(OnboardPageContext);
    const [signSignature] = useLocalStorage<string | undefined>(
        LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
        undefined,
    );

    const { signMessage, isPending } = useGetSignMessage();
    const { upsertProfile } = useSignup();

    const signMessageState = useMemo<SignMessageState>(() => {
        if (isPending || !signMessage || typeof signMessage !== "string") {
            return {
                signMessage: undefined,
                isPending: true,
            };
        } else {
            return {
                signMessage: signMessage,
                isPending: false,
            };
        }
    }, [signMessage, isPending]);

    useEffect(() => {
        const init = async () => {
            if (onboardStatusError) {
                handleUniversalError(t, onboardStatusError, {
                    onInvalidInput: async () => {
                        await signout();
                        navigate("/");
                    },
                    unauthorizedErr: {
                        title: t("onboard.error.invalidSignature"),
                        description: t("onboard.error.invalidSignatureDescription"),
                        toastType: "error",
                        name: "invalid_signature",
                        message: t("onboard.error.invalidSignatureMessage"),
                    },
                });
            }
            if (isStatusLoading) {
                return;
            }
            if (onboardStatus?.profile_id) {
                setStep(2);
                return;
            }
            if (onboardStatus) {
                setStep(1);
                return;
            }
            setStep(0);
            return;
        };
        init();
    }, [
        signMessageState.signMessage,
        onboardStatus,
        isStatusLoading,
        setStep,
        onboardStatusError,
        t,
        navigate,
        signout,
    ]);

    const handleSubmit = useCallback(async () => {
        if (isPending || !signMessage || !signSignature) {
            return;
        }
        let onboardStatus: OnboardCheckOnboardStatusResponse | null = null;
        try {
            onboardStatus = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: signSignature,
            });
        } catch (error) {
            if (error instanceof Error) {
                handleUniversalError(t, error);
            }
            return;
        }
        let account: OnboardRegisterResponse | undefined = undefined;
        try {
            if (!onboardStatus?.authentication_credential_id) {
                account = await authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodWallet,
                    signSignature: signSignature,
                });
            }
        } catch (error) {
            if (error instanceof Error) {
                handleUniversalError(t, error);
            }
            return;
        }
        if (account?.credential_id) {
            try {
                await upsertProfile({
                    method: OnboardRegistrationMethod.RegistrationMethodWallet,
                    signSignature: signSignature,
                    profile: {
                        authentication_credential_id: account.credential_id,
                    },
                });
            } catch (error) {
                if (error instanceof Error) {
                    handleUniversalError(t, error);
                }
                return;
            }
        }

        toast.success(t("flow.wallet.create_profile_success"));
        navigate("/app");
    }, [isPending, navigate, signMessage, signSignature, t, upsertProfile]);

    return (
        <WalletOnboardContext.Provider
            value={{
                ...signMessageState,
                handleSubmit,
            }}
        >
            {children}
        </WalletOnboardContext.Provider>
    );
};

export { WalletOnboardProvider, WalletOnboardContext };
