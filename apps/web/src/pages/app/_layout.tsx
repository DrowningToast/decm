import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { BaseLayout } from "@/components/layouts/BaseLayout";
import { PrivateNavbar } from "@/components/layouts/navigations/PrivateNavbar";
import { BetaFooter } from "@/components/layouts/BetaFooter";
import { Outlet } from "react-router-dom";

const AppLayout = () => {
    return (
        <ProtectedRoute>
            <BaseLayout variant="dark">
                <PrivateNavbar />
                <Outlet />
                <BetaFooter />
            </BaseLayout>
        </ProtectedRoute>
    );
};

export default AppLayout;
