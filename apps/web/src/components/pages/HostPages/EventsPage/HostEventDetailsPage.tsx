import SectionContainer from "@/components/container/SectionContainer";
import { Typography } from "@/components/typography/typography";
import { GoogleMapsEmbed } from "@/components/ui/google-maps-embed";
import {
    StyledTabs,
    StyledTabsList,
    StyledTabsTrigger,
    StyledTabsContent,
} from "@/components/ui/styled-tabs";
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "@/components/ui/accordion";
import { RequirementItem } from "@/components/ui/requirement-item";
import { TextLabelValue } from "@/components/ui/text-label-value";
import WrappedButton from "@/components/wrapper/WrappedButton";
import { CheckCircle2Icon, CloudUploadIcon, ExternalLinkIcon } from "lucide-react";
import { useTranslation } from "react-i18next";
import { DataTable } from "@/components/ui/data-table";
import { useIssuerColumns } from "./columns/issuer-columns";
import type {
    EntityEventCertificate,
    GetEventCertificateConfigData,
    GetEventContractByEventIdData,
} from "@decm/api";
import { formatEthereumAddress } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { useState, useMemo, useCallback, useEffect } from "react";
import type { SortingState } from "@tanstack/react-table";
import { useParticipantColumns, type Participant } from "./columns/useParticipantColumns";
import { useAttendeeColumns } from "./columns/useAttendeeColumns";
import { useEventAttendees } from "@/hooks/events/useEventAttendees";
import { CertificateColumns } from "./columns/CertificateColumns";
import { Separator } from "@/components/ui/separator";
import { useEventCertificates } from "@/hooks/useEventCertificates";
import { useRevokeEventCertificate } from "@/hooks/events/useRevokeEventCertificate";
import { useToggleCertificatePublished } from "@/hooks/events/useToggleCertificatePublished";
import { Link } from "@/router";
import type {
    EventRegistrationConfiguration,
    EventRegistrationInvitation,
} from "@/services/EventRegistration/EventRegistration";
import { EventStatusesViewModel, EventTypesViewModel } from "./ViewModel";
import type { EventIssuer, EventViewModel } from "@/services/EventService/EventService";

interface HostEventDetailsPageProps {
    eventId: string;
    event: EventViewModel;
    eventRegistrationConfig: EventRegistrationConfiguration;
    eventCertificateConfig?: GetEventCertificateConfigData;
    eventIssuers?: EventIssuer[];
    eventContract?: GetEventContractByEventIdData;
    eventInvitations?: EventRegistrationInvitation[];
}

export default function HostEventDetailsPage({
    eventId,
    event,
    eventRegistrationConfig,
    eventCertificateConfig,
    eventIssuers,
    eventContract,
    eventInvitations,
}: HostEventDetailsPageProps) {
    const { t } = useTranslation();

    const { revokeEventCertificate } = useRevokeEventCertificate();
    const { mutate: togglePublished, isPending: isTogglingPublished } =
        useToggleCertificatePublished();

    // Get mint readiness from certificate config
    const mintReadiness = eventCertificateConfig?.mint_readiness;

    // Check if certificate config is published - if so, disable all edit buttons
    const isCertificatePublished = eventCertificateConfig?.is_published ?? false;

    // State for client-side data management
    const [searchValue, setSearchValue] = useState("");
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [sorting, setSorting] = useState<SortingState>([]);
    const [debouncedSearch, setDebouncedSearch] = useState("");

    // Debounce search value
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(searchValue);
            setCurrentPage(1); // Reset to first page when search changes
        }, 300);

        return () => clearTimeout(timer);
    }, [searchValue]);

    // Memoized filtered and sorted data
    const processedData = useMemo(() => {
        console.log(eventInvitations);
        if (!eventInvitations) return [];

        let filtered = [...eventInvitations];

        // Filter by search value
        if (debouncedSearch) {
            const searchLower = debouncedSearch.toLowerCase();
            filtered = filtered.filter(
                (item) =>
                    item.firstName?.toLowerCase().includes(searchLower) ||
                    item.lastName?.toLowerCase().includes(searchLower) ||
                    item.email?.toLowerCase().includes(searchLower) ||
                    item.academicInstitution?.toLowerCase().includes(searchLower) ||
                    item.phoneNumber?.toLowerCase().includes(searchLower),
            );
        }

        // Apply sorting
        if (sorting.length > 0) {
            const sort = sorting[0];
            filtered.sort((a, b) => {
                const aValue = a[sort.id as keyof typeof a];
                const bValue = b[sort.id as keyof typeof b];

                if (aValue === undefined || bValue === undefined) return 0;

                let comparison = 0;
                if (aValue < bValue) comparison = -1;
                if (aValue > bValue) comparison = 1;

                return sort.desc ? comparison * -1 : comparison;
            });
        }

        return filtered;
    }, [eventInvitations, debouncedSearch, sorting]);

    // Calculate pagination and map to Participant type
    const paginatedData = useMemo(() => {
        const startIndex = (currentPage - 1) * pageSize;
        const endIndex = startIndex + pageSize;
        return processedData.slice(startIndex, endIndex).map(
            (invitation): Participant => ({
                id: invitation.id,
                firstName: invitation.firstName ?? null,
                lastName: invitation.lastName ?? null,
                email: invitation.email ?? null,
                phoneNumber: invitation.phoneNumber ?? null,
                academicInstitution: invitation.academicInstitution ?? null,
                walletAddress: "", // Not available in invitation data
                status: invitation.cancelledAt ? "rejected" : "pending",
                isAccepted: !!invitation.acceptedAt, // Convert acceptedAt timestamp to boolean
            }),
        );
    }, [processedData, currentPage, pageSize]);

    const participantColumns = useParticipantColumns(eventId);
    const attendeeColumns = useAttendeeColumns();
    const issuerColumns = useIssuerColumns();

    // Fetch actual attendees
    const { attendees, isLoading: attendeesLoading } = useEventAttendees({ eventId });

    // Callbacks for DataTable
    const handlePageChange = useCallback((page: number) => {
        setCurrentPage(page);
    }, []);

    const handlePageSizeChange = useCallback((newPageSize: number) => {
        setPageSize(newPageSize);
        setCurrentPage(1); // Reset to first page when page size changes
    }, []);

    const handleSearchChange = useCallback((value: string) => {
        setSearchValue(value);
    }, []);

    const handleSortingChange = useCallback((newSorting: SortingState) => {
        setSorting(newSorting);
    }, []);

    // Certificate state logic
    const hasCertificateConfig = !!eventCertificateConfig;

    // Check if issuers are assigned
    const hasIssuers = !!eventIssuers && eventIssuers.length > 0;

    // Fetch event certificates using the hook
    const { certificates: eventCertificates, isLoading: certificatesLoading } =
        useEventCertificates(eventId);

    // Client-side check: Must have at least 1 receiver (non-revoked certificate)
    const hasAtLeastOneReceiver = useMemo(() => {
        if (!eventCertificates) return false;
        const nonRevokedCertificates = eventCertificates.filter((cert) => !cert.revokedAt);
        return nonRevokedCertificates.length > 0;
    }, [eventCertificates]);

    // Combined readiness check (backend + client-side receiver requirement)
    const isReadyToPublish = useMemo(() => {
        if (!mintReadiness) return false;
        return mintReadiness.is_ready && hasAtLeastOneReceiver;
    }, [mintReadiness, hasAtLeastOneReceiver]);

    console.log(eventCertificates);

    return (
        <div className="flex flex-col gap-y-6">
            <SectionContainer>
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <img
                            src={event.iconPresignedUrl ?? ""}
                            alt={event.title}
                            className="w-12 h-12 rounded-full object-cover"
                        />

                        <div className="flex items-center gap-2">
                            <Typography variant="header" tag="p" size={"header"}>
                                {event.title}
                            </Typography>
                            {event.isVerified && <CheckCircle2Icon color="#eb5331" />}
                        </div>
                    </div>

                    <div className="flex items-center gap-4">
                        {/* TODO: Add confirm event button */}
                        {/* <Button size="xl" variant="secondary-light">
                            {t("events.hostDetails.actions.confirmEvent")}
                        </Button> */}
                        <Link to={`/host/events/:eventId/edit`} params={{ eventId }}>
                            <Button size="lg">{t("events.hostDetails.actions.editEvent")}</Button>
                        </Link>
                    </div>
                </div>
            </SectionContainer>

            <SectionContainer className="lg:grid lg:grid-cols-4 gap-8">
                <div className="flex flex-col gap-4 lg:col-span-3">
                    <Typography tag="p" size={"base"} color="muted">
                        {event.shortDescription}
                    </Typography>

                    <Typography tag="p" size={"base"} color="muted">
                        {event.longDescription}
                    </Typography>

                    <img
                        src={event.bannerPresignedUrl ?? ""}
                        alt={event.title}
                        className="w-full h-[350px] object-cover rounded-lg"
                    />
                </div>

                <div className="flex flex-col gap-4 mt-6 lg:mt-0 border rounded-lg p-6 ">
                    <TextLabelValue
                        label={t("events.details.status")}
                        value={
                            EventStatusesViewModel[event.eventStatus] ?? t("common.notAvailable")
                        }
                    />
                    <TextLabelValue
                        label={t("events.details.eventType")}
                        value={EventTypesViewModel[event.eventType] ?? t("common.notAvailable")}
                    />
                    <TextLabelValue
                        label={t("events.details.seatsCount")}
                        value={`${event.attendeesCount ?? 0} / ${event.maxAttendees}`}
                    />
                    <TextLabelValue
                        label={t("events.details.eventContractAddress")}
                        value={
                            eventContract?.event_contract_address
                                ? formatEthereumAddress(eventContract.event_contract_address)
                                : t("common.notAvailable")
                        }
                        endIcon={<ExternalLinkIcon size={16} />}
                        valueClassName="cursor-pointer underline"
                        href={`https://www.etherscan.io/address/${eventContract?.event_contract_address}`}
                    />
                </div>
            </SectionContainer>

            <SectionContainer>
                <StyledTabs defaultValue="event-info">
                    <StyledTabsList>
                        <StyledTabsTrigger value="event-info">
                            <Typography variant="text" tag="span" color="current">
                                {t("events.hostDetails.tabs.eventInfo")}
                            </Typography>
                        </StyledTabsTrigger>
                        <StyledTabsTrigger value="participants">
                            <Typography variant="text" tag="span" color="current">
                                {t("events.hostDetails.tabs.participants")}
                            </Typography>
                        </StyledTabsTrigger>
                        <StyledTabsTrigger value="certificates">
                            <Typography variant="text" tag="span" color="current">
                                {t("events.hostDetails.tabs.certificates")}
                            </Typography>
                        </StyledTabsTrigger>
                    </StyledTabsList>
                    <StyledTabsContent value="event-info">
                        <div className="flex flex-col gap-4 lg:flex-row">
                            <div className="flex flex-col gap-4 flex-1">
                                <TextLabelValue
                                    label={t("events.form.location")}
                                    value={event.location ?? ""}
                                />
                                <TextLabelValue
                                    label={t("events.form.googleMapQuery")}
                                    value={event.googleMapQuery ?? ""}
                                />

                                <TextLabelValue
                                    label={t("events.form.contactAddress")}
                                    value={event.contactNumber ?? ""}
                                />
                                <TextLabelValue
                                    label={t("events.hostDetails.eventInfo.contact")}
                                    value={event.contactNumber ?? ""}
                                />
                            </div>
                            <div className="flex-1">
                                <GoogleMapsEmbed query={event.googleMapQuery ?? ""} />
                            </div>
                        </div>
                    </StyledTabsContent>
                    <StyledTabsContent value="participants">
                        <div className="space-y-4">
                            {/* Event Settings Summary */}
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 my-6">
                                <TextLabelValue
                                    label={t("events.settings.eventType")}
                                    value={
                                        EventTypesViewModel[event.eventType] ??
                                        t("common.notAvailable")
                                    }
                                />
                                <TextLabelValue
                                    label={t("events.settings.bookingRequired")}
                                    value={
                                        event.isBookingRequestRequired
                                            ? t("common.yes")
                                            : t("common.no")
                                    }
                                />
                                <TextLabelValue
                                    label={t("events.settings.tokenTransferable")}
                                    value={
                                        event.isTicketTransferable
                                            ? t("common.yes")
                                            : t("common.no")
                                    }
                                />

                                <div className="flex items-center justify-end gap-4">
                                    <Link
                                        to="/host/events/:eventId/imports/participants"
                                        params={{ eventId }}
                                    >
                                        <Button size="lg" variant="secondary-light">
                                            {t("participantImport.title")}
                                        </Button>
                                    </Link>
                                    <Link
                                        to="/host/events/:eventId/settings/participant"
                                        params={{ eventId }}
                                    >
                                        <Button size="lg" variant="primary">
                                            {t("events.settings.participantSettings")}
                                        </Button>
                                    </Link>
                                </div>
                            </div>

                            {/* Participant Requirements Accordion */}
                            <Accordion
                                type="single"
                                collapsible
                                className="w-full bg-white rounded-lg "
                            >
                                <AccordionItem value="requirements">
                                    <AccordionTrigger className="hover:no-underline px-4">
                                        <div className="flex items-center justify-between w-full">
                                            <Typography
                                                variant="text"
                                                tag="span"
                                                className="font-medium text-black"
                                            >
                                                {t("events.participants.requirementsTitle")}
                                            </Typography>
                                        </div>
                                    </AccordionTrigger>
                                    <AccordionContent>
                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4">
                                            <RequirementItem
                                                label={t("events.participants.fields.firstName")}
                                                status={eventRegistrationConfig.firstName}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.lastName")}
                                                status={eventRegistrationConfig.lastName}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.email")}
                                                status={eventRegistrationConfig.email}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.bio")}
                                                status={eventRegistrationConfig.bio}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.phoneNumber")}
                                                status={eventRegistrationConfig.phoneNumber}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.address")}
                                                status={eventRegistrationConfig.address}
                                            />
                                            <RequirementItem
                                                label={t(
                                                    "events.participants.fields.academicInstitution",
                                                )}
                                                status={eventRegistrationConfig.academicInstitution}
                                            />
                                            <RequirementItem
                                                label={t(
                                                    "events.participants.fields.academicEmail",
                                                )}
                                                status={eventRegistrationConfig.academicEmail}
                                            />
                                        </div>
                                    </AccordionContent>
                                </AccordionItem>
                            </Accordion>

                            {/* Actual Attendees Section */}
                            <div className="space-y-4">
                                <Typography
                                    variant="text"
                                    tag="h3"
                                    className="text-lg font-semibold"
                                >
                                    {t("events.hostDetails.attendees.title")}
                                </Typography>
                                <Typography variant="text" tag="p" size="small" color="muted">
                                    {t("events.hostDetails.attendees.description")}
                                </Typography>
                                <DataTable
                                    columns={attendeeColumns}
                                    data={attendees}
                                    totalItems={attendees.length}
                                    currentPage={1}
                                    pageSize={attendees.length}
                                    onPageChange={() => {}}
                                    onPageSizeChange={() => {}}
                                    searchValue=""
                                    onSearchChange={() => {}}
                                    searchPlaceholder={t(
                                        "events.hostDetails.attendees.searchPlaceholder",
                                    )}
                                    sorting={[]}
                                    onSortingChange={() => {}}
                                    isLoading={attendeesLoading}
                                    disablePagination
                                />
                            </div>

                            <Separator className="my-8" />

                            {/* Invitations Section */}
                            <Typography variant="text" tag="h3" className="text-lg font-semibold">
                                {t("events.hostDetails.tabs.participants")} (Invitations)
                            </Typography>
                            <Typography variant="text" tag="p" size="small" color="muted">
                                {t("events.hostDetails.participants.description")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                size="small"
                                color="muted"
                                className="italic"
                            >
                                {t("events.hostDetails.participants.cancelWarning")}
                            </Typography>

                            <DataTable
                                columns={participantColumns}
                                data={paginatedData}
                                totalItems={processedData.length}
                                currentPage={currentPage}
                                pageSize={pageSize}
                                onPageChange={handlePageChange}
                                onPageSizeChange={handlePageSizeChange}
                                searchValue={searchValue}
                                onSearchChange={handleSearchChange}
                                searchPlaceholder={t(
                                    "events.hostDetails.participants.searchPlaceholder",
                                )}
                                sorting={sorting}
                                onSortingChange={(value) =>
                                    handleSortingChange(value as SortingState)
                                }
                                isLoading={false}
                            />
                        </div>
                    </StyledTabsContent>
                    <StyledTabsContent value="certificates">
                        {hasCertificateConfig ? (
                            <div className="space-y-6">
                                {/* Certificate Settings Section */}
                                <div className="w-full bg-white border border-white/50 rounded-lg p-6 flex flex-row items-center justify-between">
                                    <div>
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            className="font-semibold text-lg text-black"
                                        >
                                            {t("certificateSettings.pageTitle")}
                                        </Typography>
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            className="text-black/50 text-base mt-1"
                                        >
                                            {t(
                                                "events.hostDetails.certificates.summaryDescription",
                                            )}
                                        </Typography>
                                    </div>
                                    <div className="flex gap-2">
                                        <WrappedButton
                                            className="px-5 py-2 rounded-md font-medium transition bg-primary text-white hover:bg-primary/90"
                                            href={`/host/events/${eventId}/settings/certificate`}
                                        >
                                            {t("certificateSettings.pageTitle")}
                                        </WrappedButton>
                                    </div>
                                </div>

                                {/* Certificate Mint Readiness Status */}
                                {mintReadiness && (
                                    <div
                                        className={`w-full border rounded-lg p-4 flex items-start gap-3 ${
                                            isReadyToPublish
                                                ? "bg-green-50 border-green-200"
                                                : "bg-blue-50 border-blue-200"
                                        }`}
                                    >
                                        <div className="mt-0.5">
                                            {isReadyToPublish ? (
                                                <CheckCircle2Icon className="h-5 w-5 text-green-600" />
                                            ) : (
                                                <div className="h-5 w-5 rounded-full border-2 border-blue-600" />
                                            )}
                                        </div>
                                        <div className="flex-1">
                                            <Typography
                                                variant="text"
                                                tag="p"
                                                className={`font-semibold text-base ${
                                                    isReadyToPublish
                                                        ? "text-green-900"
                                                        : "text-blue-900"
                                                }`}
                                            >
                                                {isReadyToPublish
                                                    ? t("events.hostDetails.certificates.mintReady")
                                                    : t(
                                                          "events.hostDetails.certificates.mintNotReady",
                                                      )}
                                            </Typography>
                                            <div className="mt-2 space-y-1">
                                                {/* Certificate Config Status */}
                                                <div className="flex items-center gap-2 text-sm">
                                                    {mintReadiness.has_certificate_config ? (
                                                        <CheckCircle2Icon className="h-4 w-4 text-green-600 flex-shrink-0" />
                                                    ) : (
                                                        <div className="h-4 w-4 rounded-full border-2 border-gray-400 flex-shrink-0" />
                                                    )}
                                                    <span
                                                        className={
                                                            mintReadiness.has_certificate_config
                                                                ? "text-green-700"
                                                                : "text-gray-600"
                                                        }
                                                    >
                                                        {t(
                                                            "events.hostDetails.certificates.configStatus",
                                                        )}
                                                    </span>
                                                </div>
                                                {/* All Issuers Signed Status */}
                                                <div className="flex items-center gap-2 text-sm">
                                                    {mintReadiness.all_issuers_have_signed ? (
                                                        <CheckCircle2Icon className="h-4 w-4 text-green-600 flex-shrink-0" />
                                                    ) : (
                                                        <div className="h-4 w-4 rounded-full border-2 border-gray-400 flex-shrink-0" />
                                                    )}
                                                    <span
                                                        className={
                                                            mintReadiness.all_issuers_have_signed
                                                                ? "text-green-700"
                                                                : "text-gray-600"
                                                        }
                                                    >
                                                        {t(
                                                            "events.hostDetails.certificates.issuersSignedStatus",
                                                            {
                                                                signed: mintReadiness.signed_issuers_count,
                                                                total: mintReadiness.total_issuers_count,
                                                            },
                                                        )}
                                                    </span>
                                                </div>
                                                {/* Contract Deployed Status */}
                                                <div className="flex items-center gap-2 text-sm">
                                                    {mintReadiness.has_certificate_contract ? (
                                                        <CheckCircle2Icon className="h-4 w-4 text-green-600 flex-shrink-0" />
                                                    ) : (
                                                        <div className="h-4 w-4 rounded-full border-2 border-gray-400 flex-shrink-0" />
                                                    )}
                                                    <span
                                                        className={
                                                            mintReadiness.has_certificate_contract
                                                                ? "text-green-700"
                                                                : "text-gray-600"
                                                        }
                                                    >
                                                        {t(
                                                            "events.hostDetails.certificates.contractStatus",
                                                        )}
                                                    </span>
                                                </div>
                                                {/* At Least One Receiver Status (Client-side check) */}
                                                <div className="flex items-center gap-2 text-sm">
                                                    {hasAtLeastOneReceiver ? (
                                                        <CheckCircle2Icon className="h-4 w-4 text-green-600 flex-shrink-0" />
                                                    ) : (
                                                        <div className="h-4 w-4 rounded-full border-2 border-gray-400 flex-shrink-0" />
                                                    )}
                                                    <span
                                                        className={
                                                            hasAtLeastOneReceiver
                                                                ? "text-green-700"
                                                                : "text-gray-600"
                                                        }
                                                    >
                                                        {t(
                                                            "events.hostDetails.certificates.hasReceiverStatus",
                                                            {
                                                                count:
                                                                    eventCertificates?.filter(
                                                                        (c) => !c.revokedAt,
                                                                    ).length || 0,
                                                            },
                                                        )}
                                                    </span>
                                                </div>
                                            </div>
                                            {/* Missing Requirements */}
                                            {mintReadiness.missing_requirements &&
                                                mintReadiness.missing_requirements.length > 0 && (
                                                    <div className="mt-3 p-2 bg-blue-100 rounded text-sm text-blue-800">
                                                        <Typography
                                                            variant="text"
                                                            tag="p"
                                                            className="font-medium mb-1"
                                                        >
                                                            {t(
                                                                "events.hostDetails.certificates.missingRequirements",
                                                            )}
                                                            :
                                                        </Typography>
                                                        <ul className="list-disc list-inside space-y-0.5">
                                                            {mintReadiness.missing_requirements.map(
                                                                (req, idx) => (
                                                                    <li key={idx}>{req}</li>
                                                                ),
                                                            )}
                                                        </ul>
                                                    </div>
                                                )}
                                        </div>
                                    </div>
                                )}

                                {/* Certificate Published Status Indicator */}
                                {eventCertificateConfig && (
                                    <div
                                        className={`w-full border rounded-lg p-4 flex items-start gap-3 ${
                                            eventCertificateConfig.is_published
                                                ? "bg-green-50 border-green-200"
                                                : "bg-gray-50 border-gray-200"
                                        }`}
                                    >
                                        <div className="mt-0.5">
                                            {eventCertificateConfig.is_published ? (
                                                <CheckCircle2Icon className="h-5 w-5 text-green-600" />
                                            ) : (
                                                <div className="h-5 w-5 rounded-full border-2 border-gray-400" />
                                            )}
                                        </div>
                                        <div className="flex-1">
                                            <Typography
                                                variant="text"
                                                tag="p"
                                                className={`font-semibold text-base ${
                                                    eventCertificateConfig.is_published
                                                        ? "text-green-900"
                                                        : "text-gray-900"
                                                }`}
                                            >
                                                {eventCertificateConfig.is_published
                                                    ? t(
                                                          "events.hostDetails.certificates.publishedConfigStatus",
                                                      )
                                                    : t(
                                                          "events.hostDetails.certificates.notPublishedConfigStatus",
                                                      )}
                                            </Typography>
                                            {eventCertificateConfig.is_published && (
                                                <Typography
                                                    variant="text"
                                                    tag="p"
                                                    className="text-sm text-green-700 mt-1"
                                                >
                                                    {t(
                                                        "events.hostDetails.certificates.publishedWarning",
                                                    )}
                                                </Typography>
                                            )}
                                            {!eventCertificateConfig.is_published && (
                                                <Typography
                                                    variant="text"
                                                    tag="p"
                                                    className="text-sm text-gray-600 mt-1"
                                                >
                                                    {t(
                                                        "events.hostDetails.certificates.notPublishedConfigDescription",
                                                    )}
                                                </Typography>
                                            )}
                                        </div>
                                        {!eventCertificateConfig.is_published && (
                                            <Button
                                                onClick={() =>
                                                    togglePublished({
                                                        eventId,
                                                        isPublished: true,
                                                    })
                                                }
                                                disabled={isTogglingPublished || !isReadyToPublish}
                                                className="ml-auto"
                                            >
                                                {t("events.hostDetails.certificates.publishButton")}
                                            </Button>
                                        )}
                                    </div>
                                )}

                                {/* Certificate List (Receivers) */}
                                <div className="space-y-4">
                                    <div>
                                        <Typography
                                            variant="text"
                                            tag="h3"
                                            className="text-lg font-semibold"
                                        >
                                            {t(
                                                "events.hostDetails.certificates.eventCertificatesTitle",
                                                { count: eventCertificates?.length || 0 },
                                            )}
                                        </Typography>
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            className="text-sm text-muted-foreground mt-1"
                                        >
                                            {t(
                                                "events.hostDetails.certificates.eventCertificatesDescription",
                                            )}
                                        </Typography>
                                    </div>

                                    <DataTable
                                        columns={CertificateColumns(
                                            (eventCertificateId: string) => {
                                                revokeEventCertificate({
                                                    certificateIds: [eventCertificateId],
                                                    eventId,
                                                });
                                            },
                                            isCertificatePublished,
                                        )}
                                        data={
                                            (eventCertificates || [])
                                                .filter((cert) => !cert.revokedAt)
                                                .map((cert) => {
                                                    const nameParts = (cert.name || "").split(" ");
                                                    return {
                                                        id: cert.id,
                                                        event_id: cert.eventId,
                                                        event_name: cert.eventName,
                                                        receiver_credential_id:
                                                            cert.receiverCredentialId,
                                                        receiver_email: cert.receiverEmail,
                                                        name: cert.name,
                                                        academic_institution:
                                                            cert.academicInstitution,
                                                        certificate_title: cert.certificateTitle,
                                                        certificate_subtitle:
                                                            cert.certificateSubtitle,
                                                        event_contract_address:
                                                            cert.eventContractAddress,
                                                        event_certificate_address:
                                                            cert.eventCertificateAddress,
                                                        certificate_token_id:
                                                            cert.certificateTokenId,
                                                        certificate_digest: undefined,
                                                        created_at: cert.createdAt,
                                                        revoked_at: cert.revokedAt,
                                                        inbox_message_id: cert.inboxMessageId,
                                                        firstName: nameParts[0] || "",
                                                        lastName:
                                                            nameParts.slice(1).join(" ") || "",
                                                        email: cert.receiverEmail || "",
                                                        academicInstitution:
                                                            cert.academicInstitution || "",
                                                        issuedAt: cert.createdAt || "",
                                                        status: "received" as const,
                                                    };
                                                }) as EntityEventCertificate[]
                                        }
                                        totalItems={
                                            eventCertificates?.filter((c) => !c.revokedAt).length ||
                                            0
                                        }
                                        currentPage={1}
                                        pageSize={10}
                                        onPageChange={() => {}}
                                        onPageSizeChange={() => {}}
                                        searchValue=""
                                        onSearchChange={() => {}}
                                        searchPlaceholder={t(
                                            "events.hostDetails.certificates.searchCertificatesPlaceholder",
                                        )}
                                        sorting={[]}
                                        onSortingChange={() => {}}
                                        isLoading={certificatesLoading}
                                        disablePagination
                                    />
                                </div>

                                {/* Import Certificate Receivers */}
                                <a
                                    href={
                                        isCertificatePublished || !hasIssuers
                                            ? undefined
                                            : `/host/events/${eventId}/imports/certificates`
                                    }
                                    className={`flex items-center justify-center p-6 border-dashed rounded-xl border-2 gap-4 ${
                                        isCertificatePublished || !hasIssuers
                                            ? "border-gray-300 opacity-50 cursor-not-allowed pointer-events-none"
                                            : "border-white/50 cursor-pointer hover:border-primary/30"
                                    }`}
                                    onClick={
                                        isCertificatePublished || !hasIssuers
                                            ? (e) => e.preventDefault()
                                            : undefined
                                    }
                                    title={
                                        isCertificatePublished
                                            ? t(
                                                  "events.hostDetails.certificates.editDisabledMessage",
                                              )
                                            : !hasIssuers
                                              ? t(
                                                    "events.hostDetails.certificates.noIssuersMessage",
                                                )
                                              : undefined
                                    }
                                >
                                    <CloudUploadIcon />{" "}
                                    <Typography variant="text" tag="span">
                                        {t("events.hostDetails.certificates.importSectionTitle")}
                                    </Typography>
                                </a>

                                {/* Issuers Table */}
                                {eventIssuers && eventIssuers.length > 0 && (
                                    <div className="space-y-6">
                                        <div className="flex flex-col gap-4">
                                            <div className="flex items-center justify-between">
                                                <Typography
                                                    variant="text"
                                                    tag="h3"
                                                    className="text-lg font-semibold"
                                                >
                                                    {t(
                                                        "events.hostDetails.certificates.issuersTitle",
                                                        { count: eventIssuers.length },
                                                    )}
                                                </Typography>
                                            </div>

                                            <DataTable
                                                columns={issuerColumns}
                                                data={eventIssuers}
                                                totalItems={eventIssuers.length}
                                                currentPage={1}
                                                pageSize={10}
                                                onPageChange={() => {}}
                                                onPageSizeChange={() => {}}
                                                searchValue=""
                                                onSearchChange={() => {}}
                                                searchPlaceholder={t(
                                                    "events.hostDetails.certificates.searchIssuersPlaceholder",
                                                )}
                                                sorting={[]}
                                                onSortingChange={() => {}}
                                                isLoading={false}
                                                disablePagination
                                            />
                                        </div>
                                    </div>
                                )}

                                {(!eventIssuers || eventIssuers.length === 0) && (
                                    <div className="text-center py-8">
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            className="text-muted-foreground"
                                        >
                                            {t(
                                                "events.hostDetails.certificates.noIssuersConfigured",
                                            )}
                                        </Typography>
                                    </div>
                                )}
                            </div>
                        ) : (
                            <div className="w-full bg-primary/10 border border-primary/20 rounded-lg p-6 flex flex-col items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="font-semibold text-lg text-primary"
                                >
                                    {t("events.hostDetails.certificates.noConfigTitle")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-muted-foreground text-base mt-1 text-center max-w-xl"
                                >
                                    {t("events.hostDetails.certificates.noConfigDescription")}
                                </Typography>
                                <div className="flex gap-4 mt-6 items-center">
                                    <WrappedButton
                                        className="px-5 rounded-md font-medium transition bg-primary text-white hover:bg-primary/90"
                                        href={`/host/events/${eventId}/settings/certificate`}
                                    >
                                        {t("certificateSettings.pageTitle")}
                                    </WrappedButton>
                                    <a
                                        href={
                                            isCertificatePublished || !hasIssuers
                                                ? undefined
                                                : `/host/events/${eventId}/imports/certificates`
                                        }
                                        onClick={
                                            isCertificatePublished || !hasIssuers
                                                ? (e) => e.preventDefault()
                                                : undefined
                                        }
                                        className={
                                            isCertificatePublished || !hasIssuers
                                                ? "pointer-events-none"
                                                : ""
                                        }
                                        title={
                                            !hasIssuers
                                                ? t(
                                                      "events.hostDetails.certificates.noIssuersMessage",
                                                  )
                                                : undefined
                                        }
                                    >
                                        <Button
                                            size="lg"
                                            variant="secondary-light"
                                            disabled={isCertificatePublished || !hasIssuers}
                                        >
                                            {t("events.hostDetails.actions.importReceivers")}
                                        </Button>
                                    </a>
                                </div>
                            </div>
                        )}
                    </StyledTabsContent>
                </StyledTabs>
            </SectionContainer>
        </div>
    );
}
