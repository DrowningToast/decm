import React from "react";
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { ParticipatingNav } from "./ParticipatingNav";

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        className: "",
    }),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

describe("ParticipatingNav", () => {
    it("renders participating message", () => {
        render(<ParticipatingNav />);
        expect(screen.getByText("participant.events.participating")).toBeInTheDocument();
    });

    it("does not render a back button", () => {
        render(<ParticipatingNav />);
        expect(screen.queryByRole("button")).not.toBeInTheDocument();
    });
});
