import { coreApi } from "@/lib/api/api";
import { useMutation } from "@tanstack/react-query";

export const useGoogleOAuthUrl = () => {
    const { mutateAsync: requestGoogleOAuthUrl, isPending } = useMutation({
        mutationFn: async () => {
            const response = await coreApi.v1.requestGoogleOauth();
            return response.url;
        }
    })

    return {
        requestGoogleOAuthUrl,
        isPending
    }
}