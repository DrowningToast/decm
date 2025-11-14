import { ChevronLeft } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { Typography } from "@/components/typography/typography";

interface CertificateSigningNavProps {
    className?: string;
}

export const CertificateSigningNav = ({ className: propClassName }: CertificateSigningNavProps) => {
    const { onBack, className: contextClassName } = useBottomContainerContext();
    const { t } = useTranslation();

    const handleSignCertificate = () => {
        // TODO: Implement certificate signing functionality
        console.log("Sign certificate clicked");
    };

    return (
        <div
            className={cn(
                contextClassName,
                propClassName,
                "flex items-center gap-1.5 h-13 bg-primary rounded-xl p-1.5",
            )}
        >
            {/* Back Button */}
            <button
                onClick={onBack}
                className="cursor-pointer flex items-center justify-center w-10 h-10 bg-primary rounded-lg hover:bg-primary/90 transition-colors flex-shrink-0"
                aria-label="Go back"
            >
                <ChevronLeft className="w-5 h-5 text-white" />
            </button>

            {/* Message Box */}
            <button
                onClick={handleSignCertificate}
                className="flex-1 h-10 bg-white rounded-lg flex items-center justify-center px-4 hover:bg-white/90 transition-colors"
                aria-label={t("participant.certificates.sign") || "Sign your certificate"}
            >
                <Typography
                    variant="text"
                    tag="span"
                    color="background-alt"
                    className="text-xs font-normal leading-normal tracking-[0.06px] text-center whitespace-nowrap"
                >
                    {t("participant.certificates.sign") || "Click here to sign your certificate"}
                </Typography>
            </button>
        </div>
    );
};
