import React from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import type { CoreApiInternalHandlerEventconfigEventCertificateConfigResponse } from "@decm/api";
import { useCertificateFontFamilies } from "@/hooks/events/useCertificateFontFamilies";
import type { DetectedKeyword } from "@/hooks/useCertificateTemplate";

interface CertificateFontSettingsProps {
    eventCertificateConfig?: CoreApiInternalHandlerEventconfigEventCertificateConfigResponse;
    detectedKeywords?: DetectedKeyword[];
    onChange?: (fontConfig: CertificateFontConfig) => void;
}

export interface CertificateFontConfig {
    event_name_font_family_id?: number;
    event_name_font_weight?: number;
    name_font_family_id?: number;
    name_font_weight?: number;
    academic_institution_font_family_id?: number;
    academic_institution_font_weight?: number;
    certificate_title_font_family_id?: number;
    certificate_title_font_weight?: number;
    certificate_subtitle_font_family_id?: number;
    certificate_subtitle_font_weight?: number;
}

export const CertificateFontSettings: React.FC<CertificateFontSettingsProps> = ({
    eventCertificateConfig,
    detectedKeywords = [],
    onChange,
}) => {
    const { t } = useTranslation();
    const { fontFamilies, isLoading: isLoadingFonts } = useCertificateFontFamilies();

    // Get default font family ID (fallback to first font if no default found)
    const defaultFontFamilyId = React.useMemo(() => {
        const defaultFont = fontFamilies.find((f) => f.is_default);
        return defaultFont?.id || fontFamilies[0]?.id || 1;
    }, [fontFamilies]);

    // Initialize state from existing config or defaults
    const [fontConfig, setFontConfig] = React.useState<CertificateFontConfig>({
        event_name_font_family_id:
            eventCertificateConfig?.event_name_font_family_id || defaultFontFamilyId,
        event_name_font_weight: eventCertificateConfig?.event_name_font_weight || 700,
        name_font_family_id: eventCertificateConfig?.name_font_family_id || defaultFontFamilyId,
        name_font_weight: eventCertificateConfig?.name_font_weight || 700,
        academic_institution_font_family_id:
            eventCertificateConfig?.academic_institution_font_family_id || defaultFontFamilyId,
        academic_institution_font_weight:
            eventCertificateConfig?.academic_institution_font_weight || 700,
        certificate_title_font_family_id:
            eventCertificateConfig?.certificate_title_font_family_id || defaultFontFamilyId,
        certificate_title_font_weight: eventCertificateConfig?.certificate_title_font_weight || 700,
        certificate_subtitle_font_family_id:
            eventCertificateConfig?.certificate_subtitle_font_family_id || defaultFontFamilyId,
        certificate_subtitle_font_weight:
            eventCertificateConfig?.certificate_subtitle_font_weight || 700,
    });

    // Update state when config or default font changes
    React.useEffect(() => {
        if (eventCertificateConfig) {
            const newConfig = {
                event_name_font_family_id:
                    eventCertificateConfig.event_name_font_family_id || defaultFontFamilyId,
                event_name_font_weight: eventCertificateConfig.event_name_font_weight || 700,
                name_font_family_id:
                    eventCertificateConfig.name_font_family_id || defaultFontFamilyId,
                name_font_weight: eventCertificateConfig.name_font_weight || 700,
                academic_institution_font_family_id:
                    eventCertificateConfig.academic_institution_font_family_id ||
                    defaultFontFamilyId,
                academic_institution_font_weight:
                    eventCertificateConfig.academic_institution_font_weight || 700,
                certificate_title_font_family_id:
                    eventCertificateConfig.certificate_title_font_family_id || defaultFontFamilyId,
                certificate_title_font_weight:
                    eventCertificateConfig.certificate_title_font_weight || 700,
                certificate_subtitle_font_family_id:
                    eventCertificateConfig.certificate_subtitle_font_family_id ||
                    defaultFontFamilyId,
                certificate_subtitle_font_weight:
                    eventCertificateConfig.certificate_subtitle_font_weight || 700,
            };
            setFontConfig(newConfig);
        } else if (!isLoadingFonts && defaultFontFamilyId && fontFamilies.length > 0) {
            // When there's no config but we have detected keywords and fonts are loaded,
            // ensure all fields have default values
            setFontConfig((prevConfig) => {
                const needsUpdate =
                    !prevConfig.event_name_font_family_id ||
                    !prevConfig.name_font_family_id ||
                    !prevConfig.academic_institution_font_family_id ||
                    !prevConfig.certificate_title_font_family_id ||
                    !prevConfig.certificate_subtitle_font_family_id;

                if (!needsUpdate) return prevConfig;

                return {
                    event_name_font_family_id:
                        prevConfig.event_name_font_family_id || defaultFontFamilyId,
                    event_name_font_weight: prevConfig.event_name_font_weight || 700,
                    name_font_family_id: prevConfig.name_font_family_id || defaultFontFamilyId,
                    name_font_weight: prevConfig.name_font_weight || 700,
                    academic_institution_font_family_id:
                        prevConfig.academic_institution_font_family_id || defaultFontFamilyId,
                    academic_institution_font_weight:
                        prevConfig.academic_institution_font_weight || 700,
                    certificate_title_font_family_id:
                        prevConfig.certificate_title_font_family_id || defaultFontFamilyId,
                    certificate_title_font_weight: prevConfig.certificate_title_font_weight || 700,
                    certificate_subtitle_font_family_id:
                        prevConfig.certificate_subtitle_font_family_id || defaultFontFamilyId,
                    certificate_subtitle_font_weight:
                        prevConfig.certificate_subtitle_font_weight || 700,
                };
            });
        }
    }, [eventCertificateConfig, defaultFontFamilyId, isLoadingFonts, fontFamilies.length]);

    // Notify parent of changes
    React.useEffect(() => {
        onChange?.(fontConfig);
    }, [fontConfig, onChange]);

    // Get available weights for a given font family ID
    const getAvailableWeights = (fontFamilyId?: number) => {
        if (!fontFamilyId || !fontFamilies.length) return [];
        const fontFamily = fontFamilies.find((f) => f.id === fontFamilyId);
        return fontFamily?.available_font_weights?.filter((w) => w != null) || [];
    };

    // Check if fields have positions configured (meaning they appear in template)
    // First check eventCertificateConfig, then fall back to detectedKeywords
    // Only show if positions exist and are not (0,0)
    const hasEventName =
        (eventCertificateConfig?.event_name_pos_x != null &&
            eventCertificateConfig?.event_name_pos_y != null &&
            (eventCertificateConfig.event_name_pos_x !== 0 ||
                eventCertificateConfig.event_name_pos_y !== 0)) ||
        detectedKeywords.some((k) => k.keyword === "{{ eventName }}" && (k.x !== 0 || k.y !== 0));
    const hasParticipantName =
        (eventCertificateConfig?.name_pos_x != null &&
            eventCertificateConfig?.name_pos_y != null) ||
        detectedKeywords.some((k) => k.keyword === "{{ name }}");
    const hasAcademicInstitution =
        (eventCertificateConfig?.academic_institution_pos_x != null &&
            eventCertificateConfig?.academic_institution_pos_y != null) ||
        detectedKeywords.some((k) => k.keyword === "{{ academicInstitutionName }}");
    const hasCertificateTitle =
        (eventCertificateConfig?.certificate_title_pos_x != null &&
            eventCertificateConfig?.certificate_title_pos_y != null) ||
        detectedKeywords.some((k) => k.keyword === "{{ certificateTitle }}");
    const hasCertificateSubtitle =
        (eventCertificateConfig?.certificate_subtitle_pos_x != null &&
            eventCertificateConfig?.certificate_subtitle_pos_y != null) ||
        detectedKeywords.some((k) => k.keyword === "{{ certificateSubtitle }}");

    const hasAnyFields =
        hasEventName ||
        hasParticipantName ||
        hasAcademicInstitution ||
        hasCertificateTitle ||
        hasCertificateSubtitle;

    if (isLoadingFonts) {
        return (
            <div className="space-y-4 rounded-lg border p-4">
                <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                    {t("common.loading", "Loading...")}
                </Typography>
            </div>
        );
    }

    return (
        <div className="space-y-3 rounded-lg border p-3">
            <div className="space-y-0.5">
                <Typography variant="header" tag="h3" className="text-base font-semibold">
                    {t("certificateSettings.fontSettings.title", "Font Settings")}
                </Typography>
                <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                    {t(
                        "certificateSettings.fontSettings.description",
                        "Configure font family and weight for each text field in the certificate template.",
                    )}
                </Typography>
            </div>

            <div className="space-y-3">
                {!hasAnyFields && (
                    <div className="rounded-md bg-muted/50 p-3">
                        <Typography
                            variant="text"
                            tag="p"
                            className="text-sm text-muted-foreground"
                        >
                            {t(
                                "certificateSettings.fontSettings.noFieldsConfigured",
                                "No text fields have been positioned in the certificate template yet. Configure text positions first to manage font settings.",
                            )}
                        </Typography>
                    </div>
                )}

                {/* Event Name Font - Only show if configured in template */}
                {hasEventName && (
                    <div className="space-y-2">
                        <Typography variant="text" tag="h4" className="text-sm font-medium">
                            {t("certificateSettings.fontSettings.eventName", "Event Name")}
                        </Typography>
                        <div className="grid grid-cols-2 gap-2">
                            <div className="space-y-1">
                                <Label htmlFor="event-name-font-family" className="text-xs">
                                    {t(
                                        "certificateSettings.fontSettings.fontFamily",
                                        "Font Family",
                                    )}
                                </Label>
                                <Select
                                    value={
                                        fontConfig.event_name_font_family_id?.toString() ||
                                        defaultFontFamilyId.toString()
                                    }
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            event_name_font_family_id: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="event-name-font-family">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectFont",
                                                "Select font",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {fontFamilies
                                            .filter((font) => font?.id != null)
                                            .map((font) => (
                                                <SelectItem
                                                    key={font.id}
                                                    value={font.id.toString()}
                                                >
                                                    {font.font_family_name}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="space-y-1">
                                <Label htmlFor="event-name-font-weight" className="text-xs">
                                    {t(
                                        "certificateSettings.fontSettings.fontWeight",
                                        "Font Weight",
                                    )}
                                </Label>
                                <Select
                                    value={fontConfig.event_name_font_weight?.toString() || "700"}
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            event_name_font_weight: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="event-name-font-weight">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectWeight",
                                                "Select weight",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {getAvailableWeights(fontConfig.event_name_font_family_id)
                                            .filter((weight) => weight != null)
                                            .map((weight) => (
                                                <SelectItem key={weight} value={weight.toString()}>
                                                    {weight}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                    </div>
                )}

                {/* Participant Name Font - Only show if configured in template */}
                {hasParticipantName && (
                    <div className="space-y-2">
                        <Typography variant="text" tag="h4" className="text-sm font-medium">
                            {t(
                                "certificateSettings.fontSettings.participantName",
                                "Participant Name",
                            )}
                        </Typography>
                        <div className="grid grid-cols-2 gap-2">
                            <div className="space-y-1">
                                <Label htmlFor="name-font-family" className="text-xs">
                                    {t(
                                        "certificateSettings.fontSettings.fontFamily",
                                        "Font Family",
                                    )}
                                </Label>
                                <Select
                                    value={
                                        fontConfig.name_font_family_id?.toString() ||
                                        defaultFontFamilyId.toString()
                                    }
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            name_font_family_id: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="name-font-family">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectFont",
                                                "Select font",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {fontFamilies
                                            .filter((font) => font?.id != null)
                                            .map((font) => (
                                                <SelectItem
                                                    key={font.id}
                                                    value={font.id.toString()}
                                                >
                                                    {font.font_family_name}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="space-y-1">
                                <Label htmlFor="name-font-weight" className="text-xs">
                                    {t(
                                        "certificateSettings.fontSettings.fontWeight",
                                        "Font Weight",
                                    )}
                                </Label>
                                <Select
                                    value={fontConfig.name_font_weight?.toString() || "700"}
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            name_font_weight: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="name-font-weight">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectWeight",
                                                "Select weight",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {getAvailableWeights(fontConfig.name_font_family_id)
                                            .filter((weight) => weight != null)
                                            .map((weight) => (
                                                <SelectItem key={weight} value={weight.toString()}>
                                                    {weight}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                    </div>
                )}

                {/* Academic Institution Font - Only show if configured in template */}
                {hasAcademicInstitution && (
                    <div className="space-y-2">
                        <Typography variant="text" tag="h4" className="text-sm font-medium">
                            {t(
                                "certificateSettings.fontSettings.academicInstitution",
                                "Academic Institution",
                            )}
                        </Typography>
                        <div className="grid grid-cols-2 gap-2">
                            <div className="space-y-1">
                                <Label
                                    htmlFor="academic-institution-font-family"
                                    className="text-xs"
                                >
                                    {t(
                                        "certificateSettings.fontSettings.fontFamily",
                                        "Font Family",
                                    )}
                                </Label>
                                <Select
                                    value={
                                        fontConfig.academic_institution_font_family_id?.toString() ||
                                        defaultFontFamilyId.toString()
                                    }
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            academic_institution_font_family_id: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="academic-institution-font-family">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectFont",
                                                "Select font",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {fontFamilies
                                            .filter((font) => font?.id != null)
                                            .map((font) => (
                                                <SelectItem
                                                    key={font.id}
                                                    value={font.id.toString()}
                                                >
                                                    {font.font_family_name}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="space-y-1">
                                <Label
                                    htmlFor="academic-institution-font-weight"
                                    className="text-xs"
                                >
                                    {t(
                                        "certificateSettings.fontSettings.fontWeight",
                                        "Font Weight",
                                    )}
                                </Label>
                                <Select
                                    value={
                                        fontConfig.academic_institution_font_weight?.toString() ||
                                        "700"
                                    }
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            academic_institution_font_weight: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="academic-institution-font-weight">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectWeight",
                                                "Select weight",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {getAvailableWeights(
                                            fontConfig.academic_institution_font_family_id,
                                        )
                                            .filter((weight) => weight != null)
                                            .map((weight) => (
                                                <SelectItem key={weight} value={weight.toString()}>
                                                    {weight}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                    </div>
                )}

                {/* Certificate Title Font - Only show if configured in template */}
                {hasCertificateTitle && (
                    <div className="space-y-2">
                        <Typography variant="text" tag="h4" className="text-sm font-medium">
                            {t(
                                "certificateSettings.fontSettings.certificateTitle",
                                "Certificate Title",
                            )}
                        </Typography>
                        <div className="grid grid-cols-2 gap-2">
                            <div className="space-y-1">
                                <Label htmlFor="certificate-title-font-family" className="text-xs">
                                    {t(
                                        "certificateSettings.fontSettings.fontFamily",
                                        "Font Family",
                                    )}
                                </Label>
                                <Select
                                    value={
                                        fontConfig.certificate_title_font_family_id?.toString() ||
                                        defaultFontFamilyId.toString()
                                    }
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            certificate_title_font_family_id: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="certificate-title-font-family">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectFont",
                                                "Select font",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {fontFamilies
                                            .filter((font) => font?.id != null)
                                            .map((font) => (
                                                <SelectItem
                                                    key={font.id}
                                                    value={font.id.toString()}
                                                >
                                                    {font.font_family_name}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="space-y-1">
                                <Label htmlFor="certificate-title-font-weight" className="text-xs">
                                    {t(
                                        "certificateSettings.fontSettings.fontWeight",
                                        "Font Weight",
                                    )}
                                </Label>
                                <Select
                                    value={
                                        fontConfig.certificate_title_font_weight?.toString() ||
                                        "700"
                                    }
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            certificate_title_font_weight: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="certificate-title-font-weight">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectWeight",
                                                "Select weight",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {getAvailableWeights(
                                            fontConfig.certificate_title_font_family_id,
                                        )
                                            .filter((weight) => weight != null)
                                            .map((weight) => (
                                                <SelectItem key={weight} value={weight.toString()}>
                                                    {weight}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                    </div>
                )}

                {/* Certificate Subtitle Font - Only show if configured in template */}
                {hasCertificateSubtitle && (
                    <div className="space-y-2">
                        <Typography variant="text" tag="h4" className="text-sm font-medium">
                            {t(
                                "certificateSettings.fontSettings.certificateSubtitle",
                                "Certificate Subtitle",
                            )}
                        </Typography>
                        <div className="grid grid-cols-2 gap-2">
                            <div className="space-y-1">
                                <Label
                                    htmlFor="certificate-subtitle-font-family"
                                    className="text-xs"
                                >
                                    {t(
                                        "certificateSettings.fontSettings.fontFamily",
                                        "Font Family",
                                    )}
                                </Label>
                                <Select
                                    value={
                                        fontConfig.certificate_subtitle_font_family_id?.toString() ||
                                        defaultFontFamilyId.toString()
                                    }
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            certificate_subtitle_font_family_id: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="certificate-subtitle-font-family">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectFont",
                                                "Select font",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {fontFamilies
                                            .filter((font) => font?.id != null)
                                            .map((font) => (
                                                <SelectItem
                                                    key={font.id}
                                                    value={font.id.toString()}
                                                >
                                                    {font.font_family_name}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="space-y-1">
                                <Label
                                    htmlFor="certificate-subtitle-font-weight"
                                    className="text-xs"
                                >
                                    {t(
                                        "certificateSettings.fontSettings.fontWeight",
                                        "Font Weight",
                                    )}
                                </Label>
                                <Select
                                    value={
                                        fontConfig.certificate_subtitle_font_weight?.toString() ||
                                        "700"
                                    }
                                    onValueChange={(value) =>
                                        setFontConfig({
                                            ...fontConfig,
                                            certificate_subtitle_font_weight: parseInt(value),
                                        })
                                    }
                                >
                                    <SelectTrigger id="certificate-subtitle-font-weight">
                                        <SelectValue
                                            placeholder={t(
                                                "certificateSettings.fontSettings.selectWeight",
                                                "Select weight",
                                            )}
                                        />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {getAvailableWeights(
                                            fontConfig.certificate_subtitle_font_family_id,
                                        )
                                            .filter((weight) => weight != null)
                                            .map((weight) => (
                                                <SelectItem key={weight} value={weight.toString()}>
                                                    {weight}
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};
