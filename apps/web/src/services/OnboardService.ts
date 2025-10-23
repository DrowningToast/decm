import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import { OnboardRegistrationMethod } from "@decm/api";

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

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    async checkOnboardStatus(params?: CheckOnboardParams) {
        if (params?.method === OnboardRegistrationMethod.RegistrationMethodGoogle) {
            if (!params.accessToken || !params.expiresIn) {
                console.error("Invalid access token or expires in");
                return null;
            }
            return this._coreApi.v1.checkOnboardStatus({
                method: params.method,
                access_token: params.accessToken,
                expires_in: params.expiresIn,
            });
        } else if (params?.method === OnboardRegistrationMethod.RegistrationMethodWallet) {
            if (!params?.signSignature) {
                console.error("Invalid sign signature");
                return null;
            }
            return this._coreApi.v1.checkOnboardStatus({
                method: params.method,
                message_signature: params.signSignature,
            });
        } else if (!params) {
            // check via jwt cookie
            return this._coreApi.v1.checkOnboardStatus({});
        }
        throw new Error("Invalid method");
    }

    async getSignMessage() {
        const response = await this._coreApi.v1.getSignMessage();
        return response.message;
    }
}

const DefaultOnboardService = new OnboardService(coreApiClient);
export { DefaultOnboardService as onboardService };
