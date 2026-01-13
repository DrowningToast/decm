import { EventForm } from "@/components/forms/EventForm";
import type { EventFormData } from "@/lib/schemas/eventFormSchema";

import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import type { GetEventContractByEventIdData, UpdateEventPayload } from "@decm/api";
import { useEditEvent } from "./useEditEvent";
import { useDeleteEvent } from "./useDeleteEvent";
import type { EventViewModelExtended } from "@/services/EventService/EventService";

interface EditEventPageProps {
    event: EventViewModelExtended;
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

    const handleDeleteEvent = async (hostPassword: string) => {
        await deleteEvent(hostPassword);
    };

    return (
        <div className="flex flex-col gap-y-6">
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
                        contactNumber: event?.contactNumber ?? "",
                        name: event?.title ?? "",
                        shortDescription: event?.shortDescription ?? "",
                        description: event?.longDescription ?? "",
                        startDate: new Date(event?.startDate ?? ""),
                        endDate: new Date(event?.endDate ?? ""),
                        location: event?.location ?? "",
                        googleMapQuery: event?.googleMapQuery ?? "",
                        seatsCount: event?.maxAttendees ?? 0,
                    }}
                    previewBannerUrl={event?.bannerPresignedUrl ?? ""}
                    previewIconUrl={event?.iconPresignedUrl ?? ""}
                    eventContract={eventContract}
                />
            </SectionContainer>
        </div>
    );
};
