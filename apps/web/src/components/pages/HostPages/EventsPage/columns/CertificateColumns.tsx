import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import ConfirmModal from "@/components/ConfirmModal";
import { Typography } from "@/components/typography/typography";
import type { EntityEventCertificate } from "@decm/api";
import { useRevokeEventCertificate } from "@/hooks/events/useRevokeEventCertificate";

export interface CertificateRow {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    academicInstitution: string;
    issuedAt: string;
    status: "received" | "pending" | "rejected";
}

export function CertificateColumns(eventId: string): ColumnDef<EntityEventCertificate>[] {
    // const { cancelEventInvitation } = useCancelEventInvitation();
    const { revokeEventCertificate } = useRevokeEventCertificate();

    const certificateColumns: ColumnDef<EntityEventCertificate>[] = [
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
        {
            accessorKey: "issuedAt",
            header: "Issued At",
            enableSorting: true,
        },
        {
            accessorKey: "status",
            header: "Status",
            enableSorting: true,
            cell: ({ row }) => {
                const isReceived = row.original.receiver_credential_id;

                if (isReceived) {
                    return (
                        <Typography
                            variant="text"
                            tag="span"
                            className="text-green-500 inline-block px-2 py-1 rounded text-xs bg-green-100"
                        >
                            Collected
                        </Typography>
                    );
                } else {
                    return (
                        <Typography
                            variant="text"
                            tag="span"
                            className="text-red-500 inline-block px-2 py-1 rounded text-xs bg-red-100"
                        >
                            Not Collect
                        </Typography>
                    );
                }
            },
        },
        {
            id: "actions",
            header: "Action",
            enableSorting: false,
            cell: ({ row }) => {
                const eventCertificateId = row.original.id;
                console.log(eventCertificateId);

                // Disable revoke button if certificate ID is not available
                const isRevokeDisabled = !eventCertificateId;

                return (
                    <ConfirmModal
                        title="Revoke Certificate"
                        message="Are you sure you want to revoke this certificate?"
                        onConfirm={() => {
                            if (eventCertificateId) {
                                revokeEventCertificate({
                                    certificateIds: [eventCertificateId],
                                    eventId,
                                });
                            }
                        }}
                        onCancel={() => {}}
                        cancelText="Cancel"
                        confirmText="Revoke"
                    >
                        <Button
                            size="sm"
                            className="bg-red-400 text-sm text-white"
                            disabled={isRevokeDisabled}
                        >
                            Revoke
                        </Button>
                    </ConfirmModal>
                );
            },
        },
    ];

    return certificateColumns as ColumnDef<EntityEventCertificate>[];
}
