import { describe, it, expect } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { usePaginationState } from "./usePaginationState";

describe("usePaginationState", () => {
    it("should initialize with default values", () => {
        const { result } = renderHook(() => usePaginationState());

        expect(result.current.page).toBe(1);
        expect(result.current.rowsPerPage).toBe(10);
        expect(result.current.offset).toBe(0);
    });

    it("should initialize with custom values", () => {
        const { result } = renderHook(() => usePaginationState(3, 25));

        expect(result.current.page).toBe(3);
        expect(result.current.rowsPerPage).toBe(25);
        expect(result.current.offset).toBe(50); // (3 - 1) * 25
    });

    it("should calculate correct offset", () => {
        const { result } = renderHook(() => usePaginationState(1, 10));

        expect(result.current.offset).toBe(0); // (1 - 1) * 10

        act(() => {
            result.current.handlePageChange(2);
        });

        expect(result.current.offset).toBe(10); // (2 - 1) * 10

        act(() => {
            result.current.handlePageChange(5);
        });

        expect(result.current.offset).toBe(40); // (5 - 1) * 10
    });

    it("should handle page change", () => {
        const { result } = renderHook(() => usePaginationState());

        act(() => {
            result.current.handlePageChange(3);
        });

        expect(result.current.page).toBe(3);
        expect(result.current.offset).toBe(20); // (3 - 1) * 10
    });

    it("should handle rows per page change", () => {
        const { result } = renderHook(() => usePaginationState(3, 10));

        act(() => {
            result.current.handleRowsPerPageChange(25);
        });

        expect(result.current.rowsPerPage).toBe(25);
        expect(result.current.page).toBe(1); // Should reset to page 1
        expect(result.current.offset).toBe(0); // (1 - 1) * 25
    });

    it("should reset to page 1 when changing rows per page", () => {
        const { result } = renderHook(() => usePaginationState());

        act(() => {
            result.current.handlePageChange(5);
        });

        expect(result.current.page).toBe(5);

        act(() => {
            result.current.handleRowsPerPageChange(20);
        });

        expect(result.current.page).toBe(1);
        expect(result.current.rowsPerPage).toBe(20);
    });

    it("should maintain stable function references", () => {
        const { result, rerender } = renderHook(() => usePaginationState());

        const handlePageChange1 = result.current.handlePageChange;
        const handleRowsPerPageChange1 = result.current.handleRowsPerPageChange;

        rerender();

        const handlePageChange2 = result.current.handlePageChange;
        const handleRowsPerPageChange2 = result.current.handleRowsPerPageChange;

        expect(handlePageChange1).toBe(handlePageChange2);
        expect(handleRowsPerPageChange1).toBe(handleRowsPerPageChange2);
    });

    it("should handle complex pagination scenario", () => {
        const { result } = renderHook(() => usePaginationState(1, 10));

        // Navigate to page 3
        act(() => {
            result.current.handlePageChange(3);
        });
        expect(result.current.page).toBe(3);
        expect(result.current.offset).toBe(20);

        // Change rows per page (should reset to page 1)
        act(() => {
            result.current.handleRowsPerPageChange(50);
        });
        expect(result.current.page).toBe(1);
        expect(result.current.rowsPerPage).toBe(50);
        expect(result.current.offset).toBe(0);

        // Navigate to page 2 with new rows per page
        act(() => {
            result.current.handlePageChange(2);
        });
        expect(result.current.page).toBe(2);
        expect(result.current.offset).toBe(50);
    });
});
