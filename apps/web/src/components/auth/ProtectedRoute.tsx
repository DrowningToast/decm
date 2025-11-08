import React, { useEffect } from "react";
import { type Path } from "@/router";
import { useAuth } from "@/context/AuthContext";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";
import { Navigate } from "react-router-dom";
import { toast } from "sonner";
import { TOAST_USECASE_VIEWMODEL } from "@/constants/toast";
import { USECASE_IDS } from "@/constants/usecase";
import { useCheckRoles } from "@/hooks/useCheckRoles";

interface ProtectedRouteProps {
    children: React.ReactNode;
    redirectTo?: Path;
    unauthorizedRedirectTo?: Path;
    /**
     * Require user to be authenticated (checked by AuthContext)
     * @default true
     */
    requireAuthenticated?: boolean;
    /**
     * Require user to be a verified host/organizer
     */
    requireHost?: boolean;
    /**
     * Require user to be a verified issuer
     */
    requireIssuer?: boolean;
    fallback?: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
    children,
    redirectTo = "/signup",
    unauthorizedRedirectTo = "/unauthorized",
    requireAuthenticated = true,
    requireHost = false,
    requireIssuer = false,
    fallback,
}) => {
    const { isAuthenticated, isFetching } = useAuth();
    const { t } = useTranslation();

    // Determine if role checking should be enabled
    const shouldCheckRoles = isAuthenticated && (requireHost || requireIssuer);

    // Use the useCheckRoles hook
    const {
        hasRequiredRoles,
        isLoading: isCheckingRoles,
        isError: roleCheckError,
    } = useCheckRoles({
        requireHost,
        requireIssuer,
        enabled: shouldCheckRoles,
    });

    // Show toast error when role check fails
    useEffect(() => {
        if (shouldCheckRoles && !isCheckingRoles && !hasRequiredRoles && !roleCheckError) {
            toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.GENERIC].UNAUTHORIZED_RESPONSE));
        }
    }, [shouldCheckRoles, isCheckingRoles, hasRequiredRoles, roleCheckError, t]);

    // Show toast error when role check encounters an error
    useEffect(() => {
        if (roleCheckError) {
            toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.GENERIC].UNAUTHORIZED_RESPONSE));
        }
    }, [roleCheckError, t]);

    // Show loading state during auth check
    if (isFetching) {
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

    // Check authentication requirement
    if (requireAuthenticated && !isAuthenticated) {
        toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.GENERIC].UNAUTHENTICATED_RESPONSE));
        return <Navigate to={redirectTo} replace />;
    }

    // Show loading state during role check
    if (shouldCheckRoles && isCheckingRoles) {
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

    // Check role requirements
    if (shouldCheckRoles && (!hasRequiredRoles || roleCheckError)) {
        return <Navigate to={unauthorizedRedirectTo} replace />;
    }

    return <>{children}</>;
};
