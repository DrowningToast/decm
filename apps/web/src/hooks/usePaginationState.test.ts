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
