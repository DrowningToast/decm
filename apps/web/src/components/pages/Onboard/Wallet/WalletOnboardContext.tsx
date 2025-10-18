import type React from "react"
import { createContext, useContext, useEffect, useMemo } from "react";
import { useGetSignMessage } from "../useGetSignMessage";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { handleUniversalError } from "@/common/Err";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { useSignout } from "@/components/useSignout";

type SignMessageState = {
    signMessage?: string
    isPending: true
} | {
    signMessage: string
    isPending: false
}

type WalletOnboardContextType = {
} & SignMessageState

const WalletOnboardContext = createContext<WalletOnboardContextType>({} as WalletOnboardContextType)

const WalletOnboardProvider: React.FC<React.PropsWithChildren> = ({ children }) => {

    const { t } = useTranslation()
    const navigate = useNavigate()
    const { signout } = useSignout()
    const { onboardStatusError } = useContext(OnboardPageContext)
    const { setStep, onboardStatus, isStatusLoading } = useContext(OnboardPageContext)

    const { signMessage, isPending } = useGetSignMessage();

    const signMessageState = useMemo<SignMessageState>(() => {
        if (isPending || !signMessage) {
            return {
                signMessage: "",
                isPending: true,
            }
        } else {
            return {
                signMessage: signMessage,
                isPending: false,
            }
        }
    }, [signMessage, isPending])

    useEffect(() => {
        const init = async () => {

            if (onboardStatusError) {
                handleUniversalError(t, onboardStatusError, {
                    onInvalidInput: async () => {
                        await signout()
                        navigate("/")
                    },
                    unauthorizedErr: {
                        title: t("onboard.error.invalidSignature"),
                        description: t("onboard.error.invalidSignatureDescription"),
                        toastType: "error",
                        name: "invalid_signature",
                        message: t("onboard.error.invalidSignatureMessage"),
                    }
                })
            }
            if (isStatusLoading) {
                return
            }
            if (onboardStatus?.profile_id) {
                setStep(2)
                return
            }
            if (onboardStatus) {
                setStep(1)
                return
            }
            setStep(0)
            return
        }
        init()
    }, [signMessageState.signMessage, onboardStatus, isStatusLoading, setStep, onboardStatusError, t, navigate, signout])

    return (
        <WalletOnboardContext.Provider value={{
            ...signMessageState,
        }}>
            {children}
        </WalletOnboardContext.Provider>
    )
}

export { WalletOnboardProvider, WalletOnboardContext }