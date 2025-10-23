import { USECASE_IDS, type UseCaseId } from "@/constants/usecase";
import type { TFunction } from "i18next";
import { toast } from "sonner";
import { AxiosError } from "axios";

const ToastFn = (type?: ToastType) => {
    switch (type) {
        case "error":
            return toast.error;
        case "success":
            return toast.success;
        case "warning":
            return toast.warning;
        case "info":
            return toast.info;
        default:
            return toast.error;
    }
};

export type ErrType =
    | "INVALID_INPUT"
    | "NETWORK_ERROR"
    | "UNAUTHORIZED"
    | "FORBIDDEN"
    | "NOT_FOUND"
    | "DUPLICATE_ENTRY"
    | "INTERNAL_CLIENT_ERROR"
    | "INTERNAL_SERVER_ERROR";

export type ErrPresets = Record<ErrType, Partial<Err>>;

const ERROR_PRESETS: ErrPresets = {
    INVALID_INPUT: {
        title: "errors.invalidInput",
        description: "errors.invalidInputDescription",
        toastType: "error",
        useCaseId: USECASE_IDS.GENERIC,
    },
    NETWORK_ERROR: {
        title: "errors.network",
        description: "errors.networkDescription",
        toastType: "error",
        useCaseId: USECASE_IDS.GENERIC,
    },
    UNAUTHORIZED: {
        title: "errors.unauthorized",
        description: "errors.unauthorizedDescription",
        toastType: "error",
        useCaseId: USECASE_IDS.GENERIC,
    },
    FORBIDDEN: {
        title: "errors.forbidden",
        description: "errors.forbiddenDescription",
        toastType: "error",
        useCaseId: USECASE_IDS.GENERIC,
    },
    NOT_FOUND: {
        title: "errors.notFound",
        description: "errors.notFoundDescription",
        toastType: "error",
        useCaseId: USECASE_IDS.GENERIC,
    },
    DUPLICATE_ENTRY: {
        title: "errors.conflict",
        description: "errors.duplicateEntryDescription",
        toastType: "error",
        useCaseId: USECASE_IDS.GENERIC,
    },
    INTERNAL_CLIENT_ERROR: {
        title: "errors.internalClientError",
        description: "errors.internalClientErrorDescription",
        toastType: "error",
        useCaseId: USECASE_IDS.GENERIC,
    },
    INTERNAL_SERVER_ERROR: {
        title: "errors.serverError",
        description: "errors.internalServerErrorDescription",
        toastType: "error",
        useCaseId: USECASE_IDS.GENERIC,
    },
};

type ToastType = "error" | "success" | "warning" | "info";

export class Err extends Error {
    title?: string;
    description?: string;
    toastType: ToastType;
    useCaseId?: UseCaseId;

    constructor(
        message: string,
        title?: string,
        description?: string,
        toastType?: ToastType,
        useCaseId?: UseCaseId,
    ) {
        super(message);
        this.title = title;
        this.description = description;
        this.toastType = toastType || "error";
        this.useCaseId = useCaseId;
    }
}

export const ToastFromError = (t: TFunction, err: Error, type?: ErrType) => {
    // If already an Err instance, toast the given data
    if (err instanceof Err) {
        ToastFn(err.toastType)(t(err.title || "errors.generic"), {
            description: t(err.description || "errors.genericDescription"),
        });
        return;
    }

    // If not an Err instance, but provided a type, use the preset
    const preset = type ? ERROR_PRESETS[type] : undefined;
    if (preset) {
        ToastFn(preset.toastType)(t(preset.title || "errors.generic"), {
            description: t(preset.description || "errors.genericDescription"),
        });
        return;
    }

    // Try to determine the error type
    if (err instanceof AxiosError) {
        return ToastFromAxiosError(t, err);
    } else {
        return ToastFn("error")(t("errors.generic"), {
            description: t("errors.genericDescription"),
        });
    }

    return err;
};

export const ToastFromAxiosError = (
    t: TFunction,
    err: AxiosError,
    presets: ErrPresets = ERROR_PRESETS,
) => {
    let errorPresetType: ErrType | undefined;
    switch (err.response?.status) {
        case 400:
            errorPresetType = "INVALID_INPUT";
            break;
        case 401:
            errorPresetType = "UNAUTHORIZED";
            break;
        case 403:
            errorPresetType = "FORBIDDEN";
            break;
        case 404:
            errorPresetType = "NOT_FOUND";
            break;
        case 409:
            errorPresetType = "DUPLICATE_ENTRY";
            break;
        case 500:
            errorPresetType = "INTERNAL_SERVER_ERROR";
            break;
        default:
            errorPresetType = "INTERNAL_SERVER_ERROR";
            break;
    }
    if (errorPresetType) {
        ToastFn(presets[errorPresetType].toastType)(
            t(presets[errorPresetType].title || "errors.generic"),
            {
                description: t(presets[errorPresetType].description || "errors.genericDescription"),
            },
        );
    } else {
        ToastFn("error")(t("errors.generic"), {
            description: t("errors.genericDescription"),
        });
    }
};

export const handleUniversalError = (
    t: TFunction,
    err: Error,
    hooks?: Partial<HandlerAxiosErrorHooks>,
    useCaseId?: UseCaseId,
) => {
    ToastFromError(t, err);

    if (err instanceof AxiosError) {
        handleAxiosError(t, err, hooks, useCaseId);
    }
};

interface HandlerAxiosErrorHooks {
    onInvalidInput?: (err: AxiosError) => void;
    invalidInputErr?: Err;
    onUnauthorized?: (err: AxiosError) => void;
    unauthorizedErr?: Err;
    onForbidden?: (err: AxiosError) => void;
    forbiddenErr?: Err;
    onNotFound?: (err: AxiosError) => void;
    notFoundErr?: Err;
    onDuplicateEntry?: (err: AxiosError) => void;
    duplicateEntryErr?: Err;
    onInternalServerError?: (err: AxiosError) => void;
    internalServerErrorErr?: Err;
}

export const handleAxiosError = (
    t: TFunction,
    err: AxiosError,
    hooks?: Partial<HandlerAxiosErrorHooks>,
    useCaseId?: UseCaseId,
) => {
    let customErr: Err | undefined;

    switch (err.response?.status) {
        case 400:
            customErr = hooks?.invalidInputErr;
            hooks?.onInvalidInput?.(err);
            break;
        case 401:
            customErr = hooks?.unauthorizedErr;
            hooks?.onUnauthorized?.(err);
            break;
        case 403:
            customErr = hooks?.forbiddenErr;
            hooks?.onForbidden?.(err);
            break;
        case 404:
            customErr = hooks?.notFoundErr;
            hooks?.onNotFound?.(err);
            break;
        case 409:
            customErr = hooks?.duplicateEntryErr;
            hooks?.onDuplicateEntry?.(err);
            break;
        case 500:
        default:
            customErr = hooks?.internalServerErrorErr;
            hooks?.onInternalServerError?.(err);
            break;
    }

    if (customErr) {
        ToastFn(customErr.toastType || "error")(t(customErr.title || "errors.generic"), {
            description: t(customErr.description || "errors.genericDescription"),
            id: customErr.useCaseId ?? useCaseId ?? USECASE_IDS.GENERIC,
        });
    } else {
        ToastFromAxiosError(t, err);
    }
};
