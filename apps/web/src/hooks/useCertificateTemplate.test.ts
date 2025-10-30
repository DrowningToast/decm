import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { useCertificateTemplate } from './useCertificateTemplate';

describe('useCertificateTemplate', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should initialize with default values', () => {
        const { result } = renderHook(() => useCertificateTemplate());

        expect(result.current.svgFile).toBeNull();
        expect(result.current.svgPreview).toBe('');
        expect(result.current.detectedKeywords).toEqual([]);
        expect(result.current.availableKeywords.length).toBeGreaterThan(0);
        expect(result.current.hasMissingMandatory).toBe(true);
    });

    it('should initialize with custom search keys', () => {
        const customSearchKeys = ['{{ customKey }}'];
        const { result } = renderHook(() =>
            useCertificateTemplate({ searchKeys: customSearchKeys }),
        );

        expect(result.current.availableKeywords).toBeDefined();
    });

    it('should initialize with custom certificate dimensions', () => {
        const { result } = renderHook(() =>
            useCertificateTemplate({ certWidth: 800, certHeight: 600 }),
        );

        expect(result.current.availableKeywords).toBeDefined();
    });

    it('should handle file selection', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const svgContent = '<svg><g id="{{ name }}"></g><g id="{{ eventName }}"></g></svg>';
        const file = new File([svgContent], 'test.svg', { type: 'image/svg+xml' });

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        await waitFor(() => {
            expect(result.current.svgFile).toBe(file);
            expect(result.current.svgPreview).toBeTruthy();
        });
    });

    it('should detect keywords in SVG', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const svgContent = '<svg><g id="{{ name }}"><text>Name</text></g><g id="{{ eventName }}"><text>Event</text></g></svg>';
        const file = new File([svgContent], 'test.svg', { type: 'image/svg+xml' });

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        await waitFor(() => {
            expect(result.current.detectedKeywords.length).toBeGreaterThan(0);
        });
    });

    it('should clear template when clearTemplate is called', () => {
        const { result } = renderHook(() => useCertificateTemplate());

        // Set some values first
        act(() => {
            const svgContent = '<svg><g id="{{ name }}"></g></svg>';
            const file = new File([svgContent], 'test.svg', { type: 'image/svg+xml' });
            const fileInput = document.createElement('input');
            fileInput.type = 'file';
            Object.defineProperty(fileInput, 'files', {
                value: [file],
                writable: false,
            });
            const event = {
                target: fileInput,
            } as unknown as React.ChangeEvent<HTMLInputElement>;
            result.current.handleFileSelect(event);
        });

        act(() => {
            result.current.clearTemplate();
        });

        expect(result.current.svgFile).toBeNull();
        expect(result.current.svgPreview).toBe('');
        expect(result.current.detectedKeywords).toEqual([]);
    });

    it('should identify missing mandatory keywords', () => {
        const { result } = renderHook(() => useCertificateTemplate());

        // Initially, no keywords detected, so mandatory ones are missing
        expect(result.current.hasMissingMandatory).toBe(true);
        expect(result.current.missingMandatoryKeywords.length).toBeGreaterThan(0);
    });

    it('should handle invalid SVG files gracefully', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const invalidContent = 'not an svg';
        const file = new File([invalidContent], 'test.txt', { type: 'text/plain' });

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        // Should not crash, but may not set preview
        expect(result.current.svgFile).toBe(file);
    });

    it('should handle empty file selection', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        // Should not crash and should not update state
        expect(result.current.svgFile).toBeNull();
    });

    it('should provide file input ref', () => {
        const { result } = renderHook(() => useCertificateTemplate());

        expect(result.current.fileInputRef).toBeDefined();
        expect(result.current.fileInputRef.current).toBeNull();
    });

    it('should provide SVG temp ref', () => {
        const { result } = renderHook(() => useCertificateTemplate());

        expect(result.current.svgTempRef).toBeDefined();
        expect(result.current.svgTempRef.current).toBeNull();
    });
});

describe('useCertificateTemplate - Additional Edge Cases', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should handle SVG with all mandatory keywords present', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const svgContent = '<svg><g id="{{ name }}"></g><g id="{{ eventName }}"></g></svg>';
        const file = new File([svgContent], 'test.svg', { type: 'image/svg+xml' });

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        await waitFor(() => {
            expect(result.current.hasMissingMandatory).toBe(false);
            expect(result.current.missingMandatoryKeywords.length).toBe(0);
        });
    });

    it('should handle SVG with optional keywords', async () => {
        const { result } = renderHook(() =>
            useCertificateTemplate({
                searchKeys: ['{{ name }}', '{{ eventName }}', '{{ academicInstitutionName }}'],
            }),
        );

        const svgContent =
            '<svg><g id="{{ name }}"></g><g id="{{ eventName }}"></g><g id="{{ academicInstitutionName }}"></g></svg>';
        const file = new File([svgContent], 'test.svg', { type: 'image/svg+xml' });

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        await waitFor(() => {
            expect(result.current.detectedKeywords.length).toBe(3);
        });
    });

    it('should handle large SVG files', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const largeContent = '<svg><g id="{{ name }}"></g>'.repeat(1000) + '</svg>';
        const file = new File([largeContent], 'large.svg', { type: 'image/svg+xml' });

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        expect(result.current.svgFile).toBe(file);
    });

    it('should handle SVG with special characters in keywords', async () => {
        const { result } = renderHook(() =>
            useCertificateTemplate({
                searchKeys: ['{{ name }}', '{{ event-name }}'],
            }),
        );

        const svgContent = '<svg><g id="{{ name }}"></g><g id="{{ event-name }}"></g></svg>';
        const file = new File([svgContent], 'test.svg', { type: 'image/svg+xml' });

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        expect(result.current.svgFile).toBe(file);
    });

    it('should handle file reader error gracefully', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const file = new File(['test'], 'test.svg', { type: 'image/svg+xml' });
        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        // Mock FileReader to throw error
        const originalFileReader = window.FileReader;
        window.FileReader = vi.fn().mockImplementation(() => ({
            readAsText: vi.fn().mockImplementation(function () {
                if (this.onerror) {
                    this.onerror(new Error('Read error'));
                }
            }),
            onload: null,
            onerror: null,
        })) as any;

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        expect(result.current.svgFile).toBe(file);

        // Restore FileReader
        window.FileReader = originalFileReader;
    });

    it('should handle multiple file selections', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const svgContent1 = '<svg><g id="{{ name }}"></g></svg>';
        const file1 = new File([svgContent1], 'test1.svg', { type: 'image/svg+xml' });

        const fileInput1 = document.createElement('input');
        fileInput1.type = 'file';
        Object.defineProperty(fileInput1, 'files', {
            value: [file1],
            writable: false,
        });

        await act(async () => {
            result.current.handleFileSelect({
                target: fileInput1,
            } as unknown as React.ChangeEvent<HTMLInputElement>);
        });

        expect(result.current.svgFile).toBe(file1);

        const svgContent2 = '<svg><g id="{{ eventName }}"></g></svg>';
        const file2 = new File([svgContent2], 'test2.svg', { type: 'image/svg+xml' });

        const fileInput2 = document.createElement('input');
        fileInput2.type = 'file';
        Object.defineProperty(fileInput2, 'files', {
            value: [file2],
            writable: false,
        });

        await act(async () => {
            result.current.handleFileSelect({
                target: fileInput2,
            } as unknown as React.ChangeEvent<HTMLInputElement>);
        });

        expect(result.current.svgFile).toBe(file2);
    });

    it('should handle custom certificate dimensions', () => {
        const { result } = renderHook(() =>
            useCertificateTemplate({
                certWidth: 800,
                certHeight: 600,
            }),
        );

        expect(result.current.availableKeywords).toBeDefined();
    });

    it('should clear file input ref when clearing template', () => {
        const { result } = renderHook(() => useCertificateTemplate());

        // Create a mock file input
        const mockFileInput = document.createElement('input');
        mockFileInput.type = 'file';
        mockFileInput.value = 'test.svg';

        // Assign to ref
        if (result.current.fileInputRef.current) {
            Object.defineProperty(result.current.fileInputRef, 'current', {
                value: mockFileInput,
                writable: true,
            });
        }

        act(() => {
            result.current.clearTemplate();
        });

        expect(result.current.svgFile).toBeNull();
        expect(result.current.svgPreview).toBe('');
    });

    it('should handle non-SVG XML gracefully', async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const xmlContent = '<?xml version="1.0"?><root><item>data</item></root>';
        const file = new File([xmlContent], 'test.xml', { type: 'application/xml' });

        const fileInput = document.createElement('input');
        fileInput.type = 'file';
        Object.defineProperty(fileInput, 'files', {
            value: [file],
            writable: false,
        });

        const event = {
            target: fileInput,
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        await act(async () => {
            result.current.handleFileSelect(event);
        });

        expect(result.current.svgFile).toBe(file);
    });
});