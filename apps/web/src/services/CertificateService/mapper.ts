import type {
    CertificateGetMyCertificatesListViewModelResponse,
    CertificateClaimCertificateResponse,
    CoreApiInternalHandlerEventSignEventCertificatesResponse,
    CoreApiInternalHandlerEventImportCertificateReceiversResponse,
    CoreApiInternalHandlerEventRevokeEventCertificatesResponse,
    CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse,
    CoreApiInternalHandlerEventPublishEventCertificatesResponse,
    EventImportCertificateReceiverRequest,
    EntityEventCertificate,
    EntityCertificatePayload,
    EntityAttendeeProfileData,
    EntityCertificateRawData,
    CertificateShareHandlerCreateCertificateShareResponse,
    CertificateShareHandlerUpdateCertificateShareResponse,
    CertificateShareHandlerCertificateShareDataResponse,
    CertificateShareHandlerCertificateShareContractInfo,
    CertificateShareCertificateShareViewModel,
    EventClaimedCertificateViewModel,
    EventUnclaimedCertificateViewModel,
} from "@decm/api";

interface CertificateShareCertificateShareStatusResponse {
    status: string;
    certificate?: EntityEventCertificate;
}

export interface CertificateImage {
    url: string;
    blob: Blob;
    contentType: string;
}

export interface CertificateShare {
    id: string;
    eventCertificateId: string;
    active: boolean;
    handle: string;
    hasPassword: boolean;
    createdAt: string;
    updatedAt: string;
}

// Frontend camelCase certificate interface
export interface Certificate {
    id: string;
    eventId: string;
    eventName?: string;
    receiverCredentialId?: string;
    receiverEmail?: string;
    name?: string;
    academicInstitution?: string;
    certificateTitle?: string;
    certificateSubtitle?: string;
    eventContractAddress: string;
    eventCertificateAddress?: string;
    certificateTokenId?: string;
    createdAt: string;
    revokedAt?: string;
    issuerSignature?: string;
    certificateContractAddress?: string;
    inboxMessageId?: string;
    broadcastedAt?: string;
    estimatedDeadline?: string;
    abortedAt?: string;
    userClaimSignatureId?: string;
    certificateShare?: CertificateShare;
}

export interface ClaimedCertificate extends Certificate {
    status: string;
}

export type UnclaimedCertificate = Certificate;

export interface MyCertificatesViewModel {
    claimedCertificates: ClaimedCertificate[];
    unclaimedCertificates: UnclaimedCertificate[];
    totalClaimed: number;
    totalUnclaimed: number;
}

export interface SignedCertificate {
    certificateId: string;
    eventId: string;
    signature: string;
}

export interface SignedCertificatesResult {
    certificates: SignedCertificate[];
    totalSigned: number;
}

export type ImportCertificatesParams =
    | {
          eventId: string;
          hostPin: string;
          receivers: EventImportCertificateReceiverRequest[];
      }
    | {
          eventId: string;
          hostSignMessage: string;
          hostSignature: string;
          receivers: EventImportCertificateReceiverRequest[];
      };

export interface GetCertificateImportSignMessageParams {
    eventId: string;
    receivers: EventImportCertificateReceiverRequest[];
}

export interface GetCertificateImportSignMessageResult {
    signMessage: string;
    eventCertificateAddress: string;
}

export interface RevokeCertificatesParams {
    eventId: string;
    certificateIds: string[];
}

export interface ImportCertificatesResult {
    importedCount: number;
    failedCount: number;
    errors: string[];
}

export interface RevokeCertificatesResult {
    revokedCount: number;
    failedCount: number;
    errors: string[];
}

export interface RevokeAllCertificatesResult {
    revokedCount: number;
    message: string;
}

export interface PublishCertificatesResult {
    publishedCount: number;
    failedCount: number;
    errors: string[];
}

export interface ClaimCertificateSignMessage {
    signMessage: string;
    certificateId: string;
}

export interface ClaimCertificateResult {
    certificateId: string;
    transactionHash: string;
    claimedAt: string;
    message: string;
    status: string;
    estimatedDeadline?: string;
}

export interface ClaimCertificateWithPinParams {
    certificateId: string;
    accountPassword: string;
}

export interface ClaimCertificateWithSignatureParams {
    certificateId: string;
    signature: string;
    signMessage: string;
}

export interface CreateCertificateShareableLinkResult {
    certificateId: string;
    isPublished: boolean;
    handle: string;
}

export interface UpdateCertificateShareResult {
    certificateId: string;
    handle: string;
}

export type CertificateShareViewStatus = "READY" | "VALID_BUT_PENDING" | "PASSWORD_LOCKED";

export interface GetCertificateShareDataResult {
    payload: EntityCertificatePayload;
    contract: CertificateShareHandlerCertificateShareContractInfo;
    decryptedUserData?: EntityAttendeeProfileData;
    decryptedCertificateData?: EntityCertificateRawData;
}

export interface GetCertificateShareStatusResult {
    status: CertificateShareViewStatus;
    certificate?: Certificate;
}

/**
 * Converts a Blob to a CertificateImage object with object URL
 * @param blob - The image blob from the API
 * @returns CertificateImage with object URL
 */
export const mapBlobToCertificateImage = (blob: Blob): CertificateImage => {
    const url = URL.createObjectURL(blob);
    return {
        url,
        blob,
        contentType: blob.type || "image/png",
    };
};

/**
 * Maps API certificate share viewmodel to frontend camelCase
 */
export const mapCertificateShare = (
    share: CertificateShareCertificateShareViewModel,
): CertificateShare => ({
    id: share.id,
    eventCertificateId: share.event_certificate_id,
    active: share.active,
    handle: share.handle,
    hasPassword: share.has_password,
    createdAt: share.created_at,
    updatedAt: share.updated_at,
});

/**
 * Maps API snake_case certificate to frontend camelCase
 */
export const mapCertificate = (cert: EntityEventCertificate): Certificate => {
    return {
        id: cert.id,
        eventId: cert.event_id,
        eventName: cert.event_name,
        receiverCredentialId: cert.receiver_credential_id,
        receiverEmail: cert.receiver_email,
        name: cert.name,
        academicInstitution: cert.academic_institution,
        certificateTitle: cert.certificate_title,
        certificateSubtitle: cert.certificate_subtitle,
        eventContractAddress: cert.event_contract_address,
        eventCertificateAddress: cert.event_certificate_address,
        certificateTokenId: cert.certificate_token_id,
        createdAt: cert.created_at,
        revokedAt: cert.revoked_at,
        inboxMessageId: cert.inbox_message_id,
        broadcastedAt: cert.broadcasted_at,
        estimatedDeadline: cert.estimated_deadline,
        abortedAt: cert.aborted_at,
        userClaimSignatureId: cert.user_claim_signature_id,
    };
};

const mapClaimedCertificate = (cert: EventClaimedCertificateViewModel): ClaimedCertificate => ({
    id: cert.id,
    eventId: cert.event_id,
    eventName: cert.event_name,
    receiverCredentialId: cert.receiver_credential_id,
    receiverEmail: cert.receiver_email,
    name: cert.name,
    academicInstitution: cert.academic_institution,
    certificateTitle: cert.certificate_title,
    certificateSubtitle: cert.certificate_subtitle,
    eventContractAddress: cert.event_contract_address,
    eventCertificateAddress: cert.event_certificate_address,
    certificateTokenId: cert.certificate_token_id,
    createdAt: cert.created_at,
    revokedAt: cert.revoked_at,
    inboxMessageId: cert.inbox_message_id,
    broadcastedAt: cert.broadcasted_at,
    estimatedDeadline: cert.estimated_deadline,
    abortedAt: cert.aborted_at,
    userClaimSignatureId: cert.user_claim_signature_id,
    status: cert.status,
    certificateShare: cert.certificate_share
        ? mapCertificateShare(cert.certificate_share)
        : undefined,
});

const mapUnclaimedCertificate = (
    cert: EventUnclaimedCertificateViewModel,
): UnclaimedCertificate => ({
    id: cert.id,
    eventId: cert.event_id,
    eventName: cert.event_name,
    receiverCredentialId: cert.receiver_credential_id,
    receiverEmail: cert.receiver_email,
    name: cert.name,
    academicInstitution: cert.academic_institution,
    certificateTitle: cert.certificate_title,
    certificateSubtitle: cert.certificate_subtitle,
    eventContractAddress: cert.event_contract_address,
    eventCertificateAddress: cert.event_certificate_address,
    certificateTokenId: cert.certificate_token_id,
    createdAt: cert.created_at,
    revokedAt: cert.revoked_at,
    inboxMessageId: cert.inbox_message_id,
    broadcastedAt: cert.broadcasted_at,
    estimatedDeadline: cert.estimated_deadline,
    abortedAt: cert.aborted_at,
    userClaimSignatureId: cert.user_claim_signature_id,
    certificateShare: cert.certificate_share
        ? mapCertificateShare(cert.certificate_share)
        : undefined,
});

/**
 * Maps API response to MyCertificatesViewModel
 */
export const mapToMyCertificatesViewModel = (
    response: CertificateGetMyCertificatesListViewModelResponse,
): MyCertificatesViewModel => {
    return {
        claimedCertificates: (response.claimed_certificates || []).map(mapClaimedCertificate),
        unclaimedCertificates: (response.unclaimed_certificates || []).map(mapUnclaimedCertificate),
        totalClaimed: response.total_claimed || 0,
        totalUnclaimed: response.total_unclaimed || 0,
    };
};

/**
 * Maps sign certificates response to SignedCertificatesResult
 */
export const mapToSignedCertificatesResult = (
    response: CoreApiInternalHandlerEventSignEventCertificatesResponse,
): SignedCertificatesResult => {
    const certificates = (response.certificates || []).map((cert) => ({
        certificateId: cert.certificate?.id || "",
        eventId: cert.certificate?.event_id || "",
        signature: cert.signature || "",
    }));

    return {
        certificates,
        totalSigned: certificates.length,
    };
};

/**
 * Maps import certificates response
 */
export const mapImportCertificatesResponse = (
    response: CoreApiInternalHandlerEventImportCertificateReceiversResponse,
): ImportCertificatesResult => {
    return {
        importedCount: response.certificates?.length || 0,
        failedCount: 0, // API doesn't provide failed count
        errors: [], // API doesn't provide errors
    };
};

/**
 * Maps revoke certificates response
 */
export const mapRevokeCertificatesResponse = (
    response: CoreApiInternalHandlerEventRevokeEventCertificatesResponse,
): RevokeCertificatesResult => {
    return {
        revokedCount: response.revoked_certificates?.length || 0,
        failedCount: 0, // API doesn't provide failed count
        errors: [], // API doesn't provide errors
    };
};

/**
 * Maps revoke all certificates response
 */
export const mapRevokeAllCertificatesResponse = (
    response: CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse,
): RevokeAllCertificatesResult => {
    return {
        revokedCount: response.revoked_certificates?.length || 0,
        message: `Successfully revoked ${response.revoked_certificates?.length || 0} certificates`,
    };
};

/**
 * Maps publish certificates response
 */
export const mapPublishCertificatesResponse = (
    response: CoreApiInternalHandlerEventPublishEventCertificatesResponse,
): PublishCertificatesResult => {
    return {
        publishedCount: response.published_count || 0,
        failedCount: 0, // API doesn't provide failed count
        errors: [], // API doesn't provide errors
    };
};

/**
 * Maps claim certificate sign message response
 */
export const mapClaimCertificateSignMessage = (response: {
    sign_message: string;
}): ClaimCertificateSignMessage => {
    return {
        signMessage: response.sign_message || "",
        certificateId: "", // Certificate ID is not in the response, caller should provide it
    };
};

/**
 * Maps claim certificate response
 */
export const mapClaimCertificateResponse = (
    response: CertificateClaimCertificateResponse,
): ClaimCertificateResult => {
    const cert = response.certificate as EntityEventCertificate;
    const sig = response.user_signature as { estimated_deadline?: string } | null;
    return {
        certificateId: cert?.id || "",
        transactionHash: cert?.certificate_token_id || "", // Use token ID as transaction hash indicator
        claimedAt: cert?.created_at || "",
        message: response.message || "Certificate claimed successfully",
        status: response.status || "success",
        estimatedDeadline: sig?.estimated_deadline,
    };
};

/**
 * Map create certificate shareable link response
 */
export const mapCreateCertificateShareableLinkResponse = (
    response: CertificateShareHandlerCreateCertificateShareResponse,
): CreateCertificateShareableLinkResult => {
    return {
        certificateId: response.share.event_certificate_id,
        isPublished: response.share.active,
        handle: response.share.handle,
    };
};

/**
 * Maps update certificate share response
 */
export const mapUpdateCertificateShareResponse = (
    response: CertificateShareHandlerUpdateCertificateShareResponse,
): UpdateCertificateShareResult => ({
    certificateId: response.share.event_certificate_id,
    handle: response.share.handle,
});

/**
 * Maps get certificate share status response
 */
export const mapGetCertificateShareStatusResponse = (
    response: CertificateShareCertificateShareStatusResponse,
): GetCertificateShareStatusResult => ({
    status: response.status as CertificateShareViewStatus,
    certificate: response.certificate
        ? mapCertificate(response.certificate as EntityEventCertificate)
        : undefined,
});

/**
 * Maps get certificate share data response
 */
export const mapGetCertificateShareDataResponse = (
    response: CertificateShareHandlerCertificateShareDataResponse,
): GetCertificateShareDataResult => ({
    payload: response.data,
    contract: response.contract,
    decryptedUserData: response.decryptedUserData,
    decryptedCertificateData: response.decryptedCertificateData,
});
