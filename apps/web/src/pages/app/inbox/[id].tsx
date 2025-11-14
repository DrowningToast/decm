import { useTranslation } from "react-i18next";
import { useParams } from "@/router";
import { InboxDetailPage } from "@/components/pages/Participant/Inbox/InboxDetailPage";
import { useInboxDetailUsecase } from "@/components/pages/Participant/Inbox/useInboxDetailUsecase";
import { Typography } from "@/components/typography/typography";
import { PageContainer } from "@/components/container/PageContainer";

const InboxDetailRoute = () => {
    const { t } = useTranslation();
    const { id } = useParams("/app/inbox/:id");
    const { inboxDetail, isLoading, error } = useInboxDetailUsecase({ inboxId: id });

    if (isLoading) {
        return (
            <PageContainer className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="muted">
                    {t("participant.inbox.loadingDetails")}
                </Typography>
            </PageContainer>
        );
    }

    if (error || !inboxDetail) {
        return (
            <PageContainer className="relative z-10 w-full min-h-screen flex items-center justify-center">
                <Typography variant="text" tag="p" color="destructive">
                    {t("participant.inbox.failedToLoadDetails")}
                </Typography>
            </PageContainer>
        );
    }

    return (
        <PageContainer className="relative z-10">
            <InboxDetailPage inboxId={id} />
        </PageContainer>
    );
};

export default InboxDetailRoute;
