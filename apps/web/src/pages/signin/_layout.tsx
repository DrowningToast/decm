import { BaseLayout } from "@/components/layouts/BaseLayout";
import { PublicNavbar } from "@/components/layouts/navigations/PublicNavbar";
import { Outlet } from "react-router-dom";

const SignInLayout = () => {
    return (
        <BaseLayout className="max-h-screen w-full overflow-hidden" variant="dark">
            <PublicNavbar variant="dark" />
            <div className="py-4 md:py-14">
                <Outlet />
            </div>
        </BaseLayout>
    );
};

export default SignInLayout;
