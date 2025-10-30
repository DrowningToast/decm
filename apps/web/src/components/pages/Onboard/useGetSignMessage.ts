import { onboardService } from "@/services/OnboardService";
import { useQuery } from "@tanstack/react-query";

export interface UseGetSignMessageReturn {
    signMessage: unknown | undefined;
    isPending: boolean;
    error: unknown;
}

export const useGetSignMessage = (): UseGetSignMessageReturn => {
    const {
        data: signMessage,
        isPending,
        error,
    } = useQuery({
        queryKey: ["getSignMessage"],
        queryFn: async () => {
            const response = await onboardService.getSignMessage();
            return response;
        },
    });

    return {
        signMessage,
        isPending,
        error,
    };
};
