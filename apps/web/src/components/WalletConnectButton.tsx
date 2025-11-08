import { useAppKit } from "@reown/appkit/react";
import type { ClassValue } from "clsx";
import { cn } from "@/lib/utils";

interface Props extends React.PropsWithChildren {
    className?: ClassValue;
    onClick?: () => void;
    isLoading?: boolean;
}

export const WalletConnectButton = ({ className, onClick, children, isLoading = false }: Props) => {
    const { open } = useAppKit();

    const handleClick = () => {
        if (isLoading) {
            return;
        }
        onClick?.();
        open();
    };

    return (
        <button
            type="button"
            className={cn("border-0 bg-transparent p-0 font-inherit text-inherit", className)}
            onClick={handleClick}
            disabled={isLoading}
        >
            {children}
        </button>
    );
};
