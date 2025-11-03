import type React from "react";
import { createContext, useCallback, useContext, useEffect, useMemo } from "react";
import { useGetSignMessage } from "../useGetSignMessage";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { useSignout } from "@/components/useSignout";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { useWalletOnboardUsecase } from "../Usecase/useWalletOnboardUsecase";

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
    const { onboardStatusError, setStep, onboardStatus, isStatusLoading, profileForm } =
        useContext(OnboardPageContext);
    const [signSignature] = useLocalStorage<string | undefined>(
        LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
        undefined,
    );
    const { usecaseAsync } = useWalletOnboardUsecase();
    const { signMessage, isPending } = useGetSignMessage();

    const signMessageState = useMemo<SignMessageState>(() => {
        if (isPending || !signMessage) {
            return {
                signMessage: "",
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

        await usecaseAsync(signSignature, profileForm.getValues());
    }, [isPending, signMessage, signSignature, usecaseAsync, profileForm]);

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
