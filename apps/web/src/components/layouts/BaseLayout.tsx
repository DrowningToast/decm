import type { ClassValue } from "clsx";
import { cn } from "@/lib/utils";
import { cva, type VariantProps } from "class-variance-authority";
import type { PropsWithChildren } from "react";
import type React from "react";

const componentVariants = cva("pt-12 md:pt-[60px]", {
    variants: {
        variant: {
            dark: "bg-foreground-alt text-background-alt",
            light: "bg-background text-foreground",
        },
    },
    defaultVariants: {
        variant: "dark",
    },
});

type BaseLayoutVariants = VariantProps<typeof componentVariants>;

interface BaseLayoutProps extends React.PropsWithChildren {
    className?: ClassValue;
    variant?: BaseLayoutVariants["variant"];
}

export const BaseLayout = ({ children, className, variant = "dark" }: BaseLayoutProps) => {
    return <div className={cn(componentVariants({ variant }), className)}>{children}</div>;
};

export { componentVariants };
