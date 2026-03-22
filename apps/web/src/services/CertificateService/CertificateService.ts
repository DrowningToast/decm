import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import {
    mapBlobToCertificateImage,
    mapToMyCertificatesViewModel,
    mapToSignedCertificatesResult,
    mapImportCertificatesResponse,
    mapRevokeCertificatesResponse,
    mapRevokeAllCertificatesResponse,
    mapPublishCertificatesResponse,
    mapClaimCertificateSignMessage,
    mapClaimCertificateResponse,
    mapCertificate,
    type CertificateImage,
    type MyCertificatesViewModel,
    type SignedCertificatesResult,
    type ImportCertificatesParams,
    type RevokeCertificatesParams,
    type Certificate,
    type ImportCertificatesResult,
    type RevokeCertificatesResult,
    type RevokeAllCertificatesResult,
    type PublishCertificatesResult,
    type ClaimCertificateSignMessage,
    type ClaimCertificateResult,
    type ClaimCertificateWithPinParams,
    type ClaimCertificateWithSignatureParams,
    type CreateCertificateShareableLinkResult,
    mapCreateCertificateShareableLinkResponse,
    type UpdateCertificateShareResult,
    mapUpdateCertificateShareResponse,
    type GetCertificateShareDataResult,
    mapGetCertificateShareDataResponse,
} from "./mapper";

export interface FontFamily {
    name: string;
    displayName: string;
    weights: number[];
}

export interface CertificateMintReadiness {
    isReady: boolean;
    missingRequirements: string[];
    readinessInfo: {
        hasBaseCertificateImage: boolean;
        hasEventContract: boolean;
        hasCertificateContract: boolean;
        hasAtLeastOneIssuer: boolean;
        hasCertificateReceivers: boolean;
    };
}

export interface CertificateFontFamily {
    id: number;
    font_family_name: string;
    css_font_name: string;
    is_default: boolean;
    available_font_weights: number[];
    is_support_italic: boolean;
}

export class CertificateService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    /**
     * Fetches certificate image for the given certificate ID
     * @param certificateId - The ID of the certificate
     * @returns CertificateImage with object URL
     * @throws Error if the image cannot be fetched
     */
    public async getCertificateImage(certificateId: string): Promise<CertificateImage> {
        const response = await this._coreApi.v1.generateCertificateImage({
            certificateId,
        });

        // The response is a File/Blob
        if (!(response instanceof Blob)) {
            throw new Error("Invalid response: expected Blob");
        }

        return mapBlobToCertificateImage(response);
    }

    /**
     * Fetches the list of certificates for the authenticated user
     * @returns MyCertificatesViewModel with claimed and unclaimed certificates
     */
    public async getMyCertificatesList(): Promise<MyCertificatesViewModel> {
        const response = await this._coreApi.v1.getMyCertificatesListViewmodel();
        return mapToMyCertificatesViewModel(response);
    }

    /**
     * Fetches certificates for a specific event
     * @param eventId - The event ID
     * @returns Array of certificates in camelCase
     */
    public async getEventCertificates(eventId: string): Promise<Certificate[]> {
        const response = await this._coreApi.v1.getEventCertificates({ eventId });
        return (response.certificates || []).map(mapCertificate);
    }

    /**
     * Fetches available font families for certificates
     * @returns Array of font families in the format expected by the component
     */
    public async getFontFamilies(): Promise<CertificateFontFamily[]> {
        const response = await this._coreApi.v1.getEventCertificateFontFamilies();
        return response.font_families || [];
    }

    /**
     * Signs event certificates with issuer PIN
     * @param eventId - The event ID
     * @param issuerPin - The issuer's PIN for signing
     * @returns SignedCertificatesResult with signed certificates
     */
    public async signEventCertificates(
        eventId: string,
        issuerPin: string,
    ): Promise<SignedCertificatesResult> {
        const response = await this._coreApi.v1.signEventCertificates(
            { eventId },
            { issuer_pin: issuerPin },
        );
        return mapToSignedCertificatesResult(response);
    }

    /**
     * Imports certificate receivers for an event
     * @param params - Import parameters including event ID, host PIN, and receivers
     */
    public async importCertificates(
        params: ImportCertificatesParams,
    ): Promise<ImportCertificatesResult> {
        const response = await this._coreApi.v1.importCertificateReceivers(
            { eventId: params.eventId },
            {
                event_id: params.eventId,
                host_pin: params.hostPin,
                receivers: params.receivers,
            },
        );
        return mapImportCertificatesResponse(response);
    }

    /**
     * Revokes specific certificates for an event
     * @param params - Revoke parameters including event ID and certificate IDs
     */
    public async revokeCertificates(
        params: RevokeCertificatesParams,
    ): Promise<RevokeCertificatesResult> {
        const response = await this._coreApi.v1.revokeEventCertificates(
            { eventId: params.eventId },
            { certificate_ids: params.certificateIds },
        );
        return mapRevokeCertificatesResponse(response);
    }

    /**
     * Revokes all certificates for an event
     * @param eventId - The event ID
     */
    public async revokeAllCertificates(eventId: string): Promise<RevokeAllCertificatesResult> {
        const response = await this._coreApi.v1.revokeAllEventCertificates({ eventId });
        return mapRevokeAllCertificatesResponse(response);
    }

    /**
     * Publishes event certificates
     * @param eventId - The event ID
     */
    public async publishCertificates(eventId: string): Promise<PublishCertificatesResult> {
        const response = await this._coreApi.v1.publishEventCertificates({ eventId });
        return mapPublishCertificatesResponse(response);
    }

    /**
     * Toggles certificate published status
     * @param eventId - The event ID
     * @param isPublished - Whether to publish or unpublish
     */
    public async toggleCertificatePublished(eventId: string, isPublished: boolean) {
        const response = await this._coreApi.v1.toggleCertificatePublished(
            { eventId },
            { is_published: isPublished },
        );
        return response;
    }

    /**
     * Checks certificate mint readiness for an event
     * @param eventId - The event ID
     * @returns Certificate mint readiness information in camelCase
     */
    public async checkCertificateMintReadiness(eventId: string): Promise<CertificateMintReadiness> {
        const response = await this._coreApi.v1.checkCertificateMintReadiness({ eventId });
        return {
            isReady: response.is_ready || false,
            missingRequirements: response.missing_requirements || [],
            readinessInfo: {
                hasBaseCertificateImage: response.has_certificate_config || false,
                hasEventContract: true, // Event contract should always exist if we can check readiness
                hasCertificateContract: response.has_certificate_contract || false,
                hasAtLeastOneIssuer: response.total_issuers_count > 0,
                hasCertificateReceivers: true, // This should be checked elsewhere if needed
            },
        };
    }

    /**
     * Gets the sign message for claiming a certificate via wallet signature
     * @param certificateId - The certificate ID
     * @returns Sign message that needs to be signed by user's wallet
     */
    public async getClaimCertificateSignMessage(
        certificateId: string,
    ): Promise<ClaimCertificateSignMessage> {
        const response = await this._coreApi.v1.getClaimCertificateSignMessage({ certificateId });
        return mapClaimCertificateSignMessage(response);
    }

    /**
     * Claims a certificate with PIN/password
     * @param params - Certificate ID and account password
     * @returns Claim result with transaction hash
     */
    private async claimCertificateWithPin(
        params: ClaimCertificateWithPinParams,
    ): Promise<ClaimCertificateResult> {
        const response = await this._coreApi.v1.claimCertificate(
            { certificateId: params.certificateId },
            { account_password: params.accountPassword },
        );
        return mapClaimCertificateResponse(response);
    }

    /**
     * Claims a certificate with wallet signature
     * @param params - Certificate ID, signature, and sign message
     * @returns Claim result with transaction hash
     */
    private async claimCertificateWithSignature(
        params: ClaimCertificateWithSignatureParams,
    ): Promise<ClaimCertificateResult> {
        const response = await this._coreApi.v1.claimCertificate(
            { certificateId: params.certificateId },
            {
                signature: params.signature,
                sign_message: params.signMessage,
            },
        );
        return mapClaimCertificateResponse(response);
    }

    /**
     * Claims a certificate (unified method supporting both PIN and signature)
     * @param params - Either PIN params or signature params
     * @returns Claim result with transaction hash
     */
    public async claimCertificate(
        params: ClaimCertificateWithPinParams | ClaimCertificateWithSignatureParams,
    ): Promise<ClaimCertificateResult> {
        if ("accountPassword" in params) {
            return this.claimCertificateWithPin(params);
        } else {
            return this.claimCertificateWithSignature(params);
        }
    }

    public async createCertificateShareableLink(
        certificateId: string,
        password?: string,
    ): Promise<CreateCertificateShareableLinkResult> {
        const response = await this._coreApi.v1.createCertificateShare(
            { certificateId },
            { password },
        );
        return mapCreateCertificateShareableLinkResponse(response);
    }

    public async updateCertificateShare(
        shareId: string,
        params: { password?: string; active?: boolean },
    ): Promise<UpdateCertificateShareResult> {
        const response = await this._coreApi.v1.updateCertificateShare(
            { shareId },
            params as { active: boolean; password: string },
        );
        return mapUpdateCertificateShareResponse(response);
    }

    public async getCertificateShareData(
        handle: string,
        password?: string,
    ): Promise<GetCertificateShareDataResult> {
        const response = await this._coreApi.v1.getCertificateShareData({ handle }, { password });
        return mapGetCertificateShareDataResponse(response);
    }

    public async getCertificateShareImage(
        handle: string,
        password?: string,
    ): Promise<CertificateImage> {
        const response = await this._coreApi.v1.getCertificateShareImage({ handle }, { password });

        if (!(response instanceof Blob)) {
            throw new Error("Invalid response: expected Blob");
        }

        return mapBlobToCertificateImage(response);
    }
}

export const defaultCertificateService = new CertificateService(coreApiClient);
