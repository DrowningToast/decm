import { Outlet, useSearchParams } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { lazy, Suspense, useEffect } from "react";
import { HelmetProvider, Helmet } from "react-helmet-async";
import { queryClient } from "@/lib/api/queryClient";
import "../index.css";
import { ErrorBoundary } from "react-error-boundary";
import { ErrorPage } from "@/components/pages/Error";
import { Toaster } from "@/components/ui/sonner";
import { useMediaQuery } from "@/hooks/use-media-query";
import { AppKitProvider } from "@/config/walletConnect";
import { AuthProvider } from "@/context/AuthContext";
import { WalletProvider } from "@/context/WalletContext";
import { SignPasswordModalProvider } from "@/components/providers/SignPasswordModal/SignPasswordModaProvider";
import { useTranslation } from "react-i18next";
import i18n from "@/lib/i18n/config";
import type { Language } from "@/lib/i18n/config";

import { SystemStatusWrapper } from "@/components/SystemStatusWrapper";

// Lazy load the DevTools to avoid bundle issues
const ReactQueryDevtools = lazy(() =>
    import("@tanstack/react-query-devtools").then(({ ReactQueryDevtools }) => ({
        default: ReactQueryDevtools,
    })),
);

const Layout = () => {
    const isMobile = useMediaQuery("(max-width: 768px)");
    const { t } = useTranslation();
    const [searchParams, setSearchParams] = useSearchParams();

    // Detect and handle lang query parameter
    useEffect(() => {
        const langParam = searchParams.get("lang");
        if (langParam === "en" || langParam === "th") {
            // Change language
            i18n.changeLanguage(langParam as Language);

            // Remove lang query parameter from URL
            const newSearchParams = new URLSearchParams(searchParams);
            newSearchParams.delete("lang");
            setSearchParams(newSearchParams, { replace: true });
        }
    }, [searchParams, setSearchParams]);

    return (
        <QueryClientProvider client={queryClient}>
            <ErrorBoundary fallback={<ErrorPage />}>
                <Toaster
                    richColors
                    position={isMobile ? "top-center" : "bottom-right"}
                    toastOptions={{
                        duration: 3000,
                    }}
                />
                <AppKitProvider>
                    <WalletProvider>
                        <SystemStatusWrapper>
                            <main className="font-secondary bg-background text-foreground">
                                <HelmetProvider>
                                    <Helmet>
                                        <title>{t("common.appName")}</title>
                                    </Helmet>
                                    <AuthProvider>
                                        <SignPasswordModalProvider>
                                            <Outlet />
                                        </SignPasswordModalProvider>
                                        {process.env.NODE_ENV === "development" && (
                                            <Suspense fallback={null}>
                                                <ReactQueryDevtools initialIsOpen={false} />
                                            </Suspense>
                                        )}
                                    </AuthProvider>
                                </HelmetProvider>
                            </main>
                        </SystemStatusWrapper>
                    </WalletProvider>
                </AppKitProvider>
            </ErrorBoundary>
        </QueryClientProvider>
    );
};

export default Layout;
