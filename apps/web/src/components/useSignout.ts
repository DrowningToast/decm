import { authService } from "@/services/AuthService";
import { useMutation } from "@tanstack/react-query";

export const useSignout = () => {
    const {
        mutateAsync: signout,
        isPending,
        error,
    } = useMutation({
        mutationFn: async () => {
            await authService.signOut();
        },
    });

    return { signout, isPending, error };
};
