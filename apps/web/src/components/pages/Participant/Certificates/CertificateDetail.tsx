import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { useCertificateDetailUsecase } from "./useCertificateDetailUsecase";
import { BottomNav } from "@/components/Botto/BottomNav";
import { CircleCheckBig, ExternalLink } from "lucide-react";
import { Link } from "react-router-dom";

interface CertificateDetailProps {
    certificateId: string;
}

export const CertificateDetail = ({ certificateId }: CertificateDetailProps) => {
    const { t } = useTranslation();
    const { certificate, formattedDate } = useCertificateDetailUsecase(certificateId);

    if (!certificate) {
        return (
            <div className="relative min-h-screen w-full overflow-hidden flex items-center justify-center pb-24 md:pb-12">
                <Typography variant="text" tag="p" color="muted" className="text-lg">
                    {t("common.error")}
                </Typography>
            </div>
        );
    }

    return (
        <div className="relative w-full overflow-hidden">
            {/* Background image */}
            <div className="absolute bottom-0 right-0 w-[424px] h-[424px] md:w-[500px] md:h-[500px] opacity-20 pointer-events-none">
                <img
                    src="/assets/passport.webp"
                    alt=""
                    className="w-full h-full object-cover object-center"
                />
            </div>

            {/* Main content */}
            <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 py-4 md:py-16 flex flex-col gap-y-4">
                {/* Certificate Header with Check Icon */}
                <div className="flex flex-col gap-y-1.5">
                    <div className="flex items-center gap-2">
                        <CircleCheckBig className="w-6 h-6 text-primary shrink-0" />
                        <Typography
                            variant="header"
                            tag="h1"
                            color="foreground"
                            className="text-[28px] font-header leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {certificate.name}
                        </Typography>
                    </div>
                    <Link
                        to={`/app/events/${certificate.eventId}`}
                        className="underline text-muted"
                    >
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {certificate.event}
                        </Typography>
                    </Link>
                </div>

                {/* Certificate Image Placeholder */}
                <div className="w-full h-[172px] bg-muted rounded-lg overflow-hidden">
                    {certificate.certificateImageUrl ? (
                        <img
                            src={certificate.certificateImageUrl}
                            alt={certificate.name}
                            className="w-full h-full object-cover"
                        />
                    ) : (
                        <div className="w-full h-full bg-[#d9d9d9]" />
                    )}
                </div>

                {/* Certificate Description */}
                {certificate.description && (
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.certificates.detail.certificateDescription")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="foreground"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {certificate.description}
                        </Typography>
                    </div>
                )}

                {/* You signed certification on */}
                <div className="flex flex-col gap-y-1">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("participant.certificates.detail.youSignedCertificationOn")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="foreground"
                        className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {formattedDate}
                    </Typography>
                </div>

                {/* Certificate Contract Address */}
                {certificate.certificateContractAddress && (
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.certificates.detail.certificateContractAddress")}
                        </Typography>
                        <div className="flex items-center gap-2.5">
                            <Typography
                                variant="text"
                                tag="p"
                                color="foreground"
                                className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {certificate.certificateContractAddress}
                            </Typography>
                            <ExternalLink className="w-4 h-4 text-foreground shrink-0" />
                        </div>
                    </div>
                )}

                {/* Event Contract Address */}
                {certificate.eventContractAddress && (
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.certificates.detail.eventContractAddress")}
                        </Typography>
                        <div className="flex items-center gap-2.5">
                            <Typography
                                variant="text"
                                tag="p"
                                color="foreground"
                                className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {certificate.eventContractAddress}
                            </Typography>
                            <ExternalLink className="w-4 h-4 text-foreground shrink-0" />
                        </div>
                    </div>
                )}

                {/* Verifiable Credential Attributes */}
                {certificate.verifiableCredentialUrl && (
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.certificates.detail.verifiableCredentialAttributes")}
                        </Typography>
                        <a
                            href={certificate.verifiableCredentialUrl}
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] text-foreground underline"
                        >
                            {t("participant.certificates.detail.downloadAttributes")}
                        </a>
                    </div>
                )}
            </div>

            {/* Bottom Navigation */}
            <BottomNav variant="certificate-detail" onBack={() => window.history.back()} />
        </div>
    );
};
