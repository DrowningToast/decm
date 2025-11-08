import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Lock, Mail, X, Loader2, SearchX } from "lucide-react";
import { EventEmptyState } from "./EventEmptyState";
import type { Event } from "./useEventsListUsecase";
import { Link } from "@/router";
import { useSearchEventNavStore } from "@/components/BottomNav/stores/events";

interface EventListProps {
    events: Event[];
    isLoading?: boolean;
    filterType?: "all" | "my-events";
}

export const EventList = ({ events = [], isLoading, filterType }: EventListProps) => {
    const { t } = useTranslation();
    const { searchQuery } = useSearchEventNavStore();

    const hasContent = events?.length > 0;
    const isSearching = searchQuery.trim().length > 0;
    const isFiltering = filterType === "my-events";
    const hasActiveFilters = isSearching || isFiltering;

    // Loading state
    if (isLoading) {
        return (
            <div className="flex flex-col items-center justify-center gap-y-4 py-12 md:py-16">
                <Loader2 className="w-12 h-12 text-primary animate-spin" />
                <Typography variant="text" tag="p" color="muted" className="text-base">
                    {t("common.loading", "Loading...")}
                </Typography>
            </div>
        );
    }

    // No results found (search/filter applied but no matches)
    if (!hasContent && hasActiveFilters) {
        return (
            <div className="flex flex-col items-center justify-center gap-y-4 py-12 md:py-16">
                <div className="rounded-full bg-muted/30 p-4 md:p-6">
                    <SearchX className="w-8 h-8 md:w-12 md:h-12 text-muted-foreground" />
                </div>
                <div className="text-center flex flex-col gap-y-2">
                    <Typography
                        variant="text"
                        tag="h3"
                        className="text-lg md:text-xl font-semibold"
                    >
                        {t("participant.events.noResults.title", "No events found")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-sm md:text-base"
                    >
                        {t(
                            "participant.events.noResults.description",
                            "Try adjusting your search or filters",
                        )}
                    </Typography>
                </div>
            </div>
        );
    }

    // Truly empty state (no data at all)
    if (!hasContent) {
        return <EventEmptyState />;
    }

    return (
        <div className="flex flex-col gap-y-4 md:gap-y-6">
            <div className="flex flex-col gap-y-2.5">
                {/* Table Header */}
                <div className="flex items-center justify-between px-4 md:px-0">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-xs md:text-sm font-medium"
                    >
                        {t("participant.events.name")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-xs md:text-sm font-medium"
                    >
                        {t("participant.events.finalCall")}
                    </Typography>
                </div>

                {/* Divider */}
                <div className="h-px bg-border" />
            </div>

            {/* Event items */}
            <div className="flex flex-col gap-y-4 md:gap-y-6">
                {events.map((event) => (
                    <EventItem key={event.id} event={event} />
                ))}
            </div>
        </div>
    );
};

const EventItem = ({ event }: { event: Event }) => {
    const { t } = useTranslation();

    const getAccessIcon = () => {
        if (event.accessType === "password" || event.requiresPassword) {
            return {
                icon: Lock,
                label: t("participant.events.passwordRequired"),
                color: "text-foreground-alt",
            };
        }
        if (event.accessType === "invite-only") {
            return {
                icon: Mail,
                label: t("participant.events.inviteOnly"),
                color: "text-foreground-alt",
            };
        }
        if (event.status === "closed") {
            return {
                icon: X,
                label: t("participant.events.nolongerAccepting"),
                color: "text-destructive",
            };
        }
        return null;
    };

    const accessInfo = getAccessIcon();

    const formattedDate = event.finalCallDate
        ? new Date(event.finalCallDate).toLocaleDateString("en-US", {
              year: "numeric",
              month: "short",
              day: "numeric",
          })
        : "";

    return (
        <Link
            to="/app/events/:id"
            params={{ id: event.id }}
            className="w-full text-left flex flex-col gap-1 px-0 hover:opacity-80 transition-opacity cursor-pointer group"
        >
            {/* Row with name and date */}
            <div className="flex items-center gap-3 md:gap-4 justify-between">
                <Typography
                    variant="text"
                    tag="p"
                    className="text-base md:text-lg font-normal underline group-hover:text-primary transition-colors"
                >
                    {event.name || event.eventName}
                </Typography>

                <Typography
                    variant="text"
                    tag="span"
                    className="text-xs md:text-sm whitespace-nowrap [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                >
                    {formattedDate}
                </Typography>
            </div>

            {/* Access info - below name */}
            {accessInfo && (
                <div className="flex items-center gap-1.5 pl-0">
                    <accessInfo.icon className={`w-5 h-5 flex-shrink-0 ${accessInfo.color}`} />
                    <Typography
                        variant="text"
                        tag="p"
                        color={event.status === "closed" ? "destructive" : "foreground-alt"}
                        className="text-sm"
                    >
                        {accessInfo.label}
                    </Typography>
                </div>
            )}
        </Link>
    );
};
