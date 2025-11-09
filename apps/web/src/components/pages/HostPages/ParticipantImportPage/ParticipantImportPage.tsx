import { useState } from "react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { ExcelUpload } from "./ExcelUpload";
import { ExcelPreview } from "./ExcelPreview";
import PageContainer from "@/components/container/PageContainer";
import type {
    EventEventResponse,
    EventRegistrationInvitationImportEventParticipantsRequest,
    EventRegistrationInvitationParticipantRequestItem,
} from "@decm/api";
import * as XLSX from "xlsx";
import { useImportParticipants } from "@/hooks/events/useImportParticipants";

interface ParticipantImportPageProps {
    eventId: string;
    event: EventEventResponse;
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

    const handleImport = (participantData: EventRegistrationInvitationParticipantRequestItem[]) => {
        const request: EventRegistrationInvitationImportEventParticipantsRequest = {
            event_id: eventId,
            participants: participantData,
        };

        importParticipants(request);
    };

    const downloadTemplate = () => {
        // Create a simple Excel template with the required columns
        const templateData = [
            {
                [t("participantImport.templateColumns.firstName")]: "",
                [t("participantImport.templateColumns.lastName")]: "",
                [t("participantImport.templateColumns.email")]: "",
                [t("participantImport.templateColumns.phoneNumber")]: "",
                [t("participantImport.templateColumns.academicInstitution")]: "",
            },
            {
                [t("participantImport.templateColumns.firstName")]: "John",
                [t("participantImport.templateColumns.lastName")]: "Doe",
                [t("participantImport.templateColumns.email")]: "john.doe@example.com",
                [t("participantImport.templateColumns.phoneNumber")]: "+1234567890",
                [t("participantImport.templateColumns.academicInstitution")]: "Example University",
            },
        ];

        // Create a new workbook and worksheet
        const wb = XLSX.utils.book_new();
        const ws = XLSX.utils.json_to_sheet(templateData);
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
        <PageContainer title="Participant Import">
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
                                {event.start_date
                                    ? new Date(event.start_date).toLocaleDateString()
                                    : "N/A"}{" "}
                                -{" "}
                                {event.end_date
                                    ? new Date(event.end_date).toLocaleDateString()
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
                                {event.max_attendees}
                            </Typography>
                        </div>
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

                        <div className="mb-4">
                            <Typography
                                variant="text"
                                tag="p"
                                color="background"
                                className="text-sm mb-2"
                            >
                                {t("participantImport.downloadTemplateDescription")}
                            </Typography>
                            <Button
                                variant="secondary-light"
                                onClick={downloadTemplate}
                                className="mb-4"
                            >
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="background"
                                    className="font-medium"
                                >
                                    {t("participantImport.downloadTemplate")}
                                </Typography>
                            </Button>
                        </div>

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
        </PageContainer>
    );
};
