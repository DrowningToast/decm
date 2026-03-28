import { describe, it, expect } from "vitest";
import { screen } from "@testing-library/react";
import { renderWithProviders } from "../../../test/utils";
import { Row } from "./Row";

// <tr> must live inside <table><tbody> to be valid DOM
function renderRow(props: Parameters<typeof Row>[0]) {
    return renderWithProviders(
        <table>
            <tbody>
                <Row {...props} />
            </tbody>
        </table>,
    );
}

describe("Row", () => {
    describe("label and value rendering", () => {
        it("renders the label", () => {
            renderRow({ label: "Full Name", value: "Alice" });
            expect(screen.getByText("Full Name")).toBeInTheDocument();
        });

        it("renders a string value", () => {
            renderRow({ label: "Full Name", value: "Alice" });
            expect(screen.getByText("Alice")).toBeInTheDocument();
        });

        it("renders a ReactNode value", () => {
            renderRow({ label: "Status", value: <span data-testid="custom">custom node</span> });
            expect(screen.getByTestId("custom")).toBeInTheDocument();
        });
    });

    describe("verified badge behaviour (existing)", () => {
        it("shows Verified badge when verified is true", () => {
            renderRow({ label: "Name", value: "Alice", verified: true });
            expect(screen.getByText("Verified")).toBeInTheDocument();
        });

        it("shows Invalid badge when verified is false and no warningTooltip", () => {
            renderRow({ label: "Name", value: "Alice", verified: false });
            expect(screen.getByText("Invalid")).toBeInTheDocument();
        });

        it("shows spinner when verified is null", () => {
            renderRow({ label: "Name", value: "Alice", verified: null });
            expect(document.querySelector(".animate-spin")).toBeInTheDocument();
        });

        it("shows no badge when verified is undefined", () => {
            renderRow({ label: "Name", value: "Alice" });
            expect(screen.queryByText("Verified")).not.toBeInTheDocument();
            expect(screen.queryByText("Invalid")).not.toBeInTheDocument();
            expect(document.querySelector(".animate-spin")).not.toBeInTheDocument();
        });
    });

    describe("warningTooltip behaviour", () => {
        it("shows warning icon instead of Invalid badge when verified is false and warningTooltip is set", () => {
            renderRow({
                label: "Name",
                value: "Alice",
                verified: false,
                warningTooltip: "Names do not match.",
            });
            expect(screen.queryByText("Invalid")).not.toBeInTheDocument();
            expect(document.querySelector("svg")).toBeInTheDocument();
        });

        it("shows tooltip content on hover when warningTooltip is provided", async () => {
            const { userEvent } = await import("@testing-library/user-event");
            const user = userEvent.setup();
            renderRow({
                label: "Name",
                value: "Alice",
                verified: false,
                warningTooltip: "Names do not match.",
            });
            // TooltipTrigger renders as a button (no asChild) so pointer events work correctly
            const trigger = screen.getByRole("button");
            await user.hover(trigger);
            // Radix renders text in both the visible div and an aria role="tooltip" span
            const tooltip = await screen.findByRole("tooltip");
            expect(tooltip).toHaveTextContent("Names do not match.");
        });

        it("ignores warningTooltip when verified is true and shows Verified badge", () => {
            renderRow({
                label: "Name",
                value: "Alice",
                verified: true,
                warningTooltip: "Should be ignored",
            });
            expect(screen.getByText("Verified")).toBeInTheDocument();
            expect(screen.queryByText("Invalid")).not.toBeInTheDocument();
        });

        it("ignores warningTooltip when verified is null and shows spinner", () => {
            renderRow({
                label: "Name",
                value: "Alice",
                verified: null,
                warningTooltip: "Should be ignored",
            });
            expect(document.querySelector(".animate-spin")).toBeInTheDocument();
        });
    });
});
