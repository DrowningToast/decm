import { CertificateSettingsPage } from "@/components/pages/HostPages/EventPages/CertificateSettingsPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function Page() {
    return (
        <ProtectedRoute>
            <CertificateSettingsPage />
        </ProtectedRoute>
    );
}
