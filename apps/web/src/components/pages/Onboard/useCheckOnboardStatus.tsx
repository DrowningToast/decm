import { useMutation, useQuery } from "@tanstack/react-query"
import { onboardService, type CheckOnboardParams } from "../../../services/OnboardService";
import { OnboardRegistrationMethod } from "@decm/api";

const getQueryKey = (param?: CheckOnboardParams) => {
    switch (param?.method) {
        case OnboardRegistrationMethod.RegistrationMethodGoogle:
            return ["onboardStatus", param.accessToken, param.expiresIn]
        case OnboardRegistrationMethod.RegistrationMethodWallet:
            return ["onboardStatus", param.signMessage]
    }

    return ["onboardStatus"]
}

export const useCheckOnboardStatus = (param?: CheckOnboardParams) => {
    const { mutateAsync: checkOnboardStatus, isPending } = useMutation({
        mutationFn: async (param?: CheckOnboardParams) => {
            const response = await onboardService.checkOnboardStatus(param);
            return response;
        },
    });

    const { data: onboardStatus, isLoading, error } = useQuery({
        queryKey: getQueryKey(param),
        queryFn: async () => {
            const response = await checkOnboardStatus(param);

            return response;
        },
    });

    return { checkOnboardStatus, isPending, isLoading, onboardStatus, error };
};