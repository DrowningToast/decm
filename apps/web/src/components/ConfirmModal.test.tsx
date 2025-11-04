import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import ConfirmModal from "./ConfirmModal";

// Mock AlertDialog components
vi.mock("./ui/alert-dialog", () => ({
    AlertDialog: ({
        children,
        open,
    }: React.PropsWithChildren<{ open?: boolean; onOpenChange?: (open: boolean) => void }>) => (
        <div data-testid="alert-dialog" data-open={open}>
            {children}
        </div>
    ),
    AlertDialogTrigger: ({ children }: React.PropsWithChildren) => (
        <button data-testid="trigger">{children}</button>
    ),
    AlertDialogContent: ({ children }: React.PropsWithChildren) => (
        <div data-testid="content">{children}</div>
    ),
    AlertDialogHeader: ({ children }: React.PropsWithChildren) => (
        <div data-testid="header">{children}</div>
    ),
    AlertDialogTitle: ({ children }: React.PropsWithChildren) => (
        <h2 data-testid="title">{children}</h2>
    ),
    AlertDialogDescription: ({ children }: React.PropsWithChildren) => (
        <p data-testid="description">{children}</p>
    ),
    AlertDialogFooter: ({ children }: React.PropsWithChildren) => (
        <div data-testid="footer">{children}</div>
    ),
    AlertDialogCancel: ({
        children,
        onClick,
    }: React.PropsWithChildren<{ onClick?: () => void }>) => (
        <button data-testid="cancel-button" onClick={onClick}>
            {children}
        </button>
    ),
    AlertDialogAction: ({
        children,
        onClick,
    }: React.PropsWithChildren<{ onClick?: () => void }>) => (
        <button data-testid="action-button" onClick={onClick}>
            {children}
        </button>
    ),
}));

vi.mock("lucide-react", () => ({
    TriangleAlert: () => <div data-testid="default-icon">⚠️</div>,
}));

describe("ConfirmModal Component", () => {
    it("should render with required props", () => {
        render(
            <ConfirmModal
                title="Confirm Action"
                message="Are you sure?"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
            />,
        );

        expect(screen.getByTestId("title")).toHaveTextContent("Confirm Action");
        expect(screen.getByTestId("description")).toHaveTextContent("Are you sure?");
    });

    it("should render cancel and confirm buttons", () => {
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={vi.fn()}
                cancelText="No"
                confirmText="Yes"
            />,
        );

        expect(screen.getByTestId("cancel-button")).toHaveTextContent("No");
        expect(screen.getByTestId("action-button")).toHaveTextContent("Yes");
    });

    it("should call onConfirm when confirm button is clicked", async () => {
        const onConfirm = vi.fn();
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={onConfirm}
                cancelText="Cancel"
                confirmText="Confirm"
            />,
        );

        const confirmButton = screen.getByTestId("action-button");
        fireEvent.click(confirmButton);

        await waitFor(() => {
            expect(onConfirm).toHaveBeenCalled();
        });
    });

    it("should call onCancel when cancel button is clicked", async () => {
        const onCancel = vi.fn();
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={vi.fn()}
                onCancel={onCancel}
                cancelText="Cancel"
                confirmText="Confirm"
            />,
        );

        const cancelButton = screen.getByTestId("cancel-button");
        fireEvent.click(cancelButton);

        await waitFor(() => {
            expect(onCancel).toHaveBeenCalled();
        });
    });

    it("should render default icon when not provided", () => {
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
            />,
        );

        expect(screen.getByTestId("default-icon")).toBeInTheDocument();
    });

    it("should render custom icon when provided", () => {
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
                icon={<AlertCircle data-testid="custom-icon" />}
            />,
        );

        expect(screen.getByTestId("custom-icon")).toBeInTheDocument();
    });

    it("should render trigger children", () => {
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
            >
                <span>Click me</span>
            </ConfirmModal>,
        );

        expect(screen.getByText("Click me")).toBeInTheDocument();
    });

    it("should handle controlled open state", () => {
        const onOpenChange = vi.fn();
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
                open={true}
                onOpenChange={onOpenChange}
            />,
        );

        const dialog = screen.getByTestId("alert-dialog");
        expect(dialog).toHaveAttribute("data-open", "true");
    });

    it("should apply destructive styling when destructive prop is true", () => {
        render(
            <ConfirmModal
                title="Delete Item"
                message="This action cannot be undone"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Delete"
                destructive={true}
            />,
        );

        const actionButton = screen.getByTestId("action-button");
        // Check for destructive styling classes
        expect(actionButton.className).toContain("bg-red-500");
    });

    it("should not apply destructive styling when destructive prop is false", () => {
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
                destructive={false}
            />,
        );

        const actionButton = screen.getByTestId("action-button");
        expect(actionButton.className).not.toContain("bg-red-500");
    });

    it("should handle long title and message", () => {
        const longTitle = "A".repeat(100);
        const longMessage = "B".repeat(200);

        render(
            <ConfirmModal
                title={longTitle}
                message={longMessage}
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
            />,
        );

        expect(screen.getByTestId("title")).toHaveTextContent(longTitle);
        expect(screen.getByTestId("description")).toHaveTextContent(longMessage);
    });

    it("should handle special characters in text", () => {
        render(
            <ConfirmModal
                title='Delete "Important" item?'
                message="This will remove <item> permanently"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
            />,
        );

        expect(screen.getByTestId("title")).toBeInTheDocument();
        expect(screen.getByTestId("description")).toBeInTheDocument();
    });

    it("should not call onCancel if not provided", async () => {
        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={vi.fn()}
                cancelText="Cancel"
                confirmText="Confirm"
            />,
        );

        const cancelButton = screen.getByTestId("cancel-button");
        fireEvent.click(cancelButton);

        // Should not throw
        expect(true).toBe(true);
    });

    it("should handle multiple confirm/cancel cycles", async () => {
        const onConfirm = vi.fn();
        const onCancel = vi.fn();

        render(
            <ConfirmModal
                title="Test"
                message="Test message"
                onConfirm={onConfirm}
                onCancel={onCancel}
                cancelText="Cancel"
                confirmText="Confirm"
            />,
        );

        const confirmButton = screen.getByTestId("action-button");
        const cancelButton = screen.getByTestId("cancel-button");

        fireEvent.click(confirmButton);
        fireEvent.click(cancelButton);
        fireEvent.click(confirmButton);

        await waitFor(() => {
            expect(onConfirm).toHaveBeenCalledTimes(2);
            expect(onCancel).toHaveBeenCalledTimes(1);
        });
    });
});
