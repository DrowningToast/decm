import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { LogoutButton } from "@/components/LogoutButton";

const AppPage = () => {
    return (
        <ProtectedRoute>
            <LogoutButton type="signout" />
        </ProtectedRoute>
    );
};

export default AppPage;
