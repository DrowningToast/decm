import { EditProfilePage } from "@/components/pages/Participant/Settings/EditProfilePage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

const SettingsRoute = () => {
    return (
        <ProtectedRoute>
            <EditProfilePage />
        </ProtectedRoute>
    );
};

export default SettingsRoute;
