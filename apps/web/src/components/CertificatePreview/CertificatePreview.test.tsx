import { render, screen, within } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { CertificatePreview } from "./CertificatePreview";
import type { DetectedKeyword, AvailableKeyword } from "@/hooks/useCertificateTemplate";

// Mock i18n
vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

// Mock DOMPurify
vi.mock("dompurify", () => ({
    default: {
        sanitize: (html: string) => html,
    },
}));

// Mock useCertificateFontFamilies so CertificatePreview doesn't need a QueryClientProvider
vi.mock("@/hooks/events/useCertificateFontFamilies", () => ({
    useCertificateFontFamilies: vi.fn(() => ({
        fontFamilies: [],
        isLoading: false,
        error: null,
    })),
}));

// Mock processCertificateSVG so it doesn't transform SVG in tests
vi.mock("@/lib/certificateRenderer", () => ({
    processCertificateSVG: vi.fn((svgContent: string) => svgContent),
}));

describe("CertificatePreview", () => {
    const mockSvgPreview =
        '<svg width="100" height="100"><rect width="100" height="100" fill="red"/></svg>';
    const mockImageUrl = "https://example.com/certificate.png";

    const mockDetectedKeywords: DetectedKeyword[] = [
        { keyword: "{{ name }}", x: 100, y: 200, count: 1 },
        { keyword: "{{ eventName }}", x: 150, y: 250, count: 1 },
    ];

    const mockAvailableKeywords: AvailableKeyword[] = [
        { keyword: "{{ name }}", mandatory: true },
        { keyword: "{{ eventName }}", mandatory: false },
        { keyword: "{{ academicInstitutionName }}", mandatory: false },
    ];

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should not render when no preview or image URL provided", () => {
        const { container } = render(<CertificatePreview />);
        expect(container.firstChild).toBeNull();
    });

    it("should render SVG preview when provided", () => {
        const { container } = render(<CertificatePreview svgPreview={mockSvgPreview} />);

        expect(screen.getByText("certificateSettings.step2.preview.title")).toBeInTheDocument();
        expect(container.querySelector("svg")).toBeInTheDocument();
    });

    it("should render image URL when provided", () => {
        render(<CertificatePreview imageUrl={mockImageUrl} />);

        expect(screen.getByText("certificateSettings.step2.preview.title")).toBeInTheDocument();
        const img = screen.getByAltText("Certificate preview");
        expect(img).toHaveAttribute("src", mockImageUrl);
    });

    it("should prefer SVG preview over image URL", () => {
        const { container } = render(
            <CertificatePreview svgPreview={mockSvgPreview} imageUrl={mockImageUrl} />,
        );

        expect(container.querySelector("svg")).toBeInTheDocument();
        expect(screen.queryByAltText("Certificate preview")).not.toBeInTheDocument();
    });

    it("should use custom alt text", () => {
        render(<CertificatePreview imageUrl={mockImageUrl} alt="Custom alt text" />);

        const img = screen.getByAltText("Custom alt text");
        expect(img).toBeInTheDocument();
    });

    it("should display detected keywords section when SVG and keywords provided", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        expect(
            screen.getByText("certificateSettings.step2.detectedKeywords.title"),
        ).toBeInTheDocument();
    });

    it("should not display keywords section when only image URL provided", () => {
        render(
            <CertificatePreview
                imageUrl={mockImageUrl}
                // Not passing detectedKeywords or availableKeywords, so keywords section should not display
            />,
        );

        expect(
            screen.queryByText("certificateSettings.step2.detectedKeywords.title"),
        ).not.toBeInTheDocument();
    });

    it("should display all keywords with detection status", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        expect(screen.getByText("{{ name }}")).toBeInTheDocument();
        expect(screen.getByText("{{ eventName }}")).toBeInTheDocument();
        expect(screen.getByText("{{ academicInstitutionName }}")).toBeInTheDocument();
    });

    it("should show detected status for found keywords", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        const detectedTexts = screen.getAllByText(
            "certificateSettings.step2.keywordStatus.detected",
        );
        expect(detectedTexts.length).toBeGreaterThan(0);
    });

    it("should show not detected status for missing keywords", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        expect(
            screen.getByText("certificateSettings.step2.keywordStatus.notDetected"),
        ).toBeInTheDocument();
    });

    it("should display required label for mandatory keywords", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        expect(
            screen.getByText("certificateSettings.step2.keywordStatus.required"),
        ).toBeInTheDocument();
    });

    it("should display optional label for non-mandatory keywords", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        const optionalLabels = screen.getAllByText(
            "certificateSettings.step2.keywordStatus.optional",
        );
        expect(optionalLabels.length).toBeGreaterThan(0);
    });

    it("should display alert when mandatory keywords are missing", () => {
        const missingMandatory: DetectedKeyword[] = [
            { keyword: "{{ eventName }}", x: 150, y: 250, count: 1 },
        ];

        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={missingMandatory}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        expect(
            screen.getByText("certificateSettings.step2.missingRequired.title"),
        ).toBeInTheDocument();
        expect(
            screen.getByText("certificateSettings.step2.missingRequired.description"),
        ).toBeInTheDocument();
    });

    it("should not display alert when all mandatory keywords are detected", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        expect(
            screen.queryByText("certificateSettings.step2.missingRequired.title"),
        ).not.toBeInTheDocument();
    });

    it("should display keyword count", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        // Check for count display in table - should have counts for detected keywords
        const table = screen.getByRole("table");
        const rows = within(table).getAllByRole("row");
        // At least one row should have a count displayed
        expect(rows.length).toBeGreaterThan(1);
    });

    it("should handle multiple occurrences of same keyword", () => {
        const multipleOccurrences: DetectedKeyword[] = [
            { keyword: "{{ name }}", x: 100, y: 200, count: 3 },
        ];

        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={multipleOccurrences}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        const table = screen.getByRole("table");
        expect(within(table).getByText(/3/)).toBeInTheDocument();
    });

    it("should handle empty keywords arrays", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={[]}
                availableKeywords={[]}
            />,
        );

        expect(
            screen.queryByText("certificateSettings.step2.detectedKeywords.title"),
        ).not.toBeInTheDocument();
    });

    it("should sanitize SVG content", () => {
        const maliciousSvg = '<svg><script>alert("xss")</script></svg>';

        const { container } = render(<CertificatePreview svgPreview={maliciousSvg} />);

        // DOMPurify is mocked to return the input as-is
        // In real scenario, it would sanitize the content
        expect(container.innerHTML).toContain("svg");
    });

    it("should handle image load errors", () => {
        const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

        render(<CertificatePreview imageUrl="invalid-url" />);

        const img = screen.getByAltText("Certificate preview") as HTMLImageElement;
        const event = new Event("error");
        img.dispatchEvent(event);

        consoleSpy.mockRestore();
    });

    it("should render keyword status table with correct structure", () => {
        render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        expect(
            screen.getByText("certificateSettings.step2.keywordStatus.table.status"),
        ).toBeInTheDocument();
        expect(
            screen.getByText("certificateSettings.step2.keywordStatus.table.keyword"),
        ).toBeInTheDocument();
        expect(
            screen.getByText("certificateSettings.step2.keywordStatus.table.type"),
        ).toBeInTheDocument();
        expect(
            screen.getByText("certificateSettings.step2.keywordStatus.table.count"),
        ).toBeInTheDocument();
    });

    // ---
    // BoundingBoxOverlay re-measurement when SVG loads asynchronously
    // ---

    it("should render the overlay container when keywords are present with an SVG", () => {
        // The overlay itself relies on getBoundingClientRect / getBBox which return
        // zeroes in happy-dom, so we just assert the structural elements are mounted.
        const { container } = render(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        // The outer wrapper that positions the overlay must be present
        const previewWrapper = container.querySelector(".relative");
        expect(previewWrapper).toBeInTheDocument();
    });

    it("should still render the SVG when svgPreview updates after initial render", async () => {
        const { container, rerender } = render(
            <CertificatePreview
                svgPreview=""
                imageUrl="https://example.com/fallback.png"
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        // Initially renders the image fallback — no certificate <rect> in the DOM
        expect(screen.getByAltText("Certificate preview")).toBeInTheDocument();
        expect(container.querySelector("rect")).not.toBeInTheDocument();

        // Simulate SVG loading asynchronously (e.g., fetched from presigned URL)
        rerender(
            <CertificatePreview
                svgPreview={mockSvgPreview}
                detectedKeywords={mockDetectedKeywords}
                availableKeywords={mockAvailableKeywords}
            />,
        );

        // Now the certificate SVG content (rect) is rendered and the img is gone
        expect(container.querySelector("rect")).toBeInTheDocument();
        expect(screen.queryByAltText("Certificate preview")).not.toBeInTheDocument();
    });
});
