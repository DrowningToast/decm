import { useParams } from "@/router";
import { EventDetailPage } from "@/components/pages/Participant/Events/Detail/ParticipantEventDetailPage";
import { useEventViewModelUsecase } from "@/components/pages/Participant/Events/Detail/useEventViewModelUsecase";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";
import { PageContainer } from "@/components/container/PageContainer";

const EventDetailRoute = () => {
    const { t } = useTranslation();
    const { id } = useParams("/app/events/:id");
    const { event, isLoading, error } = useEventViewModelUsecase({ eventId: id });

    if (isLoading) {
        return (
            <PageContainer className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="muted">
                    {t("participant.events.detail.loadingDetails")}
                </Typography>
            </PageContainer>
        );
    }

    if (error || !event) {
        // Determine error type
        const isNotFound =
            error &&
            "response" in error &&
            error.response &&
            typeof error.response === "object" &&
            "status" in error.response &&
            error.response.status === 404;

        const errorMessage = isNotFound
            ? t("participant.events.detail.notFoundError")
            : t("participant.events.detail.loadError");

        return (
            <PageContainer className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="destructive">
                    {errorMessage}
                </Typography>
            </PageContainer>
        );
    }

    return (
        <PageContainer className="relative z-10">
            <EventDetailPage eventId={id} />
        </PageContainer>
    );
};

export default EventDetailRoute;
