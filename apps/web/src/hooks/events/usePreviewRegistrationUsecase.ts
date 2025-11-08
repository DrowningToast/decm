import { useState, useCallback, useMemo, useEffect } from "react";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { useEventPasswordNavStore } from "@/components/BottomNav/stores/event-password";
import { useEventInvitationNavStore } from "@/components/BottomNav/stores/event-invitation";
import { useCheckPasswordMutation } from "./useCheckPasswordMutation";
import { useEventViewModelUsecase } from "@/components/pages/Participant/Events/Detail/useEventViewModelUsecase";

/**
 * Mock PII data interface
 * This represents the participant information fields that might be required
 */
export interface ParticipantPIIData {
    firstName?: string;
    lastName?: string;
    email?: string;
    bio?: string;
    phoneNumber?: string;
    address?: string;
    academicInstitution?: string;
    academicEmail?: string;
}

/**
 * Configuration data returned from preview
 */
export interface RegistrationConfig {
    requireFirstName?: boolean;
    requireLastName?: boolean;
    requireEmail?: boolean;
    requireBio?: boolean;
    requirePhoneNumber?: boolean;
    requireAddress?: boolean;
    requireAcademicInstitution?: boolean;
    requireAcademicEmail?: boolean;
}

/**
 * Hook for submitting password-protected event registration
 *
 * Prerequisites:
 * - The event must be password-required
 * - The user must be authenticated
 *
 * Features:
 * - preview(): Checks password and displays requirements on success
 * - confirm(): Submits final registration with PII data
 * - Shows destructive toast on password error
 * - Shows alert modal with registration requirements on password success
 *
 * @param eventId - The ID of the event to register for
 */
export function usePreviewRegistrationUsecase(eventId: string) {
    const { t } = useTranslation();

    const { event } = useEventViewModelUsecase({ eventId });

    const [isPasswordValid, setIsPasswordValid] = useState(false);
    const [isAcceptingInvitation, setIsAcceptingInvitation] = useState(false);

    const showPreviewModal = useMemo(() => {
        return isPasswordValid || isAcceptingInvitation;
    }, [isPasswordValid, isAcceptingInvitation]);
    const closePreviewModal = useCallback(() => {
        setIsPasswordValid(false);
        setIsAcceptingInvitation(false);
    }, [setIsPasswordValid, setIsAcceptingInvitation]);

    const { mutateAsync: checkPassword } = useCheckPasswordMutation();

    const { setOnSubmitCallback } = useEventPasswordNavStore();
    const { setOnAcceptCallback } = useEventInvitationNavStore();

    const previewWithInvitation = useCallback(async () => {
        if (event?.isFull) {
            toast.error(t("event.registration.eventFull"), {
                description: t("event.registration.eventFullDescription"),
            });
            return;
        }
        if (event?.eventType !== "invite") {
            toast.error(t("event.registration.inviteOnly"), {
                description: t("event.registration.inviteOnlyDescription"),
            });
            return;
        }
        if (!event?.isInvited) {
            toast.error(t("event.registration.inviteOnly"), {
                description: t("event.registration.inviteOnlyDescription"),
            });
            return;
        }
        setIsAcceptingInvitation(true);
    }, [event?.eventType, event?.isFull, event?.isInvited, t]);

    /**
     * Preview function - checks the password with the backend
     */
    const previewWithPassword = useCallback(
        async (password: string = "") => {
            if (event?.isFull) {
                toast.error(t("event.registration.eventFull"), {
                    description: t("event.registration.eventFullDescription"),
                });
                return;
            }
            if (event?.eventType !== "private") {
                toast.error(t("event.registration.privateOnly"), {
                    description: t("event.registration.privateOnlyDescription"),
                });
                return;
            }

            const isPasswordValid = await checkPassword({ eventId, password });
            setIsPasswordValid(isPasswordValid);
            if (isPasswordValid) {
                toast.success(t("event.registration.success"), {
                    description: t("event.registration.successDescription"),
                });
            } else {
                toast.error(t("event.registration.incorrectPassword"), {
                    description: t("event.registration.incorrectPasswordDescription"),
                });
            }
        },
        [event?.isFull, event?.eventType, checkPassword, eventId, t],
    );

    // attach the callbacks to the stores
    useEffect(() => {
        setOnSubmitCallback((password: string) => {
            previewWithPassword(password);
        });
        setOnAcceptCallback(() => {
            previewWithInvitation();
        });
    }, [setOnSubmitCallback, previewWithPassword, setOnAcceptCallback, previewWithInvitation]);

    return {
        previewWithInvitation,
        previewWithPassword,
        isPasswordValid,
        showPreviewModal,
        closePreviewModal,
    };
}
