import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Typography } from "@/components/typography/typography";

export interface Participant {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    phoneNumber: string;
    academicInstitution: string;
    walletAddress: string;
    status: "confirmed" | "pending" | "rejected";
}

export const participantColumns: ColumnDef<Participant>[] = [
    {
        accessorKey: "first_name",
        header: "First Name",
        enableSorting: true,
    },
    {
        accessorKey: "last_name",
        header: "Last Name",
        enableSorting: true,
    },
    {
        accessorKey: "email",
        header: "Email",
        enableSorting: true,
    },
    {
        accessorKey: "phone_number",
        header: "Phone Number",
        enableSorting: false,
    },
    {
        accessorKey: "academic_institution",
        header: "Academic Institution",
        enableSorting: true,
    },
    // {
    //     accessorKey: "walletAddress",
    //     header: "Wallet Address",
    //     enableSorting: false,
    //     cell: ({ row }) => {
    //         const address = row.getValue("walletAddress") as string;
    //         return (
    //             <Typography variant="text" tag="span" className="font-mono text-xs">
    //                 {address}
    //             </Typography>
    //         );
    //     },
    // },
    // {
    //     accessorKey: "status",
    //     header: "Status",
    //     enableSorting: true,
    //     cell: ({ row }) => {
    //         const status = row.getValue("status") as string;
    //         const statusColors = {
    //             confirmed: "bg-green-100 text-green-800",
    //             pending: "bg-yellow-100 text-yellow-800",
    //             rejected: "bg-red-100 text-red-800",
    //         };
    //         const statusLabels = {
    //             confirmed: "Confirmed",
    //             pending: "Pending",
    //             rejected: "Rejected",
    //         };
    //         return (
    //             <Typography
    //                 variant="text"
    //                 tag="span"
    //                 className={`inline-block px-2 py-1 rounded text-xs ${statusColors[status as keyof typeof statusColors]}`}
    //             >
    //                 {statusLabels[status as keyof typeof statusLabels]}
    //             </Typography>
    //         );
    //     },
    // },
    {
        id: "actions",
        header: "Action",
        enableSorting: false,
        cell: ({ row }) => {
            return (
                <Button
                    size="sm"
                    onClick={() => {
                        // TODO: Implement view participant details
                        console.log("View participant:", row.original.id);
                    }}
                    className="bg-red-400 text-sm text-white"
                >
                    Cancel
                </Button>
            );
        },
    },
];
