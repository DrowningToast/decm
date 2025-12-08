import { useTranslation } from "react-i18next";
import DOMPurify from "dompurify";
import { Typography } from "@/components/typography/typography";
import { AlertCircle, Eye } from "lucide-react";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import type { DetectedKeyword, AvailableKeyword } from "@/hooks/useCertificateTemplate";
import { useState, useRef, useEffect } from "react";

interface BoundingBoxOverlayProps {
    detectedKeywords: DetectedKeyword[];
    containerRef: React.RefObject<HTMLDivElement>;
}

const BoundingBoxOverlay = ({ detectedKeywords, containerRef }: BoundingBoxOverlayProps) => {
    const [svgDimensions, setSvgDimensions] = useState({
        width: 0,
        height: 0,
        viewBoxWidth: 1920,
        viewBoxHeight: 1080,
        offsetLeft: 0,
        offsetTop: 0,
    });

    useEffect(() => {
        if (!containerRef.current) return;

        const updateDimensions = () => {
            const svg = containerRef.current?.querySelector("svg");
            if (!svg) {
                console.log("SVG not found in container");
                return;
            }

            // Get actual rendered dimensions
            const rect = svg.getBoundingClientRect();
            const containerRect = containerRef.current!.getBoundingClientRect();

            // Get viewBox dimensions (template coordinate space)
            const viewBox = svg.getAttribute("viewBox");
            let viewBoxWidth = 1920;
            let viewBoxHeight = 1080;

            if (viewBox) {
                const parts = viewBox.split(/\s+/);
                if (parts.length === 4) {
                    viewBoxWidth = parseFloat(parts[2]);
                    viewBoxHeight = parseFloat(parts[3]);
                }
            }

            console.log("SVG Dimensions:", {
                width: rect.width,
                height: rect.height,
                viewBoxWidth,
                viewBoxHeight,
                offsetLeft: rect.left - containerRect.left,
                offsetTop: rect.top - containerRect.top,
            });

            setSvgDimensions({
                width: rect.width,
                height: rect.height,
                viewBoxWidth,
                viewBoxHeight,
                offsetLeft: rect.left - containerRect.left,
                offsetTop: rect.top - containerRect.top,
            });
        };

        // Wait for SVG to be rendered
        const timeoutId = setTimeout(updateDimensions, 100);

        // Update on resize
        const resizeObserver = new ResizeObserver(updateDimensions);
        if (containerRef.current) {
            resizeObserver.observe(containerRef.current);
        }

        return () => {
            clearTimeout(timeoutId);
            resizeObserver.disconnect();
        };
    }, [containerRef]);

    if (svgDimensions.width === 0 || svgDimensions.height === 0) {
        console.log("Waiting for SVG dimensions...");
        return null;
    }

    console.log("Rendering overlays for keywords:", detectedKeywords);

    // Calculate scale factors to convert template coordinates to screen coordinates
    const scaleX = svgDimensions.width / svgDimensions.viewBoxWidth;
    const scaleY = svgDimensions.height / svgDimensions.viewBoxHeight;

    return (
        <div className="absolute inset-0 pointer-events-none" style={{ zIndex: 10 }}>
            {detectedKeywords.map((kw, idx) => {
                // Convert template coordinates to screen coordinates
                const screenX = kw.x * scaleX + svgDimensions.offsetLeft;
                const screenY = kw.y * scaleY + svgDimensions.offsetTop;

                // Estimate box dimensions (you may want to make this configurable)
                const boxWidth = 200 * scaleX;
                const boxHeight = 40 * scaleY;

                console.log(`Keyword ${kw.keyword}:`, {
                    templateX: kw.x,
                    templateY: kw.y,
                    screenX,
                    screenY,
                    boxWidth,
                    boxHeight,
                });

                return (
                    <div
                        key={`${kw.keyword}-${idx}`}
                        className="absolute"
                        style={{
                            left: `${screenX - boxWidth / 2}px`,
                            top: `${screenY - boxHeight / 2}px`,
                            width: `${boxWidth}px`,
                            height: `${boxHeight}px`,
                            border: "3px dashed #ef4444",
                            backgroundColor: "rgba(239, 68, 68, 0.2)",
                            borderRadius: "4px",
                            animation: "pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite",
                        }}
                    >
                        <div
                            className="absolute left-0 bg-red-600 text-white text-[10px] px-1.5 py-0.5 rounded font-mono whitespace-nowrap"
                            style={{ top: "-20px" }}
                        >
                            {kw.keyword} ({kw.x.toFixed(0)}, {kw.y.toFixed(0)})
                        </div>
                    </div>
                );
            })}
        </div>
    );
};

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
    const [showPositions, setShowPositions] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);

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
                <div className="flex items-center justify-between">
                    <Typography variant="text" tag="p" className="text-sm font-medium">
                        {t("certificateSettings.step2.preview.title")}
                    </Typography>
                    {detectedKeywords.length > 0 && (
                        <div className="flex items-center gap-2">
                            <Switch
                                id="show-positions"
                                checked={showPositions}
                                onCheckedChange={setShowPositions}
                            />
                            <Label
                                htmlFor="show-positions"
                                className="text-xs cursor-pointer flex items-center gap-1"
                            >
                                <Eye className="w-3 h-3" />
                                {t(
                                    "certificateSettings.step2.preview.showPositions",
                                    "Show Text Positions",
                                )}
                            </Label>
                        </div>
                    )}
                </div>
                <div className="rounded-md border bg-muted/30 p-4">
                    <div
                        className="w-full flex items-center justify-center relative"
                        ref={containerRef}
                    >
                        {renderContent()}
                        {showPositions && detectedKeywords.length > 0 && (
                            <BoundingBoxOverlay
                                detectedKeywords={detectedKeywords}
                                containerRef={containerRef}
                            />
                        )}
                    </div>
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
