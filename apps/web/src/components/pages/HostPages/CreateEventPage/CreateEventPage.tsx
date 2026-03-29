import { useTranslation } from "react-i18next";
import { EventForm } from "@/components/forms/EventForm";
import type { EventFormData } from "@/lib/schemas/eventFormSchema";
import { toast } from "sonner";
import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import { useCreateEvent } from "./useCreateEvent";
import type { PasswordPinModalSuccessResult } from "@/components/ui/password-pin-modal";
import { useWallet } from "@/hooks/useWallet";

export const CreateEventPage = () => {
    const { t } = useTranslation();
    const { createEventWithPassword, createEventWithSignature, isCreatingEvent } = useCreateEvent();
    const { isConnected } = useWallet();

    const handleCreateEvent = async (data: EventFormData, auth: PasswordPinModalSuccessResult) => {
        if (!data.eventBanner || !data.eventIcon) {
            toast.error(t("errors.generic"));
            return;
        }

        const commonPayload = {
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

        if (auth.type === "wallet") {
            await createEventWithSignature({
                ...commonPayload,
                signature: auth.signature,
                signMessage: auth.signMessage,
            });
        } else {
            await createEventWithPassword({
                ...commonPayload,
                hostPassword: auth.value,
            });
        }
    };

    return (
        <div className="space-y-6">
            {/* Page Header */}
            <SectionContainer>
                <TitleSubtitle
                    title={t("host.events.createEventPage.title")}
                    subtitle={t("host.events.createEventPage.subtitle")}
                />
            </SectionContainer>

            <SectionContainer>
                <EventForm
                    onSubmit={handleCreateEvent}
                    mode="create"
                    isLoading={isCreatingEvent}
                    allowWalletSigning={isConnected}
                />
            </SectionContainer>
        </div>
    );
};
