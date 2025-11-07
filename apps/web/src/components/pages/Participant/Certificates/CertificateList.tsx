import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { CircleCheckBig, Loader } from "lucide-react";
import { CertificateEmptyState } from "./CertificateEmptyState";
import type { Certificate } from "./useCertificatesListUsecase";
import { Link } from "@/router";

interface CertificateListProps {
    certificates: Certificate[];
}

export const CertificateList = ({ certificates = [] }: CertificateListProps) => {
    const { t } = useTranslation();

    const displayHasContent = certificates?.length ?? 0 > 0;

    if (!displayHasContent) {
        return <CertificateEmptyState />;
    }

    return (
        <div className="flex flex-col gap-y-4 md:gap-y-6">
            <div className="flex flex-col gap-y-2.5">
                {/* Table Header */}
                <div className="flex items-center justify-between px-4 md:px-0">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-xs md:text-sm font-medium"
                    >
                        {t("participant.certificates.title")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-xs md:text-sm font-medium"
                    >
                        {t("common.date")}
                    </Typography>
                </div>

                {/* Divider */}
                <div className="h-px bg-border" />
            </div>

            {/* Certificate items */}
            <div className="flex flex-col gap-y-4 md:gap-y-6">
                {certificates.map((certificate) => (
                    <CertificateItem key={certificate.id} certificate={certificate} />
                ))}
            </div>
        </div>
    );
};

const CertificateItem = ({ certificate }: { certificate: Certificate }) => {
    const isCompleted = certificate.status === "completed";

    const formattedDate = new Date(certificate.issuedDate).toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
    });

    return (
        <Link
            to="/app/certificates/:id"
            params={{ id: certificate.id }}
            className="w-full text-left flex flex-col gap-1 px-0 hover:opacity-80 transition-opacity cursor-pointer group"
        >
            {/* Row with icon, name and date */}
            <div className="flex items-center gap-3 md:gap-4">
                <div className="flex-shrink-0">
                    {isCompleted ? (
                        <CircleCheckBig className="w-5 h-5 text-green-500" />
                    ) : (
                        <Loader className="w-5 h-5 text-yellow-500 animate-spin" />
                    )}
                </div>

                <div className="flex-1 min-w-0">
                    <Typography
                        variant="text"
                        tag="p"
                        className="text-base md:text-lg font-normal underline truncate group-hover:text-primary transition-colors"
                    >
                        {certificate.name}
                    </Typography>
                </div>

                <div className="flex-shrink-0">
                    <Typography
                        variant="text"
                        tag="span"
                        className="text-xs md:text-sm whitespace-nowrap"
                    >
                        {formattedDate}
                    </Typography>
                </div>
            </div>

            {/* Issuer info - below name */}
            <div className="pl-8 md:pl-9">
                <Typography variant="text" tag="p" color="muted" className="text-xs md:text-sm">
                    {certificate.issuer}
                </Typography>
            </div>
        </Link>
    );
};
