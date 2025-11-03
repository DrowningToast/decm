import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { BottomNav } from "@/components/BottomNav/BottomNav";
import { useInboxDetailUsecase } from "./useInboxDetailUsecase";
import { useInboxNavStore } from "@/components/BottomNav/stores/inbox";
import { useNavigate } from "@/router";
import { Loader2, AlertCircle } from "lucide-react";

interface InboxDetailPageProps {
    inboxId: string;
}

export const InboxDetailPage: React.FC<InboxDetailPageProps> = ({ inboxId }) => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { inboxDetail, isLoading, error } = useInboxDetailUsecase({ inboxId });
    const { setOnViewEventCallback, setOnViewCertificateCallback } = useInboxNavStore();

    // Set up view event callback
    useEffect(() => {
        if (inboxDetail?.eventId) {
            setOnViewEventCallback(() => {
                navigate(`/app/events/:id`, { params: { id: inboxDetail.eventId! } });
            });
        }
    }, [inboxDetail?.eventId, setOnViewEventCallback, navigate]);

    // Set up view certificate callback
    useEffect(() => {
        if (inboxDetail?.certificateId) {
            setOnViewCertificateCallback(() => {
                navigate(`/app/certificates/:id`, { params: { id: inboxDetail.certificateId! } });
            });
        }
    }, [inboxDetail?.certificateId, setOnViewCertificateCallback, navigate]);

    // Loading state
    if (isLoading) {
        return (
            <section className="relative z-10 w-full">
                <div className="relative min-h-screen w-full flex items-center justify-center">
                    <div className="flex flex-col items-center gap-y-4">
                        <Loader2 className="w-12 h-12 text-primary animate-spin" />
                        <Typography variant="text" tag="p" color="muted" className="text-base">
                            {t("common.loading", "Loading...")}
                        </Typography>
                    </div>
                </div>
            </section>
        );
    }

    // Error state
    if (error || !inboxDetail) {
        return (
            <section className="relative z-10 w-full">
                <div className="relative min-h-screen w-full flex items-center justify-center px-4">
                    <div className="flex flex-col items-center gap-y-4 max-w-sm text-center">
                        {/* Error Icon */}
                        <div className="rounded-full bg-destructive/10 p-4 md:p-6">
                            <AlertCircle className="w-8 h-8 md:w-12 md:h-12 text-destructive" />
                        </div>

                        {/* Error Content */}
                        <div className="flex flex-col gap-y-2">
                            <Typography
                                variant="text"
                                tag="h3"
                                className="text-lg md:text-xl font-semibold"
                            >
                                {t("errors.inbox.title", "Something went wrong")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted"
                                className="text-sm md:text-base"
                            >
                                {t(
                                    "errors.inbox.description",
                                    "We couldn't load this message. Please try again later.",
                                )}
                            </Typography>
                        </div>
                    </div>
                </div>
                <BottomNav variant="search-notification" onBack={() => window.history.back()} />
            </section>
        );
    }

    const getBottomNavVariant = ():
        | "inbox-view-event"
        | "inbox-view-certificate"
        | "inbox-missing-event"
        | "search-notification" => {
        if (inboxDetail.contentType === "event-invitation") {
            return "inbox-view-event";
        } else if (inboxDetail.contentType === "certificate") {
            return inboxDetail.isUserInEvent ? "inbox-view-certificate" : "inbox-missing-event";
        }
        return "search-notification";
    };

    return (
        <section className="relative z-10 w-full">
            <div className="relative min-h-screen w-full overflow-hidden pb-24">
                {/* Background image */}
                <div className="absolute bottom-0 right-0 -translate-x-1/2 aspect-auto w-[229px] h-auto opacity-30 pointer-events-none">
                    <img
                        src="/assets/mailbox.webp"
                        alt=""
                        className="w-full h-full object-cover object-center"
                    />
                </div>

                {/* Main content */}
                <div className="relative z-10 w-full px-4 py-4 flex flex-col gap-y-6">
                    {/* Header section */}
                    <div className="flex flex-col gap-2">
                        <Typography
                            variant="header"
                            tag="h1"
                            color="primary"
                            className="text-[28px]/[34px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] font-header"
                        >
                            {inboxDetail.title}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base/[22px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.inbox.from", "From")}: {inboxDetail.sender}
                        </Typography>
                    </div>

                    {/* Content section */}
                    <div className="flex flex-col gap-4">
                        {/* Date info */}
                        <div className="flex items-center justify-between px-0">
                            <Typography variant="text" tag="p" color="muted" className="text-xs">
                                {t("participant.inbox.date", "Date")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-xs [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {inboxDetail.date}
                            </Typography>
                        </div>

                        {/* Divider */}
                        <div className="h-px bg-border" />

                        {/* Status info */}
                        <div className="flex items-center justify-between px-0">
                            <Typography variant="text" tag="p" color="muted" className="text-xs">
                                {t("participant.inbox.statusLabel", "Status")}
                            </Typography>
                            <Typography variant="text" tag="p" className="text-xs capitalize">
                                {inboxDetail.status}
                            </Typography>
                        </div>

                        {/* Divider */}
                        <div className="h-px bg-border" />

                        {/* Description */}
                        <Typography
                            variant="text"
                            tag="p"
                            color="foreground"
                            className="text-base leading-relaxed pt-4"
                        >
                            {inboxDetail.description}
                        </Typography>
                    </div>
                </div>
            </div>

            {/* Bottom Navigation */}
            <BottomNav variant={getBottomNavVariant()} onBack={() => window.history.back()} />
        </section>
    );
};
