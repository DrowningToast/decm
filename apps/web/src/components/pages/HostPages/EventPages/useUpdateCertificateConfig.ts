import { coreApiClient } from "@/lib/api/api";
import type { EventconfigUpdateEventCertificateConfigRequest } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { useCallback } from "react";

export function useUpdateCertificateConfig(eventId: string) {
    const {
        mutateAsync: _updateCertificateConfig,
        isPending: isUpdatingCertificateConfig,
        error,
    } = useMutation({
        mutationKey: ["certificate-config", eventId],
        mutationFn: (data: EventconfigUpdateEventCertificateConfigRequest) =>
            coreApiClient.v1.updateEventCertificateConfig(
                {
                    eventId,
                },
                data,
            ),
    });

    const updateCertificateConfig = useCallback(
        async (data: EventconfigUpdateEventCertificateConfigRequest) => {
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
