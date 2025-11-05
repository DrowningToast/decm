import React, { createContext, useContext } from "react";
import type { ReactNode } from "react";
import type { EntityProfile } from "@decm/api";
import { useMyProfile } from "@/hooks/useMyProfile";
interface AuthContextType {
    user: EntityProfile | null;
    isPending: boolean;
    isAuthenticated: boolean;
    refetch: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const { data: user, isPending, refetch } = useMyProfile();

    const isAuthenticated = !!user;

    const contextValue: AuthContextType = {
        user: user || null,
        isPending,
        isAuthenticated,
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
