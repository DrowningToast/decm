import { useState, useEffect, useCallback } from "react";
import type { SortingState } from "@tanstack/react-table";

export interface UseDataTableProps<TData> {
    fetchData: (params: {
        page: number;
        pageSize: number;
        search: string;
        sortBy?: string;
        sortOrder?: "asc" | "desc";
    }) => Promise<{
        data: TData[];
        total: number;
    }>;
    initialPageSize?: number;
}

export interface UseDataTableReturn<TData> {
    data: TData[];
    totalItems: number;
    currentPage: number;
    pageSize: number;
    searchValue: string;
    sorting: SortingState;
    isLoading: boolean;
    error: Error | null;
    setCurrentPage: (page: number) => void;
    setPageSize: (size: number) => void;
    setSearchValue: (search: string) => void;
    setSorting: (sorting: SortingState | ((old: SortingState) => SortingState)) => void;
    refetch: () => void;
}

export function useDataTable<TData>({
    fetchData,
    initialPageSize = 10,
}: UseDataTableProps<TData>): UseDataTableReturn<TData> {
    const [data, setData] = useState<TData[]>([]);
    const [totalItems, setTotalItems] = useState(0);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(initialPageSize);
    const [searchValue, setSearchValue] = useState("");
    const [sorting, setSorting] = useState<SortingState>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    // Debounce search
    const [debouncedSearch, setDebouncedSearch] = useState(searchValue);

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(searchValue);
            // Reset to first page when search changes
            if (searchValue !== debouncedSearch) {
                setCurrentPage(1);
            }
        }, 500);

        return () => clearTimeout(timer);
    }, [searchValue, debouncedSearch]);

    const loadData = useCallback(async () => {
        setIsLoading(true);
        setError(null);

        try {
            const sortBy = sorting[0]?.id;
            const sortOrder = sortBy ? (sorting[0]?.desc ? "desc" : "asc") : undefined;

            const result = await fetchData({
                page: currentPage,
                pageSize,
                search: debouncedSearch,
                sortBy,
                sortOrder,
            });

            setData(result.data);
            setTotalItems(result.total);
        } catch (err) {
            setError(err as Error);
            console.error("Error fetching data:", err);
        } finally {
            setIsLoading(false);
        }
    }, [currentPage, pageSize, debouncedSearch, sorting, fetchData]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    return {
        data,
        totalItems,
        currentPage,
        pageSize,
        searchValue,
        sorting,
        isLoading,
        error,
        setCurrentPage,
        setPageSize,
        setSearchValue,
        setSorting,
        refetch: loadData,
    };
}
