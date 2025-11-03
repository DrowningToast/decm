import React, { createContext, useContext, useEffect } from "react";
import type { ReactNode } from "react";
import type { EntityProfile } from "@decm/api";
import { AxiosError } from "axios";
import { useTranslation } from "react-i18next";
import { useSignout } from "@/components/useSignout";
import { useMyProfile } from "@/hooks/useMyProfile";
import { handleAxiosError } from "@/common/Err";
interface AuthContextType {
    user: EntityProfile | null;
    isPending: boolean;
    isAuthenticated: boolean;
    refetch: () => Promise<unknown>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const { t } = useTranslation();
    const { signout } = useSignout();
    const {
        data: user,
        error,
        isPending,
        refetch,
    } = useMyProfile();

    useEffect(() => {
        if (!error) {
            return;
        }

        console.error(error);
        if (error instanceof AxiosError) {
            handleAxiosError(t, error);
            void signout();
        }
    }, [error, signout, t]);

    const contextValue: AuthContextType = {
        user: user || null,
        isPending,
        isAuthenticated: !!user,
        refetch,
    };

    return <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>;
};

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};
