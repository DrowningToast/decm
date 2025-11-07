import { useState, useEffect } from "react";
import * as XLSX from "xlsx";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useTranslation } from "react-i18next";
import type { EventCertificateImportRequestItem } from "@decm/api";

interface ExcelPreviewProps {
    file: File;
    onConfirm: (certificates: EventCertificateImportRequestItem[]) => void;
    onCancel: () => void;
    disabled?: boolean;
}

interface PreviewData {
    [key: string]: string;
}

// Fixed column names that Excel files must have
const REQUIRED_COLUMNS = {
    firstName: "first_name",
    lastName: "last_name",
    email: "email",
    academicInstitution: "academic_institution",
    certificateTitle: "certificate_title",
    certificateSubtitle: "certificate_subtitle",
};

export const ExcelPreview = ({
    file,
    onConfirm,
    onCancel,
    disabled = false,
}: ExcelPreviewProps) => {
    const { t } = useTranslation();
    const [previewData, setPreviewData] = useState<PreviewData[]>([]);
    const [validationError, setValidationError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);

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
                    const jsonData = XLSX.utils.sheet_to_json(worksheet) as PreviewData[];

                    if (jsonData.length > 0) {
                        const excelColumns = Object.keys(jsonData[0]);

                        // Validate that all required columns exist
                        const missingColumns = Object.values(REQUIRED_COLUMNS).filter(
                            (col) => !excelColumns.includes(col),
                        );

                        if (missingColumns.length > 0) {
                            setValidationError(
                                t("certificateImport.missingColumns", {
                                    columns: missingColumns.join(", "),
                                }),
                            );
                        } else {
                            // Show only first 10 rows for preview
                            // const preview = jsonData.slice(0, 10);
                            console.log(jsonData);
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

    const handleConfirm = () => {
        if (validationError) return;

        const request = previewData.map((row) => ({
            first_name: row.first_name,
            last_name: row.last_name,
            email: row.email,
            academic_institution: row.academic_institution,
            certificate_title: row.certificate_title,
            certificate_subtitle: row.certificate_subtitle,
        }));

        onConfirm(request);
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
                    {Object.entries(REQUIRED_COLUMNS).map(([key, value]) => (
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
                                {Object.values(REQUIRED_COLUMNS).map((column) => (
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
                                    {Object.values(REQUIRED_COLUMNS).map((column) => (
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
                <Button onClick={handleConfirm} disabled={disabled || !!validationError}>
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
    );
};
