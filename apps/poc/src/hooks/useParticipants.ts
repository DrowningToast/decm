import { useQuery } from "@tanstack/react-query";
import { usePublicClient } from "wagmi";
import type { Address } from "viem";
import { EventABI } from "../abi/Event";
import { chainId, EVENT_ADDRESS } from "../lib/addresses";

export function useParticipantMessage(
  participant: Address,
  eventAddress: Address = EVENT_ADDRESS,
) {
  const client = usePublicClient({ chainId });

  return useQuery({
    queryKey: ["event-signing-message", eventAddress, participant],
    queryFn: async () => {
      if (!client) throw new Error("Public client unavailable");
      if (!participant) return "0x";
      return client.readContract({
        address: eventAddress,
        abi: EventABI,
        functionName: "getSigningMessage",
        args: [participant],
      });
    },
    enabled: Boolean(client && participant),
  });
}
