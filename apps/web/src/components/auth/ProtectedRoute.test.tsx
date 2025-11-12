import { describe, it, expect, beforeEach, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { ProtectedRoute } from "./ProtectedRoute";
import * as AuthContextModule from "@/context/AuthContext";
import { BrowserRouter } from "react-router-dom";
import { toast } from "sonner";
import * as UseCheckRolesModule from "@/hooks/useCheckRoles";

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

vi.mock("@/hooks/useCheckRoles", () => ({
    useCheckRoles: vi.fn(),
}));

describe("ProtectedRoute", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    const wrapper = ({ children }: { children: React.ReactNode }) => (
        <BrowserRouter>{children}</BrowserRouter>
    );

    describe("Loading State", () => {
        it("should show default loading state when isFetching is true", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isFetching: true,
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

        it("should show custom fallback when isFetching is true and fallback is provided", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isFetching: true,
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
                isFetching: true,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: true,
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
                isFetching: false,
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
                isFetching: true,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
                isFetching: false,
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
            expect(navigate).toHaveTextContent("/signin");
        });

        it("should render custom fallback component", () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isFetching: true,
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
                isFetching: true,
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

    describe("Role-Based Protection", () => {
        it("should render children when no roles are required", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: true,
                isLoading: false,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                expect(screen.getByTestId("protected-content")).toBeInTheDocument();
            });
            expect(UseCheckRolesModule.useCheckRoles).not.toHaveBeenCalled();
        });

        it("should check host role when requireHost is true", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: true,
                isHost: true,
                isLoading: false,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute requireHost={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                expect(UseCheckRolesModule.useCheckRoles).toHaveBeenCalledWith({
                    requireHost: true,
                    requireIssuer: false,
                    enabled: true,
                });
            });

            await waitFor(() => {
                expect(screen.getByTestId("protected-content")).toBeInTheDocument();
            });
        });

        it("should check issuer role when requireIssuer is true", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: true,
                isIssuer: true,
                isLoading: false,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute requireIssuer={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                expect(UseCheckRolesModule.useCheckRoles).toHaveBeenCalledWith({
                    requireHost: false,
                    requireIssuer: true,
                    enabled: true,
                });
            });

            await waitFor(() => {
                expect(screen.getByTestId("protected-content")).toBeInTheDocument();
            });
        });

        it("should check both roles when requireHost and requireIssuer are true", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: true,
                isHost: true,
                isIssuer: true,
                isLoading: false,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute requireHost={true} requireIssuer={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                expect(UseCheckRolesModule.useCheckRoles).toHaveBeenCalledWith({
                    requireHost: true,
                    requireIssuer: true,
                    enabled: true,
                });
            });

            await waitFor(() => {
                expect(screen.getByTestId("protected-content")).toBeInTheDocument();
            });
        });

        it("should redirect to unauthorized when host role check fails", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: false,
                isHost: false,
                isLoading: false,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute requireHost={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                const navigate = screen.getByTestId("navigate");
                expect(navigate).toBeInTheDocument();
                expect(navigate).toHaveTextContent("/unauthorized");
            });

            expect(toast.error).toHaveBeenCalledWith("flow.generic.unauthorized_response");
            expect(screen.queryByTestId("protected-content")).not.toBeInTheDocument();
        });

        it("should redirect to unauthorized when issuer role check fails", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: false,
                isIssuer: false,
                isLoading: false,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute requireIssuer={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                const navigate = screen.getByTestId("navigate");
                expect(navigate).toBeInTheDocument();
                expect(navigate).toHaveTextContent("/unauthorized");
            });

            expect(toast.error).toHaveBeenCalledWith("flow.generic.unauthorized_response");
            expect(screen.queryByTestId("protected-content")).not.toBeInTheDocument();
        });

        it("should redirect to custom unauthorized path when provided", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: false,
                isHost: false,
                isLoading: false,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute requireHost={true} unauthorizedRedirectTo="/error">
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                const navigate = screen.getByTestId("navigate");
                expect(navigate).toBeInTheDocument();
                expect(navigate).toHaveTextContent("/error");
            });
        });

        it("should show loading state during role check", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: false,
                isLoading: true,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute requireHost={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            // Should show loading during role check
            expect(screen.getByText("common.loading")).toBeInTheDocument();

            // After role check completes, update mock to return success
            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: true,
                isHost: true,
                isLoading: false,
                isError: false,
                error: null,
            });

            // Since we can't easily re-render with different hook results in this test structure,
            // we verify the loading state is shown during role check
            expect(screen.queryByTestId("protected-content")).not.toBeInTheDocument();
        });

        it("should handle role check error", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: false,
                isLoading: false,
                isError: true,
                error: new Error("Role check failed"),
            });

            render(
                <ProtectedRoute requireHost={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                const navigate = screen.getByTestId("navigate");
                expect(navigate).toBeInTheDocument();
                expect(navigate).toHaveTextContent("/unauthorized");
            });

            expect(toast.error).toHaveBeenCalledWith("flow.generic.unauthorized_response");
        });

        it("should not call useCheckRoles when user is not authenticated", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isFetching: false,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute requireHost={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                const navigate = screen.getByTestId("navigate");
                expect(navigate).toBeInTheDocument();
                expect(navigate).toHaveTextContent("/signup");
            });

            expect(UseCheckRolesModule.useCheckRoles).not.toHaveBeenCalled();
        });

        it("should render custom fallback during role check", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: false,
                isLoading: true,
                isError: false,
                error: null,
            });

            const customFallback = <div data-testid="custom-fallback">Checking roles...</div>;

            render(
                <ProtectedRoute requireHost={true} fallback={customFallback}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            // Should show custom fallback during role check
            expect(screen.getByTestId("custom-fallback")).toBeInTheDocument();
            expect(screen.queryByTestId("protected-content")).not.toBeInTheDocument();
        });

        it("should handle partial role satisfaction (host but not issuer)", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: true,
                isFetching: false,
                refetch: vi.fn(),
                user: { id: "user-123", first_name: "Test" },
            });

            vi.spyOn(UseCheckRolesModule, "useCheckRoles").mockReturnValue({
                hasRequiredRoles: false,
                isHost: true,
                isIssuer: false,
                isLoading: false,
                isError: false,
                error: null,
            });

            render(
                <ProtectedRoute requireHost={true} requireIssuer={true}>
                    <div data-testid="protected-content">Protected Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                const navigate = screen.getByTestId("navigate");
                expect(navigate).toBeInTheDocument();
                expect(navigate).toHaveTextContent("/unauthorized");
            });

            expect(toast.error).toHaveBeenCalledWith("flow.generic.unauthorized_response");
        });

        it("should allow requireAuthenticated to be set to false", async () => {
            vi.spyOn(AuthContextModule, "useAuth").mockReturnValue({
                isAuthenticated: false,
                isFetching: false,
                refetch: vi.fn(),
                user: null,
            });

            render(
                <ProtectedRoute requireAuthenticated={false}>
                    <div data-testid="protected-content">Public Content</div>
                </ProtectedRoute>,
                { wrapper },
            );

            await waitFor(() => {
                expect(screen.getByTestId("protected-content")).toBeInTheDocument();
            });

            expect(toast.error).not.toHaveBeenCalled();
            expect(UseCheckRolesModule.useCheckRoles).not.toHaveBeenCalled();
        });
    });
});
