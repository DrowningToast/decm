import { type SignOutParams } from "@/services/AuthService/AuthService";
import { authService } from "@/services/services";
import { useMutation } from "@tanstack/react-query";
import { useCallback } from "react";

export const useSignout = () => {
    const {
        mutateAsync: _signout,
        isPending,
        error,
    } = useMutation({
        mutationFn: async (params?: SignOutParams | undefined) => {
            await authService.signOut(params || {});
        },
    });

    const signout = useCallback(
        async (params?: SignOutParams | undefined) => {
            await _signout(params);
        },
        [_signout],
    );

    return { signout, isPending, error };
};
