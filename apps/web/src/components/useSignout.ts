import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { authService } from "@/services/AuthService";
import { useDisconnect } from "@reown/appkit/react";
import { useMutation } from "@tanstack/react-query";
import { useWalletClient } from "wagmi";

export const useSignout = () => {
	const [, setAccessToken] = useLocalStorage<string | undefined>(
		LOCAL_STORAGE_KEYS.ACCESS_TOKEN,
		undefined
	);
	const [, setExpiresIn] = useLocalStorage<number | undefined>(
		LOCAL_STORAGE_KEYS.EXPIRES_IN,
		undefined
	);
	const [, setAuthSignSignature] = useLocalStorage<string | undefined>(
		LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
		undefined
	);
	const { data: walletClient } = useWalletClient();
	const { disconnect } = useDisconnect();

	const {
		mutateAsync: signout,
		isPending,
		error,
	} = useMutation({
		mutationFn: async () => {
			await authService.signOut();

			setAccessToken(undefined);
			setExpiresIn(undefined);
			setAuthSignSignature(undefined);
			if (walletClient) {
				await disconnect();
			}
		},
	});

	return { signout, isPending, error };
};
