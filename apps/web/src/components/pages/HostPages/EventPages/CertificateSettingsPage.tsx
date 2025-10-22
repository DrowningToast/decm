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
import { useIssuerManagement } from "@/hooks/useIssuerManagement";
import { useCertificateTemplate } from "@/hooks/useCertificateTemplate";
import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import { useUpdateCertificateConfig } from "./useUpdateCertificateConfig";
import type {
    CoreApiInternalHandlerEventconfigEventCertificateConfigResponse,
    GetEventIssuersByEventIdData,
    GetVerifiedIssuersData,
    UpdateEventCertificateConfigPayload,
} from "@decm/api";
import { toast } from "sonner";
import { useNavigate } from "@/router";
import { useUpdateEventIssuer } from "./useUpdateEventIssuer";
import { useDeleteEventIssuer } from "./useDeleteEventIssuer";

interface CertificateSettingsPageProps {
    eventId: string;
    eventCertificateConfig?: CoreApiInternalHandlerEventconfigEventCertificateConfigResponse;
    verifiedIssuers?: GetVerifiedIssuersData;
    eventIssuers?: GetEventIssuersByEventIdData;
}

export const CertificateSettingsPage = ({
    eventId,
    eventCertificateConfig,
    verifiedIssuers,
    eventIssuers,
}: CertificateSettingsPageProps) => {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const { updateCertificateConfig, isUpdatingCertificateConfig } = useUpdateCertificateConfig(
        eventId!,
    );

    const { updateEventIssuer, isUpdatingEventIssuer } = useUpdateEventIssuer(eventId!);
    const { deleteEventIssuerAsync } = useDeleteEventIssuer();

    // Use custom hooks for state management
    const issuerManagement = useIssuerManagement({
        verifiedIssuers,
        selectedIssuers: eventIssuers,
    });
    const certificateTemplate = useCertificateTemplate();

    // Handle form submission
    const handleSubmit = async () => {
        try {
            // TODO: Implement API call to save certificate settings
            console.log("Event ID:", eventId);
            console.log("Selected Issuers:", issuerManagement.selectedIssuers);
            console.log("SVG File:", certificateTemplate.svgFile);
            console.log("Detected Keywords:", certificateTemplate.detectedKeywords);

            const name = certificateTemplate.detectedKeywords.find(
                (keyword) => keyword.keyword === "{{ name }}",
            );

            const eventName = certificateTemplate.detectedKeywords.find(
                (keyword) => keyword.keyword === "{{ eventName }}",
            );

            const acedmicInstitutionName = certificateTemplate.detectedKeywords.find(
                (keyword) => keyword.keyword === "{{ academicInstitutionName }}",
            );

            if (certificateTemplate.svgFile && !name) {
                toast.error(t("certificateSettings.nameNotFound"));
                return;
            }

            if (certificateTemplate.svgFile && !eventName) {
                toast.error(t("certificateSettings.eventNameNotFound"));
                return;
            }

            if (certificateTemplate.svgFile && !acedmicInstitutionName) {
                toast.error(t("certificateSettings.svgFileNotFound"));
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

            await updateCertificateConfig(req);
            await updateEventIssuer([
                {
                    event_id: eventId,
                    issuer_credential_id: "456e0d34-8497-4e2f-9992-ecb66af8f05a",
                },
            ]);

            toast.success(t("certificateSettings.saveSuccess"));
        } catch (error) {
            console.error("Error saving certificate settings:", error);
            toast.error(t("certificateSettings.saveError"));
        }
    };

    const handleRemoveIssuer = async (issuerId: string) => {
        try {
            await deleteEventIssuerAsync({ eventId, issuerId });
            issuerManagement.handleRemoveIssuer(issuerId);
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
    const isUpdateFormValid = true;

    const isFormValid = !eventCertificateConfig ? isCreateFormValid : isUpdateFormValid;

    return (
        <PageContainer
            title={t("certificateSettings.pageTitle")}
            description={t("certificateSettings.pageDescription")}
        >
            <SectionContainer>
                <div className="space-y-8">
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
                                <div className="flex gap-2 mt-2">
                                    <Input
                                        id="issuer-search"
                                        type="text"
                                        placeholder={t(
                                            "certificateSettings.step1.searchPlaceholder",
                                        )}
                                        value={issuerManagement.searchQuery}
                                        onChange={(e) =>
                                            issuerManagement.handleSearchQueryChange(e.target.value)
                                        }
                                        onKeyDown={(e) => {
                                            if (e.key === "Enter") {
                                                issuerManagement.handleSearch();
                                            }
                                        }}
                                        disabled={issuerManagement.isSearching}
                                    />
                                    <Button
                                        type="button"
                                        onClick={issuerManagement.handleSearch}
                                        disabled={
                                            issuerManagement.isSearching ||
                                            !issuerManagement.searchQuery.trim()
                                        }
                                        className="min-w-[100px]"
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
                            </div>

                            {/* Selected Issuers Table */}
                            <SelectedIssuersTable
                                selectedIssuers={issuerManagement.selectedIssuers}
                                onRemoveIssuer={handleRemoveIssuer}
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
                        />
                    </div>

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
        </PageContainer>
    );
};
