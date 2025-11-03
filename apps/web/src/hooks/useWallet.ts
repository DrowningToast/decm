import { useContext } from "react";
import { WalletContext } from "@/context/WalletContext";
import type { WalletContextType } from "@/context/WalletContext";

/**
 * Hook to use wallet context
 *
 * Provides access to wallet connection state and public client throughout the app.
 *
 * Usage:
 * const { address, publicClient, isConnected, isReady, isLoading } = useWallet();
 *
 * Properties:
 * - address: Connected wallet address
 * - isConnected: Wallet is connected
 * - isConnecting: Wallet is currently connecting
 * - isDisconnected: Wallet is disconnected
 * - chainId: Current chain ID
 * - publicClient: Viem PublicClient instance (for blockchain reads)
 * - isLoading: Wallet is connecting OR public client is initializing
 * - isReady: Wallet is fully connected AND public client is ready (use this before accessing wallet data)
 *
 * @throws Error if used outside of WalletProvider
 * @returns WalletContextType with wallet state and public client
 */
export const useWallet = (): WalletContextType => {
    const context = useContext(WalletContext);
    if (context === undefined) {
        throw new Error("useWallet must be used within a WalletProvider");
    }
    return context;
};
