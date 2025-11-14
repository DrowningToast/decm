import { ChevronLeft } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useSearchIdentitiesNavStore } from "../stores/identities";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";

interface SearchIdentitiesNavProps {
    className?: string;
}

export const SearchIdentitiesNav = ({ className: propClassName }: SearchIdentitiesNavProps) => {
    const { searchQuery, setSearchQuery } = useSearchIdentitiesNavStore();

    const { onBack, className: contextClassName } = useBottomContainerContext();
    const { t } = useTranslation();

    return (
        <div
            className={cn(
                contextClassName,
                propClassName,
                `flex items-center gap-1.5 h-13 bg-primary rounded-xl p-1.5`,
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

            {/* Search Input */}
            <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder={
                    t("participant.identities.searchPlaceholder") || "Search for identities"
                }
                className="flex-1 h-10 px-4 rounded-lg bg-white text-sm text-primary placeholder-primary/60 border border-border focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
        </div>
    );
};
