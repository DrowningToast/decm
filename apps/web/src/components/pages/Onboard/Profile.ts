import { z } from "zod";

export const ProfileSchema = z.object({
	firstName: z.string().min(1),
	lastName: z.string().min(1),
	email: z.email(),
	phone: z.string().min(1).max(10),
	address: z.string().min(1),
	academicInstitution: z.string().min(1),
	academicEmail: z.email(),
});

export type Profile = z.infer<typeof ProfileSchema>;
