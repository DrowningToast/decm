import { describe, it, expect, vi, beforeEach } from "vitest";
import {
    useSearchCertificateNavStore,
    useCertificateDetailNavStore,
    useCertificateDetailsSharedNavStore,
} from "./certificates";

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

    it("onClickShareable is a no-op function by default", () => {
        expect(() => useCertificateDetailNavStore.getState().onClickShareable()).not.toThrow();
    });

    it("isShareableLoading defaults to false", () => {
        expect(useCertificateDetailNavStore.getState().isShareableLoading).toBe(false);
    });

    it("after setOnClickShareable, calling onClickShareable invokes the callback", async () => {
        let called = false;
        useCertificateDetailNavStore.getState().setOnClickShareable(async () => {
            called = true;
        });
        await useCertificateDetailNavStore.getState().onClickShareable();
        expect(called).toBe(true);
    });

    it("isShareableLoading returns to false after callback completes", async () => {
        let called = false;
        useCertificateDetailNavStore.getState().setOnClickShareable(async () => {
            called = true;
        });
        await useCertificateDetailNavStore.getState().onClickShareable();
        expect(called).toBe(true);
        expect(useCertificateDetailNavStore.getState().isShareableLoading).toBe(false);
    });

    it("isShareableLoading stays true while async callback is running", async () => {
        let resolveCallback!: () => void;
        const asyncCallback = vi.fn(
            () =>
                new Promise<void>((resolve) => {
                    resolveCallback = resolve;
                }),
        );

        useCertificateDetailNavStore.getState().setOnClickShareable(asyncCallback);
        const clickPromise = useCertificateDetailNavStore.getState().onClickShareable();

        expect(useCertificateDetailNavStore.getState().isShareableLoading).toBe(true);

        resolveCallback();
        await clickPromise;

        expect(useCertificateDetailNavStore.getState().isShareableLoading).toBe(false);
    });

    it("isShareableLoading resets to false even when async callback rejects", async () => {
        const asyncCallback = vi.fn(() => Promise.reject(new Error("fail")));

        useCertificateDetailNavStore.getState().setOnClickShareable(asyncCallback);

        await expect(useCertificateDetailNavStore.getState().onClickShareable()).rejects.toThrow(
            "fail",
        );

        expect(useCertificateDetailNavStore.getState().isShareableLoading).toBe(false);
    });
});

describe("useCertificateDetailsSharedNavStore", () => {
    beforeEach(() => {
        useCertificateDetailsSharedNavStore.setState({
            shareableUrl: null,
            shareableHandle: null,
            isPasswordProtected: false,
            isPublished: false,
            isPasswordLoading: false,
            isPublishLoading: false,
        });
    });

    it("initializes with default values", () => {
        const state = useCertificateDetailsSharedNavStore.getState();
        expect(state.shareableUrl).toBeNull();
        expect(state.shareableHandle).toBeNull();
        expect(state.isPasswordProtected).toBe(false);
        expect(state.isPublished).toBe(false);
        expect(state.isPasswordLoading).toBe(false);
        expect(state.isPublishLoading).toBe(false);
    });

    it("setOnClickPassword — after calling, onClickPassword() invokes the callback", async () => {
        let called = false;
        useCertificateDetailsSharedNavStore.getState().setOnClickPassword(async () => {
            called = true;
        });
        await useCertificateDetailsSharedNavStore.getState().onClickPassword(false);
        expect(called).toBe(true);
    });

    it("isPasswordLoading stays true while async password callback is running", async () => {
        let resolveCallback!: () => void;
        const asyncCallback = vi.fn(
            () =>
                new Promise<void>((resolve) => {
                    resolveCallback = resolve;
                }),
        );

        useCertificateDetailsSharedNavStore.getState().setOnClickPassword(asyncCallback);
        const clickPromise = useCertificateDetailsSharedNavStore.getState().onClickPassword(true);

        expect(useCertificateDetailsSharedNavStore.getState().isPasswordLoading).toBe(true);

        resolveCallback();
        await clickPromise;

        expect(useCertificateDetailsSharedNavStore.getState().isPasswordLoading).toBe(false);
    });

    it("isPasswordLoading resets to false even when async password callback rejects", async () => {
        const asyncCallback = vi.fn(() => Promise.reject(new Error("fail")));

        useCertificateDetailsSharedNavStore.getState().setOnClickPassword(asyncCallback);

        await expect(
            useCertificateDetailsSharedNavStore.getState().onClickPassword(false),
        ).rejects.toThrow("fail");

        expect(useCertificateDetailsSharedNavStore.getState().isPasswordLoading).toBe(false);
    });

    it("setOnChangePublish — after calling, onChangePublish(true) invokes callback with true", () => {
        let receivedValue: boolean | undefined;
        useCertificateDetailsSharedNavStore
            .getState()
            .setOnChangePublish(async (isPublished: boolean) => {
                receivedValue = isPublished;
            });
        useCertificateDetailsSharedNavStore.getState().onChangePublish(true);
        expect(receivedValue).toBe(true);
    });
});
