import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { ParticipantImportPage } from "@/components/pages/HostPages/ParticipantImportPage/ParticipantImportPage";
import { useParams } from "@/router";
import { useEvent } from "@/hooks/events/useEvent";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId/imports/participants");
    const { event, isLoadingEvent, isLoadingEventError } = useEvent(eventId);

    if (isLoadingEvent) {
        return <div>Loading event...</div>;
    }

    if (isLoadingEventError || !event) {
        return <div>Error loading event</div>;
    }

    return (
        <>
            <FaviconHelmet
                title={`Import Participants | ${event.title} | DECM`}
                description={`Import participants for ${event.title} event`}
            />
            <ParticipantImportPage eventId={eventId} event={event} />
        </>
    );
}
