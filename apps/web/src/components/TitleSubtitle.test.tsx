import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import TitleSubtitle from "./TitleSubtitle";

// Mock the Typography component
vi.mock("./typography/typography", () => ({
    Typography: ({ children, ...props }: React.ComponentProps<"div">) => (
        <div data-testid="typography" {...props}>
            {children}
        </div>
    ),
}));

import { vi } from "vitest";

describe("TitleSubtitle Component", () => {
    it("should render title", () => {
        render(<TitleSubtitle title="Test Title" />);
        expect(screen.getByText("Test Title")).toBeInTheDocument();
    });

    it("should render title and subtitle when subtitle is provided", () => {
        render(<TitleSubtitle title="Main Title" subtitle="Sub Title" />);

        expect(screen.getByText("Main Title")).toBeInTheDocument();
        expect(screen.getByText("Sub Title")).toBeInTheDocument();
    });

    it("should not render subtitle when not provided", () => {
        render(<TitleSubtitle title="Only Title" />);

        expect(screen.getByText("Only Title")).toBeInTheDocument();
        expect(screen.queryByText(/Sub/)).not.toBeInTheDocument();
    });

    it("should apply custom className to container", () => {
        const { container } = render(<TitleSubtitle title="Test" className="custom-class" />);

        const div = container.querySelector(".custom-class");
        expect(div).toBeInTheDocument();
    });

    it("should render with proper spacing class", () => {
        const { container } = render(<TitleSubtitle title="Test" />);

        const div = container.firstChild;
        expect(div).toHaveClass("lg:space-y-2");
    });

    it("should render Typography components for title and subtitle", () => {
        render(<TitleSubtitle title="Title" subtitle="Subtitle" />);

        const typographies = screen.getAllByTestId("typography");
        expect(typographies.length).toBeGreaterThanOrEqual(2);
    });

    it("should handle empty strings", () => {
        const { container } = render(<TitleSubtitle title="" />);
        expect(container.firstChild).toBeInTheDocument();
    });

    it("should handle long titles", () => {
        const longTitle = "A".repeat(100);
        render(<TitleSubtitle title={longTitle} />);

        expect(screen.getByText(longTitle)).toBeInTheDocument();
    });

    it("should handle special characters in title", () => {
        const specialTitle = "Title with <special> & characters!";
        render(<TitleSubtitle title={specialTitle} />);

        expect(screen.getByText(specialTitle)).toBeInTheDocument();
    });

    it("should render subtitle conditionally", () => {
        const { rerender } = render(<TitleSubtitle title="Title" subtitle="Subtitle" />);

        expect(screen.getByText("Subtitle")).toBeInTheDocument();

        rerender(<TitleSubtitle title="Title" subtitle={undefined} />);

        expect(screen.queryByText("Subtitle")).not.toBeInTheDocument();
    });

    it("should combine custom className with default classes", () => {
        const { container } = render(<TitleSubtitle title="Test" className="extra-padding" />);

        const div = container.firstChild as HTMLElement;
        expect(div.className).toContain("lg:space-y-2");
        expect(div.className).toContain("extra-padding");
    });
});
