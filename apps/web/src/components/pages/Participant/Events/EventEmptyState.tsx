import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Calendar } from "lucide-react";

export const EventEmptyState = () => {
    const { t } = useTranslation();

    return (
        <div className="flex flex-col items-center justify-center gap-y-4 py-12 md:py-24">
            {/* Icon */}
            <div className="rounded-full bg-muted/30 p-4 md:p-6">
                <Calendar className="w-8 h-8 md:w-12 md:h-12 text-muted-foreground" />
            </div>

            {/* Content */}
            <div className="text-center flex flex-col gap-y-2">
                <Typography variant="text" tag="h3" className="text-lg md:text-xl font-semibold">
                    {t("participant.events.empty.title", "No events yet")}
                </Typography>
                <Typography variant="text" tag="p" color="muted" className="text-sm md:text-base">
                    {t("participant.events.empty.description", "Events you join will appear here")}
                </Typography>
            </div>
        </div>
    );
};
