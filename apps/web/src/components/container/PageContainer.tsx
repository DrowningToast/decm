import type { ReactNode } from "react";
import { cva } from "class-variance-authority";
import { cn } from "@/lib/utils";

const pageContainerVariants = cva("py-10 min-h-dvh lg:py-16", {
  variants: {
    bgColor: {
      default: "bg-black",
      primary: "from-[#EB5331] to-[#362927] bg-gradient-to-b",
    },
  },
  defaultVariants: {
    bgColor: "default",
  },
});

interface PageContainerProps {
  children: ReactNode;
  title: string;
  description?: string;
  bgColor?: "default" | "primary";
  className?: string;
}

export default function PageContainer({
  children,
  title,
  description,
  className,
  bgColor = "default",
}: PageContainerProps) {
  const _className = cn(
    pageContainerVariants({
      bgColor,
    }),
    "space-y-8 relative",
    className
  );

  return (
    <div className={_className}>
      {/* <CustomHelmet title={`${title} | ${PAGE_TITLE_SUBFIX}`} description={description} />
      <FaviconHelmet title={`${title} | ${PAGE_TITLE_SUBFIX}`} description={description} /> */}

      {children}
    </div>
  );
}
