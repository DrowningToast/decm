import { describe, it, expect, vi, beforeEach } from "vitest";
import type { CoreApiType } from "@/lib/api/api";
import { IssuerService } from "./IssuerService";

// Mock the coreApiClient
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            getVerifiedIssuers: vi.fn(),
            getIssuerEvents: vi.fn(),
            getIssuerEventsViewmodel: vi.fn(),
        },
    },
}));

import { coreApiClient } from "@/lib/api/api";

describe("IssuerService", () => {
    let issuerService: IssuerService;
    let mockCoreApi: CoreApiType;

    beforeEach(() => {
        vi.clearAllMocks();
        mockCoreApi = coreApiClient as unknown as CoreApiType;
        issuerService = new IssuerService(mockCoreApi);
    });

    describe("searchPublicIssuers", () => {
        it("should search public issuers with query", async () => {
            const mockResponse = [
                {
                    id: "cred-1",
                    credential_id: "cred-1",
                    authentication_credential_id: "auth-cred-1",
                    wallet_address: "0x123",
                    email: "issuer@example.com",
                    google_email: "issuer@gmail.com",
                    first_name: "John",
                    last_name: "Doe",
                    solution_status: "BYOK",
                    is_academic_email_public: false,
                    is_academic_institution_public: false,
                    is_address_public: false,
                    is_bio_public: false,
                    is_email_public: true,
                    is_first_name_public: true,
                    is_last_name_public: true,
                    is_phone_number_public: false,
                    is_profile_picture_public: true,
                    is_verified: true,
                    is_wallet_address_public: true,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            ];

            vi.mocked(mockCoreApi.v1.getVerifiedIssuers).mockResolvedValue(mockResponse);

            const result = await issuerService.searchPublicIssuers("John", 10, 0);

            expect(mockCoreApi.v1.getVerifiedIssuers).toHaveBeenCalledWith({
                limit: 10,
                offset: 0,
                search: "John",
            });
            expect(result).toHaveLength(1);
            expect(result[0].walletAddress).toBe("0x123");
            expect(result[0].email).toBe("issuer@example.com");
        });

        it("should use default limit and offset when not provided", async () => {
            const mockResponse: never[] = [];
            vi.mocked(mockCoreApi.v1.getVerifiedIssuers).mockResolvedValue(mockResponse);

            await issuerService.searchPublicIssuers("test", 30, 0);

            expect(mockCoreApi.v1.getVerifiedIssuers).toHaveBeenCalledWith({
                limit: 30,
                offset: 0,
                search: "test",
            });
        });
    });

    describe("getTaskedEvents", () => {
        it("should fetch tasked events for issuer", async () => {
            const mockResponse = [
                {
                    id: "issuer-event-1",
                    event_id: "event-1",
                    issuer_credential_id: "cred-1",
                    is_signed: 0,
                    event_title: "Test Event",
                    event_short_description: "Description",
                    event_start_date: "2024-01-01T00:00:00Z",
                    event_end_date: "2024-01-02T00:00:00Z",
                    event_location: "Location",
                    event_owner_credential_id: "owner-cred-1",
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            ];

            vi.mocked(mockCoreApi.v1.getIssuerEvents).mockResolvedValue(mockResponse);

            const result = await issuerService.getTaskedEvents("cred-1", 30, 0);

            expect(mockCoreApi.v1.getIssuerEvents).toHaveBeenCalledWith({
                issuer_credential_id: "cred-1",
                limit: 30,
                offset: 0,
            });
            expect(result).toHaveLength(1);
            expect(result[0].eventId).toBe("event-1");
            expect(result[0].issuerCredentialId).toBe("cred-1");
            expect(result[0].isSigned).toBe(false);
        });

        it("should use default limit and offset", async () => {
            const mockResponse: never[] = [];
            vi.mocked(mockCoreApi.v1.getIssuerEvents).mockResolvedValue(mockResponse);

            await issuerService.getTaskedEvents("cred-1");

            expect(mockCoreApi.v1.getIssuerEvents).toHaveBeenCalledWith({
                issuer_credential_id: "cred-1",
                limit: 30,
                offset: 0,
            });
        });
    });

    describe("getIssuerEventsViewModel", () => {
        it("should fetch issuer events view model", async () => {
            const mockResponse = [
                {
                    id: "issuer-event-1",
                    event_id: "event-1",
                    event_title: "Test Event",
                    event_short_description: "Description",
                    event_start_date: "2024-01-01T00:00:00Z",
                    event_end_date: "2024-01-02T00:00:00Z",
                    event_location: "Location",
                    event_owner_credential_id: "owner-cred-1",
                    owner_name: "Owner Name",
                    owner_wallet_address: "0xowner",
                    owner_email: "owner@example.com",
                    owner_google_email: "owner@gmail.com",
                    certificate_count: 10,
                    issuer_credential_id: "cred-1",
                    is_signed: true,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            ];

            vi.mocked(mockCoreApi.v1.getIssuerEventsViewmodel).mockResolvedValue(mockResponse);

            const result = await issuerService.getIssuerEventsViewModel("cred-1", 30, 0);

            expect(mockCoreApi.v1.getIssuerEventsViewmodel).toHaveBeenCalledWith({
                issuer_credential_id: "cred-1",
                limit: 30,
                offset: 0,
            });
            expect(result).toHaveLength(1);
            expect(result[0].id).toBe("issuer-event-1");
            expect(result[0].eventId).toBe("event-1");
            expect(result[0].eventTitle).toBe("Test Event");
            expect(result[0].ownerName).toBe("Owner Name");
            expect(result[0].certificateCount).toBe(10);
            expect(result[0].isSigned).toBe(true);
        });

        it("should use default limit and offset", async () => {
            const mockResponse: never[] = [];
            vi.mocked(mockCoreApi.v1.getIssuerEventsViewmodel).mockResolvedValue(mockResponse);

            await issuerService.getIssuerEventsViewModel("cred-1");

            expect(mockCoreApi.v1.getIssuerEventsViewmodel).toHaveBeenCalledWith({
                issuer_credential_id: "cred-1",
                limit: 30,
                offset: 0,
            });
        });
    });
});
