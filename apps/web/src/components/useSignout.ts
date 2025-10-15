import { handleUniversalError } from "@/common/Err";
import { authService } from "@/services/AuthService";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

export const useSignout = () => {
	const { t } = useTranslation();
	const {
		mutateAsync: signout,
		isPending,
		error,
	} = useMutation({
		mutationFn: async () => {
			try {
				return await authService.signOut();
			} catch (error) {
				if (error instanceof Error) {
					handleUniversalError(t, error);
				}
			}
		},
	});

	return { signout, isPending, error };
};
