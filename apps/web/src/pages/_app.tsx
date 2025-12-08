import { Outlet } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { lazy, Suspense } from "react";
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

// Lazy load the DevTools to avoid bundle issues
const ReactQueryDevtools = lazy(() =>
    import("@tanstack/react-query-devtools").then(({ ReactQueryDevtools }) => ({
        default: ReactQueryDevtools,
    })),
);

const Layout = () => {
    const isMobile = useMediaQuery("(max-width: 768px)");
    const { t } = useTranslation();

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
                    </WalletProvider>
                </AppKitProvider>
            </ErrorBoundary>
        </QueryClientProvider>
    );
};

export default Layout;
