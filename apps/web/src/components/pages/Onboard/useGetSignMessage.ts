import { onboardService } from "@/services/OnboardService";
import { useQuery } from "@tanstack/react-query";

export const useGetSignMessage = () => {
    const { data: signMessage, isPending } = useQuery({
        queryKey: ["getSignMessage"],
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
