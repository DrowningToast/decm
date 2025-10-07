import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import { useAccount, usePublicClient } from "wagmi";
import type { Address } from "viem";
import { zeroAddress } from "viem";
import { EventAccessManagerABI } from "../abi/EventAccessManager";
import {
  ACCESS_MANAGER_ADDRESS,
  chainId,
  EVENT_ADDRESS,
} from "../lib/addresses";

interface WalletContextValue {
  isHost: boolean;
  isParticipant: boolean;
  isIssuer: boolean;
  refreshRoles: () => void;
}

const WalletContext = createContext<WalletContextValue | undefined>(undefined);

interface WalletProviderProps {
  children: ReactNode;
  onRolesUpdated?: (roles: RolesState) => void;
}

const defaultRoles: RolesState = {
  isHost: false,
  isIssuer: false,
  isParticipant: false,
};

export function WalletProvider({
  children,
  onRolesUpdated,
}: WalletProviderProps) {
  const { address } = useAccount();
  const publicClient = usePublicClient({ chainId });
  const [roles, setRoles] = useState<RolesState>(defaultRoles);
  const account = address ?? zeroAddress;

  const refreshRoles = async () => {
    if (!publicClient) return;

    const responses = await Promise.allSettled([
      publicClient.readContract({
        address: EVENT_ADDRESS,
        abi: EventAccessManagerABI,
        functionName: "checkIsHost",
        args: [account],
      }),
      publicClient.readContract({
        address: EVENT_ADDRESS,
        abi: EventAccessManagerABI,
        functionName: "checkIsParticipant",
        args: [account],
      }),
      publicClient.readContract({
        address: EVENT_ADDRESS,
        abi: EventAccessManagerABI,
        functionName: "checkIsIssuer",
        args: [account],
      }),
    ]);

    const nextRoles: RolesState = {
      isHost: responses[0].status === "fulfilled" ? responses[0].value : false,
      isParticipant:
        responses[1].status === "fulfilled" ? responses[1].value : false,
      isIssuer:
        responses[2].status === "fulfilled" ? responses[2].value : false,
    };

    setRoles(nextRoles);
    onRolesUpdated?.(nextRoles);
  };

  useEffect(() => {
    refreshRoles().catch(console.error);
  }, [account, publicClient]);

  return (
    <WalletContext.Provider value={{ ...roles, refreshRoles }}>
      {children}
    </WalletContext.Provider>
  );
}

export function useWalletContext() {
  const context = useContext(WalletContext);
  if (!context) {
    throw new Error("useWalletContext must be used within a WalletProvider");
  }
  return context;
}
