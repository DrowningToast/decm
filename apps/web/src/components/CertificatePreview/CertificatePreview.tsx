import { useTranslation } from "react-i18next";
import DOMPurify from "dompurify";
import { Typography } from "@/components/typography/typography";

export interface CertificatePreviewProps {
    svgPreview?: string;
    imageUrl?: string;
    alt?: string;
}

export const CertificatePreview = ({
    svgPreview,
    imageUrl,
    alt = "Certificate preview",
}: CertificatePreviewProps) => {
    const { t } = useTranslation();

    if (!svgPreview && !imageUrl) {
        return null;
    }

    const renderContent = () => {
        if (svgPreview) {
            return (
                <div
                    className="w-full max-w-4xl [&>svg]:w-full [&>svg]:h-auto [&>svg]:max-h-[500px] [&>svg]:object-contain"
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
                    className="w-full max-w-4xl h-auto max-h-[500px] object-contain"
                    onError={() => {
                        console.error("Failed to load certificate image:", imageUrl);
                    }}
                />
            );
        }

        return null;
    };

    return (
        <div className="space-y-2">
            <Typography variant="text" tag="p" className="text-sm font-medium">
                {t("certificateSettings.step2.preview.title")}
            </Typography>
            <div className="rounded-md border bg-muted/30 p-4">
                <div className="w-full flex items-center justify-center">{renderContent()}</div>
            </div>
        </div>
    );
};
