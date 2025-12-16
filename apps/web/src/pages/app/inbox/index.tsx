import { InboxListPage } from "@/components/pages/Participant/Inbox/InboxListPage";
import { PageContainer } from "@/components/container/PageContainer";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useEffect } from "react";

const InboxPage = () => {
    // Invalidate profile cache when entering the inbox page
    // This ensures the unread message count is refreshed
    useEffect(() => {
        queryClient.invalidateQueries({ queryKey: QUERY_KEY.inbox.list() });
    }, []);

    return (
        <PageContainer className="relative z-10">
            <InboxListPage />
        </PageContainer>
    );
};

export default InboxPage;
