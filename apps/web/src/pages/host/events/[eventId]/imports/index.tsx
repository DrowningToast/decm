import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { ParticipantImportPage } from "@/components/pages/HostPages/ParticipantImportPage/ParticipantImportPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useParams } from "@/router";
import { useEvent } from "@/hooks/events/useEvent";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId/imports");
    const { event, isLoadingEvent, isLoadingEventError } = useEvent(eventId);

    if (isLoadingEvent) {
        return <div>Loading event...</div>;
    }

    if (isLoadingEventError || !event) {
        return <div>Error loading event</div>;
    }

    return (
        <ProtectedRoute>
            <FaviconHelmet
                title={`Import Participants | ${event.title} | DECM`}
                description={`Import participants for ${event.title} event`}
            />
            <ParticipantImportPage eventId={eventId} event={event} />
        </ProtectedRoute>
    );
}
