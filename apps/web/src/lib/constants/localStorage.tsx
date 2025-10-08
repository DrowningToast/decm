export const LOCAL_STORAGE_KEYS = {
    ON_GOOGLE_OAUTH_SUCCESS_REDIRECT: "on_google_oauth_success_redirect",
}

export const getLocalStorageItem = (key: keyof typeof LOCAL_STORAGE_KEYS) => {
    return localStorage.getItem(LOCAL_STORAGE_KEYS[key]);
}

export const setLocalStorageItem = (key: keyof typeof LOCAL_STORAGE_KEYS, value: string) => {
    localStorage.setItem(LOCAL_STORAGE_KEYS[key], value);
}

export const removeLocalStorageItem = (key: keyof typeof LOCAL_STORAGE_KEYS) => {
    localStorage.removeItem(LOCAL_STORAGE_KEYS[key]);
}