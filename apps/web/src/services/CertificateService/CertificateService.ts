import type { CoreApiType } from "@/lib/api/api";
import { coreApiClient } from "@/lib/api/api";

export class CertificateService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    public async getMyCertificates(eventId: string) {
        const response = await this._coreApi.v1.getEventCertificates({ eventId });
        return response;
    }
}

export const defaultCertificateService = new CertificateService(coreApiClient);
