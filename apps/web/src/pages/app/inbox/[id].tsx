import { useTranslation } from "react-i18next";
import { useParams } from "@/router";
import { InboxDetailPage } from "@/components/pages/Participant/Inbox/InboxDetailPage";
import { useInboxDetailUsecase } from "@/components/pages/Participant/Inbox/useInboxDetailUsecase";
import { Typography } from "@/components/typography/typography";

const InboxDetailRoute = () => {
    const { t } = useTranslation();
    const { id } = useParams("/app/inbox/:id");
    const { inboxDetail, isLoading, error } = useInboxDetailUsecase({ inboxId: id });

    if (isLoading) {
        return (
            <section className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="muted">
                    {t("participant.inbox.loadingDetails")}
                </Typography>
            </section>
        );
    }

    if (error || !inboxDetail) {
        return (
            <section className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="destructive">
                    {t("participant.inbox.failedToLoadDetails")}
                </Typography>
            </section>
        );
    }

    return (
        <section className="relative z-10">
            <InboxDetailPage inboxId={id} />
        </section>
    );
};

export default InboxDetailRoute;
