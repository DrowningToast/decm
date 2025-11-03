import React from "react";
import { type Path } from "@/router";
import { useAuth } from "@/context/AuthContext";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";
import { Navigate } from "react-router-dom";
import { toast } from "sonner";
import { TOAST_USECASE_VIEWMODEL } from "@/constants/toast";
import { USECASE_IDS } from "@/constants/usecase";

interface ProtectedRouteProps {
    children: React.ReactNode;
    redirectTo?: Path;
    // requiredRoles?: ("ADMIN" | "ISSUER" | "PARTICIPANT" | "HOST")[];
    fallback?: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
    children,
    redirectTo = "/signup",
    // requiredRoles,
    fallback,
}) => {
    const { isAuthenticated, isPending } = useAuth();
    const { t } = useTranslation();

    if (isPending) {
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
        toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.GENERIC].UNAUTHENTICATED_RESPONSE));
        return <Navigate to={redirectTo} replace />;
    }

    // TODO: Implement role based protection
    // if (requiredRoles && requiredRoles.length > 0 && user?.role) {
    //     // Check if user's role is included in the required roles array
    //     if (!requiredRoles.includes(userRole)) {
    //         // User is authenticated but doesn't have the required role

    //         return <Navigate to="/unauthorized" replace />;
    //     }
    // }

    return <>{children}</>;
};
