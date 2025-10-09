import { z } from "zod";

export const ProfileSchema = (t: (key: string) => string) =>
	z.object({
		firstName: z
			.string()
			.max(64, { message: t("validation.firstNameMax64") })
			.optional(),
		lastName: z
			.string()
			.max(64, { message: t("validation.lastNameMax64") })
			.optional(),
		contactEmail: z
			.email({ message: t("validation.invalidEmail") })
			.max(64, { message: t("validation.emailMax64") })
			.optional(),
		phoneNumber: z
			.string()
			.max(10, { message: t("validation.phoneNumberMax10") })
			.optional(),
	});

export type Profile = z.infer<typeof ProfileSchema>;
