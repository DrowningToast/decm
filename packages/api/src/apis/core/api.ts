/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export type CheckOnboardStatusData = OnboardCheckOnboardStatusResponse;

export type CheckOnboardStatusError = CustomerrorErrResponse;

export type CreateProfileData = ProfileCreateProfileResponse;

export type CreateProfileError = CustomerrorErr;

/** Custom error type */
export interface CustomerrorErr {
	code?: string;
	http_status?: number;
	message?: string;
	reasons?: Record<string, string>;
}

/** Response for the client to sign to register */
export interface CustomerrorErrResponse {
	message?: string;
}

export interface EntityProfile {
	academic_email?: string;
	academic_institution?: string;
	address?: string;
	authentication_credential_id?: string;
	bio?: string;
	created_at?: string;
	email?: string;
	first_name?: string;
	id?: string;
	is_academic_email_public?: boolean;
	is_academic_institution_public?: boolean;
	is_address_public?: boolean;
	is_bio_public?: boolean;
	is_email_public?: boolean;
	is_first_name_public?: boolean;
	is_last_name_public?: boolean;
	is_phone_number_public?: boolean;
	is_profile_picture_public?: boolean;
	last_name?: string;
	phone_number?: string;
	profile_picture_url?: string;
	updated_at?: string;
}

export type GetMyProfileData = EntityProfile;

export type GetMyProfileError = CustomerrorErr;

export type GetSignMessageData = OnboardGetSignMessageResponse;

export type LogoutData = Record<string, string>;

export type LogoutError = CustomerrorErrResponse;

export interface OnboardCheckOnboardStatusRequest {
	access_token?: string;
	expires_in?: number;
	message_signature?: string;
	method?: "google" | "wallet";
}

export interface OnboardCheckOnboardStatusResponse {
	authentication_credential_id?: string;
	profile_id?: string;
}

/** Response for the client to sign to register */
export interface OnboardGetSignMessageResponse {
	message: string;
}

export interface OnboardRegisterResponse {
	credential_id: string;
	jwt: string;
}

export interface OnboardRegisterWithGoogleOAuthRequest {
	access_token?: string;
	/** @minLength 6 */
	password: string;
}

export interface OnboardRegisterWithWalletRequest {
	signed_message?: string;
}

export enum OnboardRegistrationMethod {
	RegistrationMethodGoogle = "google",
	RegistrationMethodWallet = "wallet",
}

export interface ProfileCreateProfileRequest {
	academic_email?: string;
	/**
	 * @minLength 3
	 * @maxLength 255
	 */
	academic_institution?: string;
	/**
	 * @minLength 10
	 * @maxLength 255
	 */
	address?: string;
	authentication_credential_id: string;
	/**
	 * @minLength 10
	 * @maxLength 255
	 */
	bio?: string;
	/** @maxLength 64 */
	email?: string;
	/**
	 * @minLength 3
	 * @maxLength 32
	 */
	first_name?: string;
	is_academic_email_public?: boolean;
	is_academic_institution_public?: boolean;
	is_address_public?: boolean;
	is_bio_public?: boolean;
	is_email_public?: boolean;
	is_first_name_public?: boolean;
	is_last_name_public?: boolean;
	is_phone_number_public?: boolean;
	is_profile_picture_public?: boolean;
	/**
	 * @minLength 3
	 * @maxLength 32
	 */
	last_name?: string;
	phone_number?: string;
	profile_picture_url?: string;
}

export interface ProfileCreateProfileResponse {
	academic_email?: string;
	academic_institution?: string;
	address?: string;
	authentication_credential_id?: string;
	bio?: string;
	created_at?: string;
	email?: string;
	first_name?: string;
	id?: string;
	is_academic_email_public?: boolean;
	is_academic_institution_public?: boolean;
	is_address_public?: boolean;
	is_bio_public?: boolean;
	is_email_public?: boolean;
	is_first_name_public?: boolean;
	is_last_name_public?: boolean;
	is_phone_number_public?: boolean;
	is_profile_picture_public?: boolean;
	last_name?: string;
	phone_number?: string;
	profile_picture_url?: string;
	updated_at?: string;
}

export interface ProfileUpdateProfileRequest {
	academic_email?: string;
	/**
	 * @minLength 3
	 * @maxLength 255
	 */
	academic_institution?: string;
	/**
	 * @minLength 10
	 * @maxLength 255
	 */
	address?: string;
	/**
	 * @minLength 10
	 * @maxLength 255
	 */
	bio?: string;
	/** @maxLength 64 */
	email?: string;
	/**
	 * @minLength 3
	 * @maxLength 32
	 */
	first_name?: string;
	is_academic_email_public?: boolean;
	is_academic_institution_public?: boolean;
	is_address_public?: boolean;
	is_bio_public?: boolean;
	is_email_public?: boolean;
	is_first_name_public?: boolean;
	is_last_name_public?: boolean;
	is_phone_number_public?: boolean;
	is_profile_picture_public?: boolean;
	/**
	 * @minLength 3
	 * @maxLength 32
	 */
	last_name?: string;
	phone_number?: string;
	profile_picture_url?: string;
}

export interface ProfileUpdateProfileResponse {
	academic_email?: string;
	academic_institution?: string;
	address?: string;
	authentication_credential_id?: string;
	bio?: string;
	created_at?: string;
	email?: string;
	first_name?: string;
	id?: string;
	is_academic_email_public?: boolean;
	is_academic_institution_public?: boolean;
	is_address_public?: boolean;
	is_bio_public?: boolean;
	is_email_public?: boolean;
	is_first_name_public?: boolean;
	is_last_name_public?: boolean;
	is_phone_number_public?: boolean;
	is_profile_picture_public?: boolean;
	last_name?: string;
	phone_number?: string;
	profile_picture_url?: string;
	updated_at?: string;
}

export type RegisterWithGoogleOauthData = OnboardRegisterResponse;

export type RegisterWithGoogleOauthError = CustomerrorErrResponse;

export type RegisterWithWalletData = OnboardRegisterResponse;

export type RegisterWithWalletError = CustomerrorErrResponse;

export type RequestGoogleOauthError = CustomerrorErrResponse;

export type UpdateProfileByCredentialIdData = ProfileUpdateProfileResponse;

export type UpdateProfileByCredentialIdError = CustomerrorErr;

export interface UpdateProfileByCredentialIdParams {
	/** Credential ID */
	credentialId: string;
}

export type VerifyGoogleOauthError = CustomerrorErrResponse;

export interface VerifyGoogleOauthParams {
	code?: string;
	state?: string;
}

import type {
	AxiosInstance,
	AxiosRequestConfig,
	HeadersDefaults,
	ResponseType,
} from "axios";
import axios from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams
	extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
	/** set parameter to `true` for call `securityWorker` for this request */
	secure?: boolean;
	/** request path */
	path: string;
	/** content type of request body */
	type?: ContentType;
	/** query params */
	query?: QueryParamsType;
	/** format of response (i.e. response.json() -> format: "json") */
	format?: ResponseType;
	/** request body */
	body?: unknown;
}

export type RequestParams = Omit<
	FullRequestParams,
	"body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown>
	extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
	securityWorker?: (
		securityData: SecurityDataType | null
	) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
	secure?: boolean;
	format?: ResponseType;
}

export enum ContentType {
	Json = "application/json",
	JsonApi = "application/vnd.api+json",
	FormData = "multipart/form-data",
	UrlEncoded = "application/x-www-form-urlencoded",
	Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
	public instance: AxiosInstance;
	private securityData: SecurityDataType | null = null;
	private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
	private secure?: boolean;
	private format?: ResponseType;

	constructor({
		securityWorker,
		secure,
		format,
		...axiosConfig
	}: ApiConfig<SecurityDataType> = {}) {
		this.instance = axios.create({
			...axiosConfig,
			baseURL: axiosConfig.baseURL || "//localhost:8080",
		});
		this.secure = secure;
		this.format = format;
		this.securityWorker = securityWorker;
	}

	public setSecurityData = (data: SecurityDataType | null) => {
		this.securityData = data;
	};

	protected mergeRequestParams(
		params1: AxiosRequestConfig,
		params2?: AxiosRequestConfig
	): AxiosRequestConfig {
		const method = params1.method || (params2 && params2.method);

		return {
			...this.instance.defaults,
			...params1,
			...(params2 || {}),
			headers: {
				...((method &&
					this.instance.defaults.headers[
						method.toLowerCase() as keyof HeadersDefaults
					]) ||
					{}),
				...(params1.headers || {}),
				...((params2 && params2.headers) || {}),
			},
		};
	}

	protected stringifyFormItem(formItem: unknown) {
		if (typeof formItem === "object" && formItem !== null) {
			return JSON.stringify(formItem);
		} else {
			return `${formItem}`;
		}
	}

	protected createFormData(input: Record<string, unknown>): FormData {
		if (input instanceof FormData) {
			return input;
		}
		return Object.keys(input || {}).reduce((formData, key) => {
			const property = input[key];
			const propertyContent: any[] =
				property instanceof Array ? property : [property];

			for (const formItem of propertyContent) {
				const isFileType = formItem instanceof Blob || formItem instanceof File;
				formData.append(
					key,
					isFileType ? formItem : this.stringifyFormItem(formItem)
				);
			}

			return formData;
		}, new FormData());
	}

	public request = async <T = any, _E = any>({
		secure,
		path,
		type,
		query,
		format,
		body,
		...params
	}: FullRequestParams): Promise<T> => {
		const secureParams =
			((typeof secure === "boolean" ? secure : this.secure) &&
				this.securityWorker &&
				(await this.securityWorker(this.securityData))) ||
			{};
		const requestParams = this.mergeRequestParams(params, secureParams);
		const responseFormat = format || this.format || undefined;

		if (
			type === ContentType.FormData &&
			body &&
			body !== null &&
			typeof body === "object"
		) {
			body = this.createFormData(body as Record<string, unknown>);
		}

		if (
			type === ContentType.Text &&
			body &&
			body !== null &&
			typeof body !== "string"
		) {
			body = JSON.stringify(body);
		}

		return this.instance
			.request({
				...requestParams,
				headers: {
					...(requestParams.headers || {}),
					...(type ? { "Content-Type": type } : {}),
				},
				params: query,
				responseType: responseFormat,
				data: body,
				url: path,
			})
			.then((response) => response.data);
	};
}

/**
 * @title DECM Core
 * @version 1.0
 * @license Apache 2.0 (http://www.apache.org/licenses/LICENSE-2.0.html)
 * @termsOfService http://swagger.io/terms/
 * @baseUrl //localhost:8080
 * @contact
 *
 * DECM (Decentralized Event Management) platform API for NFT ticketing, digital credentials, and academic identity verification.
 */
export class Api<SecurityDataType extends unknown> {
	http: HttpClient<SecurityDataType>;

	constructor(http: HttpClient<SecurityDataType>) {
		this.http = http;
	}

	v1 = {
		/**
		 * @description Logout user by clearing session and OAuth cookies
		 *
		 * @tags Auth
		 * @name Logout
		 * @summary Logout
		 * @request POST:/api/v1/auth/logout
		 */
		logout: (params: RequestParams = {}) =>
			this.http.request<LogoutData, LogoutError>({
				path: `/api/v1/auth/logout`,
				method: "POST",
				type: ContentType.Json,
				format: "json",
				...params,
			}),

		/**
		 * @description Request Google OAuth
		 *
		 * @tags Auth
		 * @name RequestGoogleOauth
		 * @summary Request Google OAuth
		 * @request GET:/api/v1/auth/request-google-oauth
		 */
		requestGoogleOauth: (params: RequestParams = {}) =>
			this.http.request<any, RequestGoogleOauthError>({
				path: `/api/v1/auth/request-google-oauth`,
				method: "GET",
				type: ContentType.Json,
				...params,
			}),

		/**
		 * @description Verify Google OAuth code
		 *
		 * @tags Auth
		 * @name VerifyGoogleOauth
		 * @summary Verify Google OAuth code
		 * @request GET:/api/v1/auth/verify-google-oauth
		 */
		verifyGoogleOauth: (
			query: VerifyGoogleOauthParams,
			params: RequestParams = {}
		) =>
			this.http.request<any, VerifyGoogleOauthError>({
				path: `/api/v1/auth/verify-google-oauth`,
				method: "GET",
				query: query,
				type: ContentType.Json,
				...params,
			}),

		/**
		 * @description Check onboard status
		 *
		 * @tags Onboard
		 * @name CheckOnboardStatus
		 * @summary Check onboard status
		 * @request POST:/api/v1/onboard/check-onboard-status
		 */
		checkOnboardStatus: (
			message_signature: OnboardCheckOnboardStatusRequest,
			params: RequestParams = {}
		) =>
			this.http.request<CheckOnboardStatusData, CheckOnboardStatusError>({
				path: `/api/v1/onboard/check-onboard-status`,
				method: "POST",
				body: message_signature,
				type: ContentType.Json,
				format: "json",
				...params,
			}),

		/**
		 * @description Register a new user with Google OAuth
		 *
		 * @tags Onboard
		 * @name RegisterWithGoogleOauth
		 * @summary Register a new user with Google OAuth
		 * @request POST:/api/v1/onboard/register-with-google-oauth
		 */
		registerWithGoogleOauth: (
			password: OnboardRegisterWithGoogleOAuthRequest,
			params: RequestParams = {}
		) =>
			this.http.request<
				RegisterWithGoogleOauthData,
				RegisterWithGoogleOauthError
			>({
				path: `/api/v1/onboard/register-with-google-oauth`,
				method: "POST",
				body: password,
				type: ContentType.Json,
				format: "json",
				...params,
			}),

		/**
		 * @description Register a new user with wallet address
		 *
		 * @tags Onboard
		 * @name RegisterWithWallet
		 * @summary Register a new user with wallet address
		 * @request POST:/api/v1/onboard/register-with-wallet
		 */
		registerWithWallet: (
			signed_message: OnboardRegisterWithWalletRequest,
			params: RequestParams = {}
		) =>
			this.http.request<RegisterWithWalletData, RegisterWithWalletError>({
				path: `/api/v1/onboard/register-with-wallet`,
				method: "POST",
				body: signed_message,
				type: ContentType.Json,
				format: "json",
				...params,
			}),

		/**
		 * @description Retrieve preset message for the client to sign to register
		 *
		 * @tags Onboard
		 * @name GetSignMessage
		 * @summary Get preset message for the client to sign to register
		 * @request GET:/api/v1/onboard/sign-message
		 */
		getSignMessage: (params: RequestParams = {}) =>
			this.http.request<GetSignMessageData, any>({
				path: `/api/v1/onboard/sign-message`,
				method: "GET",
				format: "json",
				...params,
			}),

		/**
		 * @description Create a profile
		 *
		 * @tags Profile
		 * @name CreateProfile
		 * @summary Create a profile
		 * @request POST:/api/v1/profile
		 */
		createProfile: (
			profile: ProfileCreateProfileRequest,
			params: RequestParams = {}
		) =>
			this.http.request<CreateProfileData, CreateProfileError>({
				path: `/api/v1/profile`,
				method: "POST",
				body: profile,
				type: ContentType.Json,
				format: "json",
				...params,
			}),

		/**
		 * @description Update a profile by credential ID
		 *
		 * @tags Profile
		 * @name UpdateProfileByCredentialId
		 * @summary Update a profile by credential ID
		 * @request PATCH:/api/v1/profile/credential/{credential_id}
		 */
		updateProfileByCredentialId: (
			{ credentialId, ...query }: UpdateProfileByCredentialIdParams,
			AcademicEmail: ProfileUpdateProfileRequest,
			params: RequestParams = {}
		) =>
			this.http.request<
				UpdateProfileByCredentialIdData,
				UpdateProfileByCredentialIdError
			>({
				path: `/api/v1/profile/credential/${credentialId}`,
				method: "PATCH",
				body: AcademicEmail,
				type: ContentType.Json,
				format: "json",
				...params,
			}),

		/**
		 * @description Get my profile
		 *
		 * @tags Profile
		 * @name GetMyProfile
		 * @summary Get my profile
		 * @request GET:/api/v1/profile/my
		 */
		getMyProfile: (params: RequestParams = {}) =>
			this.http.request<GetMyProfileData, GetMyProfileError>({
				path: `/api/v1/profile/my`,
				method: "GET",
				type: ContentType.Json,
				format: "json",
				...params,
			}),
	};
}
