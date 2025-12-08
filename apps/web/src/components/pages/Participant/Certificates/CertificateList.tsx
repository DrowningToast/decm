import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Loader2, SearchX, CircleCheckBig } from "lucide-react";
import { CertificateEmptyState } from "./CertificateEmptyState";
import { Link } from "@/router";
import { useSearchCertificateNavStore } from "@/components/BottomNav/stores/certificates";
import type { EntityEventCertificate } from "@decm/api";

interface CertificateListProps {
    claimedCertificates?: EntityEventCertificate[];
    unclaimedCertificates?: EntityEventCertificate[];
    isLoading?: boolean;
}

export const CertificateList = ({
    claimedCertificates = [],
    unclaimedCertificates = [],
    isLoading,
}: CertificateListProps) => {
    const { t } = useTranslation();
    const { searchQuery } = useSearchCertificateNavStore();

    // Filter certificates based on search query
    const filteredClaimedCertificates = searchQuery.trim()
        ? claimedCertificates.filter(
              (cert) =>
                  cert.certificate_title?.toLowerCase().includes(searchQuery.toLowerCase()) ||
                  cert.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
                  cert.event_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
                  cert.academic_institution?.toLowerCase().includes(searchQuery.toLowerCase()),
          )
        : claimedCertificates;

    const filteredUnclaimedCertificates = searchQuery.trim()
        ? unclaimedCertificates.filter(
              (cert) =>
                  cert.certificate_title?.toLowerCase().includes(searchQuery.toLowerCase()) ||
                  cert.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
                  cert.event_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
                  cert.academic_institution?.toLowerCase().includes(searchQuery.toLowerCase()),
          )
        : unclaimedCertificates;

    const hasContent =
        filteredClaimedCertificates.length > 0 || filteredUnclaimedCertificates.length > 0;
    const isSearching = searchQuery.trim().length > 0;

    // Loading state
    if (isLoading) {
        return (
            <div className="flex flex-col items-center justify-center gap-y-4 py-12 md:py-16">
                <Loader2 className="w-12 h-12 text-primary animate-spin" />
                <Typography variant="text" tag="p" color="muted" className="text-base">
                    {t("common.loading", "Loading...")}
                </Typography>
            </div>
        );
    }

    // No results found (search applied but no matches)
    if (!hasContent && isSearching) {
        return (
            <div className="flex flex-col items-center justify-center gap-y-4 py-12 md:py-16">
                <div className="rounded-full bg-muted/30 p-4 md:p-6">
                    <SearchX className="w-8 h-8 md:w-12 md:h-12 text-muted-foreground" />
                </div>
                <div className="text-center flex flex-col gap-y-2">
                    <Typography
                        variant="text"
                        tag="h3"
                        className="text-lg md:text-xl font-semibold"
                    >
                        {t("participant.certificates.noResults.title", "No certificates found")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-sm md:text-base"
                    >
                        {t(
                            "participant.certificates.noResults.description",
                            "Try adjusting your search",
                        )}
                    </Typography>
                </div>
            </div>
        );
    }

    // Truly empty state (no data at all)
    if (!hasContent) {
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
                {/* Claimed certificates */}
                {filteredClaimedCertificates.map((certificate) => (
                    <CertificateItem
                        key={certificate.id}
                        certificate={certificate}
                        isClaimed={true}
                    />
                ))}
                {/* Unclaimed certificates */}
                {filteredUnclaimedCertificates.map((certificate) => (
                    <CertificateItem
                        key={certificate.id}
                        certificate={certificate}
                        isClaimed={false}
                    />
                ))}
            </div>
        </div>
    );
};

const CertificateItem = ({
    certificate,
    isClaimed,
}: {
    certificate: EntityEventCertificate;
    isClaimed: boolean;
}) => {
    const { t } = useTranslation();
    const formattedDate = new Date(certificate.created_at).toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
    });

    // Fallback order: certificate_title → certificate_subtitle → default text
    const displayTitle =
        certificate.certificate_title ||
        certificate.certificate_subtitle ||
        t("participant.certificates.defaultTitle", "Certificate of Achievement");

    return (
        <Link
            to="/app/certificates/:id"
            params={{ id: certificate.id }}
            className="w-full text-left flex flex-col gap-1 px-0 hover:opacity-80 transition-opacity cursor-pointer group"
        >
            {/* Row with icon, name and date */}
            <div className="flex items-center gap-3 md:gap-4">
                {/* Status Icon */}
                <div className="flex-shrink-0">
                    {isClaimed ? (
                        <CircleCheckBig className="w-5 h-5 text-green-500" />
                    ) : (
                        <Loader2 className="w-5 h-5 text-yellow-500 animate-spin" />
                    )}
                </div>

                <div className="flex-1 min-w-0">
                    <Typography
                        variant="text"
                        tag="p"
                        className="text-base md:text-lg font-normal underline truncate group-hover:text-primary transition-colors"
                    >
                        {displayTitle}
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

            {/* Event name and Issuer info - below name */}
            <div className="ml-8 flex flex-col gap-0.5">
                {certificate.event_name && (
                    <Typography variant="text" tag="p" color="muted" className="text-xs md:text-sm">
                        {certificate.event_name}
                    </Typography>
                )}
                {certificate.academic_institution && (
                    <Typography variant="text" tag="p" color="muted" className="text-xs md:text-sm">
                        {certificate.academic_institution}
                    </Typography>
                )}
            </div>
        </Link>
    );
};
