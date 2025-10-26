import React from "react";
import { Navigate } from "@/router";
import { useAuth } from "@/context/AuthContext";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";

interface ProtectedRouteProps {
    children: React.ReactNode;
    redirectTo?: string;
    requiredRoles?: ("ADMIN" | "ISSUER" | "PARTICIPANT" | "HOST")[];
    fallback?: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
    children,
    redirectTo = "/signup",
    requiredRoles,
    fallback,
}) => {
    const { isAuthenticated, isLoading, userRole } = useAuth();
    const { t } = useTranslation();

    if (isLoading) {
        return (
            fallback || (
                <div className="flex items-center justify-center min-h-screen">
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
                        <Typography variant="text" tag="p">
                            {t("common.loading")}
                        </Typography>
                    </div>
                </div>
            )
        );
    }

    if (!isAuthenticated) {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        return <Navigate to={redirectTo as any} replace />;
    }

    if (requiredRoles && requiredRoles.length > 0 && userRole) {
        // Check if user's role is included in the required roles array
        if (!requiredRoles.includes(userRole)) {
            // User is authenticated but doesn't have the required role
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            return <Navigate to="/unauthorized" replace />;
        }
    }

    return <>{children}</>;
};
