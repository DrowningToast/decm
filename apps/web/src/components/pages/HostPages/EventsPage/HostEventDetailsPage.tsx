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
import { useDataTable } from "@/hooks/use-data-table";
import { participantColumns, type Participant } from "./columns/participant-columns";
import type {
    EventconfigEventCertificateConfigResponse,
    EventconfigEventRegistrationConfigResponse,
    EventEventResponse,
} from "@decm/api";
import { toEventRegistrationConfigStatus } from "@/lib/events/event.utils";

interface HostEventDetailsPageProps {
    eventId: string;
    event: EventEventResponse;
    eventRegistrationConfig: EventconfigEventRegistrationConfigResponse;
    eventCertificateConfig?: EventconfigEventCertificateConfigResponse;
}

// Mock API function - replace with actual API call
const mockFetchParticipants = async ({
    page,
    pageSize,
    search,
    sortBy,
    sortOrder,
}: {
    page: number;
    pageSize: number;
    search: string;
    sortBy?: string;
    sortOrder?: "asc" | "desc";
}): Promise<{ data: Participant[]; total: number }> => {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    // Mock data - replace with actual API call
    const allParticipants: Participant[] = Array.from({ length: 45 }, (_, i) => ({
        id: `participant-${i + 1}`,
        name: `${["John", "Jane", "Bob", "Alice", "Charlie", "David", "Emma", "Frank"][i % 8]} ${
            ["Doe", "Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller"][i % 8]
        }`,
        email: `user${i + 1}@email.com`,
        phoneNumber: `+1 ${Math.floor(100 + Math.random() * 900)} ${Math.floor(100 + Math.random() * 900)} ${Math.floor(1000 + Math.random() * 9000)}`,
        walletAddress: `0x${Math.random().toString(16).substring(2, 6)}...${Math.random().toString(16).substring(2, 6)}`,
        status: ["confirmed", "pending", "rejected"][i % 3] as "confirmed" | "pending" | "rejected",
    }));

    // Filter by search
    let filteredData = allParticipants;
    if (search) {
        filteredData = allParticipants.filter(
            (p) =>
                p.name.toLowerCase().includes(search.toLowerCase()) ||
                p.email.toLowerCase().includes(search.toLowerCase()),
        );
    }

    // Sort
    if (sortBy) {
        filteredData.sort((a, b) => {
            const aValue = a[sortBy as keyof Participant];
            const bValue = b[sortBy as keyof Participant];

            if (aValue < bValue) return sortOrder === "asc" ? -1 : 1;
            if (aValue > bValue) return sortOrder === "asc" ? 1 : -1;
            return 0;
        });
    }

    // Paginate
    const startIndex = (page - 1) * pageSize;
    const paginatedData = filteredData.slice(startIndex, startIndex + pageSize);

    return {
        data: paginatedData,
        total: filteredData.length,
    };
};

export default function HostEventDetailsPage({
    eventId,
    event,
    eventRegistrationConfig,
}: HostEventDetailsPageProps) {
    const { t } = useTranslation();

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

    // Use the data table hook
    const participantsTable = useDataTable<Participant>({
        fetchData: mockFetchParticipants,
        initialPageSize: 10,
    });

    const certificatesTable = useDataTable<Participant>({
        fetchData: mockFetchParticipants,
        initialPageSize: 10,
    });

    return (
        <PageContainer title="Events Details">
            <SectionContainer>
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <div className="w-8 h-8 bg-white rounded-full"></div>
                        <div className="flex items-center gap-2">
                            <Typography tag="p" size={"subheader"}>
                                {event.title}
                            </Typography>
                            {event.is_verified && <CheckCircle2Icon color="#eb5331" />}
                        </div>
                    </div>

                    <WrappedButton href={`/host/events/${eventId}/edit`}>Edit Event</WrappedButton>
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
                    <TextLabelValue label="Status" value="NA" />
                    <TextLabelValue label="Final call for request" value={"NA"} />
                    <TextLabelValue
                        label="Participation request"
                        value={event.is_booking_request_required ? "Required" : "Not Required"}
                    />
                    <TextLabelValue label="Seats count" value={`${0} / ${event.max_attendees}`} />
                    <TextLabelValue
                        label="Event Contract Address"
                        value="NA"
                        endIcon={<ExternalLinkIcon size={16} />}
                        valueClassName="cursor-pointer underline"
                        href="https://www.etherscan.io/address/0x0000000000000000000000000000000000000000"
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

                                <div className="flex items-center justify-end">
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
                                columns={participantColumns}
                                data={participantsTable.data}
                                totalItems={participantsTable.totalItems}
                                currentPage={participantsTable.currentPage}
                                pageSize={participantsTable.pageSize}
                                onPageChange={participantsTable.setCurrentPage}
                                onPageSizeChange={participantsTable.setPageSize}
                                searchValue={participantsTable.searchValue}
                                onSearchChange={participantsTable.setSearchValue}
                                searchPlaceholder="Search participants..."
                                sorting={participantsTable.sorting}
                                onSortingChange={participantsTable.setSorting}
                                isLoading={participantsTable.isLoading}
                            />
                        </div>
                    </StyledTabsContent>
                    <StyledTabsContent value="certificates">
                        <div className="w-full bg-primary/10 border border-primary/20 rounded-lg p-6 flex flex-col items-center justify-center mb-6">
                            <p className="font-semibold text-lg text-primary">
                                Add event&apos;s certificate configuration
                            </p>
                            <p className="text-muted-foreground text-base mt-1 text-center max-w-xl">
                                Set up the certificate template and rules for this event.
                                Participants will receive certificates based on your configuration.
                            </p>
                            <WrappedButton
                                className="mt-4 px-5 py-2 rounded-md bg-primary text-white font-medium hover:bg-primary/90 transition"
                                href={`/host/events/${eventId}/settings/certificate`}
                            >
                                Certificates Settings
                            </WrappedButton>
                        </div>

                        <DataTable
                            columns={participantColumns}
                            data={certificatesTable.data}
                            totalItems={certificatesTable.totalItems}
                            currentPage={certificatesTable.currentPage}
                            pageSize={certificatesTable.pageSize}
                            onPageChange={certificatesTable.setCurrentPage}
                            onPageSizeChange={certificatesTable.setPageSize}
                            searchValue={certificatesTable.searchValue}
                            onSearchChange={certificatesTable.setSearchValue}
                            searchPlaceholder="Search certificates..."
                            sorting={certificatesTable.sorting}
                            onSortingChange={certificatesTable.setSorting}
                            isLoading={certificatesTable.isLoading}
                        />
                    </StyledTabsContent>
                </StyledTabs>
            </SectionContainer>
        </PageContainer>
    );
}
