import { useState, useCallback } from "react";
import type { EntityProfile } from "@decm/api";

export interface Issuer {
    id: string;
    name: string;
    email: string;
    organization?: string;
}

export interface UseIssuerManagementProps {
    searchFunction?: (query: string) => Promise<Issuer[]>;
    verifiedIssuers?: EntityProfile[];
    selectedIssuers?: EntityProfile[];
}

export interface UseIssuerManagementReturn {
    selectedIssuers: EntityProfile[];
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

export const useIssuerManagement = ({
    searchFunction,
    verifiedIssuers,
    selectedIssuers: initialSelectedIssuers,
}: UseIssuerManagementProps = {}): UseIssuerManagementReturn => {
    const [selectedIssuers, setSelectedIssuers] = useState<EntityProfile[]>(
        initialSelectedIssuers || [],
    );
    const [searchQuery, setSearchQuery] = useState("");
    const [isSearching, setIsSearching] = useState(false);
    const [searchResults, setSearchResults] = useState<Issuer[]>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);

    // Helper function to convert EntityProfile to Issuer
    const convertProfileToIssuer = useCallback((profile: EntityProfile): Issuer => {
        const fullName =
            [profile.first_name, profile.last_name].filter(Boolean).join(" ").trim() ||
            "Unknown Name";

        return {
            id: profile.authentication_credential_id || "",
            name: fullName,
            email: profile.email || "",
            organization: profile.academic_institution || profile.bio || undefined,
        };
    }, []);

    // Default search function using verified issuers
    const defaultSearchFunction = useCallback(
        async (query: string): Promise<Issuer[]> => {
            // Simulate API delay for better UX
            await new Promise((resolve) => setTimeout(resolve, 300));

            if (!verifiedIssuers || verifiedIssuers.length === 0) {
                return [];
            }

            // Convert EntityProfile to Issuer format
            const issuers = verifiedIssuers.map(convertProfileToIssuer);

            // Filter based on query
            const filteredIssuers = issuers.filter(
                (issuer) =>
                    issuer.name.toLowerCase().includes(query.toLowerCase()) ||
                    issuer.email.toLowerCase().includes(query.toLowerCase()) ||
                    issuer.organization?.toLowerCase().includes(query.toLowerCase()),
            );

            return filteredIssuers;
        },
        [verifiedIssuers, convertProfileToIssuer],
    );

    const handleSearch = useCallback(async () => {
        if (!searchQuery.trim()) return;

        setIsSearching(true);
        try {
            const searchFn = searchFunction || defaultSearchFunction;
            const results = await searchFn(searchQuery.trim());
            setSearchResults(results);
            setIsModalOpen(true);
        } catch (error) {
            console.error("Error searching issuers:", error);
            // Could add error handling here (toast, alert, etc.)
        } finally {
            setIsSearching(false);
        }
    }, [searchQuery, searchFunction, defaultSearchFunction]);

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
            const merged = [...selectedIssuers];
            newIssuers.forEach((issuer) => {
                if (
                    !merged.some((existing) => existing.authentication_credential_id === issuer.id)
                ) {
                    merged.push(issuer as unknown as EntityProfile);
                }
            });

            setSelectedIssuers(merged);
            handleCloseModal();
        },
        [selectedIssuers, handleCloseModal],
    );

    const handleRemoveIssuer = useCallback((issuerId: string) => {
        setSelectedIssuers((prev) =>
            prev.filter((issuer) => issuer.authentication_credential_id !== issuerId),
        );
    }, []);

    const handleClearSelection = useCallback(() => {
        setSelectedIssuers([]);
    }, []);

    const getSelectedIssuerIds = useCallback(() => {
        const issuerIds = selectedIssuers
            .map((issuer) => issuer.authentication_credential_id)
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
