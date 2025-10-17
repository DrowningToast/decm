import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { LogoutButton } from "@/components/pages/Onboard/OAuth/LogoutButton";

const AppPage = () => {
    return (
        <ProtectedRoute>
            <h1>app page</h1>
            {/* PH */}
            <LogoutButton />
        </ProtectedRoute>
    );
};

export default AppPage;
