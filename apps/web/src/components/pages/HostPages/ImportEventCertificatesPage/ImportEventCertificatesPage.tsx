import { useState } from "react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { ExcelUpload } from "./ExcelUpload";
import { ExcelPreview } from "./ExcelPreview";
import PageContainer from "@/components/container/PageContainer";
// import type {
//     EventEventResponse,
//     EventCertificateImportRequest,
//     EventCertificateImportRequestItem,
// } from "@decm/api";
import * as XLSX from "xlsx";
import type { EventEventResponse } from "@decm/api";
// import { useImportCertificates } from "@/hooks/events/useImportCertificates";

interface ImportEventCertificatesPageProps {
    eventId: string;
    event: EventEventResponse;
}

export const ImportEventCertificatesPage = ({
    eventId,
    event,
}: ImportEventCertificatesPageProps) => {
    const { t } = useTranslation();
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [showPreview, setShowPreview] = useState(false);

    // const { importCertificates, isImportingCertificates } = useImportCertificates(eventId);

    const handleFileSelect = (file: File) => {
        setSelectedFile(file);
        setShowPreview(true);
    };

    const handleCancel = () => {
        setSelectedFile(null);
        setShowPreview(false);
    };

    // const handleImport = (certificateData: EventCertificateImportRequestItem[]) => {
    //     const request: EventCertificateImportRequest = {
    //         event_id: eventId,
    //         certificates: certificateData,
    //     };

    //     importCertificates(request);
    // };

    const handleImport = () => {
        console.log("handleImport");
    };

    const downloadTemplate = () => {
        // Create a simple Excel template with the required columns
        const templateData = [
            {
                [t("certificateImport.templateColumns.firstName")]: "",
                [t("certificateImport.templateColumns.lastName")]: "",
                [t("certificateImport.templateColumns.email")]: "",
                [t("certificateImport.templateColumns.academicInstitution")]: "",
                [t("certificateImport.templateColumns.certificateTitle")]: "",
                [t("certificateImport.templateColumns.certificateSubtitle")]: "",
            },
            {
                [t("certificateImport.templateColumns.firstName")]: "John",
                [t("certificateImport.templateColumns.lastName")]: "Doe",
                [t("certificateImport.templateColumns.email")]: "john.doe@example.com",
                [t("certificateImport.templateColumns.academicInstitution")]: "Example University",
                [t("certificateImport.templateColumns.certificateTitle")]:
                    "Certificate of Achievement",
                [t("certificateImport.templateColumns.certificateSubtitle")]:
                    "For outstanding performance in the event",
            },
        ];

        // Create a new workbook and worksheet
        const wb = XLSX.utils.book_new();
        const ws = XLSX.utils.json_to_sheet(templateData);
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
        <PageContainer title="Certificate Import">
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
                            {t("certificateImport.uploadFile")}
                        </Typography>

                        <div className="mb-4">
                            <Typography
                                variant="text"
                                tag="p"
                                color="background"
                                className="text-sm mb-2"
                            >
                                {t("certificateImport.downloadTemplateDescription")}
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
                                    {t("certificateImport.downloadTemplate")}
                                </Typography>
                            </Button>
                        </div>

                        <ExcelUpload
                            onFileSelect={handleFileSelect}
                            selectedFile={selectedFile}
                            disabled={false}
                        />
                    </div>
                ) : (
                    <ExcelPreview
                        file={selectedFile!}
                        onConfirm={handleImport}
                        onCancel={handleCancel}
                        disabled={false}
                    />
                )}
            </div>
        </PageContainer>
    );
};
