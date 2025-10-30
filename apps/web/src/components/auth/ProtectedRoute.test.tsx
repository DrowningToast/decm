import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ProtectedRoute } from './auth/ProtectedRoute';
import { useAuth } from '@/context/AuthContext';

// Mock AuthContext
vi.mock('@/context/AuthContext', () => ({
    useAuth: vi.fn(),
}));

// Mock router
vi.mock('@/router', () => ({
    Navigate: ({ to }: { to: string }) => <div data-testid="navigate" data-to={to} />,
}));

// Mock i18next
vi.mock('react-i18next', () => ({
    useTranslation: () => ({
        t: (key: string) => {
            const translations: Record<string, string> = {
                'common.loading': 'Loading...',
            };
            return translations[key] || key;
        },
    }),
}));

describe('ProtectedRoute', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should render children when authenticated', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: undefined,
        });

        render(
            <ProtectedRoute>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });

    it('should show loading state when isLoading is true', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: false,
            isLoading: true,
            userRole: undefined,
        });

        render(
            <ProtectedRoute>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Loading...')).toBeInTheDocument();
        expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });

    it('should show custom fallback when loading', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: false,
            isLoading: true,
            userRole: undefined,
        });

        render(
            <ProtectedRoute fallback={<div>Custom Loading...</div>}>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Custom Loading...')).toBeInTheDocument();
        expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });

    it('should redirect to signup when not authenticated', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: false,
            isLoading: false,
            userRole: undefined,
        });

        render(
            <ProtectedRoute>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        const navigate = screen.getByTestId('navigate');
        expect(navigate).toBeInTheDocument();
        expect(navigate).toHaveAttribute('data-to', '/signup');
    });

    it('should redirect to custom path when not authenticated', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: false,
            isLoading: false,
            userRole: undefined,
        });

        render(
            <ProtectedRoute redirectTo="/custom-login">
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        const navigate = screen.getByTestId('navigate');
        expect(navigate).toHaveAttribute('data-to', '/custom-login');
    });

    it('should allow access when user has required role', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: 'ADMIN',
        });

        render(
            <ProtectedRoute requiredRoles={['ADMIN', 'ISSUER']}>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });

    it('should redirect when user lacks required role', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: 'PARTICIPANT',
        });

        render(
            <ProtectedRoute requiredRoles={['ADMIN', 'ISSUER']}>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        const navigate = screen.getByTestId('navigate');
        expect(navigate).toHaveAttribute('data-to', '/unauthorized');
    });

    it('should allow access when no roles specified', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: 'PARTICIPANT',
        });

        render(
            <ProtectedRoute>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });

    it('should handle multiple required roles', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: 'ISSUER',
        });

        render(
            <ProtectedRoute requiredRoles={['ADMIN', 'ISSUER', 'HOST']}>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
});

describe('ProtectedRoute - Additional Edge Cases', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should handle userRole as undefined when authenticated', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: undefined,
        });

        render(
            <ProtectedRoute>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });

    it('should handle empty requiredRoles array', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: 'PARTICIPANT',
        });

        render(
            <ProtectedRoute requiredRoles={[]}>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });

    it('should transition from loading to authenticated', () => {
        const { rerender } = render(<ProtectedRoute><div>Content</div></ProtectedRoute>);

        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: false,
            isLoading: true,
            userRole: undefined,
        });
        rerender(<ProtectedRoute><div>Content</div></ProtectedRoute>);
        expect(screen.getByText('Loading...')).toBeInTheDocument();

        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: 'ADMIN',
        });
        rerender(<ProtectedRoute><div>Content</div></ProtectedRoute>);
        expect(screen.getByText('Content')).toBeInTheDocument();
    });

    it('should handle case-sensitive role matching', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: 'admin' as any, // lowercase
        });

        render(
            <ProtectedRoute requiredRoles={['ADMIN']}>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        const navigate = screen.getByTestId('navigate');
        expect(navigate).toHaveAttribute('data-to', '/unauthorized');
    });

    it('should render complex children correctly', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: 'ADMIN',
        });

        render(
            <ProtectedRoute>
                <div>
                    <h1>Title</h1>
                    <p>Paragraph</p>
                    <button>Action</button>
                </div>
            </ProtectedRoute>,
        );

        expect(screen.getByText('Title')).toBeInTheDocument();
        expect(screen.getByText('Paragraph')).toBeInTheDocument();
        expect(screen.getByText('Action')).toBeInTheDocument();
    });

    it('should support custom redirectTo paths with query params', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: false,
            isLoading: false,
            userRole: undefined,
        });

        render(
            <ProtectedRoute redirectTo="/login?returnUrl=/dashboard">
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        const navigate = screen.getByTestId('navigate');
        expect(navigate).toHaveAttribute('data-to', '/login?returnUrl=/dashboard');
    });

    it('should handle multiple required roles with userRole null', () => {
        (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
            isAuthenticated: true,
            isLoading: false,
            userRole: null,
        });

        render(
            <ProtectedRoute requiredRoles={['ADMIN', 'ISSUER']}>
                <div>Protected Content</div>
            </ProtectedRoute>,
        );

        // Should allow access since requiredRoles is checked but userRole is null
        expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
});