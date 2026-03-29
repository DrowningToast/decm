import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import { Typography } from "@/components/typography/typography";
import { Link, type Params, type Path } from "@/router";
import { useTranslation } from "react-i18next";

export default function HostHomePage() {
    const { t } = useTranslation();

    // Get greeting based on time of day
    const getGreeting = () => {
        const hour = new Date().getHours();
        if (hour < 12) return t("host.home.greeting.morning");
        if (hour < 18) return t("host.home.greeting.afternoon");
        return t("host.home.greeting.evening");
    };

    const actionCount = 0; // TODO: Replace with actual action count

    return (
        <div className="flex flex-col gap-y-4" title={t("host.home.dashboardTitle")}>
            <img
                src="/justice.png"
                alt={t("host.home.dashboardAlt")}
                className="absolute bottom-0 lg:right-0 w-[400px] h-[400px] opacity-50 m-0"
            />

            <SectionContainer>
                <TitleSubtitle
                    title={getGreeting()}
                    subtitle={t("host.home.actionRequired", { count: actionCount })}
                />
            </SectionContainer>

            <SectionContainer className="space-y-2">
                <Typography size={"base"} tag="p" color="muted">
                    {t("host.home.hostServices")}
                </Typography>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-2 lg:gap-6 lg:mt-6">
                    <MenuItem title={t("host.home.services.event")} to="/host/events" />
                </div>
            </SectionContainer>

            <SectionContainer>
                <Typography size={"base"} tag="p" color="muted">
                    {t("host.home.generalServices")}
                </Typography>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-2 lg:gap-6 lg:mt-6">
                    <MenuItem title={t("host.home.services.verifyCertificate")} to="/host/home" />
                </div>
            </SectionContainer>
        </div>
    );
}

interface MenuItemProps {
    title: string;
    to: Exclude<Path, keyof Params>;
}
function MenuItem({ title, to }: MenuItemProps) {
    return (
        <Link
            to={to}
            className="lg:py-5 lg:px-6 lg:border lg:border-[#D9D9D91A] lg:rounded-lg lg:bg-[#D9D9D905]"
        >
            <Typography
                size={"header"}
                tag="p"
                color="secondary"
                className="underline lg:no-underline lg:text-2xl font-cormorant"
            >
                {title}
            </Typography>
        </Link>
    );
}
