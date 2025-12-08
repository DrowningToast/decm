import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import ConfirmModal from "@/components/ConfirmModal";
import { Typography } from "@/components/typography/typography";
import type { EntityEventCertificate } from "@decm/api";

export interface CertificateRow {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    academicInstitution: string;
    issuedAt: string;
    status: "received" | "pending" | "rejected";
}

export function CertificateColumns(
    onClickRevoke?: (eventCertificateId: string) => void,
): ColumnDef<EntityEventCertificate>[] {
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
            cell: ({ row }) => {
                const issuedAt = row.original.created_at;
                if (!issuedAt) {
                    return <span className="text-muted-foreground">-</span>;
                }
                const date = new Date(issuedAt);
                const formattedDate = date.toLocaleDateString("en-US", {
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                });
                const formattedTime = date.toLocaleTimeString("en-US", {
                    hour: "2-digit",
                    minute: "2-digit",
                });
                return (
                    <span>
                        {formattedDate} {formattedTime}
                    </span>
                );
            },
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
    ];

    if (onClickRevoke) {
        certificateColumns.push({
            id: "actions",
            header: "Action",
            enableSorting: false,
            cell: ({ row }) => {
                const eventCertificateId = row.original.id;
                const isRevokeDisabled = !eventCertificateId;
                return (
                    <ConfirmModal
                        title="Revoke Certificate"
                        message="Are you sure you want to revoke this certificate?"
                        onConfirm={() => {
                            if (eventCertificateId) {
                                onClickRevoke(eventCertificateId);
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
        });
    }

    return certificateColumns as ColumnDef<EntityEventCertificate>[];
}
