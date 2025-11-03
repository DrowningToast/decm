import { useMemo } from "react";
import { SearchCertificateNav } from "./variants/SearchCertificateNav";
import { CertificateDetailNav } from "./variants/CertificateDetailNav";
import { SearchEventNav } from "./variants/SearchEventNav";
import { SearchIdentitiesNav } from "./variants/SearchIdentitiesNav";
import { EventPasswordNav } from "./variants/EventPasswordNav";
import { InvitedNav } from "./variants/InvitedNav";
import { ShortlistedNav } from "./variants/ShortlistedNav";
import { ParticipatingNav } from "./variants/ParticipatingNav";
import { InvitationRequiredNav } from "./variants/InvitationRequiredNav";
import { SearchNotificationNav } from "./variants/SearchNotificationNav";
import { CertificateSigningNav } from "./variants/CertificateSigningNav";
import { InboxViewEventNav } from "./variants/InboxViewEventNav";
import { InboxViewCertificateNav } from "./variants/InboxViewCertificateNav";
import { InboxMissingEventNav } from "./variants/InboxMissingEventNav";
import { BottomContainerProvider } from "./context";

type BottomNavVariant =
    | "search-certificate"
    | "search-notification"
    | "certificate-detail"
    | "search-event"
    | "search-identities"
    | "event-password"
    | "invited"
    | "shortlisted"
    | "participating"
    | "invitation-required"
    | "certificate-signing"
    | "inbox-view-event"
    | "inbox-view-certificate"
    | "inbox-missing-event";

interface BottomNavProps {
    variant?: BottomNavVariant;
    onBack?: () => void;
}

export const BottomNav = ({ variant = "search-certificate", onBack }: BottomNavProps) => {
    const content = useMemo(() => {
        switch (variant) {
            case "search-certificate":
                return <SearchCertificateNav />;
            case "search-event":
                return <SearchEventNav />;
            case "certificate-detail":
                return <CertificateDetailNav />;
            case "search-notification":
                return <SearchNotificationNav />;
            case "search-identities":
                return <SearchIdentitiesNav />;
            case "event-password":
                return <EventPasswordNav />;
            case "invited":
                return <InvitedNav />;
            case "shortlisted":
                return <ShortlistedNav />;
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

    return (
        <>
            {/* Mobile - Full Width */}
            <div className="md:hidden fixed bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-background to-transparent">
                <BottomContainerProvider onBack={onBack} className="w-full">
                    <div className="flex flex-col gap-1 w-full">{content}</div>
                </BottomContainerProvider>
            </div>

            {/* Desktop - Fixed Width */}
            <div className="hidden md:flex fixed bottom-12 left-1/2 transform -translate-x-1/2 justify-center z-50 pointer-events-auto">
                <BottomContainerProvider onBack={onBack} className="w-[343px]">
                    <div className="flex flex-col gap-1 w-[343px]">{content}</div>
                </BottomContainerProvider>
            </div>
        </>
    );
};
