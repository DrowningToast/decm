import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { LogoutButton } from './LogoutButton';
import { useNavigate } from '@/router';

// Mock the router
vi.mock('@/router', () => ({
    useNavigate: vi.fn(),
}));

// Mock i18next
vi.mock('react-i18next', () => ({
    useTranslation: () => ({
        t: (key: string) => {
            const translations: Record<string, string> = {
                'onboard.logout': 'Logout',
                'onboard.disconnect': 'Disconnect',
            };
            return translations[key] || key;
        },
    }),
}));

describe('LogoutButton', () => {
    const mockNavigate = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useNavigate as ReturnType<typeof vi.fn>).mockReturnValue(mockNavigate);
    });

    it('should render signout button', () => {
        render(<LogoutButton type="signout" />);

        const button = screen.getByRole('button');
        expect(button).toBeInTheDocument();
        expect(button).toHaveTextContent('Logout');
    });

    it('should render disconnect button', () => {
        render(<LogoutButton type="disconnect" />);

        const button = screen.getByRole('button');
        expect(button).toBeInTheDocument();
        expect(button).toHaveTextContent('Disconnect');
    });

    it('should navigate to signout on click', async () => {
        const user = userEvent.setup();
        render(<LogoutButton type="signout" />);

        const button = screen.getByRole('button');
        await user.click(button);

        expect(mockNavigate).toHaveBeenCalledWith('/signout');
    });

    it('should apply custom className', () => {
        const { container } = render(<LogoutButton type="signout" className="custom-class" />);

        const button = container.querySelector('button');
        expect(button?.className).toContain('custom-class');
    });

    it('should have correct button type', () => {
        render(<LogoutButton type="signout" />);

        const button = screen.getByRole('button');
        expect(button).toHaveAttribute('type', 'button');
    });
});

describe('LogoutButton - Additional Edge Cases', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        (useNavigate as ReturnType<typeof vi.fn>).mockReturnValue(mockNavigate);
    });

    it('should handle multiple rapid clicks gracefully', async () => {
        const user = userEvent.setup();
        render(<LogoutButton type="signout" />);

        const button = screen.getByRole('button');
        
        // Click multiple times rapidly
        await user.click(button);
        await user.click(button);
        await user.click(button);

        // Navigate should be called for each click
        expect(mockNavigate).toHaveBeenCalledTimes(3);
        expect(mockNavigate).toHaveBeenCalledWith('/signout');
    });

    it('should maintain correct button state when type prop changes', () => {
        const { rerender } = render(<LogoutButton type="signout" />);
        
        expect(screen.getByRole('button')).toHaveTextContent('Logout');
        
        rerender(<LogoutButton type="disconnect" />);
        
        expect(screen.getByRole('button')).toHaveTextContent('Disconnect');
    });

    it('should have correct accessibility attributes', () => {
        render(<LogoutButton type="signout" />);
        
        const button = screen.getByRole('button');
        expect(button).toHaveAttribute('type', 'button');
        expect(button).not.toBeDisabled();
    });

    it('should render with empty className when not provided', () => {
        const { container } = render(<LogoutButton type="signout" />);
        
        const button = container.querySelector('button');
        expect(button?.className).toBeTruthy(); // Has base classes
    });

    it('should merge multiple className values correctly', () => {
        const { container } = render(
            <LogoutButton type="signout" className="custom-1 custom-2" />
        );
        
        const button = container.querySelector('button');
        expect(button?.className).toContain('custom-1');
        expect(button?.className).toContain('custom-2');
    });
});