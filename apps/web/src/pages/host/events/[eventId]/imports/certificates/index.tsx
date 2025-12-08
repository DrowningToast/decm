import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { ImportEventCertificatesPage } from "@/components/pages/HostPages/ImportEventCertificatesPage";
import { useTranslation } from "react-i18next";
import { useParams } from "@/router";
import { useEvent } from "@/hooks/events/useEvent";

const CertificateImportPage = () => {
    const { t } = useTranslation();
    const { eventId } = useParams("/host/events/:eventId/imports/certificates");
    const { event, isLoadingEvent, isLoadingEventError } = useEvent(eventId);

    if (isLoadingEvent) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
            </div>
        );
    }

    if (isLoadingEventError || !event) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="text-center">
                    <p className="text-red-500">{t("errors.failedToLoadEvent")}</p>
                </div>
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
