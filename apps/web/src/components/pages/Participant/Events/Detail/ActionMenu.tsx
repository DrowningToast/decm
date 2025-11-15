import React, { useMemo, useState } from "react";
import { BottomNav } from "@/components/BottomNav/BottomNav";
import { Typography } from "@/components/typography/typography";
import { useEventViewModelUsecase } from "./useEventViewModelUsecase";
import { usePreviewRegistrationUsecase } from "@/hooks/events/usePreviewRegistrationUsecase";
import { useTranslation } from "react-i18next";
import { useEventInvitationByUserAndEvent } from "./useEventInvitationByUserAndEvent";
import { useAuth } from "@/context/AuthContext";
import { RegistrationConfirmForm } from "./ConfirmForm";
import type { RegistrationConfirmDataForm } from "./RegistrationConfirmDataFormSchema";
import { toast } from "sonner";
import { useSignPasswordModalStore } from "@/components/providers/SignPasswordModal/store";
import { useEventContract } from "@/hooks/events/useEventContracts";
import { usePasswordPrompt } from "@/hooks/usePassowordPrompt";
import { eventRegistrationService, eventService } from "@/services/services";

interface ActionMenuProps {
    eventId: string;
}

export const ActionMenu: React.FC<ActionMenuProps> = ({ eventId }) => {
    const { t } = useTranslation();
    const { user } = useAuth();
    const [isSubmiting, setIsSubmiting] = useState(false);

    const { bottomNavVariant } = useEventViewModelUsecase({ eventId });
    const { showPreviewModal: _showPreviewModal, closePreviewModal } =
        usePreviewRegistrationUsecase(eventId);
    const { invitation, isLoading: isInvitationLoading } = useEventInvitationByUserAndEvent(
        eventId,
        user?.walletAddress,
    );
    const { event: eventViewModel } = useEventViewModelUsecase({ eventId });
    const { mutateAsync: openPasswordPrompt } = usePasswordPrompt();

    const { open: openConfirmSignModal } = useSignPasswordModalStore();

    const registrationInvitation = invitation?.registrationInvitation;

    // Use invitation data first, then profile data, then empty object
    const prefilledProfile = useMemo<RegistrationConfirmDataForm>(() => {
        return {
            firstName: registrationInvitation?.firstName ?? user?.firstName,
            lastName: registrationInvitation?.lastName ?? user?.lastName,
            bio: user?.bio,
            email: registrationInvitation?.email ?? user?.email,
            phoneNumber: registrationInvitation?.phoneNumber ?? user?.phoneNumber,
            address: user?.address,
            academicEmail: user?.academicEmail,
            academicInstitution: user?.academicInstitution,
        };
    }, [
        registrationInvitation?.email,
        registrationInvitation?.firstName,
        registrationInvitation?.lastName,
        registrationInvitation?.phoneNumber,
        user?.academicEmail,
        user?.academicInstitution,
        user?.address,
        user?.bio,
        user?.email,
        user?.firstName,
        user?.lastName,
        user?.phoneNumber,
    ]);

    const instructionText = useMemo(() => {
        switch (bottomNavVariant) {
            case "participating":
                return t("participant.events.detail.instruction.participating");
            case "invited":
                return t("participant.events.detail.instruction.invited");
            case "invitation-required":
                return t("participant.events.detail.instruction.invitationRequired");
            case "event-password":
                return t("participant.events.detail.instruction.passwordRequired");
        }
        return "";
    }, [bottomNavVariant, t]);

    const handleSubmit = async (data: RegistrationConfirmDataForm) => {
        try {
            setIsSubmiting(true);
            // password check
            const password = await openPasswordPrompt({
                eventContractAddress: eventViewModel?.eventContractAddress ?? "",
                transactionType: "Registration Confirmation",
                description: t("participant.events.detail.instruction.signatureRequest"),
                details: `Confirming registration for ${eventViewModel?.title}`,
            });

            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 1000));

            await eventRegistrationService.join;

            toast.success(t("events.registration.piiForm.submitSuccess"));
            closePreviewModal();
        } catch (error) {
            console.error("Failed to submit Registration Confirm data:", error);
            toast.error(t("events.registration.piiForm.submitError"));
        } finally {
            setIsSubmiting(false);
        }
    };

    const handleCancel = () => {
        closePreviewModal();
    };

    const showPreviewModal = useMemo(() => {
        if (openConfirmSignModal) {
            return;
        }
        return !isInvitationLoading && invitation !== undefined && _showPreviewModal;
    }, [_showPreviewModal, invitation, isInvitationLoading, openConfirmSignModal]);

    // Always show ActionMenu with conditional PII form
    return (
        <div className="relative">
            <div className="flex flex-col gap-y-4 md:sticky md:top-1/2">
                {/* Show PII form when showPreviewModal is true */}
                {showPreviewModal && (
                    <RegistrationConfirmForm
                        eventId={eventId}
                        onSubmit={handleSubmit}
                        onCancel={handleCancel}
                        isSubmitting={isSubmiting}
                        profileData={prefilledProfile}
                        invitationData={registrationInvitation}
                    />
                )}

                {/* Always show instruction text and BottomNav */}
                <div className="hidden md:inline-block w-full text-center">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-center text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {instructionText}
                    </Typography>
                </div>
                {bottomNavVariant && (
                    <BottomNav
                        className="md:min-w-[220px] md:w-full md:relative md:translate-y-0 md:top-0 md:bottom-0 md:translate-x-0 md:left-0 z-60 flex-col items-stretch"
                        variant={bottomNavVariant}
                        onBack={() => window.history.back()}
                    />
                )}
            </div>
        </div>
    );
};
