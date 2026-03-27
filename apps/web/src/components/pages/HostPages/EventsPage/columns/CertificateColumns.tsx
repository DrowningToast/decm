import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import ConfirmModal from "@/components/ConfirmModal";
import { Typography } from "@/components/typography/typography";
import type { EntityEventCertificate } from "@decm/api";
import { useTranslation } from "react-i18next";

export interface CertificateRow {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    academicInstitution: string;
    issuedAt: string;
    status: "received" | "pending" | "rejected";
}

export function useCertificateColumns(
    onClickRevoke?: (eventCertificateId: string) => void,
    isCertificatePublished?: boolean,
): ColumnDef<EntityEventCertificate>[] {
    const { t, i18n } = useTranslation();
    const certificateColumns: ColumnDef<EntityEventCertificate>[] = [
        {
            accessorKey: "firstName",
            header: t("events.participants.fields.firstName"),
            enableSorting: true,
        },
        {
            accessorKey: "lastName",
            header: t("events.participants.fields.lastName"),
            enableSorting: true,
        },
        {
            accessorKey: "email",
            header: t("events.participants.fields.email"),
            enableSorting: true,
        },
        {
            accessorKey: "academicInstitution",
            header: t("events.participants.fields.academicInstitution"),
            enableSorting: true,
        },
        {
            accessorKey: "receiver_wallet_address",
            header: t("events.participants.fields.walletAddress"),
            enableSorting: true,
            cell: ({ row }) => {
                const walletAddress = row.original.receiver_wallet_address;
                if (!walletAddress) {
                    return <span className="text-muted-foreground">-</span>;
                }
                const shortened = `${walletAddress.slice(0, 6)}...${walletAddress.slice(-4)}`;
                return (
                    <span title={walletAddress} className="font-mono text-xs">
                        {shortened}
                    </span>
                );
            },
        },
        {
            accessorKey: "issuedAt",
            header: t("events.hostDetails.certificates.table.issuedAt"),
            enableSorting: true,
            cell: ({ row }) => {
                const issuedAt = row.original.created_at;
                if (!issuedAt) {
                    return <span className="text-muted-foreground">-</span>;
                }
                const date = new Date(issuedAt);
                // Use i18n language for date formatting (e.g., 'en' or 'th')
                const locale = i18n.language === "th" ? "th-TH" : "en-US";
                const formattedDate = date.toLocaleDateString(locale, {
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                });
                const formattedTime = date.toLocaleTimeString(locale, {
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
            header: t("events.hostDetails.certificates.table.status"),
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
                            {t("events.hostDetails.certificates.table.statusCollected")}
                        </Typography>
                    );
                } else {
                    return (
                        <Typography
                            variant="text"
                            tag="span"
                            className="text-red-500 inline-block px-2 py-1 rounded text-xs bg-red-100"
                        >
                            {t("events.hostDetails.certificates.table.statusNotCollected")}
                        </Typography>
                    );
                }
            },
        },
    ];

    if (onClickRevoke) {
        certificateColumns.push({
            id: "actions",
            header: t("events.hostDetails.certificates.table.action"),
            enableSorting: false,
            cell: ({ row }) => {
                const eventCertificateId = row.original.id;
                const isRevokeDisabled = !eventCertificateId || isCertificatePublished;
                return (
                    <ConfirmModal
                        title={t("events.hostDetails.certificates.table.revokeTitle")}
                        message={
                            isCertificatePublished
                                ? t("events.hostDetails.certificates.table.revokeMessagePublished")
                                : t("events.hostDetails.certificates.table.revokeMessage")
                        }
                        onConfirm={() => {
                            if (eventCertificateId && !isCertificatePublished) {
                                onClickRevoke(eventCertificateId);
                            }
                        }}
                        onCancel={() => {}}
                        cancelText={t("common.cancel")}
                        confirmText={
                            isCertificatePublished
                                ? t("common.confirm")
                                : t("events.hostDetails.certificates.table.revokeButton")
                        }
                    >
                        <Button
                            size="sm"
                            className="bg-red-400 text-sm text-white"
                            disabled={isRevokeDisabled}
                            title={
                                isCertificatePublished
                                    ? t("events.hostDetails.certificates.table.revokeDisabledTitle")
                                    : t("events.hostDetails.certificates.table.revokeTitle")
                            }
                        >
                            {t("events.hostDetails.certificates.table.revokeButton")}
                        </Button>
                    </ConfirmModal>
                );
            },
        });
    }

    return certificateColumns as ColumnDef<EntityEventCertificate>[];
}
