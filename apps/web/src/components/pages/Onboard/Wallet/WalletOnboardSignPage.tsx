import { Button } from "@/components/ui/button"
import { useCallback, useContext, } from "react"
import { WalletOnboardContext } from "./WalletOnboardContext"
import { useSignMessage, useWalletClient } from "wagmi"
import { ErrorPage } from "../../Error"
import { useLocalStorage } from "@/hooks/use-local-storage"
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage"
import { LogoutButton } from "@/components/LogoutButton"

export const WalletOnboardSignPage = () => {

    const { data: walletClient } = useWalletClient()
    const { signMessageAsync } = useSignMessage()
    const { signMessage, isPending, } = useContext(WalletOnboardContext)
    const [, setSignSignature] = useLocalStorage<string | undefined>(LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE, undefined)

    const handleSignSignature = useCallback(async () => {
        if (!walletClient || !signMessage) {
            return
        }
        const signature = await signMessageAsync({ message: signMessage })
        setSignSignature(signature)
    }, [walletClient, signMessage, signMessageAsync, setSignSignature])

    if (isPending) {
        return <div>Loading...</div>
    }

    if (!walletClient) {
        return <ErrorPage />
    }

    return (
        <>
            <Button onClick={handleSignSignature}>
                Sign the message
            </Button>
            <LogoutButton type="disconnect" />
        </>
    )
}