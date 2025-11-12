import React from "react";
import { useTranslation } from "react-i18next";
import {
    Table,
    TableHeader,
    TableRow,
    TableHead,
    TableBody,
    TableCell,
} from "@/components/ui/table";
import { EventStatusBadge } from "./EventStatusBadge";
import { EventActions } from "./EventActions";
import type { GetIssuerEventsData } from "@decm/api";

interface EventTableProps {
    events: GetIssuerEventsData;
    type: "pending" | "signed";
    onActionClick?: (eventId: string) => void;
}

export const EventTable: React.FC<EventTableProps> = ({ events, type, onActionClick }) => {
    const { t } = useTranslation();

    return (
        <div className="border border-primary rounded-lg overflow-hidden">
            <Table className="min-w-full">
                <TableHeader className="bg-[#0a0a0a]">
                    <TableRow>
                        <TableHead className="px-6 py-3 text-left text-xs font-medium text-secondary uppercase tracking-wider">
                            {t("issuer.sign.table.eventName")}
                        </TableHead>
                        <TableHead className="px-6 py-3 text-left text-xs font-medium text-secondary uppercase tracking-wider">
                            {t("issuer.sign.table.host")}
                        </TableHead>
                        <TableHead className="px-6 py-3 text-left text-xs font-medium text-secondary uppercase tracking-wider">
                            {t("issuer.sign.table.certificates")}
                        </TableHead>
                        <TableHead className="px-6 py-3 text-left text-xs font-medium text-secondary uppercase tracking-wider">
                            {t("issuer.sign.table.status")}
                        </TableHead>
                        <TableHead className="px-6 py-3 text-left text-xs font-medium text-secondary uppercase tracking-wider">
                            {t("issuer.sign.table.action")}
                        </TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody className=" divide-y divide-primary">
                    {events.map((event) => (
                        <TableRow key={event.id}>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-primary">
                                {event.event_title}
                            </TableCell>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-primary">
                                {/* TODO: Get host name from profile API */}
                                {t("issuer.sign.table.hostPlaceholder")}
                            </TableCell>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-primary">
                                {/* TODO: Get certificate count from API */}
                                {type === "pending" ? 160 : 45}
                            </TableCell>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm">
                                <EventStatusBadge status={type} />
                            </TableCell>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm">
                                <EventActions
                                    type={type}
                                    eventId={event.event_id || ""}
                                    onActionClick={onActionClick}
                                />
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    );
};
