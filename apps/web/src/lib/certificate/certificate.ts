const SEARCH_KEYS = [
    "{{ name }}",
    "{{ eventName }}",
    "{{ academicInstitutionName }}",
    "{{ startDate }}",
    "{{ endDate }}",
];

// const outputContainer = document.getElementById("outputContainer");
// const statusMessage = document.getElementById("statusMessage");
// const svgTempContainer = document.getElementById("svg-temp-container");
// const generateButton = document.getElementById("generateButton");
// const outputCanvas = document.getElementById("outputCanvas");
// const downloadLink = document.getElementById("downloadLink");

// const CERT_WIDTH = 1920;
// const CERT_HEIGHT = 1080;

export function loadSvgTemplateFile(file: File): Document | null {
    const reader = new FileReader();
    let currentSvgDoc: Document | null = null;

    reader.onload = (e) => {
        const svgString = e?.target?.result;
        if (!svgString) return "";

        try {
            const parser = new DOMParser();
            const svgDoc = parser.parseFromString(svgString.toString(), "image/svg+xml");
            const svgElement = svgDoc.documentElement;

            if (svgElement.tagName.toLowerCase() !== "svg") {
                console.error("[ERROR] Root Element ไม่ใช่ <svg>");
                return;
            }

            console.log(`[PARSE] SVG Parsing สำเร็จ. Root element: <${svgElement.tagName}>`);
            currentSvgDoc = svgDoc;
        } catch (error) {
            console.error("[ERROR] SVG Parsing Error:", error);
        }
    };

    reader.readAsText(file);
    return currentSvgDoc;
}

type DetectTemplateKeysResponse = {
    key: keyof typeof SEARCH_KEYS;
    position: {
        x: number;
        y: number;
    };
}[];

export async function detectTemplateKeys(svgDoc: Document): Promise<DetectTemplateKeysResponse> {
    if (!svgDoc) throw new Error("SVG document is required");

    const originalSvgElement = svgDoc.documentElement;
    const clonedSvgElement = originalSvgElement.cloneNode(true) as HTMLElement;

    const CERT_WIDTH = originalSvgElement.getAttribute("width");
    const CERT_HEIGHT = originalSvgElement.getAttribute("height");

    if (!CERT_WIDTH || !CERT_HEIGHT) throw new Error("Certificate width or height is not set");

    clonedSvgElement.setAttribute("width", CERT_WIDTH);
    clonedSvgElement.setAttribute("height", CERT_HEIGHT);
    clonedSvgElement.setAttribute("viewBox", `0 0 ${CERT_WIDTH} ${CERT_HEIGHT}`);

    const svgTempContainer = document.createElement("div");
    svgTempContainer.innerHTML = "";
    svgTempContainer.appendChild(clonedSvgElement);

    const keys = SEARCH_KEYS.map((key) => {
        const escapedKey = CSS.escape(key);
        const x = clonedSvgElement.querySelector(`[data-key="${escapedKey}"]`);
        return {
            key,
            position: {
                x: x?.getBoundingClientRect().x ?? 0,
                y: x?.getBoundingClientRect().y ?? 0,
            },
        };
    });

    console.log(keys);

    return keys as DetectTemplateKeysResponse;
}
