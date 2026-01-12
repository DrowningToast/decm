import { useTranslation } from "react-i18next";
import DOMPurify from "dompurify";
import { Typography } from "@/components/typography/typography";
import { AlertCircle, Eye, Type } from "lucide-react";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import type { DetectedKeyword, AvailableKeyword } from "@/hooks/useCertificateTemplate";
import { useState, useRef, useEffect, useCallback, useMemo } from "react";
import { processCertificateSVG, type CertificateFontConfig } from "@/lib/certificateRenderer";
import { useCertificateFontFamilies } from "@/hooks/events/useCertificateFontFamilies";

interface BoundingBoxOverlayProps {
    detectedKeywords: DetectedKeyword[];
    containerRef: React.RefObject<HTMLDivElement | null>;
    onKeywordMove?: (keyword: string, x: number, y: number) => void;
    /** Changing this value re-triggers dimension measurement (e.g. when SVG loads asynchronously). */
    svgContent?: string;
}

const BoundingBoxOverlay = ({
    detectedKeywords,
    containerRef,
    onKeywordMove,
    svgContent,
}: BoundingBoxOverlayProps) => {
    const [svgDimensions, setSvgDimensions] = useState({
        width: 0,
        height: 0,
        viewBoxWidth: 1920,
        viewBoxHeight: 1080,
        offsetLeft: 0,
        offsetTop: 0,
    });

    const [draggingKeyword, setDraggingKeyword] = useState<string | null>(null);
    const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 });

    useEffect(() => {
        if (!containerRef.current) return;

        const updateDimensions = () => {
            const container = containerRef.current;
            if (!container) return;

            const svg = container.querySelector("svg");
            if (svg) {
                // Get actual rendered dimensions
                const rect = svg.getBoundingClientRect();
                const containerRect = container.getBoundingClientRect();

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

                setSvgDimensions({
                    width: rect.width,
                    height: rect.height,
                    viewBoxWidth,
                    viewBoxHeight,
                    offsetLeft: rect.left - containerRect.left,
                    offsetTop: rect.top - containerRect.top,
                });
                return;
            }

            // Fall back to <img> when the SVG was fetched via presigned URL and
            // the raw SVG string is unavailable (e.g. CORS prevents the fetch).
            // Use the image's natural dimensions as the coordinate space so the
            // overlay works correctly for any SVG size, not just 1920×1080.
            const img = container.querySelector("img");
            if (!img || img.clientWidth === 0 || img.clientHeight === 0) {
                return;
            }

            const rect = img.getBoundingClientRect();
            const containerRect = container.getBoundingClientRect();

            // naturalWidth/naturalHeight reflect the SVG's intrinsic viewBox dimensions
            // (available after the image has loaded; fall back to rendered size otherwise)
            const viewBoxWidth = img.naturalWidth || rect.width;
            const viewBoxHeight = img.naturalHeight || rect.height;

            setSvgDimensions({
                width: rect.width,
                height: rect.height,
                viewBoxWidth,
                viewBoxHeight,
                offsetLeft: rect.left - containerRect.left,
                offsetTop: rect.top - containerRect.top,
            });
        };

        // Wait for SVG/img to be rendered
        const timeoutId = setTimeout(updateDimensions, 100);

        // Update on resize
        const resizeObserver = new ResizeObserver(updateDimensions);
        if (containerRef.current) {
            resizeObserver.observe(containerRef.current);
        }

        // If content is an <img>, re-measure once it finishes loading
        const img = containerRef.current?.querySelector("img");
        if (img) {
            img.addEventListener("load", updateDimensions);
        }

        return () => {
            clearTimeout(timeoutId);
            resizeObserver.disconnect();
            img?.removeEventListener("load", updateDimensions);
        };
    }, [containerRef, svgContent]);

    const handleMouseDown = (e: React.MouseEvent, keyword: string, kwX: number, kwY: number) => {
        if (!onKeywordMove) return;
        e.preventDefault();
        e.stopPropagation();

        const scaleX = svgDimensions.width / svgDimensions.viewBoxWidth;
        const scaleY = svgDimensions.height / svgDimensions.viewBoxHeight;

        const screenX = kwX * scaleX + svgDimensions.offsetLeft;
        const screenY = kwY * scaleY + svgDimensions.offsetTop;

        setDraggingKeyword(keyword);
        setDragOffset({
            x: e.clientX - screenX,
            y: e.clientY - screenY,
        });
    };

    useEffect(() => {
        if (!draggingKeyword) return;

        const handleMouseMove = (e: MouseEvent) => {
            if (!containerRef.current || !onKeywordMove) return;

            const scaleX = svgDimensions.width / svgDimensions.viewBoxWidth;
            const scaleY = svgDimensions.height / svgDimensions.viewBoxHeight;

            // Calculate new screen position
            const newScreenX = e.clientX - dragOffset.x;
            const newScreenY = e.clientY - dragOffset.y;

            // Convert back to template coordinates
            let newTemplateX = (newScreenX - svgDimensions.offsetLeft) / scaleX;
            let newTemplateY = (newScreenY - svgDimensions.offsetTop) / scaleY;

            // Constrain to SVG viewBox dimensions (0 to viewBoxWidth/Height)
            newTemplateX = Math.max(0, Math.min(newTemplateX, svgDimensions.viewBoxWidth));
            newTemplateY = Math.max(0, Math.min(newTemplateY, svgDimensions.viewBoxHeight));

            onKeywordMove(draggingKeyword, newTemplateX, newTemplateY);
        };

        const handleMouseUp = () => {
            setDraggingKeyword(null);
        };

        window.addEventListener("mousemove", handleMouseMove);
        window.addEventListener("mouseup", handleMouseUp);

        return () => {
            window.removeEventListener("mousemove", handleMouseMove);
            window.removeEventListener("mouseup", handleMouseUp);
        };
    }, [draggingKeyword, dragOffset, svgDimensions, onKeywordMove, containerRef]);

    if (svgDimensions.width === 0 || svgDimensions.height === 0) {
        return null;
    }

    // Calculate scale factors to convert template coordinates to screen coordinates
    const scaleX = svgDimensions.width / svgDimensions.viewBoxWidth;
    const scaleY = svgDimensions.height / svgDimensions.viewBoxHeight;

    return (
        <div className="absolute inset-0 pointer-events-auto" style={{ zIndex: 10 }}>
            {detectedKeywords.map((kw, idx) => {
                // Convert template coordinates to screen coordinates
                const screenX = kw.x * scaleX + svgDimensions.offsetLeft;
                const screenY = kw.y * scaleY + svgDimensions.offsetTop;

                // Estimate box dimensions
                const boxWidth = 200 * scaleX;
                const boxHeight = 40 * scaleY;

                return (
                    <div
                        key={`${kw.keyword}-${idx}`}
                        className={`absolute cursor-move group ${
                            draggingKeyword === kw.keyword ? "z-20" : "z-10"
                        }`}
                        style={{
                            left: `${screenX - boxWidth / 2}px`,
                            top: `${screenY - boxHeight / 2}px`,
                            width: `${boxWidth}px`,
                            height: `${boxHeight}px`,
                            border:
                                draggingKeyword === kw.keyword
                                    ? "3px solid #3b82f6"
                                    : "3px dashed #ef4444",
                            backgroundColor:
                                draggingKeyword === kw.keyword
                                    ? "rgba(59, 130, 246, 0.3)"
                                    : "rgba(239, 68, 68, 0.2)",
                            borderRadius: "4px",
                            transition: draggingKeyword ? "none" : "all 0.2s ease",
                        }}
                        onMouseDown={(e) => handleMouseDown(e, kw.keyword, kw.x, kw.y)}
                    >
                        <div
                            className={`absolute left-0 text-white text-[10px] px-1.5 py-0.5 rounded font-mono whitespace-nowrap transition-colors ${
                                draggingKeyword === kw.keyword ? "bg-blue-600" : "bg-red-600"
                            }`}
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
    fontConfig?: CertificateFontConfig;
    previewData?: {
        name?: string;
        eventName?: string;
        academicInstitution?: string;
        certificateTitle?: string;
        certificateSubtitle?: string;
    };
    onKeywordMove?: (keyword: string, x: number, y: number) => void;
}

export const CertificatePreview = ({
    svgPreview,
    imageUrl,
    alt = "Certificate preview",
    detectedKeywords = [],
    availableKeywords = [],
    fontConfig,
    previewData,
    onKeywordMove,
}: CertificatePreviewProps) => {
    const { t } = useTranslation();
    const [showPositions, setShowPositions] = useState(true);
    const [showActualContent, setShowActualContent] = useState(true);
    const containerRef = useRef<HTMLDivElement>(null);
    const { fontFamilies } = useCertificateFontFamilies();

    const handleTogglePositions = useCallback((checked: boolean) => {
        setShowPositions(checked);
    }, []);

    const handleToggleActualContent = useCallback((checked: boolean) => {
        setShowActualContent(checked);
    }, []);

    const handleKeywordMove = useCallback(
        (keyword: string, x: number, y: number) => {
            onKeywordMove?.(keyword, x, y);
        },
        [onKeywordMove],
    );

    const processedSvg = useMemo(() => {
        if (!svgPreview || !showActualContent || !fontConfig) return svgPreview;

        // Map fontFamilies to the format expected by processCertificateSVG
        const mappedFontFamilies = fontFamilies.map((f) => ({
            id: f.id!,
            css_font_name: f.css_font_name!,
            font_family_name: f.font_family_name!,
        }));

        return processCertificateSVG(
            svgPreview,
            fontConfig,
            detectedKeywords,
            mappedFontFamilies,
            previewData || {},
        );
    }, [svgPreview, showActualContent, fontConfig, detectedKeywords, fontFamilies, previewData]);

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
        const contentToRender = processedSvg || svgPreview || imageUrl;

        if (processedSvg || svgPreview) {
            return (
                <div
                    className="w-full [&>svg]:w-full [&>svg]:h-auto [&>svg]:min-h-[300px] [&>svg]:object-contain"
                    dangerouslySetInnerHTML={{
                        __html: DOMPurify.sanitize(contentToRender!, {
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
                                "visibility",
                                "font-family",
                                "font-size",
                                "font-weight",
                                "fill",
                                "text-anchor",
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
                    <div className="flex items-center gap-4">
                        {svgPreview && fontConfig && (
                            <div className="flex items-center gap-2">
                                <Switch
                                    id="show-actual-content"
                                    checked={showActualContent}
                                    onCheckedChange={handleToggleActualContent}
                                />
                                <Label
                                    htmlFor="show-actual-content"
                                    className="text-xs cursor-pointer flex items-center gap-1"
                                >
                                    <Type className="w-3 h-3" />
                                    {t(
                                        "certificateSettings.step2.preview.showActualContent",
                                        "Show Actual Content",
                                    )}
                                </Label>
                            </div>
                        )}
                        {detectedKeywords.length > 0 && (
                            <div className="flex items-center gap-2">
                                <Switch
                                    id="show-positions"
                                    checked={showPositions}
                                    onCheckedChange={handleTogglePositions}
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
                                onKeywordMove={handleKeywordMove}
                                svgContent={processedSvg || svgPreview}
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
                                                "certificateSettings.step2.keywordStatus.table.status",
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
                                                {kw.isDetected ? (
                                                    <span className="px-2 py-1 bg-green-100 border border-green-300 text-green-900 rounded text-xs font-semibold">
                                                        {t(
                                                            "certificateSettings.step2.keywordStatus.detected",
                                                        )}
                                                    </span>
                                                ) : (
                                                    <span className="px-2 py-1 bg-gray-100 border border-gray-300 text-gray-900 rounded text-xs font-semibold">
                                                        {t(
                                                            "certificateSettings.step2.keywordStatus.notDetected",
                                                        )}
                                                    </span>
                                                )}
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
