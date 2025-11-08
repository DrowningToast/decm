import { useTranslation } from "react-i18next";
import {
    StyledTabs,
    StyledTabsList,
    StyledTabsTrigger,
    StyledTabsContent,
} from "@/components/ui/styled-tabs";
import { motion } from "motion/react";

type EventFilterType = "all" | "my-events";

interface EventFiltersProps {
    onFilterChange?: (filterType: EventFilterType) => void;
}

export const EventFilters = ({ onFilterChange }: EventFiltersProps) => {
    const { t } = useTranslation();

    const handleTabChange = (value: string) => {
        onFilterChange?.(value as EventFilterType);
    };

    return (
        <StyledTabs defaultValue="all" onValueChange={handleTabChange} className="h-auto">
            <StyledTabsList className="p-1.5 h-auto">
                <StyledTabsTrigger value="all" className="py-[9px]">
                    <motion.span layout className="text-base/[22px] font-semibold">
                        {t("participant.events.filters.allEvents")}
                    </motion.span>
                </StyledTabsTrigger>
                <StyledTabsTrigger value="my-events" className="py-[9px]">
                    <motion.span layout className="text-base/[22px] font-semibold">
                        {t("participant.events.filters.myEvents")}
                    </motion.span>
                </StyledTabsTrigger>
            </StyledTabsList>
            <StyledTabsContent value="all" className="mt-0" />
            <StyledTabsContent value="my-events" className="mt-0" />
        </StyledTabs>
    );
};
