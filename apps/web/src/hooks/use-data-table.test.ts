import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, waitFor, act } from "@testing-library/react";
import { useDataTable } from "./use-data-table";

describe("useDataTable", () => {
    beforeEach(() => {
        vi.useRealTimers();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should initialize with default values", () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        expect(result.current.data).toEqual([]);
        expect(result.current.totalItems).toBe(0);
        expect(result.current.currentPage).toBe(1);
        expect(result.current.pageSize).toBe(10);
        expect(result.current.searchValue).toBe("");
        expect(result.current.sorting).toEqual([]);
        expect(result.current.isLoading).toBe(true);
        expect(result.current.error).toBeNull();
    });

    it("should initialize with custom page size", () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData, initialPageSize: 20 }));

        expect(result.current.pageSize).toBe(20);
    });

    it("should fetch data on mount", async () => {
        const mockData = [{ id: 1, name: "Item 1" }];
        const fetchData = vi.fn().mockResolvedValue({ data: mockData, total: 1 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        await waitFor(
            () => {
                expect(result.current.isLoading).toBe(false);
            },
            { timeout: 1000 },
        );

        expect(fetchData).toHaveBeenCalledWith({
            page: 1,
            pageSize: 10,
            search: "",
            sortBy: undefined,
            sortOrder: undefined,
        });
        expect(result.current.data).toEqual(mockData);
        expect(result.current.totalItems).toBe(1);
    });

    it("should debounce search input", async () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        // Wait for initial load
        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        vi.clearAllMocks();

        // Set search value
        act(() => {
            result.current.setSearchValue("test");
        });

        // Should not fetch immediately
        expect(fetchData).not.toHaveBeenCalled();

        // Wait for debounce (500ms) + fetch to complete
        await waitFor(
            () => {
                expect(fetchData).toHaveBeenCalled();
            },
            { timeout: 2000 },
        );

        expect(fetchData).toHaveBeenCalledWith(
            expect.objectContaining({
                search: "test",
            }),
        );
    });

    it("should reset to page 1 when search changes", async () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        // Wait for initial load
        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // Change page
        act(() => {
            result.current.setCurrentPage(3);
        });
        await waitFor(() => {
            expect(result.current.currentPage).toBe(3);
        });

        vi.clearAllMocks();

        // Change search
        act(() => {
            result.current.setSearchValue("test");
        });

        // Wait for debounce and page reset
        await waitFor(
            () => {
                expect(result.current.currentPage).toBe(1);
            },
            { timeout: 2000 },
        );
    });

    it("should update page and fetch data", async () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        // Wait for initial load
        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        vi.clearAllMocks();

        act(() => {
            result.current.setCurrentPage(2);
        });

        await waitFor(
            () => {
                expect(fetchData).toHaveBeenCalledWith(
                    expect.objectContaining({
                        page: 2,
                    }),
                );
            },
            { timeout: 1000 },
        );
    });

    it("should update page size and fetch data", async () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        // Wait for initial load
        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        vi.clearAllMocks();

        act(() => {
            result.current.setPageSize(20);
        });

        await waitFor(
            () => {
                expect(fetchData).toHaveBeenCalledWith(
                    expect.objectContaining({
                        pageSize: 20,
                    }),
                );
            },
            { timeout: 1000 },
        );
        expect(result.current.pageSize).toBe(20);
    });

    it("should handle sorting", async () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        // Wait for initial load
        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        vi.clearAllMocks();

        act(() => {
            result.current.setSorting([{ id: "name", desc: false }]);
        });

        await waitFor(
            () => {
                expect(fetchData).toHaveBeenCalledWith(
                    expect.objectContaining({
                        sortBy: "name",
                        sortOrder: "asc",
                    }),
                );
            },
            { timeout: 1000 },
        );
    });

    it("should handle descending sort", async () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        // Wait for initial load
        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        vi.clearAllMocks();

        act(() => {
            result.current.setSorting([{ id: "name", desc: true }]);
        });

        await waitFor(
            () => {
                expect(fetchData).toHaveBeenCalledWith(
                    expect.objectContaining({
                        sortBy: "name",
                        sortOrder: "desc",
                    }),
                );
            },
            { timeout: 1000 },
        );
    });

    it("should handle fetch errors", async () => {
        const error = new Error("Fetch failed");
        const fetchData = vi.fn().mockRejectedValue(error);

        const { result } = renderHook(() => useDataTable({ fetchData }));

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.error).toBe(error);
        expect(result.current.data).toEqual([]);
    });

    it("should refetch data when refetch is called", async () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        // Wait for initial load
        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        vi.clearAllMocks();

        act(() => {
            result.current.refetch();
        });

        await waitFor(
            () => {
                expect(fetchData).toHaveBeenCalled();
            },
            { timeout: 1000 },
        );
    });

    it("should handle multiple sorting states", async () => {
        const fetchData = vi.fn().mockResolvedValue({ data: [], total: 0 });

        const { result } = renderHook(() => useDataTable({ fetchData }));

        // Wait for initial load
        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        vi.clearAllMocks();

        // Set multiple sort columns (only first one should be used)
        act(() => {
            result.current.setSorting([
                { id: "name", desc: false },
                { id: "date", desc: true },
            ]);
        });

        await waitFor(
            () => {
                expect(fetchData).toHaveBeenCalledWith(
                    expect.objectContaining({
                        sortBy: "name",
                        sortOrder: "asc",
                    }),
                );
            },
            { timeout: 1000 },
        );
    });
});
