import { cn } from "@/lib/utils";
import { ThemisBlack, ThemisWhite } from "@/components/icons";

interface PublicNavbarProps {
    className?: string;
    variant?: "light" | "dark";
}

export function PublicNavbar({ className, variant = "dark" }: PublicNavbarProps) {
    return (
        <nav
            className={cn(
                "fixed top-0 left-0 right-0 z-50 transition-all backdrop-blur-[2px] shadow-[0px_4px_16px_0px_rgba(0,0,0,0.1)] overflow-hidden bg-[rgba(217,217,217,0.1)] flex justify-center",
                className,
            )}
        >
            <div className="flex items-center justify-between px-3 py-1.5 md:px-6 max-w-[1512px] mx-auto">
                {/* Logo */}
                <a href="/" className="flex items-center transition-transform hover:scale-105">
                    <div className="size-9 md:size-12 grid place-items-center relative">
                        {variant === "light" ? (
                            <ThemisBlack className="w-full h-full" />
                        ) : (
                            <ThemisWhite className="w-full h-full" />
                        )}
                    </div>
                </a>
            </div>
        </nav>
    );
}
