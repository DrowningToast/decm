import { useState, useCallback } from "react";
import type { Profile } from "@/services/AuthService/AuthService";
import { useSearchIssuer } from "./useSearchIssuer";

export interface Issuer {
    id: string;
    name: string;
    email: string;
    organization?: string;
}

export interface UseIssuerManagementProps {
    selectedIssuers?: Issuer[];
}

export interface UseIssuerManagementReturn {
    selectedIssuers: Issuer[];
    searchQuery: string;
    isSearching: boolean;
    searchResults: Issuer[];
    isModalOpen: boolean;
    handleSearch: () => Promise<void>;
    handleSearchQueryChange: (query: string) => void;
    handleOpenModal: () => void;
    handleCloseModal: () => void;
    handleConfirmSelection: (newIssuers: Issuer[]) => void;
    handleRemoveIssuer: (issuerId: string) => void;
    handleClearSelection: () => void;
    getSelectedIssuerIds: () => Set<string>;
}

// Helper function to convert VerifiedIssuer to Issuer
export const convertProfileToIssuer = (profile: Profile): Issuer => {
    const fullName =
        [profile.firstName, profile.lastName].filter(Boolean).join(" ").trim() || "Unknown Name";

    return {
        id: profile.authenticationCredentialId || "",
        name: fullName,
        email: profile.email || "",
        organization: profile.academicInstitution || profile.bio || undefined,
    };
};

export const useIssuerManagement = ({
    selectedIssuers: initialSelectedIssuers,
}: UseIssuerManagementProps = {}): UseIssuerManagementReturn => {
    const [selectedIssuers, setSelectedIssuers] = useState<Issuer[]>(initialSelectedIssuers || []);
    const [searchQuery, setSearchQuery] = useState("");
    const [isSearching, setIsSearching] = useState(false);
    const [searchResults, setSearchResults] = useState<Issuer[]>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);

    // Use the backend search hook
    const { data: verifiedIssuers } = useSearchIssuer({
        search: searchQuery,
    });

    const handleSearch = useCallback(async () => {
        if (!searchQuery.trim()) return;

        setIsSearching(true);
        try {
            // Wait for the backend search results
            // The results will be available through the verifiedIssuers from the hook
            // Add a small delay to ensure the query completes
            await new Promise((resolve) => setTimeout(resolve, 500));

            if (verifiedIssuers && verifiedIssuers.length > 0) {
                const results = verifiedIssuers.map(convertProfileToIssuer);
                setSearchResults(results);
            } else {
                setSearchResults([]);
            }
            setIsModalOpen(true);
        } catch (error) {
            console.error("Error searching issuers:", error);
            setSearchResults([]);
        } finally {
            setIsSearching(false);
        }
    }, [searchQuery, verifiedIssuers]);

    const handleSearchQueryChange = useCallback((query: string) => {
        setSearchQuery(query);
    }, []);

    const handleOpenModal = useCallback(() => {
        setIsModalOpen(true);
    }, []);

    const handleCloseModal = useCallback(() => {
        setIsModalOpen(false);
        setSearchResults([]);
    }, []);

    const handleConfirmSelection = useCallback(
        (newIssuers: Issuer[]) => {
            // Merge with existing selections (avoid duplicates)
            const merged: Issuer[] = [...newIssuers];
            setSelectedIssuers(merged);
            handleCloseModal();
        },
        [handleCloseModal],
    );

    const handleRemoveIssuer = useCallback((issuerId: string) => {
        setSelectedIssuers((prev) => prev.filter((issuer) => issuer.id !== issuerId));
    }, []);

    const handleClearSelection = useCallback(() => {
        setSelectedIssuers([]);
    }, []);

    const getSelectedIssuerIds = useCallback(() => {
        const issuerIds = selectedIssuers
            .map((issuer) => issuer.id)
            .filter((id): id is string => typeof id === "string" && id.length > 0);

        return new Set(issuerIds);
    }, [selectedIssuers]);

    return {
        selectedIssuers,
        searchQuery,
        isSearching,
        searchResults,
        isModalOpen,
        handleSearch,
        handleSearchQueryChange,
        handleOpenModal,
        handleCloseModal,
        handleConfirmSelection,
        handleRemoveIssuer,
        handleClearSelection,
        getSelectedIssuerIds,
    };
};
