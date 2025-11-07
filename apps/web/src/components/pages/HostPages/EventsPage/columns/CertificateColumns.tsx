import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import ConfirmModal from "@/components/ConfirmModal";

export interface CertificateRow {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    academicInstitution: string;
    issuedAt: string;
    status: "received" | "pending" | "rejected";
}

export function CertificateColumns() {
    // const { cancelEventInvitation } = useCancelEventInvitation();

    const certificateColumns: ColumnDef<CertificateRow>[] = [
        {
            accessorKey: "firstName",
            header: "First Name",
            enableSorting: true,
        },
        {
            accessorKey: "lastName",
            header: "Last Name",
            enableSorting: true,
        },
        {
            accessorKey: "email",
            header: "Email",
            enableSorting: true,
        },
        {
            accessorKey: "academicInstitution",
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
                // const eventCertificateId = row.original.id;

                return (
                    <ConfirmModal
                        title="Revoke Certificate"
                        message="Are you sure you want to revoke this certificate?"
                        onConfirm={() => {}}
                        onCancel={() => {}}
                        cancelText="Cancel"
                        confirmText="Revoke"
                    >
                        <Button size="sm" className="bg-red-400 text-sm text-white">
                            Revoke
                        </Button>
                    </ConfirmModal>
                );
            },
        },
    ];

    return certificateColumns;
}
