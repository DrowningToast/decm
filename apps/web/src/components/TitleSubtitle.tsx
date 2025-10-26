import { cn } from "@/lib/utils";
import { Typography } from "./typography/typography";

interface TitleSubtitleProps {
  title: string;
  subtitle?: string;
  className?: string;
}

export default function TitleSubtitle({ title, subtitle, className }: TitleSubtitleProps) {
  return (
    <div className={cn("lg:space-y-2", className)}>
      <Typography
        size={"header"}
        tag="p"
        color="primary"
        fontFamily={"cormorant"}
        weight={"semibold"}
        className="lg:text-6xl"
      >
        {title}
      </Typography>

      {subtitle && (
        <Typography size={"base"} tag="p" color="muted">
          {subtitle}
        </Typography>
      )}
    </div>
  );
}
