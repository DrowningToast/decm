import { EditEventPage } from "@/components/pages/HostPages/EditEventPage/EditEventPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function Page() {
    return (
        <ProtectedRoute>
            <EditEventPage />
        </ProtectedRoute>
    );
}
