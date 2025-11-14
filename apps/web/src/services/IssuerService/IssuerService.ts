import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import type { Profile } from "../AuthService/AuthService";
import { mapProfileViewModel } from "../AuthService/mapper";
import { mapToIssuerEvent } from "./mapper";

export interface IssuerEvent {
    eventEndDate: Date;
    eventId: string;
    eventLocation: string;
    eventOwnerCredentialId: string;
    eventShortDescription: string;
    eventStartDate: Date;
    eventTitle: string;
    id: string;
    isSigned: boolean;
    issuerCredentialId: string;
    signMessage: string;
    signature: string;
}

export class IssuerService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    public async searchPublicIssuers(
        searchQuery: string,
        limit: number,
        offset: number,
    ): Promise<Profile[]> {
        const response = await this._coreApi.v1.getVerifiedIssuers({
            limit,
            offset,
            search: searchQuery,
        });
        return response.map(mapProfileViewModel);
    }

    public async getTaskedEvents(
        issuerCredentialID: string,
        limit: number = 30,
        offset: number = 0,
    ): Promise<IssuerEvent[]> {
        const response = await this._coreApi.v1.getIssuerEvents({
            issuer_credential_id: issuerCredentialID,
            limit: limit,
            offset: offset,
        });
        return response.map(mapToIssuerEvent);
    }
}

export const defaultIssuerService = new IssuerService(coreApiClient);
