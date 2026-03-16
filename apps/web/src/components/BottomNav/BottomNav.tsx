import { useMemo } from "react";
import { SearchCertificateNav } from "./variants/SearchCertificateNav";
import { CertificateDetailNav } from "./variants/CertificateDetailNav";
import { SearchEventNav } from "./variants/SearchEventNav";
import { SearchIdentitiesNav } from "./variants/SearchIdentitiesNav";
import { EventPasswordNav } from "./variants/EventPasswordNav";
import { InvitedNav } from "./variants/InvitedNav";
import { ParticipatingNav } from "./variants/ParticipatingNav";
import { InvitationRequiredNav } from "./variants/InvitationRequiredNav";
import { SearchNotificationNav } from "./variants/SearchNotificationNav";
import { CertificateSigningNav } from "./variants/CertificateSigningNav";
import { InboxViewEventNav } from "./variants/InboxViewEventNav";
import { InboxViewCertificateNav } from "./variants/InboxViewCertificateNav";
import { InboxMissingEventNav } from "./variants/InboxMissingEventNav";
import { BottomContainerProvider } from "./context";
import type { ClassValue } from "clsx";
import { cn } from "@/lib/utils";
import { CertificateDetailSharedNav } from "./variants/CertificateDetailSharedNav";

export type BottomNavVariant =
    | "search-certificate"
    | "search-notification"
    | "certificate-detail"
    | "certificate-details-shared"
    | "search-event"
    | "search-identities"
    | "event-password"
    | "invited"
    | "participating"
    | "invitation-required"
    | "certificate-signing"
    | "inbox-view-event"
    | "inbox-view-certificate"
    | "inbox-missing-event";

interface BottomNavProps {
    variant?: BottomNavVariant;
    onBack?: () => void;
    className?: ClassValue;
}

export const BottomNav = ({
    variant = "search-certificate",
    onBack,
    className,
}: BottomNavProps) => {
    const content = useMemo(() => {
        switch (variant) {
            case "search-certificate":
                return <SearchCertificateNav />;
            case "search-event":
                return <SearchEventNav />;
            case "certificate-detail":
                return <CertificateDetailNav />;
            case "certificate-details-shared":
                return <CertificateDetailSharedNav />;
            case "search-notification":
                return <SearchNotificationNav />;
            case "search-identities":
                return <SearchIdentitiesNav />;
            case "event-password":
                return <EventPasswordNav />;
            case "invited":
                return <InvitedNav />;
            case "participating":
                return <ParticipatingNav />;
            case "invitation-required":
                return <InvitationRequiredNav />;
            case "certificate-signing":
                return <CertificateSigningNav />;
            case "inbox-view-event":
                return <InboxViewEventNav />;
            case "inbox-view-certificate":
                return <InboxViewCertificateNav />;
            case "inbox-missing-event":
                return <InboxMissingEventNav />;
            default:
                return null;
        }
    }, [variant]);

    const isCompact = variant === "certificate-detail" || variant === "certificate-details-shared";

    return (
        <>
            {/* Mobile - Full Width (or compact for certificate variants) */}
            <div
                className={cn(
                    "md:hidden fixed bottom-0 left-0 right-0 z-50 p-4 bg-gradient-to-t from-background to-transparent",
                    isCompact && "flex justify-center",
                    className,
                )}
            >
                <BottomContainerProvider
                    onBack={onBack}
                    className={isCompact ? "w-auto" : "w-full"}
                >
                    <div className={cn("flex flex-col gap-1", isCompact ? "w-auto" : "w-full")}>
                        {content}
                    </div>
                </BottomContainerProvider>
            </div>

            {/* Desktop - Fixed Width (or compact for certificate variants) */}
            <div
                className={cn(
                    "hidden md:flex fixed bottom-12 left-1/2 transform -translate-x-1/2 justify-center z-50 pointer-events-auto",
                    isCompact ? "w-auto" : "w-[700px]",
                    className,
                )}
            >
                <BottomContainerProvider
                    onBack={onBack}
                    className={isCompact ? "w-auto" : "w-full"}
                >
                    <div className={cn("flex flex-col gap-1", isCompact ? "w-auto" : "w-full")}>
                        {content}
                    </div>
                </BottomContainerProvider>
            </div>
        </>
    );
};
