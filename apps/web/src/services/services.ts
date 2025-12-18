import { defaultAuthService } from "./AuthService/AuthService";
import { defaultEventRegistrationService } from "./EventRegistration/EventRegistration";
import { defaultEventService } from "./EventService/EventService";
import { defaultOnboardService } from "./OnboardService/OnboardService";
import { defaultIssuerService } from "./IssuerService/IssuerService";
import { defaultCertificateService } from "./CertificateService/CertificateService";
import { defaultSystemStatusService } from "./SystemStatusService/SystemStatusService";

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
export const eventRegistrationService = import.meta.env.VITE_USE_MOCK_API
    ? defaultEventRegistrationService
    : defaultEventRegistrationService;
export const issuerService = import.meta.env.VITE_USE_MOCK_API
    ? defaultIssuerService
    : defaultIssuerService;
export const certificateService = import.meta.env.VITE_USE_MOCK_API
    ? defaultCertificateService
    : defaultCertificateService;
export const systemStatusService = import.meta.env.VITE_USE_MOCK_API
    ? defaultSystemStatusService
    : defaultSystemStatusService;
// export const inboxService = import.meta.env.VITE_USE_MOCK_API ? defaultInboxService : defaultInboxService;
