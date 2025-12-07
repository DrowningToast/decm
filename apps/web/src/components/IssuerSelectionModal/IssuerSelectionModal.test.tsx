import { render, screen, fireEvent, within } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { IssuerSelectionModal } from "./IssuerSelectionModal";
import type { Issuer } from "./IssuerSelectionModal";

// Mock i18n
vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

describe("IssuerSelectionModal", () => {
    const mockIssuers: Issuer[] = [
        {
            id: "issuer-1",
            name: "John Doe",
            email: "john@example.com",
            organization: "MIT",
        },
        {
            id: "issuer-2",
            name: "Jane Smith",
            email: "jane@example.com",
            organization: "Stanford",
        },
        {
            id: "issuer-3",
            name: "Bob Johnson",
            email: "bob@example.com",
        },
    ];

    const defaultProps = {
        isOpen: true,
        onClose: vi.fn(),
        onConfirm: vi.fn(),
        searchResults: mockIssuers,
        initialSelectedIds: new Set<string>(),
        isLoading: false,
        searchQuery: "test",
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should render modal when open", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        expect(screen.getByText("certificateSettings.step1.modal.title")).toBeInTheDocument();
    });

    it("should not render modal when closed", () => {
        render(<IssuerSelectionModal {...defaultProps} isOpen={false} />);

        expect(screen.queryByText("certificateSettings.step1.modal.title")).not.toBeInTheDocument();
    });

    it("should display loading state", () => {
        render(<IssuerSelectionModal {...defaultProps} isLoading={true} />);

        expect(screen.getByText("common.loading")).toBeInTheDocument();
    });

    it("should display search results", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        expect(screen.getByText("John Doe")).toBeInTheDocument();
        expect(screen.getByText("jane@example.com")).toBeInTheDocument();
        expect(screen.getByText("Stanford")).toBeInTheDocument();
    });

    it("should display (empty) for missing organization", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        const rows = screen.getAllByRole("row");
        const bobRow = rows.find((row) => within(row).queryByText("Bob Johnson"));
        expect(bobRow).toBeDefined();
        if (bobRow) {
            expect(within(bobRow).getByText("(empty)")).toBeInTheDocument();
        }
    });

    it("should display no results message when search results are empty", () => {
        render(
            <IssuerSelectionModal {...defaultProps} searchResults={[]} searchQuery="nonexistent" />,
        );

        expect(screen.getByText("certificateSettings.step1.modal.noResults")).toBeInTheDocument();
    });

    it("should display enter search message when no search query", () => {
        render(<IssuerSelectionModal {...defaultProps} searchResults={[]} searchQuery="" />);

        expect(screen.getByText("certificateSettings.step1.modal.enterSearch")).toBeInTheDocument();
    });

    it("should toggle issuer selection on checkbox click", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        const checkboxes = screen.getAllByRole("checkbox");
        const firstCheckbox = checkboxes[0];

        expect(firstCheckbox).not.toBeChecked();

        fireEvent.click(firstCheckbox);

        expect(firstCheckbox).toBeChecked();

        fireEvent.click(firstCheckbox);

        expect(firstCheckbox).not.toBeChecked();
    });

    it("should toggle issuer selection on row click", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        const rows = screen.getAllByRole("row");
        const dataRow = rows[1]; // First data row (skip header)
        const checkbox = within(dataRow).getByRole("checkbox");

        expect(checkbox).not.toBeChecked();

        fireEvent.click(dataRow);

        expect(checkbox).toBeChecked();
    });

    it("should display selected count", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        const checkboxes = screen.getAllByRole("checkbox");
        fireEvent.click(checkboxes[0]);
        fireEvent.click(checkboxes[1]);

        // Should show "2 selected" in description and confirm button
        const selectedTexts = screen.getAllByText(/2/);
        expect(selectedTexts.length).toBeGreaterThan(0);
    });

    it("should initialize with pre-selected issuers", () => {
        const initialSelectedIds = new Set(["issuer-1", "issuer-2"]);

        render(<IssuerSelectionModal {...defaultProps} initialSelectedIds={initialSelectedIds} />);

        const checkboxes = screen.getAllByRole("checkbox");
        expect(checkboxes[0]).toBeChecked();
        expect(checkboxes[1]).toBeChecked();
        expect(checkboxes[2]).not.toBeChecked();
    });

    it("should call onConfirm with selected issuers", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        const checkboxes = screen.getAllByRole("checkbox");
        fireEvent.click(checkboxes[0]);
        fireEvent.click(checkboxes[1]);

        const confirmButton = screen.getByText(/certificateSettings.step1.modal.chooseButton/);
        fireEvent.click(confirmButton);

        expect(defaultProps.onConfirm).toHaveBeenCalledWith([mockIssuers[0], mockIssuers[1]]);
        expect(defaultProps.onClose).toHaveBeenCalled();
    });

    it("should call onCancel and reset selection", () => {
        const initialSelectedIds = new Set(["issuer-1"]);

        render(<IssuerSelectionModal {...defaultProps} initialSelectedIds={initialSelectedIds} />);

        const checkboxes = screen.getAllByRole("checkbox");
        // Uncheck the first one
        fireEvent.click(checkboxes[0]);
        // Check the second one
        fireEvent.click(checkboxes[1]);

        const cancelButton = screen.getByText("common.cancel");
        fireEvent.click(cancelButton);

        expect(defaultProps.onClose).toHaveBeenCalled();
        expect(defaultProps.onConfirm).not.toHaveBeenCalled();
    });

    it("should disable confirm button when no issuers selected", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        const confirmButton = screen
            .getByText(/certificateSettings.step1.modal.chooseButton/)
            .closest("button");

        expect(confirmButton).toBeDisabled();
    });

    it("should enable confirm button when issuers are selected", () => {
        render(<IssuerSelectionModal {...defaultProps} />);

        const checkboxes = screen.getAllByRole("checkbox");
        fireEvent.click(checkboxes[0]);

        const confirmButton = screen
            .getByText(/certificateSettings.step1.modal.chooseButton/)
            .closest("button");

        expect(confirmButton).not.toBeDisabled();
    });

    it("should disable buttons during loading", () => {
        render(<IssuerSelectionModal {...defaultProps} isLoading={true} />);

        const cancelButton = screen.getByText("common.cancel").closest("button");
        const confirmButton = screen
            .getByText(/certificateSettings.step1.modal.chooseButton/)
            .closest("button");

        expect(cancelButton).toBeDisabled();
        expect(confirmButton).toBeDisabled();
    });

    it("should reset selection when modal opens", () => {
        const { rerender } = render(<IssuerSelectionModal {...defaultProps} isOpen={false} />);

        // Open modal with initial selection
        rerender(
            <IssuerSelectionModal
                {...defaultProps}
                isOpen={true}
                initialSelectedIds={new Set(["issuer-1"])}
            />,
        );

        const checkboxes = screen.getAllByRole("checkbox");
        expect(checkboxes[0]).toBeChecked();
    });
});
