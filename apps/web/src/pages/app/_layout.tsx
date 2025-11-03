import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { BaseLayout } from "@/components/layouts/BaseLayout";
import { PrivateNavbar } from "@/components/layouts/navigations/PrivateNavbar";
import { Outlet } from "react-router-dom";

const AppLayout = () => {
    return (
        <ProtectedRoute>
            <BaseLayout>
                <Outlet />
                <PrivateNavbar />
            </BaseLayout>
        </ProtectedRoute>
    );
};

export default AppLayout;
