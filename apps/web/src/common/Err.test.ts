import { describe, it, expect, vi } from "vitest";
import { AxiosError } from "axios";
import { Err, ToastFromError, ToastFromAxiosError, handleAxiosError } from "./Err";
import { toast } from "sonner";

// Mock sonner toast
vi.mock("sonner", () => ({
    toast: {
        error: vi.fn(),
        success: vi.fn(),
        warning: vi.fn(),
        info: vi.fn(),
    },
}));

describe("Err Class", () => {
    it("should create an Err instance with message", () => {
        const err = new Err("Test error");

        expect(err).toBeInstanceOf(Error);
        expect(err.message).toBe("Test error");
        expect(err.toastType).toBe("error");
    });

    it("should create an Err instance with all properties", () => {
        const err = new Err("Test error", "Test Title", "Test Description", "warning");

        expect(err.message).toBe("Test error");
        expect(err.title).toBe("Test Title");
        expect(err.description).toBe("Test Description");
        expect(err.toastType).toBe("warning");
    });

    it("should set default toastType to error", () => {
        const err = new Err("Test error");

        expect(err.toastType).toBe("error");
    });

    it("should handle different toast types", () => {
        const errorErr = new Err("msg", "title", "desc", "error");
        const successErr = new Err("msg", "title", "desc", "success");
        const warningErr = new Err("msg", "title", "desc", "warning");
        const infoErr = new Err("msg", "title", "desc", "info");

        expect(errorErr.toastType).toBe("error");
        expect(successErr.toastType).toBe("success");
        expect(warningErr.toastType).toBe("warning");
        expect(infoErr.toastType).toBe("info");
    });

    it("should handle useCaseId", () => {
        const mockUseCase = "FORM_SUBMISSION" as string;
        const err = new Err("Test", "Title", "Desc", "error", mockUseCase);

        expect(err.useCaseId).toBe(mockUseCase);
    });
});

describe("ToastFromError", () => {
    it("should toast Err instance with its own properties", () => {
        const t = vi.fn((key: string) => key);
        const err = new Err("message", "errors.custom", "errors.customDesc", "success");

        ToastFromError(t, err);

        expect(toast.success).toHaveBeenCalledWith("errors.custom", expect.any(Object));
    });

    it("should use preset for error type", () => {
        const t = vi.fn((key: string) => key);
        const regularError = new Error("Regular error");

        ToastFromError(t, regularError, "INVALID_INPUT");

        expect(toast.error).toHaveBeenCalledWith("errors.invalidInput", expect.any(Object));
    });

    it("should fallback to generic error", () => {
        const t = vi.fn((key: string) => key);
        const unknownError = new Error("Unknown error");

        ToastFromError(t, unknownError);

        expect(toast.error).toHaveBeenCalledWith("errors.generic", expect.any(Object));
    });

    it("should handle AxiosError", () => {
        const t = vi.fn((key: string) => key);
        const axiosError = new AxiosError("Network error", "NETWORK_ERROR", undefined, undefined, {
            status: 400,
            data: null,
            statusText: "Bad Request",
            headers: {},
            config: {} as never,
        });

        ToastFromError(t, axiosError);

        // Should call ToastFromAxiosError internally
        expect(toast.error).toHaveBeenCalled();
    });

    it("should use error type if Err instance doesn't match presets", () => {
        const t = vi.fn((key: string) => key);
        const err = new Err("message");

        ToastFromError(t, err);

        expect(toast.error).toHaveBeenCalled();
    });
});

describe("ToastFromAxiosError", () => {
    it("should toast 400 error as INVALID_INPUT", () => {
        const t = vi.fn((key: string) => key);
        const error = new AxiosError("Bad request", "400", undefined, undefined, {
            status: 400,
            data: null,
            statusText: "Bad Request",
            headers: {},
            config: {} as never,
        });

        ToastFromAxiosError(t, error);

        expect(toast.error).toHaveBeenCalledWith("errors.invalidInput", expect.any(Object));
    });

    it("should toast 401 error as UNAUTHORIZED", () => {
        const t = vi.fn((key: string) => key);
        const error = new AxiosError("Unauthorized", "401", undefined, undefined, {
            status: 401,
            data: null,
            statusText: "Unauthorized",
            headers: {},
            config: {} as never,
        });

        ToastFromAxiosError(t, error);

        expect(toast.error).toHaveBeenCalledWith("errors.unauthorized", expect.any(Object));
    });

    it("should toast 403 error as FORBIDDEN", () => {
        const t = vi.fn((key: string) => key);
        const error = new AxiosError("Forbidden", "403", undefined, undefined, {
            status: 403,
            data: null,
            statusText: "Forbidden",
            headers: {},
            config: {} as never,
        });

        ToastFromAxiosError(t, error);

        expect(toast.error).toHaveBeenCalledWith("errors.forbidden", expect.any(Object));
    });

    it("should toast 404 error as NOT_FOUND", () => {
        const t = vi.fn((key: string) => key);
        const error = new AxiosError("Not found", "404", undefined, undefined, {
            status: 404,
            data: null,
            statusText: "Not Found",
            headers: {},
            config: {} as never,
        });

        ToastFromAxiosError(t, error);

        expect(toast.error).toHaveBeenCalledWith("errors.notFound", expect.any(Object));
    });

    it("should toast 409 error as DUPLICATE_ENTRY", () => {
        const t = vi.fn((key: string) => key);
        const error = new AxiosError("Conflict", "409", undefined, undefined, {
            status: 409,
            data: null,
            statusText: "Conflict",
            headers: {},
            config: {} as never,
        });

        ToastFromAxiosError(t, error);

        expect(toast.error).toHaveBeenCalledWith("errors.conflict", expect.any(Object));
    });

    it("should toast 500 error as INTERNAL_SERVER_ERROR", () => {
        const t = vi.fn((key: string) => key);
        const error = new AxiosError("Server error", "500", undefined, undefined, {
            status: 500,
            data: null,
            statusText: "Internal Server Error",
            headers: {},
            config: {} as never,
        });

        ToastFromAxiosError(t, error);

        expect(toast.error).toHaveBeenCalledWith("errors.serverError", expect.any(Object));
    });

    it("should toast unknown status as INTERNAL_SERVER_ERROR", () => {
        const t = vi.fn((key: string) => key);
        const error = new AxiosError("Unknown", "999", undefined, undefined, {
            status: 999,
            data: null,
            statusText: "Unknown",
            headers: {},
            config: {} as never,
        });

        ToastFromAxiosError(t, error);

        expect(toast.error).toHaveBeenCalledWith("errors.serverError", expect.any(Object));
    });

    it("should use custom presets when provided", () => {
        const t = vi.fn((key: string) => key);
        const customPresets = {
            INVALID_INPUT: {
                title: "custom.invalidInput",
                description: "custom.desc",
                toastType: "error" as const,
            },
        };

        const error = new AxiosError("Bad request", "400", undefined, undefined, {
            status: 400,
            data: null,
            statusText: "Bad Request",
            headers: {},
            config: {} as never,
        });

        ToastFromAxiosError(t, error, customPresets);

        expect(toast.error).toHaveBeenCalledWith("custom.invalidInput", expect.any(Object));
    });

    it("should handle error without response status", () => {
        const t = vi.fn((key: string) => key);
        const error = new AxiosError("Network error");

        ToastFromAxiosError(t, error);

        expect(toast.error).toHaveBeenCalledWith("errors.serverError", expect.any(Object));
    });
});

describe("handleAxiosError", () => {
    it("should call onInvalidInput hook for 400 error", () => {
        const t = vi.fn((key: string) => key);
        const onInvalidInput = vi.fn();
        const error = new AxiosError("Bad request", "400", undefined, undefined, {
            status: 400,
            data: null,
            statusText: "Bad Request",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { onInvalidInput });

        expect(onInvalidInput).toHaveBeenCalledWith(error);
    });

    it("should call onUnauthorized hook for 401 error", () => {
        const t = vi.fn((key: string) => key);
        const onUnauthorized = vi.fn();
        const error = new AxiosError("Unauthorized", "401", undefined, undefined, {
            status: 401,
            data: null,
            statusText: "Unauthorized",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { onUnauthorized });

        expect(onUnauthorized).toHaveBeenCalledWith(error);
    });

    it("should call onForbidden hook for 403 error", () => {
        const t = vi.fn((key: string) => key);
        const onForbidden = vi.fn();
        const error = new AxiosError("Forbidden", "403", undefined, undefined, {
            status: 403,
            data: null,
            statusText: "Forbidden",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { onForbidden });

        expect(onForbidden).toHaveBeenCalledWith(error);
    });

    it("should call onNotFound hook for 404 error", () => {
        const t = vi.fn((key: string) => key);
        const onNotFound = vi.fn();
        const error = new AxiosError("Not found", "404", undefined, undefined, {
            status: 404,
            data: null,
            statusText: "Not Found",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { onNotFound });

        expect(onNotFound).toHaveBeenCalledWith(error);
    });

    it("should call onDuplicateEntry hook for 409 error", () => {
        const t = vi.fn((key: string) => key);
        const onDuplicateEntry = vi.fn();
        const error = new AxiosError("Conflict", "409", undefined, undefined, {
            status: 409,
            data: null,
            statusText: "Conflict",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { onDuplicateEntry });

        expect(onDuplicateEntry).toHaveBeenCalledWith(error);
    });

    it("should call onInternalServerError hook for 500 error", () => {
        const t = vi.fn((key: string) => key);
        const onInternalServerError = vi.fn();
        const error = new AxiosError("Server error", "500", undefined, undefined, {
            status: 500,
            data: null,
            statusText: "Internal Server Error",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { onInternalServerError });

        expect(onInternalServerError).toHaveBeenCalledWith(error);
    });

    it("should use custom error from hook", () => {
        const t = vi.fn((key: string) => key);
        const customErr = new Err("Custom", "custom.title", "custom.desc", "warning");
        const error = new AxiosError("Bad request", "400", undefined, undefined, {
            status: 400,
            data: null,
            statusText: "Bad Request",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { invalidInputErr: customErr });

        expect(toast.warning).toHaveBeenCalledWith("custom.title", expect.any(Object));
    });

    it("should handle multiple hooks at once", () => {
        const t = vi.fn((key: string) => key);
        const onInvalidInput = vi.fn();
        const invalidInputErr = new Err("Custom", "custom.title", "custom.desc", "warning");
        const error = new AxiosError("Bad request", "400", undefined, undefined, {
            status: 400,
            data: null,
            statusText: "Bad Request",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, {
            onInvalidInput,
            invalidInputErr,
        });

        expect(onInvalidInput).toHaveBeenCalledWith(error);
        expect(toast.warning).toHaveBeenCalled();
    });

    it("should fall back to toast if no custom error is provided", () => {
        const t = vi.fn((key: string) => key);
        const onInvalidInput = vi.fn();
        const error = new AxiosError("Bad request", "400", undefined, undefined, {
            status: 400,
            data: null,
            statusText: "Bad Request",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { onInvalidInput });

        expect(onInvalidInput).toHaveBeenCalled();
        expect(toast.error).toHaveBeenCalled();
    });

    it("should use useCaseId from custom error if provided", () => {
        const t = vi.fn((key: string) => key);
        const customErr = new Err("Custom", "custom.title", "custom.desc", "error", "CUSTOM_CASE");
        const error = new AxiosError("Bad request", "400", undefined, undefined, {
            status: 400,
            data: null,
            statusText: "Bad Request",
            headers: {},
            config: {} as never,
        });

        handleAxiosError(t, error, { invalidInputErr: customErr });

        expect(toast.error).toHaveBeenCalledWith("custom.title", {
            description: expect.any(String),
            id: "CUSTOM_CASE",
        });
    });
});
