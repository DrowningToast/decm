import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { usePaginationState } from './usePaginationState';

describe('usePaginationState', () => {
    it('should initialize with default values', () => {
        const { result } = renderHook(() => usePaginationState());

        expect(result.current.page).toBe(1);
        expect(result.current.rowsPerPage).toBe(10);
        expect(result.current.offset).toBe(0);
    });

    it('should initialize with custom values', () => {
        const { result } = renderHook(() => usePaginationState(3, 20));

        expect(result.current.page).toBe(3);
        expect(result.current.rowsPerPage).toBe(20);
        expect(result.current.offset).toBe(40); // (3 - 1) * 20
    });

    it('should update page correctly', () => {
        const { result } = renderHook(() => usePaginationState());

        act(() => {
            result.current.handlePageChange(5);
        });

        expect(result.current.page).toBe(5);
        expect(result.current.offset).toBe(40); // (5 - 1) * 10
    });

    it('should update rows per page and reset to page 1', () => {
        const { result } = renderHook(() => usePaginationState(5, 10));

        act(() => {
            result.current.handleRowsPerPageChange(25);
        });

        expect(result.current.rowsPerPage).toBe(25);
        expect(result.current.page).toBe(1);
        expect(result.current.offset).toBe(0);
    });

    it('should calculate offset correctly', () => {
        const { result } = renderHook(() => usePaginationState(4, 15));

        expect(result.current.offset).toBe(45); // (4 - 1) * 15

        act(() => {
            result.current.handlePageChange(7);
        });

        expect(result.current.offset).toBe(90); // (7 - 1) * 15
    });

    it('should handle page changes independently', () => {
        const { result } = renderHook(() => usePaginationState(1, 10));

        act(() => {
            result.current.handlePageChange(2);
        });

        expect(result.current.page).toBe(2);
        expect(result.current.rowsPerPage).toBe(10);

        act(() => {
            result.current.handlePageChange(10);
        });

        expect(result.current.page).toBe(10);
        expect(result.current.rowsPerPage).toBe(10);
    });

    it('should maintain stable function references', () => {
        const { result, rerender } = renderHook(() => usePaginationState());

        const initialHandlePageChange = result.current.handlePageChange;
        const initialHandleRowsPerPageChange = result.current.handleRowsPerPageChange;

        rerender();

        expect(result.current.handlePageChange).toBe(initialHandlePageChange);
        expect(result.current.handleRowsPerPageChange).toBe(initialHandleRowsPerPageChange);
    });
});

describe('usePaginationState - Additional Edge Cases', () => {
    it('should handle page 0 (edge case)', () => {
        const { result } = renderHook(() => usePaginationState(0, 10));

        expect(result.current.page).toBe(0);
        expect(result.current.offset).toBe(-10); // (0 - 1) * 10
    });

    it('should handle very large page numbers', () => {
        const { result } = renderHook(() => usePaginationState(1000, 10));

        expect(result.current.page).toBe(1000);
        expect(result.current.offset).toBe(9990); // (1000 - 1) * 10
    });

    it('should handle zero rows per page', () => {
        const { result } = renderHook(() => usePaginationState(1, 0));

        expect(result.current.rowsPerPage).toBe(0);
        expect(result.current.offset).toBe(0); // (1 - 1) * 0
    });

    it('should handle negative rows per page', () => {
        const { result } = renderHook(() => usePaginationState(5, -10));

        expect(result.current.rowsPerPage).toBe(-10);
        expect(result.current.offset).toBe(-40); // (5 - 1) * -10
    });

    it('should recalculate offset when both page and rowsPerPage change', () => {
        const { result } = renderHook(() => usePaginationState(1, 10));

        act(() => {
            result.current.handlePageChange(3);
        });

        expect(result.current.offset).toBe(20); // (3 - 1) * 10

        act(() => {
            result.current.handleRowsPerPageChange(25);
        });

        // After changing rows per page, page resets to 1
        expect(result.current.page).toBe(1);
        expect(result.current.offset).toBe(0); // (1 - 1) * 25
    });

    it('should handle rapid page changes', () => {
        const { result } = renderHook(() => usePaginationState(1, 10));

        act(() => {
            result.current.handlePageChange(2);
            result.current.handlePageChange(3);
            result.current.handlePageChange(4);
            result.current.handlePageChange(5);
        });

        expect(result.current.page).toBe(5);
        expect(result.current.offset).toBe(40); // (5 - 1) * 10
    });

    it('should handle fractional page numbers', () => {
        const { result } = renderHook(() => usePaginationState(1, 10));

        act(() => {
            result.current.handlePageChange(2.5);
        });

        expect(result.current.page).toBe(2.5);
        expect(result.current.offset).toBe(15); // (2.5 - 1) * 10
    });

    it('should handle fractional rows per page', () => {
        const { result } = renderHook(() => usePaginationState(2, 10));

        act(() => {
            result.current.handleRowsPerPageChange(7.5);
        });

        expect(result.current.rowsPerPage).toBe(7.5);
        expect(result.current.page).toBe(1);
        expect(result.current.offset).toBe(0);
    });

    it('should calculate correct offset for typical pagination scenarios', () => {
        const { result } = renderHook(() => usePaginationState(1, 20));

        // Page 1
        expect(result.current.offset).toBe(0);

        // Page 2
        act(() => {
            result.current.handlePageChange(2);
        });
        expect(result.current.offset).toBe(20);

        // Page 5
        act(() => {
            result.current.handlePageChange(5);
        });
        expect(result.current.offset).toBe(80);
    });

    it('should maintain state immutability', () => {
        const { result } = renderHook(() => usePaginationState(1, 10));

        const initialPage = result.current.page;
        const initialRowsPerPage = result.current.rowsPerPage;

        // Attempt to mutate (should not affect internal state)
        act(() => {
            result.current.handlePageChange(5);
        });

        expect(result.current.page).not.toBe(initialPage);
        expect(result.current.rowsPerPage).toBe(initialRowsPerPage);
    });

    it('should work with common pagination values', () => {
        const testCases = [
            { page: 1, rowsPerPage: 10, expectedOffset: 0 },
            { page: 1, rowsPerPage: 25, expectedOffset: 0 },
            { page: 1, rowsPerPage: 50, expectedOffset: 0 },
            { page: 1, rowsPerPage: 100, expectedOffset: 0 },
            { page: 5, rowsPerPage: 20, expectedOffset: 80 },
            { page: 10, rowsPerPage: 50, expectedOffset: 450 },
        ];

        testCases.forEach(({ page, rowsPerPage, expectedOffset }) => {
            const { result } = renderHook(() => usePaginationState(page, rowsPerPage));
            expect(result.current.offset).toBe(expectedOffset);
        });
    });
});