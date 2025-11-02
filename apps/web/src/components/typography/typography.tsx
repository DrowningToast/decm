import { cn } from "@/lib/utils";
import { cva, type VariantProps } from "class-variance-authority";
import type React from "react";

// eslint-disable-next-line
export const typographyVariants = cva("text-wrap", {
    variants: {
        size: {
            small: "text-xs",
            base: "text-base",
            subheader: "text-2xl",
            header: "text-4xl",
        },
        color: {
            foreground: "text-foreground",
            background: "text-background",
            "background-alt": "text-background-alt",
            primary: "text-primary",
            secondary: "text-secondary",
            muted: "text-muted",
            "muted-foreground": "text-muted-foreground",
        },
        variant: {
            header: "font-semibold",
            text: "font-normal",
        },
        fontFamily: {
            inter: "font-inter",
            taviraj: "font-taviraj",
            anuphan: "font-anuphan",
            cormorant: "font-cormorant",
        },
        weight: {
            normal: "font-normal",
            semibold: "font-semibold",
            bold: "font-bold",
            medium: "font-medium",
        },
    },
    defaultVariants: {
        color: "foreground",
        variant: "text",
        size: "base",
        fontFamily: "inter",
        weight: "normal",
    },
});

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const TAGS = ["h1", "h2", "h3", "h4", "h5", "h6", "span", "p", "div", "label"] as const;

type TypographyProps<T extends (typeof TAGS)[number] = "span"> = React.HTMLAttributes<HTMLElement> &
    VariantProps<typeof typographyVariants> & {
        tag: T;
    };

export const Typography = <T extends (typeof TAGS)[number] = "span">({
    className,
    color,
    size,
    variant: level,
    tag,
    children,
    fontFamily,
    weight,
    ...props
}: TypographyProps<T>) => {
    const _className = cn(
        typographyVariants({
            size,
            color,
            variant: level,
            fontFamily,
            weight,
        }),
        className,
    );

    switch (tag) {
        case "h1":
            return (
                <h1 className={_className} {...props}>
                    {children}
                </h1>
            );
        case "h2":
            return (
                <h2 className={_className} {...props}>
                    {children}
                </h2>
            );
        case "h3":
            return (
                <h3 className={_className} {...props}>
                    {children}
                </h3>
            );
        case "h4":
            return (
                <h4 className={_className} {...props}>
                    {children}
                </h4>
            );
        case "h5":
            return (
                <h5 className={_className} {...props}>
                    {children}
                </h5>
            );
        case "h6":
            return (
                <h6 className={_className} {...props}>
                    {children}
                </h6>
            );
        case "span":
            return (
                <span className={_className} {...props}>
                    {children}
                </span>
            );
        case "p":
            return (
                <p className={_className} {...props}>
                    {children}
                </p>
            );
        case "div":
            return (
                <div className={_className} {...props}>
                    {children}
                </div>
            );
        case "label":
            return (
                <label className={_className} {...props}>
                    {children}
                </label>
            );
        default:
            break;
    }
};
