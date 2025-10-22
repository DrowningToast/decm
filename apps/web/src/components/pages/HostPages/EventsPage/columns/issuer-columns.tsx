import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { ExternalLinkIcon } from "lucide-react";
import type { GetEventIssuersByEventIdData } from "@decm/api";

export type Issuer = GetEventIssuersByEventIdData[0];

export const issuerColumns: ColumnDef<Issuer>[] = [
    {
        accessorKey: "name",
        header: "Name",
        enableSorting: true,
        cell: ({ row }) => {
            const issuer = row.original;
            const firstName = issuer.issuer_profile?.first_name || "";
            const lastName = issuer.issuer_profile?.last_name || "";
            const fullName = firstName || lastName ? `${firstName} ${lastName}`.trim() : "N/A";
            return <span>{fullName}</span>;
        },
    },
    {
        accessorKey: "email",
        header: "Email",
        enableSorting: true,
        cell: ({ row }) => {
            const email = row.original.issuer_profile?.email || "N/A";
            return <span>{email}</span>;
        },
    },
    {
        accessorKey: "organization",
        header: "Organization",
        enableSorting: true,
        cell: ({ row }) => {
            const organization = row.original.issuer_profile?.academic_institution || "N/A";
            return <span>{organization}</span>;
        },
    },
    {
        accessorKey: "signingStatus",
        header: "Signing Status",
        enableSorting: true,
        cell: ({ row }) => {
            const isSigned = row.original.is_signed;
            const statusColors = {
                0: "bg-yellow-100 text-yellow-800",
                1: "bg-green-100 text-green-800",
            };
            const statusLabels = {
                0: "Not Signed",
                1: "Signed",
            };
            return (
                <span
                    className={`inline-block px-2 py-1 rounded text-xs ${
                        statusColors[isSigned as keyof typeof statusColors]
                    }`}
                >
                    {statusLabels[isSigned as keyof typeof statusLabels]}
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
            if (issuer.is_signed === 1 && issuer.updated_at) {
                const date = new Date(issuer.updated_at);
                return (
                    <span>
                        {date.toLocaleDateString()} {date.toLocaleTimeString()}
                    </span>
                );
            }
            return <span className="text-muted-foreground">-</span>;
        },
    },
    {
        id: "transaction",
        header: "Transaction",
        enableSorting: false,
        cell: ({ row }) => {
            const issuer = row.original;
            if (issuer.is_signed === 1 && issuer.signature) {
                // Extract transaction hash from signature or use signature directly
                const txHash = issuer.signature.startsWith("0x")
                    ? issuer.signature
                    : `0x${issuer.signature}`;
                const etherscanUrl = `https://sepolia.etherscan.io/tx/${txHash}`;

                return (
                    <Button
                        variant="secondary-light"
                        size="sm"
                        onClick={() => {
                            window.open(etherscanUrl, "_blank");
                        }}
                    >
                        <ExternalLinkIcon size={14} className="mr-1" />
                        View
                    </Button>
                );
            }
            return <span className="text-muted-foreground">-</span>;
        },
    },
];
