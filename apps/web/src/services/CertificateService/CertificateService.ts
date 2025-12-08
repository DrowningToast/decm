import type { CoreApiType } from "@/lib/api/api";
import { coreApiClient } from "@/lib/api/api";
import { env } from "@/config/env";
import { mapBlobToCertificateImage } from "./mapper";
import type { CertificateImage } from "./mapper";

export class CertificateService {
    private _coreApi: CoreApiType;
    private _baseUrl: string;

    constructor(coreApi: CoreApiType, baseUrl: string = env.VITE_CORE_BACKEND_API) {
        this._coreApi = coreApi;
        this._baseUrl = baseUrl;
    }

    public async getMyCertificates(eventId: string) {
        const response = await this._coreApi.v1.getEventCertificates({ eventId });
        return response;
    }

    /**
     * Fetch certificate image with authentication
     * Returns a blob URL that can be used as an image src
     */
    public async getCertificateImage(certificateId: string): Promise<CertificateImage> {
        const response = await fetch(
            `${this._baseUrl}/api/v1/certificates/${certificateId}/image`,
            {
                method: "GET",
                credentials: "include", // Include cookies for authentication
            },
        );

        if (!response.ok) {
            throw new Error(`Failed to fetch certificate image: ${response.statusText}`);
        }

        const blob = await response.blob();
        return mapBlobToCertificateImage(blob, certificateId);
    }
}

export const defaultCertificateService = new CertificateService(coreApiClient);
