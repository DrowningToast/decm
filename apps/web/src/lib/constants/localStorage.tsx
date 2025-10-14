export const LOCAL_STORAGE_KEYS = {
    ON_GOOGLE_OAUTH_SUCCESS_REDIRECT: "on_google_oauth_success_redirect",
    JWT: "jwt",
} as const

export type LocalStorageKeys = typeof LOCAL_STORAGE_KEYS[keyof typeof LOCAL_STORAGE_KEYS];

export const getLocalStorageItem = (key: LocalStorageKeys) => {
    return localStorage.getItem(key);
}

export const setLocalStorageItem = (key: LocalStorageKeys, value: string) => {
    localStorage.setItem(key, value);
}

export const removeLocalStorageItem = (key: LocalStorageKeys) => {
    localStorage.removeItem(key);
}