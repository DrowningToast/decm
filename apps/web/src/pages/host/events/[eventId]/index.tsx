import HostEventDetailsPage from "@/components/pages/HostPages/EventsPage/HostEventDetailsPage";
import { useParams } from "@/router";
import { useEvent } from "@/hooks/events/useEvent";
import { useEventRegistrationConfig } from "@/hooks/events/useEventRegistrationConfig";
import { useEventCertificateConfig } from "@/components/pages/HostPages/EventPages/useEventCertificateConfig";
import { useEventIssuers } from "@/components/pages/HostPages/EventPages/useEventIssuers";
import { useEventContract } from "@/hooks/events/useEventContracts";
import { useEventInvitedParticipants } from "@/hooks/events/useEventInvitedParticipants";

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

    const { data: eventContract, isLoading: isLoadingEventContract } = useEventContract(eventId!);

    const { eventIssuers, isLoadingEventIssuers } = useEventIssuers(eventId!);
    const { invitations, isLoading: isLoadingEventInvitations } = useEventInvitedParticipants(
        eventId!,
    );

    const isLoading =
        isLoadingEvent ||
        isLoadingEventRegistrationConfig ||
        isLoadingEventCertificateConfig ||
        isLoadingEventIssuers ||
        isLoadingEventContract ||
        isLoadingEventInvitations;

    if (isLoading) {
        return <div>Loading event...</div>;
    }

    if (
        isLoadingEventError ||
        isErrorEventRegistrationConfig ||
        !event ||
        !eventRegistrationConfig ||
        !eventContract
    ) {
        return <div>Error loading event</div>;
    }

    return (
        <HostEventDetailsPage
            eventId={eventId}
            event={event}
            eventRegistrationConfig={eventRegistrationConfig}
            eventCertificateConfig={eventCertificateConfig}
            eventIssuers={eventIssuers}
            eventContract={eventContract}
            eventInvitations={invitations}
        />
    );
}
