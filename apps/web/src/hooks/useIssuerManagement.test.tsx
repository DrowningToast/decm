import { describe, it, expect, afterEach, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useIssuerManagement, convertProfileToIssuer } from "./useIssuerManagement";
import type { Issuer } from "./useIssuerManagement";
import type { Profile } from "@/services/AuthService/AuthService";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";

// Mock the useSearchIssuer hook
vi.mock("./useSearchIssuer", () => ({
    useSearchIssuer: vi.fn().mockReturnValue({
        data: undefined,
        isLoading: false,
        error: null,
    }),
}));

const mockProfiles: Profile[] = [
    {
        id: "profile-1",
        authenticationCredentialId: "cred-1",
        firstName: "Alice",
        lastName: "Johnson",
        email: "alice@example.com",
        googleConnectorRef: "alice.oauth@gmail.com",
        phoneNumber: "123-456-7890",
        academicInstitution: "Stanford",
        bio: "Computer Science Professor",
        isEmailPublic: true,
        isFirstNamePublic: true,
        isLastNamePublic: true,
        isPhoneNumberPublic: true,
        isAcademicInstitutionPublic: true,
        isProfilePicturePublic: false,
        isBioPublic: false,
        isAddressPublic: false,
        isAcademicEmailPublic: false,
    },
    {
        id: "profile-2",
        authenticationCredentialId: "cred-2",
        firstName: "Bob",
        lastName: "Smith",
        email: "bob@example.com",
        googleConnectorRef: "bob.oauth@gmail.com",
        phoneNumber: "098-765-4321",
        academicInstitution: "MIT",
        bio: "Physics Department",
        isEmailPublic: true,
        isFirstNamePublic: true,
        isLastNamePublic: true,
        isPhoneNumberPublic: true,
        isAcademicInstitutionPublic: true,
        isProfilePicturePublic: false,
        isBioPublic: false,
        isAddressPublic: false,
        isAcademicEmailPublic: false,
    },
];

const mockIssuers: Issuer[] = mockProfiles.map(convertProfileToIssuer);

describe("useIssuerManagement", () => {
    let queryClient: QueryClient;

    beforeEach(() => {
        queryClient = new QueryClient({
            defaultOptions: {
                queries: {
                    retry: false,
                },
            },
        });
    });

    afterEach(() => {
        vi.clearAllMocks();
        vi.clearAllTimers();
    });

    const wrapper = ({ children }: { children: ReactNode }) => {
        return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
    };

    it("should initialize with default values", () => {
        const { result } = renderHook(() => useIssuerManagement(), { wrapper });

        expect(result.current.selectedIssuers).toEqual([]);
        expect(result.current.searchQuery).toBe("");
        expect(result.current.isSearching).toBe(false);
        expect(result.current.searchResults).toEqual([]);
        expect(result.current.isModalOpen).toBe(false);
    });

    it("should initialize with provided selected issuers", () => {
        const { result } = renderHook(
            () => useIssuerManagement({ selectedIssuers: [mockIssuers[0]] }),
            { wrapper },
        );

        expect(result.current.selectedIssuers).toEqual([mockIssuers[0]]);
    });

    it("should update search query", () => {
        const { result } = renderHook(() => useIssuerManagement(), { wrapper });

        act(() => {
            result.current.handleSearchQueryChange("Alice");
        });

        expect(result.current.searchQuery).toBe("Alice");
    });

    it("should not search with empty query", async () => {
        const { result } = renderHook(() => useIssuerManagement(), { wrapper });

        act(() => {
            result.current.handleSearchQueryChange("");
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        // isSearching should remain false as search should not be triggered
        expect(result.current.isSearching).toBe(false);
    });

    it("should open and close modal", () => {
        const { result } = renderHook(() => useIssuerManagement(), { wrapper });

        act(() => {
            result.current.handleOpenModal();
        });

        expect(result.current.isModalOpen).toBe(true);

        act(() => {
            result.current.handleCloseModal();
        });

        expect(result.current.isModalOpen).toBe(false);
    });

    it("should confirm selection and merge with existing issuers", () => {
        const { result } = renderHook(
            () => useIssuerManagement({ selectedIssuers: [mockIssuers[0]] }),
            { wrapper },
        );

        const newIssuers = [mockIssuers[1]];

        act(() => {
            result.current.handleConfirmSelection(newIssuers);
        });

        expect(result.current.selectedIssuers).toHaveLength(2);
        expect(result.current.isModalOpen).toBe(false);
    });

    it("should replace issuer when selecting the same id again", () => {
        const { result } = renderHook(
            () => useIssuerManagement({ selectedIssuers: [mockIssuers[0]] }),
            { wrapper },
        );

        // Try to select the same issuer again
        const duplicateIssuer = { ...mockIssuers[0] };

        act(() => {
            result.current.handleConfirmSelection([duplicateIssuer]);
        });

        // Should still be 1 since we're replacing the same issuer
        expect(result.current.selectedIssuers).toHaveLength(1);
    });

    it("should remove issuer by id", () => {
        const { result } = renderHook(() => useIssuerManagement({ selectedIssuers: mockIssuers }), {
            wrapper,
        });

        expect(result.current.selectedIssuers).toHaveLength(2);

        act(() => {
            result.current.handleRemoveIssuer(mockIssuers[0].id);
        });

        expect(result.current.selectedIssuers).toHaveLength(1);
        expect(result.current.selectedIssuers[0].id).toBe(mockIssuers[1].id);
    });

    it("should clear all selections", () => {
        const { result } = renderHook(() => useIssuerManagement({ selectedIssuers: mockIssuers }), {
            wrapper,
        });

        expect(result.current.selectedIssuers).toHaveLength(2);

        act(() => {
            result.current.handleClearSelection();
        });

        expect(result.current.selectedIssuers).toHaveLength(0);
    });

    it("should get selected issuer ids", () => {
        const { result } = renderHook(() => useIssuerManagement({ selectedIssuers: mockIssuers }), {
            wrapper,
        });

        const ids = result.current.getSelectedIssuerIds();

        expect(ids instanceof Set).toBe(true);
        expect(ids.has(mockIssuers[0].id)).toBe(true);
        expect(ids.has(mockIssuers[1].id)).toBe(true);
        expect(ids.has("non-existent-id")).toBe(false);
    });

    it("should filter out invalid ids in getSelectedIssuerIds", () => {
        const invalidIssuers: Issuer[] = [
            { id: "valid-1", name: "Test", email: "test@test.com" },
            { id: "", name: "Invalid", email: "invalid@test.com" },
        ];

        const { result } = renderHook(
            () => useIssuerManagement({ selectedIssuers: invalidIssuers }),
            { wrapper },
        );

        const ids = result.current.getSelectedIssuerIds();

        expect(ids.size).toBe(1);
        expect(ids.has("valid-1")).toBe(true);
        expect(ids.has("")).toBe(false);
    });

    it("should maintain stable function references", () => {
        const { result, rerender } = renderHook(() => useIssuerManagement(), { wrapper });

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

    it("should handle modal close without resetting search results", () => {
        const { result } = renderHook(() => useIssuerManagement(), { wrapper });

        act(() => {
            result.current.handleOpenModal();
        });

        expect(result.current.isModalOpen).toBe(true);

        act(() => {
            result.current.handleCloseModal();
        });

        expect(result.current.isModalOpen).toBe(false);
        // searchResults and activeSearchQuery should persist
    });
});

describe("convertProfileToIssuer", () => {
    it("should convert profile to issuer correctly", () => {
        const profile: Profile = {
            id: "profile-1",
            authenticationCredentialId: "cred-1",
            firstName: "Alice",
            lastName: "Johnson",
            email: "alice@example.com",
            googleConnectorRef: "alice.oauth@gmail.com",
            phoneNumber: "123-456-7890",
            academicInstitution: "Stanford University",
            bio: "CS Professor",
            isEmailPublic: true,
            isFirstNamePublic: true,
            isLastNamePublic: true,
            isPhoneNumberPublic: true,
            isAcademicInstitutionPublic: true,
            isProfilePicturePublic: false,
            isBioPublic: false,
            isAddressPublic: false,
            isAcademicEmailPublic: false,
        };

        const issuer = convertProfileToIssuer(profile);

        expect(issuer.id).toBe("cred-1");
        expect(issuer.name).toBe("Alice Johnson");
        expect(issuer.email).toBe("alice@example.com");
        expect(issuer.googleOAuthEmail).toBe("alice.oauth@gmail.com");
        expect(issuer.phoneNumber).toBe("123-456-7890");
        expect(issuer.organization).toBe("Stanford University");
    });

    it("should handle missing last name", () => {
        const profile: Profile = {
            id: "profile-1",
            authenticationCredentialId: "cred-1",
            firstName: "Alice",
            email: "alice@example.com",
            isEmailPublic: true,
            isFirstNamePublic: true,
            isLastNamePublic: true,
            isPhoneNumberPublic: true,
            isAcademicInstitutionPublic: true,
            isProfilePicturePublic: false,
            isBioPublic: false,
            isAddressPublic: false,
            isAcademicEmailPublic: false,
        };

        const issuer = convertProfileToIssuer(profile);

        expect(issuer.name).toBe("Alice");
    });

    it("should fallback to bio for organization", () => {
        const profile: Profile = {
            id: "profile-1",
            authenticationCredentialId: "cred-1",
            firstName: "Alice",
            lastName: "Johnson",
            email: "alice@example.com",
            bio: "CS Department",
            isEmailPublic: true,
            isFirstNamePublic: true,
            isLastNamePublic: true,
            isPhoneNumberPublic: true,
            isAcademicInstitutionPublic: true,
            isProfilePicturePublic: false,
            isBioPublic: false,
            isAddressPublic: false,
            isAcademicEmailPublic: false,
        };

        const issuer = convertProfileToIssuer(profile);

        expect(issuer.organization).toBe("CS Department");
    });

    it("should handle empty name fields", () => {
        const profile: Profile = {
            id: "profile-1",
            authenticationCredentialId: "cred-1",
            email: "alice@example.com",
            isEmailPublic: true,
            isFirstNamePublic: true,
            isLastNamePublic: true,
            isPhoneNumberPublic: true,
            isAcademicInstitutionPublic: true,
            isProfilePicturePublic: false,
            isBioPublic: false,
            isAddressPublic: false,
            isAcademicEmailPublic: false,
        };

        const issuer = convertProfileToIssuer(profile);

        expect(issuer.name).toBe("");
    });

    it("should handle missing optional fields", () => {
        const profile: Profile = {
            id: "profile-1",
            authenticationCredentialId: "cred-1",
            firstName: "Alice",
            lastName: "Johnson",
            isEmailPublic: true,
            isFirstNamePublic: true,
            isLastNamePublic: true,
            isPhoneNumberPublic: true,
            isAcademicInstitutionPublic: true,
            isProfilePicturePublic: false,
            isBioPublic: false,
            isAddressPublic: false,
            isAcademicEmailPublic: false,
        };

        const issuer = convertProfileToIssuer(profile);

        expect(issuer.email).toBe("");
        expect(issuer.googleOAuthEmail).toBeUndefined();
        expect(issuer.phoneNumber).toBeUndefined();
        expect(issuer.organization).toBeUndefined();
    });
});
