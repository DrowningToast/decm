import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { ImportEventCertificatesPage } from "@/components/pages/HostPages/ImportEventCertificatesPage";
import { useTranslation } from "react-i18next";
import { useState, useEffect } from "react";
import type { EventEventResponse } from "@decm/api";

import { useParams } from "@/router";

// Mock event data for now since backend isn't ready
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const mockEvent: any = {
    id: "1",
    title: "Sample Event",
    short_description: "This is a sample event for certificate import",
    long_description:
        "This is a detailed description of the sample event for certificate import functionality",
    location: "Sample Location",
    start_date: "2023-12-01T00:00:00Z",
    end_date: "2023-12-31T23:59:59Z",
    max_attendees: 100,
    is_public: true,
    is_verified: true,
    is_booking_request_required: false,
    is_ticket_transferable: true,
    icon_presigned_url: "",
    banner_presigned_url: "",
    google_map_query: "sample location",
    contact_number: "123-456-7890",
};

const CertificateImportPage = () => {
    const { t } = useTranslation();
    const { eventId } = useParams("/host/events/:eventId/imports/certificates");
    const [event, setEvent] = useState<EventEventResponse | null>(null);

    useEffect(() => {
        // In a real implementation, we would fetch event data based on eventId
        // For now, we'll use mock data
        setEvent(mockEvent);
    }, [eventId]);

    if (!event) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
            </div>
        );
    }

    return (
        <>
            <FaviconHelmet
                title={`${t("certificateImport.title")} | ${t("common.appName")}`}
                description={t("certificateImport.description")}
            />
            <ImportEventCertificatesPage eventId={eventId} event={event} />
        </>
    );
};

export default CertificateImportPage;
