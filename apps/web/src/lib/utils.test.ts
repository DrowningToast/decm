import { describe, it, expect } from 'vitest';
import { cn, delay } from './utils';

describe('utils', () => {
    describe('cn', () => {
        it('should merge class names correctly', () => {
            expect(cn('foo', 'bar')).toBe('foo bar');
        });

        it('should handle conditional classes', () => {
            expect(cn('foo', false && 'bar', 'baz')).toBe('foo baz');
        });

        it('should merge Tailwind classes correctly', () => {
            expect(cn('px-2 py-1', 'px-4')).toBe('py-1 px-4');
        });

        it('should handle empty arrays', () => {
            expect(cn('foo', [])).toBe('foo');
        });

        it('should handle undefined and null', () => {
            expect(cn('foo', undefined, null)).toBe('foo');
        });

        it('should handle objects', () => {
            expect(cn({ foo: true, bar: false })).toBe('foo');
        });

        it('should handle mixed inputs', () => {
            expect(cn('foo', { bar: true }, 'baz', false && 'qux')).toBe('foo bar baz');
        });
    });

    describe('delay', () => {
        it('should resolve after specified milliseconds', async () => {
            const start = Date.now();
            await delay(100);
            const end = Date.now();

            expect(end - start).toBeGreaterThanOrEqual(90); // Allow some margin
        });

        it('should handle zero delay', async () => {
            const start = Date.now();
            await delay(0);
            const end = Date.now();

            expect(end - start).toBeLessThan(10);
        });
    });
});
