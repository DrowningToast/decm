import type { RequirementStatus } from "@/components/ui/requirement-item";

export function toEventRegistrationConfigStatus(status?: number): RequirementStatus {
    if (!status) {
        return "not_required";
    }

    switch (status) {
        case 0:
            return "not_required";
        case 1:
            return "required";
        case 2:
            return "optional";
        default:
            return "not_required";
    }
}

export function toEventRegistrationConfigStatusNumber(status?: RequirementStatus): number {
    if (!status) {
        return 0;
    }

    switch (status) {
        case "not_required":
            return 0;
        case "required":
            return 1;
        case "optional":
            return 2;
        default:
            return 0;
    }
}
