import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { NotFound } from "@/components/pages/NotFound";
import { useTranslation } from "react-i18next";

const NotFoundPage = () => {
    const { t } = useTranslation();

    return (
        <>
            <FaviconHelmet
                title={`${t("notFound.title")} | ${t("common.appName")}`}
                description={t("notFound.description")}
            />
            <NotFound />
        </>
    );
};

export default NotFoundPage;
