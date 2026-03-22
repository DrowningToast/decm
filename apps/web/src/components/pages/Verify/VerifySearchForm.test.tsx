import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { VerifySearchForm } from "./VerifySearchForm";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
        i18n: { language: "en" },
    }),
}));

describe("VerifySearchForm", () => {
    const defaultProps = {
        inputCode: "",
        onInputCodeChange: vi.fn(),
        isFetching: false,
        onSubmit: vi.fn(),
    };

    it("renders the code input and verify button", () => {
        render(<VerifySearchForm {...defaultProps} />);
        expect(
            screen.getByPlaceholderText("certificateVerify.codePlaceholder"),
        ).toBeInTheDocument();
        expect(
            screen.getByRole("button", { name: /certificateVerify\.verifyButton/i }),
        ).toBeInTheDocument();
    });

    it("calls onInputCodeChange when input changes", () => {
        const onInputCodeChange = vi.fn();
        render(<VerifySearchForm {...defaultProps} onInputCodeChange={onInputCodeChange} />);
        fireEvent.change(screen.getByPlaceholderText("certificateVerify.codePlaceholder"), {
            target: { value: "abc123" },
        });
        expect(onInputCodeChange).toHaveBeenCalledWith("abc123");
    });

    it("calls onSubmit when form is submitted", () => {
        const onSubmit = vi.fn((e) => e.preventDefault());
        render(<VerifySearchForm {...defaultProps} inputCode="abc123" onSubmit={onSubmit} />);
        fireEvent.submit(screen.getByRole("button").closest("form")!);
        expect(onSubmit).toHaveBeenCalledTimes(1);
    });

    it("disables button when inputCode is empty", () => {
        render(<VerifySearchForm {...defaultProps} inputCode="" />);
        expect(
            screen.getByRole("button", { name: /certificateVerify\.verifyButton/i }),
        ).toBeDisabled();
    });

    it("disables button when fetching", () => {
        render(<VerifySearchForm {...defaultProps} inputCode="abc" isFetching={true} />);
        expect(
            screen.getByRole("button", { name: /certificateVerify\.verifyButton/i }),
        ).toBeDisabled();
    });

    it("enables button when inputCode is not empty and not fetching", () => {
        render(<VerifySearchForm {...defaultProps} inputCode="abc123" isFetching={false} />);
        expect(
            screen.getByRole("button", { name: /certificateVerify\.verifyButton/i }),
        ).not.toBeDisabled();
    });
});
