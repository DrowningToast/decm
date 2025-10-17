import HostHomePage from "@/components/pages/HostPages/HomePage/HostHomePage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function Page() {
    return (
        <ProtectedRoute>
            <HostHomePage />
        </ProtectedRoute>
    );
}
