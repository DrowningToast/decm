import { InboxListPage } from "@/components/pages/Participant/Inbox/InboxListPage";
import { PageContainer } from "@/components/container/PageContainer";

const InboxPage = () => {
    return (
        <PageContainer className="relative z-10">
            <InboxListPage />
        </PageContainer>
    );
};

export default InboxPage;
