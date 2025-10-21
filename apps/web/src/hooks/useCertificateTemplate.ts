import { useState, useRef, useCallback } from "react";

export interface DetectedKeyword {
    keyword: string;
    x: number;
    y: number;
    count: number;
}

export interface AvailableKeyword {
    keyword: string;
    mandatory: boolean;
}

export interface UseCertificateTemplateProps {
    searchKeys?: string[];
    certWidth?: number;
    certHeight?: number;
}

export interface UseCertificateTemplateReturn {
    svgFile: File | null;
    svgPreview: string;
    detectedKeywords: DetectedKeyword[];
    availableKeywords: AvailableKeyword[];
    fileInputRef: React.RefObject<HTMLInputElement | null>;
    svgTempRef: React.RefObject<HTMLDivElement | null>;
    handleFileSelect: (event: React.ChangeEvent<HTMLInputElement>) => void;
    clearTemplate: () => void;
    hasMissingMandatory: boolean;
    missingMandatoryKeywords: AvailableKeyword[];
}

const DEFAULT_SEARCH_KEYS = ["{{ name }}", "{{ eventName }}"];
const DEFAULT_CERT_WIDTH = 1920;
const DEFAULT_CERT_HEIGHT = 1080;

const DEFAULT_AVAILABLE_KEYWORDS: AvailableKeyword[] = [
    { keyword: "{{ eventName }}", mandatory: true },
    { keyword: "{{ name }}", mandatory: true },
    { keyword: "{{ academicInstitutionName }}", mandatory: false },
];

export const useCertificateTemplate = ({
    searchKeys = DEFAULT_SEARCH_KEYS,
    certWidth = DEFAULT_CERT_WIDTH,
    certHeight = DEFAULT_CERT_HEIGHT,
}: UseCertificateTemplateProps = {}): UseCertificateTemplateReturn => {
    const [svgFile, setSvgFile] = useState<File | null>(null);
    const [svgPreview, setSvgPreview] = useState<string>("");
    const [detectedKeywords, setDetectedKeywords] = useState<DetectedKeyword[]>([]);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const svgTempRef = useRef<HTMLDivElement>(null);

    const generateCertificate = useCallback(
        async (_svgDoc: Document): Promise<void> => {
            const originalSvgElement = _svgDoc.documentElement;
            const clonedSvgElement = originalSvgElement.cloneNode(true) as SVGSVGElement;

            const tempContainer = svgTempRef.current;
            if (!tempContainer) return;

            tempContainer.innerHTML = "";
            tempContainer.appendChild(clonedSvgElement);

            clonedSvgElement.setAttribute("width", certWidth.toString());
            clonedSvgElement.setAttribute("height", certHeight.toString());
            clonedSvgElement.setAttribute("viewBox", `0 0 ${certWidth} ${certHeight}`);

            const newDetectedKeywords: DetectedKeyword[] = [];

            searchKeys.forEach((key) => {
                const escapedKey = CSS.escape(key);
                const placeholder = clonedSvgElement.querySelector(
                    `#${escapedKey}`,
                ) as SVGGraphicsElement | null;

                if (placeholder) {
                    const bbox = placeholder.getBBox();

                    const centerX = bbox.x + bbox.width / 2;
                    const y = bbox.y + bbox.height * 0.9;

                    newDetectedKeywords.push({
                        keyword: key,
                        x: centerX,
                        y: y,
                        count: 1,
                    });
                } else {
                    console.log(`[INFO] ไม่พบ Element ID: ${key}`);
                }
            });

            setDetectedKeywords(newDetectedKeywords);
        },
        [searchKeys, certWidth, certHeight],
    );

    const handleFileSelect = useCallback(
        async (event: React.ChangeEvent<HTMLInputElement>) => {
            const file = event.target.files?.[0];

            if (!file) return;
            setSvgFile(file);

            const reader = new FileReader();

            reader.onload = (e: ProgressEvent<FileReader>) => {
                const svgString = e.target?.result as string;
                setSvgPreview(svgString);

                try {
                    const parser = new DOMParser();
                    const svgDoc = parser.parseFromString(svgString, "image/svg+xml");
                    const svgElement = svgDoc.documentElement;

                    if (svgElement.tagName.toLowerCase() !== "svg") return;
                    generateCertificate(svgDoc);
                } catch (error) {
                    console.error("[ERROR] SVG Parsing Error:", error);
                }
            };

            reader.readAsText(file);
        },
        [generateCertificate],
    );

    const clearTemplate = useCallback(() => {
        setSvgFile(null);
        setSvgPreview("");
        setDetectedKeywords([]);
        if (fileInputRef.current) {
            fileInputRef.current.value = "";
        }
    }, []);

    // Check if mandatory keywords are detected
    const mandatoryKeywords = DEFAULT_AVAILABLE_KEYWORDS.filter((kw) => kw.mandatory);
    const detectedKeywordStrings = detectedKeywords.map((k) => k.keyword);
    const missingMandatoryKeywords = mandatoryKeywords.filter(
        (kw) => !detectedKeywordStrings.includes(kw.keyword),
    );
    const hasMissingMandatory = missingMandatoryKeywords.length > 0;

    return {
        svgFile,
        svgPreview,
        detectedKeywords,
        availableKeywords: DEFAULT_AVAILABLE_KEYWORDS,
        fileInputRef,
        svgTempRef,
        handleFileSelect,
        clearTemplate,
        hasMissingMandatory,
        missingMandatoryKeywords,
    };
};
