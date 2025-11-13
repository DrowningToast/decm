import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { BottomNav } from "@/components/BottomNav/BottomNav";

export const InboxNotFoundPage = () => {
    const { t } = useTranslation();

    return (
        <section className="relative z-10 w-full">
            <div className="relative w-full overflow-hidden pb-24">
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

                    {/* Not found state */}
                    <div className="flex flex-col items-start pt-4 px-4">
                        <Typography variant="text" tag="p" color="muted" className="text-base">
                            {t("participant.inbox.noResults", "No item matches your search")}
                        </Typography>
                    </div>
                </div>
            </div>

            {/* Bottom Navigation */}
            <BottomNav variant="search-notification" onBack={() => window.history.back()} />
        </section>
    );
};
