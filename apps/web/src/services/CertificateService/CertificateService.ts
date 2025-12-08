import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import { mapBlobToCertificateImage, type CertificateImage } from "./mapper";

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
}

export const defaultCertificateService = new CertificateService(coreApiClient);
