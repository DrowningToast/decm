import { CertificateSettingsPage } from "@/components/pages/HostPages/EventPages/CertificateSettingsPage";
import { useEventCertificateConfig } from "@/components/pages/HostPages/EventPages/useEventCertificateConfig";
import { useEventIssuers } from "@/components/pages/HostPages/EventPages/useEventIssuers";
import { useVerifiedIssuers } from "@/hooks/events/useVerifiedIssuers";
import { useParams } from "react-router-dom";

export default function Page() {
    const { eventId } = useParams<{ eventId: string }>();

    const { data: eventCertificateConfig, isLoading: isLoadingEventCertificateConfig } =
        useEventCertificateConfig(eventId!);

    const { eventIssuers, isLoadingEventIssuers } = useEventIssuers(eventId!);
    const { verifiedIssuers, isLoadingVerifiedIssuers } = useVerifiedIssuers();

    if (isLoadingEventCertificateConfig || isLoadingVerifiedIssuers || isLoadingEventIssuers) {
        return <div>Loading...</div>;
    }

    return (
        <CertificateSettingsPage
            eventId={eventId!}
            eventCertificateConfig={eventCertificateConfig}
            verifiedIssuers={verifiedIssuers}
            eventIssuers={eventIssuers}
        />
    );
}
