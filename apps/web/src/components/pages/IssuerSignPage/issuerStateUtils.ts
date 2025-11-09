// Issuer signing states
export const ISSUER_STATES = {
    PENDING: 0,
    SIGNED: 1,
} as const;

export type IssuerState = (typeof ISSUER_STATES)[keyof typeof ISSUER_STATES];

// Helper functions to determine issuer state
export const getIssuerStateText = (isSigned: number, t?: (key: string) => string): string => {
    const translationKey =
        isSigned === ISSUER_STATES.SIGNED
            ? "issuer.sign.status.signed"
            : isSigned === ISSUER_STATES.PENDING
              ? "issuer.sign.status.waiting"
              : "common.unknown";

    return t ? t(translationKey) : translationKey;
};

export const getIssuerStateColor = (isSigned: number): string => {
    switch (isSigned) {
        case ISSUER_STATES.SIGNED:
            return "text-green-500";
        case ISSUER_STATES.PENDING:
            return "text-orange-500";
        default:
            return "text-gray-500";
    }
};

export const getIssuerStateBgColor = (isSigned: number): string => {
    switch (isSigned) {
        case ISSUER_STATES.SIGNED:
            return "bg-green-100";
        case ISSUER_STATES.PENDING:
            return "bg-orange-100";
        default:
            return "bg-gray-100";
    }
};

export const isIssuerSigned = (isSigned: number): boolean => {
    return isSigned === ISSUER_STATES.SIGNED;
};

export const isIssuerPending = (isSigned: number): boolean => {
    return isSigned === ISSUER_STATES.PENDING;
};

// Calculate signing progress percentage
export const calculateSigningProgress = (issuers: Array<{ is_signed: number }>): number => {
    if (!issuers || issuers.length === 0) return 0;

    const signedCount = issuers.filter((issuer) => isIssuerSigned(issuer.is_signed)).length;
    return Math.round((signedCount / issuers.length) * 100);
};
