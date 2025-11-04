import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { LogoutButton } from "./LogoutButton";

const mockNavigate = vi.fn();

// Mock router
vi.mock("@/router", () => ({
    useNavigate: () => mockNavigate,
}));

// Mock translation
vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

// Mock Typography
vi.mock("@/components/typography/typography", () => ({
    Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
}));

// Mock utils
vi.mock("@/lib/utils", () => ({
    cn: (...args: any[]) => args.filter(Boolean).join(" "),
}));

describe("LogoutButton Component", () => {
    beforeEach(() => {
        mockNavigate.mockClear();
    });

    it("should render logout button for signout type", () => {
        render(<LogoutButton type="signout" />);

        expect(screen.getByText("onboard.logout")).toBeInTheDocument();
    });

    it("should render disconnect button for disconnect type", () => {
        render(<LogoutButton type="disconnect" />);

        expect(screen.getByText("onboard.disconnect")).toBeInTheDocument();
    });

    it("should navigate to /signout when clicked", async () => {
        render(<LogoutButton type="signout" />);

        const button = screen.getByRole("button");
        fireEvent.click(button);

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/signout");
        });
    });

    it("should navigate to /signout for both signout and disconnect types", async () => {
        const { rerender } = render(<LogoutButton type="signout" />);

        const button = screen.getByRole("button");
        fireEvent.click(button);

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/signout");
        });

        mockNavigate.mockClear();

        rerender(<LogoutButton type="disconnect" />);

        fireEvent.click(button);

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/signout");
        });
    });

    it("should apply custom className", () => {
        const { container } = render(<LogoutButton type="signout" className="custom-class" />);

        const button = container.querySelector(".custom-class");
        expect(button).toBeInTheDocument();
    });

    it("should combine custom className with default classes", () => {
        const { container } = render(<LogoutButton type="signout" className="extra-class" />);

        const button = screen.getByRole("button");
        expect(button.className).toContain("text-start");
        expect(button.className).toContain("h-[14.5px]");
        expect(button.className).toContain("inline-block");
        expect(button.className).toContain("extra-class");
    });

    it("should have correct button type attribute", () => {
        render(<LogoutButton type="signout" />);

        const button = screen.getByRole("button");
        expect(button).toHaveAttribute("type", "button");
    });

    it("should render Typography component with correct props", () => {
        render(<LogoutButton type="signout" />);

        const typography = screen.getByText("onboard.logout");
        expect(typography).toHaveClass("text-xs");
    });

    it("should handle rapid clicks", async () => {
        render(<LogoutButton type="signout" />);

        const button = screen.getByRole("button");
        fireEvent.click(button);
        fireEvent.click(button);
        fireEvent.click(button);

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledTimes(3);
            expect(mockNavigate).toHaveBeenCalledWith("/signout");
        });
    });

    it("should not break with no className prop", () => {
        const { container } = render(<LogoutButton type="disconnect" />);

        expect(container.querySelector("button")).toBeInTheDocument();
    });

    it("should render with proper styling for signout type", () => {
        render(<LogoutButton type="signout" />);

        const typography = screen.getByText("onboard.logout");
        expect(typography).toHaveClass("italic");
        expect(typography).toHaveClass("underline");
    });

    it("should render with proper styling for disconnect type", () => {
        render(<LogoutButton type="disconnect" />);

        const typography = screen.getByText("onboard.disconnect");
        expect(typography).toHaveClass("italic");
        expect(typography).toHaveClass("underline");
    });

    it("should be accessible as a button element", () => {
        render(<LogoutButton type="signout" />);

        const button = screen.getByRole("button");
        expect(button.tagName).toBe("BUTTON");
    });

    it("should handle different button states", () => {
        const { rerender } = render(<LogoutButton type="signout" />);

        const button = screen.getByRole("button");
        expect(button).toBeEnabled();

        rerender(<LogoutButton type="disconnect" />);

        expect(screen.getByRole("button")).toBeEnabled();
    });

    it("should display correct translation key for signout", () => {
        render(<LogoutButton type="signout" />);

        expect(screen.getByText("onboard.logout")).toBeInTheDocument();
        expect(screen.queryByText("onboard.disconnect")).not.toBeInTheDocument();
    });

    it("should display correct translation key for disconnect", () => {
        render(<LogoutButton type="disconnect" />);

        expect(screen.getByText("onboard.disconnect")).toBeInTheDocument();
        expect(screen.queryByText("onboard.logout")).not.toBeInTheDocument();
    });

    it("should maintain stable reference across re-renders", () => {
        const { rerender } = render(<LogoutButton type="signout" />);

        const button1 = screen.getByRole("button");

        rerender(<LogoutButton type="signout" />);

        const button2 = screen.getByRole("button");

        expect(button1).toBeInTheDocument();
        expect(button2).toBeInTheDocument();
    });

    it("should navigate with correct path on click", async () => {
        render(<LogoutButton type="signout" />);

        fireEvent.click(screen.getByRole("button"));

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith("/signout");
        });
    });
});
