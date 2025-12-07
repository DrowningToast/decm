import type { FieldStatus } from "../ConfirmForm";

export const FieldRequirementBadge: React.FC<{
    status: FieldStatus;
    t: (key: string) => string;
}> = ({ status, t }) => {
    if (status === "locked") {
        return (
            <span className="ml-2 inline-flex items-center rounded-md bg-secondary/50 px-2 py-0.5 text-xs font-medium text-muted-foreground ring-1 ring-inset ring-border">
                {t("common.locked")}
            </span>
        );
    }
    if (status === "required") {
        return (
            <span className="ml-2 inline-flex items-center rounded-md bg-destructive/10 px-2 py-0.5 text-xs font-medium text-destructive ring-1 ring-inset ring-destructive/20">
                {t("common.required")}
            </span>
        );
    }
    if (status === "optional") {
        return (
            <span className="ml-2 inline-flex items-center rounded-md bg-secondary/50 px-2 py-0.5 text-xs font-medium text-muted-foreground ring-1 ring-inset ring-border">
                {t("common.optional")}
            </span>
        );
    }
    return null;
};
