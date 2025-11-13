import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { ExternalLink } from "lucide-react";
import { BottomNav } from "@/components/BottomNav/BottomNav";
import { useEventViewModelUsecase } from "./useEventViewModelUsecase";
import { useEventPasswordNavStore } from "@/components/BottomNav/stores/event-password";
import { useEventInvitationNavStore } from "@/components/BottomNav/stores/event-invitation";
import { EventStatus, EventType } from "@/services/EventService";
import { Skeleton } from "@/components/ui/skeleton";

interface EventDetailPageProps {
    eventId: string;
}

export const EventDetailPage: React.FC<EventDetailPageProps> = ({ eventId }) => {
    const { t } = useTranslation();
    const { event, isLoading, error, bottomNavVariant } = useEventViewModelUsecase({ eventId });

    const { setOnSubmitCallback, resetPassword } = useEventPasswordNavStore();
    const { setOnAcceptCallback } = useEventInvitationNavStore();

    const eventType = useMemo(() => {
        switch (event?.eventType) {
            case EventType.EventTypeInvite: {
                return t("participant.events.detail.inviteOnly");
            }
            case EventType.EventTypePrivate: {
                return t("participant.events.detail.passwordRequired");
            }
        }
        return t("participant.events.detail.open");
    }, [event?.eventType, t]);

    const formattedDate = event?.finalCallDate
        ? new Date(event?.finalCallDate).toLocaleDateString("en-US", {
              year: "numeric",
              month: "short",
              day: "numeric",
          })
        : "";

    const isShowFinalCallDate =
        event?.finalCallDate && event?.eventStatus !== EventStatus.EventStatusClosed;

    const isClosed = event?.eventStatus === EventStatus.EventStatusClosed;

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
        <div className="relative w-full overflow-hidden">
            {/* Background image */}
            <div className="absolute bottom-0 right-0 w-[424px] h-[424px] md:w-[500px] md:h-[500px] opacity-30 pointer-events-none">
                <img
                    src="/assets/scale.webp"
                    alt=""
                    className="w-full h-full object-cover object-center"
                />
            </div>

            {/* Main content */}
            <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 py-4 md:pt-16 md:pb-24 flex flex-col gap-y-4 md:gap-y-6">
                {/* Header section with logo and title */}
                <div className="flex flex-col gap-y-4">
                    <div className="flex items-start gap-3 md:gap-4">
                        <Skeleton className="w-8 h-8 md:w-10 md:h-10 bg-muted rounded flex-shrink-0">
                            <img
                                src={event.icon_presigned_url}
                                alt={event.title}
                                className="w-full h-full object-cover object-center"
                            />
                        </Skeleton>
                        <Typography
                            variant="header"
                            tag="h1"
                            color="foreground"
                            className="text-[28px]/[34px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] font-header"
                        >
                            {event.title}
                        </Typography>
                    </div>

                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-base/[22px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {event.longDescription}
                    </Typography>
                </div>

                <div>
                    <Skeleton className="w-full h-[172px] md:h-[300px] bg-muted rounded-lg flex-shrink-0">
                        <img
                            src={event.banner_presigned_url}
                            alt={event.title}
                            className="w-full h-full object-cover object-center"
                        />
                    </Skeleton>
                </div>

                {/* Event Details */}
                <div className="flex flex-col gap-y-4">
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
                            {isClosed
                                ? t("participant.events.detail.closed")
                                : event.eventStatus === EventStatus.EventStatusActive
                                  ? t("participant.events.detail.acceptingRequests")
                                  : t("participant.events.detail.inactive")}
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
                                {event.maxAttendees}
                            </Typography>
                        </div>
                    )}

                    {/* Event Contact Address */}
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.events.detail.eventContactAddress")}
                        </Typography>
                        <div className="flex items-center gap-2">
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {event.contactAddress || t("participant.events.detail.noAddress")}
                            </Typography>
                            <ExternalLink className="w-4 h-4 text-foreground-alt" />
                        </div>
                    </div>
                </div>
            </div>

            {/* Bottom Navigation */}
            {bottomNavVariant && (
                <BottomNav variant={bottomNavVariant} onBack={() => window.history.back()} />
            )}
        </div>
    );
};
