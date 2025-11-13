import { ParticipantEventListPage } from "@/components/pages/Participant/Events/List/ParticipantEventListPage";
import { PageContainer } from "@/components/container/PageContainer";

const EventsPage = () => {
    return (
        <PageContainer className="relative z-10">
            <ParticipantEventListPage />
        </PageContainer>
    );
};

export default EventsPage;
