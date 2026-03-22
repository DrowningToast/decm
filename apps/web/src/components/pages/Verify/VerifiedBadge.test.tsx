import { describe, it, expect } from "vitest";
import { screen } from "@testing-library/react";
import { renderWithProviders } from "../../../test/utils";
import { VerifiedBadge } from "./VerifiedBadge";

describe("VerifiedBadge", () => {
    it("shows spinner when verified is null", () => {
        renderWithProviders(<VerifiedBadge verified={null} />);
        expect(document.querySelector(".animate-spin")).toBeInTheDocument();
    });

    it("shows Verified badge when verified is true", () => {
        renderWithProviders(<VerifiedBadge verified={true} />);
        expect(screen.getByText("Verified")).toBeInTheDocument();
    });

    it("shows Invalid badge when verified is false", () => {
        renderWithProviders(<VerifiedBadge verified={false} />);
        expect(screen.getByText("Invalid")).toBeInTheDocument();
    });
});
