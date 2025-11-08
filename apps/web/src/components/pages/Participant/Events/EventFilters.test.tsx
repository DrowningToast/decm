import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { EventFilters } from "./EventFilters";

// Mock react-i18next
vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => {
            const translations: Record<string, string> = {
                "participant.events.filters.allEvents": "All Events",
                "participant.events.filters.myEvents": "My Events",
            };
            return translations[key] || key;
        },
    }),
}));

describe("EventFilters", () => {
    it("renders all filter tabs", () => {
        render(<EventFilters />);

        expect(screen.getByText("All Events")).toBeInTheDocument();
        expect(screen.getByText("My Events")).toBeInTheDocument();
    });

    it("defaults to 'all' tab", () => {
        render(<EventFilters />);

        const allEventsTab = screen.getByRole("tab", { name: /all events/i });
        expect(allEventsTab).toHaveAttribute("data-state", "active");
    });

    it("calls onFilterChange when tab is clicked", async () => {
        const user = userEvent.setup();
        const onFilterChange = vi.fn();

        render(<EventFilters onFilterChange={onFilterChange} />);

        const myEventsTab = screen.getByRole("tab", { name: /my events/i });
        await user.click(myEventsTab);

        expect(onFilterChange).toHaveBeenCalledWith("my-events");
    });

    it("calls onFilterChange with 'all' when all events tab is clicked", async () => {
        const user = userEvent.setup();
        const onFilterChange = vi.fn();

        render(<EventFilters onFilterChange={onFilterChange} />);

        const myEventsTab = screen.getByRole("tab", { name: /my events/i });
        await user.click(myEventsTab);

        const allEventsTab = screen.getByRole("tab", { name: /all events/i });
        await user.click(allEventsTab);

        expect(onFilterChange).toHaveBeenCalledWith("all");
    });

    it("does not crash when onFilterChange is not provided", async () => {
        const user = userEvent.setup();

        render(<EventFilters />);

        const myEventsTab = screen.getByRole("tab", { name: /my events/i });
        await user.click(myEventsTab);

        // No error should be thrown
        expect(myEventsTab).toHaveAttribute("data-state", "active");
    });

    it("switches tab state when clicked", async () => {
        const user = userEvent.setup();

        render(<EventFilters />);

        const allEventsTab = screen.getByRole("tab", { name: /all events/i });
        const myEventsTab = screen.getByRole("tab", { name: /my events/i });

        // Initially, all events should be active
        expect(allEventsTab).toHaveAttribute("data-state", "active");
        expect(myEventsTab).toHaveAttribute("data-state", "inactive");

        // Click my events
        await user.click(myEventsTab);

        // Now my events should be active
        expect(allEventsTab).toHaveAttribute("data-state", "inactive");
        expect(myEventsTab).toHaveAttribute("data-state", "active");
    });
});
