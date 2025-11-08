import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

interface PageContainerProps {
    children: ReactNode;
    className?: string;
}

export const PageContainer: React.FC<PageContainerProps> = ({
    children,
    className,
}: PageContainerProps) => {
    return (
        <div className={cn("max-w-[1440px] py-4 md:py-14 box-content mx-auto", className)}>
            {children}
        </div>
    );
};
