import { Link } from "react-router-dom";
import { Typography } from "@/components/typography/typography";

interface ServiceCardProps {
    title: string;
    href: string;
}

export const ServiceCard: React.FC<ServiceCardProps> = ({ title, href }) => {
    return (
        <Link
            to={href}
            className="bg-muted border border-border rounded-lg px-6 py-5 transition-all hover:bg-muted/50 hover:border-border/20 cursor-pointer w-auto md:w-full inline-block"
        >
            <Typography
                variant="header"
                tag="h3"
                color="foreground"
                className="text-2xl/[34px] md:text-[24px] font-bold [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] font-header"
            >
                {title}
            </Typography>
        </Link>
    );
};
