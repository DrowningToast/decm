import { Outlet } from "react-router-dom";
import { QueryClientProvider } from '@tanstack/react-query';
import { lazy, Suspense } from 'react';
import { HelmetProvider } from 'react-helmet-async';
import { queryClient } from '@/lib/api/queryClient';
import "../index.css"

// Lazy load the DevTools to avoid bundle issues
const ReactQueryDevtools = lazy(() =>
    import('@tanstack/react-query-devtools').then(({ ReactQueryDevtools }) => ({
        default: ReactQueryDevtools,
    }))
);

const Layout = () => {
    return (
        <main className="font-secondary bg-background text-foreground">
            <HelmetProvider>
                <QueryClientProvider client={queryClient}>
                    <Outlet />
                    {process.env.NODE_ENV === 'development' && (
                        <Suspense fallback={null}>
                            <ReactQueryDevtools initialIsOpen={false} />
                        </Suspense>
                    )}
                </QueryClientProvider>
            </HelmetProvider>
        </main>
    )
}

export default Layout;