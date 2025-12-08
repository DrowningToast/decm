import { Typography } from "@/components/typography/typography";
import { IssuerStatusBadge } from "./IssuerStatusBadge";
import { useTranslation } from "react-i18next";

interface IssuersStatusProps {
    issuers: Array<{
        id: string;
        issuer_credential_id: string;
        is_signed: number;
        issuerProfile?: {
            firstName?: string;
            lastName?: string;
            walletAddress?: string;
            googleConnectorRef?: string;
        };
        // Add other issuer fields as needed
    }>;
    currentIssuerId?: string;
    className?: string;
}

export function IssuersStatus({ issuers, currentIssuerId, className = "" }: IssuersStatusProps) {
    const { t } = useTranslation();

    return (
        <div className={`bg-[#1a1a1a] border border-[#333333] rounded-lg p-6 ${className}`}>
            <Typography variant="header" tag="h3" className="text-xl font-semibold text-white mb-4">
                {t("issuer.sign.issuersStatus")}
            </Typography>

            {issuers.length === 0 ? (
                <Typography
                    variant="text"
                    tag="p"
                    color="muted-foreground"
                    className="text-center py-4"
                >
                    {t("issuer.sign.noIssuersFound")}
                </Typography>
            ) : (
                <div className="space-y-3">
                    {issuers.map((issuer) => (
                        <div
                            key={issuer.id}
                            className={`flex items-center justify-between p-3 rounded-md border ${
                                issuer.issuer_credential_id === currentIssuerId
                                    ? "border-[#ff6a39] bg-[#ff6a3910]"
                                    : "border-[#333333] bg-[#2a2a2a]"
                            }`}
                        >
                            <div className="flex items-center space-x-3">
                                <div className="w-10 h-10 rounded-full bg-gray-600 flex items-center justify-center text-white font-semibold">
                                    {issuer.issuerProfile?.firstName
                                        ? issuer.issuerProfile.firstName
                                              .substring(0, 2)
                                              .toUpperCase()
                                        : issuer.issuer_credential_id
                                              ?.substring(0, 2)
                                              .toUpperCase() || "IS"}
                                </div>
                                <div className="flex-1 min-w-0">
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        className="font-medium text-white truncate"
                                    >
                                        {issuer.issuerProfile?.firstName ||
                                        issuer.issuerProfile?.lastName
                                            ? `${issuer.issuerProfile.firstName || ""} ${issuer.issuerProfile.lastName || ""}`.trim()
                                            : t("issuer.sign.issuer") +
                                              (issuer.id?.substring(0, 8) || "Unknown")}
                                    </Typography>
                                    <div className="flex flex-col gap-0.5 mt-1">
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            className="text-xs text-muted-foreground truncate font-mono"
                                            title={issuer.issuerProfile?.walletAddress}
                                        >
                                            {issuer.issuerProfile?.walletAddress || ""}
                                        </Typography>
                                        {issuer.issuerProfile?.googleConnectorRef && (
                                            <Typography
                                                variant="text"
                                                tag="p"
                                                className="text-xs text-muted-foreground truncate"
                                                title={issuer.issuerProfile.googleConnectorRef}
                                            >
                                                {issuer.issuerProfile.googleConnectorRef}
                                            </Typography>
                                        )}
                                    </div>
                                    {issuer.issuer_credential_id === currentIssuerId && (
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            className="text-sm text-[#ff6a39] mt-1"
                                        >
                                            {t("issuer.sign.you")}
                                        </Typography>
                                    )}
                                </div>
                            </div>

                            <div className="flex items-center space-x-2">
                                <IssuerStatusBadge isSigned={Boolean(issuer.is_signed)} />
                                {Boolean(issuer.is_signed) && (
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        className="text-xs text-green-500"
                                    >
                                        {t("issuer.sign.completed")}
                                    </Typography>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
