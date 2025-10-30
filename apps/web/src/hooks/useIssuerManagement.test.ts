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
