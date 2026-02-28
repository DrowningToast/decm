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

export interface DetectedKeyword {
    keyword: string;
    x: number;
    y: number;
    count: number;
}

export interface FontFamily {
    id: number;
    css_font_name: string;
    font_family_name: string;
}

/**
 * Replicates backend logic for hiding template placeholders in SVG
 */
export const hideTemplatePlaceholders = (svgContent: string): string => {
    let result = svgContent;

    const placeholderMap: Record<string, string> = {
        name: "{{ name }}",
        eventName: "{{ eventName }}",
        academicInstitutionName: "{{ academicInstitutionName }}",
        certificateTitle: "{{ certificateTitle }}",
        certificateSubtitle: "{{ certificateSubtitle }}",
    };

    // First pass: Hide elements by ID matching
    for (const [name, placeholder] of Object.entries(placeholderMap)) {
        // Pattern 1: id="{{ name }}"
        const pattern1 = new RegExp(
            `(<text[^>]*id="${placeholder.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}"[^>]*>)`,
            "g",
        );
        result = result.replace(pattern1, (match) => {
            if (!match.includes("visibility")) {
                return match.replace(/>$/, ' visibility="hidden">');
            }
            return match;
        });

        // Pattern 2: id="name"
        const pattern2 = new RegExp(`(<text[^>]*id="${name}"[^>]*>)`, "g");
        result = result.replace(pattern2, (match) => {
            if (!match.includes("visibility")) {
                return match.replace(/>$/, ' visibility="hidden">');
            }
            return match;
        });
    }

    // Second pass: Hide elements by content matching
    for (const placeholder of Object.values(placeholderMap)) {
        const escapedPlaceholder = placeholder.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
        const pattern = new RegExp(`(<text[^>]*>)([^<]*${escapedPlaceholder}[^<]*)(</text>)`, "g");

        result = result.replace(pattern, (match, openingTag, content, closingTag) => {
            if (openingTag.includes('visibility="hidden"')) {
                return match;
            }
            const modifiedOpeningTag = openingTag.replace(/>$/, ' visibility="hidden">');
            return modifiedOpeningTag + content + closingTag;
        });
    }

    return result;
};

/**
 * Replicates backend logic for creating SVG text elements
 */
export const createTextElement = (
    text: string,
    x: number,
    y: number,
    fontFamily: string,
    fontWeight: string | number,
    fontSize: number = 16,
): string => {
    const escapedText = text
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&apos;");

    return `<text x="${x.toFixed(2)}" y="${y.toFixed(2)}" font-family="${fontFamily}" font-size="${fontSize}" font-weight="${fontWeight}" fill="#090909" text-anchor="middle">${escapedText}</text>`;
};

/**
 * Processes SVG template by hiding placeholders and adding actual text overlays
 * matching the backend rendering logic.
 */
export const processCertificateSVG = (
    svgContent: string,
    fontConfig: CertificateFontConfig,
    detectedKeywords: DetectedKeyword[],
    fontFamilies: FontFamily[],
    previewData: {
        name?: string;
        eventName?: string;
        academicInstitution?: string;
        certificateTitle?: string;
        certificateSubtitle?: string;
    },
): string => {
    // 1. Hide placeholders
    const result = hideTemplatePlaceholders(svgContent);

    // 2. Build text overlays
    let textOverlays = "";

    const getFontName = (id?: number) => {
        if (!id) return "Prompt";
        return fontFamilies.find((f) => f.id === id)?.css_font_name || "Prompt";
    };

    const keywordMap = new Map(detectedKeywords.map((kw) => [kw.keyword, kw]));

    const fields = [
        {
            data: previewData.name || "Sample Name",
            keyword: "{{ name }}",
            fontFamilyId: fontConfig.name_font_family_id,
            fontWeight: fontConfig.name_font_weight || 700,
        },
        {
            data: previewData.eventName || "Sample Event",
            keyword: "{{ eventName }}",
            fontFamilyId: fontConfig.event_name_font_family_id,
            fontWeight: fontConfig.event_name_font_weight || 700,
        },
        {
            data: previewData.academicInstitution || "Sample Institution",
            keyword: "{{ academicInstitutionName }}",
            fontFamilyId: fontConfig.academic_institution_font_family_id,
            fontWeight: fontConfig.academic_institution_font_weight || 700,
        },
        {
            data: previewData.certificateTitle || "Sample Title",
            keyword: "{{ certificateTitle }}",
            fontFamilyId: fontConfig.certificate_title_font_family_id,
            fontWeight: fontConfig.certificate_title_font_weight || 700,
        },
        {
            data: previewData.certificateSubtitle || "Sample Subtitle",
            keyword: "{{ certificateSubtitle }}",
            fontFamilyId: fontConfig.certificate_subtitle_font_family_id,
            fontWeight: fontConfig.certificate_subtitle_font_weight || 700,
        },
    ];

    for (const field of fields) {
        const kw = keywordMap.get(field.keyword);
        if (kw && field.data) {
            textOverlays += createTextElement(
                field.data,
                kw.x,
                kw.y,
                getFontName(field.fontFamilyId),
                field.fontWeight,
            );
        }
    }

    // 3. Insert before closing </svg>
    const closingSvgIndex = result.lastIndexOf("</svg>");
    if (closingSvgIndex === -1) {
        return result + textOverlays;
    }

    return result.slice(0, closingSvgIndex) + textOverlays + result.slice(closingSvgIndex);
};
