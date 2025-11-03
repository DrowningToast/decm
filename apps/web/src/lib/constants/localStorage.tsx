export const LOCAL_STORAGE_KEYS = {
    ON_GOOGLE_OAUTH_SUCCESS_REDIRECT: "on_google_oauth_success_redirect",
    JWT: "jwt",
    ACCESS_TOKEN: "accessToken",
    EXPIRES_IN: "expiresIn",
    AUTH_SIGN_SIGNATURE: "authSignSignature",
} as const;

export type LocalStorageKeys = (typeof LOCAL_STORAGE_KEYS)[keyof typeof LOCAL_STORAGE_KEYS];

export type LocalStorageType = {
    [LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT]: string;
    [LOCAL_STORAGE_KEYS.JWT]: string;
    [LOCAL_STORAGE_KEYS.ACCESS_TOKEN]: string;
    [LOCAL_STORAGE_KEYS.EXPIRES_IN]: number;
    [LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE]: string;
} & Record<keyof typeof LOCAL_STORAGE_KEYS, string | number | undefined>;

export const getLocalStorageItem = <T extends LocalStorageKeys>(
    key: T,
): LocalStorageType[T] | undefined => {
    const value = localStorage.getItem(key);
    if (value === null) {
        return undefined;
    }
    if (key === LOCAL_STORAGE_KEYS.EXPIRES_IN) {
        try {
            return JSON.parse(value) as LocalStorageType[T];
        } catch {
            const num = Number(value);
            return (Number.isNaN(num) ? undefined : num) as LocalStorageType[T];
        }
    }
    return value as LocalStorageType[T];
};

export const setLocalStorageItem = (
    key: LocalStorageKeys,
    value: LocalStorageType[keyof LocalStorageType],
) => {
    if (typeof value !== "string") {
        value = JSON.stringify(value);
    }
    localStorage.setItem(key, value);
};

export const removeLocalStorageItem = (key: LocalStorageKeys) => {
    localStorage.removeItem(key);
};
