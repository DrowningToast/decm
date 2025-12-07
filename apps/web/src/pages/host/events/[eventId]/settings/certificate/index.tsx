import { CertificateSettingsPage } from "@/components/pages/HostPages/EventPages/CertificateSettingsPage";
import { useEventCertificateConfig } from "@/components/pages/HostPages/EventPages/useEventCertificateConfig";
import { useEventIssuers } from "@/components/pages/HostPages/EventPages/useEventIssuers";
import { useParams } from "react-router-dom";

export default function Page() {
    const { eventId } = useParams<{ eventId: string }>();

    const { data: eventCertificateConfig, isLoading: isLoadingEventCertificateConfig } =
        useEventCertificateConfig(eventId!);

    const { eventIssuers, isLoadingEventIssuers } = useEventIssuers(eventId!);

    if (isLoadingEventCertificateConfig || isLoadingEventIssuers) {
        return <div>Loading...</div>;
    }

    return (
        <CertificateSettingsPage
            eventId={eventId!}
            eventCertificateConfig={eventCertificateConfig}
            eventIssuers={eventIssuers}
        />
    );
}
