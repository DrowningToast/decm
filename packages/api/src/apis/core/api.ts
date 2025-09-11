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

export interface AuthVerifyGoogleOAuthResponse {
  access_token?: string;
  expires_in?: number;
  refresh_token?: string;
}

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

export type GetRegisterSignMessageData = OnboardGetRegisterSignMessageResponse;

/** Response for the client to sign to register */
export interface OnboardGetRegisterSignMessageResponse {
  message?: string;
}

export interface OnboardRegisterWithGoogleOAuthResponse {
  mnemonic?: string[];
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

export type RegisterWithGoogleOauthData =
  OnboardRegisterWithGoogleOAuthResponse;

export type RegisterWithGoogleOauthError = CustomerrorErrResponse;

/** Access token */
export type RegisterWithGoogleOauthPayload = string;

export type RegisterWithWalletData = any;

export type RegisterWithWalletError = CustomerrorErrResponse;

/** Signed message */
export type RegisterWithWalletPayload = string;

export type RequestGoogleOauthError = CustomerrorErrResponse;

export type V1ProfileCreateData = ProfileCreateProfileResponse;

export type V1ProfileCreateError = CustomerrorErr;

export type V1ProfileCredentialPartialUpdateData = ProfileUpdateProfileResponse;

export type V1ProfileCredentialPartialUpdateError = CustomerrorErr;

export interface V1ProfileCredentialPartialUpdateParams {
  /** Credential ID */
  credentialId: string;
}

export type V1ProfileMyListData = EntityProfile;

export type V1ProfileMyListError = CustomerrorErr;

export type VerifyGoogleOauthData = AuthVerifyGoogleOAuthResponse;

export type VerifyGoogleOauthError = CustomerrorErrResponse;

/** Code */
export type VerifyGoogleOauthPayload = string;

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
    securityData: SecurityDataType | null,
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
    params2?: AxiosRequestConfig,
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
          isFileType ? formItem : this.stringifyFormItem(formItem),
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
     * @request POST:/api/v1/auth/verify-google-oauth
     */
    verifyGoogleOauth: (
      state: VerifyGoogleOauthPayload,
      params: RequestParams = {},
    ) =>
      this.http.request<VerifyGoogleOauthData, VerifyGoogleOauthError>({
        path: `/api/v1/auth/verify-google-oauth`,
        method: "POST",
        body: state,
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
      password: RegisterWithGoogleOauthPayload,
      params: RequestParams = {},
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
      wallet_address: RegisterWithWalletPayload,
      params: RequestParams = {},
    ) =>
      this.http.request<RegisterWithWalletData, RegisterWithWalletError>({
        path: `/api/v1/onboard/register-with-wallet`,
        method: "POST",
        body: wallet_address,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Retrieve preset message for the client to sign to register
     *
     * @name GetRegisterSignMessage
     * @summary Get preset message for the client to sign to register
     * @request GET:/api/v1/onboard/sign-message
     */
    getRegisterSignMessage: (params: RequestParams = {}) =>
      this.http.request<GetRegisterSignMessageData, any>({
        path: `/api/v1/onboard/sign-message`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Create a profile
     *
     * @tags Profile
     * @name V1ProfileCreate
     * @summary Create a profile
     * @request POST:/api/v1/profile
     */
    v1ProfileCreate: (
      profile: ProfileCreateProfileRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<V1ProfileCreateData, V1ProfileCreateError>({
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
     * @name V1ProfileCredentialPartialUpdate
     * @summary Update a profile by credential ID
     * @request PATCH:/api/v1/profile/credential/{credential_id}
     */
    v1ProfileCredentialPartialUpdate: (
      { credentialId, ...query }: V1ProfileCredentialPartialUpdateParams,
      profile: ProfileUpdateProfileRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<
        V1ProfileCredentialPartialUpdateData,
        V1ProfileCredentialPartialUpdateError
      >({
        path: `/api/v1/profile/credential/${credentialId}`,
        method: "PATCH",
        body: profile,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get my profile
     *
     * @tags Profile
     * @name V1ProfileMyList
     * @summary Get my profile
     * @request GET:/api/v1/profile/my
     */
    v1ProfileMyList: (params: RequestParams = {}) =>
      this.http.request<V1ProfileMyListData, V1ProfileMyListError>({
        path: `/api/v1/profile/my`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
