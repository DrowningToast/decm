import { describe, it, expect, beforeEach } from "vitest";
import { useSearchCertificateNavStore, useCertificateDetailNavStore } from "./certificates";

describe("useSearchCertificateNavStore", () => {
    beforeEach(() => {
        useSearchCertificateNavStore.setState({ searchQuery: "" });
    });

    it("initializes with empty search query", () => {
        expect(useSearchCertificateNavStore.getState().searchQuery).toBe("");
    });

    it("updates search query", () => {
        useSearchCertificateNavStore.getState().setSearchQuery("test");
        expect(useSearchCertificateNavStore.getState().searchQuery).toBe("test");
    });

    it("clears search query", () => {
        useSearchCertificateNavStore.getState().setSearchQuery("test");
        useSearchCertificateNavStore.getState().setSearchQuery("");
        expect(useSearchCertificateNavStore.getState().searchQuery).toBe("");
    });
});

describe("useCertificateDetailNavStore", () => {
    beforeEach(() => {
        useCertificateDetailNavStore.setState({
            certificateId: null,
            isClaimed: false,
            imageUrl: null,
        });
    });

    it("initializes with default values", () => {
        const state = useCertificateDetailNavStore.getState();
        expect(state.certificateId).toBeNull();
        expect(state.isClaimed).toBe(false);
        expect(state.imageUrl).toBeNull();
    });

    it("sets certificate id", () => {
        useCertificateDetailNavStore.getState().setCertificateId("cert-1");
        expect(useCertificateDetailNavStore.getState().certificateId).toBe("cert-1");
    });

    it("clears certificate id", () => {
        useCertificateDetailNavStore.getState().setCertificateId("cert-1");
        useCertificateDetailNavStore.getState().setCertificateId(null);
        expect(useCertificateDetailNavStore.getState().certificateId).toBeNull();
    });

    it("sets isClaimed", () => {
        useCertificateDetailNavStore.getState().setIsClaimed(true);
        expect(useCertificateDetailNavStore.getState().isClaimed).toBe(true);
    });

    it("sets imageUrl", () => {
        useCertificateDetailNavStore.getState().setImageUrl("https://example.com/img.png");
        expect(useCertificateDetailNavStore.getState().imageUrl).toBe(
            "https://example.com/img.png",
        );
    });

    it("clears imageUrl", () => {
        useCertificateDetailNavStore.getState().setImageUrl("https://example.com/img.png");
        useCertificateDetailNavStore.getState().setImageUrl(null);
        expect(useCertificateDetailNavStore.getState().imageUrl).toBeNull();
    });
});
