import { BaseLayout } from "@/components/layouts/BaseLayout";
import { PublicNavbar } from "@/components/layouts/navigations/PublicNavbar";
import { Outlet } from "react-router-dom";

const Layout = () => {
    return (
        <BaseLayout className="max-h-screen w-full overflow-hidden" variant="light">
            <PublicNavbar variant="light" />
            <Outlet />
        </BaseLayout>
    );
};

export default Layout;
