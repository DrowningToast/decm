export const USECASE_IDS = {
	GENERIC: "generic",
	OAUTH_SIGNUP: "oauth_signup",
	CHECK_ONBOARD_STATUS: "check_onboard_status",
} as const;

export type UseCaseId = (typeof USECASE_IDS)[keyof typeof USECASE_IDS];
