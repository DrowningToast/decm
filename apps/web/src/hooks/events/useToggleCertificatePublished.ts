import { useMutation, useQueryClient } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { toast } from "sonner";
import { QUERY_KEY } from "@/lib/queryKeys";

interface ToggleCertificatePublishedParams {
    eventId: string;
    isPublished: boolean;
}

export const useToggleCertificatePublished = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ eventId, isPublished }: ToggleCertificatePublishedParams) => {
            return await coreApiClient.v1.toggleCertificatePublished(
                { eventId },
                { is_published: isPublished },
            );
        },
        onSuccess: (_data, variables) => {
            // Invalidate and refetch certificate config
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.certificate.config(variables.eventId),
            });

            if (variables.isPublished) {
                toast.success("Certificate configuration published successfully");
            }
        },
        onError: (error: unknown) => {
            const errorMessage =
                error && typeof error === "object" && "response" in error
                    ? (error.response as { data?: { message?: string } })?.data?.message
                    : undefined;
            toast.error(errorMessage || "Failed to update certificate published status");
        },
    });
};
