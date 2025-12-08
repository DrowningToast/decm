import { useState, useCallback, useEffect } from "react";
import type { Profile } from "@/services/AuthService/AuthService";
import { useSearchIssuer } from "./useSearchIssuer";

export interface Issuer {
    id: string;
    name: string;
    email: string;
    googleOAuthEmail?: string;
    phoneNumber?: string;
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
    const fullName = [profile.firstName, profile.lastName].filter(Boolean).join(" ").trim() || "";

    return {
        id: profile.authenticationCredentialId || "",
        name: fullName,
        email: profile.email || "",
        googleOAuthEmail: profile.googleConnectorRef || undefined,
        phoneNumber: profile.phoneNumber || undefined,
        organization: profile.academicInstitution || profile.bio || undefined,
    };
};

export const useIssuerManagement = ({
    selectedIssuers: initialSelectedIssuers,
}: UseIssuerManagementProps = {}): UseIssuerManagementReturn => {
    const [selectedIssuers, setSelectedIssuers] = useState<Issuer[]>(initialSelectedIssuers || []);
    const [searchQuery, setSearchQuery] = useState("");
    const [activeSearchQuery, setActiveSearchQuery] = useState<string>("");
    const [isSearching, setIsSearching] = useState(false);
    const [searchResults, setSearchResults] = useState<Issuer[]>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);

    // Use the backend search hook - only fetches when activeSearchQuery changes (on button click)
    const { data: verifiedIssuers, isLoading: isQueryLoading } = useSearchIssuer({
        search: activeSearchQuery,
    });

    // Update isSearching based on query loading state
    useEffect(() => {
        setIsSearching(isQueryLoading);
    }, [isQueryLoading]);

    const handleSearch = useCallback(async () => {
        if (!searchQuery.trim()) return;

        const trimmedQuery = searchQuery.trim();
        setIsSearching(true);

        // If same query and we have cached data, show it immediately
        if (activeSearchQuery === trimmedQuery && verifiedIssuers !== undefined) {
            if (verifiedIssuers && verifiedIssuers.length > 0) {
                const results = verifiedIssuers.map(convertProfileToIssuer);
                setSearchResults(results);
            } else {
                setSearchResults([]);
            }
            setIsModalOpen(true);
            setIsSearching(false);
        } else {
            // Set the active search query to trigger the useSearchIssuer hook
            setActiveSearchQuery(trimmedQuery);
        }
    }, [searchQuery, activeSearchQuery, verifiedIssuers]);

    // Effect to handle search results when the query completes
    useEffect(() => {
        if (activeSearchQuery && !isQueryLoading && verifiedIssuers !== undefined) {
            if (verifiedIssuers && verifiedIssuers.length > 0) {
                const results = verifiedIssuers.map(convertProfileToIssuer);
                setSearchResults(results);
            } else {
                setSearchResults([]);
            }
            setIsModalOpen(true);
            setIsSearching(false);
        }
    }, [activeSearchQuery, verifiedIssuers, isQueryLoading]);

    const handleSearchQueryChange = useCallback((query: string) => {
        setSearchQuery(query);
    }, []);

    const handleOpenModal = useCallback(() => {
        setIsModalOpen(true);
    }, []);

    const handleCloseModal = useCallback(() => {
        setIsModalOpen(false);
        // Don't reset activeSearchQuery or searchResults - keep them for next time modal opens
    }, []);

    const handleConfirmSelection = useCallback(
        (newIssuers: Issuer[]) => {
            // Merge with existing selections (avoid duplicates)
            // Get IDs of issuers that were in the search results (to be updated/replaced)
            const searchResultIds = new Set(newIssuers.map((issuer) => issuer.id));

            // Keep existing issuers that were NOT in the search results
            const existingIssuersNotInSearch = selectedIssuers.filter(
                (issuer) => !searchResultIds.has(issuer.id),
            );

            // Merge: existing issuers not in search + new/updated issuers from search
            const merged: Issuer[] = [...existingIssuersNotInSearch, ...newIssuers];
            setSelectedIssuers(merged);
            handleCloseModal();
        },
        [handleCloseModal, selectedIssuers],
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
