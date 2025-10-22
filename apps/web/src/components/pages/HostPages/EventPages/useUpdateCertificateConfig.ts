import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import type { UpdateEventCertificateConfigPayload } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { useCallback } from "react";

export function useUpdateCertificateConfig(eventId: string) {
    const {
        mutateAsync: _updateCertificateConfig,
        isPending: isUpdatingCertificateConfig,
        error,
    } = useMutation({
        mutationKey: ["certificate-config", eventId],
        mutationFn: (data: UpdateEventCertificateConfigPayload) =>
            coreApiClient.v1.updateEventCertificateConfig(
                {
                    eventId,
                },
                data,
            ),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["event"] });
        },
    });

    const updateCertificateConfig = useCallback(
        async (data: UpdateEventCertificateConfigPayload) => {
            return await _updateCertificateConfig(data);
        },
        [_updateCertificateConfig],
    );

    return {
        updateCertificateConfig,
        isUpdatingCertificateConfig,
        error,
    };
}
