import { useState } from "react";
import { flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import type {
    ColumnDef,
    SortingState,
    ColumnFiltersState,
    VisibilityState,
    OnChangeFn,
} from "@tanstack/react-table";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { DataTablePagination } from "./data-table-pagination";
import { DataTableToolbar } from "./data-table-toolbar";
import { Typography } from "@/components/typography/typography";
import { ArrowUpDown, ArrowUp, ArrowDown } from "lucide-react";

export interface DataTableProps<TData, TValue> {
    columns: ColumnDef<TData, TValue>[];
    data: TData[];
    // Pagination props
    totalItems: number;
    currentPage: number;
    pageSize: number;
    onPageChange: (page: number) => void;
    onPageSizeChange: (pageSize: number) => void;
    disablePagination?: boolean;
    // Search props
    searchValue: string;
    onSearchChange: (value: string) => void;
    searchPlaceholder?: string;
    // Sorting props
    sorting: SortingState;
    onSortingChange: OnChangeFn<SortingState>;
    // Loading state
    isLoading?: boolean;
}

export function DataTable<TData, TValue>({
    columns,
    data,
    totalItems,
    currentPage,
    pageSize,
    onPageChange,
    onPageSizeChange,
    disablePagination = false,
    searchValue,
    onSearchChange,
    searchPlaceholder,
    sorting,
    onSortingChange,
    isLoading = false,
}: DataTableProps<TData, TValue>) {
    const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
    const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
    const [rowSelection, setRowSelection] = useState({});

    const table = useReactTable({
        data,
        columns,
        state: {
            sorting,
            columnFilters,
            columnVisibility,
            rowSelection,
        },
        enableRowSelection: true,
        onRowSelectionChange: setRowSelection,
        onSortingChange,
        onColumnFiltersChange: setColumnFilters,
        onColumnVisibilityChange: setColumnVisibility,
        getCoreRowModel: getCoreRowModel(),
        manualPagination: true,
        manualSorting: true,
        manualFiltering: true,
        pageCount: Math.ceil(totalItems / pageSize),
    });

    console.log(data);

    return (
        <div className="space-y-4">
            <DataTableToolbar
                searchValue={searchValue}
                onSearchChange={onSearchChange}
                searchPlaceholder={searchPlaceholder}
            />

            <div className="rounded-md border bg-white">
                <Table>
                    <TableHeader>
                        {table.getHeaderGroups().map((headerGroup) => (
                            <TableRow key={headerGroup.id}>
                                {headerGroup.headers.map((header) => {
                                    const canSort = header.column.getCanSort();
                                    const sortDirection = header.column.getIsSorted();

                                    return (
                                        <TableHead
                                            key={header.id}
                                            className={canSort ? "cursor-pointer select-none" : ""}
                                            onClick={
                                                canSort
                                                    ? header.column.getToggleSortingHandler()
                                                    : undefined
                                            }
                                        >
                                            {header.isPlaceholder ? null : (
                                                <div className="flex items-center space-x-2 text-black">
                                                    <span>
                                                        {flexRender(
                                                            header.column.columnDef.header,
                                                            header.getContext(),
                                                        )}
                                                    </span>
                                                    {canSort && (
                                                        <span className="ml-2">
                                                            {sortDirection === "asc" ? (
                                                                <ArrowUp className="h-4 w-4" />
                                                            ) : sortDirection === "desc" ? (
                                                                <ArrowDown className="h-4 w-4" />
                                                            ) : (
                                                                <ArrowUpDown className="h-4 w-4 opacity-50" />
                                                            )}
                                                        </span>
                                                    )}
                                                </div>
                                            )}
                                        </TableHead>
                                    );
                                })}
                            </TableRow>
                        ))}
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            <TableRow>
                                <TableCell colSpan={columns.length} className="h-24 text-center">
                                    <Typography variant="text" tag="p" color="background">
                                        Loading...
                                    </Typography>
                                </TableCell>
                            </TableRow>
                        ) : table.getRowModel().rows?.length ? (
                            table.getRowModel().rows.map((row) => (
                                <TableRow
                                    key={row.id}
                                    data-state={row.getIsSelected() && "selected"}
                                    className="text-black"
                                >
                                    {row.getVisibleCells().map((cell) => (
                                        <TableCell key={cell.id}>
                                            {flexRender(
                                                cell.column.columnDef.cell,
                                                cell.getContext(),
                                            )}
                                        </TableCell>
                                    ))}
                                </TableRow>
                            ))
                        ) : (
                            <TableRow>
                                <TableCell colSpan={columns.length} className="h-24 text-center">
                                    <Typography variant="text" tag="p" color="background-alt">
                                        No results found.
                                    </Typography>
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>

            {!disablePagination && (
                <DataTablePagination
                    totalItems={totalItems}
                    currentPage={currentPage}
                    pageSize={pageSize}
                    onPageChange={onPageChange}
                    onPageSizeChange={onPageSizeChange}
                />
            )}
        </div>
    );
}
