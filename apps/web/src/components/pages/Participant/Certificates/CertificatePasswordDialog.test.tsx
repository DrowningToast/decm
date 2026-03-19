import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { CertificatePasswordDialog } from "./CertificatePasswordDialog";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback ?? _key,
    }),
}));

vi.mock("./useCertificateUpdateSharePasswordUsecase", () => ({
    useCertificateUpdateSharePasswordUsecase: () => ({
        updateSharePassword: vi.fn(),
        isPending: false,
    }),
}));

vi.mock("@/components/ui/alert-dialog", () => ({
    AlertDialog: ({
        children,
        open,
    }: React.PropsWithChildren<{ open?: boolean; onOpenChange?: (open: boolean) => void }>) => (
        <div data-testid="alert-dialog" data-open={open}>
            {children}
        </div>
    ),
    AlertDialogContent: ({ children }: React.PropsWithChildren) => (
        <div data-testid="alert-dialog-content">{children}</div>
    ),
    AlertDialogHeader: ({ children }: React.PropsWithChildren) => (
        <div data-testid="alert-dialog-header">{children}</div>
    ),
    AlertDialogTitle: ({ children }: React.PropsWithChildren) => (
        <h2 data-testid="alert-dialog-title">{children}</h2>
    ),
    AlertDialogDescription: ({ children }: React.PropsWithChildren) => (
        <p data-testid="alert-dialog-description">{children}</p>
    ),
    AlertDialogFooter: ({ children }: React.PropsWithChildren) => (
        <div data-testid="alert-dialog-footer">{children}</div>
    ),
    AlertDialogCancel: ({
        children,
        onClick,
    }: React.PropsWithChildren<{ onClick?: () => void }>) => (
        <button data-testid="cancel-button" onClick={onClick}>
            {children}
        </button>
    ),
}));

vi.mock("@/components/ui/button", () => ({
    Button: ({
        children,
        onClick,
    }: React.PropsWithChildren<{ onClick?: () => void; [key: string]: unknown }>) => (
        <button data-testid="save-button" onClick={onClick}>
            {children}
        </button>
    ),
    buttonVariants: () => "",
}));

const defaultProps = {
    open: true,
    onOpenChange: vi.fn(),
    shareId: "share-123",
    hasPassword: false,
};

describe("CertificatePasswordDialog", () => {
    it("renders with open=true", () => {
        render(<CertificatePasswordDialog {...defaultProps} open={true} />);
        expect(screen.getByTestId("alert-dialog")).toHaveAttribute("data-open", "true");
    });

    it("renders with open=false", () => {
        render(<CertificatePasswordDialog {...defaultProps} open={false} />);
        expect(screen.getByTestId("alert-dialog")).toHaveAttribute("data-open", "false");
    });

    it("renders title and description", () => {
        render(<CertificatePasswordDialog {...defaultProps} />);
        expect(screen.getByTestId("alert-dialog-title")).toBeInTheDocument();
        expect(screen.getByTestId("alert-dialog-description")).toBeInTheDocument();
    });

    it("renders Cancel and Save buttons", () => {
        render(<CertificatePasswordDialog {...defaultProps} />);
        expect(screen.getByTestId("cancel-button")).toHaveTextContent("Cancel");
        expect(screen.getByTestId("save-button")).toHaveTextContent("Save");
    });

    it("calls onOpenChange(false) when Cancel is clicked", () => {
        const onOpenChange = vi.fn();
        render(<CertificatePasswordDialog {...defaultProps} onOpenChange={onOpenChange} />);
        fireEvent.click(screen.getByTestId("cancel-button"));
        expect(onOpenChange).toHaveBeenCalledWith(false);
    });

    it("does not call onOpenChange when Save is clicked", () => {
        const onOpenChange = vi.fn();
        render(<CertificatePasswordDialog {...defaultProps} onOpenChange={onOpenChange} />);
        fireEvent.click(screen.getByTestId("save-button"));
        expect(onOpenChange).not.toHaveBeenCalled();
    });
});
