import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { CertificateList } from "./CertificateList";
import { useCertificatesListUsecase } from "./useCertificatesListUsecase";
import { BottomNav } from "@/components/BottomNav/BottomNav";

export const CertificateListPage = () => {
    const { t } = useTranslation();
    const { certificates, isLoading } = useCertificatesListUsecase();

    return (
        <section className="relative z-10 w-full">
            <div className="relative min-h-screen w-full overflow-hidden">
                {/* Background image */}
                <div className="absolute bottom-0 right-0 w-[424px] h-[424px] md:w-[500px] md:h-[500px] opacity-20 pointer-events-none">
                    <img
                        src="/assets/passport.webp"
                        alt=""
                        className="w-full h-full object-cover object-center"
                    />
                </div>

                {/* Main content */}
                <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 py-4 md:pt-16 md:pb-24 flex flex-col gap-y-4 md:gap-y-6">
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
                    <CertificateList certificates={certificates} isLoading={isLoading} />
                </div>

                {/* Bottom Navigation */}
                <BottomNav variant="search-certificate" onBack={() => window.history.back()} />
            </div>
        </section>
    );
};
