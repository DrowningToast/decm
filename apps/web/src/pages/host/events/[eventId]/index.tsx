import HostEventDetailsPage from "@/components/pages/HostPages/EventsPage/HostEventDetailsPage";
import { useParams } from "@/router";

export default function Page() {
  const { eventId } = useParams("/host/events/:eventId");

  return <HostEventDetailsPage eventId={eventId} />;
}
