import { defaultAuthService } from "./AuthService";
import { defaultEventService } from "./EventService";
import { defaultOnboardService } from "./OnboardService";

// TODO: Switch to use the mock service if the environment variable is set to true
// WHEN: Implementing acceptance tests
export const eventService = import.meta.env.VITE_USE_MOCK_API
    ? defaultEventService
    : defaultEventService;
export const authService = import.meta.env.VITE_USE_MOCK_API
    ? defaultAuthService
    : defaultAuthService;
export const onboardService = import.meta.env.VITE_USE_MOCK_API
    ? defaultOnboardService
    : defaultOnboardService;
// export const issuerService = import.meta.env.VITE_USE_MOCK_API ? defaultIssuerService : defaultIssuerService;
// export const inboxService = import.meta.env.VITE_USE_MOCK_API ? defaultInboxService : defaultInboxService;
