import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { useLocalStorage } from './use-local-storage';

describe('useLocalStorage', () => {
    beforeEach(() => {
        localStorage.clear();
        vi.clearAllMocks();
    });

    it('should initialize with initial value when localStorage is empty', () => {
        const { result } = renderHook(() => useLocalStorage('test-key', 'initial-value'));

        expect(result.current[0]).toBe('initial-value');
    });

    it('should initialize with value from localStorage when it exists', () => {
        localStorage.setItem('test-key', JSON.stringify('stored-value'));

        const { result } = renderHook(() => useLocalStorage('test-key', 'initial-value'));

        expect(result.current[0]).toBe('stored-value');
    });

    it('should update localStorage when value changes', () => {
        const { result } = renderHook(() => useLocalStorage('test-key', 'initial-value'));

        act(() => {
            result.current[1]('new-value');
        });

        expect(result.current[0]).toBe('new-value');
        expect(localStorage.getItem('test-key')).toBe(JSON.stringify('new-value'));
    });

    it('should handle function updates', () => {
        const { result } = renderHook(() => useLocalStorage('test-key', 0));

        act(() => {
            result.current[1]((prev) => prev + 1);
        });

        expect(result.current[0]).toBe(1);
    });

    it('should remove value from localStorage when removeValue is called', () => {
        const { result } = renderHook(() => useLocalStorage('test-key', 'initial-value'));

        act(() => {
            result.current[1]('updated-value');
        });

        act(() => {
            result.current[2](); // removeValue
        });

        expect(result.current[0]).toBe('initial-value');
        expect(localStorage.getItem('test-key')).toBeNull();
    });

    it('should handle complex objects', () => {
        const initialValue = { name: 'Test', count: 0 };
        const { result } = renderHook(() => useLocalStorage('test-key', initialValue));

        act(() => {
            result.current[1]({ name: 'Updated', count: 1 });
        });

        expect(result.current[0]).toEqual({ name: 'Updated', count: 1 });
    });

    it('should handle custom serializer and deserializer', () => {
        const serializer = (value: number) => value.toString();
        const deserializer = (value: string) => parseInt(value, 10);

        const { result } = renderHook(() =>
            useLocalStorage('test-key', 0, { serializer, deserializer }),
        );

        act(() => {
            result.current[1](42);
        });

        expect(result.current[0]).toBe(42);
        expect(localStorage.getItem('test-key')).toBe('42');
    });

    it('should return initial value when JSON parsing fails', () => {
        localStorage.setItem('test-key', 'invalid-json');

        const { result } = renderHook(() => useLocalStorage('test-key', 'fallback-value'));

        expect(result.current[0]).toBe('fallback-value');
    });

    it('should handle undefined value', () => {
        localStorage.setItem('test-key', 'undefined');

        const { result } = renderHook(() => useLocalStorage<string | undefined>('test-key', 'default'));

        expect(result.current[0]).toBeUndefined();
    });

    it('should sync across multiple instances with same key', () => {
        const { result: result1 } = renderHook(() => useLocalStorage('shared-key', 'value1'));
        const { result: result2 } = renderHook(() => useLocalStorage('shared-key', 'value1'));

        act(() => {
            result1.current[1]('updated-value');
        });

        // Wait for storage event to propagate
        waitFor(() => {
            expect(result2.current[0]).toBe('updated-value');
        });
    });
});
