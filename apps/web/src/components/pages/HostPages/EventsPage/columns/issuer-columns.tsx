import type { ColumnDef } from "@tanstack/react-table";
import type { EventIssuer } from "@/services/EventService/EventService";
import { t } from "i18next";

export const issuerColumns: ColumnDef<EventIssuer>[] = [
    {
        accessorKey: "name",
        header: "Name",
        enableSorting: true,
        cell: ({ row }) => {
            const issuer = row.original;
            const firstName = issuer.issuerProfile?.firstName || "";
            const lastName = issuer.issuerProfile?.lastName || "";
            const fullName = firstName || lastName ? `${firstName} ${lastName}`.trim() : "N/A";
            return <span>{fullName}</span>;
        },
    },
    {
        accessorKey: "email",
        header: "Email",
        enableSorting: true,
        cell: ({ row }) => {
            const email = row.original.issuerProfile?.email || "N/A";
            return <span>{email}</span>;
        },
    },
    {
        accessorKey: "googleOAuthEmail",
        header: t("certificateSettings.step1.table.googleOAuthEmail"),
        enableSorting: true,
        cell: ({ row }) => {
            const googleOAuthEmail = row.original.issuerProfile?.googleConnectorRef || "N/A";
            return <span>{googleOAuthEmail}</span>;
        },
    },
    {
        accessorKey: "organization",
        header: "Organization",
        enableSorting: true,
        cell: ({ row }) => {
            const organization = row.original.issuerProfile?.academicInstitution || "N/A";
            return <span>{organization}</span>;
        },
    },
    {
        accessorKey: "signingStatus",
        header: "Signing Status",
        enableSorting: true,
        cell: ({ row }) => {
            const isSigned = row.original.isSigned;
            const statusColors = {
                false: "bg-yellow-100 text-yellow-800",
                true: "bg-green-100 text-green-800",
            };
            const statusLabels = {
                false: t("issuer.sign.status.waiting"),
                true: t("issuer.sign.status.signed"),
            };
            return (
                <span
                    className={`inline-block px-2 py-1 rounded text-xs ${isSigned ? statusColors.true : statusColors.false}`}
                >
                    {isSigned ? statusLabels.true : statusLabels.false}
                </span>
            );
        },
    },
    {
        accessorKey: "signedAt",
        header: "Signed At",
        enableSorting: true,
        cell: ({ row }) => {
            const issuer = row.original;
            if (issuer.isSigned && issuer.updatedAt) {
                const date = new Date(issuer.updatedAt);
                return (
                    <span>
                        {date.toLocaleDateString()} {date.toLocaleTimeString()}
                    </span>
                );
            }
            return <span className="text-muted-foreground">-</span>;
        },
    },
];
