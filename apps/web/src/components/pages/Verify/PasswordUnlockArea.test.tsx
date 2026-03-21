import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { PasswordUnlockArea } from "./PasswordUnlockArea";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
        i18n: { language: "en" },
    }),
}));

describe("PasswordUnlockArea", () => {
    const defaultProps = {
        password: "",
        onPasswordChange: vi.fn(),
        passwordError: "",
        onUnlock: vi.fn(),
        isUnlocking: false,
    };

    it("renders password input and unlock button", () => {
        render(<PasswordUnlockArea {...defaultProps} />);
        expect(
            screen.getByPlaceholderText("certificateVerify.passwordPlaceholder"),
        ).toBeInTheDocument();
        expect(
            screen.getByRole("button", { name: "certificateVerify.unlockButton" }),
        ).toBeInTheDocument();
    });

    it("calls onPasswordChange when input changes", () => {
        const onPasswordChange = vi.fn();
        render(<PasswordUnlockArea {...defaultProps} onPasswordChange={onPasswordChange} />);
        fireEvent.change(screen.getByPlaceholderText("certificateVerify.passwordPlaceholder"), {
            target: { value: "secret" },
        });
        expect(onPasswordChange).toHaveBeenCalledWith("secret");
    });

    it("calls onUnlock when button is clicked", () => {
        const onUnlock = vi.fn();
        render(<PasswordUnlockArea {...defaultProps} onUnlock={onUnlock} />);
        fireEvent.click(screen.getByRole("button", { name: "certificateVerify.unlockButton" }));
        expect(onUnlock).toHaveBeenCalledTimes(1);
    });

    it("calls onUnlock when Enter is pressed in input", () => {
        const onUnlock = vi.fn();
        render(<PasswordUnlockArea {...defaultProps} onUnlock={onUnlock} />);
        fireEvent.keyDown(screen.getByPlaceholderText("certificateVerify.passwordPlaceholder"), {
            key: "Enter",
        });
        expect(onUnlock).toHaveBeenCalledTimes(1);
    });

    it("displays password error message with role=alert", () => {
        render(<PasswordUnlockArea {...defaultProps} passwordError="Incorrect password" />);
        const alert = screen.getByRole("alert");
        expect(alert).toHaveTextContent("Incorrect password");
    });

    it("disables button and shows loading when isUnlocking", () => {
        render(<PasswordUnlockArea {...defaultProps} isUnlocking={true} />);
        expect(
            screen.getByRole("button", { name: "certificateVerify.unlockButton" }),
        ).toBeDisabled();
    });
});
