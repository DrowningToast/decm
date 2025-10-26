import { LogoutButton } from "@/components/LogoutButton";

const AppPage = () => {
    return (
        <ProtectedRoute>
            <h1>app page</h1>
            {/* PH */}
            <LogoutButton />
        </ProtectedRoute>
    );
};

export default AppPage;
