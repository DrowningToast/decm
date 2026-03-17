import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { CertificatePasswordDialog } from "./CertificatePasswordDialog";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback ?? _key,
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
    AlertDialogAction: ({
        children,
        onClick,
    }: React.PropsWithChildren<{ onClick?: () => void }>) => (
        <button data-testid="submit-button" onClick={onClick}>
            {children}
        </button>
    ),
}));

describe("CertificatePasswordDialog", () => {
    it("renders with open=true", () => {
        render(<CertificatePasswordDialog open={true} onOpenChange={vi.fn()} />);
        expect(screen.getByTestId("alert-dialog")).toHaveAttribute("data-open", "true");
    });

    it("renders with open=false", () => {
        render(<CertificatePasswordDialog open={false} onOpenChange={vi.fn()} />);
        expect(screen.getByTestId("alert-dialog")).toHaveAttribute("data-open", "false");
    });

    it("renders title and description", () => {
        render(<CertificatePasswordDialog open={true} onOpenChange={vi.fn()} />);
        expect(screen.getByTestId("alert-dialog-title")).toBeInTheDocument();
        expect(screen.getByTestId("alert-dialog-description")).toBeInTheDocument();
    });

    it("renders Cancel and Submit buttons", () => {
        render(<CertificatePasswordDialog open={true} onOpenChange={vi.fn()} />);
        expect(screen.getByTestId("cancel-button")).toHaveTextContent("Cancel");
        expect(screen.getByTestId("submit-button")).toHaveTextContent("Submit");
    });

    it("calls onOpenChange(false) when Cancel is clicked", () => {
        const onOpenChange = vi.fn();
        render(<CertificatePasswordDialog open={true} onOpenChange={onOpenChange} />);
        fireEvent.click(screen.getByTestId("cancel-button"));
        expect(onOpenChange).toHaveBeenCalledWith(false);
    });

    it("does not call onOpenChange when Submit is clicked", () => {
        const onOpenChange = vi.fn();
        render(<CertificatePasswordDialog open={true} onOpenChange={onOpenChange} />);
        fireEvent.click(screen.getByTestId("submit-button"));
        expect(onOpenChange).not.toHaveBeenCalled();
    });
});
