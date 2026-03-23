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

export interface BlockchainGasPriceResponse {
    base_fee_gwei: number;
    current_gas_price_gwei: number;
    hard_cap_price_gwei: number;
    /** Percentage remaining until hard cap */
    hard_safety_margin: number;
    priority_fee_gwei: number;
    soft_cap_price_gwei: number;
    /** Percentage remaining until soft cap */
    soft_safety_margin: number;
}

export interface BlockchainGetSystemStatusResponse {
    closest_incoming_schedule?: EntitySystemStatusSchedule | null;
    gas_price: BlockchainGasPriceResponse;
}

export type CancelEventRegistrationInvitationData =
    EventRegistrationEventRegistrationInvitationResponse;

export type CancelEventRegistrationInvitationError = CustomerrorErrResponse;

export interface CancelEventRegistrationInvitationParams {
    /** Event Registration Invitation ID */
    event_registration_invitation_id: string;
}

export interface CertificateClaimCertificateBody {
    account_password?: string;
    certificate_password?: string;
    sign_message?: string;
    signature?: string;
}

export interface CertificateClaimCertificateResponse {
    certificate: any;
    message: string;
    status: string;
    user_signature: any;
}

export interface CertificateGetClaimCertificateSignMessageResponse {
    sign_message: string;
}

export interface CertificateGetMyCertificatesListViewModelResponse {
    claimed_certificates: EventClaimedCertificateViewModel[];
    total_claimed: number;
    total_unclaimed: number;
    unclaimed_certificates: EventUnclaimedCertificateViewModel[];
}

export interface CertificateShareCertificateShareViewModel {
    active: boolean;
    created_at: string;
    event_certificate_id: string;
    handle: string;
    has_password: boolean;
    id: string;
    updated_at: string;
}

export interface CertificateShareHandlerCertificateShareContractInfo {
    certificateTokenId: string;
    eventCertificateContractAddress: string;
}

export interface CertificateShareHandlerCertificateShareDataResponse {
    contract: CertificateShareHandlerCertificateShareContractInfo;
    data: EntityCertificatePayload;
    decryptedCertificateData?: EntityCertificateRawData;
    decryptedUserData?: EntityAttendeeProfileData;
}

export interface CertificateShareHandlerCertificateShareViewModel {
    active: boolean;
    created_at: string;
    event_certificate_id: string;
    handle: string;
    has_password: boolean;
    id: string;
    updated_at: string;
}

export interface CertificateShareHandlerCreateCertificateShareBody {
    password?: string;
}

export interface CertificateShareHandlerCreateCertificateShareResponse {
    share: CertificateShareHandlerCertificateShareViewModel;
}

export interface CertificateShareHandlerGetCertificateShareBody {
    password?: string;
}

export interface CertificateShareHandlerGetCertificateShareImageBody {
    password?: string;
}

export interface CertificateShareHandlerUpdateCertificateShareBody {
    active?: boolean;
    password?: string;
}

export interface CertificateShareHandlerUpdateCertificateShareResponse {
    share: CertificateShareHandlerCertificateShareViewModel;
}

export type CheckCertificateMintReadinessData =
    CoreApiInternalHandlerEventconfigCertificateMintReadinessResponse;

export type CheckCertificateMintReadinessError = CustomerrorErrResponse;

export interface CheckCertificateMintReadinessParams {
    /** Event ID */
    eventId: string;
}

export type CheckEventPasswordData = EventconfigCheckEventPasswordResponse;

export type CheckEventPasswordError = CustomerrorErrResponse;

export interface CheckEventPasswordParams {
    /** Event ID */
    eventId: string;
}

export type CheckOnboardStatusData = OnboardCheckOnboardStatusResponse;

export type CheckOnboardStatusError = CustomerrorErrResponse;

export type CheckRoleData = CheckRoleResponse;

export type CheckRoleError = CustomerrorErrResponse;

export interface CheckRoleParams {
    /** Check if user is authenticated */
    is_authenticated?: boolean;
    /** Check if user is a verified host/organizer */
    is_host?: boolean;
    /** Check if user is a verified issuer */
    is_issuer?: boolean;
}

/** Role verification response */
export interface CheckRoleResponse {
    is_authenticated?: boolean;
    is_host?: boolean;
    is_issuer?: boolean;
}

export type ClaimCertificateData = CertificateClaimCertificateResponse;

export type ClaimCertificateError = CustomerrorErrResponse;

export interface ClaimCertificateParams {
    /** Certificate ID */
    certificateId: string;
}

/** @format int32 */
export enum CommonSolutionStatus {
    SolutionStatusManaged = 0,
    SolutionStatusBYOK = 1,
}

export interface CoreApiInternalHandlerEventCertificateSignature {
    certificate?: EntityEventCertificate;
    signature: string;
}

export interface CoreApiInternalHandlerEventEventContractResponse {
    access_manager_contract_address: string;
    certificate_contract_address?: string;
    created_at: string;
    event_contract_address: string;
    event_id: string;
    id: string;
    ticket_contract_address?: string;
    updated_at: string;
}

export interface CoreApiInternalHandlerEventEventResponse {
    attendees_count: number;
    banner_storage_key: string;
    chain_id: number;
    contact_address: string;
    contact_number: string;
    created_at: string;
    end_date: string;
    event_status: EntityEventStatus;
    event_type: EntityEventType;
    google_map_query: string;
    icon_storage_key: string;
    id: string;
    is_booking_request_required: boolean;
    is_public: boolean;
    is_ticket_transferable: boolean;
    is_verified: boolean;
    location: string;
    long_description: string;
    max_attendees: number;
    owner_credential_id: string;
    short_description: string;
    start_date: string;
    title: string;
    updated_at: string;
}

export interface CoreApiInternalHandlerEventImportCertificateReceiversRequest {
    event_id: string;
    host_pin: string;
    /** @minItems 1 */
    receivers: EventImportCertificateReceiverRequest[];
}

export interface CoreApiInternalHandlerEventImportCertificateReceiversResponse {
    certificates: EntityEventCertificate[];
    event_certificate_address: string;
    event_id: string;
}

export interface CoreApiInternalHandlerEventPublishEventCertificatesResponse {
    event_id: string;
    inbox_messages_created: number;
    is_published: boolean;
    published_count: number;
}

export interface CoreApiInternalHandlerEventRegistrationJoinEventPayload {
    academic_email?: string;
    academic_institution?: string;
    address?: string;
    bio?: string;
    email?: string;
    first_name?: string;
    last_name?: string;
    phone_number?: string;
}

export interface CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse {
    revoked_certificates: EntityEventCertificate[];
}

export interface CoreApiInternalHandlerEventRevokeEventCertificatesRequest {
    /** @minItems 1 */
    certificate_ids: string[];
}

export interface CoreApiInternalHandlerEventRevokeEventCertificatesResponse {
    revoked_certificates: EntityEventCertificate[];
}

export interface CoreApiInternalHandlerEventSignEventCertificatesRequest {
    issuer_pin: string;
}

export interface CoreApiInternalHandlerEventSignEventCertificatesResponse {
    certificates: CoreApiInternalHandlerEventCertificateSignature[];
}

export interface CoreApiInternalHandlerEventconfigCertificateMintReadinessResponse {
    /** @example true */
    all_issuers_have_signed: boolean;
    /** @example "0x1234567890abcdef" */
    certificate_contract_address?: string;
    /** @example true */
    has_certificate_config: boolean;
    /** @example true */
    has_certificate_contract: boolean;
    /** @example true */
    is_ready: boolean;
    missing_requirements?: string[];
    /** @example 2 */
    signed_issuers_count: number;
    /** @example 2 */
    total_issuers_count: number;
}

export interface CoreApiInternalHandlerEventconfigEventCertificateConfigResponse {
    academic_institution_font_family_id?: number;
    academic_institution_font_weight?: number;
    academic_institution_pos_x?: number;
    academic_institution_pos_y?: number;
    base_certificate_presigned_url: string;
    base_certificate_storage_key: string;
    certificate_subtitle_font_family_id?: number;
    certificate_subtitle_font_weight?: number;
    certificate_subtitle_pos_x?: number;
    certificate_subtitle_pos_y?: number;
    certificate_title_font_family_id?: number;
    certificate_title_font_weight?: number;
    certificate_title_pos_x?: number;
    certificate_title_pos_y?: number;
    created_at: string;
    event_id: string;
    event_name_font_family_id?: number;
    event_name_font_weight?: number;
    event_name_pos_x?: number;
    event_name_pos_y?: number;
    id: string;
    is_published: boolean;
    mint_readiness?: CoreApiInternalHandlerEventconfigMintReadinessInfo;
    name_font_family_id?: number;
    name_font_weight?: number;
    name_pos_x: number;
    name_pos_y: number;
    updated_at: string;
}

export interface CoreApiInternalHandlerEventconfigMintReadinessInfo {
    all_issuers_have_signed: boolean;
    certificate_contract_address?: string;
    has_certificate_config: boolean;
    has_certificate_contract: boolean;
    is_ready: boolean;
    missing_requirements?: string[];
    signed_issuers_count: number;
    total_issuers_count: number;
}

export interface CoreApiInternalUsecaseEventEventContractResponse {
    access_manager_contract_address: string;
    certificate_contract_address?: string;
    event_contract_address: string;
    event_id: string;
    id: string;
    ticket_contract_address?: string;
}

export type CreateCertificateShareData = CertificateShareHandlerCreateCertificateShareResponse;

export type CreateCertificateShareError = CustomerrorErrResponse;

export interface CreateCertificateShareParams {
    /** Certificate ID */
    certificateId: string;
}

export type CreateEventContractData = CoreApiInternalHandlerEventEventContractResponse;

export type CreateEventContractError = CustomerrorErrResponse;

export interface CreateEventContractParams {
    /** Event ID */
    eventId: string;
}

export type CreateEventData = CoreApiInternalHandlerEventEventResponse;

export type CreateEventError = CustomerrorErrResponse;

export type CreateEventIssuerData = EventEventIssuerResponse;

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
    /** Host password */
    host_password: string;
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

export type CreateEventRegistrationConfigData = EventconfigEventRegistrationConfigViewModel;

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
    message: string;
}

export type DeleteEventByIdData = CoreApiInternalHandlerEventEventResponse;

export type DeleteEventByIdError = CustomerrorErrResponse;

export interface DeleteEventByIdParams {
    /** Event ID */
    eventId: string;
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

export interface EntityAttendeeProfileData {
    academic_email: string;
    academic_institution: string;
    address: string;
    bio: string;
    email: string;
    first_name: string;
    last_name: string;
    phone_number: string;
}

export interface EntityCertificatePayload {
    data: EntityCertificatePayloadData;
    header: EntityCertificatePayloadHeader;
    proof: EntityCertificatePayloadProof;
}

export interface EntityCertificatePayloadData {
    backendEncryptedUserData: string;
    certificateId: string;
    certificateSubtitle: string;
    certificateTitle: string;
    certificateTokenId: string;
    encryptedUserData: string;
    eventDescription: string;
    eventName: string;
    issuedAt: string;
    issuerAddresses: string;
    issuerId: string;
    receiverAddress: string;
    /** "VALID" | "REVOKED" */
    status: string;
    userId: string;
}

export interface EntityCertificatePayloadHeader {
    "@context": string[];
    id: string;
    issuanceDate: string;
    issuer: string;
    type: string[];
}

export interface EntityCertificatePayloadHostProof {
    publicKey: string;
    signature: string;
}

export interface EntityCertificatePayloadIssuerProof {
    issuerPublicKey: string;
    issuerSignature: string;
}

export interface EntityCertificatePayloadProof {
    encryptedByBackendRawData: string;
    encryptedByUserRawData: string;
    hash: string;
    host: EntityCertificatePayloadHostProof;
    issuers: EntityCertificatePayloadIssuerProof[];
    signMessage: string;
}

export interface EntityCertificateRawData {
    academic_institution: string;
    certificate_subtitle: string;
    certificate_title: string;
    name: string;
}

export interface EntityEvent {
    attendees_count: number;
    banner_storage_key: string;
    chain_id: number;
    contact_address: string;
    contact_number: string;
    created_at: string;
    end_date: string;
    event_status: EntityEventStatus;
    event_type: EntityEventType;
    google_map_query: string;
    icon_storage_key: string;
    id: string;
    is_booking_request_required: boolean;
    is_public: boolean;
    is_ticket_transferable: boolean;
    is_verified: boolean;
    location: string;
    long_description: string;
    max_attendees: number;
    owner_credential_id: string;
    short_description: string;
    start_date: string;
    title: string;
    updated_at: string;
}

export interface EntityEventCertificate {
    aborted_at?: string;
    academic_institution?: string;
    broadcasted_at?: string;
    certificate_digest?: string;
    certificate_subtitle?: string;
    certificate_title?: string;
    certificate_token_id?: string;
    created_at: string;
    estimated_deadline?: string;
    event_certificate_address?: string;
    event_contract_address: string;
    event_id: string;
    event_name?: string;
    id: string;
    inbox_message_id?: string;
    name?: string;
    receiver_credential_id?: string;
    receiver_email?: string;
    revoked_at?: string;
    signature_created_at?: string;
    user_claim_signature_id?: string;
}

export interface EntityEventRegistrationInvitation {
    academic_institution?: string;
    accepted_at: string;
    cancelled_at?: string;
    code?: string;
    created_at: string;
    email?: string;
    event_id: string;
    first_name?: string;
    id: string;
    inbox_message_id: string;
    last_name?: string;
    phone_number?: string;
    updated_at: string;
    valid_until?: string;
}

export enum EntityEventStatus {
    EventStatusActive = "active",
    EventStatusInactive = "inactive",
    EventStatusClosed = "closed",
}

export enum EntityEventType {
    EventTypePrivate = "private",
    EventTypeInvite = "invite",
}

export interface EntityInboxMessage {
    created_at: string;
    deleted_at?: string;
    fallback_message_content?: string;
    hidden_at?: string;
    id: string;
    is_read: number;
    message_content: string;
    message_type: EntityInboxMessageType;
    receiver_credential_id?: string;
    receiver_email?: string;
    receiver_wallet_address?: string;
    sender_credential_id?: string;
    updated_at: string;
}

export enum EntityInboxMessageType {
    InboxMessageTypeGeneral = 0,
    InboxMessageTypeEventRegistrationInvitation = 1,
    InboxMessageTypeEventCertificateInvitation = 2,
}

export interface EntityProfile {
    academic_email?: string;
    academic_institution?: string;
    address?: string;
    authentication_credential_id: string;
    bio?: string;
    created_at: string;
    email?: string;
    first_name?: string;
    github_connector_ref?: string;
    /** Connector references (from authentication_credentials table) */
    google_connector_ref?: string;
    id: string;
    is_academic_email_public: boolean;
    is_academic_institution_public: boolean;
    is_address_public: boolean;
    is_bio_public: boolean;
    is_email_public: boolean;
    is_first_name_public: boolean;
    is_last_name_public: boolean;
    is_phone_number_public: boolean;
    is_profile_picture_public: boolean;
    last_name?: string;
    phone_number?: string;
    profile_picture_url?: string;
    updated_at: string;
    wallet_address?: string;
}

export enum EntitySystemStatus {
    SystemStatusMaintenance = 0,
    SystemStatusOperating = 1,
}

export interface EntitySystemStatusSchedule {
    created_at: string;
    deleted_at?: string | null;
    id: number;
    is_planned: boolean;
    order_id: number;
    planned_end_time?: string | null;
    start_time: string;
    status: EntitySystemStatus;
    updated_at: string;
}

export interface EventClaimedCertificateViewModel {
    aborted_at?: string;
    academic_institution?: string;
    broadcasted_at?: string;
    certificate_digest?: string;
    certificate_share?: CertificateShareCertificateShareViewModel;
    certificate_subtitle?: string;
    certificate_title?: string;
    certificate_token_id?: string;
    created_at: string;
    estimated_deadline?: string;
    event_certificate_address?: string;
    event_contract_address: string;
    event_id: string;
    event_name?: string;
    id: string;
    inbox_message_id?: string;
    name?: string;
    receiver_credential_id?: string;
    receiver_email?: string;
    revoked_at?: string;
    signature_created_at?: string;
    status: string;
    user_claim_signature_id?: string;
}

export interface EventCreateEventContractRequest {
    access_manager_contract_address: string;
    certificate_contract_address: string;
    event_contract_address: string;
    ticket_contract_address: string;
}

export interface EventCreateEventIssuerRequest {
    is_signed: 0 | 1;
    issuer_credential_id: string;
}

export interface EventDeleteEventRequest {
    host_password: string;
}

export interface EventEventIssuerResponse {
    created_at: string;
    event_id: string;
    id: string;
    is_signed: number;
    issuer_credential_id: string;
    issuer_profile?: EntityProfile;
    updated_at: string;
}

export interface EventEventParticipantResponse {
    academic_institution?: string;
    attendee_credential_id: string;
    contract_address: string;
    created_at: string;
    email?: string;
    event_id: string;
    first_name?: string;
    id: string;
    is_attendee_accepted: boolean;
    last_name?: string;
    phone_number?: string;
    updated_at: string;
    wallet_address: string;
}

export interface EventEventViewModel {
    attendees_count: number;
    banner_presigned_url: string;
    broadcasted_at?: string;
    broadcasted_at_deadline_block?: number;
    broadcasted_at_estimated_deadline?: string;
    chain_id: number;
    contact_address: string;
    contact_number: string;
    created_at: string;
    end_date: string;
    event_contract: CoreApiInternalUsecaseEventEventContractResponse;
    event_status: EntityEventStatus;
    event_type: EntityEventType;
    google_map_query: string;
    icon_presigned_url: string;
    id: string;
    is_booking_request_required: boolean;
    is_full: boolean;
    is_invited: boolean;
    is_joined?: boolean;
    is_public: boolean;
    is_ticket_transferable: boolean;
    is_verified: boolean;
    joined_at?: string;
    joined_is_accepted?: boolean;
    joined_sign_message?: string;
    /**
     * JoinedTxBroadcasted            *bool      `json:"joined_tx_broadcasted,omitempty"`
     * JoinedBroadcastedAt            *time.Time `json:"joined_broadcasted_at,omitempty"`
     */
    joined_signature?: string;
    location: string;
    long_description: string;
    max_attendees: number;
    owner_credential_id: string;
    registration_config: EventRegistrationConfigResponse;
    short_description: string;
    signature_expired_at?: string;
    start_date: string;
    title: string;
    updated_at: string;
}

export interface EventGetCertificatesListViewModelResponse {
    claimed_certificates: any;
    total_claimed: number;
    total_unclaimed: number;
    unclaimed_certificates: any;
}

export interface EventGetEventCertificatesResponse {
    certificates: EntityEventCertificate[];
}

export interface EventGetEventListResponse {
    events: EntityEvent[];
}

export interface EventImportCertificateReceiverRequest {
    academic_institution: string;
    certificate_subtitle: string;
    certificate_title: string;
    email: string;
    first_name: string;
    last_name: string;
}

export interface EventRegistrationConfigResponse {
    academic_email_requirement_status: number;
    academic_institution_requirement_status: number;
    address_requirement_status: number;
    bio_requirement_status: number;
    email_requirement_status: number;
    event_id: string;
    final_call_for_registration?: string;
    first_name_requirement_status: number;
    id: string;
    is_identity_verification_required: boolean;
    is_registration_password_required: boolean;
    last_name_requirement_status: number;
    phone_number_requirement_status: number;
}

export interface EventRegistrationEventAttendeeResponse {
    academic_email?: string;
    academic_institution?: string;
    address?: string;
    attendee_credential_id: string;
    bio?: string;
    contract_address: string;
    created_at: string;
    email?: string;
    event_id: string;
    first_name?: string;
    id: string;
    is_attendee_accepted: boolean;
    last_name?: string;
    phone_number?: string;
    updated_at: string;
    user_signature_id?: string;
    wallet_address: string;
}

export interface EventRegistrationEventRegistrationInvitationResponse {
    academic_institution?: string;
    accepted_at: string;
    cancelled_at?: string;
    code?: string;
    created_at: string;
    email?: string;
    event_id: string;
    first_name?: string;
    id: string;
    inbox_message_id: string;
    last_name?: string;
    phone_number?: string;
    updated_at: string;
    valid_until?: string;
}

export interface EventRegistrationGetEventRegistrationInvitationByUserAndEventResponse {
    inbox?: EntityInboxMessage;
    registration_invitation?: EntityEventRegistrationInvitation;
}

export interface EventRegistrationGetJoinEventSignMessageResponse {
    sign_message: string;
}

export interface EventRegistrationImportEventParticipantsRequest {
    event_id: string;
    /** @minItems 1 */
    participants: EventRegistrationParticipantRequestItem[];
}

export interface EventRegistrationJoinEventBody {
    account_password?: string;
    event_password?: string;
    registration_data: CoreApiInternalHandlerEventRegistrationJoinEventPayload;
    sign_message?: string;
    signature?: string;
}

export interface EventRegistrationParticipantRequestItem {
    academic_institution?: string;
    email: string;
    first_name: string;
    last_name: string;
    phone_number?: string;
    wallet_address?: string;
}

export interface EventUnclaimedCertificateViewModel {
    aborted_at?: string;
    academic_institution?: string;
    broadcasted_at?: string;
    certificate_digest?: string;
    certificate_share?: CertificateShareCertificateShareViewModel;
    certificate_subtitle?: string;
    certificate_title?: string;
    certificate_token_id?: string;
    created_at: string;
    estimated_deadline?: string;
    event_certificate_address?: string;
    event_contract_address: string;
    event_id: string;
    event_name?: string;
    id: string;
    inbox_message_id?: string;
    name?: string;
    receiver_credential_id?: string;
    receiver_email?: string;
    revoked_at?: string;
    signature_created_at?: string;
    user_claim_signature_id?: string;
}

export interface EventUpdateEventCertificateTextConfigRequest {
    academic_institution_font_family_id: number;
    academic_institution_font_weight: number;
    certificate_subtitle_font_family_id: number;
    certificate_subtitle_font_weight: number;
    certificate_title_font_family_id: number;
    certificate_title_font_weight: number;
    event_name_font_family_id: number;
    event_name_font_weight: number;
    name_font_family_id: number;
    name_font_weight: number;
}

export interface EventUpdateEventCertificateTextConfigResponse {
    academic_institution_font_family_id?: number;
    academic_institution_font_weight?: number;
    academic_institution_pos_x?: number;
    academic_institution_pos_y?: number;
    base_certificate_presigned_url: string;
    base_certificate_storage_key: string;
    certificate_subtitle_font_family_id?: number;
    certificate_subtitle_font_weight?: number;
    certificate_subtitle_pos_x?: number;
    certificate_subtitle_pos_y?: number;
    certificate_title_font_family_id?: number;
    certificate_title_font_weight?: number;
    certificate_title_pos_x?: number;
    certificate_title_pos_y?: number;
    created_at: string;
    event_id: string;
    event_name_font_family_id?: number;
    event_name_font_weight?: number;
    event_name_pos_x: number;
    event_name_pos_y: number;
    id: string;
    is_published: boolean;
    name_font_family_id?: number;
    name_font_weight?: number;
    name_pos_x: number;
    name_pos_y: number;
    updated_at: string;
}

export interface EventUpdateEventContractRequest {
    access_manager_contract_address: string;
    certificate_contract_address: string;
    event_contract_address: string;
    ticket_contract_address: string;
}

export interface EventUpdateEventIssuerRequest {
    event_id: string;
    issuer_credential_id: string;
}

export interface EventconfigCheckEventPasswordBody {
    password: string;
}

export interface EventconfigCheckEventPasswordResponse {
    is_valid: boolean;
}

export interface EventconfigCreateEventRegistrationConfigRequest {
    academic_email_requirement_status: 0 | 1 | 2;
    academic_institution_requirement_status: 0 | 1 | 2;
    address_requirement_status: 0 | 1 | 2;
    bio_requirement_status: 0 | 1 | 2;
    email_requirement_status: 0 | 1 | 2;
    final_call_for_registration?: string;
    first_name_requirement_status: 0 | 1 | 2;
    last_name_requirement_status: 0 | 1 | 2;
    phone_number_requirement_status: 0 | 1 | 2;
    /** @minLength 8 */
    registration_password?: string;
}

/** @format int32 */
export enum EventconfigEventRegistrationConfigRequirementStatus {
    EventRegistrationConfigRequirementStatusNotRequired = 0,
    EventRegistrationConfigRequirementStatusRequired = 1,
    EventRegistrationConfigRequirementStatusOptional = 2,
}

export interface EventconfigEventRegistrationConfigViewModel {
    academic_email_requirement_status: EventconfigEventRegistrationConfigRequirementStatus;
    academic_institution_requirement_status: EventconfigEventRegistrationConfigRequirementStatus;
    address_requirement_status: EventconfigEventRegistrationConfigRequirementStatus;
    bio_requirement_status: EventconfigEventRegistrationConfigRequirementStatus;
    created_at: string;
    email_requirement_status: EventconfigEventRegistrationConfigRequirementStatus;
    event_id: string;
    final_call_for_registration?: string;
    first_name_requirement_status: EventconfigEventRegistrationConfigRequirementStatus;
    id: string;
    last_name_requirement_status: EventconfigEventRegistrationConfigRequirementStatus;
    phone_number_requirement_status: EventconfigEventRegistrationConfigRequirementStatus;
    updated_at: string;
}

export interface EventconfigFontFamilyItem {
    available_font_weights: number[];
    css_font_name: string;
    font_family_name: string;
    id: number;
    is_default: boolean;
    is_support_italic: boolean;
}

export interface EventconfigGetEventCertificateFontFamiliesResponse {
    font_families: EventconfigFontFamilyItem[];
}

export interface EventconfigToggleCertificatePublishedRequest {
    is_published: boolean;
}

export interface EventconfigUpdateEventRegistrationConfigRequest {
    academic_email_requirement_status: number;
    academic_institution_requirement_status: number;
    address_requirement_status: number;
    bio_requirement_status: number;
    email_requirement_status: number;
    event_type: EntityEventType;
    final_call_for_registration?: string;
    first_name_requirement_status: number;
    is_booking_request_required: boolean;
    is_ticket_transferable: boolean;
    last_name_requirement_status: number;
    phone_number_requirement_status: number;
    registration_password?: string;
}

/** @format binary */
export type GenerateCertificateImageData = File;

export type GenerateCertificateImageError = CustomerrorErrResponse;

export interface GenerateCertificateImageParams {
    /**
     * Certificate ID
     * @format uuid
     */
    certificateId: string;
}

export type GetCertificateShareDataData = CertificateShareHandlerCertificateShareDataResponse;

export type GetCertificateShareDataError = CustomerrorErrResponse;

export interface GetCertificateShareDataParams {
    /** Share handle */
    handle: string;
}

/** @format binary */
export type GetCertificateShareImageData = File;

export type GetCertificateShareImageError = CustomerrorErrResponse;

export interface GetCertificateShareImageParams {
    /** Share handle */
    handle: string;
}

export type GetCertificatesListViewmodelData = EventGetCertificatesListViewModelResponse;

export type GetCertificatesListViewmodelError = CustomerrorErrResponse;

export interface GetCertificatesListViewmodelParams {
    /** Event ID */
    eventId: string;
}

export type GetClaimCertificateSignMessageData = CertificateGetClaimCertificateSignMessageResponse;

export type GetClaimCertificateSignMessageError = CustomerrorErrResponse;

export interface GetClaimCertificateSignMessageParams {
    /** Certificate ID */
    certificateId: string;
}

export type GetClosestIncomingScheduleData = SystemStatusGetClosestIncomingScheduleResponse;

export type GetClosestIncomingScheduleError = CustomerrorErrResponse;

export type GetEventByIdData = CoreApiInternalHandlerEventEventResponse;

export type GetEventByIdError = CustomerrorErrResponse;

export interface GetEventByIdParams {
    /** Event ID */
    eventId: string;
}

export type GetEventCertificateConfigData =
    CoreApiInternalHandlerEventconfigEventCertificateConfigResponse;

export type GetEventCertificateConfigError = CustomerrorErrResponse;

export interface GetEventCertificateConfigParams {
    /** Event ID */
    eventId: string;
}

export type GetEventCertificateFontFamiliesData =
    EventconfigGetEventCertificateFontFamiliesResponse;

export type GetEventCertificateFontFamiliesError = CustomerrorErrResponse;

export type GetEventCertificateTemplateData = string;

export type GetEventCertificateTemplateError = CustomerrorErrResponse;

export interface GetEventCertificateTemplateParams {
    /** Event ID */
    eventId: string;
}

export type GetEventCertificatesData = EventGetEventCertificatesResponse;

export type GetEventCertificatesError = CustomerrorErrResponse;

export interface GetEventCertificatesParams {
    /** Event ID */
    eventId: string;
}

export type GetEventContractByEventIdData = CoreApiInternalHandlerEventEventContractResponse;

export type GetEventContractByEventIdError = CustomerrorErrResponse;

export interface GetEventContractByEventIdParams {
    /** Event ID */
    eventId: string;
}

export type GetEventIssuerByIdData = EventEventIssuerResponse;

export type GetEventIssuerByIdError = CustomerrorErrResponse;

export interface GetEventIssuerByIdParams {
    /** Event ID */
    eventId: string;
    /** Issuer ID */
    issuerId: string;
}

export type GetEventIssuersByEventIdData = EventEventIssuerResponse[];

export type GetEventIssuersByEventIdError = CustomerrorErrResponse;

export interface GetEventIssuersByEventIdParams {
    /** Event ID */
    eventId: string;
}

export type GetEventParticipantsData = EventEventParticipantResponse[];

export type GetEventParticipantsError = CustomerrorErrResponse;

export interface GetEventParticipantsParams {
    /** Event ID */
    eventId: string;
}

export type GetEventRegistrationConfigData = EventconfigEventRegistrationConfigViewModel;

export type GetEventRegistrationConfigError = CustomerrorErrResponse;

export interface GetEventRegistrationConfigParams {
    /** Event ID */
    eventId: string;
}

export type GetEventRegistrationInvitationByUserAndEventData =
    EventRegistrationGetEventRegistrationInvitationByUserAndEventResponse;

export type GetEventRegistrationInvitationByUserAndEventError = CustomerrorErrResponse;

export interface GetEventRegistrationInvitationByUserAndEventParams {
    eventId: string;
}

export type GetEventRegistrationInvitationsByEventIdData =
    EventRegistrationEventRegistrationInvitationResponse[];

export type GetEventRegistrationInvitationsByEventIdError = CustomerrorErrResponse;

export interface GetEventRegistrationInvitationsByEventIdParams {
    eventId: string;
}

export type GetEventViewmodelByIdData = EventEventViewModel;

export type GetEventViewmodelByIdError = CustomerrorErrResponse;

export interface GetEventViewmodelByIdParams {
    /** Event ID */
    eventId: string;
}

export type GetEventsByOwnerCredentialsIdData = CoreApiInternalHandlerEventEventResponse[];

export type GetEventsByOwnerCredentialsIdError = CustomerrorErrResponse;

export interface GetEventsByOwnerCredentialsIdParams {
    /**
     * Limit
     * @default 10
     */
    limit?: number;
    /**
     * Offset
     * @default 0
     */
    offset?: number;
    /** Owner Credentials ID */
    ownerCredentialId: string;
}

export type GetEventsListData = EventGetEventListResponse;

export type GetEventsListError = CustomerrorErrResponse;

export interface GetEventsListParams {
    /**
     * Include active events
     * @default true
     */
    include_active_events?: boolean;
    /**
     * Include closed events
     * @default false
     */
    include_closed_events?: boolean;
    /**
     * Include inactive events
     * @default true
     */
    include_inactive_events?: boolean;
    /**
     * Only user joined events
     * @default false
     */
    only_user_joined_events?: boolean;
}

export type GetGasPriceData = BlockchainGasPriceResponse;

export type GetGasPriceError = CustomerrorErrResponse;

export type GetIssuerEventsData = IssuerIssuerEventResponse[];

export type GetIssuerEventsError = CustomerrorErrResponse;

export interface GetIssuerEventsParams {
    /** Issuer credential ID */
    issuer_credential_id?: string;
    /**
     * Limit
     * @default 10
     */
    limit?: number;
    /**
     * Offset
     * @default 0
     */
    offset?: number;
}

export type GetIssuerEventsViewmodelData = IssuerIssuerEventViewModel[];

export type GetIssuerEventsViewmodelError = CustomerrorErrResponse;

export interface GetIssuerEventsViewmodelParams {
    /** Issuer credential ID */
    issuer_credential_id?: string;
    /**
     * Limit
     * @default 10
     */
    limit?: number;
    /**
     * Offset
     * @default 0
     */
    offset?: number;
}

export type GetJoinEventSignMessageData = EventRegistrationGetJoinEventSignMessageResponse;

export type GetJoinEventSignMessageError = CustomerrorErrResponse;

export interface GetJoinEventSignMessageParams {
    /** Event ID */
    eventId: string;
}

export type GetLatestSchedulesData = SystemStatusGetLatestSchedulesResponse;

export type GetLatestSchedulesError = CustomerrorErrResponse;

export interface GetLatestSchedulesParams {
    /** Number of records to return */
    page_size: number;
}

export type GetMyCertificatesListViewmodelData = CertificateGetMyCertificatesListViewModelResponse;

export type GetMyCertificatesListViewmodelError = CustomerrorErrResponse;

export type GetMyProfileData = ProfileGetMyProfileViewModel;

export type GetMyProfileError = CustomerrorErr;

export type GetPlannedMaintenanceSchedulesData = SystemStatusGetPlannedMaintenanceSchedulesResponse;

export type GetPlannedMaintenanceSchedulesError = CustomerrorErrResponse;

export type GetSchedulesBetweenData = SystemStatusGetSchedulesBetweenResponse;

export type GetSchedulesBetweenError = CustomerrorErrResponse;

export interface GetSchedulesBetweenParams {
    /** End time unix timestamp */
    end_time: number;
    /** Start time unix timestamp */
    start_time: number;
}

export type GetSignMessageData = OnboardGetSignMessageResponse;

export type GetSystemStatusData = BlockchainGetSystemStatusResponse;

export type GetSystemStatusError = CustomerrorErrResponse;

export type GetSystemStatusViewmodelData = SystemStatusSystemStatusViewModel;

export type GetSystemStatusViewmodelError = CustomerrorErrResponse;

export type GetVerifiedIssuersData = IssuerProfileResponse[];

export type GetVerifiedIssuersError = CustomerrorErr;

export interface GetVerifiedIssuersParams {
    /** Limit */
    limit?: number;
    /** Offset */
    offset?: number;
    /** Search query (searches first name, last name, email, academic email, wallet address) */
    search?: string;
}

export type ImportCertificateReceiversData =
    CoreApiInternalHandlerEventImportCertificateReceiversResponse;

export type ImportCertificateReceiversError = CustomerrorErrResponse;

export interface ImportCertificateReceiversParams {
    eventId: string;
}

export type ImportEventParticipantsData = EventRegistrationEventRegistrationInvitationResponse[];

export type ImportEventParticipantsError = CustomerrorErrResponse;

export interface ImportEventParticipantsParams {
    eventId: string;
}

export interface InboxInboxMessageCertificateInvitationViewModel {
    certificate_id: string;
    certificate_title?: string;
    created_at: string;
    deleted_at?: string;
    event_id: string;
    event_name: string;
    has_participant_joined_event: boolean;
    hidden_at?: string;
    id: string;
    is_read: number;
    message_content: string;
    message_content_fallback?: string;
    message_type: EntityInboxMessageType;
    receiver_email?: string;
    receiver_wallet_address?: string;
    revoked_at?: string;
    sender_credential_email?: string;
    sender_credential_wallet_address?: string;
    token_id?: string;
    updated_at: string;
}

export interface InboxInboxMessagesEventRegistrationInvitationViewModel {
    academic_institution?: string;
    accepted_at?: string;
    broadcasted_at?: string;
    cancelled_at?: string;
    certificate_id?: string;
    certificate_title?: string;
    code?: string;
    created_at: string;
    deleted_at?: string;
    email?: string;
    event_id: string;
    event_name?: string;
    first_name?: string;
    has_participant_joined_event?: boolean;
    hidden_at?: string;
    id: string;
    is_read: number;
    last_name?: string;
    message_content: string;
    message_content_fallback?: string;
    message_type: EntityInboxMessageType;
    phone_number?: string;
    receiver_email?: string;
    receiver_wallet_address?: string;
    sender_credential_email?: string;
    sender_credential_wallet_address?: string;
    token_id?: string;
    updated_at: string;
    valid_until?: string;
}

export interface InboxInboxMessagesViewModel {
    certificate_id?: string;
    certificate_title?: string;
    created_at: string;
    deleted_at?: string;
    /** Certificate-specific fields (only populated for certificate invitation messages) */
    event_id?: string;
    event_name?: string;
    has_participant_joined_event?: boolean;
    hidden_at?: string;
    id: string;
    is_read: number;
    message_content: string;
    message_content_fallback?: string;
    message_type: EntityInboxMessageType;
    receiver_email?: string;
    receiver_wallet_address?: string;
    sender_credential_email?: string;
    sender_credential_wallet_address?: string;
    token_id?: string;
    updated_at: string;
}

export interface InboxmessagesGetInboxMessageResponse {
    event_certificate?: InboxInboxMessageCertificateInvitationViewModel;
    event_registration_invitation?: InboxInboxMessagesEventRegistrationInvitationViewModel;
    inbox_message: InboxInboxMessagesViewModel;
    inbox_message_type: EntityInboxMessageType;
}

export interface InboxmessagesGetInboxMessagesResponse {
    inbox_messages: InboxInboxMessagesViewModel[];
}

export interface InboxmessagesMarkAllMessagesAsReadResponse {
    inbox_messages: InboxInboxMessagesViewModel[];
}

export interface InboxmessagesMarkMessageAsReadRequest {
    message_id: string;
}

export interface InboxmessagesMarkMessageAsReadResponse {
    inbox_message: InboxInboxMessagesViewModel;
}

export interface IssuerIssuerEventResponse {
    created_at: string;
    event_end_date: string;
    event_id: string;
    event_location: string;
    event_owner_credential_id: string;
    event_short_description: string;
    event_start_date: string;
    event_title: string;
    id: string;
    is_signed: number;
    issuer_credential_id: string;
    updated_at: string;
}

export interface IssuerIssuerEventViewModel {
    certificate_count: number;
    created_at: string;
    event_end_date: string;
    event_id: string;
    event_location: string;
    event_owner_credential_id: string;
    event_short_description: string;
    event_start_date: string;
    event_title: string;
    id: string;
    is_signed: boolean;
    issuer_credential_id: string;
    owner_email: string;
    owner_google_email: string;
    owner_name: string;
    owner_wallet_address: string;
    updated_at: string;
}

export interface IssuerProfileResponse {
    academic_email?: string;
    academic_institution?: string;
    address?: string;
    authentication_credential_id: string;
    bio?: string;
    created_at: string;
    email?: string;
    first_name?: string;
    github_connector_ref?: string;
    /** Connector references (from authentication_credentials table) */
    google_connector_ref?: string;
    id: string;
    is_academic_email_public: boolean;
    is_academic_institution_public: boolean;
    is_address_public: boolean;
    is_bio_public: boolean;
    is_email_public: boolean;
    is_first_name_public: boolean;
    is_last_name_public: boolean;
    is_phone_number_public: boolean;
    is_profile_picture_public: boolean;
    last_name?: string;
    phone_number?: string;
    profile_picture_url?: string;
    updated_at: string;
    wallet_address?: string;
}

export type JoinEventData = EventRegistrationEventAttendeeResponse;

export type JoinEventError = CustomerrorErrResponse;

export interface JoinEventParams {
    /** Event ID */
    eventId: string;
}

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
    access_token: string;
    /** @minLength 6 */
    password: string;
}

export interface OnboardRegisterWithWalletRequest {
    signed_message: string;
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
    authentication_credential_id: string;
    bio?: string;
    created_at: string;
    email?: string;
    first_name?: string;
    github_connector_ref?: string;
    /** Connector references (from authentication_credentials table) */
    google_connector_ref?: string;
    id: string;
    is_academic_email_public: boolean;
    is_academic_institution_public: boolean;
    is_address_public: boolean;
    is_bio_public: boolean;
    is_email_public: boolean;
    is_first_name_public: boolean;
    is_last_name_public: boolean;
    is_phone_number_public: boolean;
    is_profile_picture_public: boolean;
    last_name?: string;
    phone_number?: string;
    profile_picture_url?: string;
    updated_at: string;
    wallet_address?: string;
}

export interface ProfileGetMyProfileViewModel {
    academic_email?: string;
    academic_institution?: string;
    address?: string;
    authentication_credential_created_at: string;
    authentication_credential_id: string;
    authentication_credential_updated_at: string;
    bio?: string;
    email?: string;
    first_name?: string;
    github_connector_ref?: string;
    google_connector_ref?: string;
    is_academic_email_public: boolean;
    is_academic_institution_public: boolean;
    is_address_public: boolean;
    is_bio_public: boolean;
    is_email_public: boolean;
    is_first_name_public: boolean;
    is_last_name_public: boolean;
    is_phone_number_public: boolean;
    is_profile_picture_public: boolean;
    last_name?: string;
    phone_number?: string;
    profile_created_at: string;
    profile_id: string;
    profile_picture_url?: string;
    profile_updated_at: string;
    solution_status: CommonSolutionStatus;
    unread_inbox_message_count: number;
    wallet_address: string;
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
    authentication_credential_id: string;
    bio?: string;
    created_at: string;
    email?: string;
    first_name?: string;
    github_connector_ref?: string;
    /** Connector references (from authentication_credentials table) */
    google_connector_ref?: string;
    id: string;
    is_academic_email_public: boolean;
    is_academic_institution_public: boolean;
    is_address_public: boolean;
    is_bio_public: boolean;
    is_email_public: boolean;
    is_first_name_public: boolean;
    is_last_name_public: boolean;
    is_phone_number_public: boolean;
    is_profile_picture_public: boolean;
    last_name?: string;
    phone_number?: string;
    profile_picture_url?: string;
    updated_at: string;
    wallet_address?: string;
}

export interface ProfileVerifyPasswordRequest {
    authentication_credential_id: string;
    password: string;
}

export interface ProfileVerifyPasswordResponse {
    is_success: boolean;
    message: string;
}

export type PublishEventCertificatesData =
    CoreApiInternalHandlerEventPublishEventCertificatesResponse;

export type PublishEventCertificatesError = CustomerrorErrResponse;

export interface PublishEventCertificatesParams {
    /** Event ID */
    eventId: string;
}

export type RegisterWithGoogleOauthData = OnboardRegisterResponse;

export type RegisterWithGoogleOauthError = CustomerrorErrResponse;

export type RegisterWithWalletData = OnboardRegisterResponse;

export type RegisterWithWalletError = CustomerrorErrResponse;

export type RequestGoogleOauthError = CustomerrorErrResponse;

export type RevokeAllEventCertificatesData =
    CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse;

export type RevokeAllEventCertificatesError = CustomerrorErrResponse;

export interface RevokeAllEventCertificatesParams {
    /** Event ID */
    eventId: string;
}

export type RevokeEventCertificatesData =
    CoreApiInternalHandlerEventRevokeEventCertificatesResponse;

export type RevokeEventCertificatesError = CustomerrorErrResponse;

export interface RevokeEventCertificatesParams {
    /** Event ID */
    eventId: string;
}

export type SignEventCertificatesData = CoreApiInternalHandlerEventSignEventCertificatesResponse;

export type SignEventCertificatesError = CustomerrorErrResponse;

export interface SignEventCertificatesParams {
    /** Event ID */
    eventId: string;
}

export interface SystemStatusGetClosestIncomingScheduleResponse {
    schedule?: EntitySystemStatusSchedule | null;
}

export interface SystemStatusGetLatestSchedulesResponse {
    schedules: EntitySystemStatusSchedule[];
}

export interface SystemStatusGetPlannedMaintenanceSchedulesResponse {
    schedules: EntitySystemStatusSchedule[];
}

export interface SystemStatusGetSchedulesBetweenResponse {
    schedules: EntitySystemStatusSchedule[];
}

export interface SystemStatusSystemStatusViewModel {
    gas_price: BlockchainGasPriceResponse;
    latest_schedule?: EntitySystemStatusSchedule | null;
}

export type ToggleCertificatePublishedData =
    CoreApiInternalHandlerEventconfigEventCertificateConfigResponse;

export type ToggleCertificatePublishedError = CustomerrorErrResponse;

export interface ToggleCertificatePublishedParams {
    /** Event ID */
    eventId: string;
}

export type UpdateCertificateShareData = CertificateShareHandlerUpdateCertificateShareResponse;

export type UpdateCertificateShareError = CustomerrorErrResponse;

export interface UpdateCertificateShareParams {
    /** Share ID */
    shareId: string;
}

export type UpdateEventCertificateConfigData =
    CoreApiInternalHandlerEventconfigEventCertificateConfigResponse;

export type UpdateEventCertificateConfigError = CustomerrorErrResponse;

export interface UpdateEventCertificateConfigParams {
    /** Event ID */
    eventId: string;
}

export interface UpdateEventCertificateConfigPayload {
    /**
     * Academic institution position x
     * @format float64
     */
    academic_institution_pos_x?: number;
    /**
     * Academic institution position y
     * @format float64
     */
    academic_institution_pos_y?: number;
    /** Base certificate image */
    base_certificate_image?: File;
    /**
     * Certificate subtitle position x
     * @format float64
     */
    certificate_subtitle_pos_x?: number;
    /**
     * Certificate subtitle position y
     * @format float64
     */
    certificate_subtitle_pos_y?: number;
    /**
     * Certificate title position x
     * @format float64
     */
    certificate_title_pos_x?: number;
    /**
     * Certificate title position y
     * @format float64
     */
    certificate_title_pos_y?: number;
    /**
     * Event name position x
     * @format float64
     */
    event_name_pos_x?: number;
    /**
     * Event name position y
     * @format float64
     */
    event_name_pos_y?: number;
    /**
     * Name position x
     * @format float64
     */
    name_pos_x: number;
    /**
     * Name position y
     * @format float64
     */
    name_pos_y: number;
}

export type UpdateEventCertificateTextConfigData = EventUpdateEventCertificateTextConfigResponse;

export type UpdateEventCertificateTextConfigError = CustomerrorErrResponse;

export interface UpdateEventCertificateTextConfigParams {
    /** Event ID */
    eventId: string;
}

export type UpdateEventContractData = CoreApiInternalHandlerEventEventContractResponse;

export type UpdateEventContractError = CustomerrorErrResponse;

export interface UpdateEventContractParams {
    /** Event ID */
    eventId: string;
}

export type UpdateEventData = CoreApiInternalHandlerEventEventResponse;

export type UpdateEventError = CustomerrorErrResponse;

export type UpdateEventIssuerData = EventEventIssuerResponse;

export type UpdateEventIssuerError = CustomerrorErrResponse;

export interface UpdateEventIssuerParams {
    /** Event ID */
    eventId: string;
}

/** Event issuer data */
export type UpdateEventIssuerPayload = EventUpdateEventIssuerRequest[];

export interface UpdateEventParams {
    eventId: string;
}

export interface UpdateEventPayload {
    /**
     * Event banner image (JPEG, PNG, WebP, max 10MB) - optional
     * @format binary
     */
    banner?: File;
    /** Contact address */
    contact_address?: string;
    /** Contact number */
    contact_number?: string;
    /** Event description */
    description?: string;
    /** End date (RFC3339 format) */
    end_date: string;
    /** Google map query */
    google_map_query?: string;
    /** Host password (required for contract update) */
    host_password: string;
    /**
     * Event icon image (JPEG, PNG, WebP, max 10MB) - optional
     * @format binary
     */
    icon?: File;
    /** Location */
    location?: string;
    /** Event name */
    name?: string;
    /** Number of seats */
    seats_count: number;
    /** Event short description */
    short_description?: string;
    /** Start date (RFC3339 format) */
    start_date: string;
}

export type UpdateEventRegistrationConfigData = EventconfigEventRegistrationConfigViewModel;

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

export type V1InboxMessagesDetailData = InboxmessagesGetInboxMessageResponse;

export interface V1InboxMessagesDetailParams {
    inboxMessageId: string;
}

export type V1InboxMessagesListData = InboxmessagesGetInboxMessagesResponse;

export type V1InboxMessagesReadAllUpdateData = InboxmessagesMarkAllMessagesAsReadResponse;

export type V1InboxMessagesReadUpdateData = InboxmessagesMarkMessageAsReadResponse;

export interface V1InboxMessagesReadUpdateParams {
    /** Message ID */
    messageId: string;
}

export type VerifyGoogleOauthError = CustomerrorErrResponse;

export interface VerifyGoogleOauthParams {
    code: string;
    state: string;
}

export type VerifyPasswordData = ProfileVerifyPasswordResponse;

export type VerifyPasswordError = CustomerrorErr;

import type { AxiosInstance, AxiosRequestConfig, HeadersDefaults, ResponseType } from "axios";
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

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

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
            const propertyContent: any[] = property instanceof Array ? property : [property];

            for (const formItem of propertyContent) {
                const isFileType = formItem instanceof Blob || formItem instanceof File;
                formData.append(key, isFileType ? formItem : this.stringifyFormItem(formItem));
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

        if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
            body = this.createFormData(body as Record<string, unknown>);
        }

        if (type === ContentType.Text && body && body !== null && typeof body !== "string") {
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
         * @description Check if the current user has specific roles (authenticated, host, issuer). Only returns requested fields.
         *
         * @tags Auth
         * @name CheckRole
         * @summary Check user roles
         * @request GET:/api/v1/auth/check-role
         */
        checkRole: (query: CheckRoleParams, params: RequestParams = {}) =>
            this.http.request<CheckRoleData, CheckRoleError>({
                path: `/api/v1/auth/check-role`,
                method: "GET",
                query: query,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

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
        verifyGoogleOauth: (query: VerifyGoogleOauthParams, params: RequestParams = {}) =>
            this.http.request<any, VerifyGoogleOauthError>({
                path: `/api/v1/auth/verify-google-oauth`,
                method: "GET",
                query: query,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Get current blockchain gas price and the configured maximum allowed fee
         *
         * @tags Blockchain
         * @name GetGasPrice
         * @summary Get current blockchain gas price
         * @request GET:/api/v1/blockchain/gas-price
         */
        getGasPrice: (params: RequestParams = {}) =>
            this.http.request<GetGasPriceData, GetGasPriceError>({
                path: `/api/v1/blockchain/gas-price`,
                method: "GET",
                format: "json",
                ...params,
            }),

        /**
         * @description Get current blockchain gas price and upcoming system status schedule in a single request
         *
         * @tags Blockchain
         * @name GetSystemStatus
         * @summary Get system status including gas price
         * @request GET:/api/v1/blockchain/status
         */
        getSystemStatus: (params: RequestParams = {}) =>
            this.http.request<GetSystemStatusData, GetSystemStatusError>({
                path: `/api/v1/blockchain/status`,
                method: "GET",
                format: "json",
                ...params,
            }),

        /**
         * @description Create a shareable link for a certificate. Only the certificate owner may call this endpoint. An optional password can be set to restrict access.
         *
         * @tags CertificateShares
         * @name CreateCertificateShare
         * @summary Create certificate share link
         * @request POST:/api/v1/certificate-shares/config/{certificate_id}
         * @secure
         */
        createCertificateShare: (
            { certificateId, ...query }: CreateCertificateShareParams,
            body: CertificateShareHandlerCreateCertificateShareBody,
            params: RequestParams = {},
        ) =>
            this.http.request<CreateCertificateShareData, CreateCertificateShareError>({
                path: `/api/v1/certificate-shares/config/${certificateId}`,
                method: "POST",
                body: body,
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Update the password of an existing share link. Only the certificate owner may call this endpoint. Omit the password field to leave existing protection unchanged; set to empty string or null to remove protection; set to a non-empty string to set a new password.
         *
         * @tags CertificateShares
         * @name UpdateCertificateShare
         * @summary Update certificate share link
         * @request PATCH:/api/v1/certificate-shares/config/{share_id}
         * @secure
         */
        updateCertificateShare: (
            { shareId, ...query }: UpdateCertificateShareParams,
            body: CertificateShareHandlerUpdateCertificateShareBody,
            params: RequestParams = {},
        ) =>
            this.http.request<UpdateCertificateShareData, UpdateCertificateShareError>({
                path: `/api/v1/certificate-shares/config/${shareId}`,
                method: "PATCH",
                body: body,
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Retrieve the full on-chain Verifiable Credential data for a share link. For password-protected shares, include the password in the request body.
         *
         * @tags CertificateShares
         * @name GetCertificateShareData
         * @summary Get on-chain certificate share data
         * @request POST:/api/v1/certificate-shares/{handle}
         */
        getCertificateShareData: (
            { handle, ...query }: GetCertificateShareDataParams,
            body: CertificateShareHandlerGetCertificateShareBody,
            params: RequestParams = {},
        ) =>
            this.http.request<GetCertificateShareDataData, GetCertificateShareDataError>({
                path: `/api/v1/certificate-shares/${handle}`,
                method: "POST",
                body: body,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Returns a PNG certificate image for a share link. For password-protected shares, include the password in the request body.
         *
         * @tags CertificateShares
         * @name GetCertificateShareImage
         * @summary Get certificate share image
         * @request POST:/api/v1/certificate-shares/{handle}/image
         */
        getCertificateShareImage: (
            { handle, ...query }: GetCertificateShareImageParams,
            body: CertificateShareHandlerGetCertificateShareImageBody,
            params: RequestParams = {},
        ) =>
            this.http.request<GetCertificateShareImageData, GetCertificateShareImageError>({
                path: `/api/v1/certificate-shares/${handle}/image`,
                method: "POST",
                body: body,
                type: ContentType.Json,
                format: "blob",
                ...params,
            }),

        /**
         * @description Queue a certificate claim using account password or wallet signature (requires authentication). The certificate will be minted asynchronously.
         *
         * @tags Certificates
         * @name ClaimCertificate
         * @summary Claim certificate
         * @request POST:/api/v1/certificates/claim/{certificate_id}
         * @secure
         */
        claimCertificate: (
            { certificateId, ...query }: ClaimCertificateParams,
            claimCertificateBody: CertificateClaimCertificateBody,
            params: RequestParams = {},
        ) =>
            this.http.request<ClaimCertificateData, ClaimCertificateError>({
                path: `/api/v1/certificates/claim/${certificateId}`,
                method: "POST",
                body: claimCertificateBody,
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get sign message for claiming a certificate (requires authentication)
         *
         * @tags Certificates
         * @name GetClaimCertificateSignMessage
         * @summary Get claim certificate sign message
         * @request GET:/api/v1/certificates/claim/{certificate_id}/sign-message
         * @secure
         */
        getClaimCertificateSignMessage: (
            { certificateId, ...query }: GetClaimCertificateSignMessageParams,
            params: RequestParams = {},
        ) =>
            this.http.request<
                GetClaimCertificateSignMessageData,
                GetClaimCertificateSignMessageError
            >({
                path: `/api/v1/certificates/claim/${certificateId}/sign-message`,
                method: "GET",
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get current user's certificates separated by claimed and unclaimed status. Claimed certificates have token_id populated, unclaimed certificates have token_id null and certificate config is published. Returns all certificates for the authenticated user across all events.
         *
         * @tags Certificates
         * @name GetMyCertificatesListViewmodel
         * @summary Get my certificates list viewmodel
         * @request GET:/api/v1/certificates/my-list-viewmodel
         */
        getMyCertificatesListViewmodel: (params: RequestParams = {}) =>
            this.http.request<
                GetMyCertificatesListViewmodelData,
                GetMyCertificatesListViewmodelError
            >({
                path: `/api/v1/certificates/my-list-viewmodel`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Generates a PNG certificate image for the authenticated user's certificate
         *
         * @tags certificates
         * @name GenerateCertificateImage
         * @summary Generate certificate image for participant
         * @request GET:/api/v1/certificates/{certificate_id}/image
         * @secure
         */
        generateCertificateImage: (
            { certificateId, ...query }: GenerateCertificateImageParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GenerateCertificateImageData, GenerateCertificateImageError>({
                path: `/api/v1/certificates/${certificateId}/image`,
                method: "GET",
                secure: true,
                format: "blob",
                ...params,
            }),

        /**
         * @description Cancel an event registration invitation by ID
         *
         * @tags Event Registration Invitation
         * @name CancelEventRegistrationInvitation
         * @summary Cancel event registration invitation
         * @request DELETE:/api/v1/event-registration/invitation
         */
        cancelEventRegistrationInvitation: (
            query: CancelEventRegistrationInvitationParams,
            params: RequestParams = {},
        ) =>
            this.http.request<
                CancelEventRegistrationInvitationData,
                CancelEventRegistrationInvitationError
            >({
                path: `/api/v1/event-registration/invitation`,
                method: "DELETE",
                query: query,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Import a list of participants for an event, creating inbox messages and event registration invitations
         *
         * @tags Event Registration Invitation
         * @name ImportEventParticipants
         * @summary Import event participants
         * @request POST:/api/v1/event-registration/invitation/:event_id/import
         */
        importEventParticipants: (
            { eventId, ...query }: ImportEventParticipantsParams,
            request: EventRegistrationImportEventParticipantsRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<ImportEventParticipantsData, ImportEventParticipantsError>({
                path: `/api/v1/event-registration/invitation/${eventId}/import`,
                method: "POST",
                body: request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get all event registration invitations for a specific event
         *
         * @tags Event Registration Invitation
         * @name GetEventRegistrationInvitationsByEventId
         * @summary Get event registration invitations by event ID
         * @request GET:/api/v1/event-registration/invitation/{eventId}
         */
        getEventRegistrationInvitationsByEventId: (
            { eventId, ...query }: GetEventRegistrationInvitationsByEventIdParams,
            params: RequestParams = {},
        ) =>
            this.http.request<
                GetEventRegistrationInvitationsByEventIdData,
                GetEventRegistrationInvitationsByEventIdError
            >({
                path: `/api/v1/event-registration/invitation/${eventId}`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Join an event
         *
         * @tags Event Registration
         * @name JoinEvent
         * @summary Join event
         * @request POST:/api/v1/event-registration/join/{event_id}
         */
        joinEvent: (
            { eventId, ...query }: JoinEventParams,
            joinEventBody: EventRegistrationJoinEventBody,
            params: RequestParams = {},
        ) =>
            this.http.request<JoinEventData, JoinEventError>({
                path: `/api/v1/event-registration/join/${eventId}`,
                method: "POST",
                body: joinEventBody,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get join event sign message
         *
         * @tags Event Registration
         * @name GetJoinEventSignMessage
         * @summary Get join event sign message
         * @request GET:/api/v1/event-registration/join/{event_id}/sign-message
         */
        getJoinEventSignMessage: (
            { eventId, ...query }: GetJoinEventSignMessageParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetJoinEventSignMessageData, GetJoinEventSignMessageError>({
                path: `/api/v1/event-registration/join/${eventId}/sign-message`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get event registration invitation of user and by event id
         *
         * @tags Event Registration Invitation
         * @name GetEventRegistrationInvitationByUserAndEvent
         * @summary Get event registration invitation of user and by event id
         * @request GET:/api/v1/event-registration/my/{event_id}
         */
        getEventRegistrationInvitationByUserAndEvent: (
            { eventId, ...query }: GetEventRegistrationInvitationByUserAndEventParams,
            params: RequestParams = {},
        ) =>
            this.http.request<
                GetEventRegistrationInvitationByUserAndEventData,
                GetEventRegistrationInvitationByUserAndEventError
            >({
                path: `/api/v1/event-registration/my/${eventId}`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Retrieves all font families that can be used in certificate templates, including their available weights and italic support
         *
         * @name GetEventCertificateFontFamilies
         * @summary Get all available font families for certificates
         * @request GET:/api/v1/eventconfig/certificate-font-families
         */
        getEventCertificateFontFamilies: (params: RequestParams = {}) =>
            this.http.request<
                GetEventCertificateFontFamiliesData,
                GetEventCertificateFontFamiliesError
            >({
                path: `/api/v1/eventconfig/certificate-font-families`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get events list
         *
         * @tags Event
         * @name GetEventsList
         * @summary Get events list
         * @request GET:/api/v1/events
         */
        getEventsList: (query: GetEventsListParams, params: RequestParams = {}) =>
            this.http.request<GetEventsListData, GetEventsListError>({
                path: `/api/v1/events`,
                method: "GET",
                query: query,
                format: "json",
                ...params,
            }),

        /**
         * @description Create a new event with banner image upload. Chain ID is automatically set from environment configuration.
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
         * @description Get events by owner credentials ID
         *
         * @name GetEventsByOwnerCredentialsId
         * @summary Get events by owner credentials ID
         * @request GET:/api/v1/events/owner-credentials/{owner_credential_id}
         */
        getEventsByOwnerCredentialsId: (
            { ownerCredentialId, ...query }: GetEventsByOwnerCredentialsIdParams,
            params: RequestParams = {},
        ) =>
            this.http.request<
                GetEventsByOwnerCredentialsIdData,
                GetEventsByOwnerCredentialsIdError
            >({
                path: `/api/v1/events/owner-credentials/${ownerCredentialId}`,
                method: "GET",
                query: query,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get event by ID
         *
         * @tags Events
         * @name GetEventById
         * @summary Get event by ID
         * @request GET:/api/v1/events/{event_id}
         */
        getEventById: ({ eventId, ...query }: GetEventByIdParams, params: RequestParams = {}) =>
            this.http.request<GetEventByIdData, GetEventByIdError>({
                path: `/api/v1/events/${eventId}`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Update an existing event with optional banner and icon image upload
         *
         * @tags Event
         * @name UpdateEvent
         * @summary Update an event
         * @request PUT:/api/v1/events/{event_id}
         */
        updateEvent: (
            { eventId, ...query }: UpdateEventParams,
            data: UpdateEventPayload,
            params: RequestParams = {},
        ) =>
            this.http.request<UpdateEventData, UpdateEventError>({
                path: `/api/v1/events/${eventId}`,
                method: "PUT",
                body: data,
                type: ContentType.FormData,
                format: "json",
                ...params,
            }),

        /**
         * @description Delete event by ID
         *
         * @tags Event
         * @name DeleteEventById
         * @summary Delete event by ID
         * @request DELETE:/api/v1/events/{event_id}
         */
        deleteEventById: (
            { eventId, ...query }: DeleteEventByIdParams,
            delete_event_request: EventDeleteEventRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<DeleteEventByIdData, DeleteEventByIdError>({
                path: `/api/v1/events/${eventId}`,
                method: "DELETE",
                body: delete_event_request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Delete event certificate configuration for an event
         *
         * @tags EventConfig
         * @name DeleteEventCertificateConfig
         * @summary Delete event certificate config
         * @request DELETE:/api/v1/events/{event_id}/certificate-config
         */
        deleteEventCertificateConfig: (
            { eventId, ...query }: DeleteEventCertificateConfigParams,
            params: RequestParams = {},
        ) =>
            this.http.request<DeleteEventCertificateConfigData, DeleteEventCertificateConfigError>({
                path: `/api/v1/events/${eventId}/certificate-config`,
                method: "DELETE",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get all certificates for an event
         *
         * @tags Event Certificates
         * @name GetEventCertificates
         * @summary Get event certificates
         * @request GET:/api/v1/events/{event_id}/certificates
         */
        getEventCertificates: (
            { eventId, ...query }: GetEventCertificatesParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetEventCertificatesData, GetEventCertificatesError>({
                path: `/api/v1/events/${eventId}/certificates`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Import certificate receivers for an event, deploys event certificate contract, and creates certificate records
         *
         * @tags Event Certificates
         * @name ImportCertificateReceivers
         * @summary Import certificate receivers for an event
         * @request POST:/api/v1/events/{event_id}/certificates/import
         */
        importCertificateReceivers: (
            { eventId, ...query }: ImportCertificateReceiversParams,
            request: CoreApiInternalHandlerEventImportCertificateReceiversRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<ImportCertificateReceiversData, ImportCertificateReceiversError>({
                path: `/api/v1/events/${eventId}/certificates/import`,
                method: "POST",
                body: request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get event certificates separated by claimed and unclaimed status. Claimed certificates have token_id populated, unclaimed certificates have token_id null and certificate config is published. Only event hosts and issuers can access this endpoint.
         *
         * @tags Event Certificates
         * @name GetCertificatesListViewmodel
         * @summary Get certificates list viewmodel
         * @request GET:/api/v1/events/{event_id}/certificates/list-viewmodel
         */
        getCertificatesListViewmodel: (
            { eventId, ...query }: GetCertificatesListViewmodelParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetCertificatesListViewmodelData, GetCertificatesListViewmodelError>({
                path: `/api/v1/events/${eventId}/certificates/list-viewmodel`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Publish certificates for an event. This will create inbox messages for all certificate receivers. All issuers must have signed before publishing.
         *
         * @name PublishEventCertificates
         * @summary Publish event certificates
         * @request POST:/api/v1/events/{event_id}/certificates/publish
         * @secure
         */
        publishEventCertificates: (
            { eventId, ...query }: PublishEventCertificatesParams,
            params: RequestParams = {},
        ) =>
            this.http.request<PublishEventCertificatesData, PublishEventCertificatesError>({
                path: `/api/v1/events/${eventId}/certificates/publish`,
                method: "POST",
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Revoke event certificates by certificate IDs
         *
         * @tags Event Certificates
         * @name RevokeEventCertificates
         * @summary Revoke event certificates
         * @request POST:/api/v1/events/{event_id}/certificates/revoke
         */
        revokeEventCertificates: (
            { eventId, ...query }: RevokeEventCertificatesParams,
            request: CoreApiInternalHandlerEventRevokeEventCertificatesRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<RevokeEventCertificatesData, RevokeEventCertificatesError>({
                path: `/api/v1/events/${eventId}/certificates/revoke`,
                method: "POST",
                body: request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Revoke all certificates for an event by event ID
         *
         * @tags Event Certificates
         * @name RevokeAllEventCertificates
         * @summary Revoke all event certificates
         * @request POST:/api/v1/events/{event_id}/certificates/revoke-all
         */
        revokeAllEventCertificates: (
            { eventId, ...query }: RevokeAllEventCertificatesParams,
            params: RequestParams = {},
        ) =>
            this.http.request<RevokeAllEventCertificatesData, RevokeAllEventCertificatesError>({
                path: `/api/v1/events/${eventId}/certificates/revoke-all`,
                method: "POST",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Sign all event certificates for an event by issuer
         *
         * @tags Event Certificates
         * @name SignEventCertificates
         * @summary Sign event certificates
         * @request POST:/api/v1/events/{event_id}/certificates/sign
         */
        signEventCertificates: (
            { eventId, ...query }: SignEventCertificatesParams,
            request: CoreApiInternalHandlerEventSignEventCertificatesRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<SignEventCertificatesData, SignEventCertificatesError>({
                path: `/api/v1/events/${eventId}/certificates/sign`,
                method: "POST",
                body: request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Update font family and font weight for all text templates in the certificate. This endpoint allows customization of fonts for event name, participant name, academic institution, certificate title, and certificate subtitle.
         *
         * @name UpdateEventCertificateTextConfig
         * @summary Update certificate text configuration
         * @request PUT:/api/v1/events/{event_id}/certificates/text-config
         */
        updateEventCertificateTextConfig: (
            { eventId, ...query }: UpdateEventCertificateTextConfigParams,
            body: EventUpdateEventCertificateTextConfigRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<
                UpdateEventCertificateTextConfigData,
                UpdateEventCertificateTextConfigError
            >({
                path: `/api/v1/events/${eventId}/certificates/text-config`,
                method: "PUT",
                body: body,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get the event certificate configuration for an event. Accessible by verified organizers or issuers assigned to the event.
         *
         * @name GetEventCertificateConfig
         * @summary Get event certificate config
         * @request GET:/api/v1/events/{event_id}/config/certificate
         */
        getEventCertificateConfig: (
            { eventId, ...query }: GetEventCertificateConfigParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetEventCertificateConfigData, GetEventCertificateConfigError>({
                path: `/api/v1/events/${eventId}/config/certificate`,
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
         * @request PUT:/api/v1/events/{event_id}/config/certificate
         */
        updateEventCertificateConfig: (
            { eventId, ...query }: UpdateEventCertificateConfigParams,
            data: UpdateEventCertificateConfigPayload,
            params: RequestParams = {},
        ) =>
            this.http.request<UpdateEventCertificateConfigData, UpdateEventCertificateConfigError>({
                path: `/api/v1/events/${eventId}/config/certificate`,
                method: "PUT",
                body: data,
                type: ContentType.FormData,
                format: "json",
                ...params,
            }),

        /**
         * @description Check if an event certificate is ready to be minted. Returns detailed status about configuration, signed issuers, and contract deployment. Accessible by verified organizers or issuers assigned to the event.
         *
         * @name CheckCertificateMintReadiness
         * @summary Check certificate mint readiness
         * @request GET:/api/v1/events/{event_id}/config/certificate/mint-readiness
         */
        checkCertificateMintReadiness: (
            { eventId, ...query }: CheckCertificateMintReadinessParams,
            params: RequestParams = {},
        ) =>
            this.http.request<
                CheckCertificateMintReadinessData,
                CheckCertificateMintReadinessError
            >({
                path: `/api/v1/events/${eventId}/config/certificate/mint-readiness`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Toggle the published status of event certificate configuration. When setting to published, inbox messages will be created for all certificate receivers. Only accessible by verified organizers.
         *
         * @name ToggleCertificatePublished
         * @summary Toggle certificate published status
         * @request PATCH:/api/v1/events/{event_id}/config/certificate/published
         */
        toggleCertificatePublished: (
            { eventId, ...query }: ToggleCertificatePublishedParams,
            request: EventconfigToggleCertificatePublishedRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<ToggleCertificatePublishedData, ToggleCertificatePublishedError>({
                path: `/api/v1/events/${eventId}/config/certificate/published`,
                method: "PATCH",
                body: request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Proxy the base certificate SVG template from storage so the frontend can access it without CORS restrictions from presigned URLs. Accessible by the event owner or issuers assigned to the event.
         *
         * @name GetEventCertificateTemplate
         * @summary Get event certificate template SVG
         * @request GET:/api/v1/events/{event_id}/config/certificate/template
         */
        getEventCertificateTemplate: (
            { eventId, ...query }: GetEventCertificateTemplateParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetEventCertificateTemplateData, GetEventCertificateTemplateError>({
                path: `/api/v1/events/${eventId}/config/certificate/template`,
                method: "GET",
                format: "blob",
                ...params,
            }),

        /**
         * @description Check if the password is correct for an event
         *
         * @tags Event Config
         * @name CheckEventPassword
         * @summary Check event password
         * @request POST:/api/v1/events/{event_id}/config/password-check
         */
        checkEventPassword: (
            { eventId, ...query }: CheckEventPasswordParams,
            request: EventconfigCheckEventPasswordBody,
            params: RequestParams = {},
        ) =>
            this.http.request<CheckEventPasswordData, CheckEventPasswordError>({
                path: `/api/v1/events/${eventId}/config/password-check`,
                method: "POST",
                body: request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get event registration configuration for an event
         *
         * @tags Events
         * @name GetEventRegistrationConfig
         * @summary Get event registration config view model
         * @request GET:/api/v1/events/{event_id}/config/registration
         */
        getEventRegistrationConfig: (
            { eventId, ...query }: GetEventRegistrationConfigParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetEventRegistrationConfigData, GetEventRegistrationConfigError>({
                path: `/api/v1/events/${eventId}/config/registration`,
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
         * @request PUT:/api/v1/events/{event_id}/config/registration
         */
        updateEventRegistrationConfig: (
            { eventId, ...query }: UpdateEventRegistrationConfigParams,
            request: EventconfigUpdateEventRegistrationConfigRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<
                UpdateEventRegistrationConfigData,
                UpdateEventRegistrationConfigError
            >({
                path: `/api/v1/events/${eventId}/config/registration`,
                method: "PUT",
                body: request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get the event contract for an event
         *
         * @tags Events
         * @name GetEventContractByEventId
         * @summary Get event contract by event ID
         * @request GET:/api/v1/events/{event_id}/contracts
         */
        getEventContractByEventId: (
            { eventId, ...query }: GetEventContractByEventIdParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetEventContractByEventIdData, GetEventContractByEventIdError>({
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
            request: EventUpdateEventContractRequest,
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
         * @tags Event
         * @name CreateEventContract
         * @summary Create event contract
         * @request POST:/api/v1/events/{event_id}/contracts
         */
        createEventContract: (
            { eventId, ...query }: CreateEventContractParams,
            request: EventCreateEventContractRequest,
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
         * @description Delete event contract for an event
         *
         * @tags Event Contracts
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
         * @tags Events
         * @name GetEventIssuersByEventId
         * @summary Get event issuers by event ID
         * @request GET:/api/v1/events/{event_id}/issuers
         */
        getEventIssuersByEventId: (
            { eventId, ...query }: GetEventIssuersByEventIdParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetEventIssuersByEventIdData, GetEventIssuersByEventIdError>({
                path: `/api/v1/events/${eventId}/issuers`,
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
         * @request PUT:/api/v1/events/{event_id}/issuers
         */
        updateEventIssuer: (
            { eventId, ...query }: UpdateEventIssuerParams,
            request: UpdateEventIssuerPayload,
            params: RequestParams = {},
        ) =>
            this.http.request<UpdateEventIssuerData, UpdateEventIssuerError>({
                path: `/api/v1/events/${eventId}/issuers`,
                method: "PUT",
                body: request,
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
            request: EventCreateEventIssuerRequest,
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
         * @description Delete an event issuer
         *
         * @tags Events
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
         * @description Get all participants (attendees) for an event
         *
         * @tags Events
         * @name GetEventParticipants
         * @summary Get event participants
         * @request GET:/api/v1/events/{event_id}/participants
         */
        getEventParticipants: (
            { eventId, ...query }: GetEventParticipantsParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetEventParticipantsData, GetEventParticipantsError>({
                path: `/api/v1/events/${eventId}/participants`,
                method: "GET",
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
            request: EventconfigCreateEventRegistrationConfigRequest,
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
         * @description Get event viewmodel by ID
         *
         * @tags Events
         * @name GetEventViewmodelById
         * @summary Get event viewmodel by ID
         * @request GET:/api/v1/events/{event_id}/viewmodel
         */
        getEventViewmodelById: (
            { eventId, ...query }: GetEventViewmodelByIdParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetEventViewmodelByIdData, GetEventViewmodelByIdError>({
                path: `/api/v1/events/${eventId}/viewmodel`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get my inbox messages
         *
         * @tags Inbox Messages
         * @name V1InboxMessagesList
         * @summary Get my inbox messages
         * @request GET:/api/v1/inbox-messages
         * @secure
         */
        v1InboxMessagesList: (params: RequestParams = {}) =>
            this.http.request<V1InboxMessagesListData, any>({
                path: `/api/v1/inbox-messages`,
                method: "GET",
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * No description
         *
         * @tags Inbox Messages
         * @name V1InboxMessagesReadUpdate
         * @request PUT:/api/v1/inbox-messages/read
         * @secure
         */
        v1InboxMessagesReadUpdate: (
            { messageId, ...query }: V1InboxMessagesReadUpdateParams,
            request: InboxmessagesMarkMessageAsReadRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<V1InboxMessagesReadUpdateData, any>({
                path: `/api/v1/inbox-messages/read`,
                method: "PUT",
                body: request,
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * No description
         *
         * @tags Inbox Messages
         * @name V1InboxMessagesReadAllUpdate
         * @request PUT:/api/v1/inbox-messages/read-all
         * @secure
         */
        v1InboxMessagesReadAllUpdate: (params: RequestParams = {}) =>
            this.http.request<V1InboxMessagesReadAllUpdateData, any>({
                path: `/api/v1/inbox-messages/read-all`,
                method: "PUT",
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get inbox message
         *
         * @tags Inbox Messages
         * @name V1InboxMessagesDetail
         * @summary Get inbox message
         * @request GET:/api/v1/inbox-messages/{inbox_message_id}
         * @secure
         */
        v1InboxMessagesDetail: (
            { inboxMessageId, ...query }: V1InboxMessagesDetailParams,
            params: RequestParams = {},
        ) =>
            this.http.request<V1InboxMessagesDetailData, any>({
                path: `/api/v1/inbox-messages/${inboxMessageId}`,
                method: "GET",
                secure: true,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get verified issuers with optional search query
         *
         * @tags Issuer
         * @name GetVerifiedIssuers
         * @summary Get verified issuers
         * @request GET:/api/v1/issuers
         */
        getVerifiedIssuers: (query: GetVerifiedIssuersParams, params: RequestParams = {}) =>
            this.http.request<GetVerifiedIssuersData, GetVerifiedIssuersError>({
                path: `/api/v1/issuers`,
                method: "GET",
                query: query,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get events assigned to the authenticated issuer for signing
         *
         * @tags Issuer
         * @name GetIssuerEvents
         * @summary Get events for issuer signing
         * @request GET:/api/v1/issuers/events
         */
        getIssuerEvents: (query: GetIssuerEventsParams, params: RequestParams = {}) =>
            this.http.request<GetIssuerEventsData, GetIssuerEventsError>({
                path: `/api/v1/issuers/events`,
                method: "GET",
                query: query,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get events assigned to the authenticated issuer for signing with view model data
         *
         * @tags Issuer
         * @name GetIssuerEventsViewmodel
         * @summary Get issuer events view model
         * @request GET:/api/v1/issuers/events/viewmodel
         */
        getIssuerEventsViewmodel: (
            query: GetIssuerEventsViewmodelParams,
            params: RequestParams = {},
        ) =>
            this.http.request<GetIssuerEventsViewmodelData, GetIssuerEventsViewmodelError>({
                path: `/api/v1/issuers/events/viewmodel`,
                method: "GET",
                query: query,
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
            message_signature: OnboardCheckOnboardStatusRequest,
            params: RequestParams = {},
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
            params: RequestParams = {},
        ) =>
            this.http.request<RegisterWithGoogleOauthData, RegisterWithGoogleOauthError>({
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
        createProfile: (profile: ProfileCreateProfileRequest, params: RequestParams = {}) =>
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
            request: ProfileUpdateProfileRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<UpdateProfileByCredentialIdData, UpdateProfileByCredentialIdError>({
                path: `/api/v1/profile/credential/${credentialId}`,
                method: "PATCH",
                body: request,
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Verify password
         *
         * @tags Profile
         * @name VerifyPassword
         * @summary Verify password
         * @request POST:/api/v1/profile/password/verify
         */
        verifyPassword: (
            verifyPasswordRequest: ProfileVerifyPasswordRequest,
            params: RequestParams = {},
        ) =>
            this.http.request<VerifyPasswordData, VerifyPasswordError>({
                path: `/api/v1/profile/password/verify`,
                method: "POST",
                body: verifyPasswordRequest,
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
         * @request GET:/api/v1/profile/viewmodel
         */
        getMyProfile: (params: RequestParams = {}) =>
            this.http.request<GetMyProfileData, GetMyProfileError>({
                path: `/api/v1/profile/viewmodel`,
                method: "GET",
                type: ContentType.Json,
                format: "json",
                ...params,
            }),

        /**
         * @description Get the closest incoming status update scheduled for the future
         *
         * @tags SystemStatus
         * @name GetClosestIncomingSchedule
         * @summary Get closest incoming status update
         * @request GET:/api/v1/system-status/closest-incoming
         */
        getClosestIncomingSchedule: (params: RequestParams = {}) =>
            this.http.request<GetClosestIncomingScheduleData, GetClosestIncomingScheduleError>({
                path: `/api/v1/system-status/closest-incoming`,
                method: "GET",
                format: "json",
                ...params,
            }),

        /**
         * @description Get latest system status schedules in the past
         *
         * @tags SystemStatus
         * @name GetLatestSchedules
         * @summary Get latest system status schedules
         * @request GET:/api/v1/system-status/latest
         */
        getLatestSchedules: (query: GetLatestSchedulesParams, params: RequestParams = {}) =>
            this.http.request<GetLatestSchedulesData, GetLatestSchedulesError>({
                path: `/api/v1/system-status/latest`,
                method: "GET",
                query: query,
                format: "json",
                ...params,
            }),

        /**
         * @description Get system status schedules updated within a specific time period using unix timestamps
         *
         * @tags SystemStatus
         * @name GetSchedulesBetween
         * @summary Get system status schedules between time period
         * @request GET:/api/v1/system-status/period
         */
        getSchedulesBetween: (query: GetSchedulesBetweenParams, params: RequestParams = {}) =>
            this.http.request<GetSchedulesBetweenData, GetSchedulesBetweenError>({
                path: `/api/v1/system-status/period`,
                method: "GET",
                query: query,
                format: "json",
                ...params,
            }),

        /**
         * @description Get all upcoming planned maintenance schedules
         *
         * @tags SystemStatus
         * @name GetPlannedMaintenanceSchedules
         * @summary Get planned maintenance schedules
         * @request GET:/api/v1/system-status/planned-maintenance
         */
        getPlannedMaintenanceSchedules: (params: RequestParams = {}) =>
            this.http.request<
                GetPlannedMaintenanceSchedulesData,
                GetPlannedMaintenanceSchedulesError
            >({
                path: `/api/v1/system-status/planned-maintenance`,
                method: "GET",
                format: "json",
                ...params,
            }),

        /**
         * @description Get the latest system status schedule combined with current gas price information
         *
         * @tags SystemStatus
         * @name GetSystemStatusViewmodel
         * @summary Get system status viewmodel
         * @request GET:/api/v1/system-status/viewmodel
         */
        getSystemStatusViewmodel: (params: RequestParams = {}) =>
            this.http.request<GetSystemStatusViewmodelData, GetSystemStatusViewmodelError>({
                path: `/api/v1/system-status/viewmodel`,
                method: "GET",
                format: "json",
                ...params,
            }),
    };
}
