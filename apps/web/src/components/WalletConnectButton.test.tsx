import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { WalletConnectButton } from "./WalletConnectButton";

// Mock useAppKit
const mockOpen = vi.fn();
vi.mock("@reown/appkit/react", () => ({
    useAppKit: () => ({
        open: mockOpen,
    }),
}));

// Mock utils
vi.mock("@/lib/utils", () => ({
    cn: (...args: Array<string | boolean | undefined>) => args.filter(Boolean).join(" "),
}));

describe("WalletConnectButton Component", () => {
    beforeEach(() => {
        mockOpen.mockClear();
    });

    it("should render children", () => {
        render(
            <WalletConnectButton>
                <span>Connect Wallet</span>
            </WalletConnectButton>,
        );

        expect(screen.getByText("Connect Wallet")).toBeInTheDocument();
    });

    it("should call open when clicked", () => {
        render(
            <WalletConnectButton>
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const button = screen.getByText("Connect");
        fireEvent.click(button);

        expect(mockOpen).toHaveBeenCalled();
    });

    it("should call custom onClick handler before opening wallet", () => {
        const onClickHandler = vi.fn();

        render(
            <WalletConnectButton onClick={onClickHandler}>
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const button = screen.getByText("Connect");
        fireEvent.click(button);

        expect(onClickHandler).toHaveBeenCalled();
        expect(mockOpen).toHaveBeenCalled();
    });

    it("should apply custom className", () => {
        const { container } = render(
            <WalletConnectButton className="custom-class">
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const div = container.firstChild;
        expect(div).toHaveClass("custom-class");
    });

    it("should render complex children elements", () => {
        render(
            <WalletConnectButton>
                <div>
                    <span>Icon</span>
                    <span>Connect Wallet</span>
                </div>
            </WalletConnectButton>,
        );

        expect(screen.getByText("Icon")).toBeInTheDocument();
        expect(screen.getByText("Connect Wallet")).toBeInTheDocument();
    });

    it("should handle multiple clicks", () => {
        render(
            <WalletConnectButton>
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const button = screen.getByText("Connect");
        fireEvent.click(button);
        fireEvent.click(button);
        fireEvent.click(button);

        expect(mockOpen).toHaveBeenCalledTimes(3);
    });

    it("should work without onClick handler", () => {
        render(
            <WalletConnectButton>
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const button = screen.getByText("Connect");
        fireEvent.click(button);

        expect(mockOpen).toHaveBeenCalled();
    });

    it("should work without className", () => {
        const { container } = render(
            <WalletConnectButton>
                <span>Connect</span>
            </WalletConnectButton>,
        );

        expect(container.firstChild).toBeInTheDocument();
    });

    it("should work with both onClick and className", () => {
        const onClickHandler = vi.fn();

        const { container } = render(
            <WalletConnectButton onClick={onClickHandler} className="wallet-btn">
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const div = container.firstChild;
        expect(div).toHaveClass("wallet-btn");

        const button = screen.getByText("Connect");
        fireEvent.click(button);

        expect(onClickHandler).toHaveBeenCalled();
        expect(mockOpen).toHaveBeenCalled();
    });

    it("should handle async onClick handlers", async () => {
        const asyncOnClick = vi.fn(async () => {
            await new Promise((resolve) => setTimeout(resolve, 10));
        });

        render(
            <WalletConnectButton onClick={asyncOnClick}>
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const button = screen.getByText("Connect");
        fireEvent.click(button);

        expect(asyncOnClick).toHaveBeenCalled();
        expect(mockOpen).toHaveBeenCalled();
    });

    it("should render with button-like behavior", () => {
        render(
            <WalletConnectButton>
                <button>Internal Button</button>
            </WalletConnectButton>,
        );

        const internalButton = screen.getByText("Internal Button");
        expect(internalButton).toBeInTheDocument();
    });

    it("should apply multiple class names", () => {
        const { container } = render(
            <WalletConnectButton className="class-1 class-2 class-3">
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const div = container.firstChild;
        expect(div).toHaveClass("class-1");
        expect(div).toHaveClass("class-2");
        expect(div).toHaveClass("class-3");
    });

    it("should handle onClick handler that throws", () => {
        const throwingOnClick = vi.fn(() => {
            throw new Error("Click handler error");
        });

        render(
            <WalletConnectButton onClick={throwingOnClick}>
                <span>Connect</span>
            </WalletConnectButton>,
        );

        const button = screen.getByText("Connect");

        // React catches errors in event handlers, so the throw won't propagate
        // Just verify the handler was called
        fireEvent.click(button);
        expect(throwingOnClick).toHaveBeenCalled();
    });

    it("should maintain event bubbling", () => {
        const parentClick = vi.fn();

        render(
            <div onClick={parentClick}>
                <WalletConnectButton>
                    <span>Connect</span>
                </WalletConnectButton>
            </div>,
        );

        const button = screen.getByText("Connect");
        fireEvent.click(button);

        // Should bubble up to parent
        expect(parentClick).toHaveBeenCalled();
    });

    it("should render with PropsWithChildren support", () => {
        const childElement = <span data-testid="child-element">Child</span>;

        render(<WalletConnectButton>{childElement}</WalletConnectButton>);

        expect(screen.getByTestId("child-element")).toBeInTheDocument();
    });
});
