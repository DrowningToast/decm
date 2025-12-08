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

export interface IssuerEventViewModel {
    id: string;
    eventId: string;
    eventTitle: string;
    eventShortDescription: string;
    eventStartDate: string;
    eventEndDate: string;
    eventLocation: string;
    eventOwnerCredentialId: string;
    ownerName: string;
    ownerWalletAddress: string;
    ownerEmail: string;
    ownerGoogleEmail: string;
    certificateCount: number;
    issuerCredentialId: string;
    isSigned: boolean;
    createdAt: string;
    updatedAt: string;
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

    public async getIssuerEventsViewModel(
        issuerCredentialID: string,
        limit: number = 30,
        offset: number = 0,
    ): Promise<IssuerEventViewModel[]> {
        const response = await this._coreApi.v1.getIssuerEventsViewmodel({
            issuer_credential_id: issuerCredentialID,
            limit: limit,
            offset: offset,
        });
        return response.map((item) => ({
            id: item.id,
            eventId: item.event_id,
            eventTitle: item.event_title,
            eventShortDescription: item.event_short_description,
            eventStartDate: item.event_start_date,
            eventEndDate: item.event_end_date,
            eventLocation: item.event_location,
            eventOwnerCredentialId: item.event_owner_credential_id,
            ownerName: item.owner_name,
            ownerWalletAddress: item.owner_wallet_address,
            ownerEmail: item.owner_email,
            ownerGoogleEmail: item.owner_google_email,
            certificateCount: item.certificate_count,
            issuerCredentialId: item.issuer_credential_id,
            isSigned: item.is_signed,
            createdAt: item.created_at,
            updatedAt: item.updated_at,
        }));
    }
}

export const defaultIssuerService = new IssuerService(coreApiClient);
