import type React from "react";
import { StaggeredMenu, type StaggeredMenuSection } from "@/components/bits/StaggeredMenu";
import { useTranslation } from "react-i18next";
import { useCheckRoles } from "@/hooks/useCheckRoles";

interface PrivateNavbarProps {
    className?: string;
    variant?: "light" | "dark";
    currentRole?: "Verified Organizer" | "Issuer";
}

export const PrivateNavbar: React.FC<PrivateNavbarProps> = ({ currentRole, variant = "dark" }) => {
    const { t } = useTranslation();
    const { isHost, isIssuer } = useCheckRoles({
        requireHost: true,
        requireIssuer: true,
    });

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
                    label: t("nav.home"),
                    ariaLabel: t("nav.home"),
                    link: "/app",
                },
                isHost
                    ? {
                          label: t("nav.host"),
                          ariaLabel: t("nav.host"),
                          link: "/host/home",
                      }
                    : null,
                isIssuer
                    ? {
                          label: t("nav.issuer"),
                          ariaLabel: t("nav.issuer"),
                          link: "/issuer/home",
                      }
                    : null,
                {
                    label: t("nav.preference"),
                    ariaLabel: t("nav.preference"),
                    link: "/app/preference",
                },
                {
                    label: t("nav.readDocs"),
                    ariaLabel: t("nav.readDocs"),
                    link: "/docs",
                    disabled: true,
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

    // Get translated role label
    const getRoleLabel = (): string | undefined => {
        if (!currentRole) return undefined;
        const roleKey =
            currentRole === "Verified Organizer"
                ? "nav.roles.verifiedOrganizer"
                : "nav.roles.issuer";
        const role = t(roleKey);
        return t("nav.signedInAs", { role });
    };

    return (
        <StaggeredMenu
            isFixed={true}
            sections={sections}
            position="right"
            roleLabel={getRoleLabel()}
            logoLink={getLogoLink()}
            variant={variant}
        />
    );
};
