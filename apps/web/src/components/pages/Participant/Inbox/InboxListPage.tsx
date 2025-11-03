import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { BottomNav } from "@/components/BottomNav/BottomNav";
import { useInboxListUsecase, type InboxItem } from "./useInboxListUsecase";
import { Link } from "@/router";
import { Inbox } from "lucide-react";

const getStatusColor = (status: InboxItem["status"]) => {
    switch (status) {
        case "pending":
            return "text-muted";
        case "available":
            return "text-success";
        case "expired":
            return "text-error";
        case "action-required":
            return "text-muted";
        default:
            return "text-muted";
    }
};

const getStatusLabel = (status: InboxItem["status"], t: ReturnType<typeof useTranslation>["t"]) => {
    switch (status) {
        case "pending":
            return t("participant.inbox.status.pending", "Pending");
        case "available":
            return t("participant.inbox.status.available", "Available");
        case "expired":
            return t("participant.inbox.status.expired", "Expired");
        case "action-required":
            return t("participant.inbox.status.actionRequired", "Action required");
        default:
            return "";
    }
};

export const InboxListPage = () => {
    const { t } = useTranslation();
    const { inboxItems } = useInboxListUsecase();

    const hasItems = inboxItems.length > 0;

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
                <div className="relative z-10 w-full px-4 py-4 flex flex-col gap-y-4">
                    {/* Header section */}
                    <div className="flex flex-col gap-1.5">
                        <Typography
                            variant="header"
                            tag="h1"
                            color="primary"
                            className="text-[28px]/[34px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] font-header"
                        >
                            {t("participant.inbox.title", "Inbox")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base/[22px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t(
                                "participant.inbox.subtitle",
                                "Invitation or pending certificate appear here",
                            )}
                        </Typography>
                    </div>

                    {/* Content area */}
                    <div className="flex flex-col gap-y-2.5 px-0">
                        {/* Table header */}
                        <div className="flex items-center justify-between px-4">
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted"
                                className="text-xs font-medium"
                            >
                                {t("participant.inbox.table.actions", "Actions")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted"
                                className="text-xs font-medium"
                            >
                                {t("participant.inbox.table.dateTime", "Date time")}
                            </Typography>
                        </div>

                        {/* Divider */}
                        <div className="h-px bg-border" />
                    </div>

                    {/* Content based on state */}
                    {hasItems ? (
                        <div className="flex flex-col gap-y-4 px-0">
                            {inboxItems.map((item) => (
                                <InboxItemComponent key={item.id} item={item} />
                            ))}
                        </div>
                    ) : (
                        <div className="flex flex-col items-center justify-center gap-y-4 py-12 md:py-16">
                            {/* Icon */}
                            <div className="rounded-full bg-muted/30 p-4 md:p-6">
                                <Inbox className="w-8 h-8 md:w-12 md:h-12 text-muted-foreground" />
                            </div>

                            {/* Content */}
                            <div className="text-center flex flex-col gap-y-2 max-w-sm px-4">
                                <Typography
                                    variant="text"
                                    tag="h3"
                                    className="text-lg md:text-xl font-semibold"
                                >
                                    {t("participant.inbox.empty.title", "Your inbox is empty")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="muted"
                                    className="text-sm md:text-base"
                                >
                                    {t(
                                        "participant.inbox.empty.description",
                                        "Invitations and certificates will appear here when they arrive",
                                    )}
                                </Typography>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Bottom Navigation */}
            <BottomNav variant="search-notification" onBack={() => window.history.back()} />
        </section>
    );
};

const InboxItemComponent = ({ item }: { item: InboxItem }) => {
    const { t } = useTranslation();

    return (
        <Link
            to="/app/inbox/:id"
            params={{ id: item.id }}
            className="flex flex-col gap-1 px-4 cursor-pointer hover:opacity-80 transition-opacity group text-left"
        >
            {/* Title and Date row */}
            <div className="flex items-center gap-3 justify-between">
                <Typography
                    variant="text"
                    tag="p"
                    className="text-base font-normal underline group-hover:text-primary transition-colors"
                >
                    {item.title}
                </Typography>
                <Typography
                    variant="text"
                    tag="span"
                    className="text-xs whitespace-nowrap [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                >
                    {item.date}
                </Typography>
            </div>

            {/* Sender and Status row */}
            <div className="flex items-center justify-between">
                <Typography variant="text" tag="p" color="muted" className="text-sm">
                    {item.sender}
                </Typography>
                <Typography
                    variant="text"
                    tag="p"
                    className={`text-sm ${getStatusColor(item.status)}`}
                >
                    {getStatusLabel(item.status, t)}
                </Typography>
            </div>
        </Link>
    );
};
