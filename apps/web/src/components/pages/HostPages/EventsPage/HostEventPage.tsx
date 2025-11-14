import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import WrappedButton from "@/components/wrapper/WrappedButton";
import { Link } from "@/router";
import { DataTablePagination } from "@/components/ui/pagination/Pagination";
import { useHostEvents } from "./useHostEvents";
import type { EventStatus, EventType } from "@/services/EventService/EventService";

const EventTypes: Record<EventType, string> = {
    private: "Private",
    invite: "Invite",
};

const EventStatuses: Record<EventStatus, string> = {
    active: "Active",
    inactive: "Inactive",
    closed: "Closed",
} as const;

export default function HostEventPage() {
    const {
        events,
        isLoadingEvents,
        isLoadingEventsError,
        page,
        rowsPerPage,
        hasPreviousPage,
        hasNextPage,
        handlePageChange,
        handleRowsPerPageChange,
    } = useHostEvents();

    return (
        <div title="Events">
            <SectionContainer className="flex items-center justify-between">
                <TitleSubtitle title="Events" subtitle="Create or manage your events" />
                <div className="flex justify-end">
                    <WrappedButton href="/host/events/create">Create Event</WrappedButton>
                </div>
            </SectionContainer>

            <SectionContainer>
                {isLoadingEvents ? (
                    <div className="flex justify-center py-8">
                        <div>Loading events...</div>
                    </div>
                ) : isLoadingEventsError ? (
                    <div className="flex justify-center py-8 text-red-500">
                        <div>Error loading events: {isLoadingEventsError.message}</div>
                    </div>
                ) : (
                    <>
                        <table className="w-full">
                            <thead className="border-b h-10">
                                <tr>
                                    <th className="text-start text-muted">Name</th>
                                    <th className="text-end text-muted">Status</th>
                                </tr>
                            </thead>
                            <tbody>
                                {events && events.length > 0 ? (
                                    events.flat().map((event) => (
                                        <tr key={event.id}>
                                            <td>
                                                <Link
                                                    to={"/host/events/:eventId"}
                                                    params={{
                                                        eventId: event.id || "",
                                                    }}
                                                >
                                                    <p className="underline mt-4 cursor-pointer">
                                                        {event.title}
                                                    </p>
                                                </Link>
                                                <p className="text-muted text-sm mt-0.5">
                                                    {event.startDate && event.endDate
                                                        ? `${new Date(event.startDate).toLocaleDateString()} - ${new Date(event.endDate).toLocaleDateString()}`
                                                        : "Date TBD"}
                                                </p>
                                                <p className="text-muted text-sm">
                                                    {EventTypes[event.eventType]}
                                                </p>
                                            </td>
                                            <td className="text-end">
                                                {EventStatuses[event.eventStatus]}
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={2} className="text-center py-8">
                                            No events found. Create your first event!
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>

                        {events && events.length > 0 && (
                            <div className="mt-6">
                                <DataTablePagination
                                    currentPage={page}
                                    hasNextPage={hasNextPage}
                                    hasPreviousPage={hasPreviousPage}
                                    onPageChange={handlePageChange}
                                    onRowsPerPageChange={handleRowsPerPageChange}
                                    isLoading={isLoadingEvents}
                                    rowsPerPage={rowsPerPage}
                                />
                            </div>
                        )}
                    </>
                )}
            </SectionContainer>
        </div>
    );
}
