import type { ColumnDef } from "@tanstack/react-table";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";
import type { EventEventParticipantResponse } from "@decm/api";
import { CopyButton } from "@/components/ui/copy-button";

// Helper to render cell content with "(empty)" fallback
const renderCellContent = (
    value: string | null | undefined,
    className?: string,
    emptyText?: string,
) => {
    const isEmpty =
        value === null ||
        value === undefined ||
        value === "undefined" ||
        (typeof value === "string" && value.trim() === "");

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

export function useAttendeeColumns() {
    const { t } = useTranslation();

    const attendeeColumns: ColumnDef<EventEventParticipantResponse>[] = [
        {
            accessorKey: "first_name",
            header: t("events.hostDetails.attendees.firstName"),
            enableSorting: true,
            cell: ({ row }) => {
                const firstName = row.getValue("first_name") as string | null | undefined;
                return renderCellContent(
                    firstName,
                    "font-mono text-xs min-w-[160px]",
                    t("common.empty"),
                );
            },
        },
        {
            accessorKey: "last_name",
            header: t("events.hostDetails.attendees.lastName"),
            enableSorting: true,
            cell: ({ row }) => {
                const lastName = row.getValue("last_name") as string | null | undefined;
                return renderCellContent(lastName, "font-mono text-xs", t("common.empty"));
            },
        },
        {
            accessorKey: "email",
            header: t("events.hostDetails.attendees.email"),
            enableSorting: true,
            cell: ({ row }) => {
                const email = row.getValue("email") as string | null | undefined;
                return renderCellContent(email, "font-mono text-xs", t("common.empty"));
            },
        },
        {
            accessorKey: "phone_number",
            header: t("events.hostDetails.attendees.phoneNumber"),
            enableSorting: false,
            cell: ({ row }) => {
                const phoneNumber = row.getValue("phone_number") as string | null | undefined;
                return renderCellContent(
                    phoneNumber,
                    "font-mono text-xs min-w-[128px]",
                    t("common.empty"),
                );
            },
        },
        {
            accessorKey: "academic_institution",
            header: t("events.hostDetails.attendees.academicInstitution"),
            enableSorting: true,
            cell: ({ row }) => {
                const academicInstitution = row.getValue("academic_institution") as
                    | string
                    | null
                    | undefined;
                return renderCellContent(
                    academicInstitution,
                    "font-mono text-xs",
                    t("common.empty"),
                );
            },
        },
        {
            accessorKey: "wallet_address",
            header: t("events.hostDetails.attendees.walletAddress"),
            enableSorting: false,
            cell: ({ row }) => {
                const address = row.getValue("wallet_address") as string | null | undefined;

                // Handle empty/null/undefined addresses
                if (
                    !address ||
                    address === "undefined" ||
                    (typeof address === "string" && address.trim() === "")
                ) {
                    return (
                        <Typography
                            variant="text"
                            tag="span"
                            className="font-mono text-xs"
                            color="muted-foreground"
                        >
                            <span className="italic">{t("common.empty")}</span>
                        </Typography>
                    );
                }

                // Truncate wallet address for better display
                const truncated = `${address.slice(0, 6)}...${address.slice(-4)}`;
                return (
                    <div className="flex items-center gap-2">
                        <Typography variant="text" tag="span" className="font-mono text-xs">
                            {truncated}
                        </Typography>
                        <CopyButton text={address} iconSize={14} />
                    </div>
                );
            },
        },
    ];

    return attendeeColumns;
}
