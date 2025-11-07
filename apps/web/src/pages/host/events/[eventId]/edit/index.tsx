import { EditEventPage } from "@/components/pages/HostPages/EditEventPage/EditEventPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useParams } from "@/router";
import { useEvent } from "@/hooks/events/useEvent";
import { useEventContract } from "@/hooks/events/useEventContracts";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId/edit");
    const { event, isLoadingEvent, isLoadingEventError } = useEvent(eventId);
    const { data: eventContract, isLoading: isLoadingEventContract } = useEventContract(eventId!);

    const isLoading = isLoadingEvent || isLoadingEventContract;
    const isError = isLoadingEventError;

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (isError || !event || !eventContract) {
        return <div>Error loading event</div>;
    }

    return (
        <ProtectedRoute>
            <EditEventPage event={event} eventContract={eventContract} />
        </ProtectedRoute>
    );
}
