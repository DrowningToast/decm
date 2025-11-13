import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { EventList } from "./EventList";
import { useEventsListUsecase } from "./useEventsListUsecase";
import { EventFilters } from "./EventFilters";
import { BottomNav } from "@/components/BottomNav/BottomNav";

type EventFilterType = "all" | "my-events";

export const ParticipantEventListPage = () => {
    const { t } = useTranslation();
    const [filterType, setFilterType] = useState<EventFilterType>("all");
    const { events, isLoading } = useEventsListUsecase({
        filterType,
    });

    const handleFilterChange = (filter: EventFilterType) => {
        setFilterType(filter);
    };

    return (
        <section className="relative z-10 w-full">
            <div className="relative w-full overflow-hidden">
                {/* Main content */}
                <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 pb-4 md:pb-24 flex flex-col gap-y-4 md:gap-y-6">
                    {/* Header section */}
                    <div className="flex flex-col gap-1.5">
                        <Typography
                            variant="header"
                            tag="h1"
                            color="primary"
                            className="text-[28px]/[34px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] font-header"
                        >
                            {t("participant.events.title")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base/[22px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.events.description")}
                        </Typography>
                    </div>

                    {/* Event List */}
                    <EventList
                        events={events ?? []}
                        isLoading={isLoading}
                        filterType={filterType}
                    />
                </div>

                {/* Filters section - Fixed above BottomNav */}
                <div className="md:hidden fixed bottom-[88px] left-0 right-0 px-4 pb-2 bg-gradient-to-t from-background via-background to-transparent z-40">
                    <EventFilters onFilterChange={handleFilterChange} />
                </div>
                <div className="hidden md:flex fixed bottom-[120px] left-1/2 transform -translate-x-1/2 z-40">
                    <div className="w-[343px]">
                        <EventFilters onFilterChange={handleFilterChange} />
                    </div>
                </div>

                {/* Bottom Navigation */}
                <BottomNav variant="search-event" onBack={() => window.history.back()} />
            </div>
        </section>
    );
};
