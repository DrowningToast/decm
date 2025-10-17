import HostEventDetailsPage from "@/components/pages/HostPages/EventsPage/HostEventDetailsPage";
import { useParams } from "@/router";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId");

    return (
        <ProtectedRoute>
            <HostEventDetailsPage eventId={eventId} />
        </ProtectedRoute>
    );
}
