import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { OnChainTable } from "./OnChainTable";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
        i18n: { language: "en" },
    }),
}));

const defaultProps = {
    contractAddress: "0xContractAddress",
    tokenId: "42",
    open: false,
    onToggle: vi.fn(),
    onChainData: null,
    isLoading: false,
    isError: false,
    expectedHash: "0xhash",
    expectedOwner: "0xowner",
};

describe("OnChainTable", () => {
    it("renders collapsed by default", () => {
        render(<OnChainTable {...defaultProps} />);
        expect(screen.queryByText("0xContractAddress")).not.toBeInTheDocument();
    });

    it("calls onToggle when header is clicked", () => {
        const onToggle = vi.fn();
        render(<OnChainTable {...defaultProps} onToggle={onToggle} />);
        fireEvent.click(screen.getByRole("button"));
        expect(onToggle).toHaveBeenCalledTimes(1);
    });

    it("shows loading indicator when open and loading", () => {
        render(<OnChainTable {...defaultProps} open={true} isLoading={true} />);
        expect(document.querySelector(".animate-spin")).toBeInTheDocument();
    });

    it("shows error message when open and error", () => {
        render(<OnChainTable {...defaultProps} open={true} isError={true} />);
        expect(screen.getByText("certificateVerify.onChain.error")).toBeInTheDocument();
    });

    it("shows contract and token data when open with data", () => {
        const onChainData = { owner: "0xowner", tokenData: '{"hash":"0xhash"}' };
        render(<OnChainTable {...defaultProps} open={true} onChainData={onChainData} />);
        expect(screen.getByText("#42")).toBeInTheDocument();
    });

    it("marks owner as verified when addresses match", () => {
        const onChainData = { owner: "0xowner", tokenData: "" };
        render(
            <OnChainTable
                {...defaultProps}
                open={true}
                onChainData={onChainData}
                expectedOwner="0xowner"
            />,
        );
        expect(screen.getByText("Verified")).toBeInTheDocument();
    });

    it("marks owner as invalid when addresses differ", () => {
        const onChainData = { owner: "0xdifferent", tokenData: "" };
        render(
            <OnChainTable
                {...defaultProps}
                open={true}
                onChainData={onChainData}
                expectedOwner="0xowner"
            />,
        );
        expect(screen.getByText("Invalid")).toBeInTheDocument();
    });
});
