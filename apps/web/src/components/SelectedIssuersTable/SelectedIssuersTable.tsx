import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { Trash2 } from "lucide-react";
import ConfirmModal from "../ConfirmModal";
import type { EventIssuer } from "@/services/EventService/EventService";

export interface Issuer {
    id: string;
    name: string;
    email: string;
    organization?: string;
}

// Helper to render cell content with "(empty)" fallback
const renderCellContent = (value: string | null | undefined, className?: string) => {
    const isEmpty =
        value === null ||
        value === undefined ||
        value === "undefined" ||
        (typeof value === "string" && value.trim() === "");

    if (isEmpty) {
        return (
            <Typography variant="text" tag="span" className={className} color="muted-foreground">
                <span className="italic">(empty)</span>
            </Typography>
        );
    }

    return (
        <Typography variant="text" tag="span" className={className}>
            {value}
        </Typography>
    );
};

export interface SelectedIssuersTableProps {
    selectedIssuers?: EventIssuer[];
    onRemoveIssuer: (issuerId: string) => void;
    title?: string;
    isUnsaved?: boolean;
}

export const SelectedIssuersTable = ({
    selectedIssuers,
    onRemoveIssuer,
    title,
    isUnsaved = false,
}: SelectedIssuersTableProps) => {
    const { t } = useTranslation();

    if (!selectedIssuers || selectedIssuers.length === 0) {
        return null;
    }

    const displayTitle = title || t("certificateSettings.step1.selectedIssuers");
    const isOnlyOneIssuer = selectedIssuers.length === 1;

    return (
        <div className="mt-6">
            <Typography variant="text" tag="p" className="text-sm font-medium mb-3">
                {displayTitle} ({selectedIssuers.length})
            </Typography>
            <div
                className={`rounded-md border ${
                    isUnsaved ? "border-amber-300 bg-amber-50/30" : "border-border"
                }`}
            >
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>{t("certificateSettings.step1.table.name")}</TableHead>
                            <TableHead>
                                {t("certificateSettings.step1.table.contactEmail")}
                            </TableHead>
                            <TableHead>
                                {t("certificateSettings.step1.table.googleOAuthEmail")}
                            </TableHead>
                            <TableHead>
                                {t("certificateSettings.step1.table.phoneNumber")}
                            </TableHead>
                            <TableHead>
                                {t("certificateSettings.step1.table.organization")}
                            </TableHead>
                            <TableHead className="w-[100px]">
                                {t("certificateSettings.step1.table.actions")}
                            </TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {selectedIssuers.map((issuer) => {
                            const firstName = issuer.issuerProfile?.firstName || "";
                            const lastName = issuer.issuerProfile?.lastName || "";
                            const fullName = [firstName, lastName].filter(Boolean).join(" ").trim();

                            return (
                                <TableRow key={issuer.id}>
                                    <TableCell className="font-medium">
                                        {renderCellContent(fullName || undefined)}
                                    </TableCell>
                                    <TableCell>
                                        {renderCellContent(issuer.issuerProfile?.email)}
                                    </TableCell>
                                    <TableCell>
                                        {renderCellContent(
                                            issuer.issuerProfile?.googleConnectorRef,
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {renderCellContent(issuer.issuerProfile?.phoneNumber)}
                                    </TableCell>
                                    <TableCell>
                                        {renderCellContent(
                                            issuer.issuerProfile?.academicInstitution,
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {isOnlyOneIssuer ? (
                                            <Tooltip>
                                                <TooltipTrigger asChild>
                                                    <span>
                                                        <Button
                                                            type="button"
                                                            variant="secondary-light"
                                                            size="sm"
                                                            disabled={true}
                                                        >
                                                            <Trash2 className="h-4 w-4 text-destructive" />
                                                        </Button>
                                                    </span>
                                                </TooltipTrigger>
                                                <TooltipContent>
                                                    <p className="text-white font-medium">
                                                        {t(
                                                            "certificateSettings.step1.removeIssuer.minimumIssuerRequired",
                                                        )}
                                                    </p>
                                                </TooltipContent>
                                            </Tooltip>
                                        ) : (
                                            <ConfirmModal
                                                title={t(
                                                    "certificateSettings.step1.removeIssuer.title",
                                                )}
                                                cancelText={t("common.cancel")}
                                                confirmText={t("common.remove")}
                                                message={t(
                                                    "certificateSettings.step1.removeIssuer.description",
                                                )}
                                                onConfirm={() => onRemoveIssuer(issuer.id ?? "")}
                                                onCancel={() => {}}
                                                destructive
                                            >
                                                <Button
                                                    type="button"
                                                    variant="secondary-light"
                                                    size="sm"
                                                >
                                                    <Trash2 className="h-4 w-4 text-destructive" />
                                                </Button>
                                            </ConfirmModal>
                                        )}
                                    </TableCell>
                                </TableRow>
                            );
                        })}
                    </TableBody>
                </Table>
            </div>
        </div>
    );
};
