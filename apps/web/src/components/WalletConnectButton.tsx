import { useAppKit } from "@reown/appkit/react";
import type { ClassValue } from "clsx";
import { cn } from "@/lib/utils";
import { Button, buttonVariants } from "./ui/button";
import type { VariantProps } from "class-variance-authority";

interface Props extends React.PropsWithChildren, VariantProps<typeof buttonVariants> {
    className?: ClassValue;
    onClick?: () => void;
    isLoading?: boolean;
}

export const WalletConnectButton = ({
    className,
    onClick,
    children,
    isLoading = false,
    variant = "primary",
}: Props) => {
    const { open } = useAppKit();

    const handleClick = () => {
        if (isLoading) {
            return;
        }
        onClick?.();
        open();
    };

    return (
        <Button
            type="button"
            className={cn(buttonVariants({ variant, className }))}
            onClick={handleClick}
            disabled={isLoading}
        >
            {children}
        </Button>
    );
};
