import { useTranslation } from "react-i18next";
import DOMPurify from "dompurify";
import { Typography } from "@/components/typography/typography";
import { AlertCircle } from "lucide-react";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import type { DetectedKeyword, AvailableKeyword } from "@/hooks/useCertificateTemplate";

export interface CertificatePreviewProps {
    svgPreview?: string;
    imageUrl?: string;
    alt?: string;
    detectedKeywords?: DetectedKeyword[];
    availableKeywords?: AvailableKeyword[];
}

export const CertificatePreview = ({
    svgPreview,
    imageUrl,
    alt = "Certificate preview",
    detectedKeywords = [],
    availableKeywords = [],
}: CertificatePreviewProps) => {
    const { t } = useTranslation();

    if (!svgPreview && !imageUrl) {
        return null;
    }

    // Create a map of detected keywords for quick lookup
    const detectedKeywordMap = new Map(detectedKeywords.map((kw) => [kw.keyword, kw]));

    // Get all available keywords with their detection status
    const keywordStatuses = availableKeywords.map((availableKw) => {
        const detected = detectedKeywordMap.get(availableKw.keyword);
        return {
            ...availableKw,
            isDetected: !!detected,
            detectedData: detected,
        };
    });

    // Find missing required keywords
    const missingRequiredKeywords = keywordStatuses.filter((kw) => !kw.isDetected && kw.mandatory);

    const renderContent = () => {
        if (svgPreview) {
            return (
                <div
                    className="w-full [&>svg]:w-full [&>svg]:h-auto [&>svg]:min-h-[300px] [&>svg]:object-contain"
                    dangerouslySetInnerHTML={{
                        __html: DOMPurify.sanitize(svgPreview, {
                            ADD_TAGS: ["pattern", "image", "use", "defs", "clipPath"],
                            ADD_ATTR: [
                                "xlink:href",
                                "patternContentUnits",
                                "patternTransform",
                                "width",
                                "height",
                                "x",
                                "y",
                                "viewBox",
                                "clip-path",
                                "id",
                            ],
                        }),
                    }}
                />
            );
        }

        if (imageUrl) {
            return (
                <img
                    src={imageUrl}
                    alt={alt}
                    className="w-full h-auto min-h-[300px] object-contain rounded-md"
                    onError={(e) => {
                        console.error("Failed to load certificate image:", imageUrl);
                        e.currentTarget.style.display = "none";
                    }}
                />
            );
        }

        return null;
    };

    return (
        <div className="space-y-4">
            {/* Preview Section */}
            <div className="space-y-2">
                <Typography variant="text" tag="p" className="text-sm font-medium">
                    {t("certificateSettings.step2.preview.title")}
                </Typography>
                <div className="rounded-md border bg-muted/30 p-4">
                    <div className="w-full flex items-center justify-center">{renderContent()}</div>
                </div>
            </div>

            {/* Detected Keywords Section */}
            {(svgPreview || detectedKeywords.length > 0) && availableKeywords.length > 0 && (
                <div className="space-y-3">
                    <Typography variant="text" tag="p" className="text-sm font-medium">
                        {t("certificateSettings.step2.detectedKeywords.title")}
                    </Typography>

                    {/* Missing Required Keywords Alert */}
                    {missingRequiredKeywords.length > 0 && (
                        <Alert variant="destructive">
                            <AlertCircle className="h-4 w-4" />
                            <AlertTitle>
                                {t("certificateSettings.step2.missingRequired.title")}
                            </AlertTitle>
                            <AlertDescription>
                                <Typography variant="text" tag="p" className="text-xs mt-2">
                                    {t("certificateSettings.step2.missingRequired.description")}
                                </Typography>
                                <ul className="mt-2 space-y-1">
                                    {missingRequiredKeywords.map((kw) => (
                                        <li key={kw.keyword} className="flex items-center gap-2">
                                            <code className="px-2 py-1 bg-red-100 text-red-900 rounded text-xs font-mono">
                                                {kw.keyword}
                                            </code>
                                        </li>
                                    ))}
                                </ul>
                                <Typography variant="text" tag="p" className="text-xs mt-2">
                                    {t("certificateSettings.step2.missingRequired.instruction")}
                                </Typography>
                            </AlertDescription>
                        </Alert>
                    )}

                    {/* Keyword Status Table */}
                    <div className="rounded-md border bg-card">
                        <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="border-b bg-gray-100">
                                        <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">
                                            {t(
                                                "certificateSettings.step2.keywordStatus.table.keyword",
                                            )}
                                        </th>
                                        <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">
                                            {t(
                                                "certificateSettings.step2.keywordStatus.table.type",
                                            )}
                                        </th>
                                        <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">
                                            {t(
                                                "certificateSettings.step2.keywordStatus.table.count",
                                            )}
                                        </th>
                                        <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">
                                            {t(
                                                "certificateSettings.step2.keywordStatus.table.positionX",
                                            )}
                                        </th>
                                        <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">
                                            {t(
                                                "certificateSettings.step2.keywordStatus.table.positionY",
                                            )}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {keywordStatuses.map((kw) => (
                                        <tr
                                            key={kw.keyword}
                                            className={`border-b last:border-b-0 hover:bg-muted/30 transition-colors ${
                                                !kw.isDetected && kw.mandatory ? "bg-red-50" : ""
                                            }`}
                                        >
                                            <td className="px-4 py-3">
                                                <code className="px-2 py-1 bg-gray-100 border border-gray-300 rounded text-xs font-mono text-gray-900 select-all">
                                                    {kw.keyword}
                                                </code>
                                            </td>
                                            <td className="px-4 py-3">
                                                {kw.mandatory ? (
                                                    <span className="px-2 py-1 bg-red-100 border border-red-300 text-red-900 rounded text-xs font-semibold">
                                                        {t(
                                                            "certificateSettings.step2.keywordStatus.required",
                                                        )}
                                                    </span>
                                                ) : (
                                                    <span className="px-2 py-1 bg-blue-100 border border-blue-300 text-blue-900 rounded text-xs font-semibold">
                                                        {t(
                                                            "certificateSettings.step2.keywordStatus.optional",
                                                        )}
                                                    </span>
                                                )}
                                            </td>
                                            <td className="px-4 py-3">
                                                {kw.isDetected && kw.detectedData ? (
                                                    <span className="text-sm font-medium text-gray-900">
                                                        {kw.detectedData.count}
                                                    </span>
                                                ) : (
                                                    <span className="text-sm font-medium text-gray-400">
                                                        —
                                                    </span>
                                                )}
                                            </td>
                                            <td className="px-4 py-3">
                                                {kw.isDetected && kw.detectedData ? (
                                                    <span className="text-sm font-mono text-gray-900">
                                                        {kw.detectedData.x.toFixed(2)}
                                                    </span>
                                                ) : (
                                                    <span className="text-sm font-medium text-gray-400">
                                                        —
                                                    </span>
                                                )}
                                            </td>
                                            <td className="px-4 py-3">
                                                {kw.isDetected && kw.detectedData ? (
                                                    <span className="text-sm font-mono text-gray-900">
                                                        {kw.detectedData.y.toFixed(2)}
                                                    </span>
                                                ) : (
                                                    <span className="text-sm font-medium text-gray-400">
                                                        —
                                                    </span>
                                                )}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};
