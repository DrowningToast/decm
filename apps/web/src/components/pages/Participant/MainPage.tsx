import { useAuth } from "@/context/AuthContext";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";
import { ServiceList } from "./components/ServiceList";
import { participantServices, generalServices } from "./config";

export const MainPage = () => {
    const { user } = useAuth();
    const { t } = useTranslation();

    if (!user) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <Typography variant="text" tag="p" color="muted-foreground">
                    {t("common.loading")}
                </Typography>
            </div>
        );
    }

    // Get greeting based on time of day
    const getGreeting = () => {
        const hour = new Date().getHours();
        if (hour < 12) return t("participant.home.greeting.morning");
        if (hour < 18) return t("participant.home.greeting.afternoon");
        return t("participant.home.greeting.evening");
    };

    // Get status message based on unread messages
    const getStatusMessage = () => {
        const hasUnreadMessages = user.unreadInboxMessageCount > 0;
        return hasUnreadMessages
            ? t("participant.home.unreadMessages")
            : t("participant.home.noRequiredAction");
    };

    // Map services with translated labels
    const participantServicesWithLabels = participantServices.map((service) => ({
        ...service,
        label: t(service.translationKey),
    }));

    const generalServicesWithLabels = generalServices.map((service) => ({
        ...service,
        label: t(service.translationKey),
    }));

    return (
        <div className="relative w-full overflow-hidden">
            {/* Main content */}
            <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 flex flex-col gap-y-8">
                {/* Greeting section */}
                <div className="flex flex-col gap-1">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="primary"
                        className="text-4xl md:text-[56px] leading-[40px] md:leading-[84px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] font-header"
                    >
                        {getGreeting()}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted-foreground"
                        className="text-base md:text-lg mt-2"
                    >
                        {getStatusMessage()}
                    </Typography>
                </div>

                {/* Participant Services Section */}
                <section className="flex flex-col gap-y-2 md:gap-y-6">
                    <Typography
                        variant="text"
                        tag="h2"
                        color="foreground-alt"
                        className="text-base md:text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("participant.home.participantServices")}
                    </Typography>

                    <ServiceList services={participantServicesWithLabels} layout="grid" />
                </section>

                {/* General Services Section */}
                <section className="flex flex-col gap-y-2">
                    <Typography
                        variant="text"
                        tag="h2"
                        color="muted"
                        className="text-base md:text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("participant.home.generalServices")}
                    </Typography>

                    <ServiceList services={generalServicesWithLabels} layout="single" />
                </section>
            </div>
        </div>
    );
};
