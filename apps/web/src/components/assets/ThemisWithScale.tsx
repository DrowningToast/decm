import { cn } from "@/lib/utils";
import type { ClassValue } from "clsx";

interface ThemisWithScaleProps {
    className?: ClassValue;
}

export const ThemisWithScale: React.FC<ThemisWithScaleProps> = ({
    className,
}: ThemisWithScaleProps) => {
    return (
        <div className={cn("relative opacity-30", className)}>
            <img
                src="/assets/themiswithscale.webp"
                alt="Themis with scale"
                className="w-full object-cover"
            />
        </div>
    );
};
