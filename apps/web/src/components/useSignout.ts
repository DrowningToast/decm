import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { authService } from "@/services/AuthService";
import { useMutation } from "@tanstack/react-query";

export const useSignout = () => {
	const [, setAccessToken] = useLocalStorage<string | undefined>(
		LOCAL_STORAGE_KEYS.ACCESS_TOKEN,
		undefined
	);
	const [, setExpiresIn] = useLocalStorage<number | undefined>(
		LOCAL_STORAGE_KEYS.EXPIRES_IN,
		undefined
	);

	const {
		mutateAsync: signout,
		isPending,
		error,
	} = useMutation({
		mutationFn: async () => {
			await authService.signOut();

			setAccessToken(undefined);
			setExpiresIn(undefined);
		},
	});

	return { signout, isPending, error };
};
