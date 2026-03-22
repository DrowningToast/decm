import React, { createContext, useEffect, useState } from "react";
import { useAccount, usePublicClient } from "wagmi";
import type { PublicClient } from "viem";

export interface WalletContextType {
    address: string | undefined;
    isConnected: boolean;
    isConnecting: boolean;
    isDisconnected: boolean;
    chainId: number | undefined;
    publicClient: PublicClient | undefined;
    isLoading: boolean;
    isReady: boolean;
}

export const WalletContext = createContext<WalletContextType | undefined>(undefined);

/**
 * WalletProvider Component
 *
 * Provides wallet connection state and public client to all child components.
 * Logs wallet connection status to console for debugging.
 *
 * Usage:
 * <WalletProvider>
 *   <YourApp />
 * </WalletProvider>
 */
export const WalletProvider = ({ children }: { children: React.ReactNode }) => {
    const { address, isConnected, isConnecting, isDisconnected, chainId } = useAccount();
    const publicClient = usePublicClient();
    const [prevAddress, setPrevAddress] = useState<string | undefined>(undefined);
    const [isPublicClientLoading, setIsPublicClientLoading] = useState(true);
    const [hasLoggedConnection, setHasLoggedConnection] = useState(false);

    // Check if public client is ready
    useEffect(() => {
        if (publicClient) {
            setIsPublicClientLoading(false);
        } else if (isConnected) {
            // Still loading if connected but no public client yet
            setIsPublicClientLoading(true);
        }
    }, [publicClient, isConnected]);

    const isLoading = isConnecting || isPublicClientLoading;
    const isReady = isConnected && publicClient !== undefined && !isLoading;

    // Log wallet connection status changes (development mode only)
    useEffect(() => {
        const isDev = process.env.NODE_ENV === "development";

        if (isReady && address && !hasLoggedConnection) {
            if (isDev) {
                console.log("=== WALLET CONNECTED ===");
                console.log("📍 Connected Address:", address);
                console.log("🔗 Chain ID:", chainId);
                console.log("⛓️ Public Client:", publicClient);
                console.log("📊 Public Client Details:", {
                    // Note: mode property removed - not available in current viem version
                    transport: {
                        type: publicClient?.transport?.type,
                        url: publicClient?.transport?.key === "http" ? "(HTTP)" : "(WebSocket)",
                    },
                    chain: {
                        name: publicClient?.chain?.name,
                        id: publicClient?.chain?.id,
                    },
                    key: publicClient?.key,
                });
                console.log("========================");
            }
            setHasLoggedConnection(true);
        } else if (isDisconnected && prevAddress) {
            if (isDev) {
                console.log("🔌 Wallet Disconnected");
            }
            setPrevAddress(undefined);
            setHasLoggedConnection(false);
        } else if (isConnecting) {
            if (isDev) {
                console.log("⏳ Wallet Connecting...");
            }
        } else if (isPublicClientLoading && address && !prevAddress) {
            if (isDev) {
                console.log("⏸️ Waiting for Public Client to initialize...");
            }
        }

        if (address !== prevAddress) {
            setPrevAddress(address);
        }
    }, [
        isReady,
        address,
        chainId,
        publicClient,
        isConnecting,
        isDisconnected,
        prevAddress,
        isPublicClientLoading,
        hasLoggedConnection,
    ]);

    const value: WalletContextType = {
        address,
        isConnected,
        isConnecting,
        isDisconnected,
        chainId,
        publicClient,
        isLoading,
        isReady,
    };

    return <WalletContext.Provider value={value}>{children}</WalletContext.Provider>;
};
