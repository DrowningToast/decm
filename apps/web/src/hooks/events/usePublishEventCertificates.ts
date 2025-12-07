import { useMutation, useQueryClient } from "@tanstack/react-query";
import { EventService } from "@/services/EventService/EventService";
import { QUERY_KEY } from "@/lib/queryKeys";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

export interface PublishEventCertificatesResult {
    eventId: string;
    publishedCount: number;
    isPublished: boolean;
    inboxMessagesCreated: number;
}

export function usePublishEventCertificates() {
    const queryClient = useQueryClient();
    const { t } = useTranslation();

    return useMutation({
        mutationFn: async (eventId: string): Promise<PublishEventCertificatesResult> => {
            const response = await EventService.publishEventCertificates(eventId);
            return {
                eventId: response.event_id,
                publishedCount: response.published_count,
                isPublished: response.is_published,
                inboxMessagesCreated: response.inbox_messages_created,
            };
        },
        onSuccess: (data) => {
            // Invalidate and refetch related queries
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.byId(data.eventId),
            });
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.certificate.config(data.eventId),
            });
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.certificates(data.eventId),
            });

            toast.success(
                t("events.hostDetails.certificates.publishSuccess", {
                    count: data.publishedCount,
                    messagesCount: data.inboxMessagesCreated,
                }),
            );
        },
        onError: (error: Error) => {
            toast.error(
                t("events.hostDetails.certificates.publishError", {
                    error: error.message,
                }),
            );
        },
    });
}
