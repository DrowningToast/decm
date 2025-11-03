import { describe, it, expect, afterEach, vi } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useIssuerManagement } from "./useIssuerManagement";
import type { EntityProfile } from "@decm/api";
import type { Issuer } from "./useIssuerManagement";

const mockIssuerProfiles: EntityProfile[] = [
    {
        id: "issuer-1",
        authentication_credential_id: "cred-1",
        first_name: "Alice",
        last_name: "Johnson",
        email: "alice@example.com",
        academic_institution: "Stanford",
        bio: "Computer Science Professor",
    } as EntityProfile,
    {
        id: "issuer-2",
        authentication_credential_id: "cred-2",
        first_name: "Bob",
        last_name: "Smith",
        email: "bob@example.com",
        academic_institution: "MIT",
        bio: "Physics Department",
    } as EntityProfile,
    {
        id: "issuer-3",
        authentication_credential_id: "cred-3",
        first_name: "Carol",
        last_name: "Williams",
        email: "carol@example.com",
    } as EntityProfile,
];

describe("useIssuerManagement", () => {
    afterEach(() => {
        vi.clearAllMocks();
        vi.clearAllTimers();
    });

    it("should initialize with default values", () => {
        const { result } = renderHook(() => useIssuerManagement());

        expect(result.current.selectedIssuers).toEqual([]);
        expect(result.current.searchQuery).toBe("");
        expect(result.current.isSearching).toBe(false);
        expect(result.current.searchResults).toEqual([]);
        expect(result.current.isModalOpen).toBe(false);
    });

    it("should initialize with provided selected issuers", () => {
        const initialIssuers = [mockIssuerProfiles[0]];

        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: initialIssuers }),
        );

        expect(result.current.selectedIssuers).toEqual(initialIssuers);
    });

    it("should update search query", () => {
        const { result } = renderHook(() => useIssuerManagement());

        act(() => {
            result.current.handleSearchQueryChange("Alice");
        });

        expect(result.current.searchQuery).toBe("Alice");
    });

    it("should perform search with custom search function", async () => {
        const mockSearchFunction = vi.fn().mockResolvedValueOnce([
            {
                id: "issuer-1",
                name: "Alice Johnson",
                email: "alice@example.com",
                organization: "Stanford",
            },
        ]);

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: mockSearchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange("Alice");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        expect(result.current.isSearching).toBe(false);
        expect(result.current.searchResults).toHaveLength(1);
        expect(result.current.searchResults[0].name).toBe("Alice Johnson");
    });

    it("should not search with empty query", async () => {
        const mockSearchFunction = vi.fn();

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: mockSearchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange("");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        expect(mockSearchFunction).not.toHaveBeenCalled();
    });

    it("should use default search function with verified issuers", async () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: mockIssuerProfiles }),
        );

        act(() => {
            result.current.handleSearchQueryChange("Alice");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults).toHaveLength(1);
        });

        expect(result.current.searchResults[0].name).toBe("Alice Johnson");
    });

    it("should filter search results by email", async () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: mockIssuerProfiles }),
        );

        act(() => {
            result.current.handleSearchQueryChange("bob@example.com");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults).toHaveLength(1);
        });

        expect(result.current.searchResults[0].email).toBe("bob@example.com");
    });

    it("should filter search results by organization", async () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: mockIssuerProfiles }),
        );

        act(() => {
            result.current.handleSearchQueryChange("Stanford");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults).toHaveLength(1);
        });

        expect(result.current.searchResults[0].organization).toBe("Stanford");
    });

    it("should open modal after search", async () => {
        const mockSearchFunction = vi.fn().mockResolvedValueOnce([
            {
                id: "issuer-1",
                name: "Alice Johnson",
                email: "alice@example.com",
            },
        ]);

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: mockSearchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange("Alice");
        });

        expect(result.current.isModalOpen).toBe(false);

        await act(async () => {
            await result.current.handleSearch();
        });

        expect(result.current.isModalOpen).toBe(true);
    });

    it("should close modal and clear search results", () => {
        const { result } = renderHook(() => useIssuerManagement());

        act(() => {
            result.current.handleOpenModal();
        });

        expect(result.current.isModalOpen).toBe(true);

        act(() => {
            result.current.handleCloseModal();
        });

        expect(result.current.isModalOpen).toBe(false);
        expect(result.current.searchResults).toEqual([]);
    });

    it("should confirm selection and merge with existing issuers", async () => {
        const initialIssuers = [mockIssuerProfiles[0]];

        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: initialIssuers }),
        );

        const newIssuers = [
            {
                id: "issuer-2",
                name: "Bob Smith",
                email: "bob@example.com",
                organization: "MIT",
            },
        ];

        act(() => {
            result.current.handleConfirmSelection(newIssuers as Issuer[]);
        });

        expect(result.current.selectedIssuers).toHaveLength(2);
        expect(result.current.isModalOpen).toBe(false);
    });

    it("should avoid duplicate selections", () => {
        const initialIssuers = [mockIssuerProfiles[0]];

        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: initialIssuers }),
        );

        // Create an Issuer object with the same authentication_credential_id
        const duplicateIssuer = {
            id: "cred-1", // This should match authentication_credential_id
            name: "Alice Johnson",
            email: "alice@example.com",
            organization: "Stanford",
        };

        act(() => {
            result.current.handleConfirmSelection([duplicateIssuer]);
        });

        // Should still be 1 since we're trying to add the same issuer
        expect(result.current.selectedIssuers).toHaveLength(1);
    });

    it("should remove issuer by id", () => {
        const initialIssuers = mockIssuerProfiles.slice(0, 2);

        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: initialIssuers }),
        );

        expect(result.current.selectedIssuers).toHaveLength(2);

        act(() => {
            result.current.handleRemoveIssuer("cred-1");
        });

        expect(result.current.selectedIssuers).toHaveLength(1);
        expect(result.current.selectedIssuers[0].authentication_credential_id).toBe("cred-2");
    });

    it("should clear all selections", () => {
        const initialIssuers = mockIssuerProfiles.slice(0, 3);

        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: initialIssuers }),
        );

        expect(result.current.selectedIssuers).toHaveLength(3);

        act(() => {
            result.current.handleClearSelection();
        });

        expect(result.current.selectedIssuers).toHaveLength(0);
    });

    it("should get selected issuer ids", () => {
        const initialIssuers = mockIssuerProfiles.slice(0, 2);

        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: initialIssuers }),
        );

        const ids = result.current.getSelectedIssuerIds();

        expect(ids instanceof Set).toBe(true);
        expect(ids.has("cred-1")).toBe(true);
        expect(ids.has("cred-2")).toBe(true);
        expect(ids.has("cred-3")).toBe(false);
    });

    it("should handle search error gracefully", async () => {
        const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});
        const mockSearchFunction = vi.fn().mockRejectedValueOnce(new Error("API error"));

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: mockSearchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange("test");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        expect(result.current.isSearching).toBe(false);
        expect(consoleSpy).toHaveBeenCalled();

        consoleSpy.mockRestore();
    });

    it("should provide case-insensitive search", async () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: mockIssuerProfiles }),
        );

        act(() => {
            result.current.handleSearchQueryChange("ALICE");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults).toHaveLength(1);
        });

        expect(result.current.searchResults[0].name).toBe("Alice Johnson");
    });

    it("should handle issuer without organization", async () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: mockIssuerProfiles }),
        );

        act(() => {
            result.current.handleSearchQueryChange("Carol");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults).toHaveLength(1);
        });

        expect(result.current.searchResults[0].name).toBe("Carol Williams");
        expect(result.current.searchResults[0].organization).toBeUndefined();
    });

    it("should trim search query", async () => {
        const mockSearchFunction = vi.fn().mockResolvedValueOnce([]);

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: mockSearchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange("  Alice  ");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        expect(mockSearchFunction).toHaveBeenCalledWith("Alice");
    });

    it("should maintain stable function references", () => {
        const { result, rerender } = renderHook(() => useIssuerManagement());

        const fn1 = result.current.handleSearch;
        const fn2 = result.current.handleRemoveIssuer;
        const fn3 = result.current.getSelectedIssuerIds;

        rerender();

        const fn1New = result.current.handleSearch;
        const fn2New = result.current.handleRemoveIssuer;
        const fn3New = result.current.getSelectedIssuerIds;

        // Functions should be stable across rerenders
        expect(fn1).toBe(fn1New);
        expect(fn2).toBe(fn2New);
        expect(fn3).toBe(fn3New);
    });
});
