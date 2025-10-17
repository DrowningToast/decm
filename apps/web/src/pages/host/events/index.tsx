import HostEventPage from "@/components/pages/HostPages/EventsPage/HostEventPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function Page() {
    return (
        <ProtectedRoute>
            <HostEventPage />
        </ProtectedRoute>
    );
}
