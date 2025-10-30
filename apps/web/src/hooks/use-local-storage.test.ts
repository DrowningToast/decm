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

describe('useLocalStorage - Additional Edge Cases', () => {
    beforeEach(() => {
        localStorage.clear();
        vi.clearAllMocks();
    });

    it('should handle localStorage quota exceeded error', () => {
        const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
        
        // Mock localStorage.setItem to throw quota exceeded error
        const originalSetItem = Storage.prototype.setItem;
        Storage.prototype.setItem = vi.fn().mockImplementation(() => {
            throw new DOMException('QuotaExceededError');
        });

        const { result } = renderHook(() => useLocalStorage('test-key', 'initial'));

        act(() => {
            result.current[1]('large-value');
        });

        expect(consoleSpy).toHaveBeenCalled();
        expect(result.current[0]).toBe('initial');

        // Restore original implementation
        Storage.prototype.setItem = originalSetItem;
        consoleSpy.mockRestore();
    });

    it('should handle null values correctly', () => {
        const { result } = renderHook(() => useLocalStorage<string | null>('test-key', null));

        expect(result.current[0]).toBeNull();

        act(() => {
            result.current[1]('value');
        });

        expect(result.current[0]).toBe('value');

        act(() => {
            result.current[1](null);
        });

        expect(result.current[0]).toBeNull();
    });

    it('should handle array values', () => {
        const initialArray = [1, 2, 3];
        const { result } = renderHook(() => useLocalStorage('test-key', initialArray));

        expect(result.current[0]).toEqual([1, 2, 3]);

        act(() => {
            result.current[1]([...result.current[0], 4]);
        });

        expect(result.current[0]).toEqual([1, 2, 3, 4]);
    });

    it('should handle boolean values', () => {
        const { result } = renderHook(() => useLocalStorage('test-key', false));

        expect(result.current[0]).toBe(false);

        act(() => {
            result.current[1](true);
        });

        expect(result.current[0]).toBe(true);
    });

    it('should handle number zero correctly', () => {
        const { result } = renderHook(() => useLocalStorage('test-key', 0));

        expect(result.current[0]).toBe(0);

        act(() => {
            result.current[1](42);
        });

        expect(result.current[0]).toBe(42);
    });

    it('should maintain type safety with TypeScript generics', () => {
        interface User {
            id: number;
            name: string;
        }

        const { result } = renderHook(() =>
            useLocalStorage<User>('user', { id: 1, name: 'John' }),
        );

        expect(result.current[0]).toEqual({ id: 1, name: 'John' });

        act(() => {
            result.current[1]({ id: 2, name: 'Jane' });
        });

        expect(result.current[0]).toEqual({ id: 2, name: 'Jane' });
    });

    it('should handle custom serializer errors gracefully', () => {
        const badSerializer = () => {
            throw new Error('Serialization failed');
        };

        const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});

        const { result } = renderHook(() =>
            useLocalStorage('test-key', 'initial', { serializer: badSerializer }),
        );

        act(() => {
            result.current[1]('new-value');
        });

        expect(consoleSpy).toHaveBeenCalled();
        consoleSpy.mockRestore();
    });

    it('should handle custom deserializer errors gracefully', () => {
        localStorage.setItem('test-key', JSON.stringify('stored-value'));

        const badDeserializer = () => {
            throw new Error('Deserialization failed');
        };

        const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

        const { result } = renderHook(() =>
            useLocalStorage('test-key', 'fallback', { deserializer: badDeserializer }),
        );

        expect(result.current[0]).toBe('fallback');
        consoleSpy.mockRestore();
    });

    it('should handle localStorage.getItem returning null', () => {
        const { result } = renderHook(() => useLocalStorage('non-existent-key', 'default'));

        expect(result.current[0]).toBe('default');
    });

    it('should update when key changes', () => {
        localStorage.setItem('key1', JSON.stringify('value1'));
        localStorage.setItem('key2', JSON.stringify('value2'));

        const { result, rerender } = renderHook(
            ({ key }) => useLocalStorage(key, 'default'),
            { initialProps: { key: 'key1' } },
        );

        expect(result.current[0]).toBe('value1');

        rerender({ key: 'key2' });

        waitFor(() => {
            expect(result.current[0]).toBe('value2');
        });
    });

    it('should handle empty string as initial value', () => {
        const { result } = renderHook(() => useLocalStorage('test-key', ''));

        expect(result.current[0]).toBe('');

        act(() => {
            result.current[1]('non-empty');
        });

        expect(result.current[0]).toBe('non-empty');
    });

    it('should handle deeply nested objects', () => {
        const deepObject = {
            level1: {
                level2: {
                    level3: {
                        value: 'deep',
                    },
                },
            },
        };

        const { result } = renderHook(() => useLocalStorage('test-key', deepObject));

        expect(result.current[0]).toEqual(deepObject);

        act(() => {
            result.current[1]({
                level1: {
                    level2: {
                        level3: {
                            value: 'updated',
                        },
                    },
                },
            });
        });

        expect(result.current[0].level1.level2.level3.value).toBe('updated');
    });
});