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
