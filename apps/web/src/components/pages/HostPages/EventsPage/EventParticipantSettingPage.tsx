import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import { ParticipantSettingsForm } from "@/components/forms/ParticipantSettingsForm";
import {
    type ParticipantSettingsData,
    defaultParticipantSettings,
} from "@/lib/schemas/participantSettingsSchema";
import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";

export const EventParticipantSettingPage = () => {
    const { t } = useTranslation();
    const { eventId } = useParams<{ eventId: string }>();
    const [isLoading, setIsLoading] = useState(false);

    const onSubmit = async (data: ParticipantSettingsData) => {
        try {
            setIsLoading(true);

            // TODO: Implement API call to save participant settings
            console.log("Event ID:", eventId);
            console.log("Participant Settings:", data);

            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 1000));

            // Success feedback (can be replaced with toast notification)
            alert(t("participantSettings.saveSuccess"));
        } catch (error) {
            console.error("Error saving participant settings:", error);
            alert(t("participantSettings.saveError"));
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <PageContainer
            title={t("participantSettings.pageTitle")}
            description={t("participantSettings.pageDescription")}
        >
            <SectionContainer>
                <ParticipantSettingsForm
                    onSubmit={onSubmit}
                    isLoading={isLoading}
                    defaultValues={defaultParticipantSettings}
                    showPreview={true}
                />
            </SectionContainer>
        </PageContainer>
    );
};
