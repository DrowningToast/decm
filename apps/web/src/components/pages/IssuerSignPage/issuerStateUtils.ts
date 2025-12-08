// Issuer signing states
export const ISSUER_STATES = {
    PENDING: 0,
    SIGNED: 1,
} as const;

export type IssuerState = (typeof ISSUER_STATES)[keyof typeof ISSUER_STATES];

// Helper functions to determine issuer state
export const getIssuerStateText = (isSigned: boolean, t?: (key: string) => string): string => {
    const translationKey = isSigned
        ? "issuer.sign.status.signed"
        : !isSigned
          ? "issuer.sign.status.waiting"
          : "common.unknown";

    return t ? t(translationKey) : translationKey;
};

export const getIssuerStateColor = (isSigned: boolean): string => {
    switch (isSigned) {
        case true:
            return "text-green-500";
        case false:
            return "text-orange-500";
        default:
            return "text-gray-500";
    }
};

export const getIssuerStateBgColor = (isSigned: boolean): string => {
    switch (isSigned) {
        case true:
            return "bg-green-100";
        case false:
            return "bg-orange-100";
        default:
            return "bg-gray-100";
    }
};

export const isIssuerSigned = (isSigned: boolean): boolean => {
    return isSigned;
};

// Calculate signing progress percentage
export const calculateSigningProgress = (issuers: Array<{ is_signed: boolean }>): number => {
    if (!issuers || issuers.length === 0) return 0;

    const signedCount = issuers.filter((issuer) => isIssuerSigned(issuer.is_signed)).length;
    return Math.round((signedCount / issuers.length) * 100);
};
