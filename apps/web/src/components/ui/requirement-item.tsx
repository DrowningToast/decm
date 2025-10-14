import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { CheckCircle2Icon, XCircleIcon, MinusCircleIcon } from "lucide-react";

export type RequirementStatus = "required" | "optional" | "not_required";

interface RequirementItemProps {
    /**
     * Label text for the requirement field
     */
    label: string;
    /**
     * Status of the requirement: "required", "optional", or "not_required"
     */
    status: RequirementStatus;
    /**
     * Optional custom className for the container
     */
    className?: string;
}

/**
 * RequirementItem - A reusable component to display requirement status
 *
 * Displays a requirement field with visual indicators for three states:
 * - Required: Green checkmark + "Required" text
 * - Optional: Gray minus icon + "Optional" text
 * - Not Required: Red X icon + "Not Required" text
 *
 * @example
 * ```tsx
 * <RequirementItem label="First Name" status="required" />
 * <RequirementItem label="Bio" status="optional" />
 * <RequirementItem label="Address" status="not_required" />
 * ```
 */
export function RequirementItem({ label, status, className = "" }: RequirementItemProps) {
    const { t } = useTranslation();

    // Determine styling and content based on status
    const getStatusConfig = () => {
        switch (status) {
            case "required":
                return {
                    icon: <CheckCircle2Icon className="h-4 w-4 text-green-600" />,
                    text: t("common.required"),
                    textColor: "text-green-600",
                    bgColor: "bg-green-50 dark:bg-green-950/20",
                    borderColor: "border-green-200 dark:border-green-800",
                };
            case "optional":
                return {
                    icon: <MinusCircleIcon className="h-4 w-4 text-muted-foreground" />,
                    text: t("common.optional"),
                    textColor: "text-muted-foreground",
                    bgColor: "bg-muted/10",
                    borderColor: "border-[#D9D9D91A]",
                };
            case "not_required":
                return {
                    icon: <XCircleIcon className="h-4 w-4 text-red-600" />,
                    text: t("common.notRequired"),
                    textColor: "text-red-600",
                    bgColor: "bg-red-50 dark:bg-red-950/20",
                    borderColor: "border-red-200 dark:border-red-800",
                };
        }
    };

    const config = getStatusConfig();

    return (
        <div
            className={`flex items-center justify-between p-3 rounded-lg border ${config.borderColor} ${config.bgColor} ${className}`}
        >
            <Typography variant="text" tag="span" className="text-sm text-black">
                {label}
            </Typography>
            <div className="flex items-center gap-2">
                {config.icon}
                <Typography
                    variant="text"
                    tag="span"
                    className={`text-sm font-medium ${config.textColor}`}
                >
                    {config.text}
                </Typography>
            </div>
        </div>
    );
}
