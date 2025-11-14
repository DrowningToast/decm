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
import { Trash2 } from "lucide-react";
import ConfirmModal from "../ConfirmModal";
import type { EventIssuer } from "@/services/EventService/EventService";

export interface Issuer {
    id: string;
    name: string;
    email: string;
    organization?: string;
}

export interface SelectedIssuersTableProps {
    selectedIssuers?: EventIssuer[];
    onRemoveIssuer: (issuerId: string) => void;
}

export const SelectedIssuersTable = ({
    selectedIssuers,
    onRemoveIssuer,
}: SelectedIssuersTableProps) => {
    const { t } = useTranslation();

    console.log(selectedIssuers);

    if (!selectedIssuers || selectedIssuers.length === 0) {
        return null;
    }

    return (
        <div className="mt-6">
            <Typography variant="text" tag="p" className="text-sm font-medium mb-3">
                {t("certificateSettings.step1.selectedIssuers")} ({selectedIssuers.length})
            </Typography>
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>{t("certificateSettings.step1.table.name")}</TableHead>
                            <TableHead>{t("certificateSettings.step1.table.email")}</TableHead>
                            <TableHead>
                                {t("certificateSettings.step1.table.organization")}
                            </TableHead>
                            <TableHead className="w-[100px]">
                                {t("certificateSettings.step1.table.actions")}
                            </TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {selectedIssuers.map((issuer) => (
                            <TableRow key={issuer.id}>
                                <TableCell className="font-medium">
                                    {issuer.issuerProfile?.firstName}{" "}
                                    {issuer.issuerProfile?.lastName}
                                </TableCell>
                                <TableCell>{issuer.issuerProfile?.email}</TableCell>
                                <TableCell>
                                    {issuer.issuerProfile?.academicInstitution || "-"}
                                </TableCell>
                                <TableCell>
                                    <ConfirmModal
                                        title={t("certificateSettings.step1.removeIssuer.title")}
                                        cancelText={t("common.cancel")}
                                        confirmText={t("common.remove")}
                                        message={t(
                                            "certificateSettings.step1.removeIssuer.description",
                                        )}
                                        onConfirm={() => onRemoveIssuer(issuer.id ?? "")}
                                        onCancel={() => {}}
                                        destructive
                                    >
                                        <Button type="button" variant="secondary-light" size="sm">
                                            <Trash2 className="h-4 w-4 text-destructive" />
                                        </Button>
                                    </ConfirmModal>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>
        </div>
    );
};
