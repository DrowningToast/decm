export const LOCAL_STORAGE_KEYS = {
    ON_GOOGLE_OAUTH_SUCCESS_REDIRECT: "on_google_oauth_success_redirect",
    JWT: "jwt",
    ACCESS_TOKEN: "accessToken",
    EXPIRES_IN: "expiresIn",
    AUTH_SIGN_SIGNATURE: "authSignSignature",
} as const;

export type LocalStorageType = {
    [LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT]: string;
    [LOCAL_STORAGE_KEYS.JWT]: string;
    [LOCAL_STORAGE_KEYS.ACCESS_TOKEN]: string;
    [LOCAL_STORAGE_KEYS.EXPIRES_IN]: number;
    [LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE]: string;
} & Record<keyof typeof LOCAL_STORAGE_KEYS, string | number | undefined>;

export type LocalStorageKeys = (typeof LOCAL_STORAGE_KEYS)[keyof typeof LOCAL_STORAGE_KEYS];

export const getLocalStorageItem = (
    key: LocalStorageKeys,
): LocalStorageType[keyof LocalStorageType] => {
    return localStorage.getItem(key) as LocalStorageType[keyof LocalStorageType];
};

export const setLocalStorageItem = (
    key: LocalStorageKeys,
    value: LocalStorageType[keyof LocalStorageType],
) => {
    localStorage.setItem(key, JSON.stringify(value));
};

export const removeLocalStorageItem = (key: LocalStorageKeys) => {
    localStorage.removeItem(key);
};
