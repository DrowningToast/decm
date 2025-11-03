import { useEffect } from "react";
import { useAccount, usePublicClient } from "wagmi";

/**
 * WalletConnectionWrapper Component
 *
 * This wrapper component sits at the top of the React tree and monitors
 * the currently connected wallet. It logs wallet connection status and
 * the public client to help with debugging and monitoring.
 *
 * Usage: Wrap your app content with this component inside AppKitProvider
 */
export const WalletConnectionWrapper = ({ children }: { children: React.ReactNode }) => {
    const { address, isConnected, isConnecting, isDisconnected, chainId } = useAccount();
    const publicClient = usePublicClient();

    useEffect(() => {
        if (isConnected && address) {
            console.log("=== WALLET CONNECTED ===");
            console.log("📍 Connected Address:", address);
            console.log("🔗 Chain ID:", chainId);
            console.log("⛓️ Public Client:", publicClient);
            console.log("📊 Public Client Details:", {
                // Note: mode property removed - not available in current viem version
                transport: publicClient?.transport?.type,
                chain: publicClient?.chain?.name,
                chainId: publicClient?.chain?.id,
                key: publicClient?.key,
            });
            console.log("========================");
        } else if (isDisconnected) {
            console.log("🔌 Wallet Disconnected");
        } else if (isConnecting) {
            console.log("⏳ Wallet Connecting...");
        }
    }, [address, isConnected, isDisconnected, isConnecting, chainId, publicClient]);

    return <>{children}</>;
};
