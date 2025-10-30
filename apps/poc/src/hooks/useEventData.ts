import { useQuery } from "@tanstack/react-query";
import { usePublicClient } from "wagmi";
import type { Address } from "viem";
import { EventABI } from "../abi/Event";
import { chainId, EVENT_ADDRESS } from "../lib/addresses";
import { EventStatus } from "../types/Event";

interface EventData {
    name: string;
    description: string;
    seatsCount: bigint;
    currentSeatsCount: bigint;
    status: EventStatus;
}

/**
 * Hook return type for useEventData
 * @property data - Event data including name, description, seats information, and status
 * @property isLoading - Loading state while fetching event data
 * @property error - Error object if the data fetch failed
 * @property isError - Boolean indicating if an error occurred
 */
export interface UseEventDataReturn {
    data: EventData | undefined;
    isLoading: boolean;
    error: Error | null;
    isError: boolean;
}

/**
 * Fetches event data from the blockchain contract
 *
 * @param address - The blockchain address of the event contract (defaults to EVENT_ADDRESS)
 * @returns {UseEventDataReturn} An object containing the event data, loading state, and error information
 *
 * @example
 * const { data, isLoading, error } = useEventData();
 *
 * @example
 * const { data, isLoading } = useEventData('0x1234...');
 */
export function useEventData(address: Address = EVENT_ADDRESS): UseEventDataReturn {
    const client = usePublicClient({ chainId });

    return useQuery<EventData>({
        queryKey: ["event-data", address],
        queryFn: async () => {
            if (!client) throw new Error("Public client unavailable");
            const [eventName, eventDescription, seatsCount, currentSeatsCount, status] =
                await Promise.all([
                    client.readContract({
                        address,
                        abi: EventABI,
                        functionName: "getEventName",
                    }),
                    client.readContract({
                        address,
                        abi: EventABI,
                        functionName: "getEventDescription",
                    }),
                    client.readContract({
                        address,
                        abi: EventABI,
                        functionName: "seatsCount",
                    }),
                    client.readContract({
                        address,
                        abi: EventABI,
                        functionName: "currentSeatsCount",
                    }),
                    client.readContract({
                        address,
                        abi: EventABI,
                        functionName: "eventStatus",
                    }),
                ]);

            return {
                name: eventName as string,
                description: eventDescription as string,
                seatsCount: seatsCount as bigint,
                currentSeatsCount: currentSeatsCount as bigint,
                status: Number(status) as EventStatus,
            };
        },
        enabled: Boolean(client),
        staleTime: 15_000,
    });
}
