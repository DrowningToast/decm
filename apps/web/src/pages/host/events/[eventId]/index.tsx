import HostEventDetailsPage from "@/components/pages/HostPages/EventsPage/HostEventDetailsPage";
import { useParams } from "@/router";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useEvent } from "@/hooks/events/useEvent";
import { useEventRegistrationConfig } from "@/hooks/events/useEventRegistrationConfig";
import { useEventCertificateConfig } from "@/components/pages/HostPages/EventPages/useEventCertificateConfig";
import { useEventIssuers } from "@/components/pages/HostPages/EventPages/useEventIssuers";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId");
    const { event, isLoadingEvent, isLoadingEventError } = useEvent(eventId);

    const {
        data: eventRegistrationConfig,
        isLoading: isLoadingEventRegistrationConfig,
        error: isErrorEventRegistrationConfig,
    } = useEventRegistrationConfig(eventId);

    const { data: eventCertificateConfig, isLoading: isLoadingEventCertificateConfig } =
        useEventCertificateConfig(eventId!);

    const { eventIssuers, isLoadingEventIssuers } = useEventIssuers(eventId!);

    const isLoading =
        isLoadingEvent ||
        isLoadingEventRegistrationConfig ||
        isLoadingEventCertificateConfig ||
        isLoadingEventIssuers;

    if (isLoading) {
        return <div>Loading event...</div>;
    }

    if (
        isLoadingEventError ||
        isErrorEventRegistrationConfig ||
        !event ||
        !eventRegistrationConfig
    ) {
        return <div>Error loading event</div>;
    }

    return (
        <ProtectedRoute>
            <HostEventDetailsPage
                eventId={eventId}
                event={event}
                eventRegistrationConfig={eventRegistrationConfig}
                eventCertificateConfig={eventCertificateConfig}
                eventIssuers={eventIssuers}
            />
        </ProtectedRoute>
    );
}
