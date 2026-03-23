import { useState, useEffect } from "react";
import * as XLSX from "xlsx";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useTranslation } from "react-i18next";
import type { EventRegistrationParticipantRequestItem } from "@decm/api";
import {
    ALL_COLUMNS,
    ALWAYS_REQUIRED_COLUMNS,
    EITHER_OR_COLUMNS,
    validateRow,
    buildParticipant,
    type PreviewData,
    type RowValidationResult,
} from "./ExcelPreviewUtils";

interface ExcelPreviewProps {
    file: File;
    onConfirm: (participants: EventRegistrationParticipantRequestItem[]) => void;
    onCancel: () => void;
    disabled?: boolean;
}

export const ExcelPreview = ({
    file,
    onConfirm,
    onCancel,
    disabled = false,
}: ExcelPreviewProps) => {
    const { t } = useTranslation();
    const [previewData, setPreviewData] = useState<PreviewData[]>([]);
    const [rowValidations, setRowValidations] = useState<RowValidationResult[]>([]);
    const [validationError, setValidationError] = useState<string | null>(null);
    const [missingColumns, setMissingColumns] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    const hasInvalidRows = rowValidations.some((v) => !v.isValid);

    useEffect(() => {
        if (file) {
            setIsLoading(true);
            setValidationError(null);
            setMissingColumns([]);
            setRowValidations([]);
            const reader = new FileReader();

            reader.onload = (e) => {
                try {
                    const data = new Uint8Array(e.target?.result as ArrayBuffer);
                    const workbook = XLSX.read(data, { type: "array" });
                    const sheetName = workbook.SheetNames[0];
                    const worksheet = workbook.Sheets[sheetName];
                    const jsonData = XLSX.utils.sheet_to_json(worksheet) as PreviewData[];

                    if (jsonData.length > 0) {
                        const excelColumns = Object.keys(jsonData[0]);

                        // Always-required columns must be present in the file
                        const missingAlways = Object.values(ALWAYS_REQUIRED_COLUMNS).filter(
                            (col) => !excelColumns.includes(col),
                        );

                        // At least one of the either/or columns must be present in the file
                        const hasAnyEitherOr = Object.values(EITHER_OR_COLUMNS).some((col) =>
                            excelColumns.includes(col),
                        );

                        const missing = [
                            ...missingAlways,
                            ...(hasAnyEitherOr
                                ? []
                                : [
                                      `${EITHER_OR_COLUMNS.email} or ${EITHER_OR_COLUMNS.walletAddress}`,
                                  ]),
                        ];

                        const validations = jsonData.map(validateRow);
                        setPreviewData(jsonData);
                        setRowValidations(validations);

                        if (missing.length > 0) {
                            setMissingColumns(missing);
                            setValidationError(
                                t("participantImport.missingColumns", {
                                    columns: missing.join(", "),
                                }),
                            );
                        } else if (validations.some((v) => !v.isValid)) {
                            setValidationError(t("participantImport.invalidRowsFound"));
                        }
                    } else {
                        setValidationError(t("participantImport.emptyFile"));
                    }
                } catch (error) {
                    console.error("Error parsing Excel file:", error);
                    setValidationError(t("participantImport.parseError"));
                } finally {
                    setIsLoading(false);
                }
            };

            reader.readAsArrayBuffer(file);
        }
    }, [file, t]);

    const handleConfirm = () => {
        if (validationError && !hasInvalidRows) return;
        if (hasInvalidRows) return;

        const request = previewData.map(buildParticipant);
        onConfirm(request);
    };

    if (isLoading) {
        return (
            <div className="flex justify-center items-center p-8">
                <Typography variant="text" tag="p" color="foreground">
                    {t("participantImport.parsingFile")}
                </Typography>
            </div>
        );
    }

    return (
        <div className="bg-card rounded-lg p-6 shadow-sm">
            <Typography
                variant="header"
                tag="h2"
                color="background"
                className="text-xl font-semibold mb-4"
            >
                {t("participantImport.previewData")}
            </Typography>

            {/* Missing Columns Error with Visual Indicator */}
            {missingColumns.length > 0 && (
                <Alert variant="destructive" className="mb-6">
                    <AlertDescription>
                        <Typography
                            variant="text"
                            tag="p"
                            color="destructive"
                            className="font-semibold mb-3"
                        >
                            {t("participantImport.missingColumnsTitle")}
                        </Typography>
                        <div className="flex flex-wrap gap-2">
                            {[
                                ...Object.values(ALWAYS_REQUIRED_COLUMNS),
                                `${EITHER_OR_COLUMNS.email} or ${EITHER_OR_COLUMNS.walletAddress}`,
                            ].map((col) => {
                                const isMissing = missingColumns.includes(col);
                                return (
                                    <div
                                        key={col}
                                        className={`flex items-center gap-2 px-3 py-1.5 rounded-md border ${
                                            isMissing
                                                ? "bg-destructive/10 border-destructive text-destructive"
                                                : "bg-primary/10 border-primary/30"
                                        }`}
                                    >
                                        <span className="text-sm font-medium">
                                            {isMissing ? "✗" : "✓"}
                                        </span>
                                        <span className="text-sm font-mono">{col}</span>
                                    </div>
                                );
                            })}
                        </div>
                    </AlertDescription>
                </Alert>
            )}

            {/* Row Validation Errors */}
            {hasInvalidRows && (
                <Alert variant="destructive" className="mb-4">
                    <AlertDescription>
                        <Typography
                            variant="text"
                            tag="p"
                            color="destructive"
                            className="font-semibold"
                        >
                            {t("participantImport.invalidRowsFound")}
                        </Typography>
                        <Typography variant="text" tag="p" color="destructive" className="text-sm">
                            {t("participantImport.invalidRowsDescription")}
                        </Typography>
                    </AlertDescription>
                </Alert>
            )}

            {/* Other Validation Errors (e.g., Empty File, Parse Error) */}
            {validationError && missingColumns.length === 0 && !hasInvalidRows && (
                <Alert variant="destructive" className="mb-4">
                    <AlertDescription>{validationError}</AlertDescription>
                </Alert>
            )}

            {/* Preview Table */}
            {previewData.length > 0 && (
                <div className="overflow-x-auto">
                    <Typography
                        variant="header"
                        tag="h3"
                        color="background"
                        className="text-lg font-medium mb-3"
                    >
                        {t("participantImport.previewTable")} ({previewData.length}{" "}
                        {t("participantImport.rows")})
                    </Typography>
                    <table className="w-full border-collapse border border-border">
                        <thead>
                            <tr className="bg-primary">
                                {Object.values(ALL_COLUMNS).map((column) => (
                                    <th
                                        key={column}
                                        className="border border-border p-2 text-left text-sm font-medium"
                                    >
                                        {column}
                                    </th>
                                ))}
                            </tr>
                        </thead>
                        <tbody>
                            {previewData.map((row, index) => {
                                const validation = rowValidations[index];
                                const isRowInvalid = validation && !validation.isValid;

                                return (
                                    <tr
                                        key={index}
                                        className={`${
                                            isRowInvalid ? "bg-destructive/50" : "bg-background"
                                        }`}
                                    >
                                        {Object.values(ALL_COLUMNS).map((column) => {
                                            const isFieldMissing =
                                                isRowInvalid &&
                                                validation.missingFields.includes(column);

                                            return (
                                                <td
                                                    key={column}
                                                    className={`border border-border p-2 text-sm text-foreground ${
                                                        isFieldMissing
                                                            ? "text-destructive font-bold"
                                                            : "text-foreground"
                                                    }`}
                                                >
                                                    {row[column] !== undefined && row[column] !== ""
                                                        ? String(row[column])
                                                        : "-"}
                                                </td>
                                            );
                                        })}
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Action Buttons */}
            <div className="flex justify-end space-x-4 mt-6">
                <Button variant="secondary-light" onClick={onCancel} disabled={disabled}>
                    <Typography
                        variant="text"
                        tag="span"
                        color="background"
                        className="font-medium"
                    >
                        {t("common.cancel")}
                    </Typography>
                </Button>
                <Button
                    onClick={handleConfirm}
                    disabled={disabled || !!validationError || hasInvalidRows}
                >
                    <Typography
                        variant="text"
                        tag="span"
                        color="foreground-alt"
                        className="font-medium"
                    >
                        {t("participantImport.confirmImport")}
                    </Typography>
                </Button>
            </div>
        </div>
    );
};
