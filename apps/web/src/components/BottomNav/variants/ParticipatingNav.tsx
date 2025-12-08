import { CircleCheck } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { Typography } from "@/components/typography/typography";

interface ParticipatingNavProps {
    className?: string;
}

export const ParticipatingNav = ({ className: propClassName }: ParticipatingNavProps) => {
    const { className: contextClassName } = useBottomContainerContext();
    const { t } = useTranslation();

    return (
        <div
            className={cn(
                contextClassName,
                propClassName,
                "flex items-center gap-2 h-13 bg-primary rounded-xl px-4",
            )}
        >
            {/* Check Icon */}
            <CircleCheck className="w-5 h-5 text-white flex-shrink-0" />

            {/* Message */}
            <Typography
                variant="text"
                tag="p"
                color="foreground"
                className="flex-1 text-xs font-medium leading-normal tracking-[0.06px]"
            >
                {t("participant.events.participating")}
            </Typography>
        </div>
    );
};
