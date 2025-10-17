import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Upload, Info, Image as ImageIcon } from "lucide-react";
import type { AvailableKeyword } from "@/hooks/useCertificateTemplate";

export interface CertificateTemplateUploadProps {
    svgFile: File | null;
    availableKeywords: AvailableKeyword[];
    onFileSelect: (event: React.ChangeEvent<HTMLInputElement>) => void;
    fileInputRef: React.RefObject<HTMLInputElement | null>;
}

export const CertificateTemplateUpload = ({
    svgFile,
    availableKeywords,
    onFileSelect,
    fileInputRef,
}: CertificateTemplateUploadProps) => {
    const { t } = useTranslation();

    return (
        <div className="space-y-6 rounded-lg border p-6">
            {/* Instructions */}
            <div className="space-y-4">
                <Typography variant="header" tag="h3" className="text-base font-semibold">
                    {t("certificateSettings.step2.instructions.title")}
                </Typography>

                {/* Instruction steps with mockup images */}
                <div className="space-y-4">
                    {[1, 2, 3].map((step) => (
                        <div key={step} className="flex gap-4">
                            <div className="flex-shrink-0 w-32 h-24 bg-muted rounded-md flex items-center justify-center border-2 border-dashed">
                                <ImageIcon className="h-8 w-8 text-muted-foreground" />
                            </div>
                            <div className="flex-1">
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-sm font-medium mb-1"
                                >
                                    {t(`certificateSettings.step2.instructions.step${step}.title`)}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-xs text-muted-foreground"
                                >
                                    {t(
                                        `certificateSettings.step2.instructions.step${step}.description`,
                                    )}
                                </Typography>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Available Keywords Info */}
            <Alert variant="info">
                <Info className="h-4 w-4" />
                <AlertTitle>{t("certificateSettings.step2.keywords.title")}</AlertTitle>
                <AlertDescription>
                    <div className="mt-2">
                        <Typography variant="text" tag="p" className="text-xs mb-2 text-blue-700">
                            {t("certificateSettings.step2.keywords.description")}
                        </Typography>
                        <div className="flex flex-wrap gap-2">
                            {availableKeywords.map((kw) => (
                                <div key={kw.keyword} className="inline-flex items-center gap-1">
                                    <code className="px-2 py-1 bg-blue-100 text-blue-700 dark:text-blue-400 rounded text-xs font-mono">
                                        {kw.keyword}
                                    </code>
                                    {kw.mandatory && (
                                        <span className="px-1.5 py-0.5 bg-red-100 text-red-700 rounded text-[10px] font-semibold">
                                            {t("common.required")}
                                        </span>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                </AlertDescription>
            </Alert>

            {/* SVG Upload */}
            <div className="space-y-2">
                <Label htmlFor="svg-upload">
                    <Typography variant="text" tag="span" className="text-sm font-medium">
                        {t("certificateSettings.step2.upload.label")}
                    </Typography>
                </Label>
                <div className="flex gap-2">
                    <input
                        id="svg-upload"
                        type="file"
                        ref={fileInputRef}
                        accept=".svg,image/svg+xml"
                        onChange={onFileSelect}
                        className="hidden"
                    />
                    <Button
                        type="button"
                        variant="secondary-light"
                        onClick={() => fileInputRef.current?.click()}
                        className="w-full h-30"
                    >
                        <Upload className="h-4 w-4 mr-2" />
                        <Typography
                            variant="text"
                            tag="span"
                            className="font-medium text-foreground"
                        >
                            {svgFile ? svgFile.name : t("certificateSettings.step2.upload.button")}
                        </Typography>
                    </Button>
                </div>
                <Typography variant="text" tag="p" className="text-xs text-muted-foreground">
                    {t("certificateSettings.step2.upload.hint")}
                </Typography>
            </div>
        </div>
    );
};
