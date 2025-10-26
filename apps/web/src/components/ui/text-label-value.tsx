import { Typography } from "@/components/typography/typography";
import { cn } from "@/lib/utils";

interface TextLabelValueProps {
    /**
     * Label text displayed above the value
     */
    label: string;
    /**
     * Value text to display
     */
    value: string;
    /**
     * Optional icon displayed at the end of the value
     */
    endIcon?: React.ReactNode;
    /**
     * Optional custom className for the value text
     */
    valueClassName?: string;
    /**
     * Optional href to make the value a clickable link
     */
    href?: string;
}

/**
 * TextLabelValue - A reusable component for displaying labeled values
 *
 * Displays a label above a value with optional icon and link support.
 * Commonly used in detail pages, forms, and data displays.
 *
 * @example
 * ```tsx
 * <TextLabelValue label="Email" value="user@example.com" />
 * <TextLabelValue label="Website" value="example.com" href="https://example.com" />
 * <TextLabelValue label="Status" value="Active" endIcon={<CheckIcon />} />
 * ```
 */
export function TextLabelValue({
    label,
    value,
    endIcon,
    valueClassName,
    href,
}: TextLabelValueProps) {
    return (
        <div className="flex flex-col gap-1">
            <Typography tag="p" size={"base"} color="muted" className="text-sm">
                {label}
            </Typography>
            <Typography
                tag="p"
                size={"base"}
                className={cn(valueClassName, {
                    "flex items-center gap-2": endIcon,
                })}
            >
                {href ? (
                    <a href={href} target="_blank" rel="noopener noreferrer">
                        {value}
                    </a>
                ) : (
                    value
                )}
                {endIcon}
            </Typography>
        </div>
    );
}
