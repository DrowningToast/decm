import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { CertificateSettingsPage } from "@/components/pages/HostPages/EventPages/CertificateSettingsPage";
import { useEventCertificateConfig } from "@/components/pages/HostPages/EventPages/useEventCertificateConfig";
import { useVerifiedIssuers } from "@/hooks/events/useVerifiedIssuers";
import { useParams } from "react-router-dom";

export default function Page() {
    const { eventId } = useParams<{ eventId: string }>();

    const { data: eventCertificateConfig, isLoading: isLoadingEventCertificateConfig } =
        useEventCertificateConfig(eventId!);

    const { verifiedIssuers, isLoadingVerifiedIssuers } = useVerifiedIssuers();

    if (isLoadingEventCertificateConfig || isLoadingVerifiedIssuers) {
        return <div>Loading...</div>;
    }

    return (
        <ProtectedRoute>
            <CertificateSettingsPage
                eventId={eventId!}
                eventCertificateConfig={eventCertificateConfig}
                verifiedIssuers={verifiedIssuers ?? []}
            />
        </ProtectedRoute>
    );
}
