import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { LanguageSwitcher } from "./LanguageSwitcher";

// Mock i18n
const mockI18n = {
    language: "en",
    changeLanguage: vi.fn(),
};

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        i18n: mockI18n,
        t: (key: string) => key,
    }),
}));

// Mock components
vi.mock("@/components/ui/button", () => ({
    Button: ({ children, ...props }: React.ComponentProps<"button">) => (
        <button {...props}>{children}</button>
    ),
}));

vi.mock("@/components/ui/dropdown-menu", () => ({
    DropdownMenu: ({ children }: React.PropsWithChildren) => <div>{children}</div>,
    DropdownMenuTrigger: ({ children }: React.PropsWithChildren) => <div>{children}</div>,
    DropdownMenuContent: ({ children }: React.PropsWithChildren) => <div>{children}</div>,
    DropdownMenuItem: ({
        children,
        onClick,
    }: React.PropsWithChildren<{ onClick?: () => void }>) => (
        <div onClick={onClick}>{children}</div>
    ),
}));

// Mock language data
vi.mock("@/lib/i18n", () => ({
    languages: {
        en: { label: "English", flag: "🇬🇧" },
        es: { label: "Español", flag: "🇪🇸" },
        th: { label: "ไทย", flag: "🇹🇭" },
    },
}));

describe("LanguageSwitcher Component", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockI18n.language = "en";
    });

    it("should render without crashing", () => {
        render(<LanguageSwitcher />);
        expect(screen.getByText("🇬🇧")).toBeInTheDocument();
    });

    it("should display current language flag", () => {
        mockI18n.language = "en";
        render(<LanguageSwitcher />);

        expect(screen.getByText("🇬🇧")).toBeInTheDocument();
    });

    it("should display current language label on larger screens", () => {
        mockI18n.language = "es";
        render(<LanguageSwitcher />);

        expect(screen.getByText("Español")).toBeInTheDocument();
    });

    it("should fallback to English when language is not recognized", () => {
        mockI18n.language = "fr";
        render(<LanguageSwitcher />);

        // Should fallback to English
        expect(screen.getByText("🇬🇧")).toBeInTheDocument();
    });

    it("should change language when menu item is clicked", async () => {
        render(<LanguageSwitcher />);

        // Find and click Spanish option
        const spanishOption = screen.getByText("Español");
        fireEvent.click(spanishOption);

        await waitFor(() => {
            expect(mockI18n.changeLanguage).toHaveBeenCalledWith("es");
        });
    });

    it("should render all available languages", () => {
        render(<LanguageSwitcher />);

        expect(screen.getByText("English")).toBeInTheDocument();
        expect(screen.getByText("Español")).toBeInTheDocument();
        expect(screen.getByText("ไทย")).toBeInTheDocument();
    });

    it("should render flags for all languages", () => {
        render(<LanguageSwitcher />);

        expect(screen.getAllByText("🇬🇧").length).toBeGreaterThan(0);
        expect(screen.getByText("🇪🇸")).toBeInTheDocument();
        expect(screen.getByText("🇹🇭")).toBeInTheDocument();
    });

    it("should show checkmark for current language", () => {
        mockI18n.language = "th";
        render(<LanguageSwitcher />);

        // Should show checkmark for Thai
        const checkmarks = screen.getAllByText("✓");
        expect(checkmarks.length).toBeGreaterThan(0);
    });

    it("should change language to English", async () => {
        mockI18n.language = "es";
        render(<LanguageSwitcher />);

        const englishOption = screen.getByText("English");
        fireEvent.click(englishOption);

        await waitFor(() => {
            expect(mockI18n.changeLanguage).toHaveBeenCalledWith("en");
        });
    });

    it("should change language to Thai", async () => {
        render(<LanguageSwitcher />);

        const thaiOption = screen.getByText("ไทย");
        fireEvent.click(thaiOption);

        await waitFor(() => {
            expect(mockI18n.changeLanguage).toHaveBeenCalledWith("th");
        });
    });

    it("should persist language selection", async () => {
        const { rerender } = render(<LanguageSwitcher />);

        const spanishOption = screen.getByText("Español");
        fireEvent.click(spanishOption);

        await waitFor(() => {
            expect(mockI18n.changeLanguage).toHaveBeenCalledWith("es");
        });

        mockI18n.language = "es";
        rerender(<LanguageSwitcher />);

        expect(screen.getByText("🇪🇸")).toBeInTheDocument();
    });

    it("should handle rapid language changes", async () => {
        render(<LanguageSwitcher />);

        const englishOption = screen.getByText("English");
        const spanishOption = screen.getByText("Español");
        const thaiOption = screen.getByText("ไทย");

        fireEvent.click(spanishOption);
        fireEvent.click(thaiOption);
        fireEvent.click(englishOption);

        await waitFor(() => {
            expect(mockI18n.changeLanguage).toHaveBeenCalledTimes(3);
        });
    });
});
