import { USECASE_IDS, type UseCaseId } from "./usecase";

export const TOAST_USECASE_VIEWMODEL = {
    [USECASE_IDS.OAUTH_SIGNUP]: {
        LOADING: "flow.oauth_google.create_account_pending",
        LOADING_PROFILE: "flow.oauth_google.create_profile_pending",
        SUCCESS: "flow.oauth_google.create_account_success",
        SUCCESS_PROFILE: "flow.oauth_google.create_profile_success",
        ERROR: "flow.oauth_google.create_account_error",
        ERROR_DUPLICATE: "flow.oauth_google.create_account_error_duplicate",
        ERROR_DUPLICATE_PROFILE: "flow.oauth_google.create_profile_error_duplicate",
        ERROR_PROFILE: "flow.oauth_google.create_profile_error",
        UNAUTHENTICATED_RESPONSE: "flow.oauth_google.unauthenticated_response",
        EXPIRED_TOKEN: "flow.oauth_google.expired_token",
        INTERNAL_ERROR_RESPONSE: "flow.oauth_google.internal_error_response",
    },
    [USECASE_IDS.CHECK_ONBOARD_STATUS]: {
        UNAUTHENTICATED_RESPONSE: "flow.check_onboard_status.unauthenticated_response",
        EXPIRED_TOKEN: "flow.check_onboard_status.expired_token",
        INTERNAL_ERROR_RESPONSE: "flow.check_onboard_status.internal_error_response",
    },
    [USECASE_IDS.OAUTH_GOOGLE_CREATE_ACCOUNT]: {
        SUCCESS: "flow.oauth_google.create_account_success",
        ERROR: "flow.oauth_google.create_account_error",
        ERROR_DUPLICATE: "flow.oauth_google.create_account_error_duplicate",
    },
    [USECASE_IDS.GENERIC]: {},
    [USECASE_IDS.OAUTH_GOOGLE_CREATE_PROFILE]: {
        SUCCESS: "flow.oauth_google.create_profile_success",
        ERROR: "flow.oauth_google.create_profile_error",
        ERROR_DUPLICATE: "flow.oauth_google.create_profile_error_duplicate",
    },
} satisfies Record<UseCaseId, Record<string, string>>;
