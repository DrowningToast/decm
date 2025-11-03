import { Link } from "react-router-dom";
import { ServiceCard } from "./ServiceCards";

interface ServiceLink {
    key: string;
    href: string;
    label: string;
}

interface ServiceListProps {
    services: ServiceLink[];
    layout?: "grid" | "single";
}

export const ServiceList = ({ services, layout = "grid" }: ServiceListProps) => {
    return (
        <>
            {/* Mobile: Vertical list with underlined links */}
            <div className="flex md:hidden flex-col gap-1.5 text-[28px]">
                {services.map((service) => (
                    <Link
                        key={service.key}
                        to={service.href}
                        className="h-[34px] text-foreground [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] underline decoration-solid [text-underline-position:from-font] font-header text-[28px] leading-normal hover:text-primary transition-colors"
                    >
                        {service.label}
                    </Link>
                ))}
            </div>

            {/* Desktop: Cards */}
            {layout === "grid" ? (
                <div className="hidden md:grid grid-cols-3 gap-6">
                    {services.map((service) => (
                        <ServiceCard key={service.key} title={service.label} href={service.href} />
                    ))}
                </div>
            ) : (
                <div className="hidden md:block max-w-[352px] md:max-w-full">
                    {services.map((service) => (
                        <ServiceCard key={service.key} title={service.label} href={service.href} />
                    ))}
                </div>
            )}
        </>
    );
};
