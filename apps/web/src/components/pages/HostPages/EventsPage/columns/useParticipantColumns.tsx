import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { useCancelEventInvitation } from "@/hooks/events/useCancelEventInvitation";
import ConfirmModal from "@/components/ConfirmModal";
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

export function useParticipantColumns() {
    const { cancelEventInvitation } = useCancelEventInvitation();

    const participantColumns: ColumnDef<Participant>[] = [
        {
            accessorKey: "firstName",
            header: "First Name",
            enableSorting: true,
            cell: ({ row }) => {
                const firstName = row.getValue("firstName") as string;
                return (
                    <Typography
                        variant="text"
                        tag="span"
                        className="font-mono text-xs min-w-[160px]"
                    >
                        {firstName}
                    </Typography>
                );
            },
        },
        {
            accessorKey: "lastName",
            header: "Last Name",
            enableSorting: true,
            cell: ({ row }) => {
                const lastName = row.getValue("lastName") as string;
                return (
                    <Typography variant="text" tag="span" className="font-mono text-xs">
                        {lastName}
                    </Typography>
                );
            },
        },
        {
            accessorKey: "email",
            header: "Email",
            enableSorting: true,
            cell: ({ row }) => {
                const email = row.getValue("email") as string;
                return (
                    <Typography variant="text" tag="span" className="font-mono text-xs">
                        {email}
                    </Typography>
                );
            },
        },
        {
            accessorKey: "phoneNumber",
            header: "Phone Number",
            enableSorting: false,
            cell: ({ row }) => {
                const phoneNumber = row.getValue("phoneNumber") as string;
                return (
                    <Typography
                        variant="text"
                        tag="span"
                        className="font-mono text-xs min-w-[128px]"
                    >
                        {phoneNumber}
                    </Typography>
                );
            },
        },
        {
            accessorKey: "academicInstitution",
            header: "Academic Institution",
            enableSorting: true,

            cell: ({ row }) => {
                const academicInstitution = row.getValue("academicInstitution") as string;
                return (
                    <Typography variant="text" tag="span" className="font-mono text-xs">
                        {academicInstitution}
                    </Typography>
                );
            },
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
                const eventInvitationId = row.original.id;

                return (
                    <ConfirmModal
                        title="Cancel Invitation"
                        message="Are you sure you want to cancel this invitation?"
                        onConfirm={() => cancelEventInvitation(eventInvitationId)}
                        onCancel={() => {}}
                        cancelText="Cancel"
                        confirmText="Confirm"
                    >
                        <Button size="sm" className="bg-red-400 text-sm text-white">
                            Cancel
                        </Button>
                    </ConfirmModal>
                );
            },
        },
    ];

    return participantColumns;
}
