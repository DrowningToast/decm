import { LogoutButton } from "@/components/LogoutButton";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

const AppPage = () => {
    const { t } = useTranslation();

    return (
        <ProtectedRoute>
            <Typography variant="header" tag="h1">
                {t("app.title")}
            </Typography>
            {/* PH */}
            <LogoutButton type="signout" />
        </ProtectedRoute>
    );
};

export default AppPage;
