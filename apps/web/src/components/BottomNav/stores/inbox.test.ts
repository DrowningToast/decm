import { describe, it, expect, vi, beforeEach } from "vitest";
import { useInboxNavStore } from "./inbox";

describe("useInboxNavStore", () => {
    beforeEach(() => {
        useInboxNavStore.setState({
            onViewEventCallback: undefined,
            onViewCertificateCallback: undefined,
        });
    });

    it("initializes with undefined callbacks", () => {
        const state = useInboxNavStore.getState();
        expect(state.onViewEventCallback).toBeUndefined();
        expect(state.onViewCertificateCallback).toBeUndefined();
    });

    it("sets onViewEventCallback", () => {
        const callback = vi.fn();
        useInboxNavStore.getState().setOnViewEventCallback(callback);
        expect(useInboxNavStore.getState().onViewEventCallback).toBe(callback);
    });

    it("sets onViewCertificateCallback", () => {
        const callback = vi.fn();
        useInboxNavStore.getState().setOnViewCertificateCallback(callback);
        expect(useInboxNavStore.getState().onViewCertificateCallback).toBe(callback);
    });

    it("resets all callbacks", () => {
        useInboxNavStore.getState().setOnViewEventCallback(vi.fn());
        useInboxNavStore.getState().setOnViewCertificateCallback(vi.fn());

        useInboxNavStore.getState().resetCallbacks();

        const state = useInboxNavStore.getState();
        expect(state.onViewEventCallback).toBeUndefined();
        expect(state.onViewCertificateCallback).toBeUndefined();
    });

    it("stored callbacks are callable", () => {
        const eventCb = vi.fn();
        const certCb = vi.fn();
        useInboxNavStore.getState().setOnViewEventCallback(eventCb);
        useInboxNavStore.getState().setOnViewCertificateCallback(certCb);

        useInboxNavStore.getState().onViewEventCallback!();
        useInboxNavStore.getState().onViewCertificateCallback!();

        expect(eventCb).toHaveBeenCalledOnce();
        expect(certCb).toHaveBeenCalledOnce();
    });
});
