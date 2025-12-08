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
    isCertificatePublished?: boolean,
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
                const isRevokeDisabled = !eventCertificateId || isCertificatePublished;
                return (
                    <ConfirmModal
                        title="Revoke Certificate"
                        message={
                            isCertificatePublished
                                ? "Cannot revoke certificates after the certificate configuration has been published."
                                : "Are you sure you want to revoke this certificate? This will reset all issuer signatures and require re-approval from all issuers before publishing again."
                        }
                        onConfirm={() => {
                            if (eventCertificateId && !isCertificatePublished) {
                                onClickRevoke(eventCertificateId);
                            }
                        }}
                        onCancel={() => {}}
                        cancelText="Cancel"
                        confirmText={isCertificatePublished ? "OK" : "Revoke"}
                    >
                        <Button
                            size="sm"
                            className="bg-red-400 text-sm text-white"
                            disabled={isRevokeDisabled}
                            title={
                                isCertificatePublished
                                    ? "Cannot revoke - certificate is published"
                                    : "Revoke certificate"
                            }
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
