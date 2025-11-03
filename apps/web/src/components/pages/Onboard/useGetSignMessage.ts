import { onboardService } from "@/services/OnboardService";
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";

export const useGetSignMessage = () => {
    const { data: signMessage, isPending } = useQuery({
        queryKey: queryKeys.onboard.signMessage,
        queryFn: async () => {
            const response = await onboardService.getSignMessage();
            return response;
        },
    });

    return {
        signMessage,
        isPending,
    };
};
