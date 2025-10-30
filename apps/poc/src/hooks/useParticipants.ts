import { useQuery } from "@tanstack/react-query";
import { usePublicClient } from "wagmi";
import type { Address } from "viem";
import { EventABI } from "../abi/Event";
import { chainId, EVENT_ADDRESS } from "../lib/addresses";

/**
 * Fetches the signing message for a participant in an event.
 *
 * This hook retrieves the message that a participant needs to sign for verification
 * purposes in the specified event contract via a blockchain read call.
 *
 * @param participant - The Ethereum address of the participant
 * @param eventAddress - The Ethereum address of the event contract (defaults to EVENT_ADDRESS)
 * @returns A TanStack React Query result object containing the signing message, loading state, and error information
 */
export function useParticipantMessage(
    participant: Address,
    eventAddress: Address = EVENT_ADDRESS,
): ReturnType<typeof useQuery> {
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
