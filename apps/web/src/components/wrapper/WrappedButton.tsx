import type { ReactNode } from "react";
import { Button } from "../ui/button";
import { cn } from "@/lib/utils";
import { cva, type VariantProps } from "class-variance-authority";
import { Link, type Path } from "@/router";

const wrappedButtonVariants = cva("h-12 rounded-[12px]", {
  variants: {
    variant: {
      secondary: "bg-white text-primary",
      primary: "bg-primary text-white",
      secondaryWhite: "bg-white text-primary",
      link: "bg-transparent text-primary underline h-fit",
    },
  },
  defaultVariants: {
    variant: "primary",
  },
});

type WrappedButtonProps = {
  children: ReactNode;
  className?: string;
  href?: string;
} & VariantProps<typeof wrappedButtonVariants>;

export default function WrappedButton({ children, className, variant, href }: WrappedButtonProps) {
  const _className = cn(
    wrappedButtonVariants({
      variant,
    }),
    className
  );

  const _Button = <Button className={cn(_className, className)}>{children}</Button>;

  if (href) {
    return <Link to={href as any}>{_Button}</Link>;
  }

  return _Button;
}
