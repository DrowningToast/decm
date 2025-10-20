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

export function loadSvgTemplateFile(file: File): Promise<Document | null> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = (e) => {
            try {
                const svgString = e?.target?.result;
                if (!svgString) return "";
                const parser = new DOMParser();
                const svgDoc = parser.parseFromString(svgString.toString(), "image/svg+xml");

                const originalSvgElement = svgDoc.documentElement;
                const clonedSvgElement = originalSvgElement.cloneNode(true) as HTMLElement;

                const CERT_WIDTH = "1920";
                const CERT_HEIGHT = "1080";

                const svgTempContainer = document.createElement("div");
                svgTempContainer.innerHTML = "";
                svgTempContainer.appendChild(clonedSvgElement);

                clonedSvgElement.setAttribute("width", CERT_WIDTH);
                clonedSvgElement.setAttribute("height", CERT_HEIGHT);
                clonedSvgElement.setAttribute("viewBox", `0 0 ${CERT_WIDTH} ${CERT_HEIGHT}`);
                console.log("[DOM] แทรก SVG เข้าสู่ DOM ชั่วคราว (ขนาด 1920x1080)");

                SEARCH_KEYS.forEach((key) => {
                    const escapedKey = CSS.escape(key);
                    const placeholder = clonedSvgElement.querySelector(`#${escapedKey}`);

                    if (placeholder) {
                        console.log(
                            `[FOUND] พบตัวแปร ID: ${key} ใน Element <${placeholder.tagName}>`,
                        );

                        try {
                            const bbox = (placeholder as SVGGraphicsElement).getBBox();

                            const centerX = bbox.x + bbox.width / 2;
                            const y = bbox.y + bbox.height * 0.9;

                            console.log("centerX", centerX);
                            console.log("y", y);
                        } catch (error) {
                            console.error("[ERROR] การคำนวณ BBox ของ Element ล้มเหลว:", error);
                        }
                    }
                });

                resolve(svgDoc);
            } catch (error) {
                console.error("[ERROR] SVG Parsing Error:", error);
                reject(error);
            }
        };
        reader.readAsText(file);
    });
}

type DetectTemplateKeysResponse = {
    key: keyof typeof SEARCH_KEYS;
    position: {
        x: number;
        y: number;
    };
}[];

export async function detectTemplateKeys(svgDoc: Document): Promise<DetectTemplateKeysResponse> {
    console.log("svgDoc", svgDoc);

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
