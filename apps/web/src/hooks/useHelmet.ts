import { useMemo } from "react";

interface UseHelmetProps {
    title?: string;
    description?: string;
    pageType?: "home" | "login" | "profile" | "events" | "credentials";
}

interface HelmetData {
    title: string;
    description: string;
    faviconUrl: string;
    themeColor: string;
}

export const useHelmet = ({
    title,
    description,
    pageType = "home",
}: UseHelmetProps = {}): HelmetData => {
    return useMemo(() => {
        const baseTitle = "DECM - Decentralized Event Management";
        const baseDescription =
            "Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification";

        const pageConfig = {
            home: {
                title: title || `Home | ${baseTitle}`,
                description: description || baseDescription,
                faviconUrl: "/favicon.ico",
                themeColor: "#ffffff",
            },
            login: {
                title: title || `Login | ${baseTitle}`,
                description: description || "Sign in to your DECM account",
                faviconUrl: "/favicon.ico",
                themeColor: "#3B82F6", // Blue theme for login
            },
            profile: {
                title: title || `Profile | ${baseTitle}`,
                description: description || "Manage your DECM profile and credentials",
                faviconUrl: "/favicon.ico",
                themeColor: "#10B981", // Green theme for profile
            },
            events: {
                title: title || `Events | ${baseTitle}`,
                description: description || "Browse and manage events on DECM platform",
                faviconUrl: "/favicon.ico",
                themeColor: "#F59E0B", // Orange theme for events
            },
            credentials: {
                title: title || `Credentials | ${baseTitle}`,
                description: description || "View and manage your digital credentials",
                faviconUrl: "/favicon.ico",
                themeColor: "#8B5CF6", // Purple theme for credentials
            },
        };

        return pageConfig[pageType];
    }, [title, description, pageType]);
};
