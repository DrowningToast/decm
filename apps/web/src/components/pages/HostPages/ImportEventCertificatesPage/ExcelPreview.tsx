import { useState, useEffect } from "react";
import * as XLSX from "xlsx";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useTranslation } from "react-i18next";
import type { EventImportCertificateReceiverRequest } from "@decm/api";
import {
    PasswordPinModal,
    type PasswordPinModalSuccessResult,
} from "@/components/ui/password-pin-modal";
import {
    CERTIFICATE_ALL_COLUMNS,
    CERTIFICATE_EITHER_OR_COLUMNS,
    validateCertificateRow,
    buildCertificate,
    type PreviewData,
    type RowValidationResult,
} from "./ExcelPreviewUtils";

interface ExcelPreviewProps {
    file: File;
    onConfirm: (
        certificates: EventImportCertificateReceiverRequest[],
        auth: PasswordPinModalSuccessResult,
    ) => void;
    /** Called before showing the signing modal to fetch the sign message from the backend */
    onPreSign?: (certificates: EventImportCertificateReceiverRequest[]) => Promise<string>;
    onCancel: () => void;
    disabled?: boolean;
    importError?: string | null;
    /** Pre-fetched wallet sign message for BYOK users */
    walletSignMessage?: string;
}

export const ExcelPreview = ({
    file,
    onConfirm,
    onPreSign,
    onCancel,
    disabled = false,
    importError = null,
    walletSignMessage,
}: ExcelPreviewProps) => {
    const { t } = useTranslation();
    const [previewData, setPreviewData] = useState<PreviewData[]>([]);
    const [rowValidations, setRowValidations] = useState<RowValidationResult[]>([]);
    const [missingColumns, setMissingColumns] = useState<string[]>([]);
    const [validationError, setValidationError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);

    const [showHostPinModal, setShowHostPinModal] = useState(false);

    useEffect(() => {
        if (file) {
            setIsLoading(true);
            setValidationError(null);
            setMissingColumns([]);
            const reader = new FileReader();

            reader.onload = (e) => {
                try {
                    const data = new Uint8Array(e.target?.result as ArrayBuffer);
                    const workbook = XLSX.read(data, { type: "array" });
                    const sheetName = workbook.SheetNames[0];
                    const worksheet = workbook.Sheets[sheetName];
                    const jsonData = (
                        XLSX.utils.sheet_to_json(worksheet, { blankrows: false }) as PreviewData[]
                    ).filter((row) => Object.values(row).some((v) => String(v).trim() !== ""));

                    if (jsonData.length > 0) {
                        const excelColumns = Object.keys(jsonData[0]);

                        // At least one of the either/or columns must be present in the file
                        const hasAnyEitherOr = Object.values(CERTIFICATE_EITHER_OR_COLUMNS).some(
                            (col) => excelColumns.includes(col),
                        );

                        const missing = hasAnyEitherOr
                            ? []
                            : [
                                  `${CERTIFICATE_EITHER_OR_COLUMNS.email} or ${CERTIFICATE_EITHER_OR_COLUMNS.walletAddress}`,
                              ];

                        const validations = jsonData.map(validateCertificateRow);
                        setPreviewData(jsonData);
                        setRowValidations(validations);
                        setMissingColumns(missing);
                    } else {
                        setValidationError(t("certificateImport.emptyFile"));
                    }
                } catch (error) {
                    console.error("Error parsing Excel file:", error);
                    setValidationError(t("certificateImport.parseError"));
                } finally {
                    setIsLoading(false);
                }
            };

            reader.readAsArrayBuffer(file);
        }
    }, [file, t]);

    const hasInvalidRows = rowValidations.some((v) => !v.isValid);
    const hasColumnErrors = missingColumns.length > 0;
    const canConfirm = !validationError && !hasInvalidRows && !hasColumnErrors && !importError;

    const [isPreSigning, setIsPreSigning] = useState(false);

    const handleOpenModal = async () => {
        if (!canConfirm) return;
        if (onPreSign) {
            setIsPreSigning(true);
            try {
                await onPreSign(previewData.map(buildCertificate));
            } catch {
                // onPreSign error will surface via importError prop from parent
            } finally {
                setIsPreSigning(false);
            }
        }
        setShowHostPinModal(true);
    };

    const handleConfirm = (result: PasswordPinModalSuccessResult) => {
        if (!canConfirm) return;
        const request = previewData.map(buildCertificate);
        onConfirm(request, result);
    };

    if (isLoading) {
        return (
            <div className="flex justify-center items-center p-8">
                <Typography variant="text" tag="p" color="foreground">
                    {t("certificateImport.parsingFile")}
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
                {t("certificateImport.previewData")}
            </Typography>

            {/* Validation Error */}
            {validationError && (
                <Alert variant="destructive" className="mb-4">
                    <AlertDescription>{validationError}</AlertDescription>
                </Alert>
            )}

            {/* Missing columns */}
            {hasColumnErrors && (
                <Alert variant="destructive" className="mb-4">
                    <AlertDescription>
                        <Typography
                            variant="text"
                            tag="p"
                            color="destructive"
                            className="font-semibold mb-2"
                        >
                            {t("certificateImport.missingColumnsTitle")}
                        </Typography>
                        <div className="flex flex-wrap gap-2">
                            {[
                                `${CERTIFICATE_EITHER_OR_COLUMNS.email} or ${CERTIFICATE_EITHER_OR_COLUMNS.walletAddress}`,
                            ].map((col) => {
                                const isMissing = missingColumns.includes(col);
                                return (
                                    <div
                                        key={col}
                                        className={`px-2 py-1 rounded text-xs font-mono ${
                                            isMissing
                                                ? "bg-destructive/20 text-destructive border border-destructive/30"
                                                : "bg-primary/10 text-primary"
                                        }`}
                                    >
                                        {col}
                                    </div>
                                );
                            })}
                        </div>
                    </AlertDescription>
                </Alert>
            )}

            {/* Invalid rows */}
            {hasInvalidRows && !hasColumnErrors && (
                <Alert variant="destructive" className="mb-4">
                    <AlertDescription>
                        {t("certificateImport.invalidRowsDescription")}
                    </AlertDescription>
                </Alert>
            )}

            {/* Import Error */}
            {importError && (
                <Alert variant="destructive" className="mb-4">
                    <AlertDescription>
                        <Typography
                            variant="text"
                            tag="p"
                            color="destructive"
                            className="font-semibold mb-2"
                        >
                            {t("certificateImport.importError")}
                        </Typography>
                        <Typography variant="text" tag="p" color="destructive" className="text-sm">
                            {importError}
                        </Typography>
                    </AlertDescription>
                </Alert>
            )}

            {/* Required Format Info */}
            <div className="mb-6 p-4 bg-muted/30 rounded-md">
                <Typography
                    variant="header"
                    tag="h3"
                    color="background"
                    className="text-lg font-medium mb-2"
                >
                    {t("certificateImport.requiredFormat")}
                </Typography>
                <Typography variant="text" tag="p" color="background" className="text-sm mb-2">
                    {t("certificateImport.requiredFormatDescription")}
                </Typography>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {Object.entries(CERTIFICATE_ALL_COLUMNS).map(([key, value]) => {
                        const isEitherOr = Object.values(CERTIFICATE_EITHER_OR_COLUMNS).includes(
                            value as (typeof CERTIFICATE_EITHER_OR_COLUMNS)[keyof typeof CERTIFICATE_EITHER_OR_COLUMNS],
                        );
                        return (
                            <div key={key} className="flex items-center space-x-2">
                                <div
                                    className={`w-3 h-3 rounded-full flex items-center justify-center ${isEitherOr ? "bg-yellow-500/20" : "bg-primary/10"}`}
                                >
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        color="primary"
                                        className="text-xs font-bold"
                                    >
                                        {key.charAt(0).toUpperCase()}
                                    </Typography>
                                </div>
                                <div>
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        color="background"
                                        className="font-medium"
                                    >
                                        {value}
                                    </Typography>
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        color="background-alt"
                                        className="text-xs"
                                    >
                                        {isEitherOr
                                            ? t("certificateImport.eitherOrColumn")
                                            : t("certificateImport.optionalColumn")}
                                    </Typography>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            {/* Preview Table */}
            {!validationError && (
                <div className="overflow-x-auto">
                    <Typography
                        variant="header"
                        tag="h3"
                        color="background"
                        className="text-lg font-medium mb-3"
                    >
                        {t("certificateImport.previewTable")} ({previewData.length}{" "}
                        {t("certificateImport.rows")})
                    </Typography>
                    <table className="w-full border-collapse border border-border">
                        <thead>
                            <tr className="bg-primary">
                                {Object.values(CERTIFICATE_ALL_COLUMNS).map((column) => (
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
                                const isInvalid = validation && !validation.isValid;
                                return (
                                    <tr
                                        key={index}
                                        className={
                                            isInvalid ? "bg-destructive/10" : "bg-background"
                                        }
                                    >
                                        {Object.values(CERTIFICATE_ALL_COLUMNS).map((column) => (
                                            <td
                                                key={column}
                                                className="border border-border p-2 text-sm"
                                            >
                                                {row[column] || ""}
                                            </td>
                                        ))}
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Warning Alert */}
            {canConfirm && (
                <Alert variant="warning" className="mb-6">
                    <AlertDescription>
                        <Typography
                            variant="text"
                            tag="p"
                            color="background"
                            className="font-medium"
                        >
                            {t("certificateImport.receiverReplacementWarning")}
                        </Typography>
                    </AlertDescription>
                </Alert>
            )}

            {/* Action Buttons */}
            <div className="flex justify-between items-center mt-6">
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
                <div className="flex space-x-4">
                    {(!canConfirm || importError) && (
                        <Button variant="secondary-light" onClick={onCancel} disabled={disabled}>
                            <Typography
                                variant="text"
                                tag="span"
                                color="background"
                                className="font-medium"
                            >
                                {t("certificateImport.reuploadFile") || "Re-upload File"}
                            </Typography>
                        </Button>
                    )}
                    <Button
                        onClick={handleOpenModal}
                        disabled={disabled || !canConfirm || isPreSigning}
                    >
                        <Typography
                            variant="text"
                            tag="span"
                            color="foreground-alt"
                            className="font-medium"
                        >
                            {t("certificateImport.confirmImport")}
                        </Typography>
                    </Button>
                </div>
            </div>

            <PasswordPinModal
                isOpen={showHostPinModal}
                onClose={() => setShowHostPinModal(false)}
                onSuccess={(result) => {
                    setShowHostPinModal(false);
                    handleConfirm(result);
                }}
                showSigningDetails
                signingDetails={{
                    details: t("signing.details.importCertificatesDescription"),
                    transactionType: t("signing.details.importCertificatesReceivers"),
                }}
                allowWalletSigning={!!walletSignMessage}
                walletSignMessage={walletSignMessage}
            />
        </div>
    );
};
