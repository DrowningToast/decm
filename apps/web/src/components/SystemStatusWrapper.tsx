import { useTranslation } from "react-i18next";
import { Loader2, AlertTriangle, Clock, WifiOff } from "lucide-react";
import { format } from "date-fns";
import { Button } from "./ui/button";
import { useSystemStatus } from "@/hooks/useSystemStatus";

interface SystemStatusWrapperProps {
    children: React.ReactNode;
}

export const SystemStatusWrapper = ({ children }: SystemStatusWrapperProps) => {
    const { t } = useTranslation();

    const { data, isLoading, isError, refetch } = useSystemStatus();

    if (isLoading) {
        return (
            <div className="flex h-screen w-screen items-center justify-center bg-background">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        );
    }

    if (isError) {
        return (
            <div className="flex h-screen w-screen flex-col items-center justify-center bg-background p-4 text-center">
                <div className="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-destructive/10">
                    <WifiOff className="h-10 w-10 text-destructive" />
                </div>
                <h1 className="mb-2 text-2xl font-bold text-foreground">
                    {t("common.networkError", "Connection Error")}
                </h1>
                <p className="mb-6 max-w-md text-muted-foreground">
                    {t(
                        "errors.networkDescription",
                        "We're having trouble connecting to our servers. Please check your internet connection and try again.",
                    )}
                </p>
                <Button onClick={() => refetch()} variant="primary">
                    {t("error.tryAgain", "Try Again")}
                </Button>
            </div>
        );
    }

    if (!data || !data.latestSchedule) {
        // Fallback to allowing access if we can't determine status but no explicit error
        return <>{children}</>;
    }

    const currentStatus = data.latestSchedule;

    if (currentStatus.status === "maintenance") {
        return (
            <div className="flex h-screen w-screen flex-col items-center justify-center bg-background p-4 text-center">
                <div className="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-yellow-100 dark:bg-yellow-900/20">
                    <AlertTriangle className="h-10 w-10 text-yellow-600 dark:text-yellow-500" />
                </div>

                <h1 className="mb-2 text-2xl font-bold text-foreground">
                    {currentStatus.isPlanned
                        ? t("systemStatus.maintenance.planned.title", "System Maintenance")
                        : t("systemStatus.maintenance.unplanned.title", "System Maintenance")}
                </h1>

                <p className="mb-6 max-w-md text-muted-foreground">
                    {currentStatus.isPlanned
                        ? t(
                            "systemStatus.maintenance.planned.description",
                            "We are currently performing scheduled maintenance to improve our services.",
                        )
                        : t(
                            "systemStatus.maintenance.unplanned.description",
                            "We are currently experiencing some technical difficulties and are working to resolve them.",
                        )}
                </p>

                {currentStatus.plannedEndTime && (
                    <div className="flex items-center gap-2 rounded-lg bg-secondary/50 px-4 py-2 text-sm text-foreground">
                        <Clock className="h-4 w-4" />
                        <span>
                            {t("systemStatus.maintenance.expectedEnd", "Expected completion:")}{" "}
                            <span className="font-semibold">
                                {format(currentStatus.plannedEndTime, "PP p")}
                            </span>
                        </span>
                    </div>
                )}
            </div>
        );
    }

    return <>{children}</>;
};
