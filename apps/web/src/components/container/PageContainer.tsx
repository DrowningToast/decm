import type { ReactNode } from "react";
import { cva } from "class-variance-authority";
import { cn } from "@/lib/utils";
import { PublicNavbar } from "../layouts/navigations/PublicNavbar";

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
  className,
  bgColor = "default",
}: PageContainerProps) {
  const _className = cn(
    pageContainerVariants({
      bgColor,
    }),
    "relative",
    className,
  );

  return (
    <div className={_className}>
      {/* <CustomHelmet title={`${title} | ${PAGE_TITLE_SUBFIX}`} description={description} />
      <FaviconHelmet title={`${title} | ${PAGE_TITLE_SUBFIX}`} description={description} /> */}
      <PublicNavbar />
      <div className="max-w-[1440px] mx-auto space-y-8 mt-14 lg:mt-10">{children}</div>
    </div>
  );
}
