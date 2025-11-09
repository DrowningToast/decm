import { Typography } from "@/components/typography/typography";
import { getIssuerStateText, getIssuerStateColor, getIssuerStateBgColor } from "./issuerStateUtils";

interface IssuerStatusBadgeProps {
    isSigned: number;
    className?: string;
}

export function IssuerStatusBadge({ isSigned, className = "" }: IssuerStatusBadgeProps) {
    const statusText = getIssuerStateText(isSigned);
    const textColor = getIssuerStateColor(isSigned);
    const bgColor = getIssuerStateBgColor(isSigned);

    return (
        <span
            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${bgColor} ${textColor} ${className}`}
        >
            {statusText}
        </span>
    );
}
