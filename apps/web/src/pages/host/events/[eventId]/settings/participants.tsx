import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { Typography } from "@/components/typography/typography";
import { ParticipantSettingsForm } from "@/components/forms/ParticipantSettingsForm";
import { RegistrationFormPreview } from "@/components/forms/ParticipantSettingsForm";
import {
    type ParticipantSettingsData,
    defaultParticipantSettings,
} from "@/lib/schemas/participantSettingsSchema";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Settings, Eye } from "lucide-react";

export const ParticipantSettingsPage = () => {
    const { t } = useTranslation();
    const { eventId } = useParams<{ eventId: string }>();

    // State to hold form values for preview
    const [formValues, setFormValues] = useState<ParticipantSettingsData>(
        defaultParticipantSettings,
    );

    const handleSubmit = async (data: ParticipantSettingsData) => {
        try {
            // Update form values for preview
            setFormValues(data);

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
        }
    };

    return (
        <>
            <FaviconHelmet
                title={`${t("participantSettings.pageTitle")} | ${t("common.appName")}`}
                description={t("participantSettings.pageDescription")}
            />

            <div className="container mx-auto px-4 py-8 max-w-6xl">
                {/* Page Header */}
                <div className="mb-8">
                    <Typography variant="header" tag="h1" className="text-3xl font-bold mb-2">
                        {t("participantSettings.pageTitle")}
                    </Typography>
                    <Typography variant="text" tag="p" className="text-muted-foreground">
                        {t("participantSettings.pageDescription")}
                    </Typography>
                </div>

                {/* Tabs for Settings and Preview */}
                <Tabs defaultValue="settings" className="w-full">
                    <TabsList className="grid w-full max-w-md grid-cols-2 mb-8">
                        <TabsTrigger value="settings" className="flex items-center gap-2">
                            <Settings className="h-4 w-4" />
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("participantSettings.tabs.settings")}
                            </Typography>
                        </TabsTrigger>
                        <TabsTrigger value="preview" className="flex items-center gap-2">
                            <Eye className="h-4 w-4" />
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("participantSettings.tabs.preview")}
                            </Typography>
                        </TabsTrigger>
                    </TabsList>

                    {/* Settings Tab */}
                    <TabsContent value="settings">
                        <ParticipantSettingsForm
                            defaultValues={formValues}
                            onSubmit={handleSubmit}
                        />
                    </TabsContent>

                    {/* Preview Tab */}
                    <TabsContent value="preview">
                        <RegistrationFormPreview settings={formValues} />
                    </TabsContent>
                </Tabs>
            </div>
        </>
    );
};
