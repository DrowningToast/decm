import { describe, it, expect, vi, beforeEach } from "vitest";
import { useEventInvitationNavStore } from "./event-invitation";

describe("useEventInvitationNavStore", () => {
    beforeEach(() => {
        useEventInvitationNavStore.setState({ onAcceptCallback: undefined });
    });

    it("initializes with undefined callback", () => {
        expect(useEventInvitationNavStore.getState().onAcceptCallback).toBeUndefined();
    });

    it("sets onAcceptCallback", () => {
        const callback = vi.fn();
        useEventInvitationNavStore.getState().setOnAcceptCallback(callback);
        expect(useEventInvitationNavStore.getState().onAcceptCallback).toBe(callback);
    });

    it("stored callback is callable", () => {
        const callback = vi.fn();
        useEventInvitationNavStore.getState().setOnAcceptCallback(callback);
        useEventInvitationNavStore.getState().onAcceptCallback!();
        expect(callback).toHaveBeenCalledOnce();
    });

    it("replaces previous callback", () => {
        const first = vi.fn();
        const second = vi.fn();
        useEventInvitationNavStore.getState().setOnAcceptCallback(first);
        useEventInvitationNavStore.getState().setOnAcceptCallback(second);
        useEventInvitationNavStore.getState().onAcceptCallback!();
        expect(first).not.toHaveBeenCalled();
        expect(second).toHaveBeenCalledOnce();
    });
});
