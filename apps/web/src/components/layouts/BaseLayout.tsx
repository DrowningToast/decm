import type { ClassValue } from "clsx";
import { cn } from "@/lib/utils";

interface BaseLayoutProps extends React.PropsWithChildren {
    children: React.ReactNode;
    className: ClassValue
    variant: "light" | "dark";
}

export const BaseLayout = ({ children, className, variant = "dark" }: BaseLayoutProps) => {
    return (
        <div className={cn("pt-12 md:pt-[60px]", className, {
            "bg-foreground-alt text-background-alt": variant === "dark",
            "bg-background text-foreground": variant === "light",
        })}>
            {children}
        </div>
    )
}