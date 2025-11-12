import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { BaseLayout } from "@/components/layouts/BaseLayout";
import { PrivateNavbar } from "@/components/layouts/navigations/PrivateNavbar";
import { Outlet } from "react-router-dom";

const IssuerLayout = () => {
    return (
        <ProtectedRoute>
            <BaseLayout variant="dark">
                <PrivateNavbar currentRole="Issuer" />
                <div className="py-4 md:py-14">
                    <Outlet />
                </div>
            </BaseLayout>
        </ProtectedRoute>
    );
};

export default IssuerLayout;
