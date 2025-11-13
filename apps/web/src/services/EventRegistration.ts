import type { CoreApiType } from "@/lib/api/api";

export class EventRegistrationService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    public async getConfiguration(eventId: string) {
        const response = await this._coreApi.v1.getEventRegistrationConfig({ eventId });
        return response;
    }

    public async checkPassword(eventId: string, password: string) {
        const response = await this._coreApi.v1.checkEventPassword({ eventId }, { password });
        return response;
    }
}
