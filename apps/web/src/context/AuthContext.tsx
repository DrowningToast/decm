import React, { createContext, useContext, useEffect, useState } from "react";
import type { ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import type { EntityProfile } from "@decm/api";
import { queryKeys } from "@/lib/queryKeys";
interface AuthContextType {
    user: EntityProfile | null;
    isLoading: boolean;
    isAuthenticated: boolean;
    refetch: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [isInitialized, setIsInitialized] = useState(false);

    const {
        data: user,
        isLoading,
        refetch,
    } = useQuery<EntityProfile | null>({
        queryKey: queryKeys.user.profile,
        queryFn: async () => {
            try {
                const response = await coreApiClient.v1.getMyProfile();
                return response;
            } catch (error) {
                console.error(error);
                return null;
            }
        },
        retry: false,
        refetchOnWindowFocus: false,
        staleTime: 5 * 60 * 1000,
    });

    useEffect(() => {
        if (!isLoading) {
            setIsInitialized(true);
        }
    }, [isLoading]);

    const isAuthenticated = !!user;

    const contextValue: AuthContextType = {
        user: user || null,
        isLoading: !isInitialized || isLoading,
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
