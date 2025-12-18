import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { ErrorPage } from "@/components/pages/Error";
import { useTranslation } from "react-i18next";

const Page = () => {
    const { t } = useTranslation();

    return (
        <>
            <FaviconHelmet
                title={`${t("error.title")} | ${t("common.appName")}`}
                description={t("error.description")}
            />
            <ErrorPage />
        </>
    );
};

export default Page;
