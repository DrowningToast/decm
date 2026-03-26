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
import { useSignPasswordModalStore } from "@/components/providers/SignPasswordModal/store";
import { usePasswordPrompt } from "@/hooks/usePassowordPrompt";
import { useJoinEventMutation } from "@/hooks/events/useJoinEventMutation";
import { walletSigningService, eventRegistrationService } from "@/services/services";
import {
    WalletNotConnectedError,
    WalletSigningRejectedError,
} from "@/services/WalletSigningService/WalletSigningService";
import { toast } from "sonner";

interface ActionMenuProps {
    eventId: string;
}

export const ActionMenu: React.FC<ActionMenuProps> = ({ eventId }) => {
    const { t, i18n } = useTranslation();
    const { user } = useAuth();
    const { joinWithAccountPassword, joinWithSignature } = useJoinEventMutation();
    const [isSubmitting, setIsSubmitting] = useState(false);

    const { bottomNavVariant, event: eventViewModel } = useEventViewModelUsecase({ eventId });
    const { showPreviewModal: _showPreviewModal, closePreviewModal } =
        usePreviewRegistrationUsecase(eventId);
    const { invitation, isLoading: isInvitationLoading } = useEventInvitationByUserAndEvent(
        eventId,
        user?.walletAddress,
    );

    const { mutateAsync: openPasswordPrompt } = usePasswordPrompt();

    const { isOpen: isPinModalOpen } = useSignPasswordModalStore();

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

    // Format dates for display
    const locale = i18n.language === "th" ? "th-TH" : "en-US";
    const dateFormatOptions: Intl.DateTimeFormatOptions = {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    };

    const formattedDeadline = eventViewModel?.broadcastedAtEstimatedDeadline
        ? new Date(eventViewModel.broadcastedAtEstimatedDeadline).toLocaleDateString(
              locale,
              dateFormatOptions,
          )
        : undefined;

    const formattedBroadcastedAt = eventViewModel?.broadcastedAt
        ? new Date(eventViewModel.broadcastedAt).toLocaleDateString(locale, dateFormatOptions)
        : undefined;

    // Determine broadcast status for joined users
    const broadcastStatus = useMemo(() => {
        if (!eventViewModel?.isJoined) return null;

        if (eventViewModel.broadcastedAt) {
            return "confirmed";
        }
        if (eventViewModel.signatureExpiredAt) {
            return "expired";
        }
        if (eventViewModel.broadcastedAtEstimatedDeadline) {
            const deadline = new Date(eventViewModel.broadcastedAtEstimatedDeadline);
            if (deadline < new Date()) {
                return "expired";
            }
            return "pending";
        }
        return null;
    }, [
        eventViewModel?.isJoined,
        eventViewModel?.broadcastedAt,
        eventViewModel?.signatureExpiredAt,
        eventViewModel?.broadcastedAtEstimatedDeadline,
    ]);

    // Pre-compute broadcast status display conditions
    const showBroadcastConfirmed = broadcastStatus === "confirmed" && !!formattedBroadcastedAt;
    const showBroadcastPending = broadcastStatus === "pending" && !!formattedDeadline;
    const showBroadcastExpired = broadcastStatus === "expired";
    const broadcastPendingText = showBroadcastPending
        ? t("participant.events.detail.broadcast.pending", { deadline: formattedDeadline })
        : "";
    const broadcastConfirmedText = showBroadcastConfirmed
        ? t("participant.events.detail.broadcast.confirmed", { date: formattedBroadcastedAt })
        : "";
    const broadcastExpiredText = showBroadcastExpired
        ? t("participant.events.detail.broadcast.expired")
        : "";

    const handleSubmit = async (data: RegistrationConfirmDataForm) => {
        setIsSubmitting(true);
        try {
            if (!registrationInvitation) {
                throw new Error("No registration invitation found");
            }

            if (user?.solutionStatus === "BYOK") {
                // BYOK flow: sign with wallet
                const signMessage = await eventRegistrationService.getJoinEventSignMessage(eventId);
                let signature: string;
                try {
                    signature = await walletSigningService.signMessage(signMessage);
                } catch (error) {
                    if (error instanceof WalletNotConnectedError) {
                        toast.error(t("errors.walletNotConnected", "Wallet is not connected"));
                    } else if (error instanceof WalletSigningRejectedError) {
                        toast.error(t("errors.walletSigningRejected", "Signing request rejected"));
                    } else {
                        toast.error(t("errors.generic"));
                    }
                    return;
                }
                await joinWithSignature.mutateAsync({
                    eventId,
                    originalSignMessage: signMessage,
                    signature,
                    registrationData: data,
                });
            } else {
                // SYSTEM_MANAGED flow: use account password
                const checkedPassword = await openPasswordPrompt({
                    eventContractAddress: eventViewModel?.eventContractAddress ?? "",
                    transactionType: "Registration Confirmation",
                    title: t("participant.events.detail.instruction.passwordRequired"),
                    description: t("participant.events.detail.instruction.signatureRequest"),
                    details: `Confirming registration for ${eventViewModel?.title}`,
                });
                await joinWithAccountPassword.mutateAsync({
                    eventId,
                    accountPassword: checkedPassword,
                    registrationData: data,
                });
            }

            // Close modal after successful join
            closePreviewModal();
        } catch (error) {
            console.error("Failed to submit Registration Confirm data:", error);
            // Error toast is handled by the mutation hook
            // Don't close modal on error so user can retry
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleCancel = () => {
        closePreviewModal();
    };

    const showPreviewModal = useMemo(() => {
        if (isPinModalOpen) {
            return false;
        }
        return !isInvitationLoading && invitation !== undefined && _showPreviewModal;
    }, [_showPreviewModal, invitation, isInvitationLoading, isPinModalOpen]);

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
                        isSubmitting={isSubmitting}
                        profileData={prefilledProfile}
                        invitationData={registrationInvitation}
                    />
                )}

                {/* Always show instruction text and BottomNav */}
                <div className="flex flex-col gap-y-2">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-center text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {instructionText}
                    </Typography>

                    {/* Broadcast status for joined users */}
                    {showBroadcastConfirmed && (
                        <Typography
                            variant="text"
                            tag="p"
                            color="primary"
                            className="text-center text-sm font-medium [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {broadcastConfirmedText}
                        </Typography>
                    )}
                    {showBroadcastPending && (
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-center text-xs [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {broadcastPendingText}
                        </Typography>
                    )}
                    {showBroadcastExpired && (
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-center text-xs text-destructive [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {broadcastExpiredText}
                        </Typography>
                    )}
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
