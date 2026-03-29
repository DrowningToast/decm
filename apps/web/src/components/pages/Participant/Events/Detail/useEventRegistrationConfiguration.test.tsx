import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useEventRegistrationConfiguration } from "./useEventRegistrationConfiguration";

vi.mock("@/services/services", () => ({
    eventRegistrationService: {
        getConfiguration: vi.fn(),
    },
}));

import { eventRegistrationService } from "@/services/services";

const mockConfig = {
    eventType: "invite" as const,
    firstName: "required" as const,
    lastName: "optional" as const,
    email: "required" as const,
    bio: "not_required" as const,
    phoneNumber: "not_required" as const,
    address: "not_required" as const,
    academicInstitution: "not_required" as const,
    academicEmail: "not_required" as const,
};

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false } },
    });
    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("useEventRegistrationConfiguration", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("returns configuration data when query succeeds", async () => {
        vi.mocked(eventRegistrationService.getConfiguration).mockResolvedValue(mockConfig);

        const { result } = renderHook(() => useEventRegistrationConfiguration("event-1"), {
            wrapper: createWrapper(),
        });

        await waitFor(() => expect(result.current.isLoading).toBe(false));

        expect(result.current.data).toEqual(mockConfig);
        expect(result.current.error).toBeNull();
        expect(eventRegistrationService.getConfiguration).toHaveBeenCalledWith("event-1");
    });

    it("is loading while fetching", () => {
        vi.mocked(eventRegistrationService.getConfiguration).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useEventRegistrationConfiguration("event-1"), {
            wrapper: createWrapper(),
        });

        expect(result.current.isLoading).toBe(true);
    });

    it("does not fetch when eventId is empty", () => {
        const { result } = renderHook(() => useEventRegistrationConfiguration(""), {
            wrapper: createWrapper(),
        });

        expect(result.current.isLoading).toBe(false);
        expect(eventRegistrationService.getConfiguration).not.toHaveBeenCalled();
    });

    it("returns error when query fails", async () => {
        vi.mocked(eventRegistrationService.getConfiguration).mockRejectedValue(
            new Error("network error"),
        );

        const { result } = renderHook(() => useEventRegistrationConfiguration("event-1"), {
            wrapper: createWrapper(),
        });

        await waitFor(() => expect(result.current.isLoading).toBe(false));

        expect(result.current.error).toBeTruthy();
        expect(result.current.data).toBeUndefined();
    });
});
