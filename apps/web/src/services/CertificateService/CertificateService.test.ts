import { describe, it, expect, vi, beforeEach } from "vitest";
import type { CoreApiType } from "@/lib/api/api";
import { CertificateShareCertificateShareViewStatus } from "@decm/api";
import { CertificateService } from "./CertificateService";

// Mock the coreApiClient
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            generateCertificateImage: vi.fn(),
            getMyCertificatesListViewmodel: vi.fn(),
            getEventCertificates: vi.fn(),
            getEventCertificateFontFamilies: vi.fn(),
            signEventCertificates: vi.fn(),
            importCertificateReceivers: vi.fn(),
            revokeEventCertificates: vi.fn(),
            revokeAllEventCertificates: vi.fn(),
            publishEventCertificates: vi.fn(),
            toggleCertificatePublished: vi.fn(),
            checkCertificateMintReadiness: vi.fn(),
            getClaimCertificateSignMessage: vi.fn(),
            claimCertificate: vi.fn(),
            createCertificateShare: vi.fn(),
            updateCertificateShare: vi.fn(),
            getCertificateShareStatus: vi.fn(),
            getCertificateShareData: vi.fn(),
            getCertificateShareDataWithPassword: vi.fn(),
        },
    },
}));

import { coreApiClient } from "@/lib/api/api";

describe("CertificateService", () => {
    let certificateService: CertificateService;
    let mockCoreApi: CoreApiType;

    beforeEach(() => {
        vi.clearAllMocks();
        mockCoreApi = coreApiClient as unknown as CoreApiType;
        certificateService = new CertificateService(mockCoreApi);
    });

    describe("getCertificateImage", () => {
        it("should fetch and map certificate image successfully", async () => {
            const mockBlob = new Blob(["image data"], { type: "image/png" });
            const mockFile = new File([mockBlob], "certificate.png", { type: "image/png" });
            const mockCreateObjectURL = vi
                .spyOn(URL, "createObjectURL")
                .mockReturnValue("blob:url");

            vi.mocked(mockCoreApi.v1.generateCertificateImage).mockResolvedValue(mockFile);

            const result = await certificateService.getCertificateImage("cert-123");

            expect(mockCoreApi.v1.generateCertificateImage).toHaveBeenCalledWith({
                certificateId: "cert-123",
            });
            expect(result.url).toBe("blob:url");
            expect(result.blob).toBe(mockFile);
            expect(result.contentType).toBe("image/png");
            expect(mockCreateObjectURL).toHaveBeenCalledWith(mockFile);

            mockCreateObjectURL.mockRestore();
        });

        it("should throw error when response is not a Blob", async () => {
            vi.mocked(mockCoreApi.v1.generateCertificateImage).mockResolvedValue({} as File);

            await expect(certificateService.getCertificateImage("cert-123")).rejects.toThrow(
                "Invalid response: expected Blob",
            );
        });
    });

    describe("getMyCertificatesList", () => {
        it("should fetch and map my certificates list", async () => {
            const mockResponse = {
                claimed_certificates: [
                    {
                        id: "cert-1",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                        status: "completed",
                    },
                ],
                unclaimed_certificates: [
                    {
                        id: "cert-2",
                        event_id: "event-2",
                        event_contract_address: "0x456",
                        created_at: "2024-01-02T00:00:00Z",
                    },
                ],
                total_claimed: 1,
                total_unclaimed: 1,
            };

            vi.mocked(mockCoreApi.v1.getMyCertificatesListViewmodel).mockResolvedValue(
                mockResponse,
            );

            const result = await certificateService.getMyCertificatesList();

            expect(mockCoreApi.v1.getMyCertificatesListViewmodel).toHaveBeenCalled();
            expect(result.claimedCertificates).toHaveLength(1);
            expect(result.unclaimedCertificates).toHaveLength(1);
            expect(result.totalClaimed).toBe(1);
            expect(result.totalUnclaimed).toBe(1);
            expect(result.claimedCertificates[0].id).toBe("cert-1");
            expect(result.claimedCertificates[0].eventId).toBe("event-1");
        });
    });

    describe("getEventCertificates", () => {
        it("should fetch and map event certificates", async () => {
            const mockResponse = {
                certificates: [
                    {
                        id: "cert-1",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                    },
                ],
            };

            vi.mocked(mockCoreApi.v1.getEventCertificates).mockResolvedValue(mockResponse);

            const result = await certificateService.getEventCertificates("event-1");

            expect(mockCoreApi.v1.getEventCertificates).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result).toHaveLength(1);
            expect(result[0].id).toBe("cert-1");
            expect(result[0].eventId).toBe("event-1");
        });

        it("should return empty array when certificates is undefined", async () => {
            vi.mocked(mockCoreApi.v1.getEventCertificates).mockResolvedValue({
                certificates: [],
            });

            const result = await certificateService.getEventCertificates("event-1");

            expect(result).toEqual([]);
        });
    });

    describe("getFontFamilies", () => {
        it("should fetch font families", async () => {
            const mockResponse = {
                font_families: [
                    {
                        id: 1,
                        font_family_name: "Arial",
                        css_font_name: "Arial, sans-serif",
                        is_default: true,
                        available_font_weights: [400, 700],
                        is_support_italic: true,
                    },
                ],
            };

            vi.mocked(mockCoreApi.v1.getEventCertificateFontFamilies).mockResolvedValue(
                mockResponse,
            );

            const result = await certificateService.getFontFamilies();

            expect(mockCoreApi.v1.getEventCertificateFontFamilies).toHaveBeenCalled();
            expect(result).toHaveLength(1);
            expect(result[0].id).toBe(1);
            expect(result[0].font_family_name).toBe("Arial");
        });

        it("should return empty array when font_families is undefined", async () => {
            vi.mocked(mockCoreApi.v1.getEventCertificateFontFamilies).mockResolvedValue({
                font_families: [],
            });

            const result = await certificateService.getFontFamilies();

            expect(result).toEqual([]);
        });
    });

    describe("signEventCertificates", () => {
        it("should sign event certificates with PIN", async () => {
            const mockResponse = {
                certificates: [
                    {
                        certificate: {
                            id: "cert-1",
                            event_id: "event-1",
                            event_contract_address: "0x123",
                            created_at: "2024-01-01T00:00:00Z",
                        },
                        signature: "sig-123",
                    },
                ],
            };

            vi.mocked(mockCoreApi.v1.signEventCertificates).mockResolvedValue(mockResponse);

            const result = await certificateService.signEventCertificates("event-1", "1234");

            expect(mockCoreApi.v1.signEventCertificates).toHaveBeenCalledWith(
                { eventId: "event-1" },
                { issuer_pin: "1234" },
            );
            expect(result.certificates).toHaveLength(1);
            expect(result.totalSigned).toBe(1);
            expect(result.certificates[0].certificateId).toBe("cert-1");
        });
    });

    describe("importCertificates", () => {
        it("should import certificate receivers", async () => {
            const mockResponse = {
                event_id: "event-1",
                event_certificate_address: "0x123",
                certificates: [
                    {
                        id: "cert-1",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                    },
                    {
                        id: "cert-2",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                    },
                    {
                        id: "cert-3",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                    },
                    {
                        id: "cert-4",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                    },
                    {
                        id: "cert-5",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                    },
                ],
            };

            const params = {
                eventId: "event-1",
                hostPin: "1234",
                receivers: [
                    {
                        academic_institution: "Test University",
                        certificate_subtitle: "Subtitle",
                        certificate_title: "Title",
                        email: "test@example.com",
                        first_name: "John",
                        last_name: "Doe",
                    },
                ],
            };

            vi.mocked(mockCoreApi.v1.importCertificateReceivers).mockResolvedValue(mockResponse);

            const result = await certificateService.importCertificates(params);

            expect(mockCoreApi.v1.importCertificateReceivers).toHaveBeenCalledWith(
                { eventId: "event-1" },
                {
                    event_id: "event-1",
                    host_pin: "1234",
                    receivers: params.receivers,
                },
            );
            expect(result.importedCount).toBe(5);
            expect(result.failedCount).toBe(0);
        });
    });

    describe("revokeCertificates", () => {
        it("should revoke specific certificates", async () => {
            const mockResponse = {
                revoked_certificates: [
                    {
                        id: "cert-1",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                    },
                    {
                        id: "cert-2",
                        event_id: "event-1",
                        event_contract_address: "0x123",
                        created_at: "2024-01-01T00:00:00Z",
                    },
                ],
            };

            const params = {
                eventId: "event-1",
                certificateIds: ["cert-1", "cert-2"],
            };

            vi.mocked(mockCoreApi.v1.revokeEventCertificates).mockResolvedValue(mockResponse);

            const result = await certificateService.revokeCertificates(params);

            expect(mockCoreApi.v1.revokeEventCertificates).toHaveBeenCalledWith(
                { eventId: "event-1" },
                { certificate_ids: ["cert-1", "cert-2"] },
            );
            expect(result.revokedCount).toBe(2);
        });
    });

    describe("revokeAllCertificates", () => {
        it("should revoke all certificates for an event", async () => {
            const mockResponse = {
                revoked_certificates: Array.from({ length: 10 }, (_, i) => ({
                    id: `cert-${i}`,
                    event_id: "event-1",
                    event_contract_address: "0x123",
                    created_at: "2024-01-01T00:00:00Z",
                })),
            };

            vi.mocked(mockCoreApi.v1.revokeAllEventCertificates).mockResolvedValue(mockResponse);

            const result = await certificateService.revokeAllCertificates("event-1");

            expect(mockCoreApi.v1.revokeAllEventCertificates).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result.revokedCount).toBe(10);
            expect(result.message).toBe("Successfully revoked 10 certificates");
        });
    });

    describe("publishCertificates", () => {
        it("should publish event certificates", async () => {
            const mockResponse = {
                event_id: "event-1",
                published_count: 5,
                failed_count: 0,
                errors: [],
                is_published: true,
                inbox_messages_created: 5,
            };

            vi.mocked(mockCoreApi.v1.publishEventCertificates).mockResolvedValue(mockResponse);

            const result = await certificateService.publishCertificates("event-1");

            expect(mockCoreApi.v1.publishEventCertificates).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result.publishedCount).toBe(5);
        });
    });

    describe("toggleCertificatePublished", () => {
        it("should toggle certificate published status", async () => {
            const mockResponse = {
                id: "config-1",
                event_id: "event-1",
                base_certificate_presigned_url: "https://example.com/cert.png",
                base_certificate_storage_key: "certs/base.png",
                event_name_pos_x: 100,
                event_name_pos_y: 200,
                name_pos_x: 150,
                name_pos_y: 250,
                is_published: true,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
            };
            vi.mocked(mockCoreApi.v1.toggleCertificatePublished).mockResolvedValue(mockResponse);

            const result = await certificateService.toggleCertificatePublished("event-1", true);

            expect(mockCoreApi.v1.toggleCertificatePublished).toHaveBeenCalledWith(
                { eventId: "event-1" },
                { is_published: true },
            );
            expect(result).toEqual(mockResponse);
        });
    });

    describe("checkCertificateMintReadiness", () => {
        it("should check certificate mint readiness", async () => {
            const mockResponse = {
                is_ready: true,
                missing_requirements: [],
                has_certificate_config: true,
                has_certificate_contract: true,
                total_issuers_count: 1,
                all_issuers_have_signed: true,
                signed_issuers_count: 1,
            };

            vi.mocked(mockCoreApi.v1.checkCertificateMintReadiness).mockResolvedValue(mockResponse);

            const result = await certificateService.checkCertificateMintReadiness("event-1");

            expect(mockCoreApi.v1.checkCertificateMintReadiness).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result.isReady).toBe(true);
            expect(result.readinessInfo.hasBaseCertificateImage).toBe(true);
            expect(result.readinessInfo.hasEventContract).toBe(true);
        });

        it("should handle missing readiness info gracefully", async () => {
            const mockResponse = {
                is_ready: false,
                missing_requirements: ["Missing contract"],
                has_certificate_config: false,
                has_certificate_contract: false,
                total_issuers_count: 0,
                all_issuers_have_signed: false,
                signed_issuers_count: 0,
            };

            vi.mocked(mockCoreApi.v1.checkCertificateMintReadiness).mockResolvedValue(mockResponse);

            const result = await certificateService.checkCertificateMintReadiness("event-1");

            expect(result.isReady).toBe(false);
            expect(result.readinessInfo.hasBaseCertificateImage).toBe(false);
            expect(result.missingRequirements).toEqual(["Missing contract"]);
        });
    });

    describe("getClaimCertificateSignMessage", () => {
        it("should get claim certificate sign message", async () => {
            const mockResponse = {
                sign_message: "Sign this message",
            };

            vi.mocked(mockCoreApi.v1.getClaimCertificateSignMessage).mockResolvedValue(
                mockResponse,
            );

            const result = await certificateService.getClaimCertificateSignMessage("cert-123");

            expect(mockCoreApi.v1.getClaimCertificateSignMessage).toHaveBeenCalledWith({
                certificateId: "cert-123",
            });
            expect(result.signMessage).toBe("Sign this message");
            expect(result.certificateId).toBe(""); // Certificate ID is not in the response
        });
    });

    describe("claimCertificateWithPin", () => {
        it("should claim certificate with PIN", async () => {
            const mockResponse = {
                certificate: {
                    id: "cert-123",
                    certificate_token_id: "0xabc",
                    created_at: "2024-01-01T00:00:00Z",
                    event_contract_address: "0x123",
                    event_id: "event-1",
                },
                message: "Certificate claimed successfully",
                status: "success",
                user_signature: null,
            };

            vi.mocked(mockCoreApi.v1.claimCertificate).mockResolvedValue(mockResponse);

            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const result = await (certificateService as any).claimCertificateWithPin({
                certificateId: "cert-123",
                accountPassword: "password123",
            });

            expect(mockCoreApi.v1.claimCertificate).toHaveBeenCalledWith(
                { certificateId: "cert-123" },
                { account_password: "password123" },
            );
            expect(result.certificateId).toBe("cert-123");
            expect(result.transactionHash).toBe("0xabc");
            expect(result.status).toBe("success");
            expect(result.estimatedDeadline).toBeUndefined();
        });

        it("should return status=queued and estimatedDeadline when backend queues the claim", async () => {
            const mockResponse = {
                certificate: {
                    id: "cert-123",
                    certificate_token_id: undefined,
                    created_at: "2024-01-01T00:00:00Z",
                    event_contract_address: "0x123",
                    event_id: "event-1",
                },
                message: "Certificate claim has been queued for processing.",
                status: "queued",
                user_signature: { estimated_deadline: "2024-01-02T00:00:00Z" },
            };

            vi.mocked(mockCoreApi.v1.claimCertificate).mockResolvedValue(mockResponse);

            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const result = await (certificateService as any).claimCertificateWithPin({
                certificateId: "cert-123",
                accountPassword: "password123",
            });

            expect(result.status).toBe("queued");
            expect(result.estimatedDeadline).toBe("2024-01-02T00:00:00Z");
            expect(result.transactionHash).toBe("");
        });
    });

    describe("claimCertificateWithSignature", () => {
        it("should claim certificate with signature", async () => {
            const mockResponse = {
                certificate: {
                    id: "cert-123",
                    certificate_token_id: "0xabc",
                    created_at: "2024-01-01T00:00:00Z",
                    event_contract_address: "0x123",
                    event_id: "event-1",
                },
                message: "Certificate claimed successfully",
                status: "success",
                user_signature: null,
            };

            vi.mocked(mockCoreApi.v1.claimCertificate).mockResolvedValue(mockResponse);

            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const result = await (certificateService as any).claimCertificateWithSignature({
                certificateId: "cert-123",
                signature: "sig-123",
                signMessage: "Sign this",
            });

            expect(mockCoreApi.v1.claimCertificate).toHaveBeenCalledWith(
                { certificateId: "cert-123" },
                {
                    signature: "sig-123",
                    sign_message: "Sign this",
                },
            );
            expect(result.certificateId).toBe("cert-123");
        });
    });

    describe("claimCertificate", () => {
        it("should claim certificate with PIN when accountPassword is provided", async () => {
            const mockResponse = {
                certificate: {
                    id: "cert-123",
                    certificate_token_id: "0xabc",
                    created_at: "2024-01-01T00:00:00Z",
                    event_contract_address: "0x123",
                    event_id: "event-1",
                },
                message: "Certificate claimed successfully",
                status: "success",
                user_signature: null,
            };

            vi.mocked(mockCoreApi.v1.claimCertificate).mockResolvedValue(mockResponse);

            const result = await certificateService.claimCertificate({
                certificateId: "cert-123",
                accountPassword: "password123",
            });

            expect(mockCoreApi.v1.claimCertificate).toHaveBeenCalledWith(
                { certificateId: "cert-123" },
                { account_password: "password123" },
            );
            expect(result.certificateId).toBe("cert-123");
        });

        it("should claim certificate with signature when signature is provided", async () => {
            const mockResponse = {
                certificate: {
                    id: "cert-123",
                    certificate_token_id: "0xabc",
                    created_at: "2024-01-01T00:00:00Z",
                    event_contract_address: "0x123",
                    event_id: "event-1",
                },
                message: "Certificate claimed successfully",
                status: "success",
                user_signature: null,
            };

            vi.mocked(mockCoreApi.v1.claimCertificate).mockResolvedValue(mockResponse);

            const result = await certificateService.claimCertificate({
                certificateId: "cert-123",
                signature: "sig-123",
                signMessage: "Sign this",
            });

            expect(mockCoreApi.v1.claimCertificate).toHaveBeenCalledWith(
                { certificateId: "cert-123" },
                {
                    signature: "sig-123",
                    sign_message: "Sign this",
                },
            );
            expect(result.certificateId).toBe("cert-123");
        });
    });

    describe("createCertificateShareableLink", () => {
        it("should create a share link without password", async () => {
            const mockResponse = {
                share: {
                    id: "share-1",
                    event_certificate_id: "cert-1",
                    active: false,
                    handle: "abc123",
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            };

            vi.mocked(mockCoreApi.v1.createCertificateShare).mockResolvedValue(mockResponse);

            const result = await certificateService.createCertificateShareableLink("cert-1");

            expect(mockCoreApi.v1.createCertificateShare).toHaveBeenCalledWith(
                { certificateId: "cert-1" },
                { password: undefined },
            );
            expect(result.certificateId).toBe("cert-1");
            expect(result.handle).toBe("abc123");
            expect(result.isPublished).toBe(false);
        });

        it("should create a share link with password", async () => {
            const mockResponse = {
                share: {
                    id: "share-2",
                    event_certificate_id: "cert-2",
                    active: false,
                    handle: "def456",
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            };

            vi.mocked(mockCoreApi.v1.createCertificateShare).mockResolvedValue(mockResponse);

            const result = await certificateService.createCertificateShareableLink(
                "cert-2",
                "secret",
            );

            expect(mockCoreApi.v1.createCertificateShare).toHaveBeenCalledWith(
                { certificateId: "cert-2" },
                { password: "secret" },
            );
            expect(result.handle).toBe("def456");
        });
    });

    describe("updateCertificateShare", () => {
        it("should update share password", async () => {
            const mockResponse = {
                share: {
                    id: "share-1",
                    event_certificate_id: "cert-1",
                    active: false,
                    handle: "abc123",
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-02T00:00:00Z",
                },
            };

            vi.mocked(mockCoreApi.v1.updateCertificateShare).mockResolvedValue(mockResponse);

            const result = await certificateService.updateCertificateShare("share-1", {
                password: "newpass",
            });

            expect(mockCoreApi.v1.updateCertificateShare).toHaveBeenCalledWith(
                { shareId: "share-1" },
                { password: "newpass" },
            );
            expect(result.certificateId).toBe("cert-1");
            expect(result.handle).toBe("abc123");
        });
    });

    describe("getCertificateShareStatus", () => {
        it("should return READY status with certificate", async () => {
            const mockResponse = {
                status: CertificateShareCertificateShareViewStatus.CertificateShareViewStatusReady,
                certificate: {
                    id: "cert-1",
                    event_id: "event-1",
                    event_contract_address: "0x123",
                    created_at: "2024-01-01T00:00:00Z",
                    certificate_title: "Test Certificate",
                },
            };

            vi.mocked(mockCoreApi.v1.getCertificateShareStatus).mockResolvedValue(mockResponse);

            const result = await certificateService.getCertificateShareStatus("handle-abc");

            expect(mockCoreApi.v1.getCertificateShareStatus).toHaveBeenCalledWith({
                handle: "handle-abc",
            });
            expect(result.status).toBe("READY");
            expect(result.certificate?.id).toBe("cert-1");
            expect(result.certificate?.certificateTitle).toBe("Test Certificate");
        });

        it("should return PASSWORD_LOCKED status without certificate", async () => {
            const mockResponse = {
                status: CertificateShareCertificateShareViewStatus.CertificateShareViewStatusPasswordLocked,
                certificate: undefined,
            };

            vi.mocked(mockCoreApi.v1.getCertificateShareStatus).mockResolvedValue(mockResponse);

            const result = await certificateService.getCertificateShareStatus("handle-abc");

            expect(result.status).toBe("PASSWORD_LOCKED");
            expect(result.certificate).toBeUndefined();
        });

        it("should return VALID_BUT_PENDING status", async () => {
            const mockResponse = {
                status: CertificateShareCertificateShareViewStatus.CertificateShareViewStatusValidButPending,
                certificate: undefined,
            };

            vi.mocked(mockCoreApi.v1.getCertificateShareStatus).mockResolvedValue(mockResponse);

            const result = await certificateService.getCertificateShareStatus("handle-abc");

            expect(result.status).toBe("VALID_BUT_PENDING");
        });
    });

    describe("getCertificateShareData", () => {
        it("should fetch and return on-chain VC data", async () => {
            const mockPayload = {
                header: {
                    "@context": ["https://w3.org/ns/credentials/v2"],
                    id: "vc-1",
                    issuanceDate: "2024-01-01T00:00:00Z",
                    issuer: "did:example:issuer",
                    type: ["VerifiableCredential"],
                },
                data: {
                    certificateId: "cert-1",
                    certificateTitle: "Test",
                    status: "VALID",
                } as never,
                proof: {
                    hash: "0xhash",
                    encryptedByBackendRawData: "",
                    encryptedByUserRawData: "",
                    signMessage: "",
                    host: { publicKey: "0xpk", signature: "0xsig" },
                    issuers: [],
                },
            };
            const mockResponse = { data: mockPayload };

            vi.mocked(mockCoreApi.v1.getCertificateShareData).mockResolvedValue(mockResponse);

            const result = await certificateService.getCertificateShareData("handle-abc");

            expect(mockCoreApi.v1.getCertificateShareData).toHaveBeenCalledWith({
                handle: "handle-abc",
            });
            expect(result.payload).toBe(mockPayload);
        });
    });

    describe("getCertificateShareDataWithPassword", () => {
        it("should fetch on-chain VC data with password", async () => {
            const mockPayload = {
                header: {
                    "@context": ["https://w3.org/ns/credentials/v2"],
                    id: "vc-2",
                    issuanceDate: "2024-01-01T00:00:00Z",
                    issuer: "did:example:issuer",
                    type: ["VerifiableCredential"],
                },
                data: {
                    certificateId: "cert-2",
                    certificateTitle: "Protected",
                    status: "VALID",
                } as never,
                proof: {
                    hash: "0xhash2",
                    encryptedByBackendRawData: "",
                    encryptedByUserRawData: "",
                    signMessage: "",
                    host: { publicKey: "0xpk", signature: "0xsig" },
                    issuers: [],
                },
            };
            const mockResponse = { data: mockPayload };

            vi.mocked(mockCoreApi.v1.getCertificateShareDataWithPassword).mockResolvedValue(
                mockResponse,
            );

            const result = await certificateService.getCertificateShareDataWithPassword(
                "handle-abc",
                "secret",
            );

            expect(mockCoreApi.v1.getCertificateShareDataWithPassword).toHaveBeenCalledWith(
                { handle: "handle-abc" },
                { password: "secret" },
            );
            expect(result.payload).toBe(mockPayload);
        });
    });
});
