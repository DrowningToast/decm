import { describe, it, expect } from "vitest";
import { renderHook } from "@testing-library/react";
import { useHelmet } from "./useHelmet";

describe("useHelmet", () => {
    it("should return correct title and description for home page", () => {
        const { result } = renderHook(() => useHelmet({ pageType: "home" }));

        expect(result.current.title).toContain("Home");
        expect(result.current.title).toContain("DECM");
        expect(result.current.description).toContain("Web 3.0");
        expect(result.current.faviconUrl).toBe("/favicon.ico");
        expect(result.current.themeColor).toBe("#ffffff");
    });

    it("should return correct theme color for login page", () => {
        const { result } = renderHook(() => useHelmet({ pageType: "login" }));

        expect(result.current.title).toContain("Login");
        expect(result.current.themeColor).toBe("#3B82F6"); // Blue
    });

    it("should return correct theme color for profile page", () => {
        const { result } = renderHook(() => useHelmet({ pageType: "profile" }));

        expect(result.current.title).toContain("Profile");
        expect(result.current.themeColor).toBe("#10B981"); // Green
    });

    it("should return correct theme color for events page", () => {
        const { result } = renderHook(() => useHelmet({ pageType: "events" }));

        expect(result.current.title).toContain("Events");
        expect(result.current.themeColor).toBe("#F59E0B"); // Orange
    });

    it("should return correct theme color for credentials page", () => {
        const { result } = renderHook(() => useHelmet({ pageType: "credentials" }));

        expect(result.current.title).toContain("Credentials");
        expect(result.current.themeColor).toBe("#8B5CF6"); // Purple
    });

    it("should use custom title when provided", () => {
        const customTitle = "Custom Page Title";
        const { result } = renderHook(() => useHelmet({ title: customTitle, pageType: "home" }));

        expect(result.current.title).toBe(customTitle);
    });

    it("should use custom description when provided", () => {
        const customDesc = "Custom description for testing";
        const { result } = renderHook(() =>
            useHelmet({ description: customDesc, pageType: "home" }),
        );

        expect(result.current.description).toBe(customDesc);
    });

    it("should override both title and description", () => {
        const customTitle = "My Custom Title";
        const customDesc = "My Custom Description";
        const { result } = renderHook(() =>
            useHelmet({
                title: customTitle,
                description: customDesc,
                pageType: "profile",
            }),
        );

        expect(result.current.title).toBe(customTitle);
        expect(result.current.description).toBe(customDesc);
        expect(result.current.themeColor).toBe("#10B981"); // Still profile theme
    });

    it("should use default home page when no pageType provided", () => {
        const { result } = renderHook(() => useHelmet());

        expect(result.current.title).toContain("Home");
        expect(result.current.faviconUrl).toBe("/favicon.ico");
    });

    it("should memoize results when dependencies don't change", () => {
        const { result: result1 } = renderHook(() => useHelmet({ pageType: "events" }));
        const { result: result2 } = renderHook(() => useHelmet({ pageType: "events" }));

        // Both should be equal (memoization check)
        expect(result1.current.title).toBe(result2.current.title);
    });

    it("should include base title in all pages", () => {
        const pageTypes: Array<"home" | "login" | "profile" | "events" | "credentials"> = [
            "home",
            "login",
            "profile",
            "events",
            "credentials",
        ];

        pageTypes.forEach((pageType) => {
            const { result } = renderHook(() => useHelmet({ pageType }));
            expect(result.current.title).toContain("DECM - Decentralized Event Management");
        });
    });

    it("should return valid favicon URL for all pages", () => {
        const pageTypes: Array<"home" | "login" | "profile" | "events" | "credentials"> = [
            "home",
            "login",
            "profile",
            "events",
            "credentials",
        ];

        pageTypes.forEach((pageType) => {
            const { result } = renderHook(() => useHelmet({ pageType }));
            expect(result.current.faviconUrl).toBe("/favicon.ico");
        });
    });
});
