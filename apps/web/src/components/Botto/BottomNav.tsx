import { useMemo } from "react";
import { SearchCertificateNav } from "./SearchCertificateNav";
import { CertificateDetailNav } from "./CertificateDetailNav";
import { BottomContainerProvider } from "./context";

type BottomNavVariant = "search-certificate" | "search-notification" | "certificate-detail";

interface BottomNavProps {
    variant?: BottomNavVariant;
    onBack?: () => void;
}

export const BottomNav = ({ variant = "search-certificate", onBack }: BottomNavProps) => {
    const content = useMemo(() => {
        switch (variant) {
            case "search-certificate":
                return <SearchCertificateNav />;
            case "certificate-detail":
                return <CertificateDetailNav />;
            case "search-notification":
                return "";
            // return <SearchNotificationNav />;
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
