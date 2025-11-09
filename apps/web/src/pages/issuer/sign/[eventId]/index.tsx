import IssuerSignPage from "@/components/pages/IssuerSignPage/IssuerSignPage";
import { useParams } from "@/router";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function Page() {
    const { eventId } = useParams("/issuer/sign/:eventId");

    return (
        <ProtectedRoute>
            <IssuerSignPage eventId={eventId!} />
        </ProtectedRoute>
    );
}
