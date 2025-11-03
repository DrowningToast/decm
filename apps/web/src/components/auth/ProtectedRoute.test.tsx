import { describe, it, expect, beforeEach, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { ProtectedRoute } from "./ProtectedRoute";
import * as AuthContextModule from "@/context/AuthContext";
import { BrowserRouter } from "react-router-dom";
import { toast } from "sonner";

// Mock dependencies
vi.mock("sonner", () => ({
    toast: {
        error: vi.fn(),
    },
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

vi.mock("react-router-dom", async () => {
    const actual = await vi.importActual("react-router-dom");
    return {
        ...actual,
        Navigate: ({ to }: { to: string }) => <div data-testid="navigate">{to}</div>,
    };
});

describe("ProtectedRoute", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    const wrapper = ({ children }: { children: React.ReactNode }) => (
        <BrowserRouter>{children}</BrowserRouter>
    );

    describe("Loading State", () => {
        it("should show default loading state when isPending is true", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: true,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.getByText("common.loading")).toBeInTheDocument();
            expect(screen.queryByText("Protected Content")).not.toBeInTheDocument();
        });

        it("should show custom fallback when isPending is true and fallback is provided", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: true,
                refetch: vi.fn(),
                user: null,
            });

            const customFallback = <div data-testid="custom-fallback">Custom Loading...</div>;

            render(
                <ProtectedRoute fallback={customFallback}>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.getByTestId("custom-fallback")).toBeInTheDocument();
            expect(screen.getByText("Custom Loading...")).toBeInTheDocument();
            expect(screen.queryByText("common.loading")).not.toBeInTheDocument();
            expect(screen.queryByText("Protected Content")).not.toBeInTheDocument();
        });

        it("should display spinner in default loading state", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: true,
                refetch: vi.fn(),
                user: null,
            });

            const { container } = render(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            const spinner = container.querySelector(".animate-spin");
            expect(spinner).toBeInTheDocument();
        });
    });

    describe("Unauthenticated State", () => {
        it("should redirect to default signup page when not authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: false,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            const navigate = screen.getByTestId("navigate");
            expect(navigate).toBeInTheDocument();
            expect(navigate).toHaveTextContent("/signup");
            expect(screen.queryByText("Protected Content")).not.toBeInTheDocument();
        });

        it("should redirect to custom path when redirectTo is provided", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: false,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute redirectTo="/signin">
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            const navigate = screen.getByTestId("navigate");
            expect(navigate).toBeInTheDocument();
            expect(navigate).toHaveTextContent("/signin");
        });

        it("should show toast error when not authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: false,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(toast.error).toHaveBeenCalledWith("flow.generic.unauthenticated_response");
            expect(toast.error).toHaveBeenCalledOnce();
        });

        it("should not render children when not authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: false,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.queryByTestId("protected-content")).not.toBeInTheDocument();
        });
    });

    describe("Authenticated State", () => {
        it("should render children when authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            render(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.getByTestId("protected-content")).toBeInTheDocument();
            expect(screen.getByText("Protected Content")).toBeInTheDocument();
        });

        it("should not redirect when authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            render(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.queryByTestId("navigate")).not.toBeInTheDocument();
            expect(screen.getByTestId("protected-content")).toBeInTheDocument();
        });

        it("should not show toast when authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            render(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(toast.error).not.toHaveBeenCalled();
        });

        it("should not show loading state when authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            render(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.queryByText("common.loading")).not.toBeInTheDocument();
            expect(screen.getByTestId("protected-content")).toBeInTheDocument();
        });
    });

    describe("Complex Children", () => {
        it("should render complex component tree when authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            render(
                <ProtectedRoute>
                    <div data-testid="parent">
                        <h1>Dashboard</h1>
                        <div data-testid="child">
                            <p>Nested Content</p>
                        </div>
                    </div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.getByTestId("parent")).toBeInTheDocument();
            expect(screen.getByText("Dashboard")).toBeInTheDocument();
            expect(screen.getByTestId("child")).toBeInTheDocument();
            expect(screen.getByText("Nested Content")).toBeInTheDocument();
        });

        it("should render multiple children when authenticated", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            render(
                <ProtectedRoute>
                    <>
                        <div data-testid="element-1">Element 1</div>
                        <div data-testid="element-2">Element 2</div>
                        <div data-testid="element-3">Element 3</div>
                    </>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.getByTestId("element-1")).toBeInTheDocument();
            expect(screen.getByTestId("element-2")).toBeInTheDocument();
            expect(screen.getByTestId("element-3")).toBeInTheDocument();
        });
    });

    describe("Edge Cases", () => {
        it("should handle transition from pending to authenticated", () => {
            const { rerender } = render(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            // Initially pending
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: true,
                refetch: vi.fn(),
                user: null,
            });

            rerender(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
            );

            expect(screen.getByText("common.loading")).toBeInTheDocument();

            // Then authenticated
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            rerender(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
            );

            expect(screen.queryByText("common.loading")).not.toBeInTheDocument();
            expect(screen.getByTestId("protected-content")).toBeInTheDocument();
        });

        it("should handle transition from pending to unauthenticated", () => {
            const { rerender } = render(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            // Initially pending
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: true,
                refetch: vi.fn(),
                user: null,
            });

            rerender(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
            );

            expect(screen.getByText("common.loading")).toBeInTheDocument();

            // Then unauthenticated
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: false,
                refetch: vi.fn(),
                user: null,
            });

            rerender(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
            );

            expect(screen.queryByText("common.loading")).not.toBeInTheDocument();
            expect(screen.getByTestId("navigate")).toBeInTheDocument();
        });

        it("should handle null children gracefully", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            const { container } = render(<ProtectedRoute>{null}</ProtectedRoute>, { wrapper });

            // Should render without errors
            expect(container).toBeInTheDocument();
        });

        it("should handle undefined children gracefully", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            const { container } = render(
                <ProtectedRoute>{undefined as unknown as React.ReactElement}</ProtectedRoute>,
                { wrapper },
            );

            // Should render without errors
            expect(container).toBeInTheDocument();
        });
    });

    describe("Props Validation", () => {
        it("should use default redirectTo when not provided", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: false,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            const navigate = screen.getByTestId("navigate");
            expect(navigate).toHaveTextContent("/signup");
        });

        it("should accept valid path as redirectTo", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: false,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            const navigate = screen.getByTestId("navigate");
            expect(navigate).toHaveTextContent("/signin");
        });

        it("should render custom fallback component", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: true,
                refetch: vi.fn(),
                user: null,
            });

            const CustomFallback = () => (
                <div data-testid="custom-loader">
                    <span>Please wait...</span>
                </div>
            );

            render(
                <ProtectedRoute fallback={<CustomFallback />}>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.getByTestId("custom-loader")).toBeInTheDocument();
            expect(screen.getByText("Please wait...")).toBeInTheDocument();
        });
    });

    describe("Typography Component", () => {
        it("should render Typography component in loading state", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isPending: true,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            // Typography component should render the loading text
            const loadingText = screen.getByText("common.loading");
            expect(loadingText).toBeInTheDocument();
            expect(loadingText.tagName).toBe("P");
        });
    });

    describe("Role-Based Protection (TODO)", () => {
        it("should render children when role-based protection is not implemented", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isPending: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            // Note: requiredRoles prop is commented out in component
            render(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            expect(screen.getByTestId("protected-content")).toBeInTheDocument();
        });
    });
});
