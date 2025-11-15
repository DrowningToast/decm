import { useState } from "react";
import { useNavigate } from "@/router";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { useEvent } from "@/hooks/events/useEvent";
import { useEventIssuers } from "@/components/pages/HostPages/EventPages/useEventIssuers";
import { useEventCertificates } from "@/hooks/useEventCertificates";
import { useEventCertificateConfig } from "@/components/pages/HostPages/EventPages/useEventCertificateConfig";
import { useEventContract } from "@/hooks/events/useEventContracts";
import { useSignEventCertificates } from "@/hooks/useSignEventCertificates";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { ExternalLinkIcon, CheckCircle, Clock } from "lucide-react";
import { PasswordPinModal } from "@/components/ui/password-pin-modal";
import { TextLabelValue } from "@/components/ui/text-label-value";
import { DataTable } from "@/components/ui/data-table";
import { CertificateColumns } from "@/components/pages/HostPages/EventsPage/columns/CertificateColumns";
import { formatEthereumAddress } from "@/lib/utils";
import SectionContainer from "@/components/container/SectionContainer";
import { IssuerStatusBadge } from "./IssuerStatusBadge";
import { IssuersStatus } from "./IssuersStatus";
import { useAuth } from "@/context/AuthContext";
import type { EntityEventCertificate } from "@decm/api";

interface IssuerSignPageProps {
    eventId: string;
}

export default function IssuerSignPage({ eventId }: IssuerSignPageProps) {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [showPinModal, setShowPinModal] = useState(false);
    const { user: currentUser } = useAuth();

    // Fetch event data
    const { event, isLoadingEventError } = useEvent(eventId);
    const { eventIssuers } = useEventIssuers(eventId);
    const { certificates: eventCertificates } = useEventCertificates(eventId);
    const { data: eventCertificateConfig } = useEventCertificateConfig(eventId);
    const { data: eventContract } = useEventContract(eventId);
    const { signEventCertificates, isSigning } = useSignEventCertificates();

    // Get current issuer's information
    const currentIssuer = eventIssuers?.find(
        (issuer) => issuer.issuer_credential_id === currentUser?.authenticationCredentialId,
    );

    // Determine if current issuer has already signed
    const hasCurrentIssuerSigned = currentIssuer && currentIssuer.is_signed === 1;
    const isCurrentIssuerPending = currentIssuer && currentIssuer.is_signed === 0;

    // Calculate total certificates to be signed
    const certificatesToSign =
        eventCertificates?.filter((cert) => !cert.revoked_at && currentIssuer) || [];

    const handleSignCertificates = () => {
        setShowPinModal(true);
    };

    const handleGoBack = () => {
        navigate(-1);
    };

    const handlePinModalClose = () => {
        setShowPinModal(false);
    };

    const handleSignWithPin = (pin: string) => {
        signEventCertificates({
            eventId,
            issuerPin: pin,
        });
    };

    if (isLoadingEventError || !event) {
        return (
            <ProtectedRoute>
                <div className="flex items-center justify-center min-h-screen">
                    <div className="text-center">
                        <Typography variant="text" tag="p" color="destructive">
                            {t("events.errors.generic")}
                        </Typography>
                        <Button onClick={handleGoBack} className="mt-4">
                            {t("common.back")}
                        </Button>
                    </div>
                </div>
            </ProtectedRoute>
        );
    }

    return (
        <ProtectedRoute>
            <div className="min-h-screen  text-white mt-12">
                {/* Main Content */}
                <SectionContainer>
                    <main>
                        {/* Page Header */}
                        <div className="mb-8">
                            <Typography
                                variant="header"
                                tag="h1"
                                className="text-4xl font-bold text-white mb-2"
                            >
                                {t("issuer.sign.pageTitle")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted-foreground"
                                className="text-lg"
                            >
                                {t("issuer.sign.pageDescription")}
                            </Typography>
                        </div>
                        {/* Issuers Status Section */}
                        {eventIssuers && eventIssuers.length > 0 && (
                            <IssuersStatus
                                issuers={
                                    eventIssuers?.map((issuer) => ({
                                        id: issuer.id || "",
                                        issuer_credential_id: issuer.issuer_credential_id || "",
                                        is_signed: issuer.is_signed || 0,
                                    })) || []
                                }
                                currentIssuerId={currentUser?.profileId}
                                className="mb-8"
                            />
                        )}

                        {/* Signing Details Section */}
                        <div className="bg-[#1a1a1a] border border-[#333333] rounded-lg p-6 mb-8">
                            <Typography
                                variant="header"
                                tag="h2"
                                className="text-xl font-semibold text-white mb-4"
                            >
                                {t("issuer.sign.signingDetails")}
                            </Typography>
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div>
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        color="muted-foreground"
                                        className="text-sm"
                                    >
                                        {t("issuer.sign.signingStatus")}
                                    </Typography>
                                    <div className="mt-1">
                                        {currentIssuer ? (
                                            <div className="flex items-center space-x-2">
                                                <IssuerStatusBadge
                                                    isSigned={currentIssuer.is_signed || 0}
                                                />
                                                {hasCurrentIssuerSigned && (
                                                    <Typography
                                                        variant="text"
                                                        tag="span"
                                                        className="text-sm text-green-500"
                                                    >
                                                        Completed
                                                    </Typography>
                                                )}
                                            </div>
                                        ) : (
                                            "N/A"
                                        )}
                                    </div>
                                </div>
                                <TextLabelValue
                                    label={t("issuer.sign.contractAddress")}
                                    value={
                                        eventContract?.certificate_contract_address
                                            ? formatEthereumAddress(
                                                  eventContract.certificate_contract_address,
                                              )
                                            : "N/A"
                                    }
                                    endIcon={<ExternalLinkIcon size={16} />}
                                    valueClassName="cursor-pointer underline"
                                    href={`https://www.etherscan.io/address/${eventContract?.certificate_contract_address}`}
                                />
                                <TextLabelValue
                                    label={t("issuer.sign.certificatesCount")}
                                    value={`${certificatesToSign.length} ${t("issuer.sign.certificates")}`}
                                />

                                <TextLabelValue
                                    label="Event"
                                    value={event.title ?? ""}
                                    endIcon={<ExternalLinkIcon size={16} />}
                                    valueClassName="cursor-pointer underline"
                                    href={`/events/${eventId}`}
                                />

                                <TextLabelValue
                                    label="Start Date"
                                    value={new Date(event.start_date ?? "").toLocaleDateString()}
                                />
                                <TextLabelValue
                                    label="End Date"
                                    value={new Date(event.end_date ?? "").toLocaleDateString()}
                                />
                            </div>
                        </div>

                        {/* Main Content Grid */}
                        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                            {/* Left Column - All Sections */}
                            <div className="lg:col-span-2 space-y-8">
                                {/* Certificate Preview Section */}
                                <div>
                                    <Typography
                                        variant="header"
                                        tag="h2"
                                        className="text-xl font-semibold text-white mb-4"
                                    >
                                        {t("issuer.sign.certificatePreview")}
                                    </Typography>
                                    <div className="bg-[#1a1a1a] border border-[#333333] rounded-lg p-4">
                                        <div className="aspect-w-16 aspect-h-9  rounded flex items-center justify-center border border-dashed border-[#333333]">
                                            {eventCertificateConfig?.base_certificate_presigned_url ? (
                                                <img
                                                    src={
                                                        eventCertificateConfig.base_certificate_presigned_url
                                                    }
                                                    alt="Certificate Preview"
                                                    className="w-full h-full object-contain rounded-md"
                                                />
                                            ) : (
                                                <div className="text-center text-[#a0a0a0]">
                                                    <Typography variant="text" tag="p">
                                                        {t("issuer.sign.templateNotAvailable")}
                                                    </Typography>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Right Column - Action Card */}
                            <div className="lg:col-span-1">
                                {/* Event Details Section */}
                                <div className="bg-[#1a1a1a] border border-[#333333] rounded-lg p-6 mb-8 flex flex-col">
                                    <Typography
                                        variant="header"
                                        tag="h2"
                                        className="text-xl font-semibold text-white mb-4"
                                    >
                                        {t("issuer.sign.eventDetails")}
                                    </Typography>
                                    <div className="space-y-4">
                                        <TextLabelValue
                                            label={t("events.details.status")}
                                            value={event.event_status?.toUpperCase() ?? "NA"}
                                        />
                                        <TextLabelValue
                                            label={t("events.details.participationRequest")}
                                            value={
                                                event.is_booking_request_required
                                                    ? t("common.required")
                                                    : t("common.notRequired")
                                            }
                                        />
                                        <TextLabelValue
                                            label={t("events.details.seatsCount")}
                                            value={`${event.attendees_count ?? 0} / ${event.max_attendees}`}
                                        />
                                        <TextLabelValue
                                            label={t("events.details.eventContractAddress")}
                                            value={
                                                eventContract?.event_contract_address
                                                    ? formatEthereumAddress(
                                                          eventContract.event_contract_address,
                                                      )
                                                    : "NA"
                                            }
                                            endIcon={<ExternalLinkIcon size={16} />}
                                            valueClassName="cursor-pointer underline"
                                            href={`https://www.etherscan.io/address/${eventContract?.event_contract_address}`}
                                        />
                                    </div>
                                </div>

                                <div className="bg-[#1a1a1a] border border-[#333333] p-6 rounded-lg shadow-lg space-y-6 sticky top-8">
                                    <div>
                                        <Typography
                                            variant="header"
                                            tag="h3"
                                            className="text-xl font-semibold text-white mb-4"
                                        >
                                            {t("issuer.sign.signCertificateBatch")}
                                        </Typography>
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            color="muted-foreground"
                                            className="text-sm mb-4"
                                        >
                                            {t("issuer.sign.signCertificateBatchDescription")}
                                        </Typography>

                                        {/* Sign Button */}
                                        {isCurrentIssuerPending ? (
                                            <Button
                                                onClick={handleSignCertificates}
                                                disabled={
                                                    isSigning || certificatesToSign.length === 0
                                                }
                                                className="w-full bg-[#ff6a39] text-white font-semibold py-3 px-6 rounded-lg mb-3"
                                            >
                                                {isSigning
                                                    ? t("issuer.sign.signing")
                                                    : `${t("issuer.sign.signAndApprove")} (${certificatesToSign.length})`}
                                            </Button>
                                        ) : hasCurrentIssuerSigned ? (
                                            <div className="w-full bg-green-600 text-white font-semibold py-3 px-6 rounded-lg mb-3 flex items-center justify-center">
                                                <CheckCircle className="mr-2 h-5 w-5" />
                                                {t("issuer.sign.alreadySigned")}
                                            </div>
                                        ) : (
                                            <div className="w-full bg-gray-600 text-white font-semibold py-3 px-6 rounded-lg mb-3 flex items-center justify-center">
                                                <Clock className="mr-2 h-5 w-5" />
                                                {t("issuer.sign.notAnIssuer")}
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </main>
                </SectionContainer>

                {/* Event Certificates Section */}
                <SectionContainer>
                    <div>
                        <Typography
                            variant="header"
                            tag="h2"
                            className="text-xl font-semibold text-white mb-4"
                        >
                            {isCurrentIssuerPending
                                ? t("issuer.sign.certificatesToSign")
                                : hasCurrentIssuerSigned
                                  ? t("issuer.sign.signedCertificates")
                                  : t("issuer.sign.eventCertificates")}{" "}
                            (
                            {eventCertificates?.filter((cert) => cert.revoked_at === null).length ||
                                0}
                            )
                        </Typography>

                        <DataTable
                            columns={CertificateColumns()}
                            data={
                                eventCertificates?.map((cert) => {
                                    const firstName = cert.name?.split(" ")[0] || "";
                                    const lastName = cert.name?.split(" ").slice(1).join(" ") || "";

                                    return {
                                        ...cert,
                                        firstName,
                                        lastName,
                                        email: cert.receiver_email || "",
                                        academicInstitution: cert.academic_institution || "",
                                        issuedAt: cert.created_at,
                                        status: cert.revoked_at
                                            ? "rejected"
                                            : isCurrentIssuerPending
                                              ? "pending_signature"
                                              : "received",
                                    } as EntityEventCertificate;
                                }) || []
                            }
                            totalItems={eventCertificates?.length || 0}
                            currentPage={1}
                            pageSize={10}
                            onPageChange={() => {}}
                            onPageSizeChange={() => {}}
                            searchValue=""
                            onSearchChange={() => {}}
                            searchPlaceholder="Search certificates..."
                            sorting={[]}
                            onSortingChange={() => {}}
                            isLoading={false}
                            disablePagination
                        />
                    </div>
                </SectionContainer>

                {/* Password/PIN Modal */}
                <PasswordPinModal
                    isOpen={showPinModal}
                    onClose={handlePinModalClose}
                    onSuccess={(result: { type: "pin" | "password"; value: string }) => {
                        handleSignWithPin(result.value);
                    }}
                    title={t("issuer.sign.signatureRequest")}
                    description={t("issuer.sign.signatureRequestDescription")}
                    showSigningDetails={true}
                    signingDetails={{
                        contractAddress: eventContract?.event_contract_address,
                        transactionType: "Certificate Signing",
                        details: `Signing ${certificatesToSign.length} certificates for ${event?.title}`,
                    }}
                />
            </div>
        </ProtectedRoute>
    );
}
