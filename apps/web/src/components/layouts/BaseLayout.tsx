import type { ClassValue } from "clsx";
import { cn } from "@/lib/utils";

interface BaseLayoutProps extends React.PropsWithChildren {
    children: React.ReactNode;
    className?: ClassValue;
    variant?: "light" | "dark";
}

export const BaseLayout: React.FC<BaseLayoutProps> = ({
    children,
    className,
    variant = "dark",
}) => {
    return (
        <div
            className={cn("pt-12 md:pt-[60px] box-content", className, {
                "bg-foreground-alt text-background-alt": variant === "light",
                "bg-background text-foreground": variant === "dark",
            })}
        >
            {children}
        </div>
    );
};
