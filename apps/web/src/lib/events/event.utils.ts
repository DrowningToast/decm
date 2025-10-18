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
