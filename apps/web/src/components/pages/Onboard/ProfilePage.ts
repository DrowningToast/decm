import { z } from "zod";

export const ProfileSchema = (t: (key: string) => string) =>
	z.object({
		firstName: z
			.string()
			.max(64, { message: t("validation.firstNameMax64") })
			.optional(),
		isFirstNamePublic: z.boolean().default(false).optional(),
		lastName: z
			.string()
			.max(64, { message: t("validation.lastNameMax64") })
			.optional(),
		isLastNamePublic: z.boolean().default(false).optional(),
		email: z
			.email({ message: t("validation.invalidEmail") })
			.max(64, { message: t("validation.emailMax64") })
			.optional(),
		isEmailPublic: z.boolean().default(false).optional(),
		phoneNumber: z
			.string()
			.max(10, { message: t("validation.phoneNumberMax10") })
			.optional(),
		isPhoneNumberPublic: z.boolean().default(false).optional(),
	});

export type Profile = z.infer<ReturnType<typeof ProfileSchema>>;
