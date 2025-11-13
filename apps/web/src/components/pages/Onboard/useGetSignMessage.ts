import { onboardService } from "@/services/OnboardService/OnboardService";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";

export const useGetSignMessage = () => {
    const { data: signMessage, isPending } = useQuery({
        queryKey: QUERY_KEY.onboard.signMessage,
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
