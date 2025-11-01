import { useState } from "react";
import { cn } from "@/lib/utils";

interface PublicNavbarProps {
    variant?: "primary" | "secondary-dark" | "secondary-light";
    className?: string;
}

export function PublicNavbar({ variant = "primary", className }: PublicNavbarProps) {
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    // Define navbar styles based on variant
    const variantStyles = {
        primary: "bg-[#eb5331]",
        "secondary-dark": "bg-[#362927]",
        "secondary-light": "bg-[#e9dede]",
    };

    const logoVariant = variant === "secondary-light" ? "black" : "white";

    return (
        <nav
            className={cn(
                "fixed top-0 left-0 right-0 z-50 transition-all",
                variantStyles[variant],
                className,
            )}
        >
            <div className="flex items-center justify-between px-4 py-3 md:px-8 md:py-4 max-w-[1512px] mx-auto">
                {/* Logo */}
                <a href="/" className="flex items-center transition-transform hover:scale-105">
                    <img src={`/logo.svg`} alt="Themis Logo" className="h-8 md:h-10 w-20" />
                </a>

                {/* Desktop Navigation */}
                <div className="hidden md:flex items-center gap-8">
                    <a
                        href="#features"
                        className={cn(
                            "text-base font-medium transition-all hover:opacity-80",
                            variant === "secondary-light" ? "text-[#362927]" : "text-white",
                        )}
                    >
                        Features
                    </a>
                    <a
                        href="#about"
                        className={cn(
                            "text-base font-medium transition-all hover:opacity-80",
                            variant === "secondary-light" ? "text-[#362927]" : "text-white",
                        )}
                    >
                        About
                    </a>
                    <a
                        href="/signin"
                        className={cn(
                            "text-base font-medium transition-all hover:opacity-80",
                            variant === "secondary-light" ? "text-[#362927]" : "text-white",
                        )}
                    >
                        Sign In
                    </a>
                </div>

                {/* Mobile Menu Button */}
                <button
                    onClick={() => setIsMenuOpen(!isMenuOpen)}
                    className="md:hidden p-2 transition-transform hover:scale-110"
                    aria-label="Toggle menu"
                >
                    <svg
                        className={cn(
                            "w-6 h-6",
                            variant === "secondary-light" ? "text-[#362927]" : "text-white",
                        )}
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        {isMenuOpen ? (
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M6 18L18 6M6 6l12 12"
                            />
                        ) : (
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M4 6h16M4 12h16M4 18h16"
                            />
                        )}
                    </svg>
                </button>
            </div>

            {/* Mobile Menu Dropdown */}
            {isMenuOpen && (
                <div
                    className={cn(
                        "md:hidden backdrop-blur-sm border-t",
                        variant === "secondary-light"
                            ? "bg-[#e9dede]/95 border-[#362927]/20"
                            : "bg-black/20 border-white/10",
                    )}
                >
                    <div className="flex flex-col gap-1 px-4 py-3">
                        <a
                            href="#features"
                            className={cn(
                                "px-4 py-3 rounded-lg text-base font-medium transition-all",
                                variant === "secondary-light"
                                    ? "text-[#362927] hover:bg-[#362927]/10"
                                    : "text-white hover:bg-white/10",
                            )}
                            onClick={() => setIsMenuOpen(false)}
                        >
                            Features
                        </a>
                        <a
                            href="#about"
                            className={cn(
                                "px-4 py-3 rounded-lg text-base font-medium transition-all",
                                variant === "secondary-light"
                                    ? "text-[#362927] hover:bg-[#362927]/10"
                                    : "text-white hover:bg-white/10",
                            )}
                            onClick={() => setIsMenuOpen(false)}
                        >
                            About
                        </a>
                        <a
                            href="/signin"
                            className={cn(
                                "px-4 py-3 rounded-lg text-base font-medium transition-all",
                                variant === "secondary-light"
                                    ? "text-[#362927] hover:bg-[#362927]/10"
                                    : "text-white hover:bg-white/10",
                            )}
                            onClick={() => setIsMenuOpen(false)}
                        >
                            Sign In
                        </a>
                    </div>
                </div>
            )}
        </nav>
    );
}
