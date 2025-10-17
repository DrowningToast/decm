import { useTranslation } from "react-i18next";
import { EventForm } from "@/components/forms/EventForm";
import type { EventFormData } from "@/lib/schemas/eventFormSchema";
import { toast } from "sonner";
import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import { useCreateEvent } from "./useCreateEvent";
import type { CreateEventPayload } from "@decm/api";

export const CreateEventPage = () => {
    const { t } = useTranslation();
    const { createEvent, isCreatingEvent } = useCreateEvent();

    const handleCreateEvent = async (data: EventFormData) => {
        if (!data.eventBanner || !data.eventIcon) {
            toast.error(t("errors.generic"));
            return;
        }

        const req: CreateEventPayload = {
            banner: data.eventBanner,
            contact_address: data.contactAddress,
            contact_number: data.contactNumber,
            description: data.description ?? "",
            end_date: data.endDate.toISOString(),
            start_date: data.startDate.toISOString(),
            seats_count: data.seatsCount,
            short_description: data.shortDescription,
            google_map_query: data.googleMapQuery,
            icon: data.eventIcon,
            location: data.location,
            name: data.name,
        };

        await createEvent(req);
    };

    return (
        <PageContainer title="Create Event" className="space-y-6">
            {/* Page Header */}
            <SectionContainer>
                <TitleSubtitle
                    title="Create Event"
                    subtitle="Fill in the details below to create a new event"
                />
            </SectionContainer>

            <SectionContainer>
                <EventForm onSubmit={handleCreateEvent} mode="create" isLoading={isCreatingEvent} />
            </SectionContainer>
        </PageContainer>
    );
};
