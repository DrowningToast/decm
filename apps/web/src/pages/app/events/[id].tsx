import { useParams } from "@/router";
import { EventDetailPage } from "@/components/pages/Participant/Events/Detail/ParticipantEventDetailPage";
import { useViewEventDetailUsecase } from "@/components/pages/Participant/Events/Detail/useViewEventDetailUsecase";
import { useInviteStatusUsecase } from "@/components/pages/Participant/Events/Detail/useInviteStatusUsecase";
import { Typography } from "@/components/typography/typography";

const EventDetailRoute = () => {
    const { id } = useParams("/app/events/:id");
    const { event, isLoading, error } = useViewEventDetailUsecase({ eventId: id });
    const { inviteStatus } = useInviteStatusUsecase(id);

    if (isLoading) {
        return (
            <section className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="muted">
                    Loading event details...
                </Typography>
            </section>
        );
    }

    if (error || !event) {
        return (
            <section className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="destructive">
                    Failed to load event details
                </Typography>
            </section>
        );
    }

    // For invite-only events, check if user is invited
    if (event.accessType === "invite-only" && inviteStatus && !inviteStatus.isInvited) {
        return (
            <section className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="destructive">
                    You are not invited to this event
                </Typography>
            </section>
        );
    }

    return (
        <section className="relative z-10">
            <EventDetailPage eventId={id} />
        </section>
    );
};

export default EventDetailRoute;
