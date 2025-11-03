import { useMutation, useQuery } from "@tanstack/react-query";
import { onboardService, type CheckOnboardParams } from "../../../services/OnboardService";
import { OnboardRegistrationMethod } from "@decm/api";
import { QUERY_KEY } from "@/lib/queryKeys";

export type UseCheckOnboardParams =
    | {
          method?: OnboardRegistrationMethod.RegistrationMethodGoogle;
          accessToken: string;
          expiresIn: number;
      }
    | {
          method?: OnboardRegistrationMethod.RegistrationMethodWallet;
          signSignature?: string;
      };

export const useCheckOnboardStatus = (param?: UseCheckOnboardParams, enable: boolean = true) => {
    const { mutateAsync: checkOnboardStatus, isPending } = useMutation({
        mutationFn: async (param?: CheckOnboardParams) => {
            const response = await onboardService.checkOnboardStatus(param);
            return response;
        },
    });

    const getQueryKey = (param?: UseCheckOnboardParams) => {
        switch (param?.method) {
            case OnboardRegistrationMethod.RegistrationMethodGoogle:
                return QUERY_KEY.onboard.status.google(param.accessToken, param.expiresIn);
            case OnboardRegistrationMethod.RegistrationMethodWallet:
                return QUERY_KEY.onboard.status.wallet(param.signSignature ?? "");
        }
        return QUERY_KEY.onboard.status.all;
    };

    const {
        data: onboardStatus,
        isLoading,
        error,
    } = useQuery({
        queryKey: getQueryKey(param),
        queryFn: async () => {
            if (param?.method === OnboardRegistrationMethod.RegistrationMethodWallet) {
                const { signSignature } = param ?? {};
                if (signSignature === undefined) {
                    return null;
                }
                return await checkOnboardStatus({
                    method: param.method,
                    signSignature: signSignature,
                });
            }
            return await checkOnboardStatus(param as CheckOnboardParams);
        },
        enabled: enable,
    });

    return { checkOnboardStatus, isPending, isLoading, onboardStatus, error };
};
