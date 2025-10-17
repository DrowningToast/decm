import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { LogoutButton } from "@/components/LogoutButton";

const AppPage = () => {
    return (
        <ProtectedRoute>
            <h1>app page</h1>
            {/* PH */}
            <LogoutButton type="signout" />
        </ProtectedRoute>
    );
};

export default AppPage;
