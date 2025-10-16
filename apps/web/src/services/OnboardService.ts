import { Err } from "@/common/Err";
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
			signMessage: string;
	  };

export class OnboardService {
	private _coreApi: CoreApiType;

	constructor(coreApi: CoreApiType) {
		this._coreApi = coreApi;
	}

	async checkOnboardStatus(params?: CheckOnboardParams) {
		if (params?.method === OnboardRegistrationMethod.RegistrationMethodGoogle) {
			if (!params.accessToken || !params.expiresIn) {
				throw new Err("Invalid access token or expires in");
			}
			return this._coreApi.v1.checkOnboardStatus({
				method: params.method,
				access_token: params.accessToken,
				expires_in: params.expiresIn,
			});
		} else if (
			params?.method === OnboardRegistrationMethod.RegistrationMethodWallet
		) {
			if (!params?.signMessage) {
				throw new Err("Invalid sign message");
			}
			return this._coreApi.v1.checkOnboardStatus({
				method: params.method,
				sign_message: params.signMessage,
			});
		} else if (!params) {
			// check via jwt cookie
			return this._coreApi.v1.checkOnboardStatus({});
		}
		throw new Error("Invalid method");
	}
}

const DefaultOnboardService = new OnboardService(coreApiClient);
export { DefaultOnboardService as onboardService };
