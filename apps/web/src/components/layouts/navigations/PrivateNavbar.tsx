import type React from "react";
import { StaggeredMenu, type StaggeredMenuSection } from "@/components/bits/StaggeredMenu";
import { useTranslation } from "react-i18next";

interface PrivateNavbarProps {
    className?: string;
    variant?: "light" | "dark";
    currentRole?: "Verified Organizer" | "Issuer";
}

export const PrivateNavbar: React.FC<PrivateNavbarProps> = ({ currentRole }) => {
    const { t } = useTranslation();

    const sections: StaggeredMenuSection[] = [
        {
            title: t("nav.account"),
            items: [
                {
                    label: t("nav.settings"),
                    ariaLabel: t("nav.settings"),
                    link: "/app/settings",
                },
                {
                    label: t("nav.signOut"),
                    ariaLabel: t("nav.signOut"),
                    link: "/signout",
                },
            ],
        },
        {
            title: t("nav.other"),
            items: [
                {
                    label: t("nav.preference"),
                    ariaLabel: t("nav.preference"),
                    link: "/app/preference",
                },
                {
                    label: t("nav.readDocs"),
                    ariaLabel: t("nav.readDocs"),
                    link: "/docs",
                },
            ],
        },
    ];

    // Determine logo link based on role
    const getLogoLink = (): string => {
        if (currentRole === "Verified Organizer") {
            return "/host/home";
        }
        if (currentRole === "Issuer") {
            return "/issuer/home";
        }
        return "/app";
    };

    return (
        <StaggeredMenu
            isFixed={true}
            sections={sections}
            position="right"
            roleLabel={currentRole ? `You're currently signed in as a ${currentRole}` : undefined}
            logoLink={getLogoLink()}
        />
    );
};
