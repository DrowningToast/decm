import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { SelectedIssuersTable } from "./SelectedIssuersTable";
import type { EventIssuer } from "@/services/EventService/EventService";

// Mock i18n
vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

// Mock ConfirmModal
vi.mock("../ConfirmModal", () => ({
    default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

describe("SelectedIssuersTable", () => {
    const mockIssuers: EventIssuer[] = [
        {
            id: "issuer-1",
            eventId: "event-123",
            isSigned: true,
            issuerCredentialId: "cred-1",
            issuerProfile: {
                id: "profile-1",
                authenticationCredentialId: "cred-1",
                firstName: "John",
                lastName: "Doe",
                email: "john@example.com",
                academicInstitution: "MIT",
            },
            createdAt: new Date("2024-01-01"),
            updatedAt: new Date("2024-01-01"),
        },
        {
            id: "issuer-2",
            eventId: "event-123",
            isSigned: false,
            issuerCredentialId: "cred-2",
            issuerProfile: {
                id: "profile-2",
                authenticationCredentialId: "cred-2",
                firstName: "Jane",
                lastName: "Smith",
                email: "jane@example.com",
                academicInstitution: "Stanford",
            },
            createdAt: new Date("2024-01-02"),
            updatedAt: new Date("2024-01-02"),
        },
        {
            id: "issuer-3",
            eventId: "event-123",
            isSigned: false,
            issuerCredentialId: "cred-3",
            issuerProfile: {
                id: "profile-3",
                authenticationCredentialId: "cred-3",
                firstName: "Bob",
                lastName: "",
                email: "bob@example.com",
            },
            createdAt: new Date("2024-01-03"),
            updatedAt: new Date("2024-01-03"),
        },
    ];

    const defaultProps = {
        selectedIssuers: mockIssuers,
        onRemoveIssuer: vi.fn(),
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should render table with issuers", () => {
        render(<SelectedIssuersTable {...defaultProps} />);

        expect(screen.getByText("John Doe")).toBeInTheDocument();
        expect(screen.getByText("jane@example.com")).toBeInTheDocument();
        expect(screen.getByText("Stanford")).toBeInTheDocument();
    });

    it("should not render when selectedIssuers is undefined", () => {
        const { container } = render(
            <SelectedIssuersTable selectedIssuers={undefined} onRemoveIssuer={vi.fn()} />,
        );

        expect(container.firstChild).toBeNull();
    });

    it("should not render when selectedIssuers is empty", () => {
        const { container } = render(
            <SelectedIssuersTable selectedIssuers={[]} onRemoveIssuer={vi.fn()} />,
        );

        expect(container.firstChild).toBeNull();
    });

    it("should display issuer count", () => {
        render(<SelectedIssuersTable {...defaultProps} />);

        expect(screen.getByText(/certificateSettings.step1.selectedIssuers/)).toBeInTheDocument();
    });

    it("should display full name from firstName and lastName", () => {
        render(<SelectedIssuersTable {...defaultProps} />);

        expect(screen.getByText("John Doe")).toBeInTheDocument();
        expect(screen.getByText("Jane Smith")).toBeInTheDocument();
    });

    it("should handle missing lastName", () => {
        render(<SelectedIssuersTable {...defaultProps} />);

        expect(screen.getByText("Bob")).toBeInTheDocument();
    });

    it("should display (empty) when name is missing", () => {
        const issuersWithoutName: EventIssuer[] = [
            {
                id: "issuer-4",
                eventId: "event-123",
                isSigned: false,
                issuerCredentialId: "cred-4",
                issuerProfile: {
                    id: "profile-4",
                    authenticationCredentialId: "cred-4",
                    firstName: "",
                    lastName: "",
                    email: "test@example.com",
                },
                createdAt: new Date("2024-01-04"),
                updatedAt: new Date("2024-01-04"),
            },
        ];

        render(
            <SelectedIssuersTable selectedIssuers={issuersWithoutName} onRemoveIssuer={vi.fn()} />,
        );

        const emptyElements = screen.getAllByText("(empty)");
        expect(emptyElements.length).toBeGreaterThan(0);
    });

    it("should display (empty) for missing organization", () => {
        render(<SelectedIssuersTable {...defaultProps} />);

        // Bob's issuer doesn't have academicInstitution
        const emptyElements = screen.getAllByText("(empty)");
        expect(emptyElements.length).toBeGreaterThan(0);
    });

    it("should display email addresses", () => {
        render(<SelectedIssuersTable {...defaultProps} />);

        expect(screen.getByText("john@example.com")).toBeInTheDocument();
        expect(screen.getByText("jane@example.com")).toBeInTheDocument();
        expect(screen.getByText("bob@example.com")).toBeInTheDocument();
    });

    it("should render remove buttons for each issuer", () => {
        render(<SelectedIssuersTable {...defaultProps} />);

        const removeButtons = screen.getAllByRole("button");
        expect(removeButtons.length).toBe(mockIssuers.length);
    });

    it("should display table headers", () => {
        render(<SelectedIssuersTable {...defaultProps} />);

        expect(screen.getByText("certificateSettings.step1.table.name")).toBeInTheDocument();
        expect(screen.getByText("certificateSettings.step1.table.email")).toBeInTheDocument();
        expect(
            screen.getByText("certificateSettings.step1.table.organization"),
        ).toBeInTheDocument();
        expect(screen.getByText("certificateSettings.step1.table.actions")).toBeInTheDocument();
    });

    it("should handle issuers without profile", () => {
        const issuersWithoutProfile: EventIssuer[] = [
            {
                id: "issuer-5",
                eventId: "event-123",
                isSigned: false,
                issuerCredentialId: "cred-5",
                createdAt: new Date("2024-01-05"),
                updatedAt: new Date("2024-01-05"),
            },
        ];

        render(
            <SelectedIssuersTable
                selectedIssuers={issuersWithoutProfile}
                onRemoveIssuer={vi.fn()}
            />,
        );

        // Should still render the table
        expect(screen.getByText("certificateSettings.step1.table.name")).toBeInTheDocument();

        // Should show (empty) for all fields
        const emptyElements = screen.getAllByText("(empty)");
        expect(emptyElements.length).toBeGreaterThan(0);
    });

    it("should handle single issuer", () => {
        render(
            <SelectedIssuersTable selectedIssuers={[mockIssuers[0]]} onRemoveIssuer={vi.fn()} />,
        );

        expect(screen.getByText("John Doe")).toBeInTheDocument();
        expect(screen.getByText("john@example.com")).toBeInTheDocument();
    });

    it("should handle multiple issuers from same organization", () => {
        const issuersFromMIT: EventIssuer[] = [
            mockIssuers[0],
            {
                ...mockIssuers[1],
                issuerProfile: {
                    ...mockIssuers[1].issuerProfile!,
                    academicInstitution: "MIT",
                },
            },
        ];

        render(<SelectedIssuersTable selectedIssuers={issuersFromMIT} onRemoveIssuer={vi.fn()} />);

        const mitElements = screen.getAllByText("MIT");
        expect(mitElements.length).toBe(2);
    });
});
