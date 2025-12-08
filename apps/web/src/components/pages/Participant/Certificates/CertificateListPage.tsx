import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { BottomNav } from "@/components/BottomNav/BottomNav";
import { CertificateList } from "./CertificateList";
import { useMyCertificatesListViewModel } from "@/hooks/useMyCertificatesListViewModel";

export const CertificateListPage = () => {
    const { t } = useTranslation();
    const { claimedCertificates, unclaimedCertificates, isLoading } =
        useMyCertificatesListViewModel();

    return (
        <div className="relative w-full overflow-hidden">
            {/* Main content */}
            <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 pb-4 md:pb-24 flex flex-col gap-y-4 md:gap-y-6">
                {/* Header section */}
                <div className="flex flex-col gap-1.5">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="primary"
                        className="text-[28px]/[34px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] font-header"
                    >
                        {t("participant.certificates.title")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-base/[22px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("participant.certificates.description")}
                    </Typography>
                </div>

                {/* Certificate List */}
                <CertificateList
                    claimedCertificates={claimedCertificates}
                    unclaimedCertificates={unclaimedCertificates}
                    isLoading={isLoading}
                />
            </div>

            {/* Bottom Navigation */}
            <BottomNav variant="search-certificate" onBack={() => window.history.back()} />
        </div>
    );
};
