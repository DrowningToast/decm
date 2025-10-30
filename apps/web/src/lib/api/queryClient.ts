import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            // Time in milliseconds that unused/inactive cache data remains in memory
            gcTime: 1000 * 60 * 60 * 24, // 24 hours
            // Time in milliseconds that the data is considered fresh
            staleTime: 1000 * 60 * 5, // 5 minutes
            retry: false,
            // Refetch on window focus in development, but not in production
            refetchOnWindowFocus: process.env.NODE_ENV === "development",
        },
        mutations: {
            // Retry failed mutations
            retry: false,
        },
    },
});
