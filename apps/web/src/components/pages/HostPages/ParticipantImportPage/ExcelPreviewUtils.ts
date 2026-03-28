import type { EventRegistrationParticipantRequestItem } from "@decm/api";

export interface PreviewData {
    [key: string]: string | number | undefined;
}

export interface RowValidationResult {
    isValid: boolean;
    missingFields: string[];
}

// Always-required columns — must be present and non-empty on every row
export const ALWAYS_REQUIRED_COLUMNS = {
    firstName: "first_name",
    lastName: "last_name",
} as const;

// Exactly one of these must be provided per row (not both, not neither)
export const EITHER_OR_COLUMNS = {
    email: "email",
    walletAddress: "wallet_address",
} as const;

// Optional columns
export const OPTIONAL_COLUMNS = {
    phoneNumber: "phone_number",
    academicInstitution: "academic_institution",
} as const;

// All columns for display/table purposes
export const ALL_COLUMNS = {
    ...ALWAYS_REQUIRED_COLUMNS,
    ...EITHER_OR_COLUMNS,
    ...OPTIONAL_COLUMNS,
} as const;

// Normalise optional field — returns undefined for empty/null/whitespace
export const normalizeOptional = (value: unknown): string | undefined => {
    if (value == null) return undefined;
    const trimmed = String(value).trim();
    return trimmed || undefined;
};

// Validate a single row: first_name + last_name required; exactly one of email / wallet_address
export const validateRow = (row: PreviewData): RowValidationResult => {
    const missingFields: string[] = [];

    if (!normalizeOptional(row[ALWAYS_REQUIRED_COLUMNS.firstName])) {
        missingFields.push(ALWAYS_REQUIRED_COLUMNS.firstName);
    }
    if (!normalizeOptional(row[ALWAYS_REQUIRED_COLUMNS.lastName])) {
        missingFields.push(ALWAYS_REQUIRED_COLUMNS.lastName);
    }

    const hasEmail = !!normalizeOptional(row[EITHER_OR_COLUMNS.email]);
    const hasWallet = !!normalizeOptional(row[EITHER_OR_COLUMNS.walletAddress]);

    if (hasEmail === hasWallet) {
        // both true → too many; both false → too few
        missingFields.push("email_or_wallet_address");
    }

    return {
        isValid: missingFields.length === 0,
        missingFields,
    };
};

// Build a participant request item from a parsed Excel row
export const buildParticipant = (row: PreviewData): EventRegistrationParticipantRequestItem => {
    const participant: EventRegistrationParticipantRequestItem = {
        first_name: String(row.first_name || ""),
        last_name: String(row.last_name || ""),
    };

    const email = normalizeOptional(row.email);
    const walletAddress = normalizeOptional(row.wallet_address);
    const phoneNumber = normalizeOptional(row.phone_number);
    const academicInstitution = normalizeOptional(row.academic_institution);

    if (email) participant.email = email;
    if (walletAddress) participant.wallet_address = walletAddress;
    if (phoneNumber) participant.phone_number = phoneNumber;
    if (academicInstitution) participant.academic_institution = academicInstitution;

    return participant;
};
