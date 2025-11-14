import React from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";

interface EventStatusBadgeProps {
    status: "pending" | "signed";
}

export const EventStatusBadge: React.FC<EventStatusBadgeProps> = ({ status }) => {
    const { t } = useTranslation();

    if (status === "pending") {
        return (
            <Typography
                variant="text"
                tag="span"
                className="inline-block px-3 py-1 text-xs font-semibold rounded-full text-black bg-yellow-500 uppercase tracking-wide"
            >
                {t("issuer.sign.status.waiting")}
            </Typography>
        );
    }

    return (
        <Typography
            variant="text"
            tag="span"
            className="inline-block px-3 py-1 text-xs font-semibold rounded-full text-white bg-green-500 uppercase tracking-wide"
        >
            {t("issuer.sign.status.signed")}
        </Typography>
    );
};
