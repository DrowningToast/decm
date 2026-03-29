import { EventForm } from "@/components/forms/EventForm";
import type { EventFormData } from "@/lib/schemas/eventFormSchema";

import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import type { GetEventContractByEventIdData } from "@decm/api";
import { useEditEvent } from "./useEditEvent";
import { useDeleteEvent } from "./useDeleteEvent";
import type { EventViewModelExtended } from "@/services/EventService/EventService";
import type { PasswordPinModalSuccessResult } from "@/components/ui/password-pin-modal";
import { useWallet } from "@/hooks/useWallet";

interface EditEventPageProps {
    event: EventViewModelExtended;
    eventContract?: GetEventContractByEventIdData;
}
export const EditEventPage = ({ event, eventContract }: EditEventPageProps) => {
    const { editEventWithPassword, editEventWithSignature, isEditingEvent } = useEditEvent(
        event.id ?? "",
    );
    const { deleteEventWithPassword, deleteEventWithSignature, isDeletingEvent } = useDeleteEvent(
        event.id ?? "",
    );
    const { isConnected } = useWallet();

    const handleEditEvent = async (data: EventFormData, auth: PasswordPinModalSuccessResult) => {
        const commonFields = {
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
            ...(data.eventBanner ? { banner: data.eventBanner } : {}),
            ...(data.eventIcon ? { icon: data.eventIcon } : {}),
        };

        if (auth.type === "wallet") {
            await editEventWithSignature({
                ...commonFields,
                signature: auth.signature,
                signMessage: auth.signMessage,
            });
        } else {
            await editEventWithPassword({
                ...commonFields,
                hostPassword: auth.value,
            });
        }
    };

    const handleDeleteEvent = async (auth: PasswordPinModalSuccessResult) => {
        if (auth.type === "wallet") {
            await deleteEventWithSignature({
                signature: auth.signature,
                signMessage: auth.signMessage,
            });
        } else {
            await deleteEventWithPassword(auth.value);
        }
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
                    allowWalletSigning={isConnected}
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
