import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { BaseLayout } from "@/components/layouts/BaseLayout";
import { PrivateNavbar } from "@/components/layouts/navigations/PrivateNavbar";
import { Outlet } from "react-router-dom";

const HostLayout = () => {
    return (
        <ProtectedRoute>
            <BaseLayout variant="dark">
                <PrivateNavbar currentRole="Verified Organizer" />
                <Outlet />
            </BaseLayout>
        </ProtectedRoute>
    );
};

export default HostLayout;
