import { EventForm } from "@/components/forms/EventForm";
import type { EventFormData } from "@/lib/schemas/eventFormSchema";
import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import type {
    EventEventResponse,
    GetEventContractByEventIdData,
    UpdateEventPayload,
} from "@decm/api";
import { useEditEvent } from "./useEditEvent";
import { useDeleteEvent } from "./useDeleteEvent";

interface EditEventPageProps {
    event: EventEventResponse;
    eventContract?: GetEventContractByEventIdData;
}
export const EditEventPage = ({ event, eventContract }: EditEventPageProps) => {
    const { editEvent, isEditingEvent } = useEditEvent(event.id ?? "");
    const { deleteEvent, isDeletingEvent } = useDeleteEvent(event.id ?? "");

    const handleEditEvent = async (data: EventFormData, hostPassword: string) => {
        const req: UpdateEventPayload = {
            name: data.name,
            short_description: data.shortDescription,
            description: data.description ?? "",
            start_date: data.startDate.toISOString(),
            end_date: data.endDate.toISOString(),
            location: data.location,
            google_map_query: data.googleMapQuery,
            seats_count: data.seatsCount,
            contact_address: data.contactAddress,
            contact_number: data.contactNumber,
            host_password: hostPassword,
        };

        if (data.eventBanner) {
            req.banner = data.eventBanner;
        }
        if (data.eventIcon) {
            req.icon = data.eventIcon;
        }

        await editEvent(req);
    };

    const handleDeleteEvent = async () => {
        await deleteEvent();
    };

    return (
        <PageContainer title="Edit Event" className="space-y-6">
            {/* Page Header */}
            <SectionContainer>
                <TitleSubtitle
                    title="Edit Event"
                    subtitle="Fill in the details below to edit the event"
                />
            </SectionContainer>

            <SectionContainer>
                <EventForm
                    onSubmit={handleEditEvent}
                    onDelete={handleDeleteEvent}
                    mode="edit"
                    isLoading={isEditingEvent || isDeletingEvent}
                    defaultValues={{
                        contactAddress: event?.location ?? "",
                        contactNumber: event?.contact_number ?? "",
                        name: event?.title ?? "",
                        shortDescription: event?.short_description ?? "",
                        description: event?.long_description ?? "",
                        startDate: new Date(event?.start_date ?? ""),
                        endDate: new Date(event?.end_date ?? ""),
                        location: event?.location ?? "",
                        googleMapQuery: event?.google_map_query ?? "",
                        seatsCount: event?.max_attendees ?? 0,
                    }}
                    previewBannerUrl={event?.banner_presigned_url ?? ""}
                    previewIconUrl={event?.icon_presigned_url ?? ""}
                    eventContract={eventContract}
                />
            </SectionContainer>
        </PageContainer>
    );
};
