import { describe, it, expect, vi, beforeEach } from "vitest";
import type { CoreApiType } from "@/lib/api/api";
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
            const mockCreateObjectURL = vi
                .spyOn(URL, "createObjectURL")
                .mockReturnValue("blob:url");

            vi.mocked(mockCoreApi.v1.generateCertificateImage).mockResolvedValue(mockBlob);

            const result = await certificateService.getCertificateImage("cert-123");

            expect(mockCoreApi.v1.generateCertificateImage).toHaveBeenCalledWith({
                certificateId: "cert-123",
            });
            expect(result.url).toBe("blob:url");
            expect(result.blob).toBe(mockBlob);
            expect(result.contentType).toBe("image/png");
            expect(mockCreateObjectURL).toHaveBeenCalledWith(mockBlob);

            mockCreateObjectURL.mockRestore();
        });

        it("should throw error when response is not a Blob", async () => {
            vi.mocked(mockCoreApi.v1.generateCertificateImage).mockResolvedValue({} as Blob);

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
            vi.mocked(mockCoreApi.v1.getEventCertificates).mockResolvedValue({});

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
            vi.mocked(mockCoreApi.v1.getEventCertificateFontFamilies).mockResolvedValue({});

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
                certificates: [
                    { id: "cert-1" },
                    { id: "cert-2" },
                    { id: "cert-3" },
                    { id: "cert-4" },
                    { id: "cert-5" },
                ],
            };

            const params = {
                eventId: "event-1",
                hostPin: "1234",
                receivers: [
                    {
                        receiver_credential_id: "cred-1",
                        receiver_email: "test@example.com",
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
                revoked_certificates: [{ id: "cert-1" }, { id: "cert-2" }],
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
                revoked_certificates: Array.from({ length: 10 }, (_, i) => ({ id: `cert-${i}` })),
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
                published_count: 5,
                failed_count: 0,
                errors: [],
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
            const mockResponse = { success: true };
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
                id: "cert-123",
                certificate_token_id: "0xabc",
                created_at: "2024-01-01T00:00:00Z",
            };

            vi.mocked(mockCoreApi.v1.claimCertificate).mockResolvedValue(mockResponse);

            const result = await certificateService.claimCertificateWithPin({
                certificateId: "cert-123",
                accountPassword: "password123",
            });

            expect(mockCoreApi.v1.claimCertificate).toHaveBeenCalledWith(
                { certificateId: "cert-123" },
                { account_password: "password123" },
            );
            expect(result.certificateId).toBe("cert-123");
            expect(result.transactionHash).toBe("0xabc");
        });
    });

    describe("claimCertificateWithSignature", () => {
        it("should claim certificate with signature", async () => {
            const mockResponse = {
                id: "cert-123",
                certificate_token_id: "0xabc",
                created_at: "2024-01-01T00:00:00Z",
            };

            vi.mocked(mockCoreApi.v1.claimCertificate).mockResolvedValue(mockResponse);

            const result = await certificateService.claimCertificateWithSignature({
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
                id: "cert-123",
                certificate_token_id: "0xabc",
                created_at: "2024-01-01T00:00:00Z",
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
                id: "cert-123",
                certificate_token_id: "0xabc",
                created_at: "2024-01-01T00:00:00Z",
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
});
