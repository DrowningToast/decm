import React, { createContext, useContext } from "react";
import type { ReactNode } from "react";
import type { EntityProfile } from "@decm/api";
import { AxiosError } from "axios";
import { useMyProfile } from "@/hooks/useMyProfile";
interface AuthContextType {
    user: EntityProfile | null;
    isFetching: boolean;
    isAuthenticated: boolean;
    refetch: () => Promise<unknown>;
    error?: AxiosError;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const { data: user, error, isFetching, refetch } = useMyProfile();

    const contextValue: AuthContextType = {
        user: user || null,
        isFetching,
        isAuthenticated: !!user,
        refetch,
        error: error instanceof AxiosError ? error : undefined,
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
