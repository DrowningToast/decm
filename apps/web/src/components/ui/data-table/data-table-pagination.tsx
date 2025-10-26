import { Button } from "@/components/ui/button";
import { Typography } from "@/components/typography/typography";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from "lucide-react";

interface DataTablePaginationProps {
    totalItems?: number;
    currentPage: number;
    pageSize: number;
    onPageChange: (page: number) => void;
    onPageSizeChange: (pageSize: number) => void;
}

export function DataTablePagination({
    totalItems,
    currentPage,
    pageSize,
    onPageChange,
    onPageSizeChange,
}: DataTablePaginationProps) {
    const totalPages = totalItems ? Math.ceil(totalItems / pageSize) : 0;
    const startItem = totalItems ? (currentPage - 1) * pageSize + 1 : 0;
    const endItem = totalItems ? Math.min(currentPage * pageSize, totalItems) : 0;

    return (
        <div className="flex items-center justify-between px-2 py-4">
            <div className="flex items-center space-x-2">
                <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                    Rows per page
                </Typography>
                <Select
                    value={`${pageSize}`}
                    onValueChange={(value) => {
                        onPageSizeChange(Number(value));
                        onPageChange(1); // Reset to first page when changing page size
                    }}
                >
                    <SelectTrigger className="h-8 w-[70px]">
                        <SelectValue placeholder={pageSize} />
                    </SelectTrigger>
                    <SelectContent side="top">
                        {[10, 20, 30, 40, 50].map((size) => (
                            <SelectItem key={size} value={`${size}`}>
                                {size}
                            </SelectItem>
                        ))}
                    </SelectContent>
                </Select>
            </div>

            <div className="flex items-center space-x-6 lg:space-x-8">
                <div className="flex items-center space-x-2">
                    <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                        {totalItems ? (
                            <>
                                {startItem}-{endItem} of {totalItems}
                            </>
                        ) : (
                            "No data"
                        )}
                    </Typography>
                </div>

                <div className="flex items-center space-x-2">
                    <Button
                        variant="secondary-light"
                        className="h-8 w-8 p-0"
                        onClick={() => onPageChange(1)}
                        disabled={currentPage === 1}
                    >
                        <span className="sr-only">Go to first page</span>
                        <ChevronsLeft className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="secondary-light"
                        className="h-8 w-8 p-0"
                        onClick={() => onPageChange(currentPage - 1)}
                        disabled={currentPage === 1}
                    >
                        <span className="sr-only">Go to previous page</span>
                        <ChevronLeft className="h-4 w-4" />
                    </Button>

                    <Typography variant="text" tag="p" className="text-sm font-medium">
                        Page {currentPage} of {totalPages || 1}
                    </Typography>

                    <Button
                        variant="secondary-light"
                        className="h-8 w-8 p-0"
                        onClick={() => onPageChange(currentPage + 1)}
                        disabled={currentPage === totalPages || !totalPages}
                    >
                        <span className="sr-only">Go to next page</span>
                        <ChevronRight className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="secondary-light"
                        className="h-8 w-8 p-0"
                        onClick={() => onPageChange(totalPages)}
                        disabled={currentPage === totalPages || !totalPages}
                    >
                        <span className="sr-only">Go to last page</span>
                        <ChevronsRight className="h-4 w-4" />
                    </Button>
                </div>
            </div>
        </div>
    );
}
