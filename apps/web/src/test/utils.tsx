import React from "react";
import { render, type RenderOptions } from "@testing-library/react";
import { QueryClient, QueryClientProvider, setLogger } from "@tanstack/react-query";
import { MemoryRouter, type MemoryRouterProps } from "react-router-dom";
import { I18nextProvider } from "react-i18next";
import i18n from "i18next";
import { vi } from "vitest";

// Configure query client logger to silence logs in tests
setLogger({
    log: () => {},
    warn: () => {},
    error: () => {},
});

// Initialize i18n for tests
i18n.init({
    lng: "en",
    fallbackLng: "en",
    ns: ["common"],
    defaultNS: "common",
    resources: {
        en: {
            common: {
                welcome: "Welcome",
                signin: "Sign In",
                signup: "Sign Up",
                signout: "Sign Out",
                continue: "Continue",
                email: "Email",
                firstName: "First Name",
                lastName: "Last Name",
                phoneNumber: "Phone Number",
                address: "Address",
                bio: "Bio",
                academicEmail: "Academic Email",
                academicInstitution: "Academic Institution",
                profilePicture: "Profile Picture",
            },
            onboard: {
                api: {
                    error: {
                        authentication_credential_id_already_exists: "Account already exists",
                    },
                },
            },
            errors: {
                generic: "An error occurred",
            },
        },
    },
    interpolation: {
        escapeValue: false,
    },
});

// Create a custom render function with all providers
interface CustomRenderOptions extends Omit<RenderOptions, "wrapper"> {
    routerProps?: MemoryRouterProps;
    queryClient?: QueryClient;
}

// eslint-disable-next-line react-refresh/only-export-components
export function renderWithProviders(ui: React.ReactElement, options?: CustomRenderOptions) {
    const { routerProps, queryClient, ...renderOptions } = options ?? {};

    // Create a new QueryClient for each test to ensure isolation
    const testQueryClient =
        queryClient ??
        new QueryClient({
            defaultOptions: {
                queries: {
                    retry: false, // Disable retries in tests
                    gcTime: Infinity, // Disable garbage collection in tests
                },
                mutations: {
                    retry: false,
                },
            },
        });

    function Wrapper({ children }: { children: React.ReactNode }) {
        return (
            <QueryClientProvider client={testQueryClient}>
                <I18nextProvider i18n={i18n}>
                    <MemoryRouter {...routerProps}>{children}</MemoryRouter>
                </I18nextProvider>
            </QueryClientProvider>
        );
    }

    const renderResult = render(ui, { wrapper: Wrapper, ...renderOptions });

    return {
        ...renderResult,
        queryClient: testQueryClient,
    };
}

// Helper to create mock onboard status
const createMockOnboardStatus = (overrides = {}) => ({
    authentication_credential_id: "mock-credential-id",
    profile_id: null,
    message: "Onboarding status checked",
    ...overrides,
});

// Test ID helpers
const testIds = {
    // Auth components
    signInButton: "sign-in-button",
    signUpButton: "sign-up-button",
    signOutButton: "sign-out-button",
    googleSignInButton: "google-sign-in-button",
    walletSignInButton: "wallet-sign-in-button",

    // Form fields
    emailInput: "email-input",
    firstNameInput: "first-name-input",
    lastNameInput: "last-name-input",
    phoneNumberInput: "phone-number-input",
    addressInput: "address-input",
    bioTextarea: "bio-textarea",
    academicEmailInput: "academic-email-input",
    academicInstitutionInput: "academic-institution-input",

    // Onboarding
    onboardingForm: "onboarding-form",
    createAccountButton: "create-account-button",
    createProfileButton: "create-profile-button",
    skipButton: "skip-button",

    // Loading states
    loadingSpinner: "loading-spinner",
    errorMessage: "error-message",
    successMessage: "success-message",
};

// Wait for async updates
const waitForAsync = (ms = 0) => new Promise((resolve) => setTimeout(resolve, ms));

// Mock API response helper
function mockApiResponse<T>(data: T, delay = 0): Promise<T> {
    return new Promise((resolve) => {
        setTimeout(() => resolve(data), delay);
    });
}

// Mock API error helper
function mockApiError(message: string, delay = 0): Promise<never> {
    return new Promise((_, reject) => {
        setTimeout(() => reject(new Error(message)), delay);
    });
}

// Mock localStorage helper
function createLocalStorageMock() {
    let store: Record<string, string> = {};

    return {
        getItem: vi.fn((key: string) => store[key] ?? null),
        setItem: vi.fn((key: string, value: string) => {
            store[key] = value;
        }),
        removeItem: vi.fn((key: string) => {
            delete store[key];
        }),
        clear: vi.fn(() => {
            store = {};
        }),
        get length() {
            return Object.keys(store).length;
        },
        key: vi.fn((index: number) => {
            const keys = Object.keys(store);
            return keys[index] ?? null;
        }),
        // Helper to get current store state (for testing)
        getStore: () => ({ ...store }),
    };
}

// Helper to create mock form data
const mockFormData = {
    profile: {
        email: "test@example.com",
        firstName: "Test",
        lastName: "User",
        phoneNumber: "+1234567890",
        address: "123 Test St",
        bio: "This is a test bio",
        academicEmail: "test@university.edu",
        academicInstitution: "Test University",
        profilePictureUrl: "https://example.com/avatar.jpg",
    },
    googleAuth: {
        accessToken: "mock-google-access-token",
        password: "SecurePassword123!",
        expiresIn: 3600,
    },
    walletAuth: {
        signMessage: "0xmocksignedmessage",
    },
};

// Re-export testing library utilities
// eslint-disable-next-line react-refresh/only-export-components
export { render, type RenderOptions } from "@testing-library/react";
// eslint-disable-next-line react-refresh/only-export-components
export { cleanup } from "@testing-library/react";
// eslint-disable-next-line react-refresh/only-export-components
export { default as userEvent } from "@testing-library/user-event";

// Export test utilities
/* eslint-disable react-refresh/only-export-components */
export {
    createMockOnboardStatus,
    testIds,
    waitForAsync,
    mockApiResponse,
    mockApiError,
    createLocalStorageMock,
    mockFormData,
    renderWithProviders,
};
/* eslint-enable react-refresh/only-export-components */
