import React from "react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";

interface EventActionsProps {
    type: "pending" | "signed";
    eventId: string;
    onActionClick?: (eventId: string) => void;
}

export const EventActions: React.FC<EventActionsProps> = ({ type, eventId, onActionClick }) => {
    const { t } = useTranslation();

    if (type === "pending") {
        return (
            <Button
                variant="primary"
                size="sm"
                className="py-1 px-3 text-xs font-semibold"
                onClick={() => onActionClick?.(eventId)}
            >
                {t("issuer.sign.action.review")}
            </Button>
        );
    }

    return (
        <Button
            variant="secondary-dark"
            size="sm"
            className="py-1 px-3 text-xs font-semibold"
            onClick={() => onActionClick?.(eventId)}
        >
            {t("issuer.sign.action.view")}
        </Button>
    );
};
