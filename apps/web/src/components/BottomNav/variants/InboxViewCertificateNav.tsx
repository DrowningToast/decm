import { ChevronLeft, ArrowRight } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { Typography } from "@/components/typography/typography";
import { useInboxNavStore } from "../stores/inbox";

interface InboxViewCertificateNavProps {
    className?: string;
}

export const InboxViewCertificateNav = ({
    className: propClassName,
}: InboxViewCertificateNavProps) => {
    const { onBack, className: contextClassName } = useBottomContainerContext();
    const { t } = useTranslation();
    const { onViewCertificateCallback } = useInboxNavStore();

    const handleViewCertificate = () => {
        if (onViewCertificateCallback) {
            onViewCertificateCallback();
        }
    };

    return (
        <div
            className={cn(
                contextClassName,
                propClassName,
                "flex items-center gap-1.5 h-13 bg-primary rounded-xl p-1.5",
            )}
        >
            {/* Back Button */}
            <button
                onClick={onBack}
                className="cursor-pointer flex items-center justify-center w-10 h-10 bg-primary rounded-lg hover:bg-primary/90 transition-colors flex-shrink-0"
                aria-label="Go back"
            >
                <ChevronLeft className="w-5 h-5 text-white" />
            </button>

            {/* View Certificate Button */}
            <button
                onClick={handleViewCertificate}
                className="flex-1 h-10 bg-white rounded-lg flex items-center justify-center px-4 hover:bg-white/90 transition-colors cursor-pointer"
                aria-label="View certificate"
            >
                <div className="flex items-center gap-2 justify-center">
                    <Typography
                        variant="text"
                        tag="span"
                        color="background-alt"
                        className="text-xs font-normal leading-normal tracking-[0.06px] whitespace-nowrap"
                    >
                        {t("participant.inbox.viewCertificate", "View Certificate")}
                    </Typography>
                    <ArrowRight className="w-4 h-4 text-background-alt" />
                </div>
            </button>
        </div>
    );
};
