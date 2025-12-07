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
import { Typography } from "@/components/typography/typography";
import { EventStatusBadge } from "./EventStatusBadge";
import { EventActions } from "./EventActions";
import type { IssuerEventViewModel } from "@/services/IssuerService/IssuerService";

interface EventTableProps {
    events: IssuerEventViewModel[];
    type: "pending" | "signed";
    onActionClick?: (eventId: string) => void;
}

export const EventTable: React.FC<EventTableProps> = ({ events, type, onActionClick }) => {
    const { t } = useTranslation();

    if (events.length === 0) {
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
                    <TableBody>
                        <TableRow>
                            <TableCell colSpan={5} className="px-6 py-12 text-center">
                                <div className="flex flex-col items-center justify-center space-y-2">
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        color="foreground"
                                        className="font-medium"
                                    >
                                        {t(`issuer.sign.empty.${type}.title`)}
                                    </Typography>
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        color="foreground-alt"
                                        className="text-sm"
                                    >
                                        {t(`issuer.sign.empty.${type}.description`)}
                                    </Typography>
                                </div>
                            </TableCell>
                        </TableRow>
                    </TableBody>
                </Table>
            </div>
        );
    }

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
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-white">
                                {event.eventTitle}
                            </TableCell>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-white">
                                <div className="flex flex-col gap-1">
                                    {event.ownerName ? (
                                        <span>{event.ownerName}</span>
                                    ) : (
                                        <span className="italic text-muted-foreground">
                                            (empty)
                                        </span>
                                    )}
                                    {event.ownerWalletAddress && (
                                        <span className="text-xs text-muted-foreground">
                                            {event.ownerWalletAddress}
                                        </span>
                                    )}
                                    {event.ownerGoogleEmail && (
                                        <span className="text-xs text-muted-foreground">
                                            {event.ownerGoogleEmail}
                                        </span>
                                    )}
                                    {!event.ownerWalletAddress &&
                                        !event.ownerGoogleEmail &&
                                        event.ownerName && (
                                            <span className="italic text-muted-foreground text-xs">
                                                (empty)
                                            </span>
                                        )}
                                </div>
                            </TableCell>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-white">
                                {event.certificateCount}
                            </TableCell>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm">
                                <EventStatusBadge status={type} />
                            </TableCell>
                            <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-white">
                                <EventActions
                                    type={type}
                                    eventId={event.eventId || ""}
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
