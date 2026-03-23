import { useState } from "react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { ExcelUpload } from "./ExcelUpload";
import { ExcelPreview } from "./ExcelPreview";

import type { EventRegistrationParticipantRequestItem } from "@decm/api";
import * as XLSX from "xlsx";
import { useImportParticipants } from "@/hooks/events/useImportParticipants";
import type { EventViewModelExtended } from "@/services/EventService/EventService";

interface ParticipantImportPageProps {
    eventId: string;
    event: EventViewModelExtended;
}

export const ParticipantImportPage = ({ eventId, event }: ParticipantImportPageProps) => {
    const { t } = useTranslation();
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [showPreview, setShowPreview] = useState(false);

    const { importParticipants, isImportingParticipants } = useImportParticipants(eventId);

    const handleFileSelect = (file: File) => {
        setSelectedFile(file);
        setShowPreview(true);
    };

    const handleCancel = () => {
        setSelectedFile(null);
        setShowPreview(false);
    };

    const handleImport = (participantData: EventRegistrationParticipantRequestItem[]) => {
        const request = {
            event_id: eventId,
            participants: participantData,
        };

        importParticipants(request);
    };

    const downloadTemplate = () => {
        // Create a simple Excel template with the required columns
        // Users should fill in EITHER email OR wallet_address (not both)
        const templateData = [
            {
                [t("participantImport.templateColumns.firstName")]: "John",
                [t("participantImport.templateColumns.lastName")]: "Doe",
                [t("participantImport.templateColumns.email")]: "john.doe@example.com",
                [t("participantImport.templateColumns.walletAddress")]: "",
                [t("participantImport.templateColumns.phoneNumber")]: "+1234567890",
                [t("participantImport.templateColumns.academicInstitution")]: "Example University",
            },
            {
                [t("participantImport.templateColumns.firstName")]: "Jane",
                [t("participantImport.templateColumns.lastName")]: "Smith",
                [t("participantImport.templateColumns.email")]: "",
                [t("participantImport.templateColumns.walletAddress")]:
                    "0x1234567890abcdef1234567890abcdef12345678",
                [t("participantImport.templateColumns.phoneNumber")]: "",
                [t("participantImport.templateColumns.academicInstitution")]: "",
            },
        ];

        // Create a new workbook and worksheet
        const wb = XLSX.utils.book_new();
        const ws = XLSX.utils.json_to_sheet(templateData);

        // Set column widths for better spacing
        const columnWidths = [
            { wch: 20 }, // first_name
            { wch: 20 }, // last_name
            { wch: 30 }, // email
            { wch: 44 }, // wallet_address
            { wch: 20 }, // phone_number
            { wch: 30 }, // academic_institution
        ];
        ws["!cols"] = columnWidths;

        XLSX.utils.book_append_sheet(wb, ws, t("participantImport.templateSheetName"));

        // Generate Excel file and download
        const excelBuffer = XLSX.write(wb, { type: "buffer" });
        const blob = new Blob([excelBuffer], {
            type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        });
        const url = URL.createObjectURL(blob);

        const link = document.createElement("a");
        link.href = url;
        link.download = `${t("participantImport.templateFileName")}.xlsx`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    };

    return (
        <div title="Participant Import">
            <div className=" mx-auto p-6">
                <div className="mb-8">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="primary"
                        className="text-3xl font-bold mb-2"
                    >
                        {t("participantImport.title")}
                    </Typography>
                    <Typography variant="text" tag="p" color="foreground-alt" className="text-lg">
                        {t("participantImport.description")}
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
                        {t("participantImport.eventDetails")}
                    </Typography>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background"
                                className="text-sm"
                            >
                                {t("participantImport.eventName")}
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
                                {t("participantImport.eventDate")}
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
                                {t("participantImport.eventLocation")}
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
                                {t("participantImport.maxAttendees")}
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

                {/* Required Fields Info Section */}
                <div className="bg-card rounded-lg p-6 mb-8 shadow-sm">
                    <Typography
                        variant="header"
                        tag="h2"
                        color="background"
                        className="text-xl font-semibold mb-4"
                    >
                        {t("participantImport.requiredFormat")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="background-alt"
                        className="text-sm mb-4"
                    >
                        {t("participantImport.requiredFormatDescription")}
                    </Typography>

                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
                        {/* Always-required Fields */}
                        <div className="flex items-center space-x-3 p-3 bg-primary/10 rounded-lg border border-primary/20">
                            <div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="foreground"
                                    className="text-sm font-bold"
                                >
                                    ✓
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    first_name
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="primary"
                                    className="text-xs font-medium"
                                >
                                    {t("participantImport.columnHeader")}
                                </Typography>
                            </div>
                        </div>
                        <div className="flex items-center space-x-3 p-3 bg-primary/10 rounded-lg border border-primary/20">
                            <div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="foreground"
                                    className="text-sm font-bold"
                                >
                                    ✓
                                </Typography>
                            </div>
                            <div>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background"
                                    className="font-semibold"
                                >
                                    last_name
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="primary"
                                    className="text-xs font-medium"
                                >
                                    {t("participantImport.columnHeader")}
                                </Typography>
                            </div>
                        </div>

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
                                    email
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("participantImport.eitherOrColumn")}
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
                                    wallet_address
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("participantImport.eitherOrColumn")}
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
                                    phone_number
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("participantImport.optionalColumn")}
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
                                    academic_institution
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="background-alt"
                                    className="text-xs font-medium"
                                >
                                    {t("participantImport.optionalColumn")}
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
                                {t("participantImport.downloadTemplate")}
                            </Typography>
                        </Button>
                        <Typography
                            variant="text"
                            tag="p"
                            color="background-alt"
                            className="text-sm"
                        >
                            {t("participantImport.downloadTemplateDescription")}
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
                            {t("participantImport.uploadFile")}
                        </Typography>

                        <ExcelUpload
                            onFileSelect={handleFileSelect}
                            selectedFile={selectedFile}
                            disabled={isImportingParticipants}
                        />
                    </div>
                ) : (
                    <ExcelPreview
                        file={selectedFile!}
                        onConfirm={handleImport}
                        onCancel={handleCancel}
                        disabled={isImportingParticipants}
                    />
                )}
            </div>
        </div>
    );
};
