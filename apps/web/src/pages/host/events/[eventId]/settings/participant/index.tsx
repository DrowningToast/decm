import { EventParticipantSettingPage } from "@/components/pages/HostPages/EventsPage/EventParticipantSettingPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useEvent } from "@/hooks/events/useEvent";
import { useParams } from "@/router";
import { useEventRegistrationConfig } from "@/hooks/events/useEventRegistrationConfig";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId/settings/participant");

    const { event, isLoadingEvent } = useEvent(eventId);
    const { data: eventRegistrationConfig, isLoading: isLoadingEventRegistrationConfig } =
        useEventRegistrationConfig(eventId);

    if (isLoadingEvent || isLoadingEventRegistrationConfig) {
        return <div>Loading...</div>;
    }

    if (!event || !eventRegistrationConfig) {
        return <div>Event or event registration config not found</div>;
    }

    return (
        <ProtectedRoute>
            <EventParticipantSettingPage
                eventId={eventId}
                event={event}
                eventRegistrationConfig={eventRegistrationConfig}
            />
        </ProtectedRoute>
    );
}
