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
import { issuerColumns } from "./columns/issuer-columns";
import type {
    EntityEventCertificate,
    GetEventCertificateConfigData,
    GetEventContractByEventIdData,
    GetEventIssuersByEventIdData,
} from "@decm/api";
import { formatEthereumAddress } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { useState, useMemo, useCallback, useEffect } from "react";
import type { SortingState } from "@tanstack/react-table";
import { useParticipantColumns, type Participant } from "./columns/useParticipantColumns";
import { CertificateColumns } from "./columns/CertificateColumns";
import { Separator } from "@/components/ui/separator";
import { useEventCertificates } from "@/hooks/useEventCertificates";
import { useRevokeEventCertificate } from "@/hooks/events/useRevokeEventCertificate";
import { Link } from "@/router";
import type {
    EventRegistrationConfiguration,
    EventRegistrationInvitation,
} from "@/services/EventRegistration/EventRegistration";
import { EventStatusesViewModel, EventTypesViewModel } from "./ViewModel";
import type { EventViewModel } from "@/services/EventService/EventService";

interface HostEventDetailsPageProps {
    eventId: string;
    event: EventViewModel;
    eventRegistrationConfig: EventRegistrationConfiguration;
    eventCertificateConfig?: GetEventCertificateConfigData;
    eventIssuers?: GetEventIssuersByEventIdData;
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
                firstName: invitation.firstName || "",
                lastName: invitation.lastName || "",
                email: invitation.email || "",
                phoneNumber: invitation.phoneNumber || "",
                academicInstitution: invitation.academicInstitution || "",
                walletAddress: "", // Not available in invitation data
                status: invitation.cancelledAt ? "rejected" : "pending",
            }),
        );
    }, [processedData, currentPage, pageSize]);

    const participantColumns = useParticipantColumns();

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
    const allIssuersSigned = eventIssuers?.every((issuer) => issuer.is_signed === 1) ?? false;

    // Fetch event certificates using the hook
    const { certificates: eventCertificates, isLoading: certificatesLoading } =
        useEventCertificates(eventId);

    console.log(event.googleMapQuery);

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
                                        event.isPublic
                                            ? t("participantSettings.eventTypePublic")
                                            : t("participantSettings.eventTypePrivate")
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
                                            className="px-5 py-2 rounded-md bg-primary text-white font-medium hover:bg-primary/90 transition"
                                            href={`/host/events/${eventId}/settings/certificate`}
                                        >
                                            {t("certificateSettings.pageTitle")}
                                        </WrappedButton>
                                    </div>
                                </div>

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

                                        <Separator className="my-12" />

                                        <div className="flex items-center justify-between">
                                            {/* Publish Certificates Section */}
                                            <div className="w-full bg-white border border-white/50 rounded-lg p-6 flex flex-row items-center justify-between">
                                                <div>
                                                    <Typography
                                                        variant="text"
                                                        tag="p"
                                                        className="font-semibold text-lg text-black"
                                                    >
                                                        {t(
                                                            "events.hostDetails.certificates.publishTitle",
                                                        )}
                                                    </Typography>
                                                    <Typography
                                                        variant="text"
                                                        tag="p"
                                                        className="text-black/50 text-base mt-1"
                                                    >
                                                        {t(
                                                            "events.hostDetails.certificates.publishDescription",
                                                        )}
                                                    </Typography>
                                                </div>
                                                <div className="flex gap-2">
                                                    <WrappedButton
                                                        className="px-4 py-2 rounded-md bg-green-600 text-white font-medium hover:bg-green-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed"
                                                        disabled={!allIssuersSigned}
                                                        onClick={() => {
                                                            // TODO: Implement publish certificates functionality
                                                            console.log(
                                                                "Publish certificates for event:",
                                                                eventId,
                                                            );
                                                        }}
                                                    >
                                                        {t(
                                                            allIssuersSigned
                                                                ? "events.hostDetails.actions.publishCertificates"
                                                                : "events.hostDetails.actions.waitingForSignatures",
                                                        )}
                                                    </WrappedButton>
                                                </div>
                                            </div>
                                        </div>

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

                                        <DataTable
                                            columns={CertificateColumns(
                                                (eventCertificateId: string) => {
                                                    revokeEventCertificate({
                                                        certificateIds: [eventCertificateId],
                                                        eventId,
                                                    });
                                                },
                                            )}
                                            data={
                                                eventCertificates?.map((cert) => {
                                                    const firstName =
                                                        cert.name?.split(" ")[0] || "";
                                                    const lastName =
                                                        cert.name?.split(" ").slice(1).join(" ") ||
                                                        "";

                                                    return {
                                                        ...cert,
                                                        firstName,
                                                        lastName,
                                                        email: cert.receiver_email || "",
                                                        academicInstitution:
                                                            cert.academic_institution || "",
                                                        issuedAt: cert.created_at,
                                                        status: cert.revoked_at
                                                            ? "rejected"
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
                                            searchPlaceholder={t(
                                                "events.hostDetails.certificates.searchCertificatesPlaceholder",
                                            )}
                                            sorting={[]}
                                            onSortingChange={() => {}}
                                            isLoading={certificatesLoading}
                                            disablePagination
                                        />

                                        <a
                                            href={`/host/events/${eventId}/imports/certificates`}
                                            className="flex items-center justify-center p-6 border-dashed rounded-xl border-2 border-white/50 gap-4 cursor-pointer"
                                        >
                                            <CloudUploadIcon />{" "}
                                            <Typography variant="text" tag="span">
                                                {t(
                                                    "events.hostDetails.certificates.importSectionTitle",
                                                )}
                                            </Typography>
                                        </a>
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
                                <div className="flex gap-4 mt-6">
                                    <WrappedButton
                                        className="px-5 py-2 rounded-md bg-primary text-white font-medium hover:bg-primary/90 transition"
                                        href={`/host/events/${eventId}/settings/certificate`}
                                    >
                                        {t("certificateSettings.pageTitle")}
                                    </WrappedButton>
                                    <a href={`/host/events/${eventId}/imports/certificates`}>
                                        <Button size="xl" variant="secondary-light">
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
