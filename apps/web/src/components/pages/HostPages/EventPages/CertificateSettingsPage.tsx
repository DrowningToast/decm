import React from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Search } from "lucide-react";
import { IssuerSelectionModal } from "@/components/IssuerSelectionModal";
import { SelectedIssuersTable } from "@/components/SelectedIssuersTable";
import { CertificateTemplateUpload } from "@/components/CertificateTemplateUpload";
import { CertificatePreview } from "@/components/CertificatePreview";
import {
    CertificateFontSettings,
    type CertificateFontConfig,
} from "@/components/CertificateFontSettings";
import { convertProfileToIssuer, useIssuerManagement } from "@/hooks/useIssuerManagement";
import { useCertificateTemplate, type DetectedKeyword } from "@/hooks/useCertificateTemplate";
import SectionContainer from "@/components/container/SectionContainer";
import { useUpdateCertificateConfig } from "./useUpdateCertificateConfig";
import { useUpdateCertificateTextConfig } from "./useUpdateCertificateTextConfig";
import type {
    CoreApiInternalHandlerEventconfigEventCertificateConfigResponse,
    UpdateEventCertificateConfigPayload,
} from "@decm/api";
import { toast } from "sonner";
import { useNavigate } from "@/router";
import { useUpdateEventIssuer } from "./useUpdateEventIssuer";
import { useDeleteEventIssuer } from "./useDeleteEventIssuer";
import type { EventIssuer } from "@/services/EventService/EventService";
import type { Profile } from "@/services/AuthService/AuthService";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";

interface CertificateSettingsPageProps {
    eventId: string;
    eventCertificateConfig?: CoreApiInternalHandlerEventconfigEventCertificateConfigResponse;
    eventIssuers?: EventIssuer[];
}

export const CertificateSettingsPage = ({
    eventId,
    eventCertificateConfig,
    eventIssuers,
}: CertificateSettingsPageProps) => {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const { updateCertificateConfig, isUpdatingCertificateConfig } = useUpdateCertificateConfig(
        eventId!,
    );

    const { updateCertificateTextConfig, isUpdatingCertificateTextConfig } =
        useUpdateCertificateTextConfig(eventId!);

    const { updateEventIssuer, isUpdatingEventIssuer } = useUpdateEventIssuer(eventId!);
    const { deleteEventIssuerAsync } = useDeleteEventIssuer();

    // Check if certificate is published - if so, disable issuer editing
    const isCertificatePublished = eventCertificateConfig?.is_published ?? false;

    // Track font settings changes
    const [fontSettings, setFontSettings] = React.useState<CertificateFontConfig | null>(null);

    // Use custom hooks for state management
    // Extract issuer profiles from event issuers
    const selectedIssuerProfiles = eventIssuers?.map((issuer) => issuer);

    const issuerManagement = useIssuerManagement({
        selectedIssuers: selectedIssuerProfiles
            ?.map((issuer) =>
                issuer.issuerProfile ? convertProfileToIssuer(issuer.issuerProfile) : undefined,
            )
            .filter((v) => !!v),
    });

    // Get the list of selected issuer credential IDs (both saved and unsaved)
    const selectedIssuerCredentialIds = new Set(
        issuerManagement.selectedIssuers.map((issuer) => issuer.id),
    );

    // Saved issuers: directly from eventIssuers to preserve all original data
    const savedIssuers =
        eventIssuers?.filter((eventIssuer) =>
            selectedIssuerCredentialIds.has(eventIssuer.issuerCredentialId),
        ) || [];

    // Unsaved issuers: those in selectedIssuers but not in eventIssuers
    const unsavedIssuers: EventIssuer[] = issuerManagement.selectedIssuers
        .filter((issuer) => {
            // Check if this issuer is NOT in the saved issuers
            return !eventIssuers?.some(
                (eventIssuer) => eventIssuer.issuerCredentialId === issuer.id,
            );
        })
        .map((issuer) => {
            // Create a new EventIssuer-like structure for unsaved issuers
            // Convert Issuer back to Profile format for issuerProfile
            const nameParts = issuer.name.trim().split(/\s+/);
            const firstName = nameParts[0] || undefined;
            const lastName = nameParts.slice(1).join(" ") || undefined;

            const profile: Profile = {
                id: issuer.id,
                authenticationCredentialId: issuer.id,
                firstName,
                lastName,
                email: issuer.email || undefined,
                phoneNumber: issuer.phoneNumber,
                academicInstitution: issuer.organization,
                googleConnectorRef: issuer.googleOAuthEmail,
                isEmailPublic: true,
                isFirstNamePublic: true,
                isLastNamePublic: true,
                isPhoneNumberPublic: true,
                isAcademicInstitutionPublic: true,
                isProfilePicturePublic: false,
                isBioPublic: false,
                isAddressPublic: false,
                isAcademicEmailPublic: false,
            };

            return {
                eventId: eventId!,
                id: `temp-${issuer.id}`, // Temporary ID for unsaved issuers
                issuerCredentialId: issuer.id,
                isSigned: false,
                issuerProfile: profile,
                createdAt: new Date(),
                updatedAt: new Date(),
            };
        });
    const certificateTemplate = useCertificateTemplate();

    // Reconstruct detected keywords from saved config when template already exists
    // Only include keywords that have meaningful positions (not 0,0) since required fields
    // always have values even if the keyword wasn't detected
    const savedDetectedKeywords = React.useMemo(() => {
        if (!eventCertificateConfig) return [];

        const keywords: DetectedKeyword[] = [];

        // Add name keyword if position exists and is not (0,0)
        if (
            eventCertificateConfig.name_pos_x !== undefined &&
            eventCertificateConfig.name_pos_y !== undefined &&
            (eventCertificateConfig.name_pos_x !== 0 || eventCertificateConfig.name_pos_y !== 0)
        ) {
            keywords.push({
                keyword: "{{ name }}",
                x: eventCertificateConfig.name_pos_x,
                y: eventCertificateConfig.name_pos_y,
                count: 1,
            });
        }

        // Add eventName keyword if position exists and is not (0,0)
        if (
            eventCertificateConfig.event_name_pos_x !== undefined &&
            eventCertificateConfig.event_name_pos_y !== undefined &&
            (eventCertificateConfig.event_name_pos_x !== 0 ||
                eventCertificateConfig.event_name_pos_y !== 0)
        ) {
            keywords.push({
                keyword: "{{ eventName }}",
                x: eventCertificateConfig.event_name_pos_x,
                y: eventCertificateConfig.event_name_pos_y,
                count: 1,
            });
        }

        // Add academicInstitutionName keyword if position exists (it's optional, so if it exists, it was detected)
        if (
            eventCertificateConfig.academic_institution_pos_x !== undefined &&
            eventCertificateConfig.academic_institution_pos_y !== undefined
        ) {
            keywords.push({
                keyword: "{{ academicInstitutionName }}",
                x: eventCertificateConfig.academic_institution_pos_x,
                y: eventCertificateConfig.academic_institution_pos_y,
                count: 1,
            });
        }

        // Add certificateTitle keyword if position exists (it's optional, so if it exists, it was detected)
        if (
            eventCertificateConfig.certificate_title_pos_x !== undefined &&
            eventCertificateConfig.certificate_title_pos_y !== undefined
        ) {
            keywords.push({
                keyword: "{{ certificateTitle }}",
                x: eventCertificateConfig.certificate_title_pos_x,
                y: eventCertificateConfig.certificate_title_pos_y,
                count: 1,
            });
        }

        // Add certificateSubtitle keyword if position exists (it's optional, so if it exists, it was detected)
        if (
            eventCertificateConfig.certificate_subtitle_pos_x !== undefined &&
            eventCertificateConfig.certificate_subtitle_pos_y !== undefined
        ) {
            keywords.push({
                keyword: "{{ certificateSubtitle }}",
                x: eventCertificateConfig.certificate_subtitle_pos_x,
                y: eventCertificateConfig.certificate_subtitle_pos_y,
                count: 1,
            });
        }

        return keywords;
    }, [eventCertificateConfig]);

    // Merge saved keywords with newly detected keywords (new upload takes precedence)
    const allDetectedKeywords = React.useMemo(() => {
        // If there's a new upload, use those keywords
        if (certificateTemplate.detectedKeywords.length > 0) {
            return certificateTemplate.detectedKeywords;
        }
        // Otherwise, use saved keywords
        return savedDetectedKeywords;
    }, [certificateTemplate.detectedKeywords, savedDetectedKeywords]);

    // Handle form submission
    const handleSubmit = async () => {
        try {
            // Validate that at least one issuer is selected
            if (issuerManagement.selectedIssuers.length === 0) {
                toast.error(t("certificateSettings.noIssuersError"));
                return;
            }

            // TODO: Implement API call to save certificate settings
            console.log("Event ID:", eventId);
            console.log("Selected Issuers:", issuerManagement.selectedIssuers);
            console.log("SVG File:", certificateTemplate.svgFile);
            console.log("Detected Keywords:", allDetectedKeywords);

            const name = allDetectedKeywords.find((keyword) => keyword.keyword === "{{ name }}");

            const eventName = allDetectedKeywords.find(
                (keyword) => keyword.keyword === "{{ eventName }}",
            );

            const acedmicInstitutionName = allDetectedKeywords.find(
                (keyword) => keyword.keyword === "{{ academicInstitutionName }}",
            );

            const certificateTitle = allDetectedKeywords.find(
                (keyword) => keyword.keyword === "{{ certificateTitle }}",
            );

            const certificateSubtitle = allDetectedKeywords.find(
                (keyword) => keyword.keyword === "{{ certificateSubtitle }}",
            );

            if (certificateTemplate.svgFile && !name) {
                toast.error(t("certificateSettings.nameNotFound"));
                return;
            }

            const req: UpdateEventCertificateConfigPayload = {
                name_pos_x: name?.x ?? eventCertificateConfig?.name_pos_x ?? 0,
                name_pos_y: name?.y ?? eventCertificateConfig?.name_pos_y ?? 0,
                event_name_pos_x: eventName?.x ?? eventCertificateConfig?.event_name_pos_x ?? 0,
                event_name_pos_y: eventName?.y ?? eventCertificateConfig?.event_name_pos_y ?? 0,
                base_certificate_image: certificateTemplate.svgFile ?? undefined,
            };

            if (acedmicInstitutionName) {
                req.academic_institution_pos_x = acedmicInstitutionName.x;
                req.academic_institution_pos_y = acedmicInstitutionName.y;
            }

            if (certificateTitle) {
                req.certificate_title_pos_x = certificateTitle.x;
                req.certificate_title_pos_y = certificateTitle.y;
            }

            if (certificateSubtitle) {
                req.certificate_subtitle_pos_x = certificateSubtitle.x;
                req.certificate_subtitle_pos_y = certificateSubtitle.y;
            }

            await updateCertificateConfig(req);

            // Only update issuers if certificate is not published
            if (!isCertificatePublished) {
                await updateEventIssuer(
                    issuerManagement.selectedIssuers.map((issuer) => {
                        console.log(issuer);
                        return {
                            event_id: eventId,
                            issuer_credential_id: issuer.id,
                        };
                    }),
                );
            }

            // Update font settings if they have changed
            if (fontSettings && eventCertificateConfig) {
                await updateCertificateTextConfig(fontSettings);
            }

            // Refetch all queries to ensure the page is up to date
            await Promise.all([
                queryClient.refetchQueries({
                    queryKey: QUERY_KEY.event.certificate.config(eventId),
                }),
                queryClient.refetchQueries({
                    queryKey: QUERY_KEY.event.issuers.byEventId(eventId),
                }),
                queryClient.refetchQueries({
                    queryKey: QUERY_KEY.event.byId(eventId),
                }),
            ]);

            toast.success(t("certificateSettings.saveSuccess"));
        } catch (error) {
            console.error("Error saving certificate settings:", error);
            toast.error(t("certificateSettings.saveError"));
        }
    };

    const handleRemoveIssuer = async (issuerId: string) => {
        try {
            // Check if this is a temporary ID (unsaved issuer) or a real EventIssuer ID
            if (issuerId.startsWith("temp-")) {
                // For unsaved issuers, just remove from local state
                const actualIssuerId = issuerId.replace("temp-", "");
                issuerManagement.handleRemoveIssuer(actualIssuerId);
            } else {
                // For saved issuers, delete from backend and remove from local state
                await deleteEventIssuerAsync({ eventId, issuerId });
                // Find the issuerCredentialId from the EventIssuer
                const eventIssuer = eventIssuers?.find((ei) => ei.id === issuerId);
                if (eventIssuer) {
                    issuerManagement.handleRemoveIssuer(eventIssuer.issuerCredentialId);
                }
            }
            toast.success(t("certificateSettings.removeIssuerSuccess"));
        } catch (error) {
            console.error("Error removing issuer:", error);
            toast.error(t("certificateSettings.removeIssuerError"));
        }
    };

    const handleCancel = () => {
        // Reset form or navigate back
        issuerManagement.handleClearSelection();
        certificateTemplate.clearTemplate();
        navigate("/host/events/:eventId", {
            params: {
                eventId,
            },
        });
    };

    // Check if form is valid for submission
    const isSelectedIssuer = issuerManagement.selectedIssuers.length > 0;

    const isCreateFormValid =
        certificateTemplate.svgFile !== null &&
        certificateTemplate.detectedKeywords.length > 0 &&
        !certificateTemplate.hasMissingMandatory &&
        isSelectedIssuer;
    // For update, we still need at least one issuer
    const isUpdateFormValid = isSelectedIssuer;

    const isFormValid = !eventCertificateConfig ? isCreateFormValid : isUpdateFormValid;

    return (
        <div>
            <SectionContainer>
                <div className="space-y-8">
                    {/* Page Title and Description */}
                    <div className="space-y-2">
                        <Typography variant="header" tag="h1" className="text-2xl font-bold">
                            {t("certificateSettings.pageTitle")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            className="text-sm text-muted-foreground"
                        >
                            {t("certificateSettings.pageDescription")}
                        </Typography>
                    </div>

                    {/* Step 1: Issuer Settings */}
                    <div className="space-y-6">
                        <div>
                            <Typography
                                variant="header"
                                tag="h2"
                                className="text-xl font-bold mb-2"
                            >
                                {t("certificateSettings.step1.title")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-sm text-muted-foreground"
                            >
                                {t("certificateSettings.step1.description")}
                            </Typography>
                        </div>

                        {/* Search Issuers */}
                        <div className="space-y-4 rounded-lg border p-6">
                            {isCertificatePublished && (
                                <div className="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        className="text-sm text-yellow-800"
                                    >
                                        {t("certificateSettings.step1.publishedWarning")}
                                    </Typography>
                                </div>
                            )}
                            <div>
                                <Label htmlFor="issuer-search">
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-sm font-medium"
                                    >
                                        {t("certificateSettings.step1.searchLabel")}
                                    </Typography>
                                </Label>
                                <div className="flex flex-col gap-2">
                                    <div className="flex gap-2 mt-2">
                                        <Input
                                            id="issuer-search"
                                            type="text"
                                            placeholder={t(
                                                "certificateSettings.step1.searchPlaceholder",
                                            )}
                                            value={issuerManagement.searchQuery}
                                            onChange={(e) =>
                                                issuerManagement.handleSearchQueryChange(
                                                    e.target.value,
                                                )
                                            }
                                            onKeyDown={(e) => {
                                                if (e.key === "Enter") {
                                                    issuerManagement.handleSearch();
                                                }
                                            }}
                                            disabled={
                                                issuerManagement.isSearching ||
                                                isCertificatePublished
                                            }
                                            className="h-10"
                                        />
                                        <Button
                                            type="button"
                                            onClick={issuerManagement.handleSearch}
                                            disabled={
                                                issuerManagement.isSearching ||
                                                !issuerManagement.searchQuery.trim() ||
                                                isCertificatePublished
                                            }
                                            className="min-w-[100px] h-10"
                                        >
                                            <Search className="h-4 w-4 mr-2" />
                                            <Typography
                                                variant="text"
                                                tag="span"
                                                className="font-medium"
                                            >
                                                {issuerManagement.isSearching
                                                    ? t("common.searching")
                                                    : t("certificateSettings.step1.searchButton")}
                                            </Typography>
                                        </Button>
                                    </div>
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        className="text-xs text-muted-foreground"
                                    >
                                        {t("certificateSettings.step1.searchNote")}
                                    </Typography>
                                </div>
                            </div>

                            {/* Saved Issuers Table */}
                            <SelectedIssuersTable
                                selectedIssuers={savedIssuers}
                                onRemoveIssuer={handleRemoveIssuer}
                                title={t("certificateSettings.step1.savedIssuers")}
                                disabled={isCertificatePublished}
                            />

                            {/* Unsaved Issuers Table */}
                            <SelectedIssuersTable
                                selectedIssuers={unsavedIssuers}
                                onRemoveIssuer={handleRemoveIssuer}
                                title={t("certificateSettings.step1.unsavedIssuers")}
                                isUnsaved={true}
                                disabled={isCertificatePublished}
                            />
                        </div>
                    </div>

                    {/* Step 2: Certificate Template Settings */}
                    <div className="space-y-6">
                        <div>
                            <Typography
                                variant="header"
                                tag="h2"
                                className="text-xl font-bold mb-2"
                            >
                                {t("certificateSettings.step2.title")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-sm text-muted-foreground"
                            >
                                {t("certificateSettings.step2.description")}
                            </Typography>
                        </div>

                        {/* Certificate Template Upload Component */}
                        <CertificateTemplateUpload
                            svgFile={certificateTemplate.svgFile}
                            availableKeywords={certificateTemplate.availableKeywords}
                            onFileSelect={certificateTemplate.handleFileSelect}
                            fileInputRef={certificateTemplate.fileInputRef}
                        />

                        {/* Certificate Preview Component */}
                        <CertificatePreview
                            svgPreview={certificateTemplate.svgPreview}
                            imageUrl={
                                certificateTemplate.svgPreview
                                    ? undefined
                                    : eventCertificateConfig?.base_certificate_presigned_url
                            }
                            detectedKeywords={allDetectedKeywords}
                            availableKeywords={certificateTemplate.availableKeywords}
                        />
                    </div>

                    {/* Step 3: Font Settings (only show if certificate config exists) */}
                    {eventCertificateConfig && (
                        <div className="space-y-6">
                            <div>
                                <Typography
                                    variant="header"
                                    tag="h2"
                                    className="text-xl font-bold mb-2"
                                >
                                    {t("certificateSettings.step3.title", "Step 3: Font Settings")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-sm text-muted-foreground"
                                >
                                    {t(
                                        "certificateSettings.step3.description",
                                        "Customize the font family and weight for each text field in your certificate template.",
                                    )}
                                </Typography>
                            </div>

                            <CertificateFontSettings
                                eventCertificateConfig={eventCertificateConfig}
                                onChange={(fontConfig: CertificateFontConfig) => {
                                    setFontSettings(fontConfig);
                                }}
                            />
                        </div>
                    )}

                    {/* Action Buttons */}
                    <div className="flex justify-end gap-4 pt-4">
                        <Button
                            type="button"
                            variant="secondary-dark"
                            size="lg"
                            onClick={handleCancel}
                            disabled={isUpdatingCertificateConfig}
                            className="min-w-[150px]"
                        >
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("common.cancel")}
                            </Typography>
                        </Button>
                        <Button
                            type="button"
                            variant="primary"
                            size="lg"
                            onClick={handleSubmit}
                            disabled={
                                !isFormValid || isUpdatingCertificateConfig || isUpdatingEventIssuer
                            }
                            className="min-w-[150px]"
                        >
                            <Typography variant="text" tag="span" className="font-medium">
                                {isUpdatingCertificateConfig
                                    ? t("common.loading")
                                    : t("certificateSettings.confirmButton")}
                            </Typography>
                        </Button>
                    </div>
                </div>
            </SectionContainer>

            {/* Issuer Selection Modal */}
            <IssuerSelectionModal
                isOpen={issuerManagement.isModalOpen}
                onClose={issuerManagement.handleCloseModal}
                onConfirm={issuerManagement.handleConfirmSelection}
                searchResults={issuerManagement.searchResults}
                initialSelectedIds={issuerManagement.getSelectedIssuerIds()}
                isLoading={issuerManagement.isSearching}
                searchQuery={issuerManagement.searchQuery}
            />

            {/* Hidden SVG container for processing */}
            <div
                ref={certificateTemplate.svgTempRef}
                id="svg-temp-container"
                className="absolute top-0 left-0 w-[1920px] h-[1080px] overflow-hidden opacity-0 pointer-events-none"
            />
        </div>
    );
};
