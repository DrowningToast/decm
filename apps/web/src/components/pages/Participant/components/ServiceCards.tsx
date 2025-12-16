import { Link } from "react-router-dom";
import { Typography } from "@/components/typography/typography";
import type { ReactNode } from "react";

interface ServiceCardProps {
    title: string;
    href: string;
    children?: ReactNode;
}

export const ServiceCard: React.FC<ServiceCardProps> = ({ title, href, children }) => {
    return (
        <Link
            to={href}
            className="bg-muted/40 border border-border/50 rounded-lg px-6 py-5 transition-all hover:bg-muted/60 hover:border-border/70 cursor-pointer w-auto md:w-full inline-block relative"
        >
            <div className="flex items-center justify-between gap-2">
                <Typography
                    variant="header"
                    tag="h3"
                    color="foreground"
                    className="text-2xl/[34px] md:text-[24px] font-bold [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] font-header"
                >
                    {title}
                </Typography>
                {children}
            </div>
        </Link>
    );
};
