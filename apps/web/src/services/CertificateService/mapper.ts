import type {
    EventGetMyCertificatesListViewModelResponse,
    CoreApiInternalHandlerEventSignEventCertificatesResponse,
    CoreApiInternalHandlerEventImportCertificateReceiversResponse,
    CoreApiInternalHandlerEventRevokeEventCertificatesResponse,
    CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse,
    CoreApiInternalHandlerEventPublishEventCertificatesResponse,
    EventImportCertificateReceiverRequest,
    EntityEventCertificate,
} from "@decm/api";

export interface CertificateImage {
    url: string;
    blob: Blob;
    contentType: string;
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
}

export interface MyCertificatesViewModel {
    claimedCertificates: Certificate[];
    unclaimedCertificates: Certificate[];
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

export interface ImportCertificatesParams {
    eventId: string;
    hostPin: string;
    receivers: EventImportCertificateReceiverRequest[];
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
    };
};

/**
 * Maps API response to MyCertificatesViewModel
 */
export const mapToMyCertificatesViewModel = (
    response: EventGetMyCertificatesListViewModelResponse,
): MyCertificatesViewModel => {
    return {
        claimedCertificates: (response.claimed_certificates || []).map(mapCertificate),
        unclaimedCertificates: (response.unclaimed_certificates || []).map(mapCertificate),
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
        certificateId: cert.certificate_id || "",
        eventId: cert.event_id || "",
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
        importedCount: response.imported_count || 0,
        failedCount: response.failed_count || 0,
        errors: response.errors || [],
    };
};

/**
 * Maps revoke certificates response
 */
export const mapRevokeCertificatesResponse = (
    response: CoreApiInternalHandlerEventRevokeEventCertificatesResponse,
): RevokeCertificatesResult => {
    return {
        revokedCount: response.revoked_count || 0,
        failedCount: response.failed_count || 0,
        errors: response.errors || [],
    };
};

/**
 * Maps revoke all certificates response
 */
export const mapRevokeAllCertificatesResponse = (
    response: CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse,
): RevokeAllCertificatesResult => {
    return {
        revokedCount: response.revoked_count || 0,
        message: response.message || "",
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
        failedCount: response.failed_count || 0,
        errors: response.errors || [],
    };
};

/**
 * Maps claim certificate sign message response
 */
export const mapClaimCertificateSignMessage = (response: {
    data?: { sign_message?: string; certificate_id?: string };
}): ClaimCertificateSignMessage => {
    return {
        signMessage: response.data?.sign_message || "",
        certificateId: response.data?.certificate_id || "",
    };
};

/**
 * Maps claim certificate response
 */
export const mapClaimCertificateResponse = (response: {
    data?: {
        certificate_id?: string;
        transaction_hash?: string;
        claimed_at?: string;
        message?: string;
    };
}): ClaimCertificateResult => {
    return {
        certificateId: response.data?.certificate_id || "",
        transactionHash: response.data?.transaction_hash || "",
        claimedAt: response.data?.claimed_at || "",
        message: response.data?.message || "",
    };
};
