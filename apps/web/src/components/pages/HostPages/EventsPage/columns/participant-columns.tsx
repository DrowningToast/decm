import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";

export interface Participant {
    id: string;
    name: string;
    email: string;
    phoneNumber: string;
    walletAddress: string;
    status: "confirmed" | "pending" | "rejected";
}

export const participantColumns: ColumnDef<Participant>[] = [
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
        accessorKey: "phoneNumber",
        header: "Phone Number",
        enableSorting: false,
    },
    {
        accessorKey: "walletAddress",
        header: "Wallet Address",
        enableSorting: false,
        cell: ({ row }) => {
            const address = row.getValue("walletAddress") as string;
            return <span className="font-mono text-xs">{address}</span>;
        },
    },
    {
        accessorKey: "status",
        header: "Status",
        enableSorting: true,
        cell: ({ row }) => {
            const status = row.getValue("status") as string;
            const statusColors = {
                confirmed: "bg-green-100 text-green-800",
                pending: "bg-yellow-100 text-yellow-800",
                rejected: "bg-red-100 text-red-800",
            };
            const statusLabels = {
                confirmed: "Confirmed",
                pending: "Pending",
                rejected: "Rejected",
            };
            return (
                <span
                    className={`inline-block px-2 py-1 rounded text-xs ${
                        statusColors[status as keyof typeof statusColors]
                    }`}
                >
                    {statusLabels[status as keyof typeof statusLabels]}
                </span>
            );
        },
    },
    {
        id: "actions",
        header: "Action",
        enableSorting: false,
        cell: ({ row }) => {
            return (
                <Button
                    variant="secondary-light"
                    size="sm"
                    onClick={() => {
                        // TODO: Implement view participant details
                        console.log("View participant:", row.original.id);
                    }}
                >
                    View
                </Button>
            );
        },
    },
];
