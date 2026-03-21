import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { VerifiedBadge } from "./VerifiedBadge";

describe("VerifiedBadge", () => {
    it("shows spinner when verified is null", () => {
        render(<VerifiedBadge verified={null} />);
        expect(document.querySelector(".animate-spin")).toBeInTheDocument();
    });

    it("shows Verified badge when verified is true", () => {
        render(<VerifiedBadge verified={true} />);
        expect(screen.getByText("Verified")).toBeInTheDocument();
    });

    it("shows Invalid badge when verified is false", () => {
        render(<VerifiedBadge verified={false} />);
        expect(screen.getByText("Invalid")).toBeInTheDocument();
    });
});
