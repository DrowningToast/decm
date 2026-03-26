import { useState, useEffect } from "react";
import * as XLSX from "xlsx";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useTranslation } from "react-i18next";
import type { EventImportCertificateReceiverRequest } from "@decm/api";
import { PasswordPinModal } from "@/components/ui/password-pin-modal";

interface ExcelPreviewProps {
    file: File;
    onConfirm: (certificates: EventImportCertificateReceiverRequest[], pin: string) => void;
    onCancel: () => void;
    disabled?: boolean;
    importError?: string | null;
}

interface PreviewData {
    [key: string]: string;
}

// Column names for Excel files (email is required, others are optional)
const EXPECTED_COLUMNS = {
    email: "email",
    firstName: "first_name",
    lastName: "last_name",
    academicInstitution: "academic_institution",
    certificateTitle: "certificate_title",
    certificateSubtitle: "certificate_subtitle",
};

export const ExcelPreview = ({
    file,
    onConfirm,
    onCancel,
    disabled = false,
    importError = null,
}: ExcelPreviewProps) => {
    const { t } = useTranslation();
    const [previewData, setPreviewData] = useState<PreviewData[]>([]);
    const [validationError, setValidationError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);

    const [showHostPinModal, setShowHostPinModal] = useState(false);

    useEffect(() => {
        if (file) {
            setIsLoading(true);
            setValidationError(null);
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
                        // Validate that email is present for all rows
                        const missingEmailRows: number[] = [];
                        jsonData.forEach((row, index) => {
                            if (!row.email || row.email.trim() === "") {
                                missingEmailRows.push(index + 1);
                            }
                        });

                        if (missingEmailRows.length > 0) {
                            setValidationError(
                                t("certificateImport.missingEmailError", {
                                    rows: missingEmailRows.join(", "),
                                }) ||
                                    `Email is required. Missing in rows: ${missingEmailRows.join(", ")}`,
                            );
                        } else {
                            setPreviewData(jsonData);
                        }
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

    const handleConfirm = (pin: string) => {
        if (validationError) return;

        const request = previewData.map((row) => ({
            email: row.email,
            first_name: row.first_name,
            last_name: row.last_name,
            academic_institution: row.academic_institution,
            certificate_title: row.certificate_title,
            certificate_subtitle: row.certificate_subtitle,
        }));

        onConfirm(request, pin);
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
                    {Object.entries(EXPECTED_COLUMNS).map(([key, value]) => (
                        <div key={key} className="flex items-center space-x-2">
                            <div className="w-3 h-3 rounded-full bg-primary/10 flex items-center justify-center">
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
                                    {t("certificateImport.columnHeader")}
                                </Typography>
                            </div>
                        </div>
                    ))}
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
                                {Object.values(EXPECTED_COLUMNS).map((column) => (
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
                            {previewData.map((row, index) => (
                                <tr key={index} className="bg-background">
                                    {Object.values(EXPECTED_COLUMNS).map((column) => (
                                        <td
                                            key={column}
                                            className="border border-border p-2 text-sm"
                                        >
                                            {row[column] || ""}
                                        </td>
                                    ))}
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Warning Alert */}
            {!validationError && !importError && (
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
                    {(validationError || importError) && (
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
                        onClick={() => setShowHostPinModal(true)}
                        disabled={disabled || !!validationError || !!importError}
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
                    handleConfirm(result.value);
                }}
                showSigningDetails
                signingDetails={{
                    details: t("signing.details.importCertificatesDescription"),
                    transactionType: t("signing.details.importCertificatesReceivers"),
                }}
            />
        </div>
    );
};
