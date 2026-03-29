import type { EventImportCertificateReceiverRequest } from "@decm/api";

export interface PreviewData {
    [key: string]: string | number | undefined;
}

export interface RowValidationResult {
    isValid: boolean;
    missingFields: string[];
}

// Exactly one of these must be provided per row (not both, not neither)
export const CERTIFICATE_EITHER_OR_COLUMNS = {
    email: "email",
    walletAddress: "wallet_address",
} as const;

// Optional columns
export const CERTIFICATE_OPTIONAL_COLUMNS = {
    firstName: "first_name",
    lastName: "last_name",
    academicInstitution: "academic_institution",
    certificateTitle: "certificate_title",
    certificateSubtitle: "certificate_subtitle",
} as const;

// All columns for display/table purposes
export const CERTIFICATE_ALL_COLUMNS = {
    ...CERTIFICATE_EITHER_OR_COLUMNS,
    ...CERTIFICATE_OPTIONAL_COLUMNS,
} as const;

// Normalise optional field — returns undefined for empty/null/whitespace
export const normalizeOptional = (value: unknown): string | undefined => {
    if (value == null) return undefined;
    const trimmed = String(value).trim();
    return trimmed || undefined;
};

// Validate a single row: exactly one of email / wallet_address
export const validateCertificateRow = (row: PreviewData): RowValidationResult => {
    const missingFields: string[] = [];

    const hasEmail = !!normalizeOptional(row[CERTIFICATE_EITHER_OR_COLUMNS.email]);
    const hasWallet = !!normalizeOptional(row[CERTIFICATE_EITHER_OR_COLUMNS.walletAddress]);

    if (hasEmail === hasWallet) {
        // both true → too many; both false → too few
        missingFields.push("email_or_wallet_address");
    }

    return {
        isValid: missingFields.length === 0,
        missingFields,
    };
};

// Build a certificate request item from a parsed Excel row
export const buildCertificate = (row: PreviewData): EventImportCertificateReceiverRequest => {
    const certificate: EventImportCertificateReceiverRequest = {
        academic_institution: "",
        certificate_subtitle: "",
        certificate_title: "",
        first_name: "",
        last_name: "",
    };

    const email = normalizeOptional(row.email);
    const walletAddress = normalizeOptional(row.wallet_address);
    const firstName = normalizeOptional(row.first_name);
    const lastName = normalizeOptional(row.last_name);
    const academicInstitution = normalizeOptional(row.academic_institution);
    const certificateTitle = normalizeOptional(row.certificate_title);
    const certificateSubtitle = normalizeOptional(row.certificate_subtitle);

    if (email) certificate.email = email;
    if (walletAddress) certificate.wallet_address = walletAddress;
    if (firstName) certificate.first_name = firstName;
    if (lastName) certificate.last_name = lastName;
    if (academicInstitution) certificate.academic_institution = academicInstitution;
    if (certificateTitle) certificate.certificate_title = certificateTitle;
    if (certificateSubtitle) certificate.certificate_subtitle = certificateSubtitle;

    return certificate;
};
