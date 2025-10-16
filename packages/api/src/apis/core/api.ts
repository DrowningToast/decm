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

export type CreateEventCertificateConfigData =
  EventConfigEventCertificateConfigResponse;

export type CreateEventCertificateConfigError = CustomerrorErrResponse;

export interface CreateEventCertificateConfigParams {
  /** Event ID */
  eventId: string;
}

export type CreateEventContractData = EventConfigEventContractResponse;

export type CreateEventContractError = CustomerrorErrResponse;

export interface CreateEventContractParams {
  /** Event ID */
  eventId: string;
}

export type CreateEventData = EntityEvent;

export type CreateEventError = CustomerrorErrResponse;

export type CreateEventIssuerData = EventConfigEventIssuerResponse;

export type CreateEventIssuerError = CustomerrorErrResponse;

export interface CreateEventIssuerParams {
  /** Event ID */
  eventId: string;
}

export interface CreateEventPayload {
  /**
   * Event banner image (JPEG, PNG, WebP, max 10MB)
   * @format binary
   */
  banner: File;
  /** Contact address */
  contact_address: string;
  /** Contact number */
  contact_number: string;
  /** Event description */
  description: string;
  /** End date (RFC3339 format) */
  end_date: string;
  /** Google map query */
  google_map_query: string;
  /**
   * Event icon image (JPEG, PNG, WebP, max 10MB)
   * @format binary
   */
  icon: File;
  /** Location */
  location: string;
  /** Event name */
  name: string;
  /** Number of seats */
  seats_count: number;
  /** Event short description */
  short_description: string;
  /** Start date (RFC3339 format) */
  start_date: string;
}

export type CreateEventRegistrationConfigData =
  EventConfigEventRegistrationConfigResponse;

export type CreateEventRegistrationConfigError = CustomerrorErrResponse;

export interface CreateEventRegistrationConfigParams {
  /** Event ID */
  eventId: string;
}

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

export type DeleteEventCertificateConfigData = Record<string, string>;

export type DeleteEventCertificateConfigError = CustomerrorErrResponse;

export interface DeleteEventCertificateConfigParams {
  /** Event ID */
  eventId: string;
}

export type DeleteEventContractData = Record<string, string>;

export type DeleteEventContractError = CustomerrorErrResponse;

export interface DeleteEventContractParams {
  /** Event ID */
  eventId: string;
}

export type DeleteEventIssuerData = Record<string, string>;

export type DeleteEventIssuerError = CustomerrorErrResponse;

export interface DeleteEventIssuerParams {
  /** Event ID */
  eventId: string;
  /** Issuer ID */
  issuerId: string;
}

export type DeleteEventRegistrationConfigData = Record<string, string>;

export type DeleteEventRegistrationConfigError = CustomerrorErrResponse;

export interface DeleteEventRegistrationConfigParams {
  /** Event ID */
  eventId: string;
}

export interface EntityEvent {
  event_banner_presigned_url?: string;
  event_name?: string;
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

export interface EventConfigCreateEventCertificateConfigRequest {
  academic_institution_pos_x?: number;
  academic_institution_pos_y?: number;
  base_certificate_storage_key?: string;
  event_name_pos_x?: number;
  event_name_pos_y?: number;
  name_pos_x?: number;
  name_pos_y?: number;
}

export interface EventConfigCreateEventContractRequest {
  access_manager_contract_address?: string;
  certificate_contract_address?: string;
  event_contract_address?: string;
  ticket_contract_address?: string;
}

export interface EventConfigCreateEventIssuerRequest {
  is_signed?: number;
  issuer_credential_id?: string;
  signature?: string;
}

export interface EventConfigCreateEventRegistrationConfigRequest {
  academic_email_requirement_status?: number;
  academic_institution_requirement_status?: number;
  address_requirement_status?: number;
  bio_requirement_status?: number;
  email_requirement_status?: number;
  first_name_requirement_status?: number;
  last_name_requirement_status?: number;
  phone_number_requirement_status?: number;
}

export interface EventConfigEventCertificateConfigResponse {
  academic_institution_pos_x?: number;
  academic_institution_pos_y?: number;
  base_certificate_storage_key?: string;
  created_at?: string;
  event_id?: string;
  event_name_pos_x?: number;
  event_name_pos_y?: number;
  id?: string;
  name_pos_x?: number;
  name_pos_y?: number;
  updated_at?: string;
}

export interface EventConfigEventContractResponse {
  access_manager_contract_address?: string;
  certificate_contract_address?: string;
  created_at?: string;
  event_contract_address?: string;
  event_id?: string;
  id?: string;
  ticket_contract_address?: string;
  updated_at?: string;
}

export interface EventConfigEventIssuerResponse {
  created_at?: string;
  event_id?: string;
  id?: string;
  is_signed?: number;
  issuer_credential_id?: string;
  signature?: string;
  updated_at?: string;
}

export interface EventConfigEventRegistrationConfigResponse {
  academic_email_requirement_status?: number;
  academic_institution_requirement_status?: number;
  address_requirement_status?: number;
  bio_requirement_status?: number;
  created_at?: string;
  email_requirement_status?: number;
  event_id?: string;
  first_name_requirement_status?: number;
  id?: string;
  last_name_requirement_status?: number;
  phone_number_requirement_status?: number;
  updated_at?: string;
}

export interface EventConfigUpdateEventCertificateConfigRequest {
  academic_institution_pos_x?: number;
  academic_institution_pos_y?: number;
  base_certificate_storage_key?: string;
  event_name_pos_x?: number;
  event_name_pos_y?: number;
  name_pos_x?: number;
  name_pos_y?: number;
}

export interface EventConfigUpdateEventContractRequest {
  access_manager_contract_address?: string;
  certificate_contract_address?: string;
  event_contract_address?: string;
  ticket_contract_address?: string;
}

export interface EventConfigUpdateEventIssuerRequest {
  is_signed?: number;
  signature?: string;
}

export interface EventConfigUpdateEventRegistrationConfigRequest {
  academic_email_requirement_status?: number;
  academic_institution_requirement_status?: number;
  address_requirement_status?: number;
  bio_requirement_status?: number;
  email_requirement_status?: number;
  first_name_requirement_status?: number;
  last_name_requirement_status?: number;
  phone_number_requirement_status?: number;
}

export type GetEventCertificateConfigData =
  EventConfigEventCertificateConfigResponse;

export type GetEventCertificateConfigError = CustomerrorErrResponse;

export interface GetEventCertificateConfigParams {
  /** Event ID */
  eventId: string;
}

export type GetEventContractByEventIdData = EventConfigEventContractResponse;

export type GetEventContractByEventIdError = CustomerrorErrResponse;

export interface GetEventContractByEventIdParams {
  /** Event ID */
  eventId: string;
}

export type GetEventIssuerByIdData = EventConfigEventIssuerResponse;

export type GetEventIssuerByIdError = CustomerrorErrResponse;

export interface GetEventIssuerByIdParams {
  /** Event ID */
  eventId: string;
  /** Issuer ID */
  issuerId: string;
}

export type GetEventIssuersByEventIdData = EventConfigEventIssuerResponse[];

export type GetEventIssuersByEventIdError = CustomerrorErrResponse;

export interface GetEventIssuersByEventIdParams {
  /** Event ID */
  eventId: string;
}

export type GetEventRegistrationConfigData =
  EventConfigEventRegistrationConfigResponse;

export type GetEventRegistrationConfigError = CustomerrorErrResponse;

export interface GetEventRegistrationConfigParams {
  /** Event ID */
  eventId: string;
}

export type GetMyProfileData = EntityProfile;

export type GetMyProfileError = CustomerrorErr;

export type GetRegisterSignMessageData = OnboardGetRegisterSignMessageResponse;

export type LogoutData = Record<string, string>;

export type LogoutError = CustomerrorErrResponse;

export interface OnboardCheckOnboardStatusRequest {
  access_token?: string;
  expires_in?: number;
  method?: "google" | "wallet";
  sign_message?: string;
}

export interface OnboardCheckOnboardStatusResponse {
  authentication_credential_id?: string;
  profile_id?: string;
}

/** Response for the client to sign to register */
export interface OnboardGetRegisterSignMessageResponse {
  message?: string;
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

export type UpdateEventCertificateConfigData =
  EventConfigEventCertificateConfigResponse;

export type UpdateEventCertificateConfigError = CustomerrorErrResponse;

export interface UpdateEventCertificateConfigParams {
  /** Event ID */
  eventId: string;
}

export type UpdateEventContractData = EventConfigEventContractResponse;

export type UpdateEventContractError = CustomerrorErrResponse;

export interface UpdateEventContractParams {
  /** Event ID */
  eventId: string;
}

export type UpdateEventIssuerData = EventConfigEventIssuerResponse;

export type UpdateEventIssuerError = CustomerrorErrResponse;

export interface UpdateEventIssuerParams {
  /** Event ID */
  eventId: string;
  /** Issuer ID */
  issuerId: string;
}

export type UpdateEventRegistrationConfigData =
  EventConfigEventRegistrationConfigResponse;

export type UpdateEventRegistrationConfigError = CustomerrorErrResponse;

export interface UpdateEventRegistrationConfigParams {
  /** Event ID */
  eventId: string;
}

export type UpdateProfileByCredentialIdData = ProfileUpdateProfileResponse;

export type UpdateProfileByCredentialIdError = CustomerrorErr;

export interface UpdateProfileByCredentialIdParams {
  /** Credential ID */
  credentialId: string;
}

export type VerifyGoogleOauthError = CustomerrorErrResponse;

export interface VerifyGoogleOauthParams {
  code: string;
  state: string;
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
      params: RequestParams = {},
    ) =>
      this.http.request<any, VerifyGoogleOauthError>({
        path: `/api/v1/auth/verify-google-oauth`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Create a new event with banner image upload
     *
     * @tags Event
     * @name CreateEvent
     * @summary Create a new event
     * @request POST:/api/v1/events
     */
    createEvent: (data: CreateEventPayload, params: RequestParams = {}) =>
      this.http.request<CreateEventData, CreateEventError>({
        path: `/api/v1/events`,
        method: "POST",
        body: data,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),

    /**
     * @description Get the event certificate configuration for an event
     *
     * @name GetEventCertificateConfig
     * @summary Get event certificate config
     * @request GET:/api/v1/events/{event_id}/certificate-config
     */
    getEventCertificateConfig: (
      { eventId, ...query }: GetEventCertificateConfigParams,
      params: RequestParams = {},
    ) =>
      this.http.request<
        GetEventCertificateConfigData,
        GetEventCertificateConfigError
      >({
        path: `/api/v1/events/${eventId}/certificate-config`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update the event certificate configuration for an event
     *
     * @name UpdateEventCertificateConfig
     * @summary Update event certificate config
     * @request PUT:/api/v1/events/{event_id}/certificate-config
     */
    updateEventCertificateConfig: (
      { eventId, ...query }: UpdateEventCertificateConfigParams,
      request: EventConfigUpdateEventCertificateConfigRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<
        UpdateEventCertificateConfigData,
        UpdateEventCertificateConfigError
      >({
        path: `/api/v1/events/${eventId}/certificate-config`,
        method: "PUT",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new event certificate configuration for an event
     *
     * @name CreateEventCertificateConfig
     * @summary Create event certificate config
     * @request POST:/api/v1/events/{event_id}/certificate-config
     */
    createEventCertificateConfig: (
      { eventId, ...query }: CreateEventCertificateConfigParams,
      request: EventConfigCreateEventCertificateConfigRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<
        CreateEventCertificateConfigData,
        CreateEventCertificateConfigError
      >({
        path: `/api/v1/events/${eventId}/certificate-config`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete the event certificate configuration for an event
     *
     * @name DeleteEventCertificateConfig
     * @summary Delete event certificate config
     * @request DELETE:/api/v1/events/{event_id}/certificate-config
     */
    deleteEventCertificateConfig: (
      { eventId, ...query }: DeleteEventCertificateConfigParams,
      params: RequestParams = {},
    ) =>
      this.http.request<
        DeleteEventCertificateConfigData,
        DeleteEventCertificateConfigError
      >({
        path: `/api/v1/events/${eventId}/certificate-config`,
        method: "DELETE",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get the event contract for an event
     *
     * @name GetEventContractByEventId
     * @summary Get event contract by event ID
     * @request GET:/api/v1/events/{event_id}/contracts
     */
    getEventContractByEventId: (
      { eventId, ...query }: GetEventContractByEventIdParams,
      params: RequestParams = {},
    ) =>
      this.http.request<
        GetEventContractByEventIdData,
        GetEventContractByEventIdError
      >({
        path: `/api/v1/events/${eventId}/contracts`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update the event contract for an event
     *
     * @name UpdateEventContract
     * @summary Update event contract
     * @request PUT:/api/v1/events/{event_id}/contracts
     */
    updateEventContract: (
      { eventId, ...query }: UpdateEventContractParams,
      request: EventConfigUpdateEventContractRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<UpdateEventContractData, UpdateEventContractError>({
        path: `/api/v1/events/${eventId}/contracts`,
        method: "PUT",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new event contract for an event
     *
     * @name CreateEventContract
     * @summary Create event contract
     * @request POST:/api/v1/events/{event_id}/contracts
     */
    createEventContract: (
      { eventId, ...query }: CreateEventContractParams,
      request: EventConfigCreateEventContractRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<CreateEventContractData, CreateEventContractError>({
        path: `/api/v1/events/${eventId}/contracts`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete the event contract for an event
     *
     * @name DeleteEventContract
     * @summary Delete event contract
     * @request DELETE:/api/v1/events/{event_id}/contracts
     */
    deleteEventContract: (
      { eventId, ...query }: DeleteEventContractParams,
      params: RequestParams = {},
    ) =>
      this.http.request<DeleteEventContractData, DeleteEventContractError>({
        path: `/api/v1/events/${eventId}/contracts`,
        method: "DELETE",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all event issuers for an event
     *
     * @name GetEventIssuersByEventId
     * @summary Get event issuers by event ID
     * @request GET:/api/v1/events/{event_id}/issuers
     */
    getEventIssuersByEventId: (
      { eventId, ...query }: GetEventIssuersByEventIdParams,
      params: RequestParams = {},
    ) =>
      this.http.request<
        GetEventIssuersByEventIdData,
        GetEventIssuersByEventIdError
      >({
        path: `/api/v1/events/${eventId}/issuers`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new event issuer for an event
     *
     * @name CreateEventIssuer
     * @summary Create event issuer
     * @request POST:/api/v1/events/{event_id}/issuers
     */
    createEventIssuer: (
      { eventId, ...query }: CreateEventIssuerParams,
      request: EventConfigCreateEventIssuerRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<CreateEventIssuerData, CreateEventIssuerError>({
        path: `/api/v1/events/${eventId}/issuers`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get an event issuer by its ID
     *
     * @name GetEventIssuerById
     * @summary Get event issuer by ID
     * @request GET:/api/v1/events/{event_id}/issuers/{issuer_id}
     */
    getEventIssuerById: (
      { eventId, issuerId, ...query }: GetEventIssuerByIdParams,
      params: RequestParams = {},
    ) =>
      this.http.request<GetEventIssuerByIdData, GetEventIssuerByIdError>({
        path: `/api/v1/events/${eventId}/issuers/${issuerId}`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update an event issuer
     *
     * @name UpdateEventIssuer
     * @summary Update event issuer
     * @request PUT:/api/v1/events/{event_id}/issuers/{issuer_id}
     */
    updateEventIssuer: (
      { eventId, issuerId, ...query }: UpdateEventIssuerParams,
      request: EventConfigUpdateEventIssuerRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<UpdateEventIssuerData, UpdateEventIssuerError>({
        path: `/api/v1/events/${eventId}/issuers/${issuerId}`,
        method: "PUT",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete an event issuer
     *
     * @name DeleteEventIssuer
     * @summary Delete event issuer
     * @request DELETE:/api/v1/events/{event_id}/issuers/{issuer_id}
     */
    deleteEventIssuer: (
      { eventId, issuerId, ...query }: DeleteEventIssuerParams,
      params: RequestParams = {},
    ) =>
      this.http.request<DeleteEventIssuerData, DeleteEventIssuerError>({
        path: `/api/v1/events/${eventId}/issuers/${issuerId}`,
        method: "DELETE",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get the event registration configuration for an event
     *
     * @name GetEventRegistrationConfig
     * @summary Get event registration config
     * @request GET:/api/v1/events/{event_id}/registration-config
     */
    getEventRegistrationConfig: (
      { eventId, ...query }: GetEventRegistrationConfigParams,
      params: RequestParams = {},
    ) =>
      this.http.request<
        GetEventRegistrationConfigData,
        GetEventRegistrationConfigError
      >({
        path: `/api/v1/events/${eventId}/registration-config`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update the event registration configuration for an event
     *
     * @name UpdateEventRegistrationConfig
     * @summary Update event registration config
     * @request PUT:/api/v1/events/{event_id}/registration-config
     */
    updateEventRegistrationConfig: (
      { eventId, ...query }: UpdateEventRegistrationConfigParams,
      request: EventConfigUpdateEventRegistrationConfigRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<
        UpdateEventRegistrationConfigData,
        UpdateEventRegistrationConfigError
      >({
        path: `/api/v1/events/${eventId}/registration-config`,
        method: "PUT",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new event registration configuration for an event
     *
     * @name CreateEventRegistrationConfig
     * @summary Create event registration config
     * @request POST:/api/v1/events/{event_id}/registration-config
     */
    createEventRegistrationConfig: (
      { eventId, ...query }: CreateEventRegistrationConfigParams,
      request: EventConfigCreateEventRegistrationConfigRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<
        CreateEventRegistrationConfigData,
        CreateEventRegistrationConfigError
      >({
        path: `/api/v1/events/${eventId}/registration-config`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete the event registration configuration for an event
     *
     * @name DeleteEventRegistrationConfig
     * @summary Delete event registration config
     * @request DELETE:/api/v1/events/{event_id}/registration-config
     */
    deleteEventRegistrationConfig: (
      { eventId, ...query }: DeleteEventRegistrationConfigParams,
      params: RequestParams = {},
    ) =>
      this.http.request<
        DeleteEventRegistrationConfigData,
        DeleteEventRegistrationConfigError
      >({
        path: `/api/v1/events/${eventId}/registration-config`,
        method: "DELETE",
        type: ContentType.Json,
        format: "json",
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
      sign_message: OnboardCheckOnboardStatusRequest,
      params: RequestParams = {},
    ) =>
      this.http.request<CheckOnboardStatusData, CheckOnboardStatusError>({
        path: `/api/v1/onboard/check-onboard-status`,
        method: "POST",
        body: sign_message,
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
      signed_message: OnboardRegisterWithWalletRequest,
      params: RequestParams = {},
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
     * @name CreateProfile
     * @summary Create a profile
     * @request POST:/api/v1/profile
     */
    createProfile: (
      profile: ProfileCreateProfileRequest,
      params: RequestParams = {},
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
      params: RequestParams = {},
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
