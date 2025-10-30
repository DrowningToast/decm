import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { useIssuerManagement } from './useIssuerManagement';
import type { EntityProfile } from '@decm/api';

describe('useIssuerManagement', () => {
    const mockVerifiedIssuers: EntityProfile[] = [
        {
            authentication_credential_id: '1',
            first_name: 'John',
            last_name: 'Doe',
            email: 'john.doe@example.com',
            academic_institution: 'Test University',
        } as EntityProfile,
        {
            authentication_credential_id: '2',
            first_name: 'Jane',
            last_name: 'Smith',
            email: 'jane.smith@example.com',
            academic_institution: 'Another University',
        } as EntityProfile,
    ];

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should initialize with default values', () => {
        const { result } = renderHook(() => useIssuerManagement());

        expect(result.current.selectedIssuers).toEqual([]);
        expect(result.current.searchQuery).toBe('');
        expect(result.current.isSearching).toBe(false);
        expect(result.current.searchResults).toEqual([]);
        expect(result.current.isModalOpen).toBe(false);
    });

    it('should initialize with initial selected issuers', () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: mockVerifiedIssuers }),
        );

        expect(result.current.selectedIssuers).toEqual(mockVerifiedIssuers);
    });

    it('should update search query', () => {
        const { result } = renderHook(() => useIssuerManagement());

        act(() => {
            result.current.handleSearchQueryChange('test query');
        });

        expect(result.current.searchQuery).toBe('test query');
    });

    it('should open and close modal', () => {
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

    it('should search issuers using default search function', async () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: mockVerifiedIssuers }),
        );

        act(() => {
            result.current.handleSearchQueryChange('John');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults.length).toBeGreaterThan(0);
            expect(result.current.isModalOpen).toBe(true);
        });

        expect(result.current.searchResults[0].name).toContain('John');
    });

    it('should use custom search function when provided', async () => {
        const customSearchFunction = vi.fn().mockResolvedValue([
            {
                id: '3',
                name: 'Custom Issuer',
                email: 'custom@example.com',
            },
        ]);

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: customSearchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange('custom');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        expect(customSearchFunction).toHaveBeenCalledWith('custom');
        expect(result.current.searchResults[0].name).toBe('Custom Issuer');
    });

    it('should not search when query is empty', async () => {
        const searchFunction = vi.fn();

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange('');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        expect(searchFunction).not.toHaveBeenCalled();
    });

    it('should confirm selection and merge with existing', () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: [mockVerifiedIssuers[0]] }),
        );

        const newIssuers = [
            {
                id: '2',
                name: 'Jane Smith',
                email: 'jane.smith@example.com',
            },
        ];

        act(() => {
            result.current.handleConfirmSelection(newIssuers);
        });

        expect(result.current.selectedIssuers.length).toBe(2);
        expect(result.current.isModalOpen).toBe(false);
    });

    it('should not add duplicate issuers', () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: [mockVerifiedIssuers[0]] }),
        );

        const duplicateIssuer = {
            id: mockVerifiedIssuers[0].authentication_credential_id || '1',
            name: 'John Doe',
            email: 'john.doe@example.com',
        };

        act(() => {
            result.current.handleConfirmSelection([duplicateIssuer]);
        });

        expect(result.current.selectedIssuers.length).toBe(1);
    });

    it('should remove issuer', () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: mockVerifiedIssuers }),
        );

        act(() => {
            result.current.handleRemoveIssuer('1');
        });

        expect(result.current.selectedIssuers.length).toBe(1);
        expect(result.current.selectedIssuers[0].authentication_credential_id).toBe('2');
    });

    it('should clear selection', () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: mockVerifiedIssuers }),
        );

        act(() => {
            result.current.handleClearSelection();
        });

        expect(result.current.selectedIssuers).toEqual([]);
    });

    it('should get selected issuer IDs', () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: mockVerifiedIssuers }),
        );

        const selectedIds = result.current.getSelectedIssuerIds();

        expect(selectedIds.size).toBe(2);
        expect(selectedIds.has('1')).toBe(true);
        expect(selectedIds.has('2')).toBe(true);
    });

    it('should handle search errors gracefully', async () => {
        const errorSearchFunction = vi.fn().mockRejectedValue(new Error('Search failed'));

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: errorSearchFunction }),
        );

        const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

        act(() => {
            result.current.handleSearchQueryChange('test');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.isSearching).toBe(false);
        });

        expect(consoleSpy).toHaveBeenCalled();
        consoleSpy.mockRestore();
    });

    it('should set searching state during search', async () => {
        const delayedSearchFunction = vi.fn().mockImplementation(
            async () =>
                new Promise((resolve) =>
                    setTimeout(() => resolve([{ id: '1', name: 'Test', email: 'test@example.com' }]), 100),
                ),
        );

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: delayedSearchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange('test');
        });

        const searchPromise = act(async () => {
            await result.current.handleSearch();
        });

        // Check that isSearching is true during search
        expect(result.current.isSearching).toBe(true);

        await searchPromise;

        await waitFor(() => {
            expect(result.current.isSearching).toBe(false);
        });
    });

    it('should handle empty verified issuers list', async () => {
        const { result } = renderHook(() => useIssuerManagement({ verifiedIssuers: [] }));

        act(() => {
            result.current.handleSearchQueryChange('test');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults).toEqual([]);
        });
    });
});

describe('useIssuerManagement - Additional Edge Cases', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should handle concurrent searches', async () => {
        let resolveSearch1: (value: any) => void;
        let resolveSearch2: (value: any) => void;

        const delayedSearch = vi.fn().mockImplementation(() => {
            return new Promise((resolve) => {
                if (delayedSearch.mock.calls.length === 1) {
                    resolveSearch1 = resolve;
                } else {
                    resolveSearch2 = resolve;
                }
            });
        });

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction: delayedSearch }),
        );

        act(() => {
            result.current.handleSearchQueryChange('query1');
        });

        const search1Promise = act(async () => {
            await result.current.handleSearch();
        });

        act(() => {
            result.current.handleSearchQueryChange('query2');
        });

        const search2Promise = act(async () => {
            await result.current.handleSearch();
        });

        // Resolve second search first
        act(() => {
            resolveSearch2!([{ id: '2', name: 'Result 2', email: 'test2@example.com' }]);
        });

        await search2Promise;

        // Then resolve first search
        act(() => {
            resolveSearch1!([{ id: '1', name: 'Result 1', email: 'test1@example.com' }]);
        });

        await search1Promise;

        expect(delayedSearch).toHaveBeenCalledTimes(2);
    });

    it('should handle profile with missing optional fields', () => {
        const profileWithMinimalData: EntityProfile = {
            authentication_credential_id: '1',
            first_name: 'John',
            last_name: '',
            email: 'john@example.com',
        } as EntityProfile;

        const { result } = renderHook(() =>
            useIssuerManagement({
                verifiedIssuers: [profileWithMinimalData],
                selectedIssuers: [profileWithMinimalData],
            }),
        );

        expect(result.current.selectedIssuers[0].first_name).toBe('John');
        expect(result.current.selectedIssuers[0].email).toBe('john@example.com');
    });

    it('should handle selection of issuer with empty ID', () => {
        const issuerWithEmptyId = {
            id: '',
            name: 'Test Issuer',
            email: 'test@example.com',
        };

        const { result } = renderHook(() => useIssuerManagement());

        act(() => {
            result.current.handleConfirmSelection([issuerWithEmptyId]);
        });

        const selectedIds = result.current.getSelectedIssuerIds();
        expect(selectedIds.has('')).toBe(true);
    });

    it('should handle search query with whitespace', async () => {
        const searchFunction = vi.fn().mockResolvedValue([]);

        const { result } = renderHook(() =>
            useIssuerManagement({ searchFunction }),
        );

        act(() => {
            result.current.handleSearchQueryChange('  test  ');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        expect(searchFunction).toHaveBeenCalledWith('test');
    });

    it('should handle removing non-existent issuer gracefully', () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ selectedIssuers: mockVerifiedIssuers }),
        );

        const initialCount = result.current.selectedIssuers.length;

        act(() => {
            result.current.handleRemoveIssuer('non-existent-id');
        });

        expect(result.current.selectedIssuers.length).toBe(initialCount);
    });

    it('should handle profile name construction with only first name', async () => {
        const profileWithOnlyFirstName: EntityProfile = {
            authentication_credential_id: '1',
            first_name: 'John',
            last_name: '',
            email: 'john@example.com',
        } as EntityProfile;

        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: [profileWithOnlyFirstName] }),
        );

        act(() => {
            result.current.handleSearchQueryChange('John');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults[0].name).toBe('John');
        });
    });

    it('should handle profile with no name fields', async () => {
        const profileWithNoName: EntityProfile = {
            authentication_credential_id: '1',
            first_name: '',
            last_name: '',
            email: 'unknown@example.com',
        } as EntityProfile;

        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: [profileWithNoName] }),
        );

        act(() => {
            result.current.handleSearchQueryChange('unknown');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults[0].name).toBe('Unknown Name');
        });
    });

    it('should filter by organization when searching', async () => {
        const profileWithOrganization: EntityProfile = {
            authentication_credential_id: '1',
            first_name: 'John',
            last_name: 'Doe',
            email: 'john@example.com',
            academic_institution: 'MIT',
        } as EntityProfile;

        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: [profileWithOrganization] }),
        );

        act(() => {
            result.current.handleSearchQueryChange('MIT');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults.length).toBeGreaterThan(0);
            expect(result.current.searchResults[0].organization).toBe('MIT');
        });
    });

    it('should handle case-insensitive search', async () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: mockVerifiedIssuers }),
        );

        act(() => {
            result.current.handleSearchQueryChange('JOHN');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults.length).toBeGreaterThan(0);
        });
    });

    it('should maintain modal state across multiple open/close cycles', () => {
        const { result } = renderHook(() => useIssuerManagement());

        act(() => {
            result.current.handleOpenModal();
        });
        expect(result.current.isModalOpen).toBe(true);

        act(() => {
            result.current.handleCloseModal();
        });
        expect(result.current.isModalOpen).toBe(false);

        act(() => {
            result.current.handleOpenModal();
        });
        expect(result.current.isModalOpen).toBe(true);
    });

    it('should clear search results when modal is closed', () => {
        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: mockVerifiedIssuers }),
        );

        act(() => {
            result.current.handleSearchQueryChange('test');
        });

        act(async () => {
            await result.current.handleSearch();
        });

        waitFor(() => {
            expect(result.current.searchResults.length).toBeGreaterThan(0);
        });

        act(() => {
            result.current.handleCloseModal();
        });

        expect(result.current.searchResults).toEqual([]);
    });

    it('should handle selection with bio as organization fallback', async () => {
        const profileWithBio: EntityProfile = {
            authentication_credential_id: '1',
            first_name: 'John',
            last_name: 'Doe',
            email: 'john@example.com',
            bio: 'Software Engineer at TechCorp',
        } as EntityProfile;

        const { result } = renderHook(() =>
            useIssuerManagement({ verifiedIssuers: [profileWithBio] }),
        );

        act(() => {
            result.current.handleSearchQueryChange('John');
        });

        await act(async () => {
            await result.current.handleSearch();
        });

        await waitFor(() => {
            expect(result.current.searchResults[0].organization).toBe('Software Engineer at TechCorp');
        });
    });
});