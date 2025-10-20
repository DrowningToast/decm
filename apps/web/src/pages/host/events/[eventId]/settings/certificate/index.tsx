import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { CertificateSettingsPage } from "@/components/pages/HostPages/EventPages/CertificateSettingsPage";

export default function Page() {
    return (
        <ProtectedRoute>
            <CertificateSettingsPage />
        </ProtectedRoute>
    );
}
