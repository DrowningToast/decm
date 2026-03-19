import { NotificationIndicator } from "./components/NotificationIndicator";

export interface ServiceConfig {
    key: string;
    href: string;
    translationKey: string;
    children?: React.ReactNode;
}

/**
 * Participant services configuration
 * Services available to authenticated participants
 */
export const participantServices: ServiceConfig[] = [
    {
        key: "certificates",
        href: "/app/certificates",
        translationKey: "participant.home.services.certificates",
    },
    {
        key: "events",
        href: "/app/events",
        translationKey: "participant.home.services.events",
    },
    {
        key: "inbox",
        href: "/app/inbox",
        translationKey: "participant.home.services.inbox",
        children: <NotificationIndicator />,
    },
] as const;

/**
 * General services configuration
 * Public services available to all users
 */
export const generalServices: ServiceConfig[] = [
    {
        key: "verify-certificates",
        href: "/verify",
        translationKey: "participant.home.services.verifyCertificates",
    },
] as const;
