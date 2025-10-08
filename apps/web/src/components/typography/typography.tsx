
import { cn } from "@/lib/utils"
import { cva, type VariantProps } from "class-variance-authority"
import type { PropsWithChildren } from "react"
import type React from "react"

const typographyVariants = cva(
    "text-wrap",
    {
        variants: {
            size: {
                small: "text-xs",
                base: "text-base",
                subheader: "text-2xl",
                header: "text-4xl"
            },
            color: {
                foreground: "text-foreground",
                background: "text-background",
                ["background-alt"]: "text-background-alt",
                primary: "text-primary",
                secondary: "text-secondary",
                muted: "text-muted"
            },
            variant: {
                header: "font-primary font-semibold",
                text: "font-secondary font-normal"
            },
        },
        defaultVariants: {
            color: 'foreground',
            variant: "text",
            size: "base"
        }
    }
)

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const TAGS = [
    "h1",
    "h2",
    "h3",
    "h4",
    "h5",
    "h6",
    "span",
    "p",
    "div"
] as const

type TypographyProps =
    React.ComponentProps<'span'> &
    React.ComponentProps<'p'> &
    React.ComponentProps<'h1'> &
    React.ComponentProps<'div'> &
    VariantProps<typeof typographyVariants> &
    PropsWithChildren &
    {
        tag: (typeof TAGS)[number]
    }

export const Typography: React.FC<TypographyProps> = ({
    className,
    color,
    size,
    variant: level,
    tag,
    children,
    ...props
}) => {

    const _className = cn(typographyVariants({
        size,
        color,
        variant: level
    }), className)

    switch (tag) {
        case "h1":
            return (
                <h1 className={_className} {...props}>
                    {children}
                </h1>
            )
        case "h2":
            return (
                <h2 className={_className} {...props}>
                    {children}
                </h2>
            )
        case "h3":
            return (
                <h3 className={_className} {...props}>
                    {children}
                </h3>
            )
        case "h4":
            return (
                <h4 className={_className} {...props}>
                    {children}
                </h4>
            )
        case "h5":
            return (
                <h5 className={_className} {...props}>
                    {children}
                </h5>
            )
        case "h6":
            return (
                <h6 className={_className} {...props}>
                    {children}
                </h6>
            )
        case "span":
            return (
                <span className={_className} {...props}>
                    {children}
                </span>
            )
        case "p":
            return (
                <p className={_className} {...props}>
                    {children}
                </p>
            )
        case "div":
            return (
                <div className={_className} {...props}>
                    {children}
                </div>
            )
        default:
            break;
    }

}