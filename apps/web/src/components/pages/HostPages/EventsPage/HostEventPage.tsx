import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import WrappedButton from "@/components/wrapper/WrappedButton";
import { Link } from "@/router";
import { DataTablePagination } from "@/components/ui/pagination/Pagination";
import { useHostEvents } from "./useHostEvents";
import type { EventEventResponse } from "@decm/api";

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
        <PageContainer title="Events">
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
                                    events.flat().map((event: EventEventResponse) => (
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
                                                    {event.start_date && event.end_date
                                                        ? `${new Date(event.start_date).toLocaleDateString()} - ${new Date(event.end_date).toLocaleDateString()}`
                                                        : "Date TBD"}
                                                </p>
                                                <p className="text-muted text-sm">
                                                    {event.is_public ? "Public" : "Invite Only"}
                                                </p>
                                            </td>
                                            <td className="text-end">
                                                {event.is_verified ? "Verified" : "Active"}
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
        </PageContainer>
    );
}
