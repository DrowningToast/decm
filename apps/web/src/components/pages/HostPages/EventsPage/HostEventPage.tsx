import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import WrappedButton from "@/components/wrapper/WrappedButton";
import { Link } from "@/router";
import { DataTablePagination } from "@/components/ui/pagination/Pagination";
import { useHostEvents } from "./useHostEvents";
import { EventStatusesViewModel, EventTypesViewModel } from "./ViewModel";
import { useTranslation } from "react-i18next";

export default function HostEventPage() {
    const { t } = useTranslation();
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
        <div title={t("host.events.title")}>
            <SectionContainer className="flex items-center justify-between">
                <TitleSubtitle
                    title={t("host.events.title")}
                    subtitle={t("host.events.subtitle")}
                />
                <div className="flex justify-end">
                    <WrappedButton href="/host/events/create">
                        {t("host.events.createEvent")}
                    </WrappedButton>
                </div>
            </SectionContainer>

            <SectionContainer>
                {isLoadingEvents ? (
                    <div className="flex justify-center py-8">
                        <div>{t("host.events.loading")}</div>
                    </div>
                ) : isLoadingEventsError ? (
                    <div className="flex justify-center py-8 text-red-500">
                        <div>
                            {t("host.events.loadingError", {
                                error: isLoadingEventsError.message,
                            })}
                        </div>
                    </div>
                ) : (
                    <>
                        <table className="w-full">
                            <thead className="border-b h-10">
                                <tr>
                                    <th className="text-start text-muted">
                                        {t("host.events.table.name")}
                                    </th>
                                    <th className="text-end text-muted">
                                        {t("host.events.table.status")}
                                    </th>
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
                                                        : t("host.events.dateTBD")}
                                                </p>
                                                <p className="text-muted text-sm">
                                                    {EventTypesViewModel[event.eventType]}
                                                </p>
                                            </td>
                                            <td className="text-end">
                                                {EventStatusesViewModel[event.eventStatus]}
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={2} className="text-center py-8">
                                            {t("host.events.empty")}
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
