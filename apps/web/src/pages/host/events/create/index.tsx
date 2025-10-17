import { CreateEventPage } from "@/components/pages/HostPages/CreateEventPage/CreateEventPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function Page() {
    return (
        <ProtectedRoute>
            <CreateEventPage />
        </ProtectedRoute>
    );
}
