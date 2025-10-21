import { useState, useCallback } from "react";

export interface Issuer {
    id: string;
    name: string;
    email: string;
    organization?: string;
}

export interface UseIssuerManagementProps {
    searchFunction?: (query: string) => Promise<Issuer[]>;
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

export const useIssuerManagement = ({
    searchFunction,
}: UseIssuerManagementProps = {}): UseIssuerManagementReturn => {
    const [selectedIssuers, setSelectedIssuers] = useState<Issuer[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [isSearching, setIsSearching] = useState(false);
    const [searchResults, setSearchResults] = useState<Issuer[]>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);

    // Default mock search function
    const defaultSearchFunction = async (query: string): Promise<Issuer[]> => {
        // Simulate API delay
        await new Promise((resolve) => setTimeout(resolve, 1000));

        // Mock results - replace with actual API data
        const mockResults: Issuer[] = [
            {
                id: "1",
                name: "John Doe",
                email: "john.doe@university.edu",
                organization: "Computer Science Department",
            },
            {
                id: "2",
                name: "Jane Smith",
                email: "jane.smith@university.edu",
                organization: "Engineering Department",
            },
            {
                id: "3",
                name: "Dr. Michael Brown",
                email: "m.brown@university.edu",
                organization: "Mathematics Department",
            },
            {
                id: "4",
                name: "Prof. Sarah Wilson",
                email: "s.wilson@university.edu",
                organization: "Physics Department",
            },
        ];

        // Filter mock results based on query
        return mockResults.filter(
            (issuer) =>
                issuer.name.toLowerCase().includes(query.toLowerCase()) ||
                issuer.email.toLowerCase().includes(query.toLowerCase()) ||
                issuer.organization?.toLowerCase().includes(query.toLowerCase()),
        );
    };

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
    }, [searchQuery, searchFunction]);

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
                if (!merged.some((existing) => existing.id === issuer.id)) {
                    merged.push(issuer);
                }
            });

            setSelectedIssuers(merged);
            handleCloseModal();
        },
        [selectedIssuers, handleCloseModal],
    );

    const handleRemoveIssuer = useCallback((issuerId: string) => {
        setSelectedIssuers((prev) => prev.filter((issuer) => issuer.id !== issuerId));
    }, []);

    const handleClearSelection = useCallback(() => {
        setSelectedIssuers([]);
    }, []);

    const getSelectedIssuerIds = useCallback(() => {
        return new Set(selectedIssuers.map((issuer) => issuer.id));
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
