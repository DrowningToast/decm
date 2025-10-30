import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach, vi } from "vitest";

// Cleanup after each test case (e.g., clearing jsdom)
afterEach(() => {
    cleanup();
});

// Mock environment variables
vi.stubEnv("VITE_CORE_BACKEND_API", "http://localhost:8080/api/v1");
vi.stubEnv("VITE_WALLETCONNECT_PROJECT_ID", "test-walletconnect-project-id");
vi.stubEnv("VITE_ENVIRONMENT", "test");

// Mock window.matchMedia
Object.defineProperty(window, "matchMedia", {
    writable: true,
    value: vi.fn().mockImplementation((query) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(), // deprecated
        removeListener: vi.fn(), // deprecated
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
    })),
});

// Mock IntersectionObserver
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
    root: null,
    rootMargin: "",
    thresholds: [],
    takeRecords: vi.fn(),
}));

// Mock ResizeObserver
global.ResizeObserver = vi.fn().mockImplementation(() => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
}));

// Prefer setting the base URL in vitest.config.ts via test.environmentOptions.happyDOM.url
// If you must stub, redefine via defineProperty with a minimal Location-like object:
// Object.defineProperty(window, 'location', {
//   configurable: true,
//   value: { ...window.location, href: 'http://localhost:3000', assign: vi.fn(), replace: vi.fn(), reload: vi.fn() }
// })
