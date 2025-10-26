import { useCallback, useState } from "react";

export function usePaginationState(initialPage: number = 1, initialRowsPerPage: number = 10) {
    const [page, setPage] = useState(initialPage);
    const [rowsPerPage, setRowsPerPage] = useState(initialRowsPerPage);
    const offset = (page - 1) * rowsPerPage;

    const handlePageChange = useCallback((newPage: number) => {
        setPage(newPage);
    }, []);

    const handleRowsPerPageChange = useCallback((newRowsPerPage: number) => {
        setRowsPerPage(newRowsPerPage);
        setPage(1);
    }, []);

    return { page, rowsPerPage, offset, handlePageChange, handleRowsPerPageChange };
}
