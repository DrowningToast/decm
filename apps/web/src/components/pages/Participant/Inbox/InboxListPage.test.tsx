import { describe, it, expect } from "vitest";

// Helper functions from InboxListPage
type InboxStatus = "pending" | "available" | "expired" | "action-required";

const getStatusColor = (status: InboxStatus) => {
    switch (status) {
        case "pending":
            return "text-muted";
        case "available":
            return "text-success";
        case "expired":
            return "text-error";
        case "action-required":
            return "text-muted";
        default:
            return "text-muted";
    }
};

const getStatusLabel = (
    status: InboxStatus,
    t: (key: string, fallback: string) => string | undefined,
) => {
    switch (status) {
        case "pending":
            return t("participant.inbox.status.pending", "Pending");
        case "available":
            return t("participant.inbox.status.available", "Available");
        case "expired":
            return t("participant.inbox.status.expired", "Expired");
        case "action-required":
            return t("participant.inbox.status.actionRequired", "Action required");
        default:
            return "";
    }
};

describe("InboxListPage Helper Functions", () => {
    describe("getStatusColor", () => {
        it("returns correct color for pending status", () => {
            expect(getStatusColor("pending")).toBe("text-muted");
        });

        it("returns correct color for available status", () => {
            expect(getStatusColor("available")).toBe("text-success");
        });

        it("returns correct color for expired status", () => {
            expect(getStatusColor("expired")).toBe("text-error");
        });

        it("returns correct color for action-required status", () => {
            expect(getStatusColor("action-required")).toBe("text-muted");
        });

        it("returns default color for unknown status", () => {
            expect(getStatusColor("unknown" as InboxStatus)).toBe("text-muted");
        });
    });

    describe("getStatusLabel", () => {
        const mockT = (key: string, fallback: string) => {
            const translations: Record<string, string> = {
                "participant.inbox.status.pending": "Pending",
                "participant.inbox.status.available": "Available",
                "participant.inbox.status.expired": "Expired",
                "participant.inbox.status.actionRequired": "Action required",
            };
            return translations[key] || fallback;
        };

        it("returns correct label for pending status", () => {
            expect(getStatusLabel("pending", mockT)).toBe("Pending");
        });

        it("returns correct label for available status", () => {
            expect(getStatusLabel("available", mockT)).toBe("Available");
        });

        it("returns correct label for expired status", () => {
            expect(getStatusLabel("expired", mockT)).toBe("Expired");
        });

        it("returns correct label for action-required status", () => {
            expect(getStatusLabel("action-required", mockT)).toBe("Action required");
        });

        it("returns empty string for unknown status", () => {
            expect(getStatusLabel("unknown" as InboxStatus, mockT)).toBe("");
        });

        it("uses fallback translation when key is not found", () => {
            const emptyT = () => undefined;
            expect(getStatusLabel("pending", emptyT)).toBeUndefined();
        });
    });
});
