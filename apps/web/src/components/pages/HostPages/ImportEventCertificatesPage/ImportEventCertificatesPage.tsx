import { useState } from "react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { ExcelUpload } from "./ExcelUpload";
import { ExcelPreview } from "./ExcelPreview";

import type { EventImportCertificateReceiverRequest } from "@decm/api";
import { useImportCertificates } from "@/hooks/events/useImportCertificates";
import * as XLSX from "xlsx";
import type { EventViewModelExtended } from "@/services/EventService/EventService";

interface ImportEventCertificatesPageProps {
    eventId: string;
    event: EventViewModelExtended;
}

export const ImportEventCertificatesPage = ({
    eventId,
    event,
}: ImportEventCertificatesPageProps) => {
    const { t } = useTranslation();
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [showPreview, setShowPreview] = useState(false);
    const [importError, setImportError] = useState<string | null>(null);

    const handleCancel = () => {
        setSelectedFile(null);
        setShowPreview(false);
        setImportError(null);
    };

    const { importCertificates, isImportingCertificates, getSignMessage, isGettingSignMessage } =
        useImportCertificates(eventId, {
            onError: (error: Error) => {
                // On error, keep the preview open so user can see the error and re-upload
                // The ExcelPreview component will show a "Re-upload File" button
                setImportError(error.message || "An error occurred during import");
            },
        });

    const [walletSignMessage, setWalletSignMessage] = useState<string | null>(null);

    const handleFileSelect = (file: File) => {
        setSelectedFile(file);
        setShowPreview(true);
        setImportError(null); // Clear any previous errors when selecting a new file
        setWalletSignMessage(null);
    };

    const handlePreSign = async (certificateData: EventImportCertificateReceiverRequest[]) => {
        const result = await getSignMessage(certificateData);
        setWalletSignMessage(result.signMessage);
        return result.signMessage;
    };

    const handleImport = (
        certificateData: EventImportCertificateReceiverRequest[],
        auth:
            | { type: "pin" | "password"; value: string }
            | { type: "wallet"; signature: string; signMessage: string },
    ) => {
        if (auth.type === "wallet") {
            importCertificates({
                hostSignMessage: auth.signMessage,
                hostSignature: auth.signature,
                receivers: certificateData,
            });
        } else {
            importCertificates({
                hostPin: auth.value,
                receivers: certificateData,
            });
        }
    };

    const downloadTemplate = () => {
        // Create a simple Excel template with the required columns
        // Users should fill in EITHER email OR wallet_address (not both)
        const templateData = [
            {
                [t("certificateImport.templateColumns.email")]: "john.doe@example.com",
                [t("certificateImport.templateColumns.walletAddress")]: "",
                [t("certificateImport.templateColumns.firstName")]: "John",
                [t("certificateImport.templateColumns.lastName")]: "Doe",
                [t("certificateImport.templateColumns.academicInstitution")]: "Example University",
                [t("certificateImport.templateColumns.certificateTitle")]:
                    "Certificate of Achievement",
                [t("certificateImport.templateColumns.certificateSubtitle")]:
                    "For outstanding performance in the event",
            },
            {
                [t("certificateImport.templateColumns.email")]: "",
                [t("certificateImport.templateColumns.walletAddress")]:
                    "0x1234567890abcdef1234567890abcdef12345678",
                [t("certificateImport.templateColumns.firstName")]: "Jane",
                [t("certificateImport.templateColumns.lastName")]: "Smith",
                [t("certificateImport.templateColumns.academicInstitution")]: "",
                [t("certificateImport.templateColumns.certificateTitle")]: "",
                [t("certificateImport.templateColumns.certificateSubtitle")]: "",
            },
        ];

        // Create a new workbook and worksheet
        const wb = XLSX.utils.book_new();
        const ws = XLSX.utils.json_to_sheet(templateData, { header: Object.keys(templateData[0]) });

        // Set column widths for better spacing
        const columnWidths = [
            { wch: 30 }, // email
            { wch: 44 }, // wallet_address
            { wch: 15 }, // first_name
            { wch: 15 }, // last_name
            { wch: 25 }, // academic_institution
            { wch: 30 }, // certificate_title
            { wch: 40 }, // certificate_subtitle
        ];
        ws["!cols"] = columnWidths;

        XLSX.utils.book_append_sheet(wb, ws, t("certificateImport.templateSheetName"));

        // Generate Excel file and download
        const excelBuffer = XLSX.write(wb, { type: "buffer" });
        const blob = new Blob([excelBuffer], {
            type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        });
        const url = URL.createObjectURL(blob);

        const link = document.createElement("a");
        link.href = url;
        link.download = `${t("certificateImport.templateFileName")}.xlsx`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    };

    return (
        <div title="Certificate Import">
            <div className="mx-auto p-6">
                <div className="mb-8">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="primary"
                        className="text-3xl font-bold mb-2"
                    >
                        {t("certificateImport.title")}
                    </Typography>
                    <Typography variant="text" tag="p" color="foreground-alt" className="text-lg">
                        {t("certificateImport.description")}
                    </Typography>
                </div>

                {/* Event Details Section */}
                <div className="bg-card rounded-lg p-6 mb-8 shadow-sm">
                    <Typography
                        variant="header"
                        tag="h2"
                        color="background"
                        className="text-xl font-semibold mb-4"
                    >
                        {t("certificateImport.eventDetails")}
                    </Typography>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background"
                                className="text-sm"
                            >
                                {t("certificateImport.eventName")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background-alt"
                                className="font-medium"
                            >
                                {event.title}
                            </Typography>
                        </div>
                        <div>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background"
                                className="text-sm"
                            >
                                {t("certificateImport.eventDate")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background-alt"
                                className="font-medium"
                            >
                                {event.startDate
                                    ? new Date(event.startDate).toLocaleDateString()
                                    : "N/A"}{" "}
                                -{" "}
                                {event.endDate
                                    ? new Date(event.endDate).toLocaleDateString()
                                    : "N/A"}
                            </Typography>
                        </div>
                        <div>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background"
                                className="text-sm"
                            >
                                {t("certificateImport.eventLocation")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background-alt"
                                className="font-medium"
                            >
                                {event.location}
                            </Typography>
                        </div>
                        <div>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background"
                                className="text-sm"
                            >
                                {t("certificateImport.maxAttendees")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background-alt"
                                className="font-medium"
                            >
                                {event.maxAttendees}
                            </Typography>
                        </div>
                    </div>
                </div>

                {/* Required Format Info Section */}
                <div className="bg-card rounded-lg p-6 mb-8 shadow-sm">
                    <Typography
                        variant="header"
                        tag="h2"
                        color="background"
                        className="text-xl font-semibold mb-4"
                    >
                        {t("certificateImport.requiredFormat")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="background-alt"
                        className="text-sm mb-4"
                    >
                        {t("certificateImport.requiredFormatDescription")}
                    </Typography>

                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
                        {/* Either/Or Fields — exactly one must be filled per row */}
                        <div className="flex items-start space-x-3 p-3 bg-yellow-500/10 rounded-lg border border-yellow-500/40">
                            <div className="w-8 h-8 rounded-full bg-yellow-500 flex items-center justify-center flex-shrink-0">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="foreground"
                                    className="text-sm font-bold"
                                >
                                    ⊕
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    {t("certificateImport.templateColumns.email")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("certificateImport.eitherOrColumn")}
                                </Typography>
                            </div>
                        </div>
                        <div className="flex items-start space-x-3 p-3 bg-yellow-500/10 rounded-lg border border-yellow-500/40">
                            <div className="w-8 h-8 rounded-full bg-yellow-500 flex items-center justify-center flex-shrink-0">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="foreground"
                                    className="text-sm font-bold"
                                >
                                    ⊕
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    {t("certificateImport.templateColumns.walletAddress")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("certificateImport.eitherOrColumn")}
                                </Typography>
                            </div>
                        </div>

                        {/* Optional Fields */}
                        <div className="flex items-center space-x-3 p-3 bg-muted/50 rounded-lg border border-border">
                            <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="background"
                                    className="text-sm font-bold"
                                >
                                    ~
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    {t("certificateImport.templateColumns.firstName")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("certificateImport.optionalColumn")}
                                </Typography>
                            </div>
                        </div>
                        <div className="flex items-center space-x-3 p-3 bg-muted/50 rounded-lg border border-border">
                            <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="background"
                                    className="text-sm font-bold"
                                >
                                    ~
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    {t("certificateImport.templateColumns.lastName")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("certificateImport.optionalColumn")}
                                </Typography>
                            </div>
                        </div>
                        <div className="flex items-center space-x-3 p-3 bg-muted/50 rounded-lg border border-border">
                            <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="background"
                                    className="text-sm font-bold"
                                >
                                    ~
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    {t("certificateImport.templateColumns.academicInstitution")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("certificateImport.optionalColumn")}
                                </Typography>
                            </div>
                        </div>
                        <div className="flex items-center space-x-3 p-3 bg-muted/50 rounded-lg border border-border">
                            <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="background"
                                    className="text-sm font-bold"
                                >
                                    ~
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    {t("certificateImport.templateColumns.certificateTitle")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("certificateImport.optionalColumn")}
                                </Typography>
                            </div>
                        </div>
                        <div className="flex items-center space-x-3 p-3 bg-muted/50 rounded-lg border border-border">
                            <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="background"
                                    className="text-sm font-bold"
                                >
                                    ~
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    {t("certificateImport.templateColumns.certificateSubtitle")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("certificateImport.optionalColumn")}
                                </Typography>
                            </div>
                        </div>
                    </div>

                    <div className="flex items-center gap-4">
                        <Button variant="secondary-light" onClick={downloadTemplate}>
                            <Typography
                                variant="text"
                                tag="span"
                                color="background"
                                className="font-medium"
                            >
                                {t("certificateImport.downloadTemplate")}
                            </Typography>
                        </Button>
                        <Typography
                            variant="text"
                            tag="p"
                            color="background-alt"
                            className="text-sm"
                        >
                            {t("certificateImport.downloadTemplateDescription")}
                        </Typography>
                    </div>
                </div>

                {/* File Upload Section */}
                {!showPreview ? (
                    <div className="bg-card rounded-lg p-6 shadow-sm">
                        <Typography
                            variant="header"
                            tag="h2"
                            color="background"
                            className="text-xl font-semibold mb-4"
                        >
                            {t("certificateImport.uploadFile")}
                        </Typography>

                        <ExcelUpload
                            onFileSelect={handleFileSelect}
                            selectedFile={selectedFile}
                            disabled={isImportingCertificates}
                        />
                    </div>
                ) : (
                    <ExcelPreview
                        file={selectedFile!}
                        onConfirm={handleImport}
                        onPreSign={handlePreSign}
                        onCancel={handleCancel}
                        disabled={isImportingCertificates || isGettingSignMessage}
                        importError={importError}
                        walletSignMessage={walletSignMessage ?? undefined}
                    />
                )}
            </div>
        </div>
    );
};
