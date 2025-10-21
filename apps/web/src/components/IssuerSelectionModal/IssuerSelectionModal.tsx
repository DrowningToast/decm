import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Checkbox } from "@/components/ui/checkbox";

// Type definitions
export interface Issuer {
    id: string;
    name: string;
    email: string;
    organization?: string;
}

export interface IssuerSelectionModalProps {
    isOpen: boolean;
    onClose: () => void;
    onConfirm: (selectedIssuers: Issuer[]) => void;
    searchResults: Issuer[];
    initialSelectedIds?: Set<string>;
    isLoading?: boolean;
    searchQuery?: string;
}

export const IssuerSelectionModal = ({
    isOpen,
    onClose,
    onConfirm,
    searchResults,
    initialSelectedIds = new Set(),
    isLoading = false,
    searchQuery = "",
}: IssuerSelectionModalProps) => {
    const { t } = useTranslation();
    const [tempSelectedIds, setTempSelectedIds] = useState<Set<string>>(new Set());

    // Initialize selection when modal opens or initial selection changes
    useEffect(() => {
        if (isOpen) {
            setTempSelectedIds(new Set(initialSelectedIds));
        }
    }, [isOpen, initialSelectedIds]);

    const handleToggleIssuerSelection = (issuerId: string) => {
        const newSelection = new Set(tempSelectedIds);
        if (newSelection.has(issuerId)) {
            newSelection.delete(issuerId);
        } else {
            newSelection.add(issuerId);
        }
        setTempSelectedIds(newSelection);
    };

    const handleConfirm = () => {
        const selectedIssuers = searchResults.filter((issuer) => tempSelectedIds.has(issuer.id));
        onConfirm(selectedIssuers);
        onClose();
    };

    const handleCancel = () => {
        setTempSelectedIds(new Set(initialSelectedIds));
        onClose();
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>
                        <Typography variant="header" tag="h2" className="text-xl font-bold">
                            {t("certificateSettings.step1.modal.title")}
                        </Typography>
                    </DialogTitle>
                </DialogHeader>

                <div className="mt-4">
                    {isLoading ? (
                        <Typography
                            variant="text"
                            tag="p"
                            className="text-center text-muted-foreground py-8"
                        >
                            {t("common.loading")}
                        </Typography>
                    ) : searchResults.length === 0 ? (
                        <Typography
                            variant="text"
                            tag="p"
                            className="text-center text-muted-foreground py-8"
                        >
                            {searchQuery
                                ? t("certificateSettings.step1.modal.noResults")
                                : t("certificateSettings.step1.modal.enterSearch")}
                        </Typography>
                    ) : (
                        <>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-sm text-muted-foreground mb-4"
                            >
                                {t("certificateSettings.step1.modal.description")} (
                                {tempSelectedIds.size}{" "}
                                {t("certificateSettings.step1.modal.selected")})
                            </Typography>

                            <div className="rounded-md border">
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead className="w-[50px]"></TableHead>
                                            <TableHead>
                                                {t("certificateSettings.step1.table.name")}
                                            </TableHead>
                                            <TableHead>
                                                {t("certificateSettings.step1.table.email")}
                                            </TableHead>
                                            <TableHead>
                                                {t("certificateSettings.step1.table.organization")}
                                            </TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {searchResults.map((issuer) => (
                                            <TableRow
                                                key={issuer.id}
                                                className="cursor-pointer hover:bg-muted/50"
                                                onClick={() =>
                                                    handleToggleIssuerSelection(issuer.id)
                                                }
                                            >
                                                <TableCell>
                                                    <Checkbox
                                                        checked={tempSelectedIds.has(issuer.id)}
                                                        onCheckedChange={() =>
                                                            handleToggleIssuerSelection(issuer.id)
                                                        }
                                                        onClick={(e) => e.stopPropagation()}
                                                    />
                                                </TableCell>
                                                <TableCell className="font-medium">
                                                    {issuer.name}
                                                </TableCell>
                                                <TableCell>{issuer.email}</TableCell>
                                                <TableCell>{issuer.organization || "-"}</TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </div>
                        </>
                    )}
                </div>

                <DialogFooter className="mt-6">
                    <Button
                        type="button"
                        variant="secondary-dark"
                        onClick={handleCancel}
                        disabled={isLoading}
                    >
                        <Typography variant="text" tag="span" className="font-medium">
                            {t("common.cancel")}
                        </Typography>
                    </Button>
                    <Button
                        type="button"
                        variant="primary"
                        onClick={handleConfirm}
                        disabled={tempSelectedIds.size === 0 || isLoading}
                    >
                        <Typography variant="text" tag="span" className="font-medium">
                            {t("certificateSettings.step1.modal.chooseButton")} (
                            {tempSelectedIds.size})
                        </Typography>
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};
