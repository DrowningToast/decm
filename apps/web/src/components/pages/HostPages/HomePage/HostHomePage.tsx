import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import { Typography } from "@/components/typography/typography";
import { Link, type Path } from "@/router";

export default function HostHomePage() {
    return (
        <div title="Host Dashboard">
            <img
                src="/justice.png"
                alt="Host Dashboard Background"
                className="absolute bottom-0 lg:right-0 w-[400px] h-[400px] opacity-50 m-0"
            />

            <SectionContainer>
                <TitleSubtitle title="Good Morning" subtitle="You've 0 action required." />
            </SectionContainer>

            <SectionContainer className="space-y-2">
                <Typography size={"base"} tag="p" color="muted">
                    HOST SERVICES
                </Typography>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-2 lg:gap-6 lg:mt-6">
                    <MenuItem title="Event" to="/host/events" />
                    <MenuItem title="Inbox" to="/host/home" />
                    <MenuItem title="Profile" to="/host/home" />
                </div>
            </SectionContainer>

            <SectionContainer>
                <Typography size={"base"} tag="p" color="muted">
                    GENERAL SERVICES
                </Typography>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-2 lg:gap-6 lg:mt-6">
                    <MenuItem title="Verify certificate" to="/host/home" />
                </div>
            </SectionContainer>
        </div>
    );
}

interface MenuItemProps {
    title: string;
    to: Path;
}
function MenuItem({ title, to }: MenuItemProps) {
    return (
        <Link
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            to={to as any}
            className="lg:py-5 lg:px-6 lg:border lg:border-[#D9D9D91A] lg:rounded-lg lg:bg-[#D9D9D905]"
        >
            <Typography
                size={"header"}
                tag="p"
                color="secondary"
                className="underline lg:no-underline lg:text-2xl font-cormorant"
            >
                {title}
            </Typography>
        </Link>
    );
}
