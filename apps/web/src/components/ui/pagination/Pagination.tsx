import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Typography } from "@/components/typography/typography";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";

interface DataTablePaginationProps {
    currentPage: number;
    totalPages?: number;
    hasNextPage: boolean;
    hasPreviousPage: boolean;
    onPageChange: (page: number) => void;
    onRowsPerPageChange?: (rowsPerPage: number) => void;
    isLoading?: boolean;
    className?: string;
    showPageNumbers?: boolean;
    maxVisiblePages?: number;
    rowsPerPage?: number;
    rowsPerPageOptions?: number[];
    totalItems?: number;
}

export const DataTablePagination = ({
    currentPage,
    totalPages,
    hasNextPage,
    hasPreviousPage,
    onPageChange,
    onRowsPerPageChange,
    isLoading = false,
    className = "",
    showPageNumbers = true,
    maxVisiblePages = 5,
    rowsPerPage = 10,
    rowsPerPageOptions = [10, 20, 50, 100],
    totalItems,
}: DataTablePaginationProps) => {
    const { t } = useTranslation();
    const handlePrevious = () => {
        if (hasPreviousPage && !isLoading) {
            onPageChange(currentPage - 1);
        }
    };

    const handleNext = () => {
        if (hasNextPage && !isLoading) {
            onPageChange(currentPage + 1);
        }
    };

    const handlePageClick = (page: number) => {
        if (!isLoading) {
            onPageChange(page);
        }
    };

    const handleRowsPerPageChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const newRowsPerPage = parseInt(e.target.value, 10);
        if (onRowsPerPageChange) {
            onRowsPerPageChange(newRowsPerPage);
        }
    };

    const renderPageNumbers = () => {
        if (!showPageNumbers || !totalPages) return null;

        // Calculate the range of page numbers to show
        let startPage = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2));
        const endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);

        // Adjust if we're near the end
        if (endPage - startPage + 1 < maxVisiblePages) {
            startPage = Math.max(1, endPage - maxVisiblePages + 1);
        }

        const pages = [];
        for (let i = startPage; i <= endPage; i++) {
            pages.push(
                <Button
                    key={i}
                    variant={i === currentPage ? "primary" : "secondary-dark"}
                    size="sm"
                    onClick={() => handlePageClick(i)}
                    disabled={isLoading}
                    className={cn(
                        "w-8 h-8 p-0",
                        i === currentPage && "bg-primary text-primary-foreground",
                    )}
                >
                    <Typography variant="text" tag="span">
                        {i}
                    </Typography>
                </Button>,
            );
        }

        return pages;
    };

    // Calculate total pages if not provided but totalItems is
    const calculatedTotalPages =
        totalPages || (totalItems ? Math.ceil(totalItems / rowsPerPage) : undefined);

    return (
        <div className={cn("flex items-center justify-between w-full", className)}>
            {/* Rows per page selector */}
            <div className="flex items-center space-x-2">
                <Typography variant="text" tag="span" className="text-sm">
                    {t("common.pagination.rowsPerPage")}
                </Typography>
                {onRowsPerPageChange ? (
                    <select
                        value={rowsPerPage}
                        onChange={handleRowsPerPageChange}
                        className="border border-gray-300 rounded px-2 py-1 text-sm"
                        disabled={isLoading}
                    >
                        {rowsPerPageOptions.map((option) => (
                            <option key={option} value={option}>
                                {option}
                            </option>
                        ))}
                    </select>
                ) : (
                    <Typography variant="text" tag="span" className="text-sm font-medium">
                        {rowsPerPage}
                    </Typography>
                )}
            </div>

            {/* Page information */}
            <div className="flex items-center">
                <Typography variant="text" tag="span" className="text-sm text-muted-foreground">
                    {calculatedTotalPages
                        ? t("common.pagination.pageOf", {
                              current: currentPage,
                              total: calculatedTotalPages,
                          })
                        : t("common.pagination.page", { current: currentPage })}
                </Typography>
            </div>

            {/* Pagination buttons */}
            <div className="flex items-center space-x-1">
                <Button
                    variant="secondary-dark"
                    size="sm"
                    onClick={handlePrevious}
                    disabled={!hasPreviousPage || isLoading}
                    className="flex items-center"
                >
                    <ChevronLeft className="h-4 w-4" />
                    <Typography variant="text" tag="span" className="ml-1">
                        {t("common.previous")}
                    </Typography>
                </Button>

                {showPageNumbers && calculatedTotalPages && (
                    <div className="flex items-center space-x-1">{renderPageNumbers()}</div>
                )}

                <Button
                    variant="secondary-dark"
                    size="sm"
                    onClick={handleNext}
                    disabled={!hasNextPage || isLoading}
                    className="flex items-center"
                >
                    <Typography variant="text" tag="span" className="mr-1">
                        {t("common.next")}
                    </Typography>
                    <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
        </div>
    );
};
