import { coreApiClient } from "@/lib/api/api";
import { useMutation } from "@tanstack/react-query";

export const useVerifyGoogleOAuth = () => {
    const { mutateAsync: verifyGoogleOAuth, isPending } = useMutation({
        mutationFn: async ({ code, state }: { code: string, state: string }) => {
            const response = await coreApiClient.v1.verifyGoogleOauth({ code, state });
            return response;
        }
    })

    return {
        verifyGoogleOAuth,
        isPending
    }
}