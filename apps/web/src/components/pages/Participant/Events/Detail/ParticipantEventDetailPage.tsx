import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { useEventViewModelUsecase } from "./useEventViewModelUsecase";
import { ImageLoader } from "@/components/common/ImageLoader";
import { addr } from "@/components/common/addr";
import { EthExplorerLink } from "@/components/common/EthscanLink";
import { GoogleMapsEmbed } from "@/components/ui/google-maps-embed";
import { ActionMenu } from "./ActionMenu";

interface EventDetailPageProps {
    eventId: string;
}

export const EventDetailPage: React.FC<EventDetailPageProps> = ({ eventId }) => {
    const { t, i18n } = useTranslation();
    const { event, isLoading, error } = useEventViewModelUsecase({ eventId });

    const eventType = useMemo(() => {
        switch (event?.eventType) {
            case "invite": {
                return t("participant.events.detail.inviteOnly");
            }
            case "private": {
                return t("participant.events.detail.passwordRequired");
            }
        }
        return t("participant.events.detail.open");
    }, [event?.eventType, t]);

    // Use i18n language for date formatting (e.g., 'en' or 'th')
    const locale = i18n.language === "th" ? "th-TH" : "en-US";
    const formattedDate = event?.finalCallDate
        ? new Date(event?.finalCallDate).toLocaleDateString(locale, {
              year: "numeric",
              month: "short",
              day: "numeric",
          })
        : "";

    const isShowFinalCallDate = event?.finalCallDate && event?.eventStatus !== "closed";

    const eventStatusText = useMemo(() => {
        switch (event?.eventStatus) {
            case "active":
                return t("participant.events.detail.active");
            case "inactive":
                return t("participant.events.detail.inactive");
        }
        return t("participant.events.detail.closed");
    }, [event?.eventStatus, t]);

    // Loading state
    if (isLoading) {
        return (
            <section className="relative z-10 w-full">
                <div className="relative min-h-screen w-full flex items-center justify-center">
                    <Typography variant="text" tag="p" color="muted" className="animate-pulse">
                        {t("common.loading")}
                    </Typography>
                </div>
            </section>
        );
    }

    // Error state
    if (error || !event) {
        return (
            <section className="relative z-10 w-full">
                <div className="relative min-h-screen w-full flex items-center justify-center">
                    <Typography variant="text" tag="p" color="destructive">
                        {t("errors.generic")}
                    </Typography>
                </div>
            </section>
        );
    }

    return (
        <div className="relative w-full md:grid md:grid-cols-2 lg:grid-cols-3 gap-x-2 md:gap-x-8 px-4 md:px-16">
            {/* Main content */}
            <div className="lg:col-span-2 relative z-10 w-full max-w-[1384px] mx-auto flex flex-col gap-y-4 md:gap-y-6 pb-8">
                {/* Header section with logo and title */}
                <div className="flex flex-col gap-y-4">
                    <div className="flex items-start gap-3 md:gap-4">
                        <ImageLoader
                            src={event.iconPresignedUrl || ""}
                            alt={event.title || ""}
                            className="w-8 h-8 md:w-10 md:h-10 bg-muted rounded flex-shrink-0"
                        />
                        <Typography
                            variant="header"
                            tag="h1"
                            color="foreground"
                            className="text-[28px]/[34px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] font-header"
                        >
                            {event.title}
                        </Typography>
                    </div>

                    {/* Short Description */}
                    {event.shortDescription && (
                        <div className="flex flex-col gap-y-2">
                            <Typography
                                variant="text"
                                tag="p"
                                color="foreground"
                                className="text-base/[22px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {event.shortDescription}
                            </Typography>
                        </div>
                    )}
                </div>

                <ImageLoader
                    src={event.bannerPresignedUrl || ""}
                    alt={event.title || ""}
                    className="w-full min-h-[172px] md:min-h-[300px] bg-muted rounded-lg flex-shrink-0"
                />

                {/* Long Description */}
                {event.longDescription && (
                    <div className="flex flex-col gap-y-2">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] font-medium"
                        >
                            {t("participant.events.detail.longDescription")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="foreground"
                            className="text-base/[22px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] whitespace-pre-wrap"
                        >
                            {event.longDescription}
                        </Typography>
                    </div>
                )}

                {/* Event Details */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                    {/* Status */}
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.events.detail.status")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {eventStatusText}
                        </Typography>
                    </div>

                    {/* Final Call Date */}
                    {isShowFinalCallDate && (
                        <div className="flex flex-col gap-y-1">
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted"
                                className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {t("participant.events.detail.finalCallForRequest")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {formattedDate}
                            </Typography>
                        </div>
                    )}
                    {/* Participation Request */}
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.events.detail.participationRequest")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {eventType}
                        </Typography>
                    </div>

                    {/* Seats Count */}
                    {event.maxAttendees && (
                        <div className="flex flex-col gap-y-1">
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted"
                                className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {t("participant.events.detail.seatsCount")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {`${event.attendeesCount ?? 0} / ${event.maxAttendees}`}
                            </Typography>
                        </div>
                    )}

                    {/* Location */}
                    {event.location && (
                        <div className="flex flex-col gap-y-1">
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted"
                                className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {t("participant.events.detail.location")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {event.location}
                            </Typography>
                        </div>
                    )}

                    {/* Event Contract Address */}
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.events.detail.eventContractAddress")}
                        </Typography>
                        <EthExplorerLink address={event.eventContractAddress || ""}>
                            <Typography
                                variant="text"
                                tag="p"
                                className="underline text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {addr(event.eventContractAddress || "") ||
                                    t("participant.events.detail.noContractAddress")}
                            </Typography>
                        </EthExplorerLink>
                    </div>
                </div>

                {/* Google Map */}
                {event.googleMapQuery && (
                    <div className="flex flex-col gap-y-2">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] font-medium"
                        >
                            {t("participant.events.detail.googleMap")}
                        </Typography>
                        <GoogleMapsEmbed
                            query={event.googleMapQuery}
                            height="300px"
                            className="md:h-[400px]"
                        />
                    </div>
                )}
            </div>
            <ActionMenu eventId={eventId} />
        </div>
    );
};
