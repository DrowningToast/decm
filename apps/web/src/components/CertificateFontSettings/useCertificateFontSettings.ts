import React from "react";
import { useTranslation } from "react-i18next";
import { useCertificateFontFamilies } from "@/hooks/events/useCertificateFontFamilies";
import type { CertificateFontConfig, CertificateFontSettingsProps } from "./types";

export const useCertificateFontSettings = ({
    eventCertificateConfig,
    detectedKeywords = [],
    onChange,
}: CertificateFontSettingsProps) => {
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

    const fontFields = [
        {
            key: "event-name",
            label: t("certificateSettings.fontSettings.eventName", "Event Name"),
            visible: hasEventName,
            familyId: "event-name-font-family",
            familyValue: fontConfig.event_name_font_family_id,
            weightValue: fontConfig.event_name_font_weight,
            onFamilyChange: (value: string) =>
                setFontConfig({ ...fontConfig, event_name_font_family_id: parseInt(value) }),
            onWeightChange: (value: string) =>
                setFontConfig({ ...fontConfig, event_name_font_weight: parseInt(value) }),
        },
        {
            key: "name",
            label: t("certificateSettings.fontSettings.participantName", "Participant Name"),
            visible: hasParticipantName,
            familyId: "name-font-family",
            familyValue: fontConfig.name_font_family_id,
            weightValue: fontConfig.name_font_weight,
            onFamilyChange: (value: string) =>
                setFontConfig({ ...fontConfig, name_font_family_id: parseInt(value) }),
            onWeightChange: (value: string) =>
                setFontConfig({ ...fontConfig, name_font_weight: parseInt(value) }),
        },
        {
            key: "academic-institution",
            label: t(
                "certificateSettings.fontSettings.academicInstitution",
                "Academic Institution",
            ),
            visible: hasAcademicInstitution,
            familyId: "academic-institution-font-family",
            familyValue: fontConfig.academic_institution_font_family_id,
            weightValue: fontConfig.academic_institution_font_weight,
            onFamilyChange: (value: string) =>
                setFontConfig({
                    ...fontConfig,
                    academic_institution_font_family_id: parseInt(value),
                }),
            onWeightChange: (value: string) =>
                setFontConfig({
                    ...fontConfig,
                    academic_institution_font_weight: parseInt(value),
                }),
        },
        {
            key: "certificate-title",
            label: t("certificateSettings.fontSettings.certificateTitle", "Certificate Title"),
            visible: hasCertificateTitle,
            familyId: "certificate-title-font-family",
            familyValue: fontConfig.certificate_title_font_family_id,
            weightValue: fontConfig.certificate_title_font_weight,
            onFamilyChange: (value: string) =>
                setFontConfig({
                    ...fontConfig,
                    certificate_title_font_family_id: parseInt(value),
                }),
            onWeightChange: (value: string) =>
                setFontConfig({
                    ...fontConfig,
                    certificate_title_font_weight: parseInt(value),
                }),
        },
        {
            key: "certificate-subtitle",
            label: t(
                "certificateSettings.fontSettings.certificateSubtitle",
                "Certificate Subtitle",
            ),
            visible: hasCertificateSubtitle,
            familyId: "certificate-subtitle-font-family",
            familyValue: fontConfig.certificate_subtitle_font_family_id,
            weightValue: fontConfig.certificate_subtitle_font_weight,
            onFamilyChange: (value: string) =>
                setFontConfig({
                    ...fontConfig,
                    certificate_subtitle_font_family_id: parseInt(value),
                }),
            onWeightChange: (value: string) =>
                setFontConfig({
                    ...fontConfig,
                    certificate_subtitle_font_weight: parseInt(value),
                }),
        },
    ];

    const visibleFields = fontFields.filter((f) => f.visible);

    return {
        t,
        fontFamilies,
        isLoadingFonts,
        defaultFontFamilyId,
        visibleFields,
        hasAnyFields,
        getAvailableWeights,
    };
};
