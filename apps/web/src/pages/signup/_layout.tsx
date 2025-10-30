import { BaseLayout } from "@/components/layouts/BaseLayout";
import { PublicNavbar } from "@/components/layouts/navigations/PublicNavbar";
import { Outlet } from "react-router-dom";

const Layout = () => {
    return (
        <BaseLayout className="pt-[68px]" variant="dark">
            <PublicNavbar variant="dark" />
            <div className="pt-3 md:pt-[68px]">
                <Outlet />
            </div>
        </BaseLayout>
    )
}

export default Layout;