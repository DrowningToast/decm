import PageContainer from "@/components/container/PageContainer";
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
import { CheckCircle2Icon, ExternalLinkIcon } from "lucide-react";
import { useTranslation } from "react-i18next";
import { DataTable } from "@/components/ui/data-table";
import { issuerColumns } from "./columns/issuer-columns";
import type {
    EventconfigEventRegistrationConfigResponse,
    EventEventResponse,
    GetEventCertificateConfigData,
    GetEventContractByEventIdData,
    GetEventIssuersByEventIdData,
    GetEventRegistrationInvitationsByEventIdData,
} from "@decm/api";
import { toEventRegistrationConfigStatus } from "@/lib/events/event.utils";
import { formatEthereumAddress } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { useState, useMemo, useCallback, useEffect } from "react";
import type { SortingState } from "@tanstack/react-table";
import { ParticipantColumns, type Participant } from "./columns/ParticipantColumns";

interface HostEventDetailsPageProps {
    eventId: string;
    event: EventEventResponse;
    eventRegistrationConfig: EventconfigEventRegistrationConfigResponse;
    eventCertificateConfig?: GetEventCertificateConfigData;
    eventIssuers?: GetEventIssuersByEventIdData;
    eventContract?: GetEventContractByEventIdData;
    eventInvitations?: GetEventRegistrationInvitationsByEventIdData;
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
        if (!eventInvitations) return [];

        let filtered = [...eventInvitations];

        // Filter by search value
        if (debouncedSearch) {
            const searchLower = debouncedSearch.toLowerCase();
            filtered = filtered.filter(
                (item) =>
                    item.first_name?.toLowerCase().includes(searchLower) ||
                    item.last_name?.toLowerCase().includes(searchLower) ||
                    item.email?.toLowerCase().includes(searchLower) ||
                    item.academic_institution?.toLowerCase().includes(searchLower) ||
                    item.phone_number?.toLowerCase().includes(searchLower),
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

    // Calculate pagination
    const paginatedData = useMemo(() => {
        const startIndex = (currentPage - 1) * pageSize;
        const endIndex = startIndex + pageSize;
        return processedData.slice(startIndex, endIndex);
    }, [processedData, currentPage, pageSize]);

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

    // Mock participant requirements data - replace with actual data
    // const participantRequirements: Record<string, RequirementStatus> = {
    //     firstName: "required",
    //     lastName: "required",
    //     email: "required",
    //     bio: "optional",
    //     phoneNumber: "optional",
    //     address: "not_required",
    //     academicInstitution: "required",
    //     academicEmail: "required",
    // };

    return (
        <PageContainer title="Events Details">
            <SectionContainer>
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <img
                            src={event.icon_presigned_url ?? ""}
                            alt={event.title}
                            className="w-12 h-12 rounded-full object-cover"
                        />

                        <div className="flex items-center gap-2">
                            <Typography tag="p" size={"subheader"}>
                                {event.title}
                            </Typography>
                            {event.is_verified && <CheckCircle2Icon color="#eb5331" />}
                        </div>
                    </div>

                    <div className="flex items-center gap-4">
                        <WrappedButton variant={"secondaryWhite"}>Confirm Event</WrappedButton>
                        <WrappedButton href={`/host/events/${eventId}/edit`}>
                            Edit Event
                        </WrappedButton>
                    </div>
                </div>
            </SectionContainer>

            <SectionContainer className="lg:grid lg:grid-cols-4 gap-8">
                <div className="flex flex-col gap-4 lg:col-span-3">
                    <Typography tag="p" size={"base"} color="muted">
                        {event.short_description}
                    </Typography>

                    <Typography tag="p" size={"base"} color="muted">
                        {event.long_description}
                    </Typography>

                    <img
                        src={event.banner_presigned_url ?? ""}
                        alt={event.title}
                        className="w-full h-[250px] object-cover rounded-lg"
                    />
                </div>

                <div className="flex flex-col gap-4 mt-6 lg:mt-0">
                    <TextLabelValue
                        label="Status"
                        value={event.event_status?.toUpperCase() ?? "NA"}
                    />
                    <TextLabelValue label="Final call for request" value={"NA"} />
                    <TextLabelValue
                        label="Participation request"
                        value={event.is_booking_request_required ? "Required" : "Not Required"}
                    />
                    <TextLabelValue label="Seats count" value={`${0} / ${event.max_attendees}`} />
                    <TextLabelValue
                        label="Event Contract Address"
                        value={
                            eventContract?.event_contract_address
                                ? formatEthereumAddress(eventContract.event_contract_address)
                                : "NA"
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
                        <StyledTabsTrigger value="event-info">Event Info</StyledTabsTrigger>
                        <StyledTabsTrigger value="participants">Participants</StyledTabsTrigger>
                        <StyledTabsTrigger value="certificates">Certificates</StyledTabsTrigger>
                    </StyledTabsList>
                    <StyledTabsContent value="event-info">
                        <div className="flex flex-col gap-4 lg:flex-row">
                            <div className="flex flex-col gap-4 flex-1">
                                <TextLabelValue
                                    label="Venue Location"
                                    value={event.location ?? ""}
                                />
                                <TextLabelValue
                                    label="Google Map Search"
                                    value={event.google_map_query ?? ""}
                                />

                                <TextLabelValue
                                    label="Contact Address"
                                    value={event.contact_number ?? ""}
                                />
                                <TextLabelValue
                                    label="Contact"
                                    value={event.contact_number ?? ""}
                                />
                            </div>
                            <div className="flex-1">
                                <GoogleMapsEmbed query={event.google_map_query ?? ""} />
                            </div>
                        </div>
                    </StyledTabsContent>
                    <StyledTabsContent value="participants">
                        <div className="space-y-4">
                            {/* Event Settings Summary */}
                            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 my-6">
                                <TextLabelValue
                                    label={t("events.settings.eventType")}
                                    value={event.is_public ? "Public" : "Private"}
                                />
                                <TextLabelValue
                                    label={t("events.settings.bookingRequired")}
                                    value={
                                        event.is_booking_request_required
                                            ? t("common.yes")
                                            : t("common.no")
                                    }
                                />
                                <TextLabelValue
                                    label={t("events.settings.tokenTransferable")}
                                    value={
                                        event.is_ticket_transferable
                                            ? t("common.yes")
                                            : t("common.no")
                                    }
                                />

                                <div className="flex items-center justify-end gap-4">
                                    <Button variant="secondary-dark" className="h-full">
                                        <a href={`/host/events/${eventId}/imports`}>
                                            Import Participants
                                        </a>
                                    </Button>
                                    <WrappedButton
                                        href={`/host/events/${eventId}/settings/participant`}
                                    >
                                        {t("events.settings.participantSettings")}
                                    </WrappedButton>
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
                                                status={toEventRegistrationConfigStatus(
                                                    eventRegistrationConfig.first_name_requirement_status,
                                                )}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.lastName")}
                                                status={toEventRegistrationConfigStatus(
                                                    eventRegistrationConfig.last_name_requirement_status,
                                                )}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.email")}
                                                status={toEventRegistrationConfigStatus(
                                                    eventRegistrationConfig.email_requirement_status,
                                                )}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.bio")}
                                                status={toEventRegistrationConfigStatus(
                                                    eventRegistrationConfig.bio_requirement_status,
                                                )}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.phoneNumber")}
                                                status={toEventRegistrationConfigStatus(
                                                    eventRegistrationConfig.phone_number_requirement_status,
                                                )}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.address")}
                                                status={toEventRegistrationConfigStatus(
                                                    eventRegistrationConfig.address_requirement_status,
                                                )}
                                            />
                                            <RequirementItem
                                                label={t(
                                                    "events.participants.fields.academicInstitution",
                                                )}
                                                status={toEventRegistrationConfigStatus(
                                                    eventRegistrationConfig.academic_institution_requirement_status,
                                                )}
                                            />
                                            <RequirementItem
                                                label={t(
                                                    "events.participants.fields.academicEmail",
                                                )}
                                                status={toEventRegistrationConfigStatus(
                                                    eventRegistrationConfig.academic_email_requirement_status,
                                                )}
                                            />
                                        </div>
                                    </AccordionContent>
                                </AccordionItem>
                            </Accordion>

                            <DataTable
                                columns={ParticipantColumns()}
                                data={paginatedData as Participant[]}
                                totalItems={processedData.length}
                                currentPage={currentPage}
                                pageSize={pageSize}
                                onPageChange={handlePageChange}
                                onPageSizeChange={handlePageSizeChange}
                                searchValue={searchValue}
                                onSearchChange={handleSearchChange}
                                searchPlaceholder="Search participants..."
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
                                <div className="w-full bg-primary/10 border border-primary/20 rounded-lg p-6 flex flex-row items-center justify-between">
                                    <div>
                                        <p className="font-semibold text-lg text-primary">
                                            Certificate Settings
                                        </p>
                                        <p className="text-muted-foreground text-base mt-1">
                                            Certificate template and rules are configured for this
                                            event. Manage issuers and publish certificates when
                                            ready.
                                        </p>
                                    </div>
                                    <WrappedButton
                                        className="px-5 py-2 rounded-md bg-primary text-white font-medium hover:bg-primary/90 transition"
                                        href={`/host/events/${eventId}/settings/certificate`}
                                    >
                                        Certificate Settings
                                    </WrappedButton>
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
                                                    Event Issuers ({eventIssuers.length})
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
                                                searchPlaceholder="Search issuers..."
                                                sorting={[]}
                                                onSortingChange={() => {}}
                                                isLoading={false}
                                                disablePagination
                                            />
                                        </div>

                                        <div className="flex items-center justify-between">
                                            <Typography
                                                variant="text"
                                                tag="h3"
                                                className="text-lg font-semibold"
                                            >
                                                Event Certificates ({eventIssuers.length})
                                            </Typography>
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
                                                {allIssuersSigned
                                                    ? "Publish Certificates"
                                                    : "Waiting for All Signatures"}
                                            </WrappedButton>
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
                                            No issuers configured for this event. Please add issuers
                                            in certificate settings.
                                        </Typography>
                                    </div>
                                )}
                            </div>
                        ) : (
                            <div className="w-full bg-primary/10 border border-primary/20 rounded-lg p-6 flex flex-col items-center justify-center">
                                <p className="font-semibold text-lg text-primary">
                                    Add event's certificate configuration
                                </p>
                                <p className="text-muted-foreground text-base mt-1 text-center max-w-xl">
                                    Set up a certificate template and rules for this event.
                                    Participants will receive certificates based on your
                                    configuration.
                                </p>
                                <WrappedButton
                                    className="mt-4 px-5 py-2 rounded-md bg-primary text-white font-medium hover:bg-primary/90 transition"
                                    href={`/host/events/${eventId}/settings/certificate`}
                                >
                                    Certificates Settings
                                </WrappedButton>
                            </div>
                        )}
                    </StyledTabsContent>
                </StyledTabs>
            </SectionContainer>
        </PageContainer>
    );
}
