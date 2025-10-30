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

describe('utils - Additional Edge Cases', () => {
    describe('cn - Advanced Scenarios', () => {
        it('should handle deeply nested arrays', () => {
            expect(cn(['foo', ['bar', ['baz']]])).toBe('foo bar baz');
        });

        it('should handle multiple Tailwind conflicts', () => {
            expect(cn('p-4 px-2', 'py-3 px-5')).toBe('p-4 py-3 px-5');
        });

        it('should handle responsive Tailwind classes', () => {
            const result = cn('sm:px-4 md:px-6', 'lg:px-8');
            expect(result).toContain('sm:px-4');
            expect(result).toContain('md:px-6');
            expect(result).toContain('lg:px-8');
        });

        it('should handle pseudo-class variants', () => {
            const result = cn('hover:bg-blue-500', 'focus:bg-blue-600');
            expect(result).toContain('hover:bg-blue-500');
            expect(result).toContain('focus:bg-blue-600');
        });

        it('should handle dark mode variants', () => {
            const result = cn('bg-white dark:bg-black', 'text-black dark:text-white');
            expect(result).toContain('bg-white');
            expect(result).toContain('dark:bg-black');
            expect(result).toContain('text-black');
            expect(result).toContain('dark:text-white');
        });

        it('should handle arbitrary values', () => {
            const result = cn('top-[117px]', 'w-[32rem]');
            expect(result).toContain('top-[117px]');
            expect(result).toContain('w-[32rem]');
        });

        it('should handle important modifier', () => {
            const result = cn('!text-red-500', 'text-blue-500');
            expect(result).toContain('!text-red-500');
        });

        it('should handle empty string inputs', () => {
            expect(cn('', 'foo', '')).toBe('foo');
        });

        it('should handle only falsy values', () => {
            expect(cn(false, null, undefined, '')).toBe('');
        });

        it('should handle mixed types correctly', () => {
            expect(cn('foo', 42, true && 'bar', { baz: true, qux: false })).toBe('foo 42 bar baz');
        });

        it('should handle spaces in class names', () => {
            expect(cn('foo  bar', 'baz')).toBe('foo bar baz');
        });

        it('should be composable', () => {
            const baseClasses = cn('font-sans text-base');
            const extendedClasses = cn(baseClasses, 'text-lg font-bold');
            expect(extendedClasses).toContain('font-sans');
            expect(extendedClasses).toContain('text-lg');
            expect(extendedClasses).toContain('font-bold');
        });

        it('should handle state-based conditional classes', () => {
            const isActive = true;
            const isDisabled = false;

            const result = cn(
                'base-class',
                isActive && 'active-class',
                isDisabled && 'disabled-class'
            );

            expect(result).toContain('base-class');
            expect(result).toContain('active-class');
            expect(result).not.toContain('disabled-class');
        });
    });

    describe('delay - Advanced Scenarios', () => {
        it('should handle very short delays accurately', async () => {
            const start = Date.now();
            await delay(1);
            const end = Date.now();

            expect(end - start).toBeLessThan(50);
        });

        it('should handle longer delays', async () => {
            const start = Date.now();
            await delay(200);
            const end = Date.now();

            expect(end - start).toBeGreaterThanOrEqual(190);
            expect(end - start).toBeLessThan(250);
        });

        it('should be chainable with other promises', async () => {
            const result = await delay(10).then(() => 'completed');
            expect(result).toBe('completed');
        });

        it('should work in Promise.all', async () => {
            const start = Date.now();
            await Promise.all([delay(50), delay(50), delay(50)]);
            const end = Date.now();

            // All should complete in parallel, taking roughly 50ms total
            expect(end - start).toBeLessThan(100);
        });

        it('should work in Promise.race', async () => {
            const start = Date.now();
            await Promise.race([delay(50), delay(100), delay(150)]);
            const end = Date.now();

            // Should complete after the shortest delay
            expect(end - start).toBeGreaterThanOrEqual(45);
            expect(end - start).toBeLessThan(80);
        });

        it('should handle negative delay as zero', async () => {
            const start = Date.now();
            await delay(-100);
            const end = Date.now();

            expect(end - start).toBeLessThan(10);
        });

        it('should be interruptible with AbortController', async () => {
            const controller = new AbortController();
            
            const delayPromise = Promise.race([
                delay(1000),
                new Promise((_, reject) => {
                    controller.signal.addEventListener('abort', () => reject(new Error('Aborted')));
                })
            ]);

            setTimeout(() => controller.abort(), 50);

            await expect(delayPromise).rejects.toThrow('Aborted');
        });

        it('should work with async/await in loops', async () => {
            const results: number[] = [];
            
            for (let i = 0; i < 3; i++) {
                await delay(10);
                results.push(i);
            }

            expect(results).toEqual([0, 1, 2]);
        });

        it('should maintain order in sequential calls', async () => {
            const order: number[] = [];

            await delay(20).then(() => order.push(1));
            await delay(10).then(() => order.push(2));
            await delay(5).then(() => order.push(3));

            expect(order).toEqual([1, 2, 3]);
        });
    });
});