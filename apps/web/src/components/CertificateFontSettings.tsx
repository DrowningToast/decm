import React from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import type { CoreApiInternalHandlerEventconfigEventCertificateConfigResponse } from "@decm/api";

interface CertificateFontSettingsProps {
    eventCertificateConfig?: CoreApiInternalHandlerEventconfigEventCertificateConfigResponse;
    onUpdate: (fontConfig: CertificateFontConfig) => Promise<void>;
    isUpdating?: boolean;
}

export interface CertificateFontConfig {
    event_name_font_family?: string;
    event_name_font_weight?: number;
    name_font_family?: string;
    name_font_weight?: number;
    academic_institution_font_family?: string;
    academic_institution_font_weight?: number;
    certificate_title_font_family?: string;
    certificate_title_font_weight?: number;
    certificate_subtitle_font_family?: string;
    certificate_subtitle_font_weight?: number;
}

// Available font families
const FONT_FAMILIES = [
    { value: "Prompt", label: "Prompt" },
    { value: "Sarabun", label: "Sarabun" },
    { value: "Kanit", label: "Kanit" },
    { value: "Arial", label: "Arial" },
    { value: "Helvetica", label: "Helvetica" },
    { value: "Times New Roman", label: "Times New Roman" },
    { value: "Georgia", label: "Georgia" },
    { value: "Verdana", label: "Verdana" },
    { value: "Courier New", label: "Courier New" },
] as const;

// Available font weights
const FONT_WEIGHTS = [
    { value: 100, label: "100 - Thin" },
    { value: 200, label: "200 - Extra Light" },
    { value: 300, label: "300 - Light" },
    { value: 400, label: "400 - Normal" },
    { value: 500, label: "500 - Medium" },
    { value: 600, label: "600 - Semi Bold" },
    { value: 700, label: "700 - Bold" },
    { value: 800, label: "800 - Extra Bold" },
    { value: 900, label: "900 - Black" },
] as const;

export const CertificateFontSettings: React.FC<CertificateFontSettingsProps> = ({
    eventCertificateConfig,
    onUpdate,
    isUpdating = false,
}) => {
    const { t } = useTranslation();

    // Initialize state from existing config or defaults
    const [fontConfig, setFontConfig] = React.useState<CertificateFontConfig>({
        event_name_font_family: eventCertificateConfig?.event_name_font_family || "Prompt",
        event_name_font_weight: eventCertificateConfig?.event_name_font_weight || 700,
        name_font_family: eventCertificateConfig?.name_font_family || "Prompt",
        name_font_weight: eventCertificateConfig?.name_font_weight || 700,
        academic_institution_font_family:
            eventCertificateConfig?.academic_institution_font_family || "Prompt",
        academic_institution_font_weight:
            eventCertificateConfig?.academic_institution_font_weight || 700,
        certificate_title_font_family:
            eventCertificateConfig?.certificate_title_font_family || "Prompt",
        certificate_title_font_weight: eventCertificateConfig?.certificate_title_font_weight || 700,
        certificate_subtitle_font_family:
            eventCertificateConfig?.certificate_subtitle_font_family || "Prompt",
        certificate_subtitle_font_weight:
            eventCertificateConfig?.certificate_subtitle_font_weight || 700,
    });

    // Update state when config changes
    React.useEffect(() => {
        if (eventCertificateConfig) {
            setFontConfig({
                event_name_font_family: eventCertificateConfig.event_name_font_family || "Prompt",
                event_name_font_weight: eventCertificateConfig.event_name_font_weight || 700,
                name_font_family: eventCertificateConfig.name_font_family || "Prompt",
                name_font_weight: eventCertificateConfig.name_font_weight || 700,
                academic_institution_font_family:
                    eventCertificateConfig.academic_institution_font_family || "Prompt",
                academic_institution_font_weight:
                    eventCertificateConfig.academic_institution_font_weight || 700,
                certificate_title_font_family:
                    eventCertificateConfig.certificate_title_font_family || "Prompt",
                certificate_title_font_weight:
                    eventCertificateConfig.certificate_title_font_weight || 700,
                certificate_subtitle_font_family:
                    eventCertificateConfig.certificate_subtitle_font_family || "Prompt",
                certificate_subtitle_font_weight:
                    eventCertificateConfig.certificate_subtitle_font_weight || 700,
            });
        }
    }, [eventCertificateConfig]);

    const handleSave = async () => {
        await onUpdate(fontConfig);
    };

    return (
        <div className="space-y-6 rounded-lg border p-6">
            <div className="space-y-2">
                <Typography variant="header" tag="h3" className="text-lg font-semibold">
                    {t("certificateSettings.fontSettings.title", "Font Settings")}
                </Typography>
                <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                    {t(
                        "certificateSettings.fontSettings.description",
                        "Configure font family and weight for each text field in the certificate template.",
                    )}
                </Typography>
            </div>

            <div className="space-y-6">
                {/* Event Name Font */}
                <div className="space-y-4">
                    <Typography variant="text" tag="h4" className="text-sm font-medium">
                        {t("certificateSettings.fontSettings.eventName", "Event Name")}
                    </Typography>
                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label htmlFor="event-name-font-family">
                                {t("certificateSettings.fontSettings.fontFamily", "Font Family")}
                            </Label>
                            <Select
                                value={fontConfig.event_name_font_family}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        event_name_font_family: value,
                                    })
                                }
                            >
                                <SelectTrigger id="event-name-font-family">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_FAMILIES.map((font) => (
                                        <SelectItem key={font.value} value={font.value}>
                                            {font.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="event-name-font-weight">
                                {t("certificateSettings.fontSettings.fontWeight", "Font Weight")}
                            </Label>
                            <Select
                                value={fontConfig.event_name_font_weight?.toString()}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        event_name_font_weight: parseInt(value),
                                    })
                                }
                            >
                                <SelectTrigger id="event-name-font-weight">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_WEIGHTS.map((weight) => (
                                        <SelectItem
                                            key={weight.value}
                                            value={weight.value.toString()}
                                        >
                                            {weight.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                </div>

                {/* Participant Name Font */}
                <div className="space-y-4">
                    <Typography variant="text" tag="h4" className="text-sm font-medium">
                        {t("certificateSettings.fontSettings.participantName", "Participant Name")}
                    </Typography>
                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label htmlFor="name-font-family">
                                {t("certificateSettings.fontSettings.fontFamily", "Font Family")}
                            </Label>
                            <Select
                                value={fontConfig.name_font_family}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        name_font_family: value,
                                    })
                                }
                            >
                                <SelectTrigger id="name-font-family">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_FAMILIES.map((font) => (
                                        <SelectItem key={font.value} value={font.value}>
                                            {font.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="name-font-weight">
                                {t("certificateSettings.fontSettings.fontWeight", "Font Weight")}
                            </Label>
                            <Select
                                value={fontConfig.name_font_weight?.toString()}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        name_font_weight: parseInt(value),
                                    })
                                }
                            >
                                <SelectTrigger id="name-font-weight">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_WEIGHTS.map((weight) => (
                                        <SelectItem
                                            key={weight.value}
                                            value={weight.value.toString()}
                                        >
                                            {weight.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                </div>

                {/* Academic Institution Font */}
                <div className="space-y-4">
                    <Typography variant="text" tag="h4" className="text-sm font-medium">
                        {t(
                            "certificateSettings.fontSettings.academicInstitution",
                            "Academic Institution",
                        )}
                    </Typography>
                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label htmlFor="academic-institution-font-family">
                                {t("certificateSettings.fontSettings.fontFamily", "Font Family")}
                            </Label>
                            <Select
                                value={fontConfig.academic_institution_font_family}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        academic_institution_font_family: value,
                                    })
                                }
                            >
                                <SelectTrigger id="academic-institution-font-family">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_FAMILIES.map((font) => (
                                        <SelectItem key={font.value} value={font.value}>
                                            {font.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="academic-institution-font-weight">
                                {t("certificateSettings.fontSettings.fontWeight", "Font Weight")}
                            </Label>
                            <Select
                                value={fontConfig.academic_institution_font_weight?.toString()}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        academic_institution_font_weight: parseInt(value),
                                    })
                                }
                            >
                                <SelectTrigger id="academic-institution-font-weight">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_WEIGHTS.map((weight) => (
                                        <SelectItem
                                            key={weight.value}
                                            value={weight.value.toString()}
                                        >
                                            {weight.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                </div>

                {/* Certificate Title Font */}
                <div className="space-y-4">
                    <Typography variant="text" tag="h4" className="text-sm font-medium">
                        {t(
                            "certificateSettings.fontSettings.certificateTitle",
                            "Certificate Title",
                        )}
                    </Typography>
                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label htmlFor="certificate-title-font-family">
                                {t("certificateSettings.fontSettings.fontFamily", "Font Family")}
                            </Label>
                            <Select
                                value={fontConfig.certificate_title_font_family}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        certificate_title_font_family: value,
                                    })
                                }
                            >
                                <SelectTrigger id="certificate-title-font-family">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_FAMILIES.map((font) => (
                                        <SelectItem key={font.value} value={font.value}>
                                            {font.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="certificate-title-font-weight">
                                {t("certificateSettings.fontSettings.fontWeight", "Font Weight")}
                            </Label>
                            <Select
                                value={fontConfig.certificate_title_font_weight?.toString()}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        certificate_title_font_weight: parseInt(value),
                                    })
                                }
                            >
                                <SelectTrigger id="certificate-title-font-weight">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_WEIGHTS.map((weight) => (
                                        <SelectItem
                                            key={weight.value}
                                            value={weight.value.toString()}
                                        >
                                            {weight.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                </div>

                {/* Certificate Subtitle Font */}
                <div className="space-y-4">
                    <Typography variant="text" tag="h4" className="text-sm font-medium">
                        {t(
                            "certificateSettings.fontSettings.certificateSubtitle",
                            "Certificate Subtitle",
                        )}
                    </Typography>
                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label htmlFor="certificate-subtitle-font-family">
                                {t("certificateSettings.fontSettings.fontFamily", "Font Family")}
                            </Label>
                            <Select
                                value={fontConfig.certificate_subtitle_font_family}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        certificate_subtitle_font_family: value,
                                    })
                                }
                            >
                                <SelectTrigger id="certificate-subtitle-font-family">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_FAMILIES.map((font) => (
                                        <SelectItem key={font.value} value={font.value}>
                                            {font.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="certificate-subtitle-font-weight">
                                {t("certificateSettings.fontSettings.fontWeight", "Font Weight")}
                            </Label>
                            <Select
                                value={fontConfig.certificate_subtitle_font_weight?.toString()}
                                onValueChange={(value) =>
                                    setFontConfig({
                                        ...fontConfig,
                                        certificate_subtitle_font_weight: parseInt(value),
                                    })
                                }
                            >
                                <SelectTrigger id="certificate-subtitle-font-weight">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {FONT_WEIGHTS.map((weight) => (
                                        <SelectItem
                                            key={weight.value}
                                            value={weight.value.toString()}
                                        >
                                            {weight.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                </div>
            </div>

            {/* Save Button */}
            <div className="flex justify-end pt-4">
                <Button
                    type="button"
                    variant="primary"
                    onClick={handleSave}
                    disabled={isUpdating}
                    className="min-w-[150px]"
                >
                    <Typography variant="text" tag="span" className="font-medium">
                        {isUpdating
                            ? t("common.loading", "Loading...")
                            : t(
                                  "certificateSettings.fontSettings.saveButton",
                                  "Save Font Settings",
                              )}
                    </Typography>
                </Button>
            </div>
        </div>
    );
};
