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

export function useEventData(address: Address = EVENT_ADDRESS) {
  const client = usePublicClient({ chainId });

  return useQuery<EventData>({
    queryKey: ["event-data", address],
    queryFn: async () => {
      if (!client) throw new Error("Public client unavailable");
      const [
        eventName,
        eventDescription,
        seatsCount,
        currentSeatsCount,
        status,
      ] = await Promise.all([
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
