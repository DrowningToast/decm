import { useState } from "react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { useSignup } from "./useSignup";

export interface ConfirmationItem {
    id: string;
    message: string;
}

interface BaseConfirmPageProps extends React.PropsWithChildren {
    /** The title text displayed at the top */
    title: string;
    /** Array of confirmation items to display as checkboxes */
    requiredConfirmations: ConfirmationItem[];
    /** Text for the confirm button */
    confirmButtonText: string;
    /** Text for the back button (desktop only) */
    backButtonText: string;
    /** Callback when confirm button is clicked */
    onConfirm: () => void;
    /** Callback when back button is clicked */
    onBack: () => void;
    /** Whether the confirm button should be disabled if not all checkboxes are checked */
    requireAllChecked?: boolean;
}

export const BaseConfirmPage: React.FC<BaseConfirmPageProps> = ({
    title,
    requiredConfirmations,
    confirmButtonText,
    backButtonText,
    onConfirm,
    onBack,
    requireAllChecked = true,
    children,
}) => {

    const { isLoading } = useSignup()

    const [checkedItems, setCheckedItems] = useState<Record<string, boolean>>(
        requiredConfirmations.reduce((acc, item) => ({ ...acc, [item.id]: false }), {})
    );

    const handleCheckChange = (id: string, checked: boolean) => {
        setCheckedItems((prev) => ({ ...prev, [id]: checked }));
    };

    const allChecked = Object.values(checkedItems).every((checked) => checked);
    const isConfirmDisabled = requireAllChecked && !allChecked;

    return (
        <div className="min-h-screen bg-[#e9dede] flex flex-col items-center px-6 py-16 md:py-24 relative overflow-hidden">
            {/* Background decorative image - positioned absolutely */}
            <div className="absolute inset-0 pointer-events-none opacity-30 overflow-hidden">
                <div className="absolute top-[53%] left-1/2 -translate-x-1/2 w-[96%] md:w-[31%] h-auto">
                    {/* Lady Justice placeholder - using background color for now */}
                    <div className="w-full aspect-[3/4] bg-gradient-to-b from-transparent via-muted/20 to-transparent rounded-full blur-3xl" />
                </div>
            </div>

            {/* Main content - positioned relative to stay above background */}
            <div className="relative z-10 w-full max-w-[420px] space-y-4">
                {/* Header Section */}
                <div className="space-y-1.5">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="primary"
                        className="text-[36px] leading-[40px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] tracking-[0.06px]"
                    >
                        {title}
                    </Typography>

                    {/* Description - children slot */}
                    <div className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] text-background-alt">
                        {children}
                    </div>
                </div>

                {/* Confirmation Checkboxes */}
                <div className="space-y-3 md:space-y-3">
                    {requiredConfirmations.map((item) => (
                        <div key={item.id} className="flex gap-2 items-start">
                            <Checkbox
                                id={item.id}
                                checked={checkedItems[item.id]}
                                onCheckedChange={(checked) =>
                                    handleCheckChange(item.id, checked === true)
                                }
                                className="mt-0.5 size-4 rounded-[2px] data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                            />
                            <label
                                htmlFor={item.id}
                                className="text-base leading-normal text-background tracking-[0.06px] cursor-pointer select-none flex-1"
                            >
                                {item.message}
                            </label>
                        </div>
                    ))}
                </div>

                {/* Action Buttons Section */}
                <div className="space-y-3 pt-2 md:pt-4">
                    {/* Confirm Button */}
                    <Button
                        type="button"
                        onClick={onConfirm}
                        disabled={isConfirmDisabled || isLoading}
                        variant="primary"
                        size="xl"
                        className="w-full"
                    >
                        {confirmButtonText}
                    </Button>

                    {/* Back Button - Only visible on desktop */}
                    <Button
                        type="button"
                        onClick={onBack}
                        variant="secondary-light"
                        size="xl"
                        className="w-full hidden md:flex"
                    >
                        {backButtonText}
                    </Button>
                </div>
            </div>
        </div>
    );
};