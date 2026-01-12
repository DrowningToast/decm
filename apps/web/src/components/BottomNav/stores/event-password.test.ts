import { describe, it, expect, vi, beforeEach } from "vitest";
import { useEventPasswordNavStore } from "./event-password";

describe("useEventPasswordNavStore", () => {
    beforeEach(() => {
        useEventPasswordNavStore.setState({
            password: "",
            onSubmitCallback: undefined,
        });
    });

    it("initializes with empty password and undefined callback", () => {
        const state = useEventPasswordNavStore.getState();
        expect(state.password).toBe("");
        expect(state.onSubmitCallback).toBeUndefined();
    });

    it("sets password", () => {
        useEventPasswordNavStore.getState().setPassword("secret123");
        expect(useEventPasswordNavStore.getState().password).toBe("secret123");
    });

    it("resets password to empty string", () => {
        useEventPasswordNavStore.getState().setPassword("secret123");
        useEventPasswordNavStore.getState().resetPassword();
        expect(useEventPasswordNavStore.getState().password).toBe("");
    });

    it("sets onSubmitCallback", () => {
        const callback = vi.fn();
        useEventPasswordNavStore.getState().setOnSubmitCallback(callback);
        expect(useEventPasswordNavStore.getState().onSubmitCallback).toBe(callback);
    });

    it("stored callback receives password argument", () => {
        const callback = vi.fn();
        useEventPasswordNavStore.getState().setOnSubmitCallback(callback);
        useEventPasswordNavStore.getState().onSubmitCallback!("mypass");
        expect(callback).toHaveBeenCalledWith("mypass");
    });
});
