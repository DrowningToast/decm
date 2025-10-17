import { EditEventPage } from "@/components/pages/HostPages/EditEventPage/EditEventPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useParams } from "@/router";
import { useEvent } from "@/hooks/events/useEvent";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId/edit");
    const { event, isLoadingEvent, isLoadingEventError } = useEvent(eventId);

    const isLoading = isLoadingEvent;
    const isError = isLoadingEventError;

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (isError || !event) {
        return <div>Error loading event</div>;
    }

    return (
        <ProtectedRoute>
            <EditEventPage event={event} />
        </ProtectedRoute>
    );
}
