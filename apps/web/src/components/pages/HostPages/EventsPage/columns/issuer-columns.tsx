import type { ColumnDef } from "@tanstack/react-table";
import type { EventIssuer } from "@/services/EventService/EventService";
import { useTranslation } from "react-i18next";

export function useIssuerColumns(): ColumnDef<EventIssuer>[] {
    const { t, i18n } = useTranslation();
    return [
        {
            accessorKey: "name",
            header: t("certificateSettings.step1.table.name"),
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
            header: t("certificateSettings.step1.table.email"),
            enableSorting: true,
            cell: ({ row }) => {
                const email = row.original.issuerProfile?.email || t("common.notAvailable");
                return <span>{email}</span>;
            },
        },
        {
            accessorKey: "googleOAuthEmail",
            header: t("certificateSettings.step1.table.googleOAuthEmail"),
            enableSorting: true,
            cell: ({ row }) => {
                const googleOAuthEmail =
                    row.original.issuerProfile?.googleConnectorRef || t("common.notAvailable");
                return <span>{googleOAuthEmail}</span>;
            },
        },
        {
            accessorKey: "organization",
            header: t("certificateSettings.step1.table.organization"),
            enableSorting: true,
            cell: ({ row }) => {
                const organization =
                    row.original.issuerProfile?.academicInstitution || t("common.notAvailable");
                return <span>{organization}</span>;
            },
        },
        {
            accessorKey: "signingStatus",
            header: t("events.hostDetails.certificates.table.signingStatus"),
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
            header: t("events.hostDetails.certificates.table.signedAt"),
            enableSorting: true,
            cell: ({ row }) => {
                const issuer = row.original;
                if (issuer.isSigned && issuer.updatedAt) {
                    const date = new Date(issuer.updatedAt);
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
                }
                return <span className="text-muted-foreground">-</span>;
            },
        },
    ];
}
