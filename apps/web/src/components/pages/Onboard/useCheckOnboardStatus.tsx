import { coreApi } from "@/lib/api/api";
import { useMutation } from "@tanstack/react-query"
import type { OnboardRegistrationMethod } from "@decm/api";

export const useCheckOnboardStatus = () => {

    const { mutateAsync: checkOnboardStatus, isPending } = useMutation({
        mutationFn: async ({ method, accessToken, expiresIn, signMessage }: { method: OnboardRegistrationMethod, accessToken?: string, expiresIn?: number, signMessage?: string }) => {
            const response = await coreApi.v1.checkOnboardStatus({
                method,
                access_token: accessToken,
                expires_in: expiresIn,
                sign_message: signMessage,
            });

            return response;
        },
    });

    return { checkOnboardStatus, isPending };
};