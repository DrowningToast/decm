import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { HowItWorks } from "./HowItWorks";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
        i18n: { language: "en" },
    }),
}));

describe("HowItWorks", () => {
    it("renders the section heading", () => {
        render(<HowItWorks />);
        expect(screen.getByText("certificateVerify.how.heading")).toBeInTheDocument();
    });

    it("renders all three step titles", () => {
        render(<HowItWorks />);
        expect(screen.getAllByText("certificateVerify.how.step1.title").length).toBeGreaterThan(0);
        expect(screen.getAllByText("certificateVerify.how.step2.title").length).toBeGreaterThan(0);
        expect(screen.getAllByText("certificateVerify.how.step3.title").length).toBeGreaterThan(0);
    });
});
