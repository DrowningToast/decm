import { EventParticipantSettingPage } from "@/components/pages/HostPages/EventsPage/EventParticipantSettingPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function Page() {
    return (
        <ProtectedRoute>
            <EventParticipantSettingPage />
        </ProtectedRoute>
    );
}
