import { cn } from "@/lib/utils";
import type { ReactNode } from "react";

interface SectionContainerProps {
  children: ReactNode;
  className?: string;
}

export default function SectionContainer({ children, className }: SectionContainerProps) {
  return <section className={cn("px-6 lg:px-16", className)}>{children}</section>;
}
