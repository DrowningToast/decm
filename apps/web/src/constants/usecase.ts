export const USECASE_IDS = {
    GENERIC: "usecase.generic",
    OAUTH_SIGNUP: "usecase.oauth_signup",
    OAUTH_GOOGLE_CREATE_ACCOUNT: "usecase.oauth_google_create_account",
    OAUTH_GOOGLE_CREATE_PROFILE: "usecase.oauth_google_create_profile",
    CHECK_ONBOARD_STATUS: "usecase.check_onboard_status",
    SIGN_IN: "usecase.sign_in",
} as const;

export type UseCaseId = (typeof USECASE_IDS)[keyof typeof USECASE_IDS];
