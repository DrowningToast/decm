import { Typography } from "@/components/typography/typography";
import { calculateSigningProgress } from "./issuerStateUtils";
import { useTranslation } from "react-i18next";

interface SigningProgressProps {
    issuers: Array<{ is_signed: number }>;
    className?: string;
}

export function SigningProgress({ issuers, className = "" }: SigningProgressProps) {
    const { t } = useTranslation();
    const progress = calculateSigningProgress(
        issuers.map((issuer) => ({ is_signed: issuer.is_signed === 1 })),
    );
    // const signedCount = issuers?.filter((issuer) => issuer.is_signed === 1).length || 0;
    const totalCount = issuers?.length || 0;

    return (
        <div className={`w-full ${className}`}>
            <div className="flex justify-between items-center mb-2">
                <Typography variant="text" tag="span" className="text-sm font-medium">
                    {t("issuer.sign.signingProgress")}
                </Typography>
                <Typography variant="text" tag="span" className="text-sm text-gray-500">
                    {t("issuer.sign.ofSigned", { count: totalCount })}
                </Typography>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2.5">
                <div
                    className="bg-[#ff6a39] h-2.5 rounded-full transition-all duration-300"
                    style={{ width: `${progress}%` }}
                ></div>
            </div>
            <div className="mt-1 text-right">
                <Typography variant="text" tag="span" className="text-xs text-gray-500">
                    {progress}%
                </Typography>
            </div>
        </div>
    );
}
