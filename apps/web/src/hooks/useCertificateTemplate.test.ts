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
