import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { OnboardRegistrationMethod } from "@decm/api";
import { QueryClient } from "@tanstack/react-query";

export type CheckOnboardParams =
    | {
          method?: OnboardRegistrationMethod.RegistrationMethodGoogle;
          accessToken: string;
          expiresIn: number;
      }
    | {
          method?: OnboardRegistrationMethod.RegistrationMethodWallet;
          signSignature: string;
      };

export class OnboardService {
    private _coreApi: CoreApiType;
    private _queryClient: QueryClient;

    constructor(coreApi: CoreApiType, queryClient: QueryClient) {
        this._queryClient = queryClient;
        this._coreApi = coreApi;
    }

    async checkOnboardStatus(params?: CheckOnboardParams) {
        if (params?.method === OnboardRegistrationMethod.RegistrationMethodGoogle) {
            if (!params.accessToken || !params.expiresIn) {
                console.error("Invalid access token or expires in");
                return null;
            }
            const response = await this._coreApi.v1.checkOnboardStatus({
                method: params.method,
                access_token: params.accessToken,
                expires_in: params.expiresIn,
            });
            await this._queryClient.setQueryData(QUERY_KEY.onboard.status.all, response);
            return response;
        } else if (params?.method === OnboardRegistrationMethod.RegistrationMethodWallet) {
            if (!params?.signSignature) {
                console.error("Invalid sign signature");
                return null;
            }
            const response = await this._coreApi.v1.checkOnboardStatus({
                method: params.method,
                message_signature: params.signSignature,
            });
            await this._queryClient.setQueryData(
                QUERY_KEY.onboard.status.wallet(params.signSignature),
                response,
            );
            return response;
        } else if (!params) {
            // check via jwt cookie
            const response = await this._coreApi.v1.checkOnboardStatus({});
            await this._queryClient.setQueryData(QUERY_KEY.onboard.status.all, response);
            return response;
        }
        throw new Error("Invalid method");
    }

    async getSignMessage() {
        try {
            const response = await this._coreApi.v1.getSignMessage();
            return response.message;
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : "Unknown error";
            console.error(`Failed to fetch sign message: ${errorMessage}`, error);
            throw new Error(`Failed to retrieve signing message: ${errorMessage}`);
        }
    }
}

export const defaultOnboardService = new OnboardService(coreApiClient, queryClient);
