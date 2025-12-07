import { getIssuerStateText, getIssuerStateColor, getIssuerStateBgColor } from "./issuerStateUtils";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";

interface IssuerStatusBadgeProps {
    isSigned: boolean;
    className?: string;
}

export function IssuerStatusBadge({ isSigned, className = "" }: IssuerStatusBadgeProps) {
    const { t } = useTranslation();
    const statusText = getIssuerStateText(isSigned, t);
    const textColor = getIssuerStateColor(isSigned);
    const bgColor = getIssuerStateBgColor(isSigned);

    return (
        <Typography
            variant="text"
            tag="span"
            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${bgColor} ${textColor} ${className}`}
        >
            {statusText}
        </Typography>
    );
}
