# Data Table Components

## Overview

Reusable TanStack Table (React Table) implementation with built-in pagination, search, sorting, and server-side data fetching support.

## Installation

```bash
pnpm add @tanstack/react-table
```

## Components

### 1. DataTable

Main table component that renders the table with all features.

### 2. DataTablePagination

Pagination controls component with page size selector and navigation buttons.

### 3. DataTableToolbar

Toolbar component with search input and filter reset button.

### 4. useDataTable Hook

Custom hook for managing table state and API calls with automatic debouncing.

## Features

✅ **Server-side Pagination** - Queries backend when page changes  
✅ **Search with Debouncing** - 500ms debounce on search input  
✅ **Sortable Columns** - Click headers to sort ASC/DESC  
✅ **Loading States** - Shows loading indicator during data fetch  
✅ **Flexible Page Sizes** - 10, 20, 30, 40, 50 rows per page  
✅ **Type-safe** - Full TypeScript support  
✅ **Reusable** - Easy to implement across different tables

## Usage

### Step 1: Define Column Definitions

```tsx
import { ColumnDef } from "@tanstack/react-table";

export interface YourDataType {
    id: string;
    name: string;
    email: string;
    status: string;
}

export const yourColumns: ColumnDef<YourDataType>[] = [
    {
        accessorKey: "name",
        header: "Name",
        enableSorting: true,
    },
    {
        accessorKey: "email",
        header: "Email",
        enableSorting: true,
    },
    {
        accessorKey: "status",
        header: "Status",
        enableSorting: true,
        cell: ({ row }) => {
            // Custom cell rendering
            return <span>{row.getValue("status")}</span>;
        },
    },
    {
        id: "actions",
        header: "Actions",
        enableSorting: false,
        cell: ({ row }) => {
            return <Button onClick={() => console.log(row.original.id)}>View</Button>;
        },
    },
];
```

### Step 2: Create API Fetch Function

```tsx
const fetchYourData = async ({
    page,
    pageSize,
    search,
    sortBy,
    sortOrder,
}: {
    page: number;
    pageSize: number;
    search: string;
    sortBy?: string;
    sortOrder?: "asc" | "desc";
}): Promise<{ data: YourDataType[]; total: number }> => {
    // Replace with your actual API call
    const response = await fetch(
        `/api/your-endpoint?page=${page}&pageSize=${pageSize}&search=${search}&sortBy=${sortBy}&sortOrder=${sortOrder}`,
    );
    const result = await response.json();

    return {
        data: result.items,
        total: result.total,
    };
};
```

### Step 3: Use in Your Component

```tsx
import { DataTable } from "@/components/ui/data-table";
import { useDataTable } from "@/hooks/use-data-table";
import { yourColumns, YourDataType } from "./columns/your-columns";

export function YourPage() {
    const table = useDataTable<YourDataType>({
        fetchData: fetchYourData,
        initialPageSize: 10,
    });

    return (
        <div>
            <DataTable
                columns={yourColumns}
                data={table.data}
                totalItems={table.totalItems}
                currentPage={table.currentPage}
                pageSize={table.pageSize}
                onPageChange={table.setCurrentPage}
                onPageSizeChange={table.setPageSize}
                searchValue={table.searchValue}
                onSearchChange={table.setSearchValue}
                searchPlaceholder="Search your data..."
                sorting={table.sorting}
                onSortingChange={table.setSorting}
                isLoading={table.isLoading}
            />
        </div>
    );
}
```

## API Reference

### DataTable Props

| Prop                | Type                         | Description                            |
| ------------------- | ---------------------------- | -------------------------------------- |
| `columns`           | `ColumnDef<TData, TValue>[]` | Column definitions                     |
| `data`              | `TData[]`                    | Table data array                       |
| `totalItems`        | `number`                     | Total number of items (for pagination) |
| `currentPage`       | `number`                     | Current page number (1-indexed)        |
| `pageSize`          | `number`                     | Number of items per page               |
| `onPageChange`      | `(page: number) => void`     | Callback when page changes             |
| `onPageSizeChange`  | `(size: number) => void`     | Callback when page size changes        |
| `searchValue`       | `string`                     | Current search query                   |
| `onSearchChange`    | `(value: string) => void`    | Callback when search changes           |
| `searchPlaceholder` | `string?`                    | Placeholder text for search input      |
| `sorting`           | `SortingState`               | Current sorting state                  |
| `onSortingChange`   | `OnChangeFn<SortingState>`   | Callback when sorting changes          |
| `isLoading`         | `boolean?`                   | Loading state                          |

### useDataTable Hook

```tsx
const table = useDataTable<TData>({
    fetchData: (params) => Promise<{ data: TData[]; total: number }>,
    initialPageSize?: number, // Default: 10
});
```

**Returns:**

- `data` - Current page data
- `totalItems` - Total number of items
- `currentPage` - Current page number
- `pageSize` - Current page size
- `searchValue` - Current search value
- `sorting` - Current sorting state
- `isLoading` - Loading state
- `error` - Error object (if any)
- `setCurrentPage` - Function to change page
- `setPageSize` - Function to change page size
- `setSearchValue` - Function to change search
- `setSorting` - Function to change sorting
- `refetch` - Function to manually refetch data

## Column Definition Options

```tsx
{
    accessorKey: "fieldName",      // Data field to access
    header: "Header Text",          // Column header text
    enableSorting: true,            // Enable sorting for this column
    cell: ({ row }) => {            // Custom cell renderer
        return <CustomComponent value={row.getValue("fieldName")} />;
    },
}
```

## Features in Detail

### Sorting

- Click column header to toggle sort direction
- Visual indicators: ↑ (ASC), ↓ (DESC), ↕ (unsorted)
- Only enabled columns are sortable
- Server-side sorting via API

### Search

- 500ms debounce to prevent excessive API calls
- Searches across all searchable fields
- Automatically resets to page 1 when search changes
- Reset button appears when search is active

### Pagination

- First, Previous, Next, Last page buttons
- Shows current page and total pages
- Shows item range (e.g., "1-10 of 45")
- Page size selector (10, 20, 30, 40, 50)
- Automatically adjusts when data changes

### Loading States

- Shows "Loading..." when fetching data
- Disables interactive elements during load
- Smooth transition between states

## Example: Real API Integration

```tsx
// api/participants.ts
export async function fetchParticipants({
    page,
    pageSize,
    search,
    sortBy,
    sortOrder,
}: {
    page: number;
    pageSize: number;
    search: string;
    sortBy?: string;
    sortOrder?: "asc" | "desc";
}) {
    const params = new URLSearchParams({
        page: page.toString(),
        pageSize: pageSize.toString(),
        ...(search && { search }),
        ...(sortBy && { sortBy }),
        ...(sortOrder && { sortOrder }),
    });

    const response = await fetch(`/api/v1/events/${eventId}/participants?${params}`);

    if (!response.ok) {
        throw new Error("Failed to fetch participants");
    }

    const data = await response.json();

    return {
        data: data.items,
        total: data.total,
    };
}
```

## Customization

### Custom Cell Rendering

```tsx
{
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
        const status = row.getValue("status") as string;
        return (
            <span className={`badge badge-${status}`}>
                {status.toUpperCase()}
            </span>
        );
    },
}
```

### Custom Actions Column

```tsx
{
    id: "actions",
    header: "Actions",
    enableSorting: false,
    cell: ({ row }) => {
        return (
            <div className="flex gap-2">
                <Button onClick={() => handleEdit(row.original)}>
                    Edit
                </Button>
                <Button onClick={() => handleDelete(row.original.id)}>
                    Delete
                </Button>
            </div>
        );
    },
}
```

## Best Practices

1. **Always use server-side operations** for large datasets
2. **Memoize column definitions** if they don't change
3. **Handle errors gracefully** with try-catch in fetch functions
4. **Use TypeScript** for type safety across columns and data
5. **Keep search debounce** at 500ms for good UX
6. **Provide meaningful search placeholders** for context

## Related Files

- Components: `apps/web/src/components/ui/data-table/`
- Hook: `apps/web/src/hooks/use-data-table.ts`
- Example: `apps/web/src/components/pages/HostPages/EventsPage/HostEventDetailsPage.tsx`
- Columns: `apps/web/src/components/pages/HostPages/EventsPage/columns/participant-columns.tsx`

## Dependencies

- `@tanstack/react-table` - Table core functionality
- `lucide-react` - Icons for sorting and pagination
- Built-in UI components: Button, Input, Select, Table

---

**Created**: October 14, 2025  
**Status**: ✅ Production Ready
