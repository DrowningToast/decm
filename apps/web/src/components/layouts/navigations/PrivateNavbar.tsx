import type React from "react";
import { StaggeredMenu, type StaggeredMenuItem } from "@/components/bits/StaggeredMenu";

const items: StaggeredMenuItem[] = [
    {
        label: "Home",
        ariaLabel: "Home",
        link: "/",
    },
    {
        label: "About",
        ariaLabel: "About",
        link: "/about",
    },
    {
        label: "Contact",
        ariaLabel: "Contact",
        link: "/contact",
    },
] as const;

interface PrivateNavbarProps {
    className?: string;
    variant?: "light" | "dark";
}

export const PrivateNavbar: React.FC<PrivateNavbarProps> = () => {
    return (
        <div className="max-w-screen overflowh-x-hidden">
            <StaggeredMenu isFixed={true} items={items} position="right" />
        </div>
    );
};
