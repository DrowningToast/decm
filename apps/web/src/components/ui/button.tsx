import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive cursor-pointer",
  {
    variants: {
      variant: {
        // Figma Design System Variants using CSS variables
        primary:
          "bg-primary text-secondary rounded-[12px] text-base font-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:bg-primary/90 tracking-[0.06px]",
        "secondary-dark":
          "bg-secondary-foreground text-foreground rounded-[12px] text-base font-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:bg-secondary-foreground/90 tracking-[0.06px]",
        "secondary-light":
          "bg-secondary text-primary rounded-[12px] text-base font-medium [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:bg-secondary/90 tracking-[0.06px]",
      },
      size: {
        default: "h-9 px-4 py-2 has-[>svg]:px-3",
        sm: "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5",
        lg: "h-10 rounded-md px-6 has-[>svg]:px-4",
        xl: "h-12 rounded-[12px] px-8 has-[>svg]:px-6",
        icon: "size-9",
      },
    },
    defaultVariants: {
      variant: "primary",
      size: "default",
    },
  }
)

function Button({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean
  }) {
  const Comp = asChild ? Slot : "button"

  return (
    <Comp
      data-slot="button"
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export { Button, buttonVariants }
