import { ChevronLeft } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { Typography } from "@/components/typography/typography";

export const ShortlistedNav = () => {
    const { onBack, className } = useBottomContainerContext();
    const { t } = useTranslation();

    return (
        <div
            className={cn(className, "flex items-center gap-1.5 h-13 bg-primary rounded-xl p-1.5")}
        >
            {/* Back Button */}
            <button
                onClick={onBack}
                className="cursor-pointer flex items-center justify-center w-10 h-10 bg-primary rounded-lg hover:bg-primary/90 transition-colors flex-shrink-0"
                aria-label="Go back"
            >
                <ChevronLeft className="w-5 h-5 text-white" />
            </button>

            {/* Message */}
            <Typography
                variant="text"
                tag="p"
                color="foreground"
                className="flex-1 text-xs font-normal italic leading-normal tracking-[0.06px] text-center"
            >
                {t("participant.events.shortlisted") || "Your name is shortlisted!"}
            </Typography>
        </div>
    );
};
