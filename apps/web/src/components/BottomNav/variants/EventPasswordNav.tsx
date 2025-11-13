import { ChevronLeft, Send } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useEventPasswordNavStore } from "../stores/event-password";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";

interface EventPasswordNavProps {
    className?: string;
}

export const EventPasswordNav = ({ className: propClassName }: EventPasswordNavProps) => {
    const { password, setPassword, onSubmitCallback } = useEventPasswordNavStore();

    const { onBack, className: contextClassName } = useBottomContainerContext();
    const { t } = useTranslation();

    const onSubmit = () => {
        if (onSubmitCallback) {
            onSubmitCallback(password);
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

            {/* Password Input */}
            <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder={t("participant.events.passwordPlaceholder") || "Password here"}
                className="flex-1 h-10 px-4 rounded-lg bg-white text-sm text-primary placeholder-primary/60 border border-border focus:outline-none focus:ring-2 focus:ring-primary/50"
            />

            {/* Submit Button */}
            <button
                onClick={onSubmit}
                className="cursor-pointer flex items-center justify-center w-10 h-10 bg-white rounded-lg hover:bg-white/90 transition-colors flex-shrink-0"
                aria-label="Submit password"
            >
                <Send className="w-5 h-5 text-primary" />
            </button>
        </div>
    );
};
