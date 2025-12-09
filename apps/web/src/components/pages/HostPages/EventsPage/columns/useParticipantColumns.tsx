import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { useCancelEventInvitation } from "@/hooks/events/useCancelEventInvitation";
import ConfirmModal from "@/components/ConfirmModal";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";

export interface Participant {
    id: string;
    firstName: string | null;
    lastName: string | null;
    email: string | null;
    phoneNumber: string | null;
    academicInstitution: string | null;
    walletAddress: string;
    status: "confirmed" | "pending" | "rejected";
    isAccepted: boolean;
}

// Helper to render cell content with "(empty)" fallback
const renderCellContent = (
    value: string | null | undefined,
    className?: string,
    emptyText?: string,
) => {
    const isEmpty = !value || (typeof value === "string" && value.trim() === "");

    if (isEmpty) {
        return (
            <Typography variant="text" tag="span" className={className} color="muted-foreground">
                <span className="italic">{emptyText || "(empty)"}</span>
            </Typography>
        );
    }

    return (
        <Typography variant="text" tag="span" className={className}>
            {value}
        </Typography>
    );
};

export function useParticipantColumns(eventId: string) {
    const { cancelEventInvitation } = useCancelEventInvitation();
    const { t } = useTranslation();

    const participantColumns: ColumnDef<Participant>[] = [
        {
            accessorKey: "firstName",
            header: t("events.participants.fields.firstName"),
            enableSorting: true,
            cell: ({ row }) => {
                const firstName = row.getValue("firstName") as string | null;
                return renderCellContent(
                    firstName,
                    "font-mono text-xs min-w-[160px]",
                    t("common.empty"),
                );
            },
        },
        {
            accessorKey: "lastName",
            header: t("events.participants.fields.lastName"),
            enableSorting: true,
            cell: ({ row }) => {
                const lastName = row.getValue("lastName") as string | null;
                return renderCellContent(lastName, "font-mono text-xs", t("common.empty"));
            },
        },
        {
            accessorKey: "email",
            header: t("events.participants.fields.email"),
            enableSorting: true,
            cell: ({ row }) => {
                const email = row.getValue("email") as string | null;
                return renderCellContent(email, "font-mono text-xs", t("common.empty"));
            },
        },
        {
            accessorKey: "phoneNumber",
            header: t("events.participants.fields.phoneNumber"),
            enableSorting: false,
            cell: ({ row }) => {
                const phoneNumber = row.getValue("phoneNumber") as string | null;
                return renderCellContent(
                    phoneNumber,
                    "font-mono text-xs min-w-[128px]",
                    t("common.empty"),
                );
            },
        },
        {
            accessorKey: "academicInstitution",
            header: t("events.participants.fields.academicInstitution"),
            enableSorting: true,
            cell: ({ row }) => {
                const academicInstitution = row.getValue("academicInstitution") as string | null;
                return renderCellContent(
                    academicInstitution,
                    "font-mono text-xs",
                    t("common.empty"),
                );
            },
        },
        {
            accessorKey: "isAccepted",
            header: t("events.hostDetails.participants.statusHeader"),
            enableSorting: true,
            cell: ({ row }) => {
                const isAccepted = row.original.isAccepted;
                const statusColor = isAccepted
                    ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                    : "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200";
                const statusLabel = isAccepted
                    ? t("events.hostDetails.participants.statusAccepted")
                    : t("events.hostDetails.participants.statusPending");

                return (
                    <Typography
                        variant="text"
                        tag="span"
                        className={`inline-block px-2 py-1 rounded text-xs font-medium ${statusColor}`}
                    >
                        {statusLabel}
                    </Typography>
                );
            },
        },
        {
            id: "actions",
            header: t("events.hostDetails.participants.table.actions"),
            enableSorting: false,
            cell: ({ row }) => {
                const eventInvitationId = row.original.id;
                const isAccepted = row.original.isAccepted;

                // If invitation is already accepted, show disabled button
                if (isAccepted) {
                    return (
                        <Button
                            size="sm"
                            className="bg-gray-300 text-sm text-gray-500 cursor-not-allowed"
                            disabled
                            title={t("events.hostDetails.participants.alreadyAccepted")}
                        >
                            {t("events.hostDetails.participants.revokeInvitation")}
                        </Button>
                    );
                }

                return (
                    <ConfirmModal
                        title={t("events.hostDetails.participants.revokeInvitationTitle")}
                        message={t("events.hostDetails.participants.revokeInvitationMessage")}
                        onConfirm={() => cancelEventInvitation({ eventInvitationId, eventId })}
                        onCancel={() => {}}
                        cancelText={t("common.cancel")}
                        confirmText={t("common.confirm")}
                    >
                        <Button size="sm" className="bg-red-400 text-sm text-white">
                            {t("events.hostDetails.participants.revokeInvitation")}
                        </Button>
                    </ConfirmModal>
                );
            },
        },
    ];

    return participantColumns;
}
