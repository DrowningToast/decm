import type { CoreApiInternalHandlerEventconfigEventCertificateConfigResponse } from "@decm/api";
import type { DetectedKeyword } from "@/hooks/useCertificateTemplate";

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

export interface CertificateFontSettingsProps {
    eventCertificateConfig?: CoreApiInternalHandlerEventconfigEventCertificateConfigResponse;
    detectedKeywords?: DetectedKeyword[];
    onChange?: (fontConfig: CertificateFontConfig) => void;
}
