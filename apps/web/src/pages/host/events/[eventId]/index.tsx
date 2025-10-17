import HostEventDetailsPage from "@/components/pages/HostPages/EventsPage/HostEventDetailsPage";
import { useParams } from "@/router";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useEvent } from "@/hooks/events/useEvent";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId");
    const { event, isLoadingEvent, isLoadingEventError } = useEvent(eventId);

    if (isLoadingEvent) {
        return <div>Loading event...</div>;
    }

    if (isLoadingEventError || !event) {
        return <div>Error loading event</div>;
    }

    return (
        <ProtectedRoute>
            <HostEventDetailsPage eventId={eventId} event={event} />
        </ProtectedRoute>
    );
}
