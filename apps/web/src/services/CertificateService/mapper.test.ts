import { describe, it, expect } from "vitest";
import {
    mapCertificate,
    mapCertificateShare,
    mapToMyCertificatesViewModel,
    mapToSignedCertificatesResult,
    mapImportCertificatesResponse,
    mapRevokeCertificatesResponse,
    mapRevokeAllCertificatesResponse,
    mapPublishCertificatesResponse,
    mapClaimCertificateSignMessage,
    mapClaimCertificateResponse,
    mapCreateCertificateShareableLinkResponse,
    mapUpdateCertificateShareResponse,
    mapGetCertificateShareStatusResponse,
    mapGetCertificateShareDataResponse,
} from "./mapper";
import type {
    EntityEventCertificate,
    EntityCertificateShare,
    EntityCertificatePayload,
    CertificateGetMyCertificatesListViewModelResponse,
    CoreApiInternalHandlerEventSignEventCertificatesResponse,
    CoreApiInternalHandlerEventImportCertificateReceiversResponse,
    CoreApiInternalHandlerEventRevokeEventCertificatesResponse,
    CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse,
    CoreApiInternalHandlerEventPublishEventCertificatesResponse,
    CertificateClaimCertificateResponse,
    CertificateShareCertificateShareViewModel,
    CertificateShareCreateCertificateShareResponse,
    CertificateShareUpdateCertificateShareResponse,
    CertificateShareCertificateShareStatusResponse,
    CertificateShareCertificateShareDataResponse,
    EventClaimedCertificateViewModel,
    EventUnclaimedCertificateViewModel,
} from "@decm/api";

const baseCert: EntityEventCertificate = {
    id: "cert-1",
    event_id: "evt-1",
    event_name: "Test Event",
    receiver_credential_id: "cred-1",
    receiver_email: "user@test.com",
    name: "John Doe",
    academic_institution: "MIT",
    certificate_title: "Completion",
    certificate_subtitle: "With honors",
    event_contract_address: "0xEC",
    event_certificate_address: "0xECA",
    certificate_token_id: "tok-1",
    created_at: "2024-01-01T00:00:00Z",
    revoked_at: "2024-06-01T00:00:00Z",
    inbox_message_id: "msg-1",
    broadcasted_at: "2024-06-02T00:00:00Z",
    estimated_deadline: "2024-06-03T00:00:00Z",
    user_claim_signature_id: "sig-1",
};

describe("mapCertificate", () => {
    it("maps snake_case to camelCase", () => {
        const result = mapCertificate(baseCert);

        expect(result.id).toBe("cert-1");
        expect(result.eventId).toBe("evt-1");
        expect(result.eventName).toBe("Test Event");
        expect(result.receiverCredentialId).toBe("cred-1");
        expect(result.receiverEmail).toBe("user@test.com");
        expect(result.name).toBe("John Doe");
        expect(result.academicInstitution).toBe("MIT");
        expect(result.certificateTitle).toBe("Completion");
        expect(result.certificateSubtitle).toBe("With honors");
        expect(result.eventContractAddress).toBe("0xEC");
        expect(result.eventCertificateAddress).toBe("0xECA");
        expect(result.certificateTokenId).toBe("tok-1");
        expect(result.createdAt).toBe("2024-01-01T00:00:00Z");
        expect(result.revokedAt).toBe("2024-06-01T00:00:00Z");
        expect(result.inboxMessageId).toBe("msg-1");
        expect(result.broadcastedAt).toBe("2024-06-02T00:00:00Z");
        expect(result.estimatedDeadline).toBe("2024-06-03T00:00:00Z");
        expect(result.userClaimSignatureId).toBe("sig-1");
    });

    it("handles optional fields as undefined", () => {
        const minimal: EntityEventCertificate = {
            id: "cert-2",
            event_id: "evt-2",
            event_contract_address: "0xEC2",
            created_at: "2024-01-01T00:00:00Z",
        };
        const result = mapCertificate(minimal);

        expect(result.eventName).toBeUndefined();
        expect(result.receiverCredentialId).toBeUndefined();
        expect(result.certificateTokenId).toBeUndefined();
        expect(result.revokedAt).toBeUndefined();
        expect(result.broadcastedAt).toBeUndefined();
        expect(result.estimatedDeadline).toBeUndefined();
        expect(result.userClaimSignatureId).toBeUndefined();
    });
});

describe("mapToMyCertificatesViewModel", () => {
    it("maps claimed and unclaimed certificates", () => {
        const response: CertificateGetMyCertificatesListViewModelResponse = {
            claimed_certificates: [baseCert],
            unclaimed_certificates: [{ ...baseCert, id: "cert-2" }],
            total_claimed: 1,
            total_unclaimed: 1,
        };
        const result = mapToMyCertificatesViewModel(response);

        expect(result.claimedCertificates).toHaveLength(1);
        expect(result.claimedCertificates[0].id).toBe("cert-1");
        expect(result.unclaimedCertificates).toHaveLength(1);
        expect(result.unclaimedCertificates[0].id).toBe("cert-2");
        expect(result.totalClaimed).toBe(1);
        expect(result.totalUnclaimed).toBe(1);
    });

    it("handles null/empty arrays", () => {
        const response = {
            claimed_certificates: null,
            unclaimed_certificates: null,
            total_claimed: 0,
            total_unclaimed: 0,
        } as unknown as CertificateGetMyCertificatesListViewModelResponse;
        const result = mapToMyCertificatesViewModel(response);

        expect(result.claimedCertificates).toEqual([]);
        expect(result.unclaimedCertificates).toEqual([]);
        expect(result.totalClaimed).toBe(0);
        expect(result.totalUnclaimed).toBe(0);
    });
});

describe("mapToSignedCertificatesResult", () => {
    it("maps certificates with signatures", () => {
        const response: CoreApiInternalHandlerEventSignEventCertificatesResponse = {
            certificates: [
                {
                    certificate: { ...baseCert, id: "cert-1", event_id: "evt-1" },
                    signature: "sig-abc",
                },
                {
                    certificate: { ...baseCert, id: "cert-2", event_id: "evt-1" },
                    signature: "sig-def",
                },
            ],
        };
        const result = mapToSignedCertificatesResult(response);

        expect(result.certificates).toHaveLength(2);
        expect(result.certificates[0].certificateId).toBe("cert-1");
        expect(result.certificates[0].signature).toBe("sig-abc");
        expect(result.totalSigned).toBe(2);
    });

    it("handles empty certificates array", () => {
        const response = {
            certificates: null,
        } as unknown as CoreApiInternalHandlerEventSignEventCertificatesResponse;
        const result = mapToSignedCertificatesResult(response);

        expect(result.certificates).toEqual([]);
        expect(result.totalSigned).toBe(0);
    });

    it("handles missing certificate in signature entry", () => {
        const response: CoreApiInternalHandlerEventSignEventCertificatesResponse = {
            certificates: [{ signature: "sig-abc" }],
        } as CoreApiInternalHandlerEventSignEventCertificatesResponse;
        const result = mapToSignedCertificatesResult(response);

        expect(result.certificates[0].certificateId).toBe("");
        expect(result.certificates[0].eventId).toBe("");
    });
});

describe("mapImportCertificatesResponse", () => {
    it("extracts count from certificates array", () => {
        const response: CoreApiInternalHandlerEventImportCertificateReceiversResponse = {
            certificates: [baseCert, { ...baseCert, id: "cert-2" }],
            event_certificate_address: "0xECA",
            event_id: "evt-1",
        };
        const result = mapImportCertificatesResponse(response);

        expect(result.importedCount).toBe(2);
        expect(result.failedCount).toBe(0);
        expect(result.errors).toEqual([]);
    });

    it("handles missing certificates with fallback to 0", () => {
        const response = {
            certificates: null,
            event_certificate_address: "0xECA",
            event_id: "evt-1",
        } as unknown as CoreApiInternalHandlerEventImportCertificateReceiversResponse;
        const result = mapImportCertificatesResponse(response);

        expect(result.importedCount).toBe(0);
    });
});

describe("mapRevokeCertificatesResponse", () => {
    it("extracts count from revoked_certificates array", () => {
        const response: CoreApiInternalHandlerEventRevokeEventCertificatesResponse = {
            revoked_certificates: [baseCert],
        };
        const result = mapRevokeCertificatesResponse(response);

        expect(result.revokedCount).toBe(1);
        expect(result.failedCount).toBe(0);
        expect(result.errors).toEqual([]);
    });

    it("handles missing revoked_certificates with fallback to 0", () => {
        const response = {
            revoked_certificates: null,
        } as unknown as CoreApiInternalHandlerEventRevokeEventCertificatesResponse;
        const result = mapRevokeCertificatesResponse(response);

        expect(result.revokedCount).toBe(0);
    });
});

describe("mapRevokeAllCertificatesResponse", () => {
    it("extracts count and formats message", () => {
        const response: CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse = {
            revoked_certificates: [baseCert, { ...baseCert, id: "cert-2" }],
        };
        const result = mapRevokeAllCertificatesResponse(response);

        expect(result.revokedCount).toBe(2);
        expect(result.message).toBe("Successfully revoked 2 certificates");
    });

    it("handles empty array", () => {
        const response: CoreApiInternalHandlerEventRevokeAllEventCertificatesResponse = {
            revoked_certificates: [],
        };
        const result = mapRevokeAllCertificatesResponse(response);

        expect(result.revokedCount).toBe(0);
        expect(result.message).toBe("Successfully revoked 0 certificates");
    });
});

describe("mapPublishCertificatesResponse", () => {
    it("extracts published count", () => {
        const response: CoreApiInternalHandlerEventPublishEventCertificatesResponse = {
            published_count: 5,
            event_id: "evt-1",
            inbox_messages_created: 5,
            is_published: true,
        };
        const result = mapPublishCertificatesResponse(response);

        expect(result.publishedCount).toBe(5);
        expect(result.failedCount).toBe(0);
        expect(result.errors).toEqual([]);
    });

    it("defaults published count to 0 when missing", () => {
        const response = {
            event_id: "evt-1",
            inbox_messages_created: 0,
            is_published: false,
        } as unknown as CoreApiInternalHandlerEventPublishEventCertificatesResponse;
        const result = mapPublishCertificatesResponse(response);

        expect(result.publishedCount).toBe(0);
    });
});

describe("mapClaimCertificateSignMessage", () => {
    it("extracts sign message", () => {
        const result = mapClaimCertificateSignMessage({
            sign_message: "Please sign this message",
        });

        expect(result.signMessage).toBe("Please sign this message");
        expect(result.certificateId).toBe("");
    });

    it("defaults to empty string when sign_message is empty", () => {
        const result = mapClaimCertificateSignMessage({ sign_message: "" });
        expect(result.signMessage).toBe("");
    });
});

describe("mapClaimCertificateResponse", () => {
    it("extracts certificate details from response", () => {
        const response: CertificateClaimCertificateResponse = {
            certificate: {
                ...baseCert,
                id: "cert-claimed",
                certificate_token_id: "tok-99",
                created_at: "2024-05-01T00:00:00Z",
            },
            message: "Claimed!",
            status: "success",
            user_signature: null,
        };
        const result = mapClaimCertificateResponse(response);

        expect(result.certificateId).toBe("cert-claimed");
        expect(result.transactionHash).toBe("tok-99");
        expect(result.claimedAt).toBe("2024-05-01T00:00:00Z");
        expect(result.message).toBe("Claimed!");
        expect(result.status).toBe("success");
        expect(result.estimatedDeadline).toBeUndefined();
    });

    it("handles null certificate with fallbacks", () => {
        const response: CertificateClaimCertificateResponse = {
            certificate: null,
            message: "",
            status: "success",
            user_signature: null,
        };
        const result = mapClaimCertificateResponse(response);

        expect(result.certificateId).toBe("");
        expect(result.transactionHash).toBe("");
        expect(result.claimedAt).toBe("");
        expect(result.message).toBe("Certificate claimed successfully");
        expect(result.status).toBe("success");
    });

    it("reads estimatedDeadline from user_signature when status is queued", () => {
        const response: CertificateClaimCertificateResponse = {
            certificate: {
                ...baseCert,
                id: "cert-queued",
                certificate_token_id: undefined,
                created_at: "2024-05-01T00:00:00Z",
            },
            message: "Certificate claim has been queued for processing.",
            status: "queued",
            user_signature: { estimated_deadline: "2024-05-02T00:00:00Z" },
        };
        const result = mapClaimCertificateResponse(response);

        expect(result.status).toBe("queued");
        expect(result.estimatedDeadline).toBe("2024-05-02T00:00:00Z");
        expect(result.transactionHash).toBe("");
    });
});

const baseShare: CertificateShareCertificateShareViewModel = {
    id: "share-1",
    event_certificate_id: "cert-1",
    active: true,
    handle: "my-handle",
    has_password: false,
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-02-01T00:00:00Z",
};

describe("mapCertificateShare", () => {
    it("maps all fields from snake_case to camelCase", () => {
        const result = mapCertificateShare(baseShare);

        expect(result.id).toBe("share-1");
        expect(result.eventCertificateId).toBe("cert-1");
        expect(result.active).toBe(true);
        expect(result.handle).toBe("my-handle");
        expect(result.hasPassword).toBe(false);
        expect(result.createdAt).toBe("2024-01-01T00:00:00Z");
        expect(result.updatedAt).toBe("2024-02-01T00:00:00Z");
    });

    it("maps has_password: false to hasPassword: false", () => {
        const result = mapCertificateShare({ ...baseShare, has_password: false });
        expect(result.hasPassword).toBe(false);
    });

    it("maps has_password: true to hasPassword: true", () => {
        const result = mapCertificateShare({ ...baseShare, has_password: true });
        expect(result.hasPassword).toBe(true);
    });
});

describe("mapToMyCertificatesViewModel with typed view models", () => {
    it("maps a claimed certificate including status and certificate_share", () => {
        const claimed: EventClaimedCertificateViewModel = {
            ...baseCert,
            status: "claimed",
            certificate_share: baseShare,
        };
        const response: CertificateGetMyCertificatesListViewModelResponse = {
            claimed_certificates: [claimed],
            unclaimed_certificates: [],
            total_claimed: 1,
            total_unclaimed: 0,
        };
        const result = mapToMyCertificatesViewModel(response);

        expect(result.claimedCertificates).toHaveLength(1);
        const cert = result.claimedCertificates[0];
        expect(cert.id).toBe("cert-1");
        expect(cert.status).toBe("claimed");
        expect(cert.certificateShare).toBeDefined();
        expect(cert.certificateShare?.id).toBe("share-1");
        expect(cert.certificateShare?.eventCertificateId).toBe("cert-1");
        expect(cert.certificateShare?.hasPassword).toBe(false);
    });

    it("maps an unclaimed certificate including certificate_share", () => {
        const unclaimed: EventUnclaimedCertificateViewModel = {
            ...baseCert,
            id: "cert-2",
            certificate_share: { ...baseShare, id: "share-2", event_certificate_id: "cert-2" },
        };
        const response: CertificateGetMyCertificatesListViewModelResponse = {
            claimed_certificates: [],
            unclaimed_certificates: [unclaimed],
            total_claimed: 0,
            total_unclaimed: 1,
        };
        const result = mapToMyCertificatesViewModel(response);

        expect(result.unclaimedCertificates).toHaveLength(1);
        const cert = result.unclaimedCertificates[0];
        expect(cert.id).toBe("cert-2");
        expect(cert.certificateShare).toBeDefined();
        expect(cert.certificateShare?.id).toBe("share-2");
        expect(cert.certificateShare?.eventCertificateId).toBe("cert-2");
    });

    it("maps certificate_share: undefined to certificateShare: undefined", () => {
        const claimed: EventClaimedCertificateViewModel = {
            ...baseCert,
            status: "claimed",
            certificate_share: undefined,
        };
        const unclaimed: EventUnclaimedCertificateViewModel = {
            ...baseCert,
            id: "cert-3",
            certificate_share: undefined,
        };
        const response: CertificateGetMyCertificatesListViewModelResponse = {
            claimed_certificates: [claimed],
            unclaimed_certificates: [unclaimed],
            total_claimed: 1,
            total_unclaimed: 1,
        };
        const result = mapToMyCertificatesViewModel(response);

        expect(result.claimedCertificates[0].certificateShare).toBeUndefined();
        expect(result.unclaimedCertificates[0].certificateShare).toBeUndefined();
    });
});

describe("mapCreateCertificateShareableLinkResponse", () => {
    it("maps share fields to certificateId, isPublished, and handle", () => {
        const entityShare: EntityCertificateShare = {
            id: "share-99",
            event_certificate_id: "cert-42",
            active: true,
            handle: "public-link",
            created_at: "2024-03-01T00:00:00Z",
            updated_at: "2024-03-02T00:00:00Z",
        };
        const response: CertificateShareCreateCertificateShareResponse = {
            share: entityShare,
        };
        const result = mapCreateCertificateShareableLinkResponse(response);

        expect(result.certificateId).toBe("cert-42");
        expect(result.isPublished).toBe(true);
        expect(result.handle).toBe("public-link");
    });

    it("maps share.active: false to isPublished: false", () => {
        const entityShare: EntityCertificateShare = {
            id: "share-100",
            event_certificate_id: "cert-43",
            active: false,
            handle: "draft-link",
            created_at: "2024-03-01T00:00:00Z",
            updated_at: "2024-03-02T00:00:00Z",
        };
        const response: CertificateShareCreateCertificateShareResponse = {
            share: entityShare,
        };
        const result = mapCreateCertificateShareableLinkResponse(response);

        expect(result.isPublished).toBe(false);
        expect(result.certificateId).toBe("cert-43");
        expect(result.handle).toBe("draft-link");
    });
});

describe("mapUpdateCertificateShareResponse", () => {
    it("maps share to certificateId and handle", () => {
        const response: CertificateShareUpdateCertificateShareResponse = {
            share: {
                id: "share-1",
                event_certificate_id: "cert-10",
                active: false,
                handle: "updated-handle",
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-02-01T00:00:00Z",
            },
        };
        const result = mapUpdateCertificateShareResponse(response);

        expect(result.certificateId).toBe("cert-10");
        expect(result.handle).toBe("updated-handle");
    });
});

describe("mapGetCertificateShareStatusResponse", () => {
    it("maps READY status with certificate to camelCase", () => {
        const response: CertificateShareCertificateShareStatusResponse = {
            status: "READY" as never,
            certificate: {
                ...baseCert,
                id: "cert-share-1",
                certificate_title: "Shared Cert",
            },
        };
        const result = mapGetCertificateShareStatusResponse(response);

        expect(result.status).toBe("READY");
        expect(result.certificate?.id).toBe("cert-share-1");
        expect(result.certificate?.certificateTitle).toBe("Shared Cert");
    });

    it("maps PASSWORD_LOCKED status without certificate", () => {
        const response: CertificateShareCertificateShareStatusResponse = {
            status: "PASSWORD_LOCKED" as never,
            certificate: undefined,
        };
        const result = mapGetCertificateShareStatusResponse(response);

        expect(result.status).toBe("PASSWORD_LOCKED");
        expect(result.certificate).toBeUndefined();
    });

    it("maps VALID_BUT_PENDING status", () => {
        const response: CertificateShareCertificateShareStatusResponse = {
            status: "VALID_BUT_PENDING" as never,
            certificate: undefined,
        };
        const result = mapGetCertificateShareStatusResponse(response);

        expect(result.status).toBe("VALID_BUT_PENDING");
        expect(result.certificate).toBeUndefined();
    });
});

describe("mapGetCertificateShareDataResponse", () => {
    it("maps response data to payload", () => {
        const mockPayload: EntityCertificatePayload = {
            header: {
                "@context": ["https://w3.org/ns/credentials/v2"],
                id: "vc-1",
                issuanceDate: "2024-01-01T00:00:00Z",
                issuer: "did:example:issuer",
                type: ["VerifiableCredential"],
            },
            data: {
                certificateId: "cert-1",
                certificateTitle: "Test Cert",
                certificateSubtitle: "Sub",
                certificateTokenId: "tok-1",
                eventName: "Test Event",
                eventDescription: "Desc",
                issuedAt: "2024-01-01T00:00:00Z",
                issuerAddresses: "0xissuer",
                issuerId: "issuer-1",
                receiverAddress: "0xreceiver",
                status: "VALID",
                userId: "user-1",
                backendEncryptedUserData: "",
                encryptedUserData: "",
            },
            proof: {
                hash: "0xhash",
                encryptedByBackendRawData: "",
                encryptedByUserRawData: "",
                signMessage: "sign",
                host: { publicKey: "0xpk", signature: "0xsig" },
                issuers: [],
            },
        };
        const response: CertificateShareCertificateShareDataResponse = {
            data: mockPayload,
        };
        const result = mapGetCertificateShareDataResponse(response);

        expect(result.payload).toBe(mockPayload);
        expect(result.payload.data.certificateTitle).toBe("Test Cert");
        expect(result.payload.header.id).toBe("vc-1");
    });
});
