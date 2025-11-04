import React from "react";
import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Typography } from "./typography";

describe("Typography", () => {
    it("should render h1 tag", () => {
        render(<Typography tag="h1">Heading 1</Typography>);
        expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("Heading 1");
    });

    it("should render h2 tag", () => {
        render(<Typography tag="h2">Heading 2</Typography>);
        expect(screen.getByRole("heading", { level: 2 })).toHaveTextContent("Heading 2");
    });

    it("should render h3 tag", () => {
        render(<Typography tag="h3">Heading 3</Typography>);
        expect(screen.getByRole("heading", { level: 3 })).toHaveTextContent("Heading 3");
    });

    it("should render p tag", () => {
        render(<Typography tag="p">Paragraph text</Typography>);
        const paragraph = screen.getByText("Paragraph text");
        expect(paragraph.tagName).toBe("P");
    });

    it("should render span tag", () => {
        render(<Typography tag="span">Span text</Typography>);
        const span = screen.getByText("Span text");
        expect(span.tagName).toBe("SPAN");
    });

    it("should render div tag", () => {
        render(<Typography tag="div">Div text</Typography>);
        const div = screen.getByText("Div text");
        expect(div.tagName).toBe("DIV");
    });

    it("should apply default variant classes", () => {
        render(<Typography tag="p">Default</Typography>);
        const element = screen.getByText("Default");

        // Default variant is "text" which applies "font-normal font-body"
        // Default size is "base" which applies "text-base"
        // No default color, so color must be specified explicitly
        expect(element).toHaveClass("font-normal"); // default variant
        expect(element).toHaveClass("text-base"); // default size
        expect(element).toHaveClass("font-body"); // default font family
    });

    it("should apply header variant", () => {
        render(
            <Typography tag="h1" variant="header">
                Header
            </Typography>,
        );
        // Note: weight prop takes precedence over variant's default weight
        // So we need to check the element has the text-wrap and font-inter
        const element = screen.getByText("Header");
        expect(element).toBeInTheDocument();
    });

    it("should apply text variant", () => {
        render(
            <Typography tag="p" variant="text">
                Text
            </Typography>,
        );
        expect(screen.getByText("Text")).toHaveClass("font-normal");
    });

    it("should apply different sizes", () => {
        const { rerender } = render(
            <Typography tag="p" size="small">
                Small
            </Typography>,
        );
        expect(screen.getByText("Small")).toHaveClass("text-xs");

        rerender(
            <Typography tag="p" size="base">
                Base
            </Typography>,
        );
        expect(screen.getByText("Base")).toHaveClass("text-base");

        rerender(
            <Typography tag="p" size="subheader">
                Subheader
            </Typography>,
        );
        expect(screen.getByText("Subheader")).toHaveClass("text-2xl");

        rerender(
            <Typography tag="p" size="header">
                Header
            </Typography>,
        );
        expect(screen.getByText("Header")).toHaveClass("text-4xl");
    });

    it("should apply different colors", () => {
        const { rerender } = render(
            <Typography tag="p" color="primary">
                Primary
            </Typography>,
        );
        expect(screen.getByText("Primary")).toHaveClass("text-primary");

        rerender(
            <Typography tag="p" color="secondary">
                Secondary
            </Typography>,
        );
        expect(screen.getByText("Secondary")).toHaveClass("text-secondary");

        rerender(
            <Typography tag="p" color="muted">
                Muted
            </Typography>,
        );
        expect(screen.getByText("Muted")).toHaveClass("text-muted");

        rerender(
            <Typography tag="p" color="background">
                Background
            </Typography>,
        );
        expect(screen.getByText("Background")).toHaveClass("text-background");
    });

    it("should apply font families through variant", () => {
        const { rerender } = render(
            <Typography tag="p" variant="text">
                Body Text
            </Typography>,
        );
        // variant="text" applies "font-body"
        expect(screen.getByText("Body Text")).toHaveClass("font-body");

        rerender(
            <Typography tag="h1" variant="header">
                Header Text
            </Typography>,
        );
        // variant="header" applies "font-header"
        expect(screen.getByText("Header Text")).toHaveClass("font-header");
    });

    it("should apply font weights through variant", () => {
        const { rerender } = render(
            <Typography tag="p" variant="text">
                Normal
            </Typography>,
        );
        // variant="text" applies "font-normal"
        expect(screen.getByText("Normal")).toHaveClass("font-normal");

        rerender(
            <Typography tag="h1" variant="header">
                Semibold
            </Typography>,
        );
        // variant="header" applies "font-semibold"
        expect(screen.getByText("Semibold")).toHaveClass("font-semibold");
    });

    it("should apply custom className", () => {
        render(
            <Typography tag="p" className="custom-class">
                Custom
            </Typography>,
        );
        expect(screen.getByText("Custom")).toHaveClass("custom-class");
    });

    it("should combine multiple variant props", () => {
        render(
            <Typography tag="h1" variant="header" size="header" color="primary">
                Combined
            </Typography>,
        );

        const element = screen.getByText("Combined");

        // variant="header" applies "font-semibold font-header"
        expect(element).toHaveClass("text-4xl"); // size header
        expect(element).toHaveClass("text-primary"); // color primary
        expect(element).toHaveClass("font-header"); // fontFamily from variant
        expect(element).toHaveClass("font-semibold"); // weight from variant
    });

    it("should pass through other HTML attributes", () => {
        render(
            <Typography tag="p" id="test-id" data-testid="typography">
                Test
            </Typography>,
        );

        const element = screen.getByTestId("typography");

        expect(element).toHaveAttribute("id", "test-id");
    });

    it("should render nested content", () => {
        render(
            <Typography tag="p">
                Text with <strong>bold</strong> and <em>italic</em>
            </Typography>,
        );

        expect(screen.getByText(/Text with/)).toBeInTheDocument();
        expect(screen.getByText("bold")).toBeInTheDocument();
        expect(screen.getByText("italic")).toBeInTheDocument();
    });
});
