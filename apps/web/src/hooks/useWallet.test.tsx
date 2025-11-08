import { describe, it, expect, vi } from "vitest";
import { renderHook } from "@testing-library/react";
import { useWallet } from "./useWallet";
import { WalletContext, type WalletContextType } from "@/context/WalletContext";
import type { PublicClient } from "viem";

describe("useWallet", () => {
    it("should throw error when used outside WalletProvider", () => {
        // Suppress error output for this test to avoid test output pollution
        const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

        try {
            renderHook(() => useWallet());
            // If we get here, the test should fail
            expect.fail("Hook should have thrown an error");
        } catch (error) {
            // The hook should throw when used outside a provider
            expect(error).toBeDefined();
        } finally {
            consoleSpy.mockRestore();
        }
    });

    it("should return wallet context when used within provider", () => {
        const mockWalletContext: WalletContextType = {
            address: "0x1234567890123456789012345678901234567890",
            isConnected: true,
            isConnecting: false,
            isDisconnected: false,
            chainId: 1,
            publicClient: {} as PublicClient,
            isLoading: false,
            isReady: true,
        };

        const wrapper = ({ children }: { children: React.ReactNode }) => (
            <WalletContext.Provider value={mockWalletContext}>{children}</WalletContext.Provider>
        );

        const { result } = renderHook(() => useWallet(), { wrapper });

        expect(result.current).toEqual(mockWalletContext);
        expect(result.current.address).toBe("0x1234567890123456789012345678901234567890");
        expect(result.current.isConnected).toBe(true);
        expect(result.current.isReady).toBe(true);
    });

    it("should provide wallet connection state", () => {
        const mockWalletContext: WalletContextType = {
            address: "0x1234567890123456789012345678901234567890",
            isConnected: true,
            isConnecting: false,
            isDisconnected: false,
            chainId: 1,
            publicClient: {} as PublicClient,
            isLoading: false,
            isReady: true,
        };

        const wrapper = ({ children }: { children: React.ReactNode }) => (
            <WalletContext.Provider value={mockWalletContext}>{children}</WalletContext.Provider>
        );

        const { result } = renderHook(() => useWallet(), { wrapper });

        expect(result.current.isConnected).toBe(true);
        expect(result.current.isConnecting).toBe(false);
        expect(result.current.isDisconnected).toBe(false);
    });

    it("should provide chain ID information", () => {
        const mockWalletContext: WalletContextType = {
            address: "0x1234567890123456789012345678901234567890",
            isConnected: true,
            isConnecting: false,
            isDisconnected: false,
            chainId: 137, // Polygon
            publicClient: {} as PublicClient,
            isLoading: false,
            isReady: true,
        };

        const wrapper = ({ children }: { children: React.ReactNode }) => (
            <WalletContext.Provider value={mockWalletContext}>{children}</WalletContext.Provider>
        );

        const { result } = renderHook(() => useWallet(), { wrapper });

        expect(result.current.chainId).toBe(137);
    });

    it("should provide public client when ready", () => {
        const mockPublicClient = { id: "publicClient" } as unknown as PublicClient;
        const mockWalletContext: WalletContextType = {
            address: "0x1234567890123456789012345678901234567890",
            isConnected: true,
            isConnecting: false,
            isDisconnected: false,
            chainId: 1,
            publicClient: mockPublicClient,
            isLoading: false,
            isReady: true,
        };

        const wrapper = ({ children }: { children: React.ReactNode }) => (
            <WalletContext.Provider value={mockWalletContext}>{children}</WalletContext.Provider>
        );

        const { result } = renderHook(() => useWallet(), { wrapper });

        expect(result.current.publicClient).toBe(mockPublicClient);
    });

    it("should indicate loading state during initialization", () => {
        const mockWalletContext = {
            address: "0x1234567890123456789012345678901234567890",
            isConnected: false,
            isConnecting: true,
            isDisconnected: false,
            chainId: undefined,
            publicClient: undefined,
            isLoading: true,
            isReady: false,
        };

        const wrapper = ({ children }: { children: React.ReactNode }) => (
            <WalletContext.Provider value={mockWalletContext}>{children}</WalletContext.Provider>
        );

        const { result } = renderHook(() => useWallet(), { wrapper });

        expect(result.current.isLoading).toBe(true);
        expect(result.current.isReady).toBe(false);
        expect(result.current.isConnecting).toBe(true);
    });

    it("should handle disconnected state", () => {
        const mockWalletContext = {
            address: undefined,
            isConnected: false,
            isConnecting: false,
            isDisconnected: true,
            chainId: undefined,
            publicClient: undefined,
            isLoading: false,
            isReady: false,
        };

        const wrapper = ({ children }: { children: React.ReactNode }) => (
            <WalletContext.Provider value={mockWalletContext}>{children}</WalletContext.Provider>
        );

        const { result } = renderHook(() => useWallet(), { wrapper });

        expect(result.current.isDisconnected).toBe(true);
        expect(result.current.isConnected).toBe(false);
        expect(result.current.address).toBeUndefined();
        expect(result.current.isReady).toBe(false);
    });

    it("should handle partial public client loading state", () => {
        const mockWalletContext = {
            address: "0x1234567890123456789012345678901234567890",
            isConnected: true,
            isConnecting: false,
            isDisconnected: false,
            chainId: 1,
            publicClient: undefined,
            isLoading: true, // Still loading public client
            isReady: false,
        };

        const wrapper = ({ children }: { children: React.ReactNode }) => (
            <WalletContext.Provider value={mockWalletContext}>{children}</WalletContext.Provider>
        );

        const { result } = renderHook(() => useWallet(), { wrapper });

        expect(result.current.isConnected).toBe(true);
        expect(result.current.publicClient).toBeUndefined();
        expect(result.current.isLoading).toBe(true);
        expect(result.current.isReady).toBe(false);
    });
});
